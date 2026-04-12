package docs

import (
	"strings"
	"testing"
)

const sampleModuleHTML = `<!DOCTYPE html>
<html>
<head><title>ansible.posix.firewalld – Manage firewalld rules — Ansible Documentation</title></head>
<body>
<h1>ansible.posix.firewalld – Manage firewalld rules</h1>
<section id="parameters">
  <h2>Parameters</h2>
  <ul>
    <li class="ansible-option">
      <strong>state</strong>
      <p>Enable or disable the rule.</p>
      Type: str
      Required: false
      Default: enabled
      Choices: enabled, disabled, present, absent
    </li>
    <li class="ansible-option">
      <strong>port</strong>
      <p>Name or number of the TCP or UDP port or port range to add or remove.</p>
      Type: str
      Required: false
    </li>
  </ul>
</section>
<section id="examples">
  <h2>Examples</h2>
  <pre>
- name: Enable firewalld
  ansible.posix.firewalld:
    state: enabled
  </pre>
</section>
<section id="return-values">
  <h2>Return Values</h2>
  <ul>
    <li class="ansible-option">
      <strong>changed</strong>
      <p>Whether the rule was changed.</p>
    </li>
  </ul>
</section>
<section id="notes">
  <h2>Notes</h2>
  <ul>
    <li>Requires firewalld >= 0.2.11</li>
  </ul>
</section>
</body>
</html>`

const sampleKeywordsHTML = `<!DOCTYPE html>
<html>
<body>
<dl>
  <dt>become</dt>
  <dd>Boolean that controls if privilege escalation is used or not on Task execution.</dd>
  <dt>delegate_to</dt>
  <dd>Host to execute task instead of the target host.</dd>
  <dt>when</dt>
  <dd>Conditional expression, determines if an iteration of a task is run or not.</dd>
</dl>
</body>
</html>`

const sampleSpecialVarsHTML = `<!DOCTYPE html>
<html>
<body>
<dl>
  <dt>inventory_hostname</dt>
  <dd>The inventory name for the 'current' host being iterated over in the play.</dd>
  <dt>ansible_facts</dt>
  <dd>Contains any facts gathered or cached for the inventory_hostname.</dd>
  <dt>hostvars</dt>
  <dd>A dictionary/map with all the hosts in inventory and variables assigned to them.</dd>
</dl>
</body>
</html>`

func TestParseModuleDocs(t *testing.T) {
	doc, err := ParseModuleDocs(sampleModuleHTML)
	if err != nil {
		t.Fatalf("ParseModuleDocs() error = %v", err)
	}

	if !strings.Contains(doc.ShortDescription, "firewalld") {
		t.Errorf("ShortDescription = %q, expected to contain 'firewalld'", doc.ShortDescription)
	}

	if len(doc.Parameters) == 0 {
		t.Error("expected parameters, got none")
	}

	// Check that at least one parameter has a name
	found := false
	for _, p := range doc.Parameters {
		if p.Name == "state" || p.Name == "port" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected to find 'state' or 'port' parameter, got: %v", doc.Parameters)
	}

	if doc.Examples == "" {
		t.Error("expected examples, got empty string")
	}

	if len(doc.ReturnValues) == 0 {
		t.Error("expected return values, got none")
	}

	if len(doc.Notes) == 0 {
		t.Error("expected notes, got none")
	}
}

func TestParsePlaybookKeywords(t *testing.T) {
	keywords, err := ParsePlaybookKeywords(sampleKeywordsHTML)
	if err != nil {
		t.Fatalf("ParsePlaybookKeywords() error = %v", err)
	}

	if len(keywords) < 3 {
		t.Errorf("expected >= 3 keywords, got %d", len(keywords))
	}

	names := make(map[string]bool)
	for _, kw := range keywords {
		names[kw.Name] = true
		if kw.Description == "" {
			t.Errorf("keyword %q has empty description", kw.Name)
		}
	}

	for _, expected := range []string{"become", "delegate_to", "when"} {
		if !names[expected] {
			t.Errorf("expected keyword %q not found", expected)
		}
	}
}

func TestParseSpecialVariables(t *testing.T) {
	vars, err := ParseSpecialVariables(sampleSpecialVarsHTML)
	if err != nil {
		t.Fatalf("ParseSpecialVariables() error = %v", err)
	}

	if len(vars) == 0 {
		t.Error("expected variables, got none")
	}

	names := make(map[string]bool)
	for _, v := range vars {
		names[v.Name] = true
	}

	for _, expected := range []string{"inventory_hostname", "ansible_facts", "hostvars"} {
		if !names[expected] {
			t.Errorf("expected variable %q not found in %v", expected, vars)
		}
	}
}

func TestParseGenericContent(t *testing.T) {
	html := `<html><body>
<div class="document">
<h1>Best Practices</h1>
<p>Use idempotent modules whenever possible.</p>
<p>Always test with check mode first.</p>
</div>
</body></html>`

	content, err := ParseGenericContent(html)
	if err != nil {
		t.Fatalf("ParseGenericContent() error = %v", err)
	}
	if content == "" {
		t.Error("expected non-empty content")
	}
	if !strings.Contains(content, "idempotent") {
		t.Errorf("content = %q, expected to contain 'idempotent'", content)
	}
}

func TestExtractTitle(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantSub  string
	}{
		{
			name:    "h1 tag",
			html:    `<html><body><h1>ansible.posix.firewalld module</h1></body></html>`,
			wantSub: "ansible.posix.firewalld",
		},
		{
			name:    "title tag with suffix",
			html:    `<html><head><title>Firewalld — Ansible Documentation</title></head><body></body></html>`,
			wantSub: "Firewalld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, _ := ParseModuleDocs(tt.html)
			if !strings.Contains(doc.ShortDescription, tt.wantSub) {
				t.Errorf("ShortDescription = %q, expected to contain %q", doc.ShortDescription, tt.wantSub)
			}
		})
	}
}
