# CLAUDE.md — Ansible MCP Server Devcontainer

This devcontainer is configured for building the Ansible MCP Server in Go.

## Project context

Full project context, architecture decisions, tool inventory, and implementation
instructions are in `/workspace/ansible-mcp-server/CLAUDE.md`.

Full implementation spec is in `/workspace/ansible-mcp-server/spec.md`.

Read both files at the start of any session before writing code.

## MCP servers available

- **sequential-thinking** — use for complex multi-step reasoning (e.g. fallback
  logic, HTML parser design, auth flow)
- **context7** — use to look up live API documentation for Go libraries before
  writing code against them: `mcp-go`, `go-retryablehttp`, `golang.org/x/net/html`
- **github** (via CLAUDE_GITHUB_MCP_ENABLED) — use to read
  `hashicorp/terraform-mcp-server` source when you need to match its
  architectural patterns

## Slash commands

| Command | What it does |
|---|---|
| `/build-mcp` | `go build ./...` + `go vet ./...` |
| `/test-mcp` | `go test ./... -v -race` |
| `/lint-mcp` | `golangci-lint` + `gosec` |
| `/validate-mcp` | All three in sequence; stops on first failure |
| `/create-mcp-tool <name>` | Scaffold tool file + test file from spec.md |

## Development loop

```
/loop /validate-mcp
```

Implement one tool, run `/validate-mcp`, fix failures, loop until green, commit.
