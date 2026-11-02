package artifact

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

// --- fixtures / helpers -------------------------------------------------

func mustCanonic(t *testing.T, v any) []byte {
	t.Helper()
	b, err := encodeCanonical(v)
	if err != nil {
		t.Fatalf("encodeCanonical: %v", err)
	}
	return b
}

// makeManifest builds the raw canonical bytes for a kind with the given payload
// bytes and input slot->id map.
func makeManifest(t *testing.T, kind string, payload []byte, inputs map[string]string) []byte {
	t.Helper()
	if inputs == nil {
		inputs = map[string]string{}
	}
	man := map[string]any{
		"contract": Contract,
		"kind":     kind,
		"payload": map[string]any{
			"media_type": "text/plain",
			"sha256":     HashBytes(payload),
			"size":       int64(len(payload)),
		},
		"inputs": inputs,
	}
	return mustCanonic(t, man)
}

// built accumulates validated manifests keyed by content ID.
type built struct {
	node map[string]*Manifest
}

func newBuilt(t *testing.T) *built {
	t.Helper()
	return &built{node: map[string]*Manifest{}}
}

// add stores one manifest and returns its content ID.
func (b *built) add(t *testing.T, kind string, payload []byte, inputs map[string]string) string {
	t.Helper()
	raw := makeManifest(t, kind, payload, inputs)
	id := HashBytes(raw)
	m, err := ParseManifest(raw)
	if err != nil {
		t.Fatalf("ParseManifest(%s): %v", kind, err)
	}
	b.node[id] = m
	return id
}

// graphFrom assembles a validated Graph from the built nodes.
func (b *built) graphFrom(t *testing.T) *Graph {
	t.Helper()
	submitted := make([]SubmittedNode, 0, len(b.node))
	for id, m := range b.node {
		// Rebuild the canonical raw from the manifest to recover bytes.
		man := map[string]any{
			"contract": m.Contract,
			"kind":     m.Kind,
			"payload": map[string]any{
				"media_type": m.Payload.MediaType,
				"sha256":     m.Payload.SHA256,
				"size":       m.Payload.Size,
			},
			"inputs": m.Inputs,
		}
		submitted = append(submitted, SubmittedNode{ID: id, Raw: mustCanonic(t, man)})
	}
	g, err := NewGraph(submitted)
	if err != nil {
		t.Fatalf("NewGraph: %v", err)
	}
	return g
}

// allLeavesCurrent returns a set of all leaf IDs in g marked current.
func allLeavesCurrent(t *testing.T, g *Graph) map[string]bool {
	t.Helper()
	set := map[string]bool{}
	for _, id := range g.LeafIDs() {
		set[id] = true
	}
	return set
}

// descendants returns id and every node that transitively depends on it.
func descendants(g *Graph, id string) map[string]bool {
	rev := map[string][]string{}
	for nid, m := range g.nodes {
		for _, in := range m.Inputs {
			rev[in] = append(rev[in], nid)
		}
	}
	seen := map[string]bool{id: true}
	stack := []string{id}
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, dep := range rev[cur] {
			if !seen[dep] {
				seen[dep] = true
				stack = append(stack, dep)
			}
		}
	}
	return seen
}

// buildWorkedGraph constructs the full worked pipeline graph from
// docs/design/artifact-contract.md using synthetic payload bytes.
func buildWorkedGraph(t *testing.T) *Graph {
	t.Helper()
	b := newBuilt(t)
	leaf := func(kind, label string) string {
		return b.add(t, kind, []byte(label), nil)
	}
	nonLeaf := func(kind, label string, inputs map[string]string) string {
		full := map[string]string{
			"schema":        leaf("schema", "schema:"+kind),
			"tool":          leaf("tool", "tool:"+kind),
			"configuration": leaf("configuration", "configuration:"+kind),
		}
		for k, v := range inputs {
			if _, dup := full[k]; dup {
				t.Fatalf("worked-graph: duplicate slot %q", k)
			}
			full[k] = v
		}
		return b.add(t, kind, []byte(label), full)
	}

	A := leaf("source", "candidatePDF")
	B := leaf("source", "canonicalPDF")
	R := leaf("trust-policy", "officialDomainEvidence")
	E := leaf("configuration", "retrievalEvidenceInventory")
	Q := leaf("prompt", "provenancePrompt")
	L := leaf("model", "provenanceModelDescriptor")
	MX := leaf("model", "doclingBundle")
	QU := leaf("prompt", "auditPrompt")
	MU := leaf("model", "auditModel")
	F := leaf("configuration", "fidelityFormula")
	N := leaf("taxonomy", "taxonomyRelease")
	J := leaf("configuration", "jurisdictionPolicy")
	QD := leaf("prompt", "normalizerPrompt")
	MD := leaf("model", "normalizerModel")
	QC := leaf("prompt", "claimPrompt")
	MC := leaf("model", "claimModel")
	QO := leaf("prompt", "coveragePrompt")
	MO := leaf("model", "coverageModel")
	HT := leaf("templates", "suppliedAssets")
	Z := leaf("configuration", "qaConfig")
	QV := leaf("prompt", "visualPrompt")
	MV := leaf("model", "visualModel")

	P := nonLeaf("provenance", "verifiedReport", map[string]string{
		"submitted-source": A, "canonical-source": B, "trust-policy": R,
		"evidence": E, "prompt": Q, "model": L,
	})
	S := nonLeaf("source-record", "immutableRegistration", map[string]string{
		"source": B, "provenance": P,
	})
	X := nonLeaf("extraction", "pageMappedJSON", map[string]string{
		"source-record": S, "extraction-model": MX,
	})
	Y := nonLeaf("secondary-evidence", "secondaryText", map[string]string{
		"source-record": S,
	})
	I := nonLeaf("page-images", "originalPageImageInventory", map[string]string{
		"source-record": S,
	})
	U := nonLeaf("extraction-audit", "auditReport", map[string]string{
		"extraction": X, "secondary-evidence": Y, "page-images": I,
		"prompt": QU, "model": MU,
	})
	V := nonLeaf("extraction-approval", "extractionDecision", map[string]string{
		"audit": U, "extraction": X, "formula": F,
	})
	D := nonLeaf("programme", "canonicalDraft", map[string]string{
		"extraction-approval": V, "extraction": X, "taxonomy": N,
		"jurisdiction": J, "prompt": QD, "model": MD,
	})
	C := nonLeaf("claim-report", "claimReport", map[string]string{
		"programme": D, "extraction": X, "source-record": S,
		"page-images": I, "prompt": QC, "model": MC,
	})
	O := nonLeaf("coverage-report", "coverageReport", map[string]string{
		"programme": D, "extraction": X, "source-record": S,
		"page-images": I, "prompt": QO, "model": MO,
	})
	G := nonLeaf("content-approval", "contentDecision", map[string]string{
		"programme": D, "claim-report": C, "coverage-report": O, "extraction-approval": V,
	})
	W := nonLeaf("markdown", "spanishMarkdown", map[string]string{
		"programme": D, "content-approval": G,
	})
	H := nonLeaf("html", "suppliedTemplateRendering", map[string]string{
		"programme": D, "content-approval": G, "templates": HT,
	})
	K := nonLeaf("site-deterministic-report", "deterministicQA", map[string]string{
		"programme": D, "markdown": W, "html": H, "qa": Z,
	})
	ZV := nonLeaf("site-visual-report", "visualQA", map[string]string{
		"programme": D, "html": H, "deterministic-report": K,
		"prompt": QV, "model": MV, "qa": Z,
	})
	GS := nonLeaf("site-approval", "siteDecision", map[string]string{
		"deterministic-report": K, "visual-report": ZV,
	})
	PUB := nonLeaf("publication", "distributionInventory", map[string]string{
		"source-record": S, "extraction-approval": V, "content-approval": G,
		"markdown": W, "html": H, "site-approval": GS, "trust-policy": R,
	})
	_ = PUB
	return b.graphFrom(t)
}

// --- canonicalization ---------------------------------------------------

func TestCanonicalReferenceVectors(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
		code string
	}{
		{"sorted", `{"b":2,"a":1}`, `{"a":1,"b":2}`, ""},
		{"whitespace", `{"a":1, "b":2}`, `{"a":1,"b":2}`, ""},
		{"noHTMLEscape", `{"quote":"<salud>&"}`, `{"quote":"<salud>&"}`, ""},
		{"dupMember", `{"a":1,"a":2}`, "", CodeDupMember},
		{"fractional", `{"score":0.5}`, "", CodeFractional},
		{"exponent", `{"n":1e0}`, "", CodeFractional},
		{"negZero", `{"n":-0}`, "", CodeNegZero},
		{"outOfRange", `{"n":9007199254740992}`, "", CodeOutOfRange},
		{"trailing", `{} {}`, "", CodeTrailing},
		{"twoValues", `[1] 2`, "", CodeTrailing},
		{"bom", "\ufeff{}", "", CodeMalformed},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := Canonicalize([]byte(tc.in))
			if tc.code != "" {
				if err == nil {
					t.Fatalf("want error %s, got none (out=%s)", tc.code, out)
				}
				if codeOf(t, err) != tc.code {
					t.Fatalf("want code %s, got %v", tc.code, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if string(out) != tc.want {
				t.Fatalf("canon = %s, want %s", out, tc.want)
			}
		})
	}
}

func TestCanonicalInvalidUTF8AndSurrogates(t *testing.T) {
	// Invalid raw UTF-8 byte in a string.
	if _, err := Canonicalize([]byte(`{"s":"\xff"}`)); err == nil {
		t.Error("want invalid-utf8 rejection")
	}
	// Unpaired high surrogate.
	if _, err := Canonicalize([]byte(`{"s":"\ud800"}`)); err == nil {
		t.Error("want unpaired-surrogate rejection")
	}
	// Invalid surrogate pair (high then high, not a low surrogate).
	if _, err := Canonicalize([]byte(`{"s":"\ud800\ud800"}`)); err == nil {
		t.Error("want invalid-surrogate-pair rejection")
	}
	// A valid pair encodes U+10000.
	out, err := Canonicalize([]byte(`{"s":"\ud800\udc00"}`))
	if err != nil {
		t.Fatalf("valid pair rejected: %v", err)
	}
	if !bytes.Contains(out, []byte("\U00010000")) {
		t.Fatalf("valid pair not encoded as U+10000: %s", out)
	}
}

func TestCanonicalDepthLimit(t *testing.T) {
	// 64 nested arrays is allowed (root is depth 1).
	ok := strings.Repeat("[", 64) + strings.Repeat("]", 64)
	if _, err := Canonicalize([]byte(ok)); err != nil {
		t.Fatalf("depth 64 rejected: %v", err)
	}
	// 65 nested arrays is rejected.
	too := strings.Repeat("[", 65) + strings.Repeat("]", 65)
	if _, err := Canonicalize([]byte(too)); err == nil {
		t.Fatal("depth 65 not rejected")
	}
}

func TestLimits(t *testing.T) {
	// Canonicalize rejects input over the 16 MiB payload bound.
	over := make([]byte, 16<<20+1)
	over[0] = '{'
	if _, err := Canonicalize(over); err == nil {
		t.Fatal("oversized payload accepted")
	}
	// NewGraph rejects a single manifest over the 1 MiB bound.
	big := make([]byte, 1<<20+1)
	if _, err := NewGraph([]SubmittedNode{{ID: "sha256:" + strings.Repeat("0", 64), Raw: big}}); err == nil {
		t.Fatal("oversized manifest accepted")
	}
}

func TestCanonicalIntegerRange(t *testing.T) {
	if _, err := Canonicalize([]byte(`{"n":9007199254740991}`)); err != nil {
		t.Fatalf("max int rejected: %v", err)
	}
	if _, err := Canonicalize([]byte(`{"n":-9007199254740991}`)); err != nil {
		t.Fatalf("min int rejected: %v", err)
	}
	if _, err := Canonicalize([]byte(`{"n":9007199254740992}`)); err == nil {
		t.Fatal("max+1 not rejected")
	}
}

func TestDeterminism(t *testing.T) {
	// Same logical content, different key order and indentation -> same bytes/ID.
	a := `{"kind":"source","contract":"atlas.artifact/v1","inputs":{},"payload":{"size":3,"sha256":"sha256:0000000000000000000000000000000000000000000000000000000000000000","media_type":"text/plain"}}`
	b := "{\n  \"payload\": { \"media_type\": \"text/plain\", \"sha256\": \"sha256:0000000000000000000000000000000000000000000000000000000000000000\", \"size\": 3 },\n  \"kind\": \"source\",\n  \"inputs\": {},\n  \"contract\": \"atlas.artifact/v1\"\n}"
	ca, errA := Canonicalize([]byte(a))
	cb, errB := Canonicalize([]byte(b))
	if errA != nil || errB != nil {
		t.Fatalf("canonicalize failed: %v / %v", errA, errB)
	}
	if !bytes.Equal(ca, cb) {
		t.Fatalf("canonical bytes differ:\n%s\n%s", ca, cb)
	}
	ida, _ := ManifestID([]byte(a))
	idb, _ := ManifestID([]byte(b))
	if ida != idb {
		t.Fatalf("IDs differ: %s vs %s", ida, idb)
	}
}

// --- manifest validation ------------------------------------------------

func TestManifestStructure(t *testing.T) {
	// Valid leaf.
	raw := makeManifest(t, "source", []byte("x"), nil)
	if _, err := ParseManifest(raw); err != nil {
		t.Fatalf("valid leaf rejected: %v", err)
	}
	// Unknown contract version.
	bad := makeManifest(t, "source", []byte("x"), nil)
	bad = bytes.Replace(bad, []byte(Contract), []byte("atlas.artifact/v2"), 1)
	if _, err := ParseManifest(bad); err == nil {
		t.Fatal("v2 contract accepted")
	}
	// Unknown kind.
	rawU := makeManifest(t, "bogus", []byte("x"), nil)
	if _, err := ParseManifest(rawU); err == nil {
		t.Fatal("unknown kind accepted")
	}
	// Leaf with an input -> unknown slot.
	rawLeafIn := makeManifest(t, "source", []byte("x"), map[string]string{"schema": strings.Repeat("a", 71)})
	if _, err := ParseManifest(rawLeafIn); err == nil {
		t.Fatal("leaf with input accepted")
	}
	// Non-leaf missing a required slot.
	rawMissing := makeManifest(t, "extraction", []byte("x"), map[string]string{
		"source-record": strings.Repeat("a", 71), "extraction-model": strings.Repeat("b", 71),
	})
	if _, err := ParseManifest(rawMissing); err == nil {
		t.Fatal("non-leaf missing base inputs accepted")
	}
	// Unknown slot on a non-leaf.
	rawUnknown := makeManifest(t, "extraction", []byte("x"), map[string]string{
		"schema": strings.Repeat("a", 71), "tool": strings.Repeat("b", 71),
		"configuration": strings.Repeat("c", 71), "source-record": strings.Repeat("d", 71),
		"extraction-model": strings.Repeat("e", 71), "extra": strings.Repeat("f", 71),
	})
	if _, err := ParseManifest(rawUnknown); err == nil {
		t.Fatal("unknown slot accepted")
	}
	// Bad hash.
	badHash := makeManifest(t, "source", []byte("x"), nil)
	badHash = bytes.Replace(badHash, []byte(HashBytes([]byte("x"))), []byte("sha256:ZZZ"), 1)
	if _, err := ParseManifest(badHash); err == nil {
		t.Fatal("bad payload hash accepted")
	}
	// Unknown manifest field.
	extra := makeManifest(t, "source", []byte("x"), nil)
	extra = bytes.Replace(extra, []byte(`"inputs":{}`), []byte(`"inputs":{},"bogus":1`), 1)
	if _, err := ParseManifest(extra); err == nil {
		t.Fatal("unknown manifest field accepted")
	}
}

func TestValidatePayload(t *testing.T) {
	payload := []byte("hello")
	m, err := ParseManifest(makeManifest(t, "source", payload, nil))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if err := ValidatePayload(m, payload); err != nil {
		t.Fatalf("valid payload rejected: %v", err)
	}
	if err := ValidatePayload(m, []byte("hellO")); err == nil {
		t.Fatal("changed bytes accepted")
	}
	if err := ValidatePayload(m, []byte("he")); err == nil {
		t.Fatal("wrong size accepted")
	}
}

// --- report outcomes ----------------------------------------------------

func TestReportResult(t *testing.T) {
	pass := mustCanonic(t, map[string]any{"result": "pass", "findings": []any{}, "score": 900})
	r, err := ReportResult("claim-report", pass)
	if err != nil || r != "pass" {
		t.Fatalf("pass: got %q, %v", r, err)
	}
	if !ResultAdvances(r) {
		t.Fatal("pass should advance")
	}
	fail := mustCanonic(t, map[string]any{"result": "fail", "findings": []any{map[string]any{"id": "f1"}}})
	rFail, err := ReportResult("coverage-report", fail)
	if err != nil || rFail != "fail" {
		t.Fatalf("fail: got %q, %v", rFail, err)
	}
	if ResultAdvances(rFail) {
		t.Fatal("fail must not advance")
	}
	insuf := mustCanonic(t, map[string]any{"result": "insufficient", "findings": []any{}})
	rInsuf, err := ReportResult("publication", insuf)
	if err != nil || rInsuf != "insufficient" {
		t.Fatalf("insufficient: got %q, %v", rInsuf, err)
	}
	if ResultAdvances(rInsuf) {
		t.Fatal("insufficient must not advance")
	}
	// Bad enum.
	bad := mustCanonic(t, map[string]any{"result": "maybe", "findings": []any{}})
	if _, err := ReportResult("claim-report", bad); err == nil {
		t.Fatal("bad result enum accepted")
	}
	// Missing findings.
	noFind := mustCanonic(t, map[string]any{"result": "pass"})
	if _, err := ReportResult("claim-report", noFind); err == nil {
		t.Fatal("missing findings accepted")
	}
	// Non-report kind.
	if _, err := ReportResult("source", pass); err == nil {
		t.Fatal("non-report kind accepted")
	}
}

// --- graph validation ---------------------------------------------------

func TestWorkedGraphValidates(t *testing.T) {
	g := buildWorkedGraph(t)
	// Every v1 kind is represented.
	kinds := g.Kinds()
	if len(kinds) != 26 {
		t.Fatalf("kinds present = %d, want 26: %v", len(kinds), kinds)
	}
	// Find the publication node as root.
	var pub string
	for id, m := range g.nodes {
		if m.Kind == "publication" {
			pub = id
		}
	}
	if pub == "" {
		t.Fatal("no publication node")
	}
	if err := g.Validate([]string{pub}); err != nil {
		t.Fatalf("worked graph failed validation: %v", err)
	}
}

func TestStalenessPrecise(t *testing.T) {
	g := buildWorkedGraph(t)
	current := allLeavesCurrent(t, g)
	// With all leaves current, nothing is stale.
	if stale := g.StaleIDs(current); len(stale) != 0 {
		t.Fatalf("expected no stale, got %d", len(stale))
	}
	// Find representative leaves and assert each invalidates precisely its descendants.
	byKind := map[string]string{}
	for id, m := range g.nodes {
		if leafKinds[m.Kind] {
			if _, ok := byKind[m.Kind]; !ok {
				byKind[m.Kind] = id
			}
		}
	}
	pick := []string{"model", "source", "schema", "taxonomy", "prompt", "templates", "configuration"}
	for _, kind := range pick {
		id, ok := byKind[kind]
		if !ok {
			t.Fatalf("no leaf of kind %s", kind)
		}
		got := g.StaleIDs(removeOne(current, id))
		want := descendants(g, id)
		if !setEqual(got, want) {
			t.Fatalf("kind %s: stale != descendants.\n got=%v\nwant=%v", kind, got, want)
		}
	}
}

func TestDanglingReference(t *testing.T) {
	b := newBuilt(t)
	// A non-leaf whose source-record input points at a missing node.
	missing := strings.Repeat("a", 71)
	b.add(t, "extraction", []byte("x"), map[string]string{
		"schema":        b.add(t, "schema", []byte("s"), nil),
		"tool":          b.add(t, "tool", []byte("t"), nil),
		"configuration": b.add(t, "configuration", []byte("c"), nil),
		"source-record": missing, "extraction-model": b.add(t, "model", []byte("m"), nil),
	})
	g := b.graphFrom(t)
	var root string
	for id, m := range g.nodes {
		if m.Kind == "extraction" {
			root = id
		}
	}
	err := g.Validate([]string{root})
	if err == nil {
		t.Fatal("dangling reference accepted")
	}
	if codeOf(t, err) != CodeMissingDep {
		t.Fatalf("want %s, got %v", CodeMissingDep, err)
	}
}

func TestGraphCycle(t *testing.T) {
	// A true node-level cycle is unrepresentable with content-addressed IDs (an
	// ID is a hash of the bytes that would have to contain it), so the acyclicity
	// detector is exercised directly on an abstract cyclic edge graph. This is
	// defense-in-depth: the kind DAG already forbids cycles in well-formed
	// graphs, but the validator must still reject a cycle if one is ever built.
	g := &Graph{nodes: map[string]*Manifest{
		"A": {Kind: "extraction", Inputs: map[string]string{"source-record": "B"}},
		"B": {Kind: "extraction", Inputs: map[string]string{"source-record": "A"}},
	}}
	if err := g.checkAcyclic(); err == nil {
		t.Fatal("cycle accepted")
	} else if codeOf(t, err) != CodeCycle {
		t.Fatalf("want %s, got %v", CodeCycle, err)
	}

	// A self-edge is also a cycle.
	self := &Graph{nodes: map[string]*Manifest{
		"A": {Kind: "extraction", Inputs: map[string]string{"source-record": "A"}},
	}}
	if err := self.checkAcyclic(); err == nil {
		t.Fatal("self-edge accepted")
	} else if codeOf(t, err) != CodeCycle {
		t.Fatalf("self-edge: want %s, got %v", CodeCycle, err)
	}
	// An acyclic diamond is accepted.
	ok := &Graph{nodes: map[string]*Manifest{
		"root":  {Kind: "extraction", Inputs: map[string]string{"source-record": "left", "extraction-model": "right"}},
		"left":  {Kind: "source-record", Inputs: map[string]string{"source": "leaf", "provenance": "prov"}},
		"right": {Kind: "model", Inputs: map[string]string{}},
		"leaf":  {Kind: "source", Inputs: map[string]string{}},
		"prov":  {Kind: "trust-policy", Inputs: map[string]string{}},
	}}
	if err := ok.checkAcyclic(); err != nil {
		t.Fatalf("acyclic graph rejected: %v", err)
	}
}

func TestDuplicateID(t *testing.T) {
	raw := makeManifest(t, "source", []byte("x"), nil)
	id := HashBytes(raw)
	submitted := []SubmittedNode{{ID: id, Raw: raw}, {ID: id, Raw: raw}}
	if _, err := NewGraph(submitted); err == nil {
		t.Fatal("duplicate ID accepted")
	}
}

func TestIDMismatch(t *testing.T) {
	raw := makeManifest(t, "source", []byte("x"), nil)
	submitted := []SubmittedNode{{ID: "sha256:" + strings.Repeat("f", 64), Raw: raw}}
	if _, err := NewGraph(submitted); err == nil {
		t.Fatal("ID mismatch accepted")
	}
}

func TestInputKindMismatch(t *testing.T) {
	// An extraction whose source-record input points at a model (wrong kind).
	// The reference resolves (no dangling) and the graph is acyclic, so the
	// kind check must fail.
	b := newBuilt(t)
	schemaID := b.add(t, "schema", []byte("s"), nil)
	toolID := b.add(t, "tool", []byte("t"), nil)
	confID := b.add(t, "configuration", []byte("c"), nil)
	modelID := b.add(t, "model", []byte("m"), nil)
	root := b.add(t, "extraction", []byte("x"), map[string]string{
		"schema": schemaID, "tool": toolID, "configuration": confID,
		"source-record": modelID, "extraction-model": modelID,
	})
	g := b.graphFrom(t)
	err := g.Validate([]string{root})
	if err == nil {
		t.Fatal("kind mismatch accepted")
	}
	if codeOf(t, err) != CodeInputMismatch {
		t.Fatalf("want %s, got %v", CodeInputMismatch, err)
	}
}

// --- small test utilities ------------------------------------------------

func removeOne(set map[string]bool, id string) map[string]bool {
	out := map[string]bool{}
	for k, v := range set {
		if k != id {
			out[k] = v
		}
	}
	return out
}

func setEqual(a, b map[string]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if !b[k] {
			return false
		}
	}
	return true
}

func codeOf(t *testing.T, err error) string {
	t.Helper()
	var ae *Error
	if !errors.As(err, &ae) {
		t.Fatalf("error is not *Error: %v", err)
	}
	return ae.Code
}
