# Ansible MCP Server — Claude Code Context

## What this project is

A Go MCP server that brings the full Ansible content ecosystem into
LLM-assisted playbook authoring. Users interact with it through any
MCP-compatible client (Claude Desktop, GitHub Copilot, VS Code extensions).

**Modeled on**: `github.com/hashicorp/terraform-mcp-server` — read its
source when you need architectural patterns. Use the `github` MCP server
to read it directly.

**Full implementation spec**: `spec.md` in this repo root. That is the
authoritative reference for tools, API endpoints, config schema, auth, and
Molecule test generation. Read it before writing any code.

---

## Key decisions (do not re-litigate these)

| Decision | Choice | Reason |
|---|---|---|
| Language | Go | User requirement |
| MCP library | `github.com/mark3labs/mcp-go` | Same as terraform-mcp-server |
| Transport | stdio | Multi-client compatibility |
| HTTP client | `go-retryablehttp` | Same as terraform-mcp-server; 3 retries, 10s timeout, 429-aware backoff |
| Caching | None | Stateless, same as terraform-mcp-server; live queries only |
| Rate limiting | `MCP_RATE_LIMIT_GLOBAL` / `MCP_RATE_LIMIT_SESSION` env vars | Same pattern as terraform-mcp-server |
| Hub target | `ANSIBLE_HUB_TARGET=saas\|aap\|both` | Config-driven, no restart |
| Auth (v1) | Token + Basic auth | OAuth is future work (see spec.md) |
| Galaxy | Fallback only | Used when Hub does not have full module parameter docs |
| Docs | `docs.ansible.com` — live fetch, deterministic URLs | No caching, no scraping index needed |
| Molecule | v6+ (collection-based) | User requirement |
| Distribution | Internal for now | Public open source is future work |

---

## Integration targets

### Automation Hub SaaS
- Base URL: `https://cloud.redhat.com/api/automation-hub/v3/`
- Auth: Bearer token (`ANSIBLE_HUB_SAAS_TOKEN`)

### On-prem AAP (2.x+)
- Base URL: `https://<host>/api/galaxy/v3/`
- Auth: token (`ANSIBLE_AAP_TOKEN`) or basic (`ANSIBLE_AAP_USERNAME` + `ANSIBLE_AAP_PASSWORD`)
- Set `ANSIBLE_AAP_AUTH_MODE=token` or `basic`

### Ansible Galaxy (fallback)
- Base URL: `https://galaxy.ansible.com/api/v3/`
- No auth required for reads

### docs.ansible.com
- Base: `https://docs.ansible.com/projects/ansible/latest/`
- Key entry points: `collections/index.html`, `collections/all_plugins.html`
- Module URL pattern: `collections/<namespace>/<collection>/<module>_module.html`
- All URLs are deterministic — no crawling or sitemap needed

---

## Tools to implement (12 total — see spec.md for full detail)

In this order (lowest to highest external dependency):

1. `get_playbook_keywords` — pure HTTP fetch, single known URL
2. `get_special_variables` — same pattern, validates fetcher is reusable
3. `get_module_docs` (docs.ansible.com path) — adds HTML parsing
4. `search_collections` — first Hub API call, validates auth + retryablehttp
5. `get_collection_details` — builds on search
6. `search_modules` — combines docs index + Hub
7. `get_module_docs` (Galaxy fallback path) — adds second API client
8. `search_roles` / `get_role_details` — Galaxy + Hub, same patterns
9. `get_best_practices` — topic-to-URL routing logic
10. `generate_playbook_scaffold` — pure generation, no external calls
11. `validate_playbook` — adds subprocess (ansible-lint on $PATH)
12. `generate_test_cases` — Molecule v6+ directory output

---

## Project structure

```
ansible-mcp-server/
├── main.go
├── go.mod
├── go.sum
├── spec.md                    ← full implementation reference
├── CLAUDE.md                  ← this file
├── Makefile
├── internal/
│   ├── config/config.go       ← env var schema and loading
│   ├── httpclient/client.go   ← go-retryablehttp wrapper + rate limiting
│   ├── hub/
│   │   ├── saas.go            ← SaaS Hub API client
│   │   └── aap.go             ← on-prem AAP client
│   ├── galaxy/client.go       ← Galaxy API client
│   ├── docs/
│   │   ├── fetcher.go         ← docs.ansible.com HTTP fetcher
│   │   └── parser.go          ← HTML parser for module docs pages
│   └── tools/
│       ├── registry.go
│       ├── collections.go
│       ├── modules.go
│       ├── roles.go
│       ├── bestpractices.go
│       ├── scaffold.go
│       ├── validate.go
│       └── testcases.go
```

---

## Development workflow

This project uses the devcontainer from `hashi-demo-lab/claude-code-devcontainer`
(forked and adapted for Go). Use these commands:

| Command | What it does |
|---|---|
| `/build-mcp` | `go build ./...` + `go vet ./...` |
| `/test-mcp` | `go test ./... -v -race -count=1` |
| `/lint-mcp` | `golangci-lint run` + `gosec ./...` |
| `/validate-mcp` | All three in sequence; stops on first failure |
| `/create-mcp-tool <name>` | Scaffold tool file + test file from spec.md pattern |

**Loop pattern**: `/loop /validate-mcp` — implement a tool, validate, fix
failures, loop until green, then commit before moving to the next tool.

### Per-tool workflow (test-first)

```
1. Write input/output structs and handler signature
2. Write test file (table-driven, error cases included)
3. /test-mcp → confirm failure (expected, no implementation yet)
4. Implement handler
5. /validate-mcp → must be green before next tool
6. Commit
```

---

## What NOT to do

- Do not add response caching — the design is intentionally stateless
- Do not mock the Hub or Galaxy APIs in tests — use recorded HTTP fixtures
  or a test server that replays real responses; mocks hide real API shape
- Do not add error handling for scenarios that cannot happen — trust Go
  types and framework guarantees; only validate at system boundaries
- Do not add OAuth in v1 — it is explicitly future work
- Do not support Ansible Tower — AAP 2.x+ only
- Do not add playbook execution, inventory management, or AAP workflow tools
  — out of scope for this server

---

## Runtime dependency

`ansible-lint` must be on `$PATH` for `validate_playbook` to work.
If not found, return a structured MCP error with installation instructions
(`pip install ansible-lint`). Do not panic.

---

## MCP servers available in this devcontainer

- `sequential-thinking` — use for complex multi-step reasoning during
  implementation (e.g., Galaxy fallback logic, HTML parser design)
- `context7` — use to look up live API docs for mcp-go, go-retryablehttp,
  golang.org/x/net/html before writing code against them
- `github` — use to read `hashicorp/terraform-mcp-server` source directly
  when you need to match its architectural patterns exactly

---

## Validation checklist (before marking any tool complete)

- [ ] `go build ./...` passes
- [ ] `go test ./... -race` passes (including error path tests)
- [ ] `golangci-lint run` passes
- [ ] Tool registered in `tools/registry.go`
- [ ] Tool description is accurate (MCP clients show this to users)
- [ ] Error returns are structured MCP errors, not panics
- [ ] Source URL included in all doc/hub responses
- [ ] Committed
