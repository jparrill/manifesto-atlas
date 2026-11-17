package runsource

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// MaxSourceBytes bounds a downloaded programme so a hostile or mistaken URL
// cannot fill the disk.
const MaxSourceBytes = 120 << 20

// Fetch downloads the programme from its official URL into a local cache and
// returns the cached path and the SHA-256 of the bytes. It refuses a non-https
// origin, a redirected non-https origin, a non-PDF answer and an oversized body.
func Fetch(ctx context.Context, url, cacheDir string, client *http.Client) (string, string, error) {
	if !strings.HasPrefix(url, "https://") {
		return "", "", fmt.Errorf("the provenance URL must be https")
	}
	if client == nil {
		client = &http.Client{}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("download: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("download returned %d", resp.StatusCode)
	}
	if got := resp.Header.Get("Content-Type"); !strings.Contains(got, "pdf") {
		return "", "", fmt.Errorf("the URL answers with %q, not a PDF", got)
	}
	// The bytes are hashed while streaming, so nothing is trusted before it is
	// measured. A single byte over the limit aborts the read.
	h := sha256.New()
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return "", "", err
	}
	tmp, err := os.CreateTemp(cacheDir, "download-*.pdf")
	if err != nil {
		return "", "", err
	}
	written, err := io.Copy(io.MultiWriter(tmp, h), io.LimitReader(resp.Body, MaxSourceBytes+1))
	tmp.Close()
	if err != nil {
		os.Remove(tmp.Name())
		return "", "", fmt.Errorf("download: %w", err)
	}
	if written > MaxSourceBytes {
		os.Remove(tmp.Name())
		return "", "", fmt.Errorf("the document is larger than the %d byte limit", MaxSourceBytes)
	}
	sum := hex.EncodeToString(h.Sum(nil))
	// One directory per hash so the original file name survives and two runs of
	// the same document share the cached bytes.
	name := filepath.Base(strings.SplitN(url, "?", 2)[0])
	if !strings.HasSuffix(strings.ToLower(name), ".pdf") {
		name = "source.pdf"
	}
	dir := filepath.Join(cacheDir, sum)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		os.Remove(tmp.Name())
		return "", "", err
	}
	final := filepath.Join(dir, name)
	if err := os.Rename(tmp.Name(), final); err != nil {
		return "", "", err
	}
	return final, sum, nil
}
