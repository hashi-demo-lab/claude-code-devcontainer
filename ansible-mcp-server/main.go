package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hashi-demo-lab/ansible-mcp-server/internal/config"
	"github.com/hashi-demo-lab/ansible-mcp-server/internal/docs"
	"github.com/hashi-demo-lab/ansible-mcp-server/internal/galaxy"
	"github.com/hashi-demo-lab/ansible-mcp-server/internal/hub"
	"github.com/hashi-demo-lab/ansible-mcp-server/internal/httpclient"
	"github.com/hashi-demo-lab/ansible-mcp-server/internal/tools"
	"github.com/mark3labs/mcp-go/server"
)

// version is set at build time via -ldflags="-X main.version=..."
var version = "dev"

func main() {
	cfg, err := config.Load(version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ansible-mcp-server: configuration error: %v\n", err)
		os.Exit(1)
	}

	httpClient, err := httpclient.New(
		cfg.RequestTimeout,
		cfg.MaxRetries,
		cfg.GlobalRateLimit,
		cfg.SessionRateLimit,
		cfg.Version,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ansible-mcp-server: HTTP client error: %v\n", err)
		os.Exit(1)
	}

	// Create docs fetcher
	docsFetcher := docs.NewFetcher(httpClient, cfg.DocsBaseURL)

	// Create Galaxy client (always available as fallback)
	galaxyClient := galaxy.NewClient(httpClient, cfg.GalaxyURL)

	// Create Hub clients based on target configuration
	deps := &tools.Dependencies{
		DocsBaseURL: cfg.DocsBaseURL,
		DocsFetcher: docsFetcher,
		Galaxy:      galaxyClient,
		HubTarget:   cfg.HubTarget,
	}

	switch cfg.HubTarget {
	case "saas", "both":
		if cfg.HubSaaSToken != "" {
			deps.SaaSHub = hub.NewSaaSClient(httpClient, cfg.HubSaaSURL, cfg.HubSaaSToken)
		} else if cfg.HubTarget == "saas" {
			log.Printf("WARNING: ANSIBLE_HUB_TARGET=saas but ANSIBLE_HUB_SAAS_TOKEN is not set; Hub searches will fail")
		}
		fallthrough
	case "aap":
		if cfg.HubTarget != "saas" && cfg.AAPURL != "" {
			deps.AAPHub = hub.NewAAPClient(httpClient, cfg.AAPURL, cfg.AAPAuthMode, cfg.AAPToken, cfg.AAPUsername, cfg.AAPPassword)
		}
	}

	// Build and configure the MCP server
	s := server.NewMCPServer(
		"ansible-mcp-server",
		cfg.Version,
		server.WithToolCapabilities(false),
	)

	// Register all 12 tools
	tools.RegisterAll(s, deps)

	// Serve over stdio (MCP standard transport)
	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "ansible-mcp-server: %v\n", err)
		os.Exit(1)
	}
}
