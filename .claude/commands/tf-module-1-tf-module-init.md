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
**Important**: The repository will be cloned into a directory matching the repo name

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

**CRITICAL**: Execute these steps in exact order. Stop execution if any step fails and report the error to the user.

### Step 1: Validate Prerequisites and Gather Information

**First, validate GitHub CLI is available:**

```bash
# Check if GitHub CLI is installed and authenticated
if ! command -v gh &> /dev/null; then
    echo "ERROR: GitHub CLI not found. Please install it first."
    exit 1
fi

# Verify authentication
if ! gh auth status &> /dev/null; then
    echo "ERROR: GitHub CLI not authenticated. Run 'gh auth login' first."
    exit 1
fi
```

**Then gather module information:**

- Ask user for module name (validate against terraform-<provider>-<name> pattern)
- Confirm provider and infrastructure type
- Generate appropriate description
- Verify template repository exists before proceeding

### Step 2: Create Repository

```bash
# Set variables for reliability
REPO_NAME="terraform-<provider>-<name>"
ORG_NAME="hashi-demo-lab"
TEMPLATE_REPO="hashi-demo-lab/tf-module-template"
DESCRIPTION="Terraform <provider> <name> module"

# Create repository from template (without cloning initially)
echo "Creating repository ${ORG_NAME}/${REPO_NAME}..."
gh repo create "${ORG_NAME}/${REPO_NAME}" \
  --template "${TEMPLATE_REPO}" \
  --description "${DESCRIPTION}" \
  --public

# Verify repository was created successfully
if ! gh repo view "${ORG_NAME}/${REPO_NAME}" &> /dev/null; then
    echo "ERROR: Repository creation failed or repository not accessible"
    exit 1
fi

# Clone the repository to current directory
echo "Cloning repository..."
gh repo clone "${ORG_NAME}/${REPO_NAME}"
```

### Step 3: Navigate to Module Directory

```bash
# Navigate to the cloned repository directory
cd "${REPO_NAME}"

# Verify we're in the correct directory
if [[ ! -f "main.tf" && ! -f "variables.tf" ]]; then
    echo "ERROR: Not in a Terraform module directory. Expected main.tf or variables.tf files."
    exit 1
fi

# Confirm location
echo "Successfully navigated to: $(pwd)"
ls -la
```

### Step 4: Initialize Development Tools

```bash
# Initialize TFLint (always available in devcontainer)
echo "Initializing TFLint..."
if ! tflint --init; then
    echo "WARNING: TFLint initialization failed, but continuing..."
fi

# Enable pre-commit hooks if available (optional step)
if command -v pre-commit &> /dev/null; then
    echo "Installing pre-commit hooks..."
    pre-commit install
else
    echo "Pre-commit not available - skipping (this is optional)"
fi

# Verify directory structure
echo "Module directory structure:"
ls -la
```

### Step 5: Verify Setup

```bash
# Validate Terraform configuration
echo "Initializing and validating Terraform..."
if ! terraform init; then
    echo "ERROR: Terraform initialization failed"
    exit 1
fi

if ! terraform validate; then
    echo "ERROR: Terraform validation failed"
    exit 1
fi

echo "✅ Terraform module setup completed successfully!"
```

### step 6: update CLAUDE.md to reflect what been completed

Include the following details when updating CLAUDE.md

- terraform module name
- Repository URL "https://github.com/${ORG_NAME}/${REPO_NAME}"
- Local directory path $(pwd)
- Any warnings or optional tools that weren't available
- Next development steps


### Step 7: Final Confirmation

```bash
# Display success summary
echo "==========================================="
echo "✅ MODULE CREATION COMPLETED SUCCESSFULLY"
echo "==========================================="
echo "Repository: https://github.com/${ORG_NAME}/${REPO_NAME}"
echo ""
echo "Local path: $(pwd)"
echo "Next steps:"
echo run /tf-module-planning
"
```

**Claude should report:**

- Repository URL
- Local directory path
- Any warnings or optional tools that weren't available
- Next development steps
- For module design and planning via a GitHub issue run /tf-module-planning. This will start the planning process via GitHub Issues

## Related Documentation

- [GitHub Templates](https://docs.github.com/en/repositories/creating-and-managing-repositories/creating-a-template-repository)
- [GitHub CLI documentation](https://cli.github.com/manual/)
- [Repository Creation](https://docs.github.com/en/repositories/creating-and-managing-repositories)
