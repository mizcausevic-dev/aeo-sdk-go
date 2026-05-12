package aeo

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWellKnownURL_StripsTrailingSlashes(t *testing.T) {
	cases := []struct{ in, want string }{
		{"https://example.com", "https://example.com/.well-known/aeo.json"},
		{"https://example.com/", "https://example.com/.well-known/aeo.json"},
		{"https://example.com///", "https://example.com/.well-known/aeo.json"},
	}
	for _, c := range cases {
		if got := WellKnownURL(c.in); got != c.want {
			t.Errorf("WellKnownURL(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestFetchWellKnown_HappyPath(t *testing.T) {
	body := `{
        "aeo_version": "0.1",
        "entity": {
            "id": "https://example.com/#org",
            "type": "Organization",
            "name": "Example Org",
            "canonical_url": "https://example.com/"
        },
        "authority": { "primary_sources": ["https://example.com/"] },
        "claims": [
            { "id": "tagline", "predicate": "description", "value": "test", "confidence": "high" }
        ]
    }`
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/aeo.json", func(w http.ResponseWriter, r *http.Request) {
		if accept := r.Header.Get("Accept"); accept != AcceptHeader {
			t.Errorf("Accept = %q, want %q", accept, AcceptHeader)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	doc, err := DefaultClient().FetchWellKnown(context.Background(), srv.URL)
	if err != nil {
		t.Fatalf("FetchWellKnown: %v", err)
	}
	if doc.Entity.Name != "Example Org" {
		t.Errorf("name = %q, want %q", doc.Entity.Name, "Example Org")
	}
}

func TestFetchWellKnown_404(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/aeo.json", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	_, err := DefaultClient().FetchWellKnown(context.Background(), srv.URL)
	if err == nil {
		t.Fatal("expected error for 404, got nil")
	}
	if !IsHTTPStatusError(err) {
		t.Errorf("expected HTTPStatusError, got %T: %v", err, err)
	}
}
