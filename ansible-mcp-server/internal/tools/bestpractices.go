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

func registerBestPracticesTools(s *server.MCPServer, d *Dependencies) {
	s.AddTool(
		mcp.NewTool("get_playbook_keywords",
			mcp.WithDescription("Retrieve the Ansible playbook keyword reference from docs.ansible.com. Returns definitions, types, defaults, and scope for all play and task keywords."),
			mcp.WithString("keyword",
				mcp.Description("Filter to a specific keyword (e.g., become, delegate_to, when). Omit to get the full reference."),
			),
		),
		handleGetPlaybookKeywords(d),
	)

	s.AddTool(
		mcp.NewTool("get_special_variables",
			mcp.WithDescription("Retrieve Ansible special (magic) variables reference from docs.ansible.com. Includes inventory_hostname, ansible_facts, hostvars, and all other magic variables."),
			mcp.WithString("variable",
				mcp.Description("Filter to a specific variable (e.g., inventory_hostname, ansible_facts). Omit to get all special variables."),
			),
		),
		handleGetSpecialVariables(d),
	)

	s.AddTool(
		mcp.NewTool("get_best_practices",
			mcp.WithDescription("Retrieve Ansible best practices for a specific topic from docs.ansible.com. Topics include: idempotency, error handling, variables, handlers, tags, roles, collections, loops, conditionals, templates, vault, inventory, galaxy."),
			mcp.WithString("topic",
				mcp.Required(),
				mcp.Description("Topic to look up (e.g., idempotency, error handling, variable naming, handlers, tags, roles, vault, inventory)"),
			),
		),
		handleGetBestPractices(d),
	)
}

func handleGetPlaybookKeywords(d *Dependencies) Handler {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		keyword := argString(req, "keyword")

		htmlContent, err := d.DocsFetcher.FetchPage(ctx, "reference_appendices/playbooks_keywords.html")
		if err != nil {
			return errResult(fmt.Sprintf("ERROR [docs]: %s", err))
		}

		keywords, err := docs.ParsePlaybookKeywords(htmlContent)
		if err != nil {
			return errResult(fmt.Sprintf("ERROR [docs]: parsing playbook keywords: %s", err))
		}

		sourceURL := d.DocsFetcher.PlaybookKeywordsURL()

		if keyword != "" {
			lower := strings.ToLower(keyword)
			for _, kw := range keywords {
				if strings.ToLower(kw.Name) == lower {
					out, _ := json.MarshalIndent(map[string]interface{}{
						"keyword":    kw,
						"source_url": sourceURL,
					}, "", "  ")
					return textResult(string(out))
				}
			}
			return errResult(fmt.Sprintf("ERROR [docs]: keyword %q not found in playbook keywords reference", keyword))
		}

		out, _ := json.MarshalIndent(map[string]interface{}{
			"keywords":   keywords,
			"count":      len(keywords),
			"source_url": sourceURL,
		}, "", "  ")
		return textResult(string(out))
	}
}

func handleGetSpecialVariables(d *Dependencies) Handler {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		variable := argString(req, "variable")

		htmlContent, err := d.DocsFetcher.FetchPage(ctx, "reference_appendices/special_variables.html")
		if err != nil {
			return errResult(fmt.Sprintf("ERROR [docs]: %s", err))
		}

		vars, err := docs.ParseSpecialVariables(htmlContent)
		if err != nil {
			return errResult(fmt.Sprintf("ERROR [docs]: parsing special variables: %s", err))
		}

		sourceURL := d.DocsFetcher.SpecialVariablesURL()

		if variable != "" {
			lower := strings.ToLower(variable)
			for _, v := range vars {
				if strings.ToLower(v.Name) == lower {
					out, _ := json.MarshalIndent(map[string]interface{}{
						"variable":   v,
						"source_url": sourceURL,
					}, "", "  ")
					return textResult(string(out))
				}
			}
			return errResult(fmt.Sprintf("ERROR [docs]: variable %q not found in special variables reference", variable))
		}

		out, _ := json.MarshalIndent(map[string]interface{}{
			"variables":  vars,
			"count":      len(vars),
			"source_url": sourceURL,
		}, "", "  ")
		return textResult(string(out))
	}
}

// topicToPath maps topic keywords to docs.ansible.com relative paths.
var topicToPath = []struct {
	keywords []string
	path     string
}{
	{[]string{"playbook structure", "tasks", "plays", "playbook intro"}, "playbook_guide/playbooks_intro.html"},
	{[]string{"variables", "precedence", "variable naming"}, "playbook_guide/playbooks_variables.html"},
	{[]string{"handlers"}, "playbook_guide/playbooks_handlers.html"},
	{[]string{"tags"}, "playbook_guide/playbooks_tags.html"},
	{[]string{"error handling", "ignore_errors", "rescue", "errors"}, "playbook_guide/playbooks_error_handling.html"},
	{[]string{"idempotency", "state", "testing strategies"}, "reference_appendices/test_strategies.html"},
	{[]string{"roles", "role structure", "reuse roles"}, "playbook_guide/playbooks_reuse_roles.html"},
	{[]string{"collections", "using collections"}, "collections_guide/index.html"},
	{[]string{"loops", "with_items", "loop"}, "playbook_guide/playbooks_loops.html"},
	{[]string{"conditionals", "when"}, "playbook_guide/playbooks_conditionals.html"},
	{[]string{"templates", "jinja2", "jinja"}, "playbook_guide/playbooks_templating.html"},
	{[]string{"vault", "secrets", "encrypt"}, "vault_guide/index.html"},
	{[]string{"inventory"}, "inventory_guide/index.html"},
	{[]string{"galaxy", "install collections"}, "galaxy/user_guide.html"},
	{[]string{"yaml syntax", "yaml"}, "reference_appendices/YAMLSyntax.html"},
	{[]string{"config", "ansible.cfg"}, "reference_appendices/config.html"},
	{[]string{"return values", "return"}, "reference_appendices/common_return_values.html"},
}

func resolveTopicPath(topic string) (string, bool) {
	lower := strings.ToLower(strings.TrimSpace(topic))
	for _, entry := range topicToPath {
		for _, kw := range entry.keywords {
			if strings.Contains(lower, kw) || strings.Contains(kw, lower) {
				return entry.path, true
			}
		}
	}
	return "", false
}

func handleGetBestPractices(d *Dependencies) Handler {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		topic := argString(req, "topic")
		if topic == "" {
			return errResult("ERROR: topic is required")
		}

		path, found := resolveTopicPath(topic)
		if !found {
			// Fall back to a generic search on the testing strategies page
			path = "reference_appendices/test_strategies.html"
		}

		htmlContent, err := d.DocsFetcher.FetchPage(ctx, path)
		if err != nil {
			return errResult(fmt.Sprintf("ERROR [docs]: %s", err))
		}

		content, err := docs.ParseGenericContent(htmlContent)
		if err != nil {
			return errResult(fmt.Sprintf("ERROR [docs]: parsing page: %s", err))
		}

		sourceURL := d.DocsFetcher.BestPracticesURL(path)

		out, _ := json.MarshalIndent(map[string]interface{}{
			"topic":                       topic,
			"content":                     content,
			"source_url":                  sourceURL,
			"related_keywords_reference":  d.DocsFetcher.PlaybookKeywordsURL(),
		}, "", "  ")
		return textResult(string(out))
	}
}
