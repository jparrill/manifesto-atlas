// Package artifact implements the atlas.artifact/v1 envelope and graph
// validator established by docs/design/artifact-contract.md (TQ-0008).
//
// It provides deterministic canonical JSON serialization, SHA-256 content
// identity, manifest structure validation, bounded offline graph traversal,
// payload integrity, and precise invalidation. Verification is fail-closed and
// never retrieves files, follows URLs, invokes a model, or approves domain
// content. Attestation authority is TQ-0020 and is deliberately absent here.
package artifact

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

// Contract is the artifact envelope version implemented by this package.
const Contract = "atlas.artifact/v1"

// Limits. None may be disabled by values inside an untrusted manifest.
const (
	// maxCanonicalInput bounds a JSON payload fed to Canonicalize.
	maxCanonicalInput = 16 << 20
	// maxManifestInput bounds a single manifest.
	maxManifestInput = 1 << 20
	// maxGraphInput bounds the total raw manifest input per graph.
	maxGraphInput = 64 << 20
	// maxNodes bounds the node count per graph.
	maxNodes = 10000
	// maxEdgesPerNode bounds the input edges per node.
	maxEdgesPerNode = 1024
	// maxDepth is the maximum JSON nesting depth (root is depth 1).
	maxDepth = 64
)

// v1 integer bounds: [-(2^53-1), 2^53-1].
const (
	maxInt = 9007199254740991
	minInt = -9007199254740991
)

// Error codes. They are stable and carry no private data.
const (
	CodeMalformed      = "malformed-json"
	CodeTrailing       = "trailing-value"
	CodeDupMember      = "duplicate-member"
	CodeFractional     = "fractional-number"
	CodeNegZero        = "negative-zero"
	CodeOutOfRange     = "integer-out-of-range"
	CodeTooDeep        = "too-deep"
	CodeTooLarge       = "too-large"
	CodeInvalidUTF8    = "invalid-utf8"
	CodeMissingField   = "missing-field"
	CodeUnsupportedVer = "unsupported-version"
	CodeUnknownKind    = "unknown-kind"
	CodeUnknownInput   = "unknown-input"
	CodeMissingInput   = "missing-input"
	CodeIntegrity      = "integrity"
	CodeIDMismatch     = "id-mismatch"
	CodeDupID          = "duplicate-id"
	CodeMissingDep     = "missing-dependency"
	CodeCycle          = "cycle"
	CodeInputMismatch  = "input-mismatch"
	CodeBadResult      = "bad-result"
)

// Error is a stable, fail-closed artifact error.
type Error struct {
	Code    string
	Subject string
	Err     error
}

func (e *Error) Error() string {
	s := "artifact: " + e.Code
	if e.Subject != "" {
		s += " (" + e.Subject + ")"
	}
	if e.Err != nil {
		s += ": " + e.Err.Error()
	}
	return s
}

// Unwrap exposes the underlying error when present.
func (e *Error) Unwrap() error { return e.Err }

func errf(code, subject string) error { return &Error{Code: code, Subject: subject} }

// HashBytes returns the atlas H(bytes) identity: "sha256:" plus 64 lowercase
// hex digits of the SHA-256 digest.
func HashBytes(b []byte) string {
	sum := sha256.Sum256(b)
	return "sha256:" + hex.EncodeToString(sum[:])
}

// ManifestID returns the content identity of a manifest: H(C(manifest)).
func ManifestID(raw []byte) (string, error) {
	canon, err := Canonicalize(raw)
	if err != nil {
		return "", err
	}
	return HashBytes(canon), nil
}

// ValidHash reports whether h is a full H identity: "sha256:" followed by 64
// lowercase hexadecimal digits.
func ValidHash(h string) bool { return validHash(h) }

// CanonicalValue returns the canonical bytes for a Go value restricted to the
// strict v1 profile (objects, arrays, strings, booleans, null, and integer
// json.Number). It is a convenience for building canonical statements without
// hand-serializing JSON.
func CanonicalValue(v any) ([]byte, error) { return encodeCanonical(v) }

// Canonicalize validates that data is exactly one strict JSON value under the
// atlas.artifact/v1 serialization rules and returns its canonical bytes.
//
// The rules reject a BOM, trailing data, duplicate object members, fractional
// or exponent numbers, negative zero, integers outside the v1 range, nesting
// deeper than maxDepth, invalid UTF-8, and unpaired surrogates. Keys are sorted
// by UTF-8 byte order; array order and string code points are preserved. The
// output uses Go encoding/json escaping with HTML escaping disabled and has no
// trailing newline.
func Canonicalize(data []byte) ([]byte, error) {
	if len(data) > maxCanonicalInput {
		return nil, errf(CodeTooLarge, "input")
	}
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		return nil, errf(CodeMalformed, "bom")
	}
	p := &parser{b: data}
	p.skipWS()
	v, err := p.parseValue()
	if err != nil {
		return nil, err
	}
	p.skipWS()
	if p.i != len(data) {
		return nil, errf(CodeTrailing, fmt.Sprintf("%d", p.i))
	}
	return encodeCanonical(v)
}

// parser is a strict, bounded, integer-preserving JSON scanner. It validates
// while building a Go tree so that canonicalization is a pure re-encoding.
type parser struct {
	b     []byte
	i     int
	depth int
}

func (p *parser) skipWS() {
	for p.i < len(p.b) {
		switch p.b[p.i] {
		case ' ', '\t', '\n', '\r':
			p.i++
		default:
			return
		}
	}
}

func (p *parser) parseValue() (any, error) {
	p.skipWS()
	if p.i >= len(p.b) {
		return nil, errf(CodeMalformed, "eof")
	}
	switch c := p.b[p.i]; c {
	case '{':
		return p.parseObject()
	case '[':
		return p.parseArray()
	case '"':
		return p.parseString()
	case 't':
		return p.parseLit("true", true)
	case 'f':
		return p.parseLit("false", false)
	case 'n':
		return p.parseLit("null", nil)
	default:
		return p.parseNumber()
	}
}

func (p *parser) parseLit(word string, val any) (any, error) {
	if p.i+len(word) <= len(p.b) && string(p.b[p.i:p.i+len(word)]) == word {
		p.i += len(word)
		return val, nil
	}
	return nil, errf(CodeMalformed, string(p.b[p.i]))
}

func (p *parser) parseObject() (any, error) {
	if p.depth+1 > maxDepth {
		return nil, errf(CodeTooDeep, fmt.Sprintf("%d", p.depth+1))
	}
	p.i++ // consume '{'
	p.depth++
	defer func() { p.depth-- }()
	obj := make(map[string]any)
	p.skipWS()
	if p.i < len(p.b) && p.b[p.i] == '}' {
		p.i++
		return obj, nil
	}
	for {
		p.skipWS()
		if p.i >= len(p.b) || p.b[p.i] != '"' {
			return nil, errf(CodeMalformed, "object-key")
		}
		key, err := p.parseString()
		if err != nil {
			return nil, err
		}
		if _, dup := obj[key]; dup {
			return nil, errf(CodeDupMember, key)
		}
		p.skipWS()
		if p.i >= len(p.b) || p.b[p.i] != ':' {
			return nil, errf(CodeMalformed, "object-colon")
		}
		p.i++ // consume ':'
		val, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		obj[key] = val
		p.skipWS()
		if p.i >= len(p.b) {
			return nil, errf(CodeMalformed, "object-end")
		}
		switch p.b[p.i] {
		case ',':
			p.i++
		case '}':
			p.i++
			return obj, nil
		default:
			return nil, errf(CodeMalformed, "object-sep")
		}
	}
}

func (p *parser) parseArray() (any, error) {
	if p.depth+1 > maxDepth {
		return nil, errf(CodeTooDeep, fmt.Sprintf("%d", p.depth+1))
	}
	p.i++ // consume '['
	p.depth++
	defer func() { p.depth-- }()
	arr := make([]any, 0, 4)
	p.skipWS()
	if p.i < len(p.b) && p.b[p.i] == ']' {
		p.i++
		return arr, nil
	}
	for {
		v, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		arr = append(arr, v)
		p.skipWS()
		if p.i >= len(p.b) {
			return nil, errf(CodeMalformed, "array-end")
		}
		switch p.b[p.i] {
		case ',':
			p.i++
		case ']':
			p.i++
			return arr, nil
		default:
			return nil, errf(CodeMalformed, "array-sep")
		}
	}
}

// parseString reads a JSON string into a Go string, rejecting unescaped control
// characters, invalid UTF-8, and unpaired or invalid surrogate escapes.
func (p *parser) parseString() (string, error) {
	p.i++ // consume opening quote
	var sb strings.Builder
	for {
		if p.i >= len(p.b) {
			return "", errf(CodeMalformed, "string-eof")
		}
		c := p.b[p.i]
		switch {
		case c == '"':
			p.i++
			return sb.String(), nil
		case c == '\\':
			p.i++
			if p.i >= len(p.b) {
				return "", errf(CodeMalformed, "string-escape-eof")
			}
			e := p.b[p.i]
			switch e {
			case '"', '\\', '/':
				sb.WriteByte(e)
				p.i++
			case 'b':
				sb.WriteByte('\b')
				p.i++
			case 'f':
				sb.WriteByte('\f')
				p.i++
			case 'n':
				sb.WriteByte('\n')
				p.i++
			case 'r':
				sb.WriteByte('\r')
				p.i++
			case 't':
				sb.WriteByte('\t')
				p.i++
			case 'u':
				r, err := p.parseUnicode()
				if err != nil {
					return "", err
				}
				if utf16.IsSurrogate(r) {
					// p.i is at the second escape's backslash; advance to its 'u'
					// so parseUnicode can consume the \uXXXX for the low surrogate.
					if p.i+1 < len(p.b) && p.b[p.i] == '\\' && p.b[p.i+1] == 'u' {
						p.i++
						r2, err := p.parseUnicode()
						if err != nil {
							return "", err
						}
						combined := utf16.DecodeRune(r, r2)
						if combined == utf8.RuneError {
							return "", errf(CodeInvalidUTF8, "surrogate-pair")
						}
						sb.WriteRune(combined)
					} else {
						return "", errf(CodeInvalidUTF8, "unpaired-surrogate")
					}
				} else {
					sb.WriteRune(r)
				}
			default:
				return "", errf(CodeMalformed, "string-escape")
			}
		case c < 0x20:
			return "", errf(CodeMalformed, "string-control")
		default:
			r, size := utf8.DecodeRune(p.b[p.i:])
			if r == utf8.RuneError && size == 1 {
				return "", errf(CodeInvalidUTF8, "")
			}
			sb.WriteRune(r)
			p.i += size
		}
	}
}

func (p *parser) parseUnicode() (rune, error) {
	p.i++ // consume 'u'
	if p.i+4 > len(p.b) {
		return 0, errf(CodeMalformed, "unicode-eof")
	}
	var r rune
	for j := 0; j < 4; j++ {
		h := p.b[p.i+j]
		var d rune
		switch {
		case h >= '0' && h <= '9':
			d = rune(h - '0')
		case h >= 'a' && h <= 'f':
			d = rune(h-'a') + 10
		case h >= 'A' && h <= 'F':
			d = rune(h-'A') + 10
		default:
			return 0, errf(CodeMalformed, "unicode-hex")
		}
		r = r*16 + d
	}
	p.i += 4
	return r, nil
}

// parseNumber reads a strict v1 integer, rejecting fractional/exponent tokens
// and negative zero, and returning an integer-preserving json.Number.
func (p *parser) parseNumber() (any, error) {
	start := p.i
	neg := false
	if p.i < len(p.b) && p.b[p.i] == '-' {
		neg = true
		p.i++
	}
	if p.i >= len(p.b) || p.b[p.i] < '0' || p.b[p.i] > '9' {
		return nil, errf(CodeMalformed, "number-int")
	}
	if p.b[p.i] == '0' {
		p.i++
		if p.i < len(p.b) && p.b[p.i] >= '0' && p.b[p.i] <= '9' {
			return nil, errf(CodeMalformed, "leading-zero")
		}
	} else {
		for p.i < len(p.b) && p.b[p.i] >= '0' && p.b[p.i] <= '9' {
			p.i++
		}
	}
	if p.i < len(p.b) {
		switch p.b[p.i] {
		case '.', 'e', 'E':
			return nil, errf(CodeFractional, "")
		}
	}
	tok := string(p.b[start:p.i])
	if neg {
		allZero := true
		for k := 1; k < len(tok); k++ {
			if tok[k] != '0' {
				allZero = false
				break
			}
		}
		if allZero {
			return nil, errf(CodeNegZero, "")
		}
	}
	n, err := strconv.ParseInt(tok, 10, 64)
	if err != nil || n > maxInt || n < minInt {
		return nil, errf(CodeOutOfRange, tok)
	}
	return json.Number(tok), nil
}

// encodeCanonical encodes a validated Go tree to canonical bytes: objects with
// keys sorted by UTF-8 byte order (Go's encoding/json map ordering), arrays and
// strings preserved, Go HTML-escaping disabled, and no trailing newline.
func encodeCanonical(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, &Error{Code: CodeMalformed, Err: err}
	}
	b := buf.Bytes()
	if len(b) > 0 && b[len(b)-1] == '\n' {
		b = b[:len(b)-1]
	}
	return b, nil
}

// Payload describes the immutable bytes a manifest binds.
type Payload struct {
	MediaType string
	SHA256    string
	Size      int64
}

// Manifest is a validated atlas.artifact/v1 envelope.
type Manifest struct {
	Contract string
	Kind     string
	Payload  Payload
	Inputs   map[string]string
}

type manifestWire struct {
	Contract *string            `json:"contract"`
	Kind     *string            `json:"kind"`
	Payload  *payloadWire       `json:"payload"`
	Inputs   *map[string]string `json:"inputs"`
}

type payloadWire struct {
	MediaType *string `json:"media_type"`
	SHA256    *string `json:"sha256"`
	Size      *int64  `json:"size"`
}

// ParseManifest validates the strict envelope: exactly the required fields, the
// v1 contract, a known kind, a well-formed payload, valid slot names, and the
// exact required input slot set for the kind (no missing or undeclared slots).
func ParseManifest(raw []byte) (*Manifest, error) {
	if len(raw) > maxManifestInput {
		return nil, errf(CodeTooLarge, "manifest")
	}
	canon, err := Canonicalize(raw)
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(bytes.NewReader(canon))
	dec.DisallowUnknownFields()
	var w manifestWire
	if err := dec.Decode(&w); err != nil {
		return nil, &Error{Code: CodeMalformed, Err: err}
	}
	if w.Contract == nil {
		return nil, errf(CodeMissingField, "contract")
	}
	if *w.Contract != Contract {
		return nil, errf(CodeUnsupportedVer, "contract")
	}
	if w.Kind == nil {
		return nil, errf(CodeMissingField, "kind")
	}
	if !knownKinds[*w.Kind] {
		return nil, errf(CodeUnknownKind, *w.Kind)
	}
	if w.Payload == nil {
		return nil, errf(CodeMissingField, "payload")
	}
	if w.Payload.MediaType == nil || *w.Payload.MediaType == "" {
		return nil, errf(CodeMissingField, "payload.media_type")
	}
	if w.Payload.SHA256 == nil || !validHash(*w.Payload.SHA256) {
		return nil, errf(CodeIntegrity, "payload.sha256")
	}
	if w.Payload.Size == nil || *w.Payload.Size < 0 {
		return nil, errf(CodeMissingField, "payload.size")
	}
	if w.Inputs == nil {
		return nil, errf(CodeMissingField, "inputs")
	}
	req := requiredFor(*w.Kind)
	for slot := range *w.Inputs {
		if _, ok := req[slot]; !ok {
			return nil, errf(CodeUnknownInput, slot)
		}
	}
	for slot := range req {
		if _, ok := (*w.Inputs)[slot]; !ok {
			return nil, errf(CodeMissingInput, slot)
		}
	}
	return &Manifest{
		Contract: *w.Contract,
		Kind:     *w.Kind,
		Payload:  Payload{MediaType: *w.Payload.MediaType, SHA256: *w.Payload.SHA256, Size: *w.Payload.Size},
		Inputs:   *w.Inputs,
	}, nil
}

// ValidatePayload checks that the stored payload bytes match the manifest's
// declared size and hash before they may be parsed as content.
func ValidatePayload(m *Manifest, payload []byte) error {
	if int64(len(payload)) != m.Payload.Size {
		return errf(CodeIntegrity, "size")
	}
	if HashBytes(payload) != m.Payload.SHA256 {
		return errf(CodeIntegrity, "sha256")
	}
	return nil
}

func validHash(h string) bool {
	if len(h) != 7+64 || h[:7] != "sha256:" {
		return false
	}
	for i := 7; i < len(h); i++ {
		c := h[i]
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			return false
		}
	}
	return true
}

// requiredFor returns the exact required input slots (slot -> kind) for a kind.
// Leaf kinds have none; every other kind requires schema, tool, and
// configuration plus its declared additional slots.
func requiredFor(kind string) map[string]string {
	m := map[string]string{}
	if leafKinds[kind] {
		return m
	}
	for _, ri := range baseInputs {
		m[ri.slot] = ri.kind
	}
	for _, ri := range kindRequired[kind] {
		m[ri.slot] = ri.kind
	}
	return m
}

// baseInputs are the required schema/tool/configuration edges for every
// non-leaf kind.
var baseInputs = []reqInput{
	{"schema", "schema"},
	{"tool", "tool"},
	{"configuration", "configuration"},
}

// reqInput is a required input slot and the kind it must reference.
type reqInput struct {
	slot string
	kind string
}

// leafKinds are the kinds with no inputs.
var leafKinds = map[string]bool{
	"source": true, "schema": true, "taxonomy": true, "prompt": true,
	"model": true, "tool": true, "configuration": true, "templates": true,
	"trust-policy": true,
}

// kindRequired lists the additional required input slots per non-leaf kind, on
// top of baseInputs.
var kindRequired = map[string][]reqInput{
	"provenance": {
		{"submitted-source", "source"}, {"canonical-source", "source"},
		{"trust-policy", "trust-policy"}, {"evidence", "configuration"},
		{"prompt", "prompt"}, {"model", "model"},
	},
	"source-record":      {{"source", "source"}, {"provenance", "provenance"}},
	"extraction":         {{"source-record", "source-record"}, {"extraction-model", "model"}},
	"secondary-evidence": {{"source-record", "source-record"}},
	"page-images":        {{"source-record", "source-record"}},
	"extraction-audit": {
		{"extraction", "extraction"}, {"secondary-evidence", "secondary-evidence"},
		{"page-images", "page-images"}, {"prompt", "prompt"}, {"model", "model"},
	},
	"extraction-approval": {
		{"audit", "extraction-audit"}, {"extraction", "extraction"},
		{"formula", "configuration"},
	},
	"programme": {
		{"extraction-approval", "extraction-approval"}, {"extraction", "extraction"},
		{"taxonomy", "taxonomy"}, {"jurisdiction", "configuration"},
		{"prompt", "prompt"}, {"model", "model"},
	},
	"claim-report": {
		{"programme", "programme"}, {"extraction", "extraction"},
		{"source-record", "source-record"}, {"page-images", "page-images"},
		{"prompt", "prompt"}, {"model", "model"},
	},
	"coverage-report": {
		{"programme", "programme"}, {"extraction", "extraction"},
		{"source-record", "source-record"}, {"page-images", "page-images"},
		{"prompt", "prompt"}, {"model", "model"},
	},
	"content-approval": {
		{"programme", "programme"}, {"claim-report", "claim-report"},
		{"coverage-report", "coverage-report"}, {"extraction-approval", "extraction-approval"},
	},
	"markdown": {{"programme", "programme"}, {"content-approval", "content-approval"}},
	"html": {
		{"programme", "programme"}, {"content-approval", "content-approval"},
		{"templates", "templates"},
	},
	"site-deterministic-report": {
		{"programme", "programme"}, {"markdown", "markdown"},
		{"html", "html"}, {"qa", "configuration"},
	},
	"site-visual-report": {
		{"programme", "programme"}, {"html", "html"},
		{"deterministic-report", "site-deterministic-report"},
		{"prompt", "prompt"}, {"model", "model"}, {"qa", "configuration"},
	},
	"site-approval": {
		{"deterministic-report", "site-deterministic-report"},
		{"visual-report", "site-visual-report"},
	},
	"publication": {
		{"source-record", "source-record"}, {"extraction-approval", "extraction-approval"},
		{"content-approval", "content-approval"}, {"markdown", "markdown"},
		{"html", "html"}, {"site-approval", "site-approval"},
		{"trust-policy", "trust-policy"},
	},
}

// knownKinds is every v1 kind (leaf and non-leaf).
var knownKinds = func() map[string]bool {
	m := map[string]bool{}
	for k := range leafKinds {
		m[k] = true
	}
	for k := range kindRequired {
		m[k] = true
	}
	return m
}()

// reportKinds are the kinds whose JSON payload carries a required result.
var reportKinds = map[string]bool{
	"provenance": true, "extraction-audit": true, "extraction-approval": true,
	"claim-report": true, "coverage-report": true, "content-approval": true,
	"site-deterministic-report": true, "site-visual-report": true,
	"site-approval": true, "publication": true,
}

// Result values for a report payload.
const (
	ResultPass         = "pass"
	ResultFail         = "fail"
	ResultInsufficient = "insufficient"
	ResultError        = "error"
)

var resultValues = map[string]bool{
	ResultPass: true, ResultFail: true, ResultInsufficient: true, ResultError: true,
}

// ReportResult extracts and validates the result of a report-kind payload.
// The payload must be a strict JSON object whose result is in the v1 enum and
// whose findings is an array. Stage-specific decision/evidence fields are
// permitted but never interpreted here.
func ReportResult(kind string, payloadJSON []byte) (string, error) {
	if !reportKinds[kind] {
		return "", errf(CodeUnknownKind, kind)
	}
	canon, err := Canonicalize(payloadJSON)
	if err != nil {
		return "", err
	}
	var wire struct {
		Result   *string          `json:"result"`
		Findings *json.RawMessage `json:"findings"`
	}
	if err := json.Unmarshal(canon, &wire); err != nil {
		return "", &Error{Code: CodeMalformed, Err: err}
	}
	if wire.Result == nil {
		return "", errf(CodeMissingField, "result")
	}
	if !resultValues[*wire.Result] {
		return "", errf(CodeBadResult, *wire.Result)
	}
	if wire.Findings == nil || !bytes.HasPrefix(bytes.TrimSpace(*wire.Findings), []byte("[")) {
		return "", errf(CodeMissingField, "findings")
	}
	return *wire.Result, nil
}

// ResultAdvances reports whether a report result may advance its stage. Only
// pass advances; fail, insufficient, and error are retained but block.
func ResultAdvances(result string) bool {
	return result == ResultPass
}

// Graph is a validated set of manifests keyed by content ID.
type Graph struct {
	nodes map[string]*Manifest
}

// SubmittedNode is one manifest submitted for graph validation: a caller-supplied
// ID and its raw bytes. The ID must equal the computed content ID.
type SubmittedNode struct {
	ID  string
	Raw []byte
}

// NewGraph verifies every submitted ID matches its computed content ID, rejects
// repeated IDs, validates each manifest, and applies the graph size limits.
func NewGraph(submitted []SubmittedNode) (*Graph, error) {
	if len(submitted) > maxNodes {
		return nil, errf(CodeTooLarge, "nodes")
	}
	total := 0
	seen := make(map[string]bool, len(submitted))
	g := &Graph{nodes: make(map[string]*Manifest, len(submitted))}
	for _, s := range submitted {
		if len(s.Raw) > maxManifestInput {
			return nil, errf(CodeTooLarge, "manifest")
		}
		total += len(s.Raw)
		canon, err := Canonicalize(s.Raw)
		if err != nil {
			return nil, err
		}
		id := HashBytes(canon)
		if id != s.ID {
			return nil, errf(CodeIDMismatch, s.ID)
		}
		if seen[id] {
			return nil, errf(CodeDupID, id)
		}
		seen[id] = true
		m, err := ParseManifest(canon)
		if err != nil {
			return nil, err
		}
		g.nodes[id] = m
	}
	if total > maxGraphInput {
		return nil, errf(CodeTooLarge, "graph")
	}
	return g, nil
}

// Validate checks that the selected roots are present, every reference resolves,
// the graph is acyclic, and every input references the kind its slot requires.
// It is offline and fail-closed.
func (g *Graph) Validate(roots []string) error {
	for _, r := range roots {
		if _, ok := g.nodes[r]; !ok {
			return errf(CodeMissingDep, r)
		}
	}
	// Every reference must resolve before structural checks, so a dangling
	// edge is reported as missing-dependency, not misread as a cycle.
	for id, m := range g.nodes {
		for slot, in := range m.Inputs {
			if _, ok := g.nodes[in]; !ok {
				return errf(CodeMissingDep, id+":"+slot)
			}
		}
	}
	if err := g.checkAcyclic(); err != nil {
		return err
	}
	for id, m := range g.nodes {
		req := requiredFor(m.Kind)
		for slot, in := range m.Inputs {
			if g.nodes[in].Kind != req[slot] {
				return errf(CodeInputMismatch, id+":"+slot)
			}
		}
	}
	return nil
}

// checkAcyclic runs Kahn's algorithm over the reference edges (iterative and
// bounded). It returns a cycle error, naming a node still in the cycle.
func (g *Graph) checkAcyclic() error {
	n := len(g.nodes)
	indeg := make(map[string]int, n)
	adj := make(map[string][]string, n)
	for id := range g.nodes {
		indeg[id] = 0
	}
	for id, m := range g.nodes {
		if len(m.Inputs) > maxEdgesPerNode {
			return errf(CodeTooLarge, id)
		}
		for _, in := range m.Inputs {
			adj[id] = append(adj[id], in)
			indeg[in]++
		}
	}
	queue := make([]string, 0, n)
	for id, d := range indeg {
		if d == 0 {
			queue = append(queue, id)
		}
	}
	processed := 0
	for len(queue) > 0 {
		id := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		processed++
		for _, in := range adj[id] {
			indeg[in]--
			if indeg[in] == 0 {
				queue = append(queue, in)
			}
		}
	}
	if processed != n {
		subject := ""
		for id, d := range indeg {
			if d > 0 {
				subject = id
				break
			}
		}
		return errf(CodeCycle, subject)
	}
	return nil
}

// StaleIDs returns the set of nodes that are not current for the given current
// leaf identities. A leaf is current iff it is in currentLeaves; a non-leaf is
// current iff all of its inputs are current. Changing a leaf therefore invalidates
// precisely its transitive dependents and no siblings.
func (g *Graph) StaleIDs(currentLeaves map[string]bool) map[string]bool {
	memo := make(map[string]bool, len(g.nodes))
	stale := map[string]bool{}
	for id := range g.nodes {
		if !g.isCurrent(id, currentLeaves, memo) {
			stale[id] = true
		}
	}
	return stale
}

func (g *Graph) isCurrent(id string, currentLeaves, memo map[string]bool) bool {
	if v, ok := memo[id]; ok {
		return v
	}
	m := g.nodes[id]
	var cur bool
	if leafKinds[m.Kind] {
		cur = currentLeaves[id]
	} else {
		cur = true
		for _, in := range m.Inputs {
			if !g.isCurrent(in, currentLeaves, memo) {
				cur = false
				break
			}
		}
	}
	memo[id] = cur
	return cur
}

// LeafIDs returns the IDs of every leaf-kind node in the graph.
func (g *Graph) LeafIDs() []string {
	var out []string
	for id, m := range g.nodes {
		if leafKinds[m.Kind] {
			out = append(out, id)
		}
	}
	return out
}

// Kinds returns the set of kinds present in the graph.
func (g *Graph) Kinds() map[string]bool {
	out := map[string]bool{}
	for _, m := range g.nodes {
		out[m.Kind] = true
	}
	return out
}
