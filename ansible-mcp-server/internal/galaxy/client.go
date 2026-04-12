package galaxy

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashi-demo-lab/ansible-mcp-server/internal/hub"
	"github.com/hashi-demo-lab/ansible-mcp-server/internal/httpclient"
)

// Client queries the Ansible Galaxy v3 API.
// Galaxy is used as a fallback when Hub does not have full module docs.
type Client struct {
	client  *httpclient.Client
	baseURL string
}

// NewClient creates a new Galaxy API client.
func NewClient(client *httpclient.Client, baseURL string) *Client {
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	return &Client{client: client, baseURL: baseURL}
}

// --- wire types ---

type collectionListResponse struct {
	Data []struct {
		Namespace struct {
			Name string `json:"name"`
		} `json:"namespace"`
		Name          string `json:"name"`
		HighestVersion struct {
			Version string `json:"version"`
			Href    string `json:"href"`
		} `json:"highest_version"`
		LatestVersion struct {
			Version string `json:"version"`
			Href    string `json:"href"`
		} `json:"latest_version"`
		Description string `json:"description"`
		DownloadCount int  `json:"download_count"`
	} `json:"data"`
	Links struct {
		Next string `json:"next"`
	} `json:"links"`
	Meta struct {
		Count int `json:"count"`
	} `json:"meta"`
}

type roleListResponse struct {
	Results []struct {
		ID          int    `json:"id"`
		Namespace   struct {
			Name string `json:"name"`
		} `json:"namespace"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Platforms   []struct {
			Name     string `json:"name"`
			Release  string `json:"release"`
		} `json:"platforms"`
		GithubUser string `json:"github_user"`
		GithubRepo string `json:"github_repo"`
	} `json:"results"`
	Count int `json:"count"`
}

type roleDetailResponse struct {
	ID          int    `json:"id"`
	Namespace   struct {
		Name string `json:"name"`
	} `json:"namespace"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Platforms   []struct {
		Name    string `json:"name"`
		Release string `json:"release"`
	} `json:"platforms"`
	GithubUser  string `json:"github_user"`
	GithubRepo  string `json:"github_repo"`
}

type docsBlobResponse struct {
	DocsBlob struct {
		Contents []struct {
			ContentType string `json:"content_type"`
			ContentName string `json:"content_name"`
			DocStrings  struct {
				Doc struct {
					ShortDescription string                     `json:"short_description"`
					Description      interface{}                `json:"description"`
					Options          map[string]parameterWire   `json:"options"`
					Examples         string                     `json:"examples"`
					ReturnValues     map[string]returnValueWire `json:"return_values"`
				} `json:"doc"`
			} `json:"doc_strings"`
		} `json:"contents"`
	} `json:"docs_blob"`
}

type parameterWire struct {
	Description interface{}              `json:"description"`
	Type        string                   `json:"type"`
	Required    bool                     `json:"required"`
	Default     interface{}              `json:"default"`
	Choices     []interface{}            `json:"choices"`
	Aliases     []string                 `json:"aliases"`
	Suboptions  map[string]parameterWire `json:"suboptions"`
}

type returnValueWire struct {
	Description interface{} `json:"description"`
	Returned    string      `json:"returned"`
	Type        string      `json:"type"`
	Sample      interface{} `json:"sample"`
}

// SearchCollections searches Galaxy for collections matching a keyword.
func (c *Client) SearchCollections(ctx context.Context, query, namespace string, limit int) ([]hub.CollectionSearchResult, error) {
	params := url.Values{}
	params.Set("keywords", query)
	params.Set("limit", strconv.Itoa(limit))
	if namespace != "" {
		params.Set("namespace", namespace)
	}

	apiURL := c.baseURL + "collections/?" + params.Encode()
	body, err := c.get(ctx, apiURL)
	if err != nil {
		return nil, err
	}

	var resp collectionListResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("[galaxy] parse collections response: %w", err)
	}

	results := make([]hub.CollectionSearchResult, 0, len(resp.Data))
	for _, item := range resp.Data {
		version := item.HighestVersion.Version
		if version == "" {
			version = item.LatestVersion.Version
		}
		results = append(results, hub.CollectionSearchResult{
			Namespace:    item.Namespace.Name,
			Name:         item.Name,
			Version:      version,
			Description:  item.Description,
			SupportLevel: "community",
			HubURL:       c.baseURL + fmt.Sprintf("collections/%s/%s/", item.Namespace.Name, item.Name),
			Source:       "galaxy",
		})
	}
	return results, nil
}

// GetModuleDocs retrieves module parameter documentation from the Galaxy docs-blob.
// This is the fallback path when docs.ansible.com does not have the module page.
func (c *Client) GetModuleDocs(ctx context.Context, namespace, collection, module string) (*ModuleDocsResult, error) {
	// First get the collection to find the latest version
	collURL := c.baseURL + fmt.Sprintf("collections/%s/%s/", namespace, collection)
	collBody, err := c.get(ctx, collURL)
	if err != nil {
		return nil, err
	}

	var collResp struct {
		HighestVersion struct {
			Version string `json:"version"`
		} `json:"highest_version"`
		LatestVersion struct {
			Version string `json:"version"`
		} `json:"latest_version"`
	}
	if err := json.Unmarshal(collBody, &collResp); err != nil {
		return nil, fmt.Errorf("[galaxy] parse collection: %w", err)
	}

	version := collResp.HighestVersion.Version
	if version == "" {
		version = collResp.LatestVersion.Version
	}
	if version == "" {
		return nil, fmt.Errorf("[galaxy] could not determine version for %s.%s", namespace, collection)
	}

	// Fetch the docs-blob
	blobURL := c.baseURL + fmt.Sprintf("collections/%s/%s/versions/%s/docs-blob/", namespace, collection, version)
	blobBody, err := c.get(ctx, blobURL)
	if err != nil {
		return nil, err
	}

	var blob docsBlobResponse
	if err := json.Unmarshal(blobBody, &blob); err != nil {
		return nil, fmt.Errorf("[galaxy] parse docs-blob: %w", err)
	}

	// Find the specific module in the blob
	for _, content := range blob.DocsBlob.Contents {
		if content.ContentType == "module" && content.ContentName == module {
			return convertBlobToModuleDocs(content.DocStrings.Doc.ShortDescription,
				descToStrings(content.DocStrings.Doc.Description),
				content.DocStrings.Doc.Options,
				content.DocStrings.Doc.Examples,
				content.DocStrings.Doc.ReturnValues,
				namespace, collection, module,
			), nil
		}
	}

	return nil, fmt.Errorf("[galaxy] module %s.%s.%s not found in collection docs-blob", namespace, collection, module)
}

// SearchRoles searches Galaxy for roles matching a keyword.
func (c *Client) SearchRoles(ctx context.Context, query, namespace string, limit int) ([]hub.RoleSearchResult, error) {
	// Galaxy roles use the legacy API format
	params := url.Values{}
	params.Set("keywords", query)
	params.Set("page_size", strconv.Itoa(limit))
	if namespace != "" {
		params.Set("namespace", namespace)
	}

	// Galaxy roles are at /api/v1/roles/ (legacy) but we'll use v3 format
	// Actually Galaxy v3 has roles at /api/v3/roles/
	apiURL := c.baseURL + "roles/?" + params.Encode()
	body, err := c.get(ctx, apiURL)
	if err != nil {
		return nil, err
	}

	var resp roleListResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("[galaxy] parse roles response: %w", err)
	}

	results := make([]hub.RoleSearchResult, 0, len(resp.Results))
	for _, item := range resp.Results {
		results = append(results, hub.RoleSearchResult{
			Namespace:   item.Namespace.Name,
			Name:        item.Name,
			Description: item.Description,
			Source:      "galaxy",
			URL:         fmt.Sprintf("https://galaxy.ansible.com/%s/%s", item.GithubUser, item.Name),
		})
	}
	return results, nil
}

// GetRoleDetails fetches full details for a role from Galaxy.
func (c *Client) GetRoleDetails(ctx context.Context, namespace, name string) (*hub.RoleDetails, error) {
	// Search for the role first to get its ID
	params := url.Values{}
	params.Set("namespace", namespace)
	params.Set("name", name)
	params.Set("page_size", "1")

	apiURL := c.baseURL + "roles/?" + params.Encode()
	body, err := c.get(ctx, apiURL)
	if err != nil {
		return nil, err
	}

	var resp roleListResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("[galaxy] parse role search: %w", err)
	}
	if len(resp.Results) == 0 {
		return nil, fmt.Errorf("[galaxy] role %s.%s not found", namespace, name)
	}

	item := resp.Results[0]
	details := &hub.RoleDetails{
		Namespace:   item.Namespace.Name,
		Name:        item.Name,
		Description: item.Description,
		Source:      "galaxy",
		URL:         fmt.Sprintf("https://galaxy.ansible.com/%s/%s", item.GithubUser, item.Name),
	}

	for _, p := range item.Platforms {
		// Group platforms by OS name
		found := false
		for i, existing := range details.Platforms {
			if existing.Name == p.Name {
				details.Platforms[i].Versions = append(details.Platforms[i].Versions, p.Release)
				found = true
				break
			}
		}
		if !found {
			details.Platforms = append(details.Platforms, hub.PlatformEntry{
				Name:     p.Name,
				Versions: []string{p.Release},
			})
		}
	}

	return details, nil
}

func (c *Client) get(ctx context.Context, apiURL string) ([]byte, error) {
	resp, err := c.client.Get(ctx, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("[galaxy] request %s: %w", apiURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[galaxy] request %s: HTTP %d", apiURL, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[galaxy] reading response: %w", err)
	}
	return body, nil
}

// ModuleDocsResult holds module docs retrieved from Galaxy.
type ModuleDocsResult struct {
	FQCN             string            `json:"fqcn"`
	ShortDescription string            `json:"short_description"`
	Description      []string          `json:"description,omitempty"`
	Parameters       []ParameterResult `json:"parameters,omitempty"`
	Examples         string            `json:"examples,omitempty"`
	ReturnValues     []ReturnValResult `json:"return_values,omitempty"`
	Source           string            `json:"source"`
}

// ParameterResult is a module parameter from the Galaxy docs-blob.
type ParameterResult struct {
	Name        string            `json:"name"`
	Type        string            `json:"type,omitempty"`
	Required    bool              `json:"required"`
	Default     string            `json:"default,omitempty"`
	Choices     []string          `json:"choices,omitempty"`
	Description string            `json:"description"`
	Aliases     []string          `json:"aliases,omitempty"`
	Suboptions  []ParameterResult `json:"suboptions,omitempty"`
}

// ReturnValResult is a return value from the Galaxy docs-blob.
type ReturnValResult struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Returned    string `json:"returned,omitempty"`
	Type        string `json:"type,omitempty"`
}

func convertBlobToModuleDocs(shortDesc string, desc []string, options map[string]parameterWire, examples string, returnVals map[string]returnValueWire, namespace, collection, module string) *ModuleDocsResult {
	result := &ModuleDocsResult{
		FQCN:             fmt.Sprintf("%s.%s.%s", namespace, collection, module),
		ShortDescription: shortDesc,
		Description:      desc,
		Examples:         examples,
		Source:           "galaxy-fallback",
	}

	for name, opt := range options {
		p := ParameterResult{
			Name:        name,
			Type:        opt.Type,
			Required:    opt.Required,
			Description: descToString(opt.Description),
			Aliases:     opt.Aliases,
		}
		if opt.Default != nil {
			p.Default = fmt.Sprintf("%v", opt.Default)
		}
		for _, choice := range opt.Choices {
			p.Choices = append(p.Choices, fmt.Sprintf("%v", choice))
		}
		for subName, subOpt := range opt.Suboptions {
			sub := ParameterResult{
				Name:        subName,
				Type:        subOpt.Type,
				Required:    subOpt.Required,
				Description: descToString(subOpt.Description),
			}
			p.Suboptions = append(p.Suboptions, sub)
		}
		result.Parameters = append(result.Parameters, p)
	}

	for name, rv := range returnVals {
		result.ReturnValues = append(result.ReturnValues, ReturnValResult{
			Name:        name,
			Description: descToString(rv.Description),
			Returned:    rv.Returned,
			Type:        rv.Type,
		})
	}

	return result
}

// descToString converts a description that may be a string or []interface{} to a single string.
func descToString(v interface{}) string {
	switch d := v.(type) {
	case string:
		return d
	case []interface{}:
		parts := make([]string, 0, len(d))
		for _, p := range d {
			parts = append(parts, fmt.Sprintf("%v", p))
		}
		return strings.Join(parts, " ")
	default:
		if v != nil {
			return fmt.Sprintf("%v", v)
		}
		return ""
	}
}

// descToStrings converts a description that may be a string or []interface{} to a []string.
func descToStrings(v interface{}) []string {
	switch d := v.(type) {
	case string:
		if d == "" {
			return nil
		}
		return []string{d}
	case []interface{}:
		parts := make([]string, 0, len(d))
		for _, p := range d {
			parts = append(parts, fmt.Sprintf("%v", p))
		}
		return parts
	default:
		return nil
	}
}
