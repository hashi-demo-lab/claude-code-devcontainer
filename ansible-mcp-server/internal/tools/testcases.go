package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerTestCaseTools(s *server.MCPServer, _ *Dependencies) {
	s.AddTool(
		mcp.NewTool("generate_test_cases",
			mcp.WithDescription("Generate Molecule v6+ test scenarios for an Ansible playbook or role. Returns a complete molecule/ directory structure with molecule.yml, converge.yml, verify.yml, prepare.yml, and cleanup.yml."),
			mcp.WithString("playbook",
				mcp.Description("Playbook YAML content to generate tests for (use this or task_description)"),
			),
			mcp.WithString("task_description",
				mcp.Description("Task description used to infer test assertions (used if playbook not provided)"),
			),
			mcp.WithString("role_name",
				mcp.Description("Role name for role-based test structure"),
			),
			mcp.WithString("driver",
				mcp.Description("Molecule driver: docker, podman, delegated (default: docker)"),
			),
			mcp.WithArray("platforms",
				mcp.Description("Target platforms as objects with name, image, pre_build_image fields"),
			),
		),
		handleGenerateTestCases(),
	)
}

// Platform describes a Molecule test platform.
type Platform struct {
	Name          string `json:"name"`
	Image         string `json:"image"`
	PreBuildImage bool   `json:"pre_build_image"`
}

// TestCaseOutput is the output from generate_test_cases.
type TestCaseOutput struct {
	Files map[string]string `json:"files"`
	Notes []string          `json:"notes"`
}

var defaultPlatform = Platform{
	Name:          "instance",
	Image:         "registry.access.redhat.com/ubi9/ubi-init",
	PreBuildImage: true,
}

func handleGenerateTestCases() Handler {
	return func(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		playbook := argString(req, "playbook")
		taskDesc := argString(req, "task_description")
		roleName := argString(req, "role_name")
		driver := argString(req, "driver")
		if driver == "" {
			driver = "docker"
		}

		// Validate driver
		switch driver {
		case "docker", "podman", "delegated":
		default:
			return errResult(fmt.Sprintf("ERROR: invalid driver %q: must be docker, podman, or delegated", driver))
		}

		// Parse platforms
		platforms := parsePlatforms(req)
		if len(platforms) == 0 {
			platforms = []Platform{defaultPlatform}
		}

		// Infer role name from playbook or description
		if roleName == "" {
			if taskDesc != "" {
				roleName = descToRoleName(taskDesc)
			} else {
				roleName = "my_role"
			}
		}

		// Infer what the converge play should test
		if taskDesc == "" && playbook != "" {
			taskDesc = extractTaskDescFromPlaybook(playbook)
		}

		output := &TestCaseOutput{
			Files: generateMoleculeFiles(roleName, taskDesc, playbook, driver, platforms),
			Notes: moleculeNotes(driver),
		}

		out, _ := json.MarshalIndent(output, "", "  ")
		return textResult(string(out))
	}
}

func parsePlatforms(req mcp.CallToolRequest) []Platform {
	raw, ok := req.Params.Arguments["platforms"].([]interface{})
	if !ok {
		return nil
	}
	var platforms []Platform
	for _, item := range raw {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		p := Platform{
			Name:  fmt.Sprintf("%v", m["name"]),
			Image: fmt.Sprintf("%v", m["image"]),
		}
		if pb, ok := m["pre_build_image"].(bool); ok {
			p.PreBuildImage = pb
		}
		if p.Name != "" && p.Image != "" {
			platforms = append(platforms, p)
		}
	}
	return platforms
}

func generateMoleculeFiles(roleName, taskDesc, playbook, driver string, platforms []Platform) map[string]string {
	files := make(map[string]string)

	files["molecule/default/molecule.yml"] = generateMoleculeYAML(driver, platforms)
	files["molecule/default/converge.yml"] = generateConvergeYAML(roleName, taskDesc, playbook)
	files["molecule/default/verify.yml"] = generateVerifyYAML(taskDesc, playbook)
	files["molecule/default/prepare.yml"] = generatePrepareYAML()
	files["molecule/default/cleanup.yml"] = generateCleanupYAML()

	return files
}

func generateMoleculeYAML(driver string, platforms []Platform) string {
	var b strings.Builder

	b.WriteString("---\n")
	b.WriteString("dependency:\n")
	b.WriteString("  name: galaxy\n")
	b.WriteString("\n")
	b.WriteString("driver:\n")
	b.WriteString(fmt.Sprintf("  name: %s\n", driver))
	b.WriteString("\n")
	b.WriteString("platforms:\n")
	for _, p := range platforms {
		b.WriteString(fmt.Sprintf("  - name: %s\n", p.Name))
		b.WriteString(fmt.Sprintf("    image: %s\n", p.Image))
		if p.PreBuildImage {
			b.WriteString("    pre_build_image: true\n")
		}
	}
	b.WriteString("\n")
	b.WriteString("provisioner:\n")
	b.WriteString("  name: ansible\n")
	b.WriteString("  config_options:\n")
	b.WriteString("    defaults:\n")
	b.WriteString("      interpreter_python: auto_silent\n")
	b.WriteString("\n")
	b.WriteString("verifier:\n")
	b.WriteString("  name: ansible\n")
	b.WriteString("\n")
	b.WriteString("# Idempotency check: run converge twice and assert no changes on second run\n")
	b.WriteString("lint: |\n")
	b.WriteString("  set -e\n")
	b.WriteString("  ansible-lint\n")

	return b.String()
}

func generateConvergeYAML(roleName, taskDesc, playbook string) string {
	if playbook != "" {
		// Use the provided playbook as the converge play, adjusting hosts
		return injectMoleculeHosts(playbook)
	}

	var b strings.Builder
	b.WriteString("---\n")
	b.WriteString("- name: Converge\n")
	b.WriteString("  hosts: all\n")
	b.WriteString("  gather_facts: true\n")
	b.WriteString("\n")
	b.WriteString("  roles:\n")
	b.WriteString(fmt.Sprintf("    - role: %s\n", roleName))

	return b.String()
}

// injectMoleculeHosts replaces the hosts field in a playbook with 'all' for Molecule.
func injectMoleculeHosts(playbook string) string {
	lines := strings.Split(playbook, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "hosts:") {
			indent := strings.Index(line, "hosts:")
			lines[i] = line[:indent] + "hosts: all  # molecule: replaced from original"
			break
		}
	}
	return strings.Join(lines, "\n")
}

func generateVerifyYAML(taskDesc, playbook string) string {
	assertions := inferAssertions(taskDesc, playbook)

	var b strings.Builder
	b.WriteString("---\n")
	b.WriteString("- name: Verify\n")
	b.WriteString("  hosts: all\n")
	b.WriteString("  gather_facts: true\n")
	b.WriteString("\n")
	b.WriteString("  tasks:\n")

	for _, a := range assertions {
		b.WriteString("\n    - name: " + a.name + "\n")
		b.WriteString("      " + a.module + ":\n")
		for k, v := range a.args {
			b.WriteString(fmt.Sprintf("        %s: %s\n", k, v))
		}
	}

	if len(assertions) == 0 {
		b.WriteString("\n    - name: Assert default verification passes\n")
		b.WriteString("      ansible.builtin.assert:\n")
		b.WriteString("        that:\n")
		b.WriteString("          - true\n")
		b.WriteString("        success_msg: \"Verification passed\"\n")
	}

	return b.String()
}

type assertion struct {
	name   string
	module string
	args   map[string]string
}

// inferAssertions derives verify.yml assertions from playbook content or task description.
func inferAssertions(taskDesc, playbook string) []assertion {
	var assertions []assertion
	content := strings.ToLower(taskDesc + " " + playbook)

	if strings.Contains(content, "service") || strings.Contains(content, "nginx") ||
		strings.Contains(content, "httpd") || strings.Contains(content, "apache") {
		assertions = append(assertions, assertion{
			name:   "Assert service is running",
			module: "ansible.builtin.service_facts",
			args:   map[string]string{},
		})
		assertions = append(assertions, assertion{
			name:   "Verify service state",
			module: "ansible.builtin.assert",
			args: map[string]string{
				"that": "\"ansible_facts.services['{{ service_name }}.service']['state'] == 'running'\"",
				"fail_msg": "\"Service {{ service_name }} is not running\"",
			},
		})
	}

	if strings.Contains(content, "package") || strings.Contains(content, "install") {
		assertions = append(assertions, assertion{
			name:   "Assert package is installed",
			module: "ansible.builtin.package_facts",
			args:   map[string]string{"manager": "auto"},
		})
		assertions = append(assertions, assertion{
			name:   "Verify package is present",
			module: "ansible.builtin.assert",
			args: map[string]string{
				"that":     "\"'{{ package_name }}' in ansible_facts.packages\"",
				"fail_msg": "\"Package {{ package_name }} is not installed\"",
			},
		})
	}

	if strings.Contains(content, "config") || strings.Contains(content, "template") || strings.Contains(content, "file") {
		assertions = append(assertions, assertion{
			name:   "Assert configuration file exists",
			module: "ansible.builtin.stat",
			args:   map[string]string{"path": "\"{{ config_file_path }}\"", "register": "config_stat"},
		})
		assertions = append(assertions, assertion{
			name:   "Verify configuration file is present",
			module: "ansible.builtin.assert",
			args: map[string]string{
				"that":     "\"config_stat.stat.exists\"",
				"fail_msg": "\"Configuration file {{ config_file_path }} does not exist\"",
			},
		})
	}

	if strings.Contains(content, "user") || strings.Contains(content, "account") {
		assertions = append(assertions, assertion{
			name:   "Assert application user exists",
			module: "ansible.builtin.getent",
			args:   map[string]string{"database": "passwd", "key": "\"{{ app_user }}\""},
		})
	}

	if strings.Contains(content, "port") || strings.Contains(content, "listen") {
		assertions = append(assertions, assertion{
			name:   "Assert service is listening on expected port",
			module: "ansible.builtin.wait_for",
			args: map[string]string{
				"port":    "\"{{ app_port }}\"",
				"state":   "started",
				"timeout": "5",
			},
		})
	}

	return assertions
}

func generatePrepareYAML() string {
	return `---
- name: Prepare
  hosts: all
  gather_facts: false

  tasks:
    - name: Ensure prerequisites are installed
      ansible.builtin.package:
        name:
          - python3
          - python3-pip
        state: present
      become: true
      # Add any pre-converge setup tasks here
`
}

func generateCleanupYAML() string {
	return `---
- name: Cleanup
  hosts: all
  gather_facts: false

  tasks:
    - name: Remove test artifacts
      ansible.builtin.debug:
        msg: "Cleanup complete"
      # Add teardown tasks here if needed
`
}

func extractTaskDescFromPlaybook(playbook string) string {
	for _, line := range strings.Split(playbook, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- name:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "- name:"))
		}
		if strings.HasPrefix(line, "name:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "name:"))
		}
	}
	return "ansible role"
}

func moleculeNotes(driver string) []string {
	notes := []string{
		"molecule/default/ contains the default scenario (converge + verify + idempotency)",
		"converge.yml applies the role/tasks under test",
		"verify.yml uses ansible.builtin.assert to check expected state",
		"prepare.yml runs before converge for pre-conditions",
		"cleanup.yml runs after tests for teardown",
		"Run 'molecule test' to execute the full test sequence",
		"Run 'molecule converge' then 'molecule verify' for faster iteration",
		"Idempotency is checked by running converge twice and asserting zero changes",
	}
	if driver == "docker" || driver == "podman" {
		notes = append(notes, fmt.Sprintf("Using %s driver: ensure %s daemon is running", driver, driver))
		notes = append(notes, "Images with init (ubi-init, systemd) are required for service management tests")
	}
	return notes
}
