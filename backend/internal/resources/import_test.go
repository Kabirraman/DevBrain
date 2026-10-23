package resources

import (
	"net/url"
	"testing"
)

func TestParseImportURLRejectsUnsafeTargets(t *testing.T) {
	for _, rawURL := range []string{
		"file:///etc/passwd",
		"http://127.0.0.1:8080/admin",
		"http://user:password@example.com/article",
		"not a URL",
	} {
		t.Run(rawURL, func(t *testing.T) {
			parsed, err := parseImportURL(rawURL)
			if err == nil {
				if parsedErr := validatePublicURL(parsed); parsedErr == nil {
					t.Fatalf("expected URL to be rejected: %s", rawURL)
				}
			}
		})
	}
}

func TestValidatePublicURLRejectsPrivateIP(t *testing.T) {
	for _, host := range []string{
		"10.0.0.8",
		"100.64.0.1",
		"169.254.169.254",
		"192.0.2.1",
		"198.18.0.1",
		"198.51.100.1",
		"203.0.113.1",
		"[::1]",
		"[fc00::1]",
		"[fe80::1]",
		"[2001:db8::1]",
	} {
		parsed := &url.URL{Scheme: "http", Host: host}
		if err := validatePublicURL(parsed); err == nil {
			t.Fatalf("expected %s to be rejected", host)
		}
	}
}

func TestParseYouTubeID(t *testing.T) {
	cases := map[string]string{
		"https://www.youtube.com/watch?v=abcdefghijk": "abcdefghijk",
		"https://youtu.be/abcdefghijk":                "abcdefghijk",
		"https://youtube.com/shorts/abcdefghijk":      "abcdefghijk",
		"https://m.youtube.com/embed/abcdefghijk":     "abcdefghijk",
	}
	for rawURL, expected := range cases {
		t.Run(rawURL, func(t *testing.T) {
			actual, err := parseYouTubeID(rawURL)
			if err != nil {
				t.Fatalf("parseYouTubeID returned error: %v", err)
			}
			if actual != expected {
				t.Fatalf("got %q, want %q", actual, expected)
			}
		})
	}
	if _, err := parseYouTubeID("https://example.com/watch?v=abcdefghijk"); err == nil {
		t.Fatal("expected non-YouTube host to be rejected")
	}
}
