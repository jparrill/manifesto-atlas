package runsource

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFetch(t *testing.T) {
	pdf := "%PDF-1.4 contenido"
	t.Run("When the URL serves a PDF it should cache it under its hash", func(t *testing.T) {
		srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/pdf")
			_, _ = w.Write([]byte(pdf))
		}))
		defer srv.Close()
		cache := t.TempDir()
		path, sum, err := Fetch(context.Background(), srv.URL+"/programa.pdf", cache, srv.Client())
		if err != nil {
			t.Fatalf("Fetch: %v", err)
		}
		if filepath.Base(path) != "programa.pdf" || filepath.Base(filepath.Dir(path)) != sum {
			t.Fatalf("cache path %q should be <hash>/programa.pdf", path)
		}
		raw, err := os.ReadFile(path)
		if err != nil || string(raw) != pdf {
			t.Fatalf("cached bytes differ: %v", err)
		}
	})
	t.Run("When the origin is not https it should refuse", func(t *testing.T) {
		if _, _, err := Fetch(context.Background(), "http://example.org/p.pdf", t.TempDir(), nil); err == nil {
			t.Fatalf("expected a scheme error")
		}
	})
	t.Run("When the URL answers with HTML it should refuse", func(t *testing.T) {
		srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write([]byte("<html>error</html>"))
		}))
		defer srv.Close()
		if _, _, err := Fetch(context.Background(), srv.URL+"/programa.pdf", t.TempDir(), srv.Client()); err == nil ||
			!strings.Contains(err.Error(), "not a PDF") {
			t.Fatalf("expected a content-type error, got %v", err)
		}
	})
	t.Run("When the server fails it should report the status", func(t *testing.T) {
		srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "nope", http.StatusNotFound)
		}))
		defer srv.Close()
		if _, _, err := Fetch(context.Background(), srv.URL+"/programa.pdf", t.TempDir(), srv.Client()); err == nil {
			t.Fatalf("expected a status error")
		}
	})
}
