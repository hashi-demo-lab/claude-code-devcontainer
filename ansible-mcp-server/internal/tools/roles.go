package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashi-demo-lab/ansible-mcp-server/internal/hub"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerRoleTools(s *server.MCPServer, d *Dependencies) {
	s.AddTool(
		mcp.NewTool("search_roles",
			mcp.WithDescription("Search for Ansible roles in Automation Hub and/or Ansible Galaxy."),
			mcp.WithString("query",
				mcp.Required(),
				mcp.Description("Task description or keyword (e.g., nginx, postgresql, system hardening)"),
			),
			mcp.WithString("namespace",
				mcp.Description("Filter by namespace/author (e.g., geerlingguy, redhat)"),
			),
			mcp.WithNumber("limit",
				mcp.Description("Maximum results to return (default: 10)"),
			),
			mcp.WithString("source",
				mcp.Description("Where to search: hub, galaxy, or both (default: both)"),
			),
		),
		handleSearchRoles(d),
	)

	s.AddTool(
		mcp.NewTool("get_role_details",
			mcp.WithDescription("Get full details for a specific Ansible role including variables, platforms, dependencies, and examples."),
			mcp.WithString("namespace",
				mcp.Required(),
				mcp.Description("Role namespace or author (e.g., geerlingguy, redhat)"),
			),
			mcp.WithString("name",
				mcp.Required(),
				mcp.Description("Role name (e.g., apache, postgresql)"),
			),
			mcp.WithString("source",
				mcp.Description("Where to look: hub or galaxy (default: auto-detect from namespace)"),
			),
		),
		handleGetRoleDetails(d),
	)
}

func handleSearchRoles(d *Dependencies) Handler {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query := argString(req, "query")
		if query == "" {
			return errResult("ERROR: query is required")
		}
		namespace := argString(req, "namespace")
		limit := argInt(req, "limit", 10)
		source := argString(req, "source")
		if source == "" {
			source = "both"
		}

		var allResults []hub.RoleSearchResult

		switch source {
		case "hub":
			results, err := searchHubRoles(ctx, d, query, namespace, limit)
			if err != nil {
				return errResult(fmt.Sprintf("ERROR [hub]: %s", err))
			}
			allResults = results

		case "galaxy":
			results, err := d.Galaxy.SearchRoles(ctx, query, namespace, limit)
			if err != nil {
				return errResult(fmt.Sprintf("ERROR [galaxy]: %s", err))
			}
			allResults = results

		default: // "both"
			hubResults, _ := searchHubRoles(ctx, d, query, namespace, limit/2+limit%2)
			allResults = append(allResults, hubResults...)

			galaxyResults, _ := d.Galaxy.SearchRoles(ctx, query, namespace, limit/2)
			allResults = append(allResults, galaxyResults...)
		}

		out, _ := json.MarshalIndent(map[string]interface{}{
			"roles": allResults,
			"count": len(allResults),
			"query": query,
		}, "", "  ")
		return textResult(string(out))
	}
}

func handleGetRoleDetails(d *Dependencies) Handler {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		namespace := argString(req, "namespace")
		name := argString(req, "name")
		source := argString(req, "source")

		if namespace == "" {
			return errResult("ERROR: namespace is required")
		}
		if name == "" {
			return errResult("ERROR: name is required")
		}

		if source == "" {
			// Auto-detect: try Hub first, fall back to Galaxy
			source = "auto"
		}

		var details *hub.RoleDetails
		var err error

		switch source {
		case "hub":
			details, err = getHubRoleDetails(ctx, d, namespace, name)
			if err != nil {
				return errResult(fmt.Sprintf("ERROR [hub]: %s", err))
			}
		case "galaxy":
			details, err = d.Galaxy.GetRoleDetails(ctx, namespace, name)
			if err != nil {
				return errResult(fmt.Sprintf("ERROR [galaxy]: %s", err))
			}
		default: // "auto"
			details, err = getHubRoleDetails(ctx, d, namespace, name)
			if err != nil {
				// Fallback to Galaxy
				details, err = d.Galaxy.GetRoleDetails(ctx, namespace, name)
				if err != nil {
					return errResult(fmt.Sprintf("ERROR [galaxy]: %s", err))
				}
			}
		}

		out, _ := json.MarshalIndent(details, "", "  ")
		return textResult(string(out))
	}
}

func searchHubRoles(ctx context.Context, d *Dependencies, query, namespace string, limit int) ([]hub.RoleSearchResult, error) {
	switch d.HubTarget {
	case "saas":
		if d.SaaSHub == nil {
			return nil, fmt.Errorf("SaaS Hub not configured")
		}
		return d.SaaSHub.SearchRoles(ctx, query, namespace, limit)
	case "aap":
		if d.AAPHub == nil {
			return nil, fmt.Errorf("AAP Hub not configured")
		}
		return d.AAPHub.SearchRoles(ctx, query, namespace, limit)
	case "both":
		if d.SaaSHub != nil {
			results, err := d.SaaSHub.SearchRoles(ctx, query, namespace, limit)
			if err == nil && len(results) > 0 {
				return results, nil
			}
		}
		if d.AAPHub != nil {
			return d.AAPHub.SearchRoles(ctx, query, namespace, limit)
		}
	}
	return nil, nil
}

func getHubRoleDetails(ctx context.Context, d *Dependencies, namespace, name string) (*hub.RoleDetails, error) {
	switch d.HubTarget {
	case "saas", "both":
		if d.SaaSHub != nil {
			return d.SaaSHub.GetRoleDetails(ctx, namespace, name)
		}
	case "aap":
		if d.AAPHub != nil {
			// AAPClient doesn't have GetRoleDetails yet; return not found
			return nil, fmt.Errorf("role %s.%s not found on AAP Hub", namespace, name)
		}
	}
	return nil, fmt.Errorf("no Hub client configured for role lookup")
}
