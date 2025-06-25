# Terraform Module Template

This template provides a standardized structure for Terraform modules with best practices for variables, outputs, documentation, and testing.

## Module Structure

```
module-name/
├── main.tf              # Primary resource definitions
├── variables.tf         # Input variable declarations
├── outputs.tf          # Output value declarations
├── versions.tf         # Provider version constraints
├── README.md           # Module documentation (auto-generated)
├── examples/           # Usage examples
│   ├── basic/
│   │   ├── main.tf
│   │   ├── variables.tf
│   │   └── outputs.tf
│   └── advanced/
│       ├── main.tf
│       ├── variables.tf
│       └── outputs.tf
├── tests/              # Automated tests
│   ├── basic_test.tftest.hcl
│   ├── integration_test.tftest.hcl
│   └── security_test.tftest.hcl
└── docs/               # Additional documentation
```

## Standard Variable Definitions

### variables.tf Template

```hcl
# Module metadata
variable "name" {
  description = "Name of the resource(s) to create"
  type        = string
  validation {
    condition     = can(regex("^[a-zA-Z0-9-_]+$", var.name))
    error_message = "Name must contain only alphanumeric characters, hyphens, and underscores."
  }
}

variable "environment" {
  description = "Environment name (e.g., dev, staging, prod)"
  type        = string
  default     = "dev"
  validation {
    condition     = contains(["dev", "staging", "prod"], var.environment)
    error_message = "Environment must be one of: dev, staging, prod."
  }
}

# Common tags
variable "tags" {
  description = "A map of tags to assign to the resource"
  type        = map(string)
  default     = {}
}

variable "additional_tags" {
  description = "Additional tags to merge with default tags"
  type        = map(string)
  default     = {}
}

# Feature toggles
variable "enable_monitoring" {
  description = "Enable monitoring and logging for the resource"
  type        = bool
  default     = true
}

variable "enable_encryption" {
  description = "Enable encryption at rest and in transit"
  type        = bool
  default     = true
}

# Resource-specific variables (example for S3 bucket)
variable "bucket_name" {
  description = "Name of the S3 bucket to create"
  type        = string
  default     = null
}

variable "versioning_enabled" {
  description = "Enable versioning for the S3 bucket"
  type        = bool
  default     = true
}

variable "lifecycle_rules" {
  description = "List of lifecycle rules for the bucket"
  type = list(object({
    id                            = string
    enabled                       = bool
    abort_incomplete_multipart_upload_days = number
    expiration_days              = number
    noncurrent_version_expiration_days = number
  }))
  default = []
}

# Networking variables (if applicable)
variable "vpc_id" {
  description = "VPC ID where resources will be created"
  type        = string
  default     = null
}

variable "subnet_ids" {
  description = "List of subnet IDs for resource placement"
  type        = list(string)
  default     = []
}

variable "security_group_ids" {
  description = "List of security group IDs to associate with the resource"
  type        = list(string)
  default     = []
}

# Data source filters
variable "data_source_filters" {
  description = "Filters for data source lookups"
  type        = map(string)
  default     = {}
}
```

## Standard Output Specifications

### outputs.tf Template

```hcl
# Resource identifiers
output "id" {
  description = "The ID of the created resource"
  value       = aws_s3_bucket.this.id
}

output "arn" {
  description = "The ARN of the created resource"
  value       = aws_s3_bucket.this.arn
}

output "name" {
  description = "The name of the created resource"
  value       = aws_s3_bucket.this.bucket
}

# Resource attributes
output "region" {
  description = "The region where the resource is created"
  value       = aws_s3_bucket.this.region
}

output "domain_name" {
  description = "The domain name of the resource (if applicable)"
  value       = aws_s3_bucket.this.bucket_domain_name
}

# Security and networking outputs
output "security_group_id" {
  description = "The ID of the security group created for the resource"
  value       = try(aws_security_group.this[0].id, null)
}

output "subnet_ids" {
  description = "List of subnet IDs where the resource is deployed"
  value       = var.subnet_ids
}

# Monitoring and logging outputs
output "cloudwatch_log_group_name" {
  description = "Name of the CloudWatch log group (if monitoring is enabled)"
  value       = try(aws_cloudwatch_log_group.this[0].name, null)
}

output "monitoring_enabled" {
  description = "Whether monitoring is enabled for the resource"
  value       = var.enable_monitoring
}

# Configuration outputs
output "configuration" {
  description = "Configuration summary of the created resource"
  value = {
    name                = aws_s3_bucket.this.bucket
    environment         = var.environment
    versioning_enabled  = var.versioning_enabled
    encryption_enabled  = var.enable_encryption
    monitoring_enabled  = var.enable_monitoring
    tags               = local.tags
  }
  sensitive = false
}

# Sensitive outputs (if any)
output "sensitive_data" {
  description = "Sensitive configuration data (marked as sensitive)"
  value = {
    access_key = try(aws_iam_access_key.this[0].id, null)
  }
  sensitive = true
}
```

## Documentation Requirements

### README.md Template (Auto-generated by terraform-docs)

````markdown
<!-- BEGIN_TF_DOCS -->

# Module Name

Brief description of what this module does and its primary use case.

## Usage

Basic usage example:

```hcl
module "example" {
  source = "./path/to/module"

  name        = "my-resource"
  environment = "dev"

  tags = {
    Owner   = "platform-team"
    Project = "infrastructure"
  }
}
```
````

Advanced usage example:

```hcl
module "advanced_example" {
  source = "./path/to/module"

  name        = "my-advanced-resource"
  environment = "prod"

  enable_monitoring = true
  enable_encryption = true

  lifecycle_rules = [
    {
      id                            = "cleanup"
      enabled                       = true
      abort_incomplete_multipart_upload_days = 7
      expiration_days              = 90
      noncurrent_version_expiration_days = 30
    }
  ]

  tags = {
    Owner       = "platform-team"
    Project     = "infrastructure"
    Environment = "production"
  }
}
```

## Requirements

| Name                                                                     | Version |
| ------------------------------------------------------------------------ | ------- |
| <a name="requirement_terraform"></a> [terraform](#requirement_terraform) | >= 1.8  |
| <a name="requirement_aws"></a> [aws](#requirement_aws)                   | >= 6.0  |

## Providers

| Name                                             | Version |
| ------------------------------------------------ | ------- |
| <a name="provider_aws"></a> [aws](#provider_aws) | >= 6.0  |

## Resources

| Name                                                                                                        | Type     |
| ----------------------------------------------------------------------------------------------------------- | -------- |
| [aws_s3_bucket.this](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/s3_bucket) | resource |

## Inputs

| Name                                          | Description                       | Type     | Default | Required |
| --------------------------------------------- | --------------------------------- | -------- | ------- | :------: |
| <a name="input_name"></a> [name](#input_name) | Name of the resource(s) to create | `string` | n/a     |   yes    |

## Outputs

| Name                                      | Description                    |
| ----------------------------------------- | ------------------------------ |
| <a name="output_id"></a> [id](#output_id) | The ID of the created resource |

<!-- END_TF_DOCS -->

````

## Testing Patterns

Terraform offers a built-in testing framework with `terraform test` that automatically discovers and executes tests in the `tests` directory. Tests are written in HCL and provide a more native testing experience. Preference using plan basic unit tests whilst iterating on development. Only test with apply based Integration tests to confirm the final solution before PR

### Example of Basic Unit Test (tests/basic_test.tftest.hcl)

```hcl
# Basic functionality test
variables {
  name        = "test-resource"
  environment = "test"
}

# Run the module with basic configuration
run "basic_configuration" {
  command = plan

  # Override variables
  variables {
    name        = var.name
    environment = var.environment
  }

  # Assertions for planned outputs
  assert {
    condition     = module.name == var.name
    error_message = "Resource name should match the input name variable"
  }

  assert {
    condition     = length(module.id) > 0
    error_message = "Resource ID should not be empty"
  }
}

````

### Example Integration Test (tests/integration_test.tftest.hcl)

```hcl
# Integration test with actual cloud resources
variables {
  name               = "integration-test"
  environment        = "test"
  enable_monitoring  = true
  region             = "us-west-2"
}

# Run the advanced example
run "advanced_configuration" {
  command = apply

  variables {
    name              = var.name
    environment       = var.environment
    enable_monitoring = var.enable_monitoring
  }

  # Verify AWS-specific configurations
  assert {
    condition     = output.configuration.versioning_enabled == "true"
    error_message = "Versioning should be enabled for the S3 bucket"
  }

  # Check that monitoring is enabled
  assert {
    condition     = output.monitoring_enabled == var.enable_monitoring
    error_message = "Monitoring should be enabled"
  }

  # Destroy after test (built-in behavior with terraform test)
}
```

## Best Practices Checklist

### Module Design

- [ ] Single responsibility principle - module does one thing well
- [ ] Composable - can be used with other modules
- [ ] Configurable - sensible defaults with override options

### Code Quality

- [ ] Variables have descriptions and validation rules
- [ ] Outputs are comprehensive and well-documented
- [ ] Resource naming follows conventions
- [ ] Tags are consistent and comprehensive

### Security

- [ ] Encryption enabled by default where applicable
- [ ] Least privilege access principles
- [ ] Sensitive data marked appropriately
- [ ] Security scanning passes (Checkov, TFSec, Terrascan)

### Documentation

- [ ] README.md auto-generated with terraform-docs
- [ ] Usage examples provided
- [ ] Variable descriptions are clear
- [ ] Output descriptions are helpful

### Testing

- [ ] Use Terraform validate to confirm valid code
- [ ] Tests written using Terraform's native test framework (.tftest.hcl files)
- [ ] Integration tests verify cloud provider interaction
- [ ] Security compliance tests validate security settings
- [ ] All tests pass with `terraform test` command
