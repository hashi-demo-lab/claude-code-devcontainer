package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad_Defaults(t *testing.T) {
	// Clear relevant env vars
	for _, k := range []string{
		"ANSIBLE_HUB_SAAS_URL", "ANSIBLE_HUB_SAAS_TOKEN",
		"ANSIBLE_AAP_URL", "ANSIBLE_AAP_AUTH_MODE",
		"ANSIBLE_HUB_TARGET", "ANSIBLE_GALAXY_URL",
		"ANSIBLE_REQUEST_TIMEOUT", "ANSIBLE_MAX_RETRIES",
		"MCP_RATE_LIMIT_GLOBAL", "MCP_RATE_LIMIT_SESSION",
		"ANSIBLE_DOCS_BASE_URL",
	} {
		os.Unsetenv(k)
	}

	cfg, err := Load("1.0.0")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	tests := []struct {
		name string
		got  interface{}
		want interface{}
	}{
		{"HubSaaSURL", cfg.HubSaaSURL, "https://cloud.redhat.com/api/automation-hub/v3/"},
		{"HubTarget", cfg.HubTarget, "saas"},
		{"GalaxyURL", cfg.GalaxyURL, "https://galaxy.ansible.com/api/v3/"},
		{"MaxRetries", cfg.MaxRetries, 3},
		{"RequestTimeout", cfg.RequestTimeout, 10 * time.Second},
		{"GlobalRateLimit", cfg.GlobalRateLimit, "10:20"},
		{"SessionRateLimit", cfg.SessionRateLimit, "5:10"},
		{"Version", cfg.Version, "1.0.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %v, want %v", tt.got, tt.want)
			}
		})
	}

	// DocsBaseURL should end with /
	if cfg.DocsBaseURL[len(cfg.DocsBaseURL)-1] != '/' {
		t.Errorf("DocsBaseURL %q does not end with /", cfg.DocsBaseURL)
	}
}

func TestLoad_EnvOverrides(t *testing.T) {
	os.Setenv("ANSIBLE_HUB_TARGET", "aap")
	os.Setenv("ANSIBLE_REQUEST_TIMEOUT", "30")
	os.Setenv("ANSIBLE_MAX_RETRIES", "5")
	defer func() {
		os.Unsetenv("ANSIBLE_HUB_TARGET")
		os.Unsetenv("ANSIBLE_REQUEST_TIMEOUT")
		os.Unsetenv("ANSIBLE_MAX_RETRIES")
	}()

	cfg, err := Load("dev")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.HubTarget != "aap" {
		t.Errorf("HubTarget = %q, want aap", cfg.HubTarget)
	}
	if cfg.RequestTimeout != 30*time.Second {
		t.Errorf("RequestTimeout = %v, want 30s", cfg.RequestTimeout)
	}
	if cfg.MaxRetries != 5 {
		t.Errorf("MaxRetries = %d, want 5", cfg.MaxRetries)
	}
}

func TestLoad_InvalidHubTarget(t *testing.T) {
	os.Setenv("ANSIBLE_HUB_TARGET", "invalid")
	defer os.Unsetenv("ANSIBLE_HUB_TARGET")

	_, err := Load("dev")
	if err == nil {
		t.Error("Load() expected error for invalid ANSIBLE_HUB_TARGET")
	}
}

func TestLoad_InvalidTimeout(t *testing.T) {
	os.Setenv("ANSIBLE_REQUEST_TIMEOUT", "notanumber")
	defer os.Unsetenv("ANSIBLE_REQUEST_TIMEOUT")

	_, err := Load("dev")
	if err == nil {
		t.Error("Load() expected error for invalid ANSIBLE_REQUEST_TIMEOUT")
	}
}

func TestLoad_DocsBaseURLTrailingSlash(t *testing.T) {
	os.Setenv("ANSIBLE_DOCS_BASE_URL", "https://docs.ansible.com/projects/ansible/latest")
	defer os.Unsetenv("ANSIBLE_DOCS_BASE_URL")

	cfg, err := Load("dev")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.DocsBaseURL != "https://docs.ansible.com/projects/ansible/latest/" {
		t.Errorf("DocsBaseURL = %q, want trailing slash", cfg.DocsBaseURL)
	}
}
