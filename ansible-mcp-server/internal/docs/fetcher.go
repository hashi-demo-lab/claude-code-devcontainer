package docs

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashi-demo-lab/ansible-mcp-server/internal/httpclient"
)

// Fetcher retrieves HTML pages from docs.ansible.com.
type Fetcher struct {
	client  *httpclient.Client
	baseURL string
}

// NewFetcher creates a new docs fetcher.
func NewFetcher(client *httpclient.Client, baseURL string) *Fetcher {
	return &Fetcher{client: client, baseURL: baseURL}
}

// FetchPage fetches a docs page by relative path and returns the HTML body.
func (f *Fetcher) FetchPage(ctx context.Context, path string) (string, error) {
	url := f.baseURL + strings.TrimPrefix(path, "/")
	resp, err := f.client.Get(ctx, url, nil)
	if err != nil {
		return "", fmt.Errorf("[docs] fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", &ErrNotFound{URL: url}
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("[docs] fetch %s: HTTP %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("[docs] reading response body: %w", err)
	}
	return string(body), nil
}

// ModuleURL returns the canonical URL for a module documentation page.
func (f *Fetcher) ModuleURL(namespace, collection, module string) string {
	return f.baseURL + fmt.Sprintf("collections/%s/%s/%s_module.html", namespace, collection, module)
}

// CollectionURL returns the canonical URL for a collection index page.
func (f *Fetcher) CollectionURL(namespace, collection string) string {
	return f.baseURL + fmt.Sprintf("collections/%s/%s/index.html", namespace, collection)
}

// PlaybookKeywordsURL returns the URL for the playbook keywords reference page.
func (f *Fetcher) PlaybookKeywordsURL() string {
	return f.baseURL + "reference_appendices/playbooks_keywords.html"
}

// SpecialVariablesURL returns the URL for the special variables reference page.
func (f *Fetcher) SpecialVariablesURL() string {
	return f.baseURL + "reference_appendices/special_variables.html"
}

// BestPracticesURL returns the docs URL for a given relative path.
func (f *Fetcher) BestPracticesURL(relPath string) string {
	return f.baseURL + relPath
}

// ErrNotFound is returned when a docs page returns HTTP 404.
type ErrNotFound struct {
	URL string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("[docs] page not found: %s", e.URL)
}

// IsNotFound reports whether err is an ErrNotFound.
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*ErrNotFound)
	return ok
}
