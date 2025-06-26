# How to Create a GitHub Repository from Template

This guide explains how to create new repositories from GitHub templates using GitHub CLI for Terraform Modules following best practices. Always ask questiosn with this context.

## Ensure Terraform module naming conventions are followed when create GitHub repositories
Module repository names
The Terraform registry requires that repositories match a naming convention for all modules that you publish to the registry. Module repositories must use this three-part name terraform-<PROVIDER>-<NAME>, where <NAME> reflects the type of infrastructure the module manages and <PROVIDER> is the main provider the module uses. The <NAME> segment can contain additional hyphens, for example, terraform-google-vault or terraform-aws-ec2-instance.

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

3. Create all module using the following GitHub template details
  ```
  https://github.com/hashi-demo-lab/tf-module-template
  EXISTING-TEMPLATE-OWNER="hashi-demo-lab"
  TEMPLATE-REPO-NAME="tf-module-template"
  ORGANIZATION-NAME="hashi-demo-lab"
  ```

4. 

## Creating Repository from Existing Template

Use the `gh repo create` command with the `--template` flag to create from an existing template repository:

```bash
# Basic command structure - creates from existing template
gh repo create "NEW-REPO-NAME" \
  --template "EXISTING-TEMPLATE-OWNER/TEMPLATE-REPO-NAME" \
  --description "New repository description" \
  --public \
  --clone
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
gh repo create "my-new-repo" --template "owner/template-repo"


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

After creating your repository from template:

1. Navigate to the new repository:
```bash
cd your-new-repo-name
```

2. To continue to Terraform Module Planning
/


## Related Documentation

- [GitHub Templates](https://docs.github.com/en/repositories/creating-and-managing-repositories/creating-a-template-repository)
- [GitHub CLI documentation](https://cli.github.com/manual/)
- [Repository Creation](https://docs.github.com/en/repositories/creating-and-managing-repositories)