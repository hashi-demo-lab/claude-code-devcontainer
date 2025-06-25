# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

# Repository Overview

This is a DevContainer configuration repository for Claude Code development. It provides a secure, isolated Docker-based development environment with pre-configured tools optimized for AI-assisted coding and infrastructure development.

# Architecture

## Core Components

- **DevContainer Setup**: `.devcontainer/` contains Docker configuration with Node.js 20 base image
- **MCP Server Integration**: `.mcp.json` configures Model Context Protocol servers for enhanced Claude Code functionality
- **Security Layer**: Network firewall restrictions for isolated development (commented out in current Dockerfile)
- **Tool Chain**: Pre-installed Terraform ecosystem, GitHub CLI, and development utilities

## MCP Servers

The repository includes three configured MCP servers:

1. **Terraform MCP** (`hashicorp/terraform-mcp-server`): Docker-based server for Terraform operations
2. **AWS Labs Terraform MCP** (`awslabs.terraform-mcp-server`): Advanced Terraform server with AWS integration via uvx
3. **Context7 MCP** (`@upstash/context7-mcp`): Documentation context server via npx

## Development Environment

- **Base**: Node.js 20 with zsh shell and development tools
- **User**: Runs as `node` user with sudo privileges
- **Workspace**: Mounted at `/workspace` with persistent volumes for bash history and Claude config
- **Extensions**: VS Code configured with ESLint, Prettier, Terraform, Claude Code, and Azure GitHub Copilot

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
- `GH_TOKEN`: GitHub Personal Access Token (set externally for authentication)

# Key Features

## Security
- Privileged container with NET_ADMIN/NET_RAW capabilities for firewall management
- Network isolation capabilities (firewall scripts available but currently disabled)
- Restricted outbound traffic to necessary domains only when firewall is active

## Docker Integration
- Docker-in-Docker feature enabled for running containerized MCP servers
- Docker socket mounted from host for Docker command execution
- Supports both local and containerized development workflows

## VS Code Integration
- Format on save enabled with Prettier
- ESLint auto-fix on save
- Terraform and HCL syntax support
- Terminal defaults to zsh with enhanced prompt

# Setup Requirements

1. Docker Desktop
2. VS Code with Remote-Containers extension  
3. GitHub Personal Access Token (for GitHub CLI integration)
4. Open repository in VS Code and select "Reopen in Container"

# Post-Creation Setup

The `post-create.sh` script automatically:
- Authenticates GitHub CLI if `GH_TOKEN` is provided
- Installs Claude GitHub extension (`anthropic/gh-claude`)
- Provides manual setup instructions if token is not available