package tools

import (
	"testing"
)

func TestResolveTopicPath(t *testing.T) {
	tests := []struct {
		topic    string
		wantPath string
		wantOK   bool
	}{
		{"handlers", "playbook_guide/playbooks_handlers.html", true},
		{"tags", "playbook_guide/playbooks_tags.html", true},
		{"error handling", "playbook_guide/playbooks_error_handling.html", true},
		{"ignore_errors", "playbook_guide/playbooks_error_handling.html", true},
		{"idempotency", "reference_appendices/test_strategies.html", true},
		{"roles", "playbook_guide/playbooks_reuse_roles.html", true},
		{"vault", "vault_guide/index.html", true},
		{"secrets", "vault_guide/index.html", true},
		{"loops", "playbook_guide/playbooks_loops.html", true},
		{"with_items", "playbook_guide/playbooks_loops.html", true},
		{"conditionals", "playbook_guide/playbooks_conditionals.html", true},
		{"when", "playbook_guide/playbooks_conditionals.html", true},
		{"templates", "playbook_guide/playbooks_templating.html", true},
		{"jinja2", "playbook_guide/playbooks_templating.html", true},
		{"inventory", "inventory_guide/index.html", true},
		{"galaxy", "galaxy/user_guide.html", true},
		{"variables", "playbook_guide/playbooks_variables.html", true},
		{"collections", "collections_guide/index.html", true},
		{"unknown topic xyz", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.topic, func(t *testing.T) {
			path, ok := resolveTopicPath(tt.topic)
			if ok != tt.wantOK {
				t.Errorf("resolveTopicPath(%q) ok = %v, want %v", tt.topic, ok, tt.wantOK)
			}
			if tt.wantOK && path != tt.wantPath {
				t.Errorf("resolveTopicPath(%q) = %q, want %q", tt.topic, path, tt.wantPath)
			}
		})
	}
}

func TestDescToRoleName(t *testing.T) {
	tests := []struct {
		desc string
		want string
	}{
		{"install nginx", "install_nginx"},
		{"install and configure nginx on RHEL", "install_configure_nginx"},
		{"", "my_role"},
		{"  ", "my_role"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got := descToRoleName(tt.desc)
			if got != tt.want {
				t.Errorf("descToRoleName(%q) = %q, want %q", tt.desc, got, tt.want)
			}
		})
	}
}
