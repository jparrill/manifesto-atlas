package provenance_test

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jparrill/manifesto-atlas/internal/artifact"
	"github.com/jparrill/manifesto-atlas/internal/provenance"
)

// fixtureHost is the fake public hostname the fixture dialer routes locally to.
const fixtureHost = "fixture.example.com"

// newFixture starts a local TLS server and returns a Policy whose resolver
// reports a public IP and whose dialer routes to the local server. This
// exercises the full policy (URL, SSRF, redirect, size, status, MIME) while I/O
// stays on the fixture — no real network.
func newFixture(t *testing.T, handler http.Handler) provenance.Policy {
	t.Helper()
	ts := httptest.NewTLSServer(handler)
	t.Cleanup(ts.Close)
	addr := ts.Listener.Addr().String()
	return provenance.Policy{
		MaxBytes:     1 << 20,
		MaxRedirects: 5,
		Timeout:      10 * time.Second,
		Resolver: func(ctx context.Context, host string) ([]net.IP, error) {
			return []net.IP{net.ParseIP("93.184.216.34")}, nil
		},
		DialContext: func(ctx context.Context, network, a string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, network, addr)
		},
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}

func tcode(t *testing.T, err error) string {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	ce, ok := err.(*provenance.Error)
	if !ok {
		t.Fatalf("expected provenance error, got %T: %v", err, err)
	}
	return ce.Code
}

func TestDirectFetch(t *testing.T) {
	body := []byte("%PDF-1.7 fixture bytes")
	pol := newFixture(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/pdf")
		_, _ = w.Write(body)
	}))
	ev, err := provenance.FetchEvidence(context.Background(),
		provenance.FetchRequest{URL: "https://" + fixtureHost + "/doc.pdf"}, pol)
	if err != nil {
		t.Fatalf("direct fetch rejected: %v", err)
	}
	if ev.Checksum != artifact.HashBytes(body) {
		t.Fatalf("checksum mismatch: %s", ev.Checksum)
	}
	if ev.Status != 200 || ev.Filename != "doc.pdf" || ev.OriginalURL != "https://"+fixtureHost+"/doc.pdf" {
		t.Fatalf("unexpected evidence: %+v", ev)
	}
	if len(ev.Redirects) != 0 {
		t.Fatalf("unexpected redirects: %+v", ev.Redirects)
	}
}

func TestMultiHopRedirect(t *testing.T) {
	pol := newFixture(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/a":
			http.Redirect(w, r, "https://"+fixtureHost+"/b", http.StatusFound)
		case "/b":
			http.Redirect(w, r, "https://"+fixtureHost+"/c", http.StatusFound)
		default:
			w.Header().Set("Content-Type", "application/pdf")
			_, _ = w.Write([]byte("final-body"))
		}
	}))
	ev, err := provenance.FetchEvidence(context.Background(),
		provenance.FetchRequest{URL: "https://" + fixtureHost + "/a"}, pol)
	if err != nil {
		t.Fatalf("redirect fetch rejected: %v", err)
	}
	if len(ev.Redirects) != 2 {
		t.Fatalf("redirects = %d, want 2: %+v", len(ev.Redirects), ev.Redirects)
	}
	if ev.FinalURL != "https://"+fixtureHost+"/c" {
		t.Fatalf("final URL = %q", ev.FinalURL)
	}
	if ev.Checksum != artifact.HashBytes([]byte("final-body")) {
		t.Fatal("checksum mismatch")
	}
}

func TestReferralEvidence(t *testing.T) {
	pol := newFixture(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/pdf")
		_, _ = w.Write([]byte("pdf"))
	}))
	ref := "https://party.example.com/programa"
	ev, err := provenance.FetchEvidence(context.Background(),
		provenance.FetchRequest{URL: "https://" + fixtureHost + "/doc.pdf", ReferralPageURL: &ref}, pol)
	if err != nil {
		t.Fatalf("referral fetch rejected: %v", err)
	}
	if ev.ReferralPageURL == nil || *ev.ReferralPageURL != ref {
		t.Fatalf("referral not retained: %+v", ev.ReferralPageURL)
	}
}

func TestDNSAddressPolicy(t *testing.T) {
	cases := []struct {
		ip  string
		pub bool
	}{
		{"93.184.216.34", true},    // public
		{"127.0.0.1", false},       // loopback
		{"10.0.0.1", false},        // private
		{"172.16.0.1", false},      // private
		{"192.168.1.1", false},     // private
		{"169.254.10.10", false},   // link-local
		{"100.64.0.1", false},      // carrier NAT
		{"192.0.2.1", false},       // TEST-NET-1
		{"0.0.0.0", false},         // unspecified
		{"255.255.255.255", false}, // reserved
		{"::1", false},             // IPv6 loopback
		{"fe80::1", false},         // IPv6 link-local
	}
	for _, c := range cases {
		if got := provenance.IsPublicIP(net.ParseIP(c.ip)); got != c.pub {
			t.Errorf("IsPublicIP(%s) = %v, want %v", c.ip, got, c.pub)
		}
	}

	// Mixed public/private answer set fails closed.
	_, err := provenance.ResolvePublicIP(context.Background(), "mixed.example.com",
		func(ctx context.Context, h string) ([]net.IP, error) {
			return []net.IP{net.ParseIP("93.184.216.34"), net.ParseIP("127.0.0.1")}, nil
		})
	if tcode(t, err) != provenance.CodeNonPublic {
		t.Fatal("mixed answers not rejected")
	}

	// A loopback answer blocks the fetch before any dial.
	pol := provenance.Policy{
		MaxBytes: 1 << 20, MaxRedirects: 5, Timeout: 5 * time.Second,
		Resolver: func(ctx context.Context, h string) ([]net.IP, error) {
			return []net.IP{net.ParseIP("127.0.0.1")}, nil
		},
		DialContext: func(ctx context.Context, n, a string) (net.Conn, error) {
			t.Fatal("dialer must not be called for a non-public address")
			return nil, nil
		},
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if _, err := provenance.FetchEvidence(context.Background(),
		provenance.FetchRequest{URL: "https://private.example.com/"}, pol); tcode(t, err) != provenance.CodeNonPublic {
		t.Fatal("loopback destination not blocked")
	}
}

func TestDisallowedScheme(t *testing.T) {
	for _, raw := range []string{"http://example.com/", "ftp://example.com/", "example.com/"} {
		if _, err := provenance.ValidateRequestURL(raw); tcode(t, err) != provenance.CodeBadScheme {
			t.Errorf("scheme not rejected for %q", raw)
		}
	}
}

func TestCredentialBearing(t *testing.T) {
	if _, err := provenance.ValidateRequestURL("https://user:pass@example.com/"); tcode(t, err) != provenance.CodeUserinfo {
		t.Fatal("userinfo not rejected")
	}
	if _, err := provenance.ValidateRequestURL("https://example.com/?token=abc"); tcode(t, err) != provenance.CodeBearerQuery {
		t.Fatal("bearer query not rejected")
	}
}

func TestRedirectLoop(t *testing.T) {
	pol := newFixture(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://"+fixtureHost+"/loop", http.StatusFound)
	}))
	if _, err := provenance.FetchEvidence(context.Background(),
		provenance.FetchRequest{URL: "https://" + fixtureHost + "/loop"}, pol); tcode(t, err) != provenance.CodeTooManyRedirects {
		t.Fatal("redirect loop not bounded")
	}
}

func TestOversized(t *testing.T) {
	pol := newFixture(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/pdf")
		_, _ = w.Write(make([]byte, 100))
	}))
	pol.MaxBytes = 10
	if _, err := provenance.FetchEvidence(context.Background(),
		provenance.FetchRequest{URL: "https://" + fixtureHost + "/big.pdf"}, pol); tcode(t, err) != provenance.CodeOversized {
		t.Fatal("oversized body not rejected")
	}
}

func TestSlowResponse(t *testing.T) {
	pol := newFixture(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		_, _ = w.Write([]byte("late"))
	}))
	pol.Timeout = 300 * time.Millisecond
	if _, err := provenance.FetchEvidence(context.Background(),
		provenance.FetchRequest{URL: "https://" + fixtureHost + "/slow.pdf"}, pol); tcode(t, err) != provenance.CodeTimeout {
		t.Fatal("slow response not bounded")
	}
}

func TestCancellation(t *testing.T) {
	pol := newFixture(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
		case <-time.After(5 * time.Second):
		}
	}))
	ctx, cancel := context.WithCancel(context.Background())
	go func() { time.Sleep(50 * time.Millisecond); cancel() }()
	if _, err := provenance.FetchEvidence(ctx,
		provenance.FetchRequest{URL: "https://" + fixtureHost + "/block.pdf"}, pol); tcode(t, err) != provenance.CodeCancellation {
		t.Fatal("cancellation not surfaced")
	}
}

func TestFilenameTraversal(t *testing.T) {
	pol := newFixture(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", `filename="../../etc/passwd"`)
		_, _ = w.Write([]byte("pdf"))
	}))
	if _, err := provenance.FetchEvidence(context.Background(),
		provenance.FetchRequest{URL: "https://" + fixtureHost + "/doc.pdf"}, pol); tcode(t, err) != provenance.CodeBadFilename {
		t.Fatal("filename traversal not rejected")
	}
}

func TestInvalidMIME(t *testing.T) {
	pol := newFixture(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte("<html>not a pdf</html>"))
	}))
	if _, err := provenance.FetchEvidence(context.Background(),
		provenance.FetchRequest{URL: "https://" + fixtureHost + "/doc.pdf"}, pol); tcode(t, err) != provenance.CodeBadMIME {
		t.Fatal("declared non-PDF MIME not rejected")
	}
}

func TestBadStatus(t *testing.T) {
	pol := newFixture(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	if _, err := provenance.FetchEvidence(context.Background(),
		provenance.FetchRequest{URL: "https://" + fixtureHost + "/missing.pdf"}, pol); tcode(t, err) != provenance.CodeBadStatus {
		t.Fatal("non-200 status not rejected")
	}
}
