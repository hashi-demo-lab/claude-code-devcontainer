package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerValidateTools(s *server.MCPServer, _ *Dependencies) {
	s.AddTool(
		mcp.NewTool("validate_playbook",
			mcp.WithDescription("Lint and validate an Ansible playbook using ansible-lint. Returns violations with rule IDs, line numbers, severity, and remediation hints. Requires ansible-lint on PATH."),
			mcp.WithString("playbook",
				mcp.Required(),
				mcp.Description("YAML playbook content to validate"),
			),
			mcp.WithString("profile",
				mcp.Description("ansible-lint profile: min, basic, moderate, safety, shared, production (default: basic)"),
			),
		),
		handleValidatePlaybook(),
	)
}

// lintResult is the output from validate_playbook.
type lintResult struct {
	Passed     bool            `json:"passed"`
	Violations []lintViolation `json:"violations"`
	Summary    lintSummary     `json:"summary"`
}

// lintViolation is a single ansible-lint finding.
type lintViolation struct {
	RuleID      string `json:"rule_id"`
	Description string `json:"description"`
	Line        int    `json:"line"`
	Severity    string `json:"severity"`
	Remediation string `json:"remediation"`
}

// lintSummary totals violations by severity.
type lintSummary struct {
	Warnings int `json:"warnings"`
	Errors   int `json:"errors"`
	Total    int `json:"total"`
}

// ansibleLintJSONOutput is the JSON schema emitted by ansible-lint --format sarif or json.
type ansibleLintJSONOutput struct {
	Runs []struct {
		Results []struct {
			RuleID  string `json:"ruleId"`
			Level   string `json:"level"` // "warning" | "error" | "note"
			Message struct {
				Text string `json:"text"`
			} `json:"message"`
			Locations []struct {
				PhysicalLocation struct {
					Region struct {
						StartLine int `json:"startLine"`
					} `json:"region"`
				} `json:"physicalLocation"`
			} `json:"locations"`
		} `json:"results"`
	} `json:"runs"`
}

func handleValidatePlaybook() Handler {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		playbook := argString(req, "playbook")
		if playbook == "" {
			return errResult("ERROR: playbook content is required")
		}
		profile := argString(req, "profile")
		if profile == "" {
			profile = "basic"
		}

		// Validate profile value
		validProfiles := map[string]bool{
			"min": true, "basic": true, "moderate": true,
			"safety": true, "shared": true, "production": true,
		}
		if !validProfiles[profile] {
			return errResult(fmt.Sprintf("ERROR: invalid profile %q: must be one of min, basic, moderate, safety, shared, production", profile))
		}

		// Check ansible-lint is available
		if _, err := exec.LookPath("ansible-lint"); err != nil {
			return errResult("ERROR [ansible-lint]: ansible-lint not found on PATH. Install it with: pip install ansible-lint")
		}

		// Write playbook to a temp file
		tmpFile, err := os.CreateTemp("", "ansible-mcp-playbook-*.yml")
		if err != nil {
			return errResult(fmt.Sprintf("ERROR: creating temp file: %s", err))
		}
		defer os.Remove(tmpFile.Name())

		if _, err := tmpFile.WriteString(playbook); err != nil {
			tmpFile.Close()
			return errResult(fmt.Sprintf("ERROR: writing temp file: %s", err))
		}
		tmpFile.Close()

		// Run ansible-lint
		cmd := exec.CommandContext(ctx, "ansible-lint",
			"--format", "json",
			"--profile", profile,
			"--nocolor",
			tmpFile.Name(),
		)

		output, err := cmd.Output()
		exitErr, isExitErr := err.(*exec.ExitError)

		// ansible-lint exits non-zero when violations are found — that's expected
		if err != nil && !isExitErr {
			return errResult(fmt.Sprintf("ERROR [ansible-lint]: failed to run: %s", err))
		}

		result := &lintResult{}

		// Try to parse JSON output
		if len(output) > 0 {
			violations, parseErr := parseAnsibleLintJSON(output, tmpFile.Name())
			if parseErr == nil {
				result.Violations = violations
			} else {
				// Fall back to parsing stderr/stdout as text
				stderr := ""
				if isExitErr && exitErr != nil {
					stderr = string(exitErr.Stderr)
				}
				result.Violations = parseAnsibleLintText(string(output) + "\n" + stderr)
			}
		}

		// Compute summary
		for _, v := range result.Violations {
			result.Summary.Total++
			if v.Severity == "error" {
				result.Summary.Errors++
			} else {
				result.Summary.Warnings++
			}
		}

		result.Passed = result.Summary.Errors == 0 && result.Summary.Total == 0

		out, _ := json.MarshalIndent(result, "", "  ")
		return textResult(string(out))
	}
}

func parseAnsibleLintJSON(data []byte, tmpPath string) ([]lintViolation, error) {
	// ansible-lint JSON output format (SARIF-like)
	var raw ansibleLintJSONOutput
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	var violations []lintViolation
	for _, run := range raw.Runs {
		for _, result := range run.Results {
			v := lintViolation{
				RuleID:      result.RuleID,
				Description: result.Message.Text,
			}
			if len(result.Locations) > 0 {
				v.Line = result.Locations[0].PhysicalLocation.Region.StartLine
			}
			switch result.Level {
			case "error":
				v.Severity = "error"
			default:
				v.Severity = "warning"
			}
			v.Remediation = remediationForRule(result.RuleID)
			violations = append(violations, v)
		}
	}
	return violations, nil
}

func parseAnsibleLintText(output string) []lintViolation {
	var violations []lintViolation
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Format: path:line: rule-id msg
		// or: WARNING: ... [rule-id]
		if strings.Contains(line, ": ") {
			v := lintViolation{
				Description: line,
				Severity:    "warning",
			}
			if strings.Contains(strings.ToLower(line), "error") {
				v.Severity = "error"
			}
			// Extract rule ID from [brackets] if present
			if start := strings.LastIndex(line, "["); start != -1 {
				if end := strings.LastIndex(line, "]"); end > start {
					v.RuleID = line[start+1 : end]
					v.Remediation = remediationForRule(v.RuleID)
				}
			}
			violations = append(violations, v)
		}
	}
	return violations
}

// remediationForRule returns a human-readable remediation hint for known ansible-lint rules.
func remediationForRule(ruleID string) string {
	hints := map[string]string{
		"command-instead-of-module":  "Replace the shell/command task with the appropriate Ansible module (e.g., use ansible.builtin.service instead of 'systemctl restart')",
		"no-free-form":               "Use the canonical YAML key: value syntax instead of the free-form shorthand",
		"yaml":                       "Fix YAML formatting: check indentation, trailing spaces, or missing newline at end of file",
		"name[casing]":               "Task names should use sentence case (capitalize the first word only)",
		"name[missing]":              "Add a 'name:' field to every task describing what it does",
		"no-changed-when":            "Add 'changed_when: false' or a proper changed condition to command/shell tasks",
		"risky-file-permissions":     "Specify explicit file permissions (e.g., mode: '0644') instead of relying on defaults",
		"key-order":                  "Reorder task keys to the conventional order: name, module, args, become, notify, etc.",
		"partial-become":             "When using become, also set become_user explicitly to avoid privilege escalation ambiguity",
		"fqcn":                       "Use the fully qualified collection name (e.g., ansible.builtin.copy instead of copy)",
		"fqcn[action-core]":          "Use the fully qualified collection name for built-in modules (e.g., ansible.builtin.package)",
		"var-naming[no-role-prefix]": "Role variables should be prefixed with the role name to avoid conflicts",
		"no-log-password":            "Add 'no_log: true' to tasks that handle passwords, tokens, or other secrets",
		"galaxy[version-incorrect]":  "Set a valid semantic version in meta/main.yml (e.g., 1.0.0)",
		"package-latest":             "Pin package versions with state: present instead of state: latest for reproducible builds",
	}
	if hint, ok := hints[ruleID]; ok {
		return hint
	}
	return fmt.Sprintf("See https://ansible-lint.readthedocs.io/rules/%s/ for remediation guidance", ruleID)
}
