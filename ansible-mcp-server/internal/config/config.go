package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all server configuration loaded from environment variables.
type Config struct {
	// Hub SaaS
	HubSaaSURL   string
	HubSaaSToken string

	// Hub on-prem AAP
	AAPURL      string
	AAPAuthMode string
	AAPToken    string
	AAPUsername string
	AAPPassword string

	// Hub target selection: saas | aap | both
	HubTarget string

	// Galaxy fallback
	GalaxyURL string

	// HTTP client
	RequestTimeout time.Duration
	MaxRetries     int

	// Rate limiting (rate:burst format)
	GlobalRateLimit  string
	SessionRateLimit string

	// Documentation base URL
	DocsBaseURL string

	// Build version (injected at build time via -ldflags)
	Version string
}

// Load reads all configuration from environment variables and applies defaults.
func Load(version string) (*Config, error) {
	c := &Config{
		HubSaaSURL:       getEnv("ANSIBLE_HUB_SAAS_URL", "https://cloud.redhat.com/api/automation-hub/v3/"),
		HubSaaSToken:     os.Getenv("ANSIBLE_HUB_SAAS_TOKEN"),
		AAPURL:           os.Getenv("ANSIBLE_AAP_URL"),
		AAPAuthMode:      getEnv("ANSIBLE_AAP_AUTH_MODE", "token"),
		AAPToken:         os.Getenv("ANSIBLE_AAP_TOKEN"),
		AAPUsername:      os.Getenv("ANSIBLE_AAP_USERNAME"),
		AAPPassword:      os.Getenv("ANSIBLE_AAP_PASSWORD"),
		HubTarget:        getEnv("ANSIBLE_HUB_TARGET", "saas"),
		GalaxyURL:        getEnv("ANSIBLE_GALAXY_URL", "https://galaxy.ansible.com/api/v3/"),
		GlobalRateLimit:  getEnv("MCP_RATE_LIMIT_GLOBAL", "10:20"),
		SessionRateLimit: getEnv("MCP_RATE_LIMIT_SESSION", "5:10"),
		DocsBaseURL:      getEnv("ANSIBLE_DOCS_BASE_URL", "https://docs.ansible.com/projects/ansible/latest/"),
		Version:          version,
	}

	timeoutSecs, err := strconv.Atoi(getEnv("ANSIBLE_REQUEST_TIMEOUT", "10"))
	if err != nil {
		return nil, fmt.Errorf("invalid ANSIBLE_REQUEST_TIMEOUT: %w", err)
	}
	c.RequestTimeout = time.Duration(timeoutSecs) * time.Second

	c.MaxRetries, err = strconv.Atoi(getEnv("ANSIBLE_MAX_RETRIES", "3"))
	if err != nil {
		return nil, fmt.Errorf("invalid ANSIBLE_MAX_RETRIES: %w", err)
	}

	switch strings.ToLower(c.HubTarget) {
	case "saas", "aap", "both":
		c.HubTarget = strings.ToLower(c.HubTarget)
	default:
		return nil, fmt.Errorf("invalid ANSIBLE_HUB_TARGET %q: must be saas, aap, or both", c.HubTarget)
	}

	// Ensure DocsBaseURL ends with /
	if !strings.HasSuffix(c.DocsBaseURL, "/") {
		c.DocsBaseURL += "/"
	}

	return c, nil
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
