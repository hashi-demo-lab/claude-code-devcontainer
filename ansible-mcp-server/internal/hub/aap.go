package hub

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashi-demo-lab/ansible-mcp-server/internal/httpclient"
)

// AAPClient queries an on-premises Ansible Automation Platform Hub API.
// The AAP Hub mirrors the SaaS Hub v3 API under /api/galaxy/v3/.
type AAPClient struct {
	client   *httpclient.Client
	baseURL  string
	authMode string // "token" or "basic"
	token    string
	username string
	password string
}

// NewAAPClient creates a new on-prem AAP Hub client.
func NewAAPClient(client *httpclient.Client, aapURL, authMode, token, username, password string) *AAPClient {
	baseURL := strings.TrimSuffix(aapURL, "/") + "/api/galaxy/v3/"
	return &AAPClient{
		client:   client,
		baseURL:  baseURL,
		authMode: authMode,
		token:    token,
		username: username,
		password: password,
	}
}

// SearchCollections searches the AAP Hub for collections matching a keyword.
func (c *AAPClient) SearchCollections(ctx context.Context, query, namespace string, limit int, certifiedOnly bool) ([]CollectionSearchResult, error) {
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
		return nil, fmt.Errorf("[hub-aap] parse collections response: %w", err)
	}

	results := make([]CollectionSearchResult, 0, len(resp.Data))
	for _, item := range resp.Data {
		version := item.LatestVersion.Version
		if version == "" {
			version = item.HighestVersion.Version
		}
		results = append(results, CollectionSearchResult{
			Namespace:    item.Namespace.Name,
			Name:         item.Name,
			Version:      version,
			Description:  item.LatestVersion.Metadata.Description,
			SupportLevel: "certified",
			HubURL:       c.baseURL + fmt.Sprintf("collections/%s/%s/", item.Namespace.Name, item.Name),
			Source:       "aap",
		})
	}
	return results, nil
}

// GetCollectionDetails fetches full metadata for a collection from AAP Hub.
func (c *AAPClient) GetCollectionDetails(ctx context.Context, namespace, name, version string) (*CollectionDetails, error) {
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

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("[hub-aap] parse collection details: %w", err)
	}

	details := &CollectionDetails{
		Namespace: namespace,
		Name:      name,
		Source:    "aap",
	}

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

	return details, nil
}

// SearchRoles searches the AAP Hub for roles.
func (c *AAPClient) SearchRoles(ctx context.Context, query, namespace string, limit int) ([]RoleSearchResult, error) {
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
		return nil, fmt.Errorf("[hub-aap] parse roles response: %w", err)
	}

	results := make([]RoleSearchResult, 0, len(resp.Data))
	for _, item := range resp.Data {
		results = append(results, RoleSearchResult{
			Namespace:   item.Namespace.Name,
			Name:        item.Name,
			Description: item.Description,
			Source:      "aap",
			URL:         item.HRef,
		})
	}
	return results, nil
}

func (c *AAPClient) get(ctx context.Context, apiURL string) ([]byte, error) {
	headers := make(map[string]string)
	switch c.authMode {
	case "basic":
		creds := base64.StdEncoding.EncodeToString([]byte(c.username + ":" + c.password))
		headers["Authorization"] = "Basic " + creds
	default: // "token"
		headers["Authorization"] = "Bearer " + c.token
	}

	resp, err := c.client.Get(ctx, apiURL, headers)
	if err != nil {
		return nil, fmt.Errorf("[hub-aap] request %s: %w", apiURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("[hub-aap] authentication failed: check ANSIBLE_AAP_TOKEN or credentials")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[hub-aap] request %s: HTTP %d", apiURL, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[hub-aap] reading response: %w", err)
	}
	return body, nil
}
