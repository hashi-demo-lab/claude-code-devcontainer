package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashi-demo-lab/ansible-mcp-server/internal/hub"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerCollectionTools(s *server.MCPServer, d *Dependencies) {
	s.AddTool(
		mcp.NewTool("search_collections",
			mcp.WithDescription("Search Ansible Automation Hub (or Galaxy fallback) for collections matching a keyword or task description. Returns certified collections by default."),
			mcp.WithString("query",
				mcp.Required(),
				mcp.Description("Keyword or task description to search for"),
			),
			mcp.WithString("namespace",
				mcp.Description("Filter by namespace (e.g., ansible, community, redhat)"),
			),
			mcp.WithNumber("limit",
				mcp.Description("Maximum results to return (default: 10, max: 50)"),
			),
			mcp.WithBoolean("certified_only",
				mcp.Description("Return only Hub-certified content (default: true)"),
			),
		),
		handleSearchCollections(d),
	)

	s.AddTool(
		mcp.NewTool("get_collection_details",
			mcp.WithDescription("Get full metadata and content listing for a specific Ansible collection, including modules, roles, and plugins."),
			mcp.WithString("namespace",
				mcp.Required(),
				mcp.Description("Collection namespace (e.g., ansible, community, redhat)"),
			),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Collection name (e.g., posix, general, rhel_system_roles)"),
			),
			mcp.WithString("version",
				mcp.Description("Specific version to fetch (default: latest)"),
			),
		),
		handleGetCollectionDetails(d),
	)
}

func handleSearchCollections(d *Dependencies) Handler {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query := argString(req, "query")
		if query == "" {
			return errResult("ERROR: query is required")
		}
		namespace := argString(req, "namespace")
		limit := argInt(req, "limit", 10)
		if limit > 50 {
			limit = 50
		}
		certifiedOnly := argBool(req, "certified_only", true)

		results, source, err := searchCollectionsWithFallback(ctx, d, query, namespace, limit, certifiedOnly)
		if err != nil {
			return errResult(fmt.Sprintf("ERROR [%s]: %s", source, err))
		}

		out, _ := json.MarshalIndent(map[string]interface{}{
			"collections": results,
			"count":       len(results),
			"query":       query,
		}, "", "  ")
		return textResult(string(out))
	}
}

func handleGetCollectionDetails(d *Dependencies) Handler {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := argString(req, "namespace")
		name := argString(req, "name")
		version := argString(req, "version")

		if namespace == "" {
			return errResult("ERROR: namespace is required")
		}
		if name == "" {
			return errResult("ERROR: name is required")
		}

		details, source, err := getCollectionDetailsWithFallback(ctx, d, namespace, name, version)
		if err != nil {
			return errResult(fmt.Sprintf("ERROR [%s]: %s", source, err))
		}

		details.DocsURL = d.DocsFetcher.CollectionURL(namespace, name)
		out, _ := json.MarshalIndent(details, "", "  ")
		return textResult(string(out))
	}
}

// searchCollectionsWithFallback queries Hub per the configured target,
// falling back to the secondary Hub or Galaxy as needed.
func searchCollectionsWithFallback(ctx context.Context, d *Dependencies, query, namespace string, limit int, certifiedOnly bool) ([]hub.CollectionSearchResult, string, error) {
	switch d.HubTarget {
	case "saas":
		if d.SaaSHub == nil {
			return nil, "hub-saas", fmt.Errorf("SaaS Hub not configured: set ANSIBLE_HUB_SAAS_TOKEN")
		}
		results, err := d.SaaSHub.SearchCollections(ctx, query, namespace, limit, certifiedOnly)
		return results, "hub-saas", err

	case "aap":
		if d.AAPHub == nil {
			return nil, "hub-aap", fmt.Errorf("AAP Hub not configured: set ANSIBLE_AAP_URL and credentials")
		}
		results, err := d.AAPHub.SearchCollections(ctx, query, namespace, limit, certifiedOnly)
		return results, "hub-aap", err

	case "both":
		// SaaS is primary; AAP is fallback on error or empty results
		if d.SaaSHub != nil {
			results, err := d.SaaSHub.SearchCollections(ctx, query, namespace, limit, certifiedOnly)
			if err == nil && len(results) > 0 {
				return results, "hub-saas", nil
			}
		}
		if d.AAPHub != nil {
			results, err := d.AAPHub.SearchCollections(ctx, query, namespace, limit, certifiedOnly)
			if err == nil && len(results) > 0 {
				return results, "hub-aap", nil
			}
		}
		// Final fallback: Galaxy
		results, err := d.Galaxy.SearchCollections(ctx, query, namespace, limit)
		return results, "galaxy", err
	}

	return nil, "unknown", fmt.Errorf("unknown hub target %q", d.HubTarget)
}

func getCollectionDetailsWithFallback(ctx context.Context, d *Dependencies, namespace, name, version string) (*hub.CollectionDetails, string, error) {
	switch d.HubTarget {
	case "saas":
		if d.SaaSHub == nil {
			return nil, "hub-saas", fmt.Errorf("SaaS Hub not configured")
		}
		details, err := d.SaaSHub.GetCollectionDetails(ctx, namespace, name, version)
		return details, "hub-saas", err

	case "aap":
		if d.AAPHub == nil {
			return nil, "hub-aap", fmt.Errorf("AAP Hub not configured")
		}
		details, err := d.AAPHub.GetCollectionDetails(ctx, namespace, name, version)
		return details, "hub-aap", err

	case "both":
		if d.SaaSHub != nil {
			details, err := d.SaaSHub.GetCollectionDetails(ctx, namespace, name, version)
			if err == nil {
				return details, "hub-saas", nil
			}
		}
		if d.AAPHub != nil {
			details, err := d.AAPHub.GetCollectionDetails(ctx, namespace, name, version)
			if err == nil {
				return details, "hub-aap", nil
			}
		}
		return nil, "hub", fmt.Errorf("collection %s.%s not found in configured Hub targets", namespace, name)
	}

	return nil, "unknown", fmt.Errorf("unknown hub target %q", d.HubTarget)
}
