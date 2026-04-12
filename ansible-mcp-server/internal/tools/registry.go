package tools

import (
	"context"

	"github.com/hashi-demo-lab/ansible-mcp-server/internal/docs"
	"github.com/hashi-demo-lab/ansible-mcp-server/internal/galaxy"
	"github.com/hashi-demo-lab/ansible-mcp-server/internal/hub"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Dependencies holds all external clients available to tool handlers.
type Dependencies struct {
	DocsBaseURL string
	DocsFetcher *docs.Fetcher
	SaaSHub     *hub.SaaSClient // nil when ANSIBLE_HUB_TARGET != saas|both
	AAPHub      *hub.AAPClient  // nil when ANSIBLE_HUB_TARGET != aap|both
	Galaxy      *galaxy.Client
	HubTarget   string // "saas" | "aap" | "both"
}

// RegisterAll registers all MCP tools with the server.
func RegisterAll(s *server.MCPServer, d *Dependencies) {
	registerBestPracticesTools(s, d)
	registerCollectionTools(s, d)
	registerModuleTools(s, d)
	registerRoleTools(s, d)
	registerScaffoldTools(s, d)
	registerValidateTools(s, d)
	registerTestCaseTools(s, d)
}

// --- argument helpers ---

// argString extracts a string argument from a tool request (returns "" if absent).
func argString(req mcp.CallToolRequest, key string) string {
	v, _ := req.Params.Arguments[key].(string)
	return v
}

// argInt extracts an integer argument; JSON numbers decode as float64.
func argInt(req mcp.CallToolRequest, key string, def int) int {
	if v, ok := req.Params.Arguments[key].(float64); ok {
		return int(v)
	}
	return def
}

// argBool extracts a boolean argument.
func argBool(req mcp.CallToolRequest, key string, def bool) bool {
	if v, ok := req.Params.Arguments[key].(bool); ok {
		return v
	}
	return def
}

// argStringSlice extracts a []string argument.
func argStringSlice(req mcp.CallToolRequest, key string) []string {
	raw, ok := req.Params.Arguments[key].([]interface{})
	if !ok {
		return nil
	}
	result := make([]string, 0, len(raw))
	for _, v := range raw {
		if s, ok := v.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

// errResult returns a structured MCP error result.
func errResult(msg string) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultError(msg), nil
}

// textResult returns a successful MCP text result.
func textResult(text string) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText(text), nil
}

// Handler type alias for readability.
type Handler = func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error)
