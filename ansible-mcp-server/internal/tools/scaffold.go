package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerScaffoldTools(s *server.MCPServer, _ *Dependencies) {
	s.AddTool(
		mcp.NewTool("generate_playbook_scaffold",
			mcp.WithDescription("Generate a boilerplate Ansible playbook following best practices: explicit state, task-level become, handlers for restarts, named tasks and plays, gather_facts enabled."),
			mcp.WithString("task_description",
				mcp.Required(),
				mcp.Description("What the playbook should accomplish (e.g., 'install and configure nginx on RHEL')"),
			),
			mcp.WithString("target_os",
				mcp.Description("Target OS family: rhel, debian, ubuntu, windows (default: rhel)"),
			),
			mcp.WithArray("collections",
				mcp.Description("Collections to declare in the playbook (e.g., [\"ansible.posix\", \"community.general\"])"),
			),
			mcp.WithBoolean("use_roles",
				mcp.Description("Structure as a role invocation instead of inline tasks (default: false)"),
			),
			mcp.WithString("style",
				mcp.Description("simple (single playbook file) or project (directory layout with roles/) (default: simple)"),
			),
		),
		handleGeneratePlaybookScaffold(),
	)
}

type scaffoldOutput struct {
	Playbook string            `json:"playbook,omitempty"`
	Files    map[string]string `json:"files,omitempty"`
	Notes    []string          `json:"notes"`
}

func handleGeneratePlaybookScaffold() Handler {
	return func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		taskDesc := argString(req, "task_description")
		if taskDesc == "" {
			return errResult("ERROR: task_description is required")
		}
		targetOS := argString(req, "target_os")
		if targetOS == "" {
			targetOS = "rhel"
		}
		collections := argStringSlice(req, "collections")
		useRoles := argBool(req, "use_roles", false)
		style := argString(req, "style")
		if style == "" {
			style = "simple"
		}

		var output scaffoldOutput

		if style == "project" {
			output.Files = generateProjectScaffold(taskDesc, targetOS, collections, useRoles)
			output.Notes = scaffoldNotes(useRoles, style)
		} else {
			output.Playbook = generateSimplePlaybook(taskDesc, targetOS, collections, useRoles)
			output.Notes = scaffoldNotes(useRoles, style)
		}

		out, _ := json.MarshalIndent(output, "", "  ")
		return textResult(string(out))
	}
}

func generateSimplePlaybook(taskDesc, targetOS string, collections []string, useRoles bool) string {
	var b strings.Builder

	b.WriteString("---\n")

	// Collections declaration
	if len(collections) > 0 {
		b.WriteString("# Declare collections used in this playbook\n")
		b.WriteString("collections:\n")
		for _, c := range collections {
			b.WriteString(fmt.Sprintf("  - %s\n", c))
		}
		b.WriteString("\n")
	}

	// Play header
	b.WriteString(fmt.Sprintf("- name: %s\n", sentenceCase(taskDesc)))
	b.WriteString("  hosts: all\n")
	b.WriteString("  gather_facts: true\n")
	b.WriteString("\n")

	// Variables
	b.WriteString("  vars:\n")
	b.WriteString(fmt.Sprintf("    # Target OS: %s\n", targetOS))
	b.WriteString("    # Add your variables here\n")
	b.WriteString("    # example_package: nginx\n")
	b.WriteString("\n")

	if useRoles {
		roleName := descToRoleName(taskDesc)
		b.WriteString("  roles:\n")
		b.WriteString(fmt.Sprintf("    - role: %s\n", roleName))
		return b.String()
	}

	// Handlers
	b.WriteString("  handlers:\n")
	b.WriteString("    - name: Restart service\n")
	b.WriteString("      ansible.builtin.service:\n")
	b.WriteString("        name: \"{{ service_name }}\"\n")
	b.WriteString("        state: restarted\n")
	b.WriteString("\n")

	// Tasks
	b.WriteString("  tasks:\n")
	b.WriteString(generateTasksForOS(taskDesc, targetOS))

	return b.String()
}

func generateTasksForOS(taskDesc, targetOS string) string {
	var b strings.Builder
	lower := strings.ToLower(taskDesc)

	// Generic tasks based on description keywords
	if strings.Contains(lower, "install") || strings.Contains(lower, "package") {
		b.WriteString("\n    - name: Ensure required packages are installed\n")
		if targetOS == "windows" {
			b.WriteString("      ansible.windows.win_package:\n")
			b.WriteString("        name: \"{{ package_name }}\"\n")
			b.WriteString("        state: present\n")
		} else {
			b.WriteString("      ansible.builtin.package:\n")
			b.WriteString("        name: \"{{ package_name }}\"\n")
			b.WriteString("        state: present\n")
		}
		b.WriteString("      become: true\n")
	}

	if strings.Contains(lower, "service") || strings.Contains(lower, "start") || strings.Contains(lower, "enable") {
		b.WriteString("\n    - name: Ensure service is started and enabled\n")
		if targetOS == "windows" {
			b.WriteString("      ansible.windows.win_service:\n")
			b.WriteString("        name: \"{{ service_name }}\"\n")
			b.WriteString("        state: started\n")
			b.WriteString("        start_mode: auto\n")
		} else {
			b.WriteString("      ansible.builtin.service:\n")
			b.WriteString("        name: \"{{ service_name }}\"\n")
			b.WriteString("        state: started\n")
			b.WriteString("        enabled: true\n")
		}
		b.WriteString("      become: true\n")
	}

	if strings.Contains(lower, "config") || strings.Contains(lower, "configure") || strings.Contains(lower, "template") {
		b.WriteString("\n    - name: Deploy configuration from template\n")
		if targetOS == "windows" {
			b.WriteString("      ansible.windows.win_template:\n")
			b.WriteString("        src: templates/config.j2\n")
			b.WriteString("        dest: 'C:\\ProgramData\\app\\config.cfg'\n")
		} else {
			b.WriteString("      ansible.builtin.template:\n")
			b.WriteString("        src: templates/config.j2\n")
			b.WriteString("        dest: /etc/app/config.cfg\n")
			b.WriteString("        owner: root\n")
			b.WriteString("        group: root\n")
			b.WriteString("        mode: '0644'\n")
		}
		b.WriteString("      become: true\n")
		b.WriteString("      notify: Restart service\n")
	}

	if strings.Contains(lower, "user") || strings.Contains(lower, "account") {
		b.WriteString("\n    - name: Ensure application user exists\n")
		b.WriteString("      ansible.builtin.user:\n")
		b.WriteString("        name: \"{{ app_user }}\"\n")
		b.WriteString("        state: present\n")
		b.WriteString("        system: true\n")
		b.WriteString("        shell: /sbin/nologin\n")
		b.WriteString("      become: true\n")
	}

	if strings.Contains(lower, "firewall") || strings.Contains(lower, "port") {
		b.WriteString("\n    - name: Open required firewall port\n")
		if targetOS == "rhel" {
			b.WriteString("      ansible.posix.firewalld:\n")
			b.WriteString("        port: \"{{ app_port }}/tcp\"\n")
			b.WriteString("        state: enabled\n")
			b.WriteString("        permanent: true\n")
			b.WriteString("        immediate: true\n")
		} else {
			b.WriteString("      community.general.ufw:\n")
			b.WriteString("        rule: allow\n")
			b.WriteString("        port: \"{{ app_port }}\"\n")
			b.WriteString("        proto: tcp\n")
			b.WriteString("        state: enabled\n")
		}
		b.WriteString("      become: true\n")
	}

	// Default task if nothing specific was matched
	if b.Len() == 0 {
		b.WriteString("\n    - name: " + sentenceCase(taskDesc) + "\n")
		b.WriteString("      # TODO: implement task for: " + taskDesc + "\n")
		b.WriteString("      ansible.builtin.debug:\n")
		b.WriteString("        msg: \"Placeholder for: {{ task_description }}\"\n")
		b.WriteString("      vars:\n")
		b.WriteString("        task_description: \"" + taskDesc + "\"\n")
	}

	return b.String()
}

func generateProjectScaffold(taskDesc, targetOS string, collections []string, useRoles bool) map[string]string {
	roleName := descToRoleName(taskDesc)
	files := make(map[string]string)

	// Main playbook
	files["site.yml"] = generateSimplePlaybook(taskDesc, targetOS, collections, true)

	// Role structure
	files[fmt.Sprintf("roles/%s/tasks/main.yml", roleName)] = generateTasksYAML(taskDesc, targetOS)
	files[fmt.Sprintf("roles/%s/handlers/main.yml", roleName)] = generateHandlersYAML()
	files[fmt.Sprintf("roles/%s/defaults/main.yml", roleName)] = generateDefaultsYAML(taskDesc)
	files[fmt.Sprintf("roles/%s/templates/config.j2", roleName)] = "# Template for " + taskDesc + "\n# Add your Jinja2 template here\n"
	files[fmt.Sprintf("roles/%s/meta/main.yml", roleName)] = generateRoleMetaYAML(roleName, taskDesc)

	// Inventory example
	files["inventory/hosts.yml"] = generateInventoryYAML(targetOS)

	// Group vars
	files["group_vars/all.yml"] = generateGroupVarsYAML(taskDesc)

	_ = useRoles // useRoles is always true for project layout
	return files
}

func generateTasksYAML(taskDesc, targetOS string) string {
	return "---\n# Tasks for: " + taskDesc + "\n" + generateTasksForOS(taskDesc, targetOS)
}

func generateHandlersYAML() string {
	return `---
- name: Restart service
  ansible.builtin.service:
    name: "{{ service_name }}"
    state: restarted
`
}

func generateDefaultsYAML(taskDesc string) string {
	return fmt.Sprintf(`---
# Default variables for role
# Override these in group_vars, host_vars, or playbook vars

# Example defaults
# package_name: myapp
# service_name: myapp
# app_port: 8080
# app_user: myapp

# Role description: %s
`, taskDesc)
}

func generateRoleMetaYAML(roleName, taskDesc string) string {
	return fmt.Sprintf(`---
galaxy_info:
  role_name: %s
  description: %s
  author: your_name
  license: Apache-2.0
  min_ansible_version: "2.12"
  platforms:
    - name: EL
      versions:
        - "8"
        - "9"
    - name: Ubuntu
      versions:
        - focal
        - jammy

dependencies: []
`, roleName, taskDesc)
}

func generateInventoryYAML(targetOS string) string {
	group := "webservers"
	return fmt.Sprintf(`---
all:
  children:
    %s:
      hosts:
        host1.example.com:
        host2.example.com:
      vars:
        ansible_user: ansible
        # ansible_become: true  # set at task level instead
`, group)
}

func generateGroupVarsYAML(taskDesc string) string {
	return fmt.Sprintf(`---
# Group variables for all hosts
# These are examples — replace with your actual values

# Uncomment and set your variables:
# package_name: myapp
# service_name: myapp
# app_port: 8080

# Description: %s
`, taskDesc)
}

func sentenceCase(s string) string {
	if s == "" {
		return s
	}
	words := strings.Fields(s)
	if len(words) == 0 {
		return s
	}
	words[0] = strings.ToUpper(words[0][:1]) + words[0][1:]
	return strings.Join(words, " ")
}

func descToRoleName(desc string) string {
	words := strings.Fields(strings.ToLower(desc))
	var parts []string
	for _, w := range words {
		// Remove common filler words
		switch w {
		case "and", "the", "a", "an", "on", "in", "for", "to", "with":
			continue
		}
		// Remove non-alphanumeric
		clean := strings.Map(func(r rune) rune {
			if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
				return r
			}
			return '_'
		}, w)
		if clean != "" && clean != "_" {
			parts = append(parts, clean)
		}
		if len(parts) >= 3 {
			break
		}
	}
	if len(parts) == 0 {
		return "my_role"
	}
	return strings.Join(parts, "_")
}

func scaffoldNotes(useRoles bool, style string) []string {
	notes := []string{
		"gather_facts: true is enabled by default for fact-based conditionals",
		"become is applied at the task level, not the play level",
		"Handlers are used for service restarts triggered by configuration changes",
		"Variables are declared in vars: or group_vars — not hardcoded in tasks",
		"state: is explicit on all modules that support it",
		"All tasks and plays have descriptive name: values",
	}
	if useRoles || style == "project" {
		notes = append(notes, "Role structure follows ansible-galaxy init layout")
		notes = append(notes, "defaults/main.yml contains overridable defaults")
	}
	notes = append(notes, "Add no_log: true to any task that handles passwords or secrets")
	notes = append(notes, "Run ansible-lint to validate this playbook before use")
	return notes
}
