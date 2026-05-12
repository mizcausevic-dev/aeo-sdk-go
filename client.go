package aeo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// WellKnownPath is the canonical well-known URL path for AEO declarations.
const WellKnownPath = "/.well-known/aeo.json"

// AcceptHeader is the value the SDK sends in the HTTP Accept header.
const AcceptHeader = "application/aeo+json, application/json"

// DefaultTimeout is the HTTP timeout used by the default client.
const DefaultTimeout = 10 * time.Second

// WellKnownURL returns the canonical well-known URL for the given origin.
// Trailing slashes on the origin are stripped before appending the path.
func WellKnownURL(origin string) string {
	return strings.TrimRight(origin, "/") + WellKnownPath
}

// Client is a configurable HTTP fetcher for AEO declarations.
type Client struct {
	HTTPClient *http.Client
}

// DefaultClient returns a Client with a sensible default timeout.
func DefaultClient() *Client {
	return &Client{HTTPClient: &http.Client{Timeout: DefaultTimeout}}
}

// FetchWellKnown fetches and parses the AEO declaration at origin's
// well-known URL using the package-level default client.
func FetchWellKnown(origin string) (*Document, error) {
	return DefaultClient().FetchWellKnown(context.Background(), origin)
}

// FetchWellKnown fetches and parses the AEO declaration at origin's
// well-known URL using this client.
func (c *Client) FetchWellKnown(ctx context.Context, origin string) (*Document, error) {
	url := WellKnownURL(origin)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("aeo: build request: %w", err)
	}
	req.Header.Set("Accept", AcceptHeader)

	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: DefaultTimeout}
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("aeo: fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &HTTPStatusError{Status: resp.StatusCode, URL: url}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("aeo: read response: %w", err)
	}

	return ParseDocument(body)
}

// HTTPStatusError is returned when an origin responds with a non-2xx status.
type HTTPStatusError struct {
	Status int
	URL    string
}

// Error implements the error interface.
func (e *HTTPStatusError) Error() string {
	return fmt.Sprintf("aeo: fetch failed: HTTP %d from %s", e.Status, e.URL)
}

// IsHTTPStatusError reports whether the error is an HTTPStatusError.
func IsHTTPStatusError(err error) bool {
	var target *HTTPStatusError
	return errors.As(err, &target)
}
