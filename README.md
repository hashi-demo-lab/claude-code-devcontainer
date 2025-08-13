# 🚀 Claude Code DevContainer & Terraform Module Workflow - WORK IN PROGRESS

A development container configuration for working with Claude Code and Claude GitHub integration, with comprehensive Terraform module development workflows. This repository provides a secure, isolated environment for development with pre-configured tools and settings optimized for AI-assisted coding and infrastructure development.

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

### Building and Publishing the Docker Image

If you want to build and publish your own version of the devcontainer image:

1. Navigate to the `.devcontainer` directory:

   ```bash
   cd .devcontainer
   ```

2. Build the Docker image:

   ```bash
   docker build -t your-dockerhub-username/claude-code-tf-devcontainer:latest .
   ```

3. Push to Docker Hub:

   ```bash
   docker push your-dockerhub-username/claude-code-tf-devcontainer:latest
   ```

The pre-built image is available at: `docker.io/srlynch1/claude-code-tf-devcontainer:latest`

## 🔒 Security Features

This devcontainer includes a sophisticated firewall setup that:

- Restricts outbound traffic to only necessary domains and IPs
- Allows GitHub API/Git operations
- Enables NPM registry access
- Permits connections to Anthropic APIs
- Maintains local network connectivity

## 🧰 Included Tools

### Development Tools

- Git
- GitHub CLI with Claude extension
- Pre-commit

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

## 🏗️ Terraform Module Development Workflow

This repository includes a complete workflow for creating and managing Terraform modules using slash commands and GitHub CLI.

### Quick Start

```bash
# Set your template repository
export TF_TEMPLATE_REPO="your-org/terraform-module-template"

# Create new module
/create-tf-module vpc aws "VPC module for AWS infrastructure"

# Setup and configure
cd terraform-aws-vpc
/setup-tf-module vpc aws "VPC module with subnets and routing"

# Validate and test
/validate-tf-module
/docs-tf-module

# Commit initial version
/commit-tf-module "Initial module implementation"
```

### Available Commands

#### Module Creation
- `/create-tf-module <name> <provider> [description]` - Create public module
- `/create-tf-module-private <name> <provider> [description]` - Create private module
- `/create-tf-module-org <org> <name> <provider> [description]` - Create org module
- `/new-tf-module <name> <provider> [description]` - Complete creation workflow

#### Module Development
- `/setup-tf-module <name> <provider> [description]` - Configure module from template
- `/validate-tf-module` - Format, validate, lint, and security scan
- `/docs-tf-module` - Generate documentation with terraform-docs
- `/test-tf-module` - Run tests and validate examples
- `/commit-tf-module [message]` - Commit and push changes

#### GitHub CLI Aliases
- `gh tf-new <provider> <module>` - Create public module
- `gh tf-pr` - Create pull request with template
- `gh tf-ci` - Run CI workflow
- `gh tf-release <version>` - Create release
- `gh tf-view` - View repository in browser

### Documentation
- **[tf-module-commands.md](.claude/commands/tf-module-commands.md)** - Slash commands for module creation
- **[tf-module-setup.md](.claude/commands/tf-module-setup.md)** - Post-creation setup and automation  
- **[gh-aliases.md](.claude/commands/gh-aliases.md)** - GitHub CLI aliases and shortcuts

### Prerequisites for Module Development
- GitHub CLI authenticated (`gh auth login`)
- Template repository set (`export TF_TEMPLATE_REPO="your-org/template"`)
- Access to your Terraform module template repository

### Example Workflow
```bash
# 1. Create AWS VPC module
/create-tf-module vpc aws "VPC with public/private subnets"

# 2. Enter directory and setup
cd terraform-aws-vpc
/setup-tf-module vpc aws "VPC with public/private subnets"

# 3. Develop module (edit .tf files)

# 4. Validate and document
/validate-tf-module
/docs-tf-module

# 5. Test examples
/test-tf-module

# 6. Commit initial version
/commit-tf-module "Initial VPC module with subnets and routing"

# 7. Create feature branch for development
git checkout -b feature/add-nat-gateway

# 8. After changes, create PR
gh tf-pr

# 9. After merge, create release
gh tf-release 1.0.0
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
