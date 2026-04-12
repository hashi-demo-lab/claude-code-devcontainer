package hub

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashi-demo-lab/ansible-mcp-server/internal/httpclient"
)

// SaaSClient queries the Automation Hub SaaS API.
type SaaSClient struct {
	client  *httpclient.Client
	baseURL string
	token   string
}

// NewSaaSClient creates a new Automation Hub SaaS client.
func NewSaaSClient(client *httpclient.Client, baseURL, token string) *SaaSClient {
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	return &SaaSClient{client: client, baseURL: baseURL, token: token}
}

// CollectionSearchResult is a collection returned from search.
type CollectionSearchResult struct {
	Namespace    string `json:"namespace"`
	Name         string `json:"name"`
	Version      string `json:"version"`
	Description  string `json:"description"`
	SupportLevel string `json:"support_level"`
	HubURL       string `json:"hub_url"`
	DocsURL      string `json:"docs_url"`
	Source       string `json:"source"`
}

// CollectionDetails contains full metadata for a collection.
type CollectionDetails struct {
	Namespace    string              `json:"namespace"`
	Name         string              `json:"name"`
	Version      string              `json:"version"`
	Description  string              `json:"description"`
	License      []string            `json:"license,omitempty"`
	Authors      []string            `json:"authors,omitempty"`
	Dependencies map[string]string   `json:"dependencies,omitempty"`
	Modules      []ContentSummary    `json:"modules,omitempty"`
	Roles        []ContentSummary    `json:"roles,omitempty"`
	Plugins      []ContentSummary    `json:"plugins,omitempty"`
	DocsURL      string              `json:"docs_url"`
	Source       string              `json:"source"`
}

// ContentSummary is a brief entry in a collection content list.
type ContentSummary struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
}

// RoleSearchResult is a role returned from search.
type RoleSearchResult struct {
	Namespace   string   `json:"namespace"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Platforms   []string `json:"platforms,omitempty"`
	Source      string   `json:"source"`
	URL         string   `json:"url"`
}

// RoleDetails contains full details for a role.
type RoleDetails struct {
	Namespace    string            `json:"namespace"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Version      string            `json:"version,omitempty"`
	Platforms    []PlatformEntry   `json:"platforms,omitempty"`
	Variables    []RoleVariable    `json:"variables,omitempty"`
	Dependencies []string          `json:"dependencies,omitempty"`
	Examples     []string          `json:"examples,omitempty"`
	Readme       string            `json:"readme,omitempty"`
	Source       string            `json:"source"`
	URL          string            `json:"url"`
}

// PlatformEntry describes a supported OS/version combination.
type PlatformEntry struct {
	Name     string   `json:"name"`
	Versions []string `json:"versions,omitempty"`
}

// RoleVariable describes a role default variable.
type RoleVariable struct {
	Name        string `json:"name"`
	Default     string `json:"default,omitempty"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
}

// DocsBlobContent holds a single content item from a collection docs-blob.
type DocsBlobContent struct {
	ContentType string          `json:"content_type"`
	ContentName string          `json:"content_name"`
	DocStrings  json.RawMessage `json:"doc_strings,omitempty"`
}

// --- wire types for API responses ---

type collectionListResponse struct {
	Data []struct {
		Namespace struct {
			Name string `json:"name"`
		} `json:"namespace"`
		Name          string `json:"name"`
		LatestVersion struct {
			Version  string `json:"version"`
			Metadata struct {
				Description  string   `json:"description"`
				Authors      []string `json:"authors"`
				License      []string `json:"license"`
				Dependencies map[string]string `json:"dependencies"`
			} `json:"metadata"`
			HRef string `json:"href"`
		} `json:"latest_version"`
		HighestVersion struct {
			Version string `json:"version"`
		} `json:"highest_version"`
	} `json:"data"`
	Meta struct {
		Count int `json:"count"`
	} `json:"meta"`
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
	Description interface{}            `json:"description"`
	Type        string                 `json:"type"`
	Required    bool                   `json:"required"`
	Default     interface{}            `json:"default"`
	Choices     []interface{}          `json:"choices"`
	Aliases     []string               `json:"aliases"`
	Suboptions  map[string]parameterWire `json:"suboptions"`
}

type returnValueWire struct {
	Description interface{} `json:"description"`
	Returned    string      `json:"returned"`
	Type        string      `json:"type"`
	Sample      interface{} `json:"sample"`
}

type roleListResponse struct {
	Data []struct {
		Namespace struct {
			Name string `json:"name"`
		} `json:"namespace"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Platforms   []struct {
			Name     string   `json:"name"`
			Versions []string `json:"versions"`
		} `json:"platforms"`
		HRef string `json:"href"`
	} `json:"data"`
	Meta struct {
		Count int `json:"count"`
	} `json:"meta"`
}

// SearchCollections searches the SaaS Hub for collections matching a keyword.
func (c *SaaSClient) SearchCollections(ctx context.Context, query, namespace string, limit int, certifiedOnly bool) ([]CollectionSearchResult, error) {
	params := url.Values{}
	params.Set("keywords", query)
	params.Set("limit", strconv.Itoa(limit))
	if namespace != "" {
		params.Set("namespace", namespace)
	}
	if certifiedOnly {
		params.Set("certification", "certified")
	}

	apiURL := c.baseURL + "collections/?" + params.Encode()
	body, err := c.get(ctx, apiURL)
	if err != nil {
		return nil, err
	}

	var resp collectionListResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("[hub-saas] parse collections response: %w", err)
	}

	results := make([]CollectionSearchResult, 0, len(resp.Data))
	for _, item := range resp.Data {
		version := item.LatestVersion.Version
		if version == "" {
			version = item.HighestVersion.Version
		}
		results = append(results, CollectionSearchResult{
			Namespace:   item.Namespace.Name,
			Name:        item.Name,
			Version:     version,
			Description: item.LatestVersion.Metadata.Description,
			SupportLevel: "certified",
			HubURL:      c.baseURL + fmt.Sprintf("collections/%s/%s/", item.Namespace.Name, item.Name),
			Source:      "saas",
		})
	}
	return results, nil
}

// GetCollectionDetails fetches full metadata for a specific collection.
func (c *SaaSClient) GetCollectionDetails(ctx context.Context, namespace, name, version string) (*CollectionDetails, error) {
	var apiURL string
	if version != "" {
		apiURL = c.baseURL + fmt.Sprintf("collections/%s/%s/versions/%s/", namespace, name, version)
	} else {
		apiURL = c.baseURL + fmt.Sprintf("collections/%s/%s/", namespace, name)
	}

	body, err := c.get(ctx, apiURL)
	if err != nil {
		return nil, err
	}

	// Parse the collection detail response
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("[hub-saas] parse collection details: %w", err)
	}

	details := &CollectionDetails{
		Namespace: namespace,
		Name:      name,
		Source:    "saas",
	}

	// Extract version
	if v, ok := raw["latest_version"]; ok {
		var lv struct {
			Version  string `json:"version"`
			Metadata struct {
				Description  string            `json:"description"`
				Authors      []string          `json:"authors"`
				License      []string          `json:"license"`
				Dependencies map[string]string `json:"dependencies"`
			} `json:"metadata"`
		}
		if err := json.Unmarshal(v, &lv); err == nil {
			details.Version = lv.Version
			details.Description = lv.Metadata.Description
			details.Authors = lv.Metadata.Authors
			details.License = lv.Metadata.License
			details.Dependencies = lv.Metadata.Dependencies
		}
	}

	// Fetch the docs-blob for content listing
	blobURL := c.baseURL + fmt.Sprintf("collections/%s/%s/versions/%s/docs-blob/", namespace, name, details.Version)
	blobBody, err := c.get(ctx, blobURL)
	if err == nil {
		var blob docsBlobResponse
		if err := json.Unmarshal(blobBody, &blob); err == nil {
			for _, content := range blob.DocsBlob.Contents {
				summary := ContentSummary{
					Name: content.ContentName,
					Type: content.ContentType,
				}
				switch content.ContentType {
				case "module":
					details.Modules = append(details.Modules, summary)
				case "role":
					details.Roles = append(details.Roles, summary)
				default:
					details.Plugins = append(details.Plugins, summary)
				}
			}
		}
	}

	return details, nil
}

// GetModuleDocsFromBlob fetches module parameter docs from the Hub docs-blob API.
func (c *SaaSClient) GetModuleDocsFromBlob(ctx context.Context, namespace, collection, module, version string) (*docsBlobResponse, error) {
	if version == "" {
		version = "latest"
	}
	blobURL := c.baseURL + fmt.Sprintf("collections/%s/%s/versions/%s/docs-blob/", namespace, collection, version)
	body, err := c.get(ctx, blobURL)
	if err != nil {
		return nil, err
	}
	var blob docsBlobResponse
	if err := json.Unmarshal(body, &blob); err != nil {
		return nil, fmt.Errorf("[hub-saas] parse docs-blob: %w", err)
	}
	return &blob, nil
}

// SearchRoles searches the SaaS Hub for roles matching a keyword.
func (c *SaaSClient) SearchRoles(ctx context.Context, query, namespace string, limit int) ([]RoleSearchResult, error) {
	params := url.Values{}
	params.Set("keywords", query)
	params.Set("limit", strconv.Itoa(limit))
	if namespace != "" {
		params.Set("namespace", namespace)
	}

	apiURL := c.baseURL + "roles/?" + params.Encode()
	body, err := c.get(ctx, apiURL)
	if err != nil {
		return nil, err
	}

	var resp roleListResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("[hub-saas] parse roles response: %w", err)
	}

	results := make([]RoleSearchResult, 0, len(resp.Data))
	for _, item := range resp.Data {
		results = append(results, RoleSearchResult{
			Namespace:   item.Namespace.Name,
			Name:        item.Name,
			Description: item.Description,
			Source:      "saas",
			URL:         item.HRef,
		})
	}
	return results, nil
}

// GetRoleDetails fetches full details for a role from Hub SaaS.
func (c *SaaSClient) GetRoleDetails(ctx context.Context, namespace, name string) (*RoleDetails, error) {
	apiURL := c.baseURL + fmt.Sprintf("roles/%s/%s/", namespace, name)
	body, err := c.get(ctx, apiURL)
	if err != nil {
		return nil, err
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("[hub-saas] parse role details: %w", err)
	}

	details := &RoleDetails{
		Namespace: namespace,
		Name:      name,
		Source:    "saas",
		URL:       apiURL,
	}

	if desc, ok := raw["description"]; ok {
		_ = json.Unmarshal(desc, &details.Description)
	}

	return details, nil
}

func (c *SaaSClient) get(ctx context.Context, apiURL string) ([]byte, error) {
	headers := map[string]string{
		"Authorization": "Bearer " + c.token,
	}
	resp, err := c.client.Get(ctx, apiURL, headers)
	if err != nil {
		return nil, fmt.Errorf("[hub-saas] request %s: %w", apiURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("[hub-saas] authentication failed: check ANSIBLE_HUB_SAAS_TOKEN")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[hub-saas] request %s: HTTP %d", apiURL, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[hub-saas] reading response: %w", err)
	}
	return body, nil
}
