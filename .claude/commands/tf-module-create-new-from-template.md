# Terraform Module Repository Creation from Template

This command guide enables Claude Code to create new Terraform module repositories from templates using GitHub CLI, following HashiCorp naming conventions and best practices.

## Required Terraform Module Naming Convention

**CRITICAL**: All Terraform registry modules must follow this exact naming pattern:
`terraform-<PROVIDER>-<NAME>`

- `<PROVIDER>`: Main cloud provider (aws, azure, google, etc.)
- `<NAME>`: Infrastructure type the module manages (can contain hyphens)
- Examples: `terraform-aws-vpc`, `terraform-google-vault`, `terraform-aws-ec2-instance`

Claude should ALWAYS validate module names against this pattern before creation.

## Prerequisites

1. Install GitHub CLI if you haven't already:

   ```bash
   # macOS
   brew install gh

   # Windows
   winget install --id GitHub.cli

   # Linux
   # Follow instructions at https://github.com/cli/cli/blob/trunk/docs/install_linux.md
   ```

2. Authenticate with GitHub:

   ```bash
   gh auth login
   ```

3. **Required Template Configuration:**

```bash
TEMPLATE_OWNER="hashi-demo-lab"
TEMPLATE_REPO="tf-module-template"
ORGANIZATION="hashi-demo-lab"
TEMPLATE_URL="https://github.com/hashi-demo-lab/tf-module-template"
```

## Creating Repository from Existing Template

Use the `gh repo create` command with the `--template` flag to create from an existing template repository:
when cloning don't specify any directory path this will be inherited by default

```bash
# Basic command structure - creates from existing template
# Note: Create repository first, then clone separately to avoid branch reference issues
gh repo create "NEW-REPO-NAME" \
  --template "EXISTING-TEMPLATE-OWNER/TEMPLATE-REPO-NAME" \
  --description "New repository description" \
  --public

# Then clone the repository separately
gh repo clone "ORGANIZATION-NAME/NEW-REPO-NAME"
```

### Examples with Existing Templates

```bash
# Create from HashiCorp's module template
gh repo create "hashi-demo-lab/terraform-aws-vpc" \
  --template "hashi-demo-lab/tf-module-template" \
  --description "Terraform AWS VPC module" \
  --public \
  --clone

# Create from your organization's existing template
gh repo create "hashi-demo-lab/terraform-azure-storage" \
  --template "hashi-demo-lab/tf-module-template" \
  --description "Terraform Azure Storage module" \
  --public \
  --clone
```

## Repository Creation Options

### Visibility Options

```bash
# Public repository (default)
--public

# Private repository
--private

# Internal repository (for organizations)
--internal
```

### Additional Options

```bash
# Clone after creation
--clone

# Add repository description
--description "Your description here"

# Disable issues
--disable-issues

# Disable wiki
--disable-wiki

# Include all branches from template
--include-all-branches
```

## Complete Command Examples

```bash
# Minimal command
gh repo create "my-new-repo" --template "owner/template-repo" --public --clone


# Full command with all options
gh repo create "terraform-aws-s3" \
  --template "hashicorp/terraform-module-template" \
  --description "Terraform module for AWS S3 buckets" \
  --public \
  --clone

# Private repository
gh repo create "private-module" \
  --template "myorg/private-template" \
  --description "Private Terraform module" \
  --private \
  --clone
```

## Claude Code Execution Steps

Execute these steps in order with proper error handling:

### Step 1: Gather Module Information

- Ask user for module name (validate against terraform-<provider>-<name> pattern)
- Confirm provider and infrastructure type
- Generate appropriate description

### Step 2: Create Repository

```bash
# Create repository from template (without cloning)
gh repo create "hashi-demo-lab/terraform-<provider>-<name>" \
  --template "hashi-demo-lab/tf-module-template" \
  --description "Terraform <provider> <name> module" \
  --public

# Clone the repository to current directory
gh repo clone "hashi-demo-lab/terraform-<provider>-<name>"
```

### Step 3: Navigate to Module Directory

```bash
# Verify directory exists and navigate
ls -la
cd terraform-<provider>-<name>
pwd  # Confirm we're in the right directory
```

### Step 4: Initialize Development Tools

```bash
# Check if directory structure is correct
ls -la

# Initialize TFLint (always available in devcontainer)
tflint --init

# Enable pre-commit hooks if available (optional step)
if command -v pre-commit &> /dev/null; then
    pre-commit install
else
    echo "Pre-commit not available - skipping (this is optional)"
fi
```

### Step 5: Verify Setup

```bash
# Ensure we're in the module directory
pwd
ls -la

# Validate Terraform configuration
terraform init
terraform validate


### Step 6: Confirmation and Next Steps

- Confirm all steps completed successfully
- Provide repository URL and local path
- List any optional tools that weren't initialized
- Suggest next development steps

## Related Documentation

- [GitHub Templates](https://docs.github.com/en/repositories/creating-and-managing-repositories/creating-a-template-repository)
- [GitHub CLI documentation](https://cli.github.com/manual/)
- [Repository Creation](https://docs.github.com/en/repositories/creating-and-managing-repositories)
```
