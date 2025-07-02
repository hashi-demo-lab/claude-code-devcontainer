# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

# Repository Overview

This is a DevContainer configuration repository for Claude Code development that provides a secure, isolated Docker-based development environment optimized for AI-assisted coding and Terraform infrastructure development. The repository includes a complete workflow for creating and managing Terraform modules using slash commands and GitHub CLI.

# Architecture

## Core Components

- **DevContainer Setup**: `.devcontainer/devcontainer.json` contains Docker configuration with Node.js 20 base image, privileged container access for network management
- **MCP Server Integration**: `.mcp.json` configures four Model Context Protocol servers for enhanced Claude Code functionality
- **Terraform Module Development**: Complete workflow with slash commands for module creation, validation, testing, and deployment
- **Security Layer**: Network firewall restrictions for isolated development (commented out in current Dockerfile)
- **Tool Chain**: Pre-installed Terraform ecosystem, GitHub CLI, and development utilities

## MCP Servers Configuration

The repository includes four configured MCP servers in `.mcp.json`:

1. **Terraform MCP** (`hashicorp/terraform-mcp-server`): Docker-based server for Terraform operations
2. **AWS Labs Terraform MCP** (`awslabs.terraform-mcp-server`): Advanced Terraform server with AWS integration via uvx
3. **Context7 MCP** (`@upstash/context7-mcp`): Documentation context server via npx
4. **Sequential Thinking MCP** (`@modelcontextprotocol/server-sequential-thinking`): Problem-solving assistance server

## Development Environment

- **Base**: Node.js 20 with zsh shell and development tools
- **User**: Runs as `node` user with sudo privileges
- **Workspace**: Mounted at `/workspace` with persistent volumes for bash history and Claude config
- **Extensions**: VS Code configured with ESLint, Prettier, Terraform, Claude Code, Azure GitHub Copilot, and HashiCorp HCL
- **AWS Integration**: Environment variables configured for AWS access key propagation from host

# Common Development Commands

## Terraform Module Development Workflow

The repository includes comprehensive slash commands for Terraform module development:

### Module Creation Commands
- `/create-tf-module <name> <provider> [description]` - Create public module
- `/create-tf-module-private <name> <provider> [description]` - Create private module  
- `/create-tf-module-org <org> <name> <provider> [description]` - Create org module
- `/new-tf-module <name> <provider> [description]` - Complete creation workflow

### Module Development Commands
- `/setup-tf-module <name> <provider> [description]` - Configure module from template
- `/validate-tf-module` - Format, validate, lint, and security scan with all tools
- `/docs-tf-module` - Generate documentation with terraform-docs
- `/test-tf-module` - Run tests and validate examples
- `/commit-tf-module [message]` - Commit and push changes

### GitHub CLI Aliases
- `gh tf-new <provider> <module>` - Create public module
- `gh tf-pr` - Create pull request with template
- `gh tf-ci` - Run CI workflow
- `gh tf-release <version>` - Create release
- `gh tf-view` - View repository in browser

## Terraform Validation and Security

Use `/validate-tf-module` command which runs:
- `terraform fmt -recursive` for formatting
- `terraform validate` for syntax validation
- `tflint` with AWS/Azure/GCP rulesets for best practices
- `tfsec` for security scanning
- `terrascan` for policy scanning
- `checkov` for compliance checking

# Pre-installed Tools

## Infrastructure Tools
- Terraform 1.12.1
- Terraform Docs 0.20.0  
- TFSec 1.28.13 (security scanner)
- Terrascan 1.19.9 (security scanner)
- TFLint 0.48.0 with AWS/Azure/GCP rulesets
- Infracost 0.10.41 (cost analysis)
- Checkov 3.2.439 (security scanner)

## Development Tools
- Git with Delta for enhanced diffs
- GitHub CLI with Claude extension (via post-create script)
- Python 3 with pip, venv support
- fzf, jq, and standard Unix utilities

# Environment Variables

- `NODE_OPTIONS`: Set to `--max-old-space-size=4096` for 4GB memory allocation
- `CLAUDE_CONFIG_DIR`: Points to `/home/node/.claude` for persistent Claude configuration
- `CLAUDE_GITHUB_MCP_ENABLED`: Enables Claude GitHub MCP integration
- `GITHUB_TOKEN`: GitHub Personal Access Token (propagated from host)
- `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_SESSION_TOKEN`, `AWS_SESSION_EXPIRATION`: AWS credentials (propagated from host)

# Key Features

## Security
- Privileged container with NET_ADMIN/NET_RAW capabilities for firewall management
- Network isolation capabilities (firewall scripts available but currently disabled)
- Comprehensive security scanning with multiple tools (TFSec, Terrascan, Checkov)

## Docker Integration
- Docker-in-Docker feature enabled for running containerized MCP servers
- Docker socket mounted from host for Docker command execution
- Supports both local and containerized development workflows

## VS Code Integration
- Format on save enabled with Prettier
- ESLint auto-fix on save
- Terraform and HCL syntax support with HashiCorp extensions
- Terminal defaults to zsh with enhanced prompt

# Setup Requirements

1. Docker Desktop
2. VS Code with Remote-Containers extension  
3. GitHub Personal Access Token (for GitHub CLI integration)
4. Template repository configured (`export TF_TEMPLATE_REPO="your-org/terraform-module-template"`)
5. Open repository in VS Code and select "Reopen in Container"

# Post-Creation Setup

The `post-create.sh` script automatically:
- Authenticates GitHub CLI if `GITHUB_TOKEN` is provided
- Installs Claude GitHub extension (`anthropic/gh-claude`)
- Provides manual setup instructions if token is not available