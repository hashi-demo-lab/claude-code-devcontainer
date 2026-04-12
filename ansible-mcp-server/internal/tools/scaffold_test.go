package tools

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestHandleGeneratePlaybookScaffold_Simple(t *testing.T) {
	handler := handleGeneratePlaybookScaffold()

	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]interface{}{
		"task_description": "install and configure nginx",
		"target_os":        "rhel",
	}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handler returned error: %v", result.Content)
	}

	text := extractResultText(t, result)

	var output scaffoldOutput
	if err := json.Unmarshal([]byte(text), &output); err != nil {
		t.Fatalf("unmarshal output: %v\nraw: %s", err, text)
	}

	if output.Playbook == "" {
		t.Error("expected non-empty playbook")
	}
	if len(output.Notes) == 0 {
		t.Error("expected non-empty notes")
	}

	// Verify best practices are in the scaffold
	if !strings.Contains(output.Playbook, "gather_facts: true") {
		t.Error("playbook missing 'gather_facts: true'")
	}
	if !strings.Contains(output.Playbook, "become: true") {
		t.Error("playbook missing task-level become")
	}
	if !strings.Contains(output.Playbook, "name:") {
		t.Error("playbook missing task names")
	}
	if !strings.Contains(output.Playbook, "handlers:") {
		t.Error("playbook missing handlers section")
	}
}

func TestHandleGeneratePlaybookScaffold_ProjectStyle(t *testing.T) {
	handler := handleGeneratePlaybookScaffold()

	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]interface{}{
		"task_description": "deploy postgresql database",
		"style":            "project",
		"target_os":        "rhel",
	}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handler returned error: %v", result.Content)
	}

	text := extractResultText(t, result)
	var output scaffoldOutput
	if err := json.Unmarshal([]byte(text), &output); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}

	if len(output.Files) == 0 {
		t.Error("expected files in project scaffold")
	}

	// Verify expected project files exist
	expectedFiles := []string{"site.yml", "inventory/hosts.yml", "group_vars/all.yml"}
	for _, f := range expectedFiles {
		if _, ok := output.Files[f]; !ok {
			t.Errorf("expected file %q in project scaffold, got keys: %v", f, fileKeys(output.Files))
		}
	}

	// Verify role structure
	foundRole := false
	for path := range output.Files {
		if strings.HasPrefix(path, "roles/") {
			foundRole = true
			break
		}
	}
	if !foundRole {
		t.Error("expected role directory structure in project scaffold")
	}
}

func TestHandleGeneratePlaybookScaffold_MissingDescription(t *testing.T) {
	handler := handleGeneratePlaybookScaffold()

	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]interface{}{}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for missing task_description")
	}
}

func TestHandleGeneratePlaybookScaffold_WindowsOS(t *testing.T) {
	handler := handleGeneratePlaybookScaffold()

	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]interface{}{
		"task_description": "install and start service on windows",
		"target_os":        "windows",
	}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handler returned error: %v", result.Content)
	}

	text := extractResultText(t, result)
	var output scaffoldOutput
	if err := json.Unmarshal([]byte(text), &output); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}

	if !strings.Contains(output.Playbook, "windows") {
		t.Error("expected windows-specific module references in playbook")
	}
}

func TestHandleGeneratePlaybookScaffold_WithCollections(t *testing.T) {
	handler := handleGeneratePlaybookScaffold()

	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]interface{}{
		"task_description": "configure firewall",
		"collections":      []interface{}{"ansible.posix", "community.general"},
	}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handler returned error: %v", result.Content)
	}

	text := extractResultText(t, result)
	var output scaffoldOutput
	if err := json.Unmarshal([]byte(text), &output); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}

	if !strings.Contains(output.Playbook, "ansible.posix") {
		t.Error("expected ansible.posix in collections declaration")
	}
	if !strings.Contains(output.Playbook, "community.general") {
		t.Error("expected community.general in collections declaration")
	}
}

// --- helpers ---

func extractResultText(t *testing.T, result *mcp.CallToolResult) string {
	t.Helper()
	if len(result.Content) == 0 {
		t.Fatal("result has no content")
	}
	for _, item := range result.Content {
		if tc, ok := item.(mcp.TextContent); ok {
			return tc.Text
		}
	}
	t.Fatalf("no text content in result: %v", result.Content)
	return ""
}

func fileKeys(files map[string]string) []string {
	keys := make([]string, 0, len(files))
	for k := range files {
		keys = append(keys, k)
	}
	return keys
}
