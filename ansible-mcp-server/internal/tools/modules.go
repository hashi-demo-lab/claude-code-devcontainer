package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashi-demo-lab/ansible-mcp-server/internal/docs"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerModuleTools(s *server.MCPServer, d *Dependencies) {
	s.AddTool(
		mcp.NewTool("search_modules",
			mcp.WithDescription("Search for Ansible modules matching a task description. Uses docs.ansible.com index to find modules across all collections."),
			mcp.WithString("query",
				mcp.Required(),
				mcp.Description("Task description or keyword (e.g., 'manage firewall rules', 'copy files', 'install packages')"),
			),
			mcp.WithString("namespace",
				mcp.Description("Filter by namespace (e.g., ansible, community)"),
			),
			mcp.WithString("collection",
				mcp.Description("Filter by collection name (e.g., posix, general)"),
			),
			mcp.WithNumber("limit",
				mcp.Description("Maximum results to return (default: 10, max: 50)"),
			),
		),
		handleSearchModules(d),
	)

	s.AddTool(
		mcp.NewTool("get_module_docs",
			mcp.WithDescription("Get full documentation for a specific Ansible module including all parameters, examples, and return values. Fetches from docs.ansible.com with Galaxy fallback for community modules."),
			mcp.WithString("fqcn",
				mcp.Description("Fully qualified collection name (e.g., ansible.posix.firewalld, community.general.git). Use this OR the separate namespace/collection/module fields."),
			),
			mcp.WithString("namespace",
				mcp.Description("Module namespace (e.g., ansible, community)"),
			),
			mcp.WithString("collection",
				mcp.Description("Collection name (e.g., posix, general)"),
			),
			mcp.WithString("module",
				mcp.Description("Module name (e.g., firewalld, git)"),
			),
		),
		handleGetModuleDocs(d),
	)
}

func handleSearchModules(d *Dependencies) Handler {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query := argString(req, "query")
		if query == "" {
			return errResult("ERROR: query is required")
		}
		namespace := argString(req, "namespace")
		collection := argString(req, "collection")
		limit := argInt(req, "limit", 10)
		if limit > 50 {
			limit = 50
		}

		// Fetch the all_plugins index to search across modules
		htmlContent, err := d.DocsFetcher.FetchPage(ctx, "collections/all_plugins.html")
		if err != nil {
			return errResult(fmt.Sprintf("ERROR [docs]: %s", err))
		}

		results := searchModulesInIndex(htmlContent, query, namespace, collection, limit, d.DocsBaseURL)

		out, _ := json.MarshalIndent(map[string]interface{}{
			"modules": results,
			"count":   len(results),
			"query":   query,
		}, "", "  ")
		return textResult(string(out))
	}
}

// moduleResult is returned from search_modules.
type moduleResult struct {
	FQCN             string `json:"fqcn"`
	Namespace        string `json:"namespace"`
	Collection       string `json:"collection"`
	ModuleName       string `json:"module_name"`
	ShortDescription string `json:"short_description"`
	DocsURL          string `json:"docs_url"`
	Source           string `json:"source"`
}

// searchModulesInIndex parses all_plugins.html and returns modules matching the query.
func searchModulesInIndex(htmlContent, query, namespace, collection string, limit int, docsBaseURL string) []moduleResult {
	// all_plugins.html contains a table of modules with FQCN and descriptions
	// We parse links with href matching the module URL pattern
	lowerQuery := strings.ToLower(query)
	lowerNS := strings.ToLower(namespace)
	lowerColl := strings.ToLower(collection)

	var results []moduleResult

	// Simple text-based search through the HTML for module references
	// The all_plugins page has lines like: namespace.collection.module_name - description
	lines := strings.Split(htmlContent, "\n")
	for _, line := range lines {
		if len(results) >= limit {
			break
		}

		// Look for links to _module.html pages
		if !strings.Contains(line, "_module.html") {
			continue
		}

		// Extract href
		hrefStart := strings.Index(line, `href="`)
		if hrefStart == -1 {
			continue
		}
		hrefStart += 6
		hrefEnd := strings.Index(line[hrefStart:], `"`)
		if hrefEnd == -1 {
			continue
		}
		href := line[hrefStart : hrefStart+hrefEnd]

		// Parse namespace/collection/module from path like:
		// ../ansible/posix/firewalld_module.html or
		// collections/ansible/posix/firewalld_module.html
		parts := strings.Split(strings.TrimPrefix(href, "../"), "/")
		if len(parts) < 3 {
			continue
		}

		var ns, coll, mod string
		// Check for collections/ns/coll/mod pattern
		start := 0
		if parts[0] == "collections" {
			start = 1
		}
		if len(parts) < start+3 {
			continue
		}
		ns = parts[start]
		coll = parts[start+1]
		modFile := parts[start+2]
		mod = strings.TrimSuffix(modFile, "_module.html")

		if ns == "" || coll == "" || mod == "" {
			continue
		}

		// Apply namespace/collection filters
		if lowerNS != "" && strings.ToLower(ns) != lowerNS {
			continue
		}
		if lowerColl != "" && strings.ToLower(coll) != lowerColl {
			continue
		}

		fqcn := fmt.Sprintf("%s.%s.%s", ns, coll, mod)

		// Extract description from the link text or surrounding text
		desc := extractLinkDescription(line)

		// Score against query
		lowerFQCN := strings.ToLower(fqcn)
		lowerDesc := strings.ToLower(desc)
		if !strings.Contains(lowerFQCN, lowerQuery) &&
			!strings.Contains(lowerDesc, lowerQuery) &&
			!strings.Contains(lowerQuery, strings.ToLower(mod)) {
			continue
		}

		results = append(results, moduleResult{
			FQCN:             fqcn,
			Namespace:        ns,
			Collection:       coll,
			ModuleName:       mod,
			ShortDescription: desc,
			DocsURL:          docsBaseURL + fmt.Sprintf("collections/%s/%s/%s_module.html", ns, coll, mod),
			Source:           "docs",
		})
	}

	return results
}

func extractLinkDescription(line string) string {
	// Extract text content between > and < tags
	var desc strings.Builder
	inTag := false
	for _, ch := range line {
		switch ch {
		case '<':
			inTag = true
		case '>':
			inTag = false
		default:
			if !inTag {
				desc.WriteRune(ch)
			}
		}
	}
	result := strings.TrimSpace(desc.String())
	// Remove leading separators
	result = strings.TrimLeft(result, " \t–—-")
	return strings.TrimSpace(result)
}

func handleGetModuleDocs(d *Dependencies) Handler {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Accept either fqcn or separate namespace/collection/module fields
		fqcn := argString(req, "fqcn")
		namespace := argString(req, "namespace")
		collection := argString(req, "collection")
		module := argString(req, "module")

		// Parse FQCN if provided
		if fqcn != "" {
			parts := strings.SplitN(fqcn, ".", 3)
			if len(parts) != 3 {
				return errResult(fmt.Sprintf("ERROR: invalid fqcn %q: expected namespace.collection.module", fqcn))
			}
			namespace = parts[0]
			collection = parts[1]
			module = parts[2]
		}

		if namespace == "" || collection == "" || module == "" {
			return errResult("ERROR: provide either fqcn or all of namespace, collection, and module")
		}

		docsURL := d.DocsFetcher.ModuleURL(namespace, collection, module)

		// Primary path: docs.ansible.com
		htmlContent, err := d.DocsFetcher.FetchPage(ctx, fmt.Sprintf("collections/%s/%s/%s_module.html", namespace, collection, module))
		if err == nil {
			moduleDocs, parseErr := docs.ParseModuleDocs(htmlContent)
			if parseErr != nil {
				return errResult(fmt.Sprintf("ERROR [docs]: parsing module page: %s", parseErr))
			}
			moduleDocs.FQCN = fmt.Sprintf("%s.%s.%s", namespace, collection, module)
			moduleDocs.DocsURL = docsURL
			moduleDocs.Source = "docs"

			out, _ := json.MarshalIndent(moduleDocs, "", "  ")
			return textResult(string(out))
		}

		// Check if it's a 404 (module not on docs.ansible.com)
		if !docs.IsNotFound(err) {
			return errResult(fmt.Sprintf("ERROR [docs]: %s", err))
		}

		// Fallback: Galaxy docs-blob
		galaxyDocs, galaxyErr := d.Galaxy.GetModuleDocs(ctx, namespace, collection, module)
		if galaxyErr != nil {
			return errResult(fmt.Sprintf("ERROR [galaxy-fallback]: module %s.%s.%s not found on docs.ansible.com or Galaxy. docs error: %s. galaxy error: %s",
				namespace, collection, module, err, galaxyErr))
		}

		galaxyDocs.FQCN = fmt.Sprintf("%s.%s.%s", namespace, collection, module)
		galaxyDocs.Source = "galaxy-fallback"

		// Combine into unified output format
		out, _ := json.MarshalIndent(map[string]interface{}{
			"fqcn":              galaxyDocs.FQCN,
			"short_description": galaxyDocs.ShortDescription,
			"description":       galaxyDocs.Description,
			"parameters":        galaxyDocs.Parameters,
			"examples":          galaxyDocs.Examples,
			"return_values":     galaxyDocs.ReturnValues,
			"docs_url":          docsURL,
			"source":            "galaxy-fallback",
		}, "", "  ")
		return textResult(string(out))
	}
}
