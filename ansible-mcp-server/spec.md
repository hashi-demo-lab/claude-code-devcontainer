# Ansible MCP Server — Implementation Spec

## Overview

A Go-based MCP server that brings the full Ansible content ecosystem into
LLM-assisted playbook authoring. It exposes tools for discovering certified
collections, modules, and roles from Ansible Automation Hub (SaaS and
on-prem AAP), fetching full module parameter documentation from
docs.ansible.com and Ansible Galaxy, scaffolding playbooks, linting with
ansible-lint, and generating Molecule v6+ test cases.

Modeled on the architecture of [terraform-mcp-server](https://github.com/hashicorp/terraform-mcp-server).
Compatible with Claude Desktop, GitHub Copilot, and any MCP-compatible client.

---

## Architecture

### Transport

- **stdio** — standard MCP transport, compatible with all MCP clients
- No HTTP server; the client spawns the process and communicates over stdin/stdout

### MCP Library

- [`github.com/mark3labs/mcp-go`](https://github.com/mark3labs/mcp-go) — Go MCP SDK (same as terraform-mcp-server)

### Project Structure

```
ansible-mcp-server/
├── main.go
├── go.mod
├── go.sum
├── Makefile
├── internal/
│   ├── config/
│   │   └── config.go          # Config schema, env var loading, validation
│   ├── httpclient/
│   │   └── client.go          # go-retryablehttp wrapper, rate limiting, backoff
│   ├── hub/
│   │   ├── saas.go            # Automation Hub SaaS API client
│   │   └── aap.go             # On-prem AAP Hub API client
│   ├── galaxy/
│   │   └── client.go          # Ansible Galaxy API client (fallback)
│   ├── docs/
│   │   ├── fetcher.go         # docs.ansible.com HTML fetcher
│   │   └── parser.go          # HTML parser for module docs, best practices pages
│   └── tools/
│       ├── registry.go        # Tool registration
│       ├── collections.go     # search_collections, get_collection_details
│       ├── modules.go         # search_modules, get_module_docs
│       ├── roles.go           # search_roles, get_role_details
│       ├── bestpractices.go   # get_best_practices, get_playbook_keywords,
│       │                      # get_special_variables, get_return_values
│       ├── scaffold.go        # generate_playbook_scaffold
│       ├── validate.go        # validate_playbook (ansible-lint)
│       └── testcases.go       # generate_test_cases (Molecule v6+)
```

---

## Configuration

All configuration is via environment variables. No config file is required.

### Automation Hub SaaS

| Env var | Description | Default |
|---|---|---|
| `ANSIBLE_HUB_SAAS_URL` | Base URL for SaaS Hub API | `https://cloud.redhat.com/api/automation-hub/v3/` |
| `ANSIBLE_HUB_SAAS_TOKEN` | Offline Bearer token from Red Hat SSO | — |

### On-prem AAP

| Env var | Description | Default |
|---|---|---|
| `ANSIBLE_AAP_URL` | Base URL for on-prem AAP Hub API | — |
| `ANSIBLE_AAP_AUTH_MODE` | `token` or `basic` | `token` |
| `ANSIBLE_AAP_TOKEN` | Bearer token (when `AUTH_MODE=token`) | — |
| `ANSIBLE_AAP_USERNAME` | Username (when `AUTH_MODE=basic`) | — |
| `ANSIBLE_AAP_PASSWORD` | Password (when `AUTH_MODE=basic`) | — |

### Hub Target Selection

| Env var | Description | Default |
|---|---|---|
| `ANSIBLE_HUB_TARGET` | `saas`, `aap`, or `both` | `saas` |

When `both`: SaaS Hub is the primary target; on-prem AAP is used as fallback
if the SaaS call returns no results or fails.

### Galaxy

| Env var | Description | Default |
|---|---|---|
| `ANSIBLE_GALAXY_URL` | Galaxy API base URL | `https://galaxy.ansible.com/api/v3/` |

Galaxy requires no authentication for read operations.

### HTTP Client

| Env var | Description | Default |
|---|---|---|
| `ANSIBLE_REQUEST_TIMEOUT` | Per-request timeout in seconds | `10` |
| `ANSIBLE_MAX_RETRIES` | Max retry attempts per request | `3` |

### Rate Limiting

| Env var | Description | Default |
|---|---|---|
| `MCP_RATE_LIMIT_GLOBAL` | Global rate limit as `rate:burst` | `10:20` |
| `MCP_RATE_LIMIT_SESSION` | Per-session rate limit as `rate:burst` | `5:10` |

### Documentation

| Env var | Description | Default |
|---|---|---|
| `ANSIBLE_DOCS_BASE_URL` | Base URL for Ansible docs | `https://docs.ansible.com/projects/ansible/latest/` |

---

## Authentication

### SaaS Hub (token)

All requests include:

```
Authorization: Bearer <ANSIBLE_HUB_SAAS_TOKEN>
```

The offline token is obtained from the Red Hat Hybrid Cloud Console and does
not expire unless revoked.

### On-prem AAP — token mode

```
Authorization: Bearer <ANSIBLE_AAP_TOKEN>
```

### On-prem AAP — basic auth mode

```
Authorization: Basic base64(<username>:<password>)
```

### OAuth (future)

OAuth 2.0 (Authorization Code + PKCE) support is on the roadmap for both
SaaS Hub and on-prem AAP. Not in scope for v1. When implemented it will be
addable as a third `ANSIBLE_AAP_AUTH_MODE=oauth` option.

### Galaxy

No authentication for read operations. No `Authorization` header is sent.

---

## HTTP Client

Uses [`go-retryablehttp`](https://github.com/hashicorp/go-retryablehttp),
matching terraform-mcp-server.

### Behavior

- **Timeout**: `ANSIBLE_REQUEST_TIMEOUT` seconds per request (default: 10)
- **Retries**: Up to `ANSIBLE_MAX_RETRIES` attempts (default: 3)
- **Retry conditions**: 429, 500, 502, 503, 504
- **Backoff**: On 429, reads `x-ratelimit-reset` header and waits until that
  timestamp before retrying; falls back to exponential backoff if header absent
- **User-Agent**: `ansible-mcp-server/<version>`
- **Proxy**: respects `HTTP_PROXY` / `HTTPS_PROXY` / `NO_PROXY` environment variables

### Error Handling

If all retries are exhausted, the tool returns a structured MCP error with:
- HTTP status code
- Target (saas/aap/galaxy/docs)
- Original error message

The server never panics on API failure — it returns a structured error per
tool call.

---

## docs.ansible.com URL Patterns

All documentation is fetched live from `https://docs.ansible.com/projects/ansible/latest/`.

| Content | URL pattern |
|---|---|
| Collections index | `collections/index.html` |
| All modules & plugins | `collections/all_plugins.html` |
| Individual collection | `collections/<namespace>/<collection>/index.html` |
| Individual module | `collections/<namespace>/<collection>/<module>_module.html` |
| Individual plugin (lookup) | `collections/<namespace>/<collection>/<plugin>_lookup.html` |
| Individual plugin (filter) | `collections/<namespace>/<collection>/<plugin>_filter.html` |
| Playbook keywords | `reference_appendices/playbooks_keywords.html` |
| Return values | `reference_appendices/common_return_values.html` |
| Special variables | `reference_appendices/special_variables.html` |
| YAML syntax | `reference_appendices/YAMLSyntax.html` |
| Ansible config settings | `reference_appendices/config.html` |
| Testing strategies | `reference_appendices/test_strategies.html` |
| Galaxy user guide | `galaxy/user_guide.html` |

The `all_plugins.html` and `collections/index.html` pages are used as
enumerable indexes — they list all known namespaces and modules and can seed
search and autocomplete without crawling every page.

---

## Automation Hub API Endpoints

### SaaS Hub

Base: `https://cloud.redhat.com/api/automation-hub/v3/`

| Operation | Endpoint |
|---|---|
| Search collections | `GET /collections/?keywords=<q>&limit=<n>` |
| Get collection | `GET /collections/<namespace>/<name>/` |
| Get collection version | `GET /collections/<namespace>/<name>/versions/<version>/` |
| List collection content | `GET /collections/<namespace>/<name>/versions/<version>/docs-blob/` |
| Search roles | `GET /roles/?keywords=<q>` |
| Get role | `GET /roles/<namespace>/<name>/` |

### On-prem AAP

Base: `https://<ANSIBLE_AAP_URL>/api/galaxy/v3/`

Endpoints mirror the SaaS Hub v3 API. The AAP client is a separate struct
that reads from `ANSIBLE_AAP_URL` and applies the configured auth.

---

## Ansible Galaxy API Endpoints

Base: `https://galaxy.ansible.com/api/v3/`

Used as fallback for full module parameter documentation when Hub does not
return docs-blob content.

| Operation | Endpoint |
|---|---|
| Search collections | `GET /collections/?keywords=<q>` |
| Get collection | `GET /collections/<namespace>/<name>/` |
| Get collection version docs | `GET /collections/<namespace>/<name>/versions/<version>/docs-blob/` |
| Search roles | `GET /roles/?keywords=<q>` |
| Get role | `GET /roles/<id>/` |

---

## Tool Inventory

### `search_collections`

Search Automation Hub (or Galaxy fallback) for collections matching a
keyword or task description.

**Inputs**

| Field | Type | Required | Description |
|---|---|---|---|
| `query` | string | yes | Keyword or task description |
| `namespace` | string | no | Filter by namespace (e.g., `ansible`, `community`) |
| `limit` | integer | no | Max results to return (default: 10, max: 50) |
| `certified_only` | boolean | no | Return only Hub-certified content (default: true) |

**Output** — list of collections, each with:
- `namespace`, `name`, `version`, `description`
- `support_level` (certified / community / partner)
- `hub_url`, `docs_url`
- `source` (saas / aap / galaxy)

---

### `get_collection_details`

Get full metadata and content listing for a specific collection.

**Inputs**

| Field | Type | Required | Description |
|---|---|---|---|
| `namespace` | string | yes | Collection namespace |
| `name` | string | yes | Collection name |
| `version` | string | no | Specific version (default: latest) |

**Output**
- Metadata: namespace, name, version, description, license, authors, dependencies
- Modules list with short descriptions
- Roles list
- Plugins list (lookup, filter, callback, inventory)
- `docs_url`: `<ANSIBLE_DOCS_BASE_URL>collections/<namespace>/<name>/index.html`
- `source` (saas / aap / galaxy)

---

### `search_modules`

Find modules matching a task description across Hub and docs.ansible.com.

**Inputs**

| Field | Type | Required | Description |
|---|---|---|---|
| `query` | string | yes | Task description or keyword |
| `namespace` | string | no | Filter by namespace |
| `collection` | string | no | Filter by collection |
| `limit` | integer | no | Max results (default: 10, max: 50) |

**Output** — list of modules, each with:
- `fqcn` (fully qualified collection name, e.g., `ansible.posix.firewalld`)
- `namespace`, `collection`, `module_name`
- `short_description`
- `docs_url`: constructed from `all_plugins.html` index
- `source` (hub / galaxy / docs)

---

### `get_module_docs`

Get full documentation for a specific module including all parameters.

**Inputs**

| Field | Type | Required | Description |
|---|---|---|---|
| `fqcn` | string | yes | Fully qualified collection name (e.g., `ansible.posix.firewalld`) |

Alternatively accepts separate fields:

| Field | Type | Required | Description |
|---|---|---|---|
| `namespace` | string | yes | Module namespace |
| `collection` | string | yes | Collection name |
| `module` | string | yes | Module name |

**Output**
- `fqcn`, `short_description`, `version_added`
- `parameters`: array of parameter objects, each with:
  - `name`, `type`, `required` (bool), `default`, `choices` (enum values)
  - `description`, `aliases`
  - `suboptions` (for dict/list parameters with nested keys)
- `examples`: YAML task examples
- `return_values`: what the module returns
- `notes`, `seealso`
- `docs_url`: `<ANSIBLE_DOCS_BASE_URL>collections/<namespace>/<collection>/<module>_module.html`
- `source` (docs / galaxy-fallback)

**Fallback logic**: If docs.ansible.com does not have the module page (e.g., it
is a community collection not mirrored there), the tool fetches the
docs-blob from Galaxy API to retrieve full parameter documentation.

---

### `search_roles`

Search for roles in Automation Hub and Ansible Galaxy.

**Inputs**

| Field | Type | Required | Description |
|---|---|---|---|
| `query` | string | yes | Task description or keyword |
| `namespace` | string | no | Filter by namespace |
| `limit` | integer | no | Max results (default: 10) |
| `source` | string | no | `hub`, `galaxy`, or `both` (default: `both`) |

**Output** — list of roles, each with:
- `namespace`, `name`, `description`
- `platforms`: supported OS/versions
- `source` (hub / galaxy)
- `url`

---

### `get_role_details`

Get full details for a specific role including variables and dependencies.

**Inputs**

| Field | Type | Required | Description |
|---|---|---|---|
| `namespace` | string | yes | Role namespace |
| `name` | string | yes | Role name |
| `source` | string | no | `hub` or `galaxy` (default: auto-detect from namespace) |

**Output**
- `namespace`, `name`, `description`, `version`
- `platforms`: supported OS/version matrix
- `variables`: role defaults with types and descriptions
- `dependencies`: list of required roles or collections
- `examples`: usage examples
- `readme`: extracted from Galaxy/Hub if available
- `source`, `url`

---

### `get_best_practices`

Retrieve best practices for an Ansible topic from docs.ansible.com.

**Inputs**

| Field | Type | Required | Description |
|---|---|---|---|
| `topic` | string | yes | Topic to look up (e.g., `idempotency`, `error handling`, `variable naming`, `handlers`, `tags`) |

**Output**
- `topic`
- `content`: extracted best practice content from the relevant docs page
- `source_url`: the specific docs.ansible.com URL fetched
- `related_keywords_reference`: link to Playbook Keywords page if applicable

**Topic-to-URL mapping** (representative):

| Topic keywords | Docs page |
|---|---|
| playbook structure, tasks, plays | `playbook_guide/playbooks_intro.html` |
| variables, precedence | `playbook_guide/playbooks_variables.html` |
| handlers | `playbook_guide/playbooks_handlers.html` |
| tags | `playbook_guide/playbooks_tags.html` |
| error handling, ignore_errors, rescue | `playbook_guide/playbooks_error_handling.html` |
| idempotency, state | `reference_appendices/test_strategies.html` |
| roles, role structure | `playbook_guide/playbooks_reuse_roles.html` |
| collections, using collections | `collections_guide/index.html` |
| loops, with_items | `playbook_guide/playbooks_loops.html` |
| conditionals, when | `playbook_guide/playbooks_conditionals.html` |
| templates, Jinja2 | `playbook_guide/playbooks_templating.html` |
| vault, secrets | `vault_guide/index.html` |
| inventory | `inventory_guide/index.html` |
| galaxy, collections install | `galaxy/user_guide.html` |

---

### `get_playbook_keywords`

Retrieve the Ansible playbook keyword reference.

**Inputs**

| Field | Type | Required | Description |
|---|---|---|---|
| `keyword` | string | no | Filter to a specific keyword (e.g., `become`, `delegate_to`) |

**Output**
- If `keyword` supplied: definition, type, default, scope (play/task/role), description
- If no keyword: full keyword reference list
- `source_url`: `<ANSIBLE_DOCS_BASE_URL>reference_appendices/playbooks_keywords.html`

---

### `get_special_variables`

Retrieve Ansible special (magic) variables.

**Inputs**

| Field | Type | Required | Description |
|---|---|---|---|
| `variable` | string | no | Filter to a specific variable (e.g., `inventory_hostname`, `ansible_facts`) |

**Output**
- Variable name, type, description, scope
- `source_url`: `<ANSIBLE_DOCS_BASE_URL>reference_appendices/special_variables.html`

---

### `generate_playbook_scaffold`

Generate a boilerplate Ansible playbook following best practices.

**Inputs**

| Field | Type | Required | Description |
|---|---|---|---|
| `task_description` | string | yes | What the playbook should accomplish |
| `target_os` | string | no | Target OS family (e.g., `rhel`, `debian`, `windows`) |
| `collections` | []string | no | Collections to declare (e.g., `["ansible.posix", "community.general"]`) |
| `use_roles` | boolean | no | Structure as a role invocation (default: false) |
| `style` | string | no | `simple` (single file) or `project` (directory layout with roles/) (default: `simple`) |

**Output**
- `playbook`: YAML content of the generated playbook
- If `style=project`: directory tree with file paths and content for each file
- `notes`: list of best practices applied (e.g., "become used only where needed", "handlers used for service restarts")

**Best practices enforced in scaffold**:
- `gather_facts: true` unless explicitly not needed
- `become` applied at task level, not play level, unless the whole play requires it
- Handlers for service restarts (not `service` module in a task after config change)
- Variables in `vars:` or referenced from inventory, not hardcoded in tasks
- `state:` always explicit on modules that support it
- `name:` on every task and play, using descriptive imperative sentences
- `no_log: true` on tasks that handle secrets

---

### `validate_playbook`

Lint and validate a playbook using ansible-lint.

**Inputs**

| Field | Type | Required | Description |
|---|---|---|---|
| `playbook` | string | yes | YAML playbook content to validate |
| `profile` | string | no | ansible-lint profile: `min`, `basic`, `moderate`, `safety`, `shared`, `production` (default: `basic`) |

**Output**
- `passed`: boolean
- `violations`: array of objects, each with:
  - `rule_id`: ansible-lint rule ID (e.g., `command-instead-of-module`)
  - `description`: what the rule checks
  - `line`: line number in the submitted playbook
  - `severity`: `warning` / `error`
  - `remediation`: specific corrective suggestion for this violation
- `summary`: total violations by severity

**Runtime dependency**: `ansible-lint` must be installed and available on
`$PATH` of the server process. If not found, the tool returns a structured
error with installation instructions.

---

### `generate_test_cases`

Generate Molecule v6+ test scenarios for an Ansible playbook or role.

**Inputs**

| Field | Type | Required | Description |
|---|---|---|---|
| `playbook` | string | no | Playbook YAML content to generate tests for |
| `task_description` | string | no | Task description (used if playbook not provided) |
| `role_name` | string | no | Role name (for role-based test structure) |
| `driver` | string | no | Molecule driver: `docker`, `podman`, `delegated` (default: `docker`) |
| `platforms` | []object | no | Target platforms (see below) |

**Platform object**:
```json
{
  "name": "rhel9",
  "image": "registry.access.redhat.com/ubi9/ubi-init",
  "pre_build_image": true
}
```

Default platform if not specified:
```json
{
  "name": "instance",
  "image": "registry.access.redhat.com/ubi9/ubi-init",
  "pre_build_image": true
}
```

**Output** — Molecule v6+ directory structure:

```
molecule/
└── default/
    ├── molecule.yml        # driver, platforms, provisioner, verifier config
    ├── converge.yml        # playbook that applies the role/tasks under test
    ├── verify.yml          # assertions that validate the converged state
    ├── prepare.yml         # optional: pre-converge setup tasks
    └── cleanup.yml         # optional: post-test teardown tasks
```

Each file is returned as a path + YAML content string.

**Scenarios generated**:
1. **Default**: full converge + verify
2. **Idempotency**: converge run twice; second run must report zero changes
   (this maps to terraform-mcp-server's equivalent of plan showing no diff
   after apply)

**molecule.yml structure**:
```yaml
dependency:
  name: galaxy
driver:
  name: <driver>
platforms:
  - name: <platform.name>
    image: <platform.image>
    pre_build_image: <platform.pre_build_image>
provisioner:
  name: ansible
  config_options:
    defaults:
      interpreter_python: auto_silent
verifier:
  name: ansible
```

**verify.yml approach**: Uses `ansible.builtin.assert` tasks to check
expected state (service running, file present, package installed, etc.)
inferred from the converge playbook content.

---

## Multi-client Compatibility

The server uses stdio transport and standard MCP protocol. No client-specific
code is required.

### Claude Desktop (`claude_desktop_config.json`)

```json
{
  "mcpServers": {
    "ansible": {
      "command": "/path/to/ansible-mcp-server",
      "env": {
        "ANSIBLE_HUB_SAAS_TOKEN": "...",
        "ANSIBLE_HUB_TARGET": "saas"
      }
    }
  }
}
```

### VS Code (`.vscode/mcp.json` or user settings)

```json
{
  "servers": {
    "ansible": {
      "type": "stdio",
      "command": "/path/to/ansible-mcp-server",
      "env": {
        "ANSIBLE_HUB_SAAS_TOKEN": "...",
        "ANSIBLE_HUB_TARGET": "saas"
      }
    }
  }
}
```

### GitHub Copilot (VS Code MCP extension)

Uses the same VS Code config format above. No additional changes required.

---

## Dependencies

### Go modules

| Module | Purpose |
|---|---|
| `github.com/mark3labs/mcp-go` | MCP server SDK (stdio transport, tool registry) |
| `github.com/hashicorp/go-retryablehttp` | Retryable HTTP client with backoff |
| `golang.org/x/net/html` | HTML parsing for docs.ansible.com pages |
| `golang.org/x/time/rate` | Token bucket rate limiter |

### Runtime (must be on `$PATH`)

| Tool | Required for | Install |
|---|---|---|
| `ansible-lint` | `validate_playbook` tool | `pip install ansible-lint` |

All other tools (search, docs, scaffold, Molecule generation) have no
runtime dependencies beyond the Go binary.

---

## Build and Run

```makefile
# Makefile targets

build:
    go build -o ansible-mcp-server ./...

test:
    go test ./...

lint:
    golangci-lint run

install:
    go install ./...
```

### Running locally

```bash
ANSIBLE_HUB_SAAS_TOKEN=<token> \
ANSIBLE_HUB_TARGET=saas \
./ansible-mcp-server
```

### Running against on-prem AAP

```bash
ANSIBLE_HUB_TARGET=aap \
ANSIBLE_AAP_URL=https://aap.internal \
ANSIBLE_AAP_AUTH_MODE=token \
ANSIBLE_AAP_TOKEN=<token> \
./ansible-mcp-server
```

---

## Error Handling Contract

All tools return one of two shapes:

**Success**:
```json
{
  "content": [{ "type": "text", "text": "<structured result>" }]
}
```

**Error**:
```json
{
  "isError": true,
  "content": [{ "type": "text", "text": "ERROR [<source>]: <message>" }]
}
```

Error messages include the source (`hub-saas`, `hub-aap`, `galaxy`, `docs`,
`ansible-lint`) so the caller knows which integration failed.

---

## Validation Scenarios

| Scenario | Method | Pass condition |
|---|---|---|
| Hub search accuracy | Call `search_collections` with known keyword | Returns certified collections with correct metadata from live API |
| Dual-target auth | Configure SaaS, then AAP independently | Same tool calls return correct results from both |
| Full module params | `get_module_docs` for `ansible.posix.firewalld` | All params, types, defaults, choices, examples present |
| Galaxy fallback | `get_module_docs` for a community module not on Hub | Params fetched from Galaxy docs-blob |
| Docs URL construction | `get_best_practices` for any topic | Response includes a valid, live `docs.ansible.com` URL |
| Scaffold validity | `generate_playbook_scaffold` for any task | Output passes `ansible-lint` + `ansible-playbook --syntax-check` with zero errors |
| Molecule structure | `generate_test_cases` for any playbook | Output matches Molecule v6+ layout; `converge.yml` is valid YAML; idempotency scenario present |
| Lint with remediation | `validate_playbook` with `shell` used where module exists | Returns `command-instead-of-module` rule ID + remediation hint |
| Certified preference | Search for content in both Hub and Galaxy | Hub certified version is primary; Galaxy is labeled fallback |
| Multi-client load | Load server in Claude Desktop + VS Code MCP | All tools discoverable in both clients |
| Rate limit resilience | Simulate 429 from Hub | Server waits on `x-ratelimit-reset`, retries up to 3x, returns structured error |
| ansible-lint missing | `validate_playbook` with no `ansible-lint` on PATH | Structured error with install instructions, no panic |

---

## Future Work

| Item | Notes |
|---|---|
| OAuth 2.0 auth | Authorization Code + PKCE for SaaS Hub and AAP; addable as `ANSIBLE_AAP_AUTH_MODE=oauth` without breaking existing token/basic config |
| Response caching | Optional TTL-based in-memory cache for Hub API responses; off by default, opt-in via env var |
| AAP Controller tools | Job template listing, inventory listing — extends the server into AAP orchestration, not just Hub content discovery |
| Public open source release | Licensing, contribution guide, and public docs once internal version is stable |
| additional Molecule drivers | `vagrant`, `ec2` driver support in `generate_test_cases` |
| ansible-navigator integration | Alternative to ansible-lint using `ansible-navigator lint` for EE-based environments |
