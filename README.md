# 🚀 Claude Code DevContainer

A development container configuration for working with Claude Code and Claude GitHub integration. This repository provides a secure, isolated environment for development with pre-configured tools and settings optimized for AI-assisted coding.

## 🛠️ Setup

### Prerequisites

- [Docker](https://www.docker.com/products/docker-desktop/)
- [VS Code](https://code.visualstudio.com/)
- [VS Code Remote - Containers Extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)
- [GitHub Personal Access Token](https://github.com/settings/tokens) with appropriate permissions

### Getting Started

1. Clone this repository
2. Open in VS Code
3. When prompted, click "Reopen in Container"
4. Alternatively, press F1 and select "Remote-Containers: Reopen in Container"

## 🔒 Security Features

This devcontainer includes a sophisticated firewall setup that:

- Restricts outbound traffic to only necessary domains and IPs
- Allows GitHub API/Git operations
- Enables NPM registry access
- Permits connections to Anthropic APIs
- Maintains local network connectivity

## 🧰 Included Tools

### Development Tools

- Node.js 20
- Git with Delta for improved diffs
- GitHub CLI with Claude extension
- Python 3 with pip

### Terraform Tools

- Terraform 1.12.1
- Terraform Docs
- TFSec
- Terrascan
- TFLint with AWS, Azure, and GCP rulesets
- Infracost
- Checkov

### VS Code Extensions

- ESLint
- Prettier
- Terraform
- Claude Code
- YAML
- Azure GitHub Copilot
- HashiCorp HCL

## 🔄 Environment Variables

| Variable | Description |
|----------|-------------|
| `GH_TOKEN` | GitHub Personal Access Token for authentication |
| `NODE_OPTIONS` | Node.js memory allocation (4GB by default) |
| `CLAUDE_CONFIG_DIR` | Location of Claude configuration |
| `CLAUDE_GITHUB_MCP_ENABLED` | Enables Claude GitHub integration |

## � Usage

The container automatically sets up the necessary environment for working with:

1. Claude Code in VS Code
2. GitHub CLI integration with Claude
3. Terraform projects with security scanning tools
4. Network-isolated development environment

## � MCP Servers

This devcontainer supports Model Context Protocol (MCP) servers, including Docker-based MCP servers.

### Using Docker-based MCP Servers

The devcontainer is configured with:

- Docker CLI installed
- Docker socket mounted from the host
- Support for running Docker commands inside the container

A sample Terraform MCP server is included in the `.mcp.json` file at the project root:

```json
{
  "mcpServers": {
    "terraform": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "hashicorp/terraform-mcp-server"
      ]
    }
  }
}
```

### Adding Custom MCP Servers

To add your own MCP servers:

1. Edit the `.mcp.json` file at the project root
2. Add your server configuration following the MCP protocol
3. Claude Code will automatically detect and prompt for approval

For more information, see the [Claude Code MCP documentation](https://docs.anthropic.com/en/docs/claude-code/mcp).

## �🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
