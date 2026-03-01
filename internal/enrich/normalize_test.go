package enrich

import "testing"

func TestNormalizeURLCanonicalization(t *testing.T) {
	got := NormalizeURL("HTTPS://Example.com:443/path?q=1&utm_source=mail#frag")

	if got.Canonical != "https://example.com/path?q=1" {
		t.Fatalf("unexpected canonical URL: %q", got.Canonical)
	}

	if got.Domain != "example.com" {
		t.Fatalf("unexpected domain: %q", got.Domain)
	}
}

func TestNormalizeURLAddsHTTPSWhenMissingScheme(t *testing.T) {
	got := NormalizeURL("example.com/docs?utm_campaign=abc")

	if got.Canonical != "https://example.com/docs" {
		t.Fatalf("unexpected canonical URL: %q", got.Canonical)
	}

	if got.Domain != "example.com" {
		t.Fatalf("unexpected domain: %q", got.Domain)
	}
}

func TestNormalizeURLInvalidReturnsTrimmedInput(t *testing.T) {
	got := NormalizeURL("http://[::1")

	if got.Canonical != "http://[::1" {
		t.Fatalf("unexpected canonical fallback: %q", got.Canonical)
	}

	if got.Domain != "" {
		t.Fatalf("expected empty domain, got %q", got.Domain)
	}
}
