package tools

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestHandleGenerateTestCases_Default(t *testing.T) {
	handler := handleGenerateTestCases()

	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]interface{}{
		"task_description": "install and configure nginx",
		"role_name":        "nginx",
	}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if result.IsError {
		t.Fatalf("handler returned error: %v", result.Content)
	}

	text := extractResultText(t, result)
	var output TestCaseOutput
	if err := json.Unmarshal([]byte(text), &output); err != nil {
		t.Fatalf("unmarshal output: %v\nraw: %s", err, text)
	}

	// Verify required Molecule files are present
	requiredFiles := []string{
		"molecule/default/molecule.yml",
		"molecule/default/converge.yml",
		"molecule/default/verify.yml",
		"molecule/default/prepare.yml",
		"molecule/default/cleanup.yml",
	}
	for _, f := range requiredFiles {
		if _, ok := output.Files[f]; !ok {
			t.Errorf("missing required file %q, got: %v", f, fileKeys(output.Files))
		}
	}

	if len(output.Notes) == 0 {
		t.Error("expected non-empty notes")
	}
}

func TestHandleGenerateTestCases_MoleculeYML(t *testing.T) {
	handler := handleGenerateTestCases()

	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]interface{}{
		"task_description": "install nginx",
		"driver":           "podman",
	}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}

	text := extractResultText(t, result)
	var output TestCaseOutput
	if err := json.Unmarshal([]byte(text), &output); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}

	moleculeYML := output.Files["molecule/default/molecule.yml"]

	// Validate molecule.yml structure
	if !strings.Contains(moleculeYML, "driver:") {
		t.Error("molecule.yml missing driver section")
	}
	if !strings.Contains(moleculeYML, "name: podman") {
		t.Error("molecule.yml missing podman driver name")
	}
	if !strings.Contains(moleculeYML, "platforms:") {
		t.Error("molecule.yml missing platforms section")
	}
	if !strings.Contains(moleculeYML, "provisioner:") {
		t.Error("molecule.yml missing provisioner section")
	}
	if !strings.Contains(moleculeYML, "verifier:") {
		t.Error("molecule.yml missing verifier section")
	}
	if !strings.Contains(moleculeYML, "interpreter_python: auto_silent") {
		t.Error("molecule.yml missing interpreter_python config")
	}
}

func TestHandleGenerateTestCases_DefaultPlatform(t *testing.T) {
	handler := handleGenerateTestCases()

	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]interface{}{
		"task_description": "configure service",
	}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}

	text := extractResultText(t, result)
	var output TestCaseOutput
	if err := json.Unmarshal([]byte(text), &output); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}

	moleculeYML := output.Files["molecule/default/molecule.yml"]
	if !strings.Contains(moleculeYML, "ubi9") {
		t.Error("expected default UBI9 platform when no platforms specified")
	}
}

func TestHandleGenerateTestCases_InvalidDriver(t *testing.T) {
	handler := handleGenerateTestCases()

	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]interface{}{
		"task_description": "test",
		"driver":           "vagrant",
	}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error for invalid driver")
	}
}

func TestHandleGenerateTestCases_ConvergeYMLUsesPlaybook(t *testing.T) {
	handler := handleGenerateTestCases()

	playbook := `---
- name: Install nginx web server
  hosts: webservers
  tasks:
    - name: Install nginx
      ansible.builtin.package:
        name: nginx
        state: present`

	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]interface{}{
		"playbook":  playbook,
		"role_name": "nginx",
	}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}

	text := extractResultText(t, result)
	var output TestCaseOutput
	if err := json.Unmarshal([]byte(text), &output); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}

	converge := output.Files["molecule/default/converge.yml"]
	// hosts should be replaced with 'all' for Molecule
	if strings.Contains(converge, "hosts: webservers") {
		t.Error("converge.yml should replace specific hosts with 'all'")
	}
	if !strings.Contains(converge, "hosts: all") {
		t.Error("converge.yml should set hosts: all for Molecule")
	}
}

func TestHandleGenerateTestCases_VerifyYMLAssertions(t *testing.T) {
	handler := handleGenerateTestCases()

	req := mcp.CallToolRequest{}
	req.Params.Arguments = map[string]interface{}{
		"task_description": "install nginx and start the service",
	}

	result, err := handler(context.Background(), req)
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}

	text := extractResultText(t, result)
	var output TestCaseOutput
	if err := json.Unmarshal([]byte(text), &output); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}

	verify := output.Files["molecule/default/verify.yml"]
	if !strings.Contains(verify, "ansible.builtin.assert") {
		t.Error("verify.yml should contain assert tasks")
	}
}
