# CLAUDE.md - Terraform Development Focus

This file provides specialized guidance to Claude Code for Terraform infrastructure development in this repository.

# Terraform Ecosystem Overview

This repository is optimized for Terraform module development and infrastructure as code workflows. Use the comprehensive Terraform toolchain and MCP integrations for enhanced development capabilities.

## Terraform Tool Suite

### Core Terraform Tools

- **Terraform 1.12.1**: Primary infrastructure provisioning tool
- **Terraform Docs 0.20.0**: Auto-generate documentation from Terraform modules
- **TFLint 0.48.0**: Linter with AWS/Azure/GCP provider rulesets for code quality

### Security and Compliance Tools

- **TFSec 1.28.13**: Static analysis security scanner for Terraform code
- **Terrascan 1.19.9**: Multi-cloud security vulnerability scanner
- **Checkov 3.2.439**: Comprehensive policy-as-code security scanner
- **Infracost 0.10.41**: Cloud cost estimation for Terraform changes

### MCP Server Integration for Terraform

#### AWS Labs Terraform MCP Server

**Primary tool for Terraform operations** - Use these functions:

- `ExecuteTerraformCommand`: Run terraform init, plan, validate, apply, destroy
- `SearchAwsProviderDocs`: Get AWS provider resource documentation
- `SearchAwsccProviderDocs`: Access AWS Cloud Control API provider docs
- `SearchSpecificAwsIaModules`: Find AWS-IA certified modules
- `RunCheckovScan`: Execute security scans on Terraform code
- `SearchUserProvidedModule`: Analyze custom Terraform modules

#### Standard Terraform MCP Server

**Complementary tool** - Use for:

- Additional Terraform provider documentation lookups
- Module registry searches and analysis
- Cross-reference with AWS Labs server for comprehensive coverage

## Terraform Development Workflow

### 1. Module Planning Phase

Use the structured planning framework in `terraform-planning/`:

- **Requirements Analysis** (`phases/01-requirements-analysis.md`)
- **Architecture Design** (`phases/02-architecture-design.md`)
- **Module Specification** (`phases/03-module-specification.md`)
- **Implementation Planning** (`phases/04-implementation-planning.md`)

### 2. Development Commands

Always run these commands during Terraform development:

```bash
# Initialize and validate
terraform init
terraform validate

# Security scanning
checkov -d . --framework terraform
tfsec .
terrascan scan -i terraform -d .

# Linting
tflint --init
tflint

# Documentation generation
terraform-docs markdown table . > README.md

# Cost analysis (requires AWS credentials)
infracost breakdown --path .
```

### 3. Quality Gates

Before completing any Terraform work, ensure:

1. **Security**: All security scanners pass (Checkov, TFSec, Terrascan)
2. **Linting**: TFLint passes with no errors
3. **Validation**: `terraform validate` succeeds
4. **Documentation**: Module documentation is auto-generated and current
5. **Cost Analysis**: Cost impact is understood via Infracost

## MCP Integration Best Practices

### For AWS Resources

1. **Always search AWS provider docs first** using `SearchAwsProviderDocs`
2. **Use AWS-IA modules when available** via `SearchSpecificAwsIaModules`
3. **Consider AWSCC provider** for Cloud Control API resources via `SearchAwsccProviderDocs`
4. **Run security scans** using `RunCheckovScan` before finalizing code

### For Module Development

1. **Search existing modules** using `SearchUserProvidedModule` to avoid duplication
2. **Validate against best practices** using the planning framework
3. **Document thoroughly** using terraform-docs integration
4. **Test security posture** with all available scanners

## Environment Configuration

### Required Environment Variables

- `AWS_PROFILE` or AWS credentials for provider operations
- `TF_LOG`: Set to `DEBUG` for detailed Terraform logging
- `CHECKOV_LOG_LEVEL`: Set to `INFO` for security scan details

### Recommended Aliases

```bash
alias tf='terraform'
alias tfd='terraform-docs'
alias tfscan='checkov -d . --framework terraform && tfsec . && terrascan scan -i terraform -d .'
```

## Security Guidelines

### Terraform-Specific Security

- **Never commit state files** - Use remote state backends
- **Encrypt sensitive variables** - Use AWS Parameter Store/Secrets Manager
- **Scan before apply** - Always run security tools before deployment
- **Follow least privilege** - Use minimal required IAM permissions
- **Version pin providers** - Specify exact provider versions

### DevContainer Security

- Container runs with elevated privileges for Docker-in-Docker
- Network isolation available but currently disabled
- All Terraform operations are containerized for isolation

## Module Development Templates

Use the template in `templates/module-template.md` for consistent module structure:

- Standard variable definitions
- Output specifications
- Documentation requirements
- Testing patterns

Refer to `workflows/module-development-workflow.md` for step-by-step development process.
