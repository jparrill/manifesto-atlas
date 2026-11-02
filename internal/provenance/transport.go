package provenance

// Transport for TQ-0022: a bounded, SSRF-protected HTTP download adapter for
// atlas.provenance/v1. It owns network access for canonical evidence retrieval.
//
// The security core is deterministic and tested in isolation: the document
// request URL policy, the public-address (SSRF) resolution policy, and per-hop
// redirect revalidation. Actual TCP/TLS I/O is performed through an injected
// DialContext seam, so tests run against synthetic local fixtures with a
// controlled resolver and never reach the real network. In production the
// caller supplies a dialer that connects to the pinned resolved IP.

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/jparrill/manifesto-atlas/internal/artifact"
)

// Transport error codes.
const (
	CodeBadScheme        = "bad-scheme"
	CodeUserinfo         = "userinfo"
	CodeIPLiteral        = "ip-literal"
	CodeFragment         = "fragment"
	CodeBearerQuery      = "bearer-query"
	CodeNonPublic        = "non-public-address"
	CodeTooManyRedirects = "too-many-redirects"
	CodeBadStatus        = "bad-status"
	CodeBadMIME          = "bad-mime"
	CodeOversized        = "oversized"
	CodeTimeout          = "timeout"
	CodeCancellation     = "cancellation"
	CodeEmpty            = "empty-response"
	CodeBadFilename      = "unsafe-filename"
)

// Defaults from the approved isolation budget.
const (
	defaultMaxBytes     int64 = 1 << 30
	defaultMaxRedirects       = 5
	defaultTimeout            = 120 * time.Second
	maxFilenameBytes          = 255
)

// specialIPv4Nets are special-purpose IPv4 ranges that IsGlobalUnicast does not
// exclude but the SSRF policy must reject.
var specialIPv4Nets = func() []*net.IPNet {
	cidrs := []string{
		"192.0.2.0/24", "198.51.100.0/24", "203.0.113.0/24", // TEST-NET-1/2/3
		"100.64.0.0/10", // carrier NAT
		"0.0.0.0/8",     // "this" network
		"240.0.0.0/4",   // reserved
		"192.0.0.0/29",
		"198.18.0.0/15", // benchmarking
		"192.88.99.0/24",
	}
	nets := make([]*net.IPNet, 0, len(cidrs))
	for _, c := range cidrs {
		_, n, _ := net.ParseCIDR(c)
		nets = append(nets, n)
	}
	return nets
}()

// Policy is the bounded network policy for a fetch.
type Policy struct {
	MaxBytes     int64
	MaxRedirects int
	Timeout      time.Duration
	// Resolver maps a hostname to its A/AAAA answers. Production uses system DNS;
	// tests inject a controlled resolver.
	Resolver func(ctx context.Context, host string) ([]net.IP, error)
	// DialContext performs the actual TCP connection. Production dials the pinned
	// resolved IP; tests inject a dialer that reaches a local fixture.
	DialContext func(ctx context.Context, network, addr string) (net.Conn, error)
	// TLSClientConfig configures the TLS handshake. Production uses system roots;
	// tests may trust a fixture CA or skip verification for a local fixture.
	TLSClientConfig *tls.Config
}

// Redirect records one followed hop.
type Redirect struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// FetchRequest is one bounded download attempt.
type FetchRequest struct {
	URL             string
	ReferralPageURL *string
}

// Evidence is the immutable retrieval result.
type Evidence struct {
	OriginalURL     string     `json:"original_url"`
	FinalURL        string     `json:"final_url"`
	ReferralPageURL *string    `json:"referral_page_url"`
	Redirects       []Redirect `json:"redirects"`
	Status          int        `json:"status"`
	ContentType     string     `json:"content_type"`
	Filename        string     `json:"filename"`
	Bytes           []byte     `json:"bytes"`
	Checksum        string     `json:"checksum"`
	RetrievedAt     string     `json:"retrieved_at"`
}

// IsPublicIP reports whether ip is publicly routable, rejecting loopback,
// private, link-local, multicast, unspecified, and special-purpose ranges
// (including TEST-NET and carrier NAT) that IsGlobalUnicast alone misses.
func IsPublicIP(ip net.IP) bool {
	if ip == nil {
		return false
	}
	if v4 := ip.To4(); v4 != nil {
		ip = v4
	}
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() || ip.IsInterfaceLocalMulticast() ||
		ip.IsMulticast() || ip.IsUnspecified() {
		return false
	}
	for _, n := range specialIPv4Nets {
		if n.Contains(ip) {
			return false
		}
	}
	return ip.IsGlobalUnicast()
}

// ValidateRequestURL applies the document-request URL policy: HTTPS only, no
// userinfo, no IP-literal or Unicode host, no fragment, no bearer-style query
// fields, and a valid lowercase A-label host.
func ValidateRequestURL(raw string) (*url.URL, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return nil, errf(CodeBadScheme, raw)
	}
	if u.Scheme != "https" {
		return nil, errf(CodeBadScheme, u.Scheme)
	}
	if u.User != nil {
		return nil, errf(CodeUserinfo, raw)
	}
	if u.Fragment != "" {
		return nil, errf(CodeFragment, raw)
	}
	host := u.Hostname()
	if net.ParseIP(host) != nil {
		return nil, errf(CodeIPLiteral, host)
	}
	if err := validAHostname(host); err != nil {
		return nil, err
	}
	for key := range u.Query() {
		lk := strings.ToLower(key)
		switch lk {
		case "token", "access_token", "auth", "key", "signature", "sig":
			return nil, errf(CodeBearerQuery, key)
		}
	}
	return u, nil
}

// ResolvePublicIP resolves host and returns a single pinned public address. It
// rejects the whole answer set if any address is not publicly routable, so an
// empty or mixed public/private answer fails closed.
func ResolvePublicIP(ctx context.Context, host string, resolver func(ctx context.Context, host string) ([]net.IP, error)) (net.IP, error) {
	if resolver == nil {
		resolver = defaultResolver
	}
	addrs, err := resolver(ctx, host)
	if err != nil || len(addrs) == 0 {
		return nil, errf(CodeNonPublic, host)
	}
	for _, ip := range addrs {
		if !IsPublicIP(ip) {
			return nil, errf(CodeNonPublic, host)
		}
	}
	return addrs[0], nil
}

// FetchEvidence performs one bounded GET under the policy, revalidating every
// redirect hop, and returns immutable retrieval evidence. A failure produces no
// evidence and no registered file; cancellation closes the response and any
// partial scratch.
func FetchEvidence(ctx context.Context, req FetchRequest, pol Policy) (*Evidence, error) {
	if pol.MaxBytes <= 0 {
		pol.MaxBytes = defaultMaxBytes
	}
	if pol.MaxRedirects <= 0 {
		pol.MaxRedirects = defaultMaxRedirects
	}
	if pol.Timeout <= 0 {
		pol.Timeout = defaultTimeout
	}
	if pol.DialContext == nil {
		return nil, errf(CodeMissingField, "dialer")
	}
	if pol.Resolver == nil {
		pol.Resolver = defaultResolver
	}

	initial, err := ValidateRequestURL(req.URL)
	if err != nil {
		return nil, err
	}
	if _, err := ResolvePublicIP(ctx, initial.Hostname(), pol.Resolver); err != nil {
		return nil, err
	}

	var redirects []Redirect
	transport := &http.Transport{
		DialContext:         pol.DialContext,
		TLSClientConfig:     pol.TLSClientConfig,
		DisableCompression:  true,
		DisableKeepAlives:   true,
		Proxy:               nil,
		ProxyConnectHeader:  nil,
		MaxIdleConnsPerHost: -1,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   pol.Timeout,
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			if len(via) > pol.MaxRedirects {
				return errf(CodeTooManyRedirects, r.URL.String())
			}
			if len(via) > 0 {
				redirects = append(redirects, Redirect{From: via[len(via)-1].URL.String(), To: r.URL.String()})
			}
			if _, err := ValidateRequestURL(r.URL.String()); err != nil {
				return err
			}
			_, err := ResolvePublicIP(ctx, r.URL.Hostname(), pol.Resolver)
			return err
		},
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, initial.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, mapTransportError(ctx, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		drain(resp.Body)
		return nil, errf(CodeBadStatus, strconv.Itoa(resp.StatusCode))
	}
	// A declared non-PDF MIME rejects; an absent MIME defers to parser+magic.
	if ct := resp.Header.Get("Content-Type"); ct != "" && !strings.Contains(strings.ToLower(ct), "pdf") {
		drain(resp.Body)
		return nil, errf(CodeBadMIME, ct)
	}
	body, err := readBounded(resp.Body, pol.MaxBytes)
	if err != nil {
		return nil, err
	}
	if len(body) == 0 {
		return nil, errf(CodeEmpty, req.URL)
	}
	name, err := safeFilename(resp.Request.URL.Path, resp.Header.Get("Content-Disposition"))
	if err != nil {
		return nil, err
	}
	return &Evidence{
		OriginalURL:     req.URL,
		FinalURL:        resp.Request.URL.String(),
		ReferralPageURL: req.ReferralPageURL,
		Redirects:       redirects,
		Status:          resp.StatusCode,
		ContentType:     resp.Header.Get("Content-Type"),
		Filename:        name,
		Bytes:           body,
		Checksum:        artifact.HashBytes(body),
		RetrievedAt:     time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func mapTransportError(ctx context.Context, err error) error {
	if ctx.Err() != nil {
		return errf(CodeCancellation, ctx.Err().Error())
	}
	// A policy error returned by CheckRedirect is wrapped in *url.Error.
	var ue *url.Error
	if errors.As(err, &ue) {
		if pe, ok := ue.Err.(*Error); ok {
			return pe
		}
		if errors.Is(ue.Err, context.DeadlineExceeded) {
			return errf(CodeTimeout, "deadline")
		}
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return errf(CodeTimeout, "deadline")
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return errf(CodeTimeout, "network")
	}
	return err
}

func readBounded(r io.Reader, max int64) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(r, max+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > max {
		return nil, errf(CodeOversized, strconv.FormatInt(max, 10))
	}
	return data, nil
}

func drain(r io.Reader) {
	_, _ = io.Copy(io.Discard, io.LimitReader(r, 1<<20))
}

// safeFilename derives a display-only basename from the final path or
// Content-Disposition, rejecting traversal, separators, and over-long names. It
// never produces an output path; it is metadata only.
func safeFilename(urlPath, contentDisposition string) (string, error) {
	name := ""
	if i := strings.Index(contentDisposition, "filename="); i >= 0 {
		rest := contentDisposition[i+len("filename="):]
		rest = strings.TrimPrefix(rest, "\"")
		if end := strings.Index(rest, "\""); end >= 0 {
			rest = rest[:end]
		}
		name = rest
	}
	if name == "" {
		name = path.Base(urlPath)
	}
	if name == "" || name == "." || name == ".." {
		return "", errf(CodeBadFilename, name)
	}
	if strings.ContainsAny(name, "/\\") || strings.Contains(name, "..") || strings.HasPrefix(name, "-") {
		return "", errf(CodeBadFilename, name)
	}
	if len(name) > maxFilenameBytes {
		return "", errf(CodeBadFilename, name)
	}
	return name, nil
}

func defaultResolver(ctx context.Context, host string) ([]net.IP, error) {
	addrs, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, err
	}
	ips := make([]net.IP, 0, len(addrs))
	for _, a := range addrs {
		ips = append(ips, a.IP)
	}
	return ips, nil
}
