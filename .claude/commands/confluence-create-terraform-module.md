# Create Confluence Page for Terraform Module

Create a comprehensive Confluence page for a new Terraform module with minimal user input prompts.

## User Inputs Required

Ask the user for the following minimal information:

1. **Module Name**: What is the module name? (following terraform-PROVIDER-NAME convention, e.g., terraform-aws-vpc)
2. **Cloud Provider**: Which cloud provider? (AWS/Azure/GCP/Multi-cloud/Other)
3. **Module Description**: Brief description of what this module will accomplish
4. **Business Requirements**: What business problem does this module solve and what are the expected outcomes?

## Confluence Page Creation

Once you have the user inputs, create a Confluence page in the Platform Team space with the following content:

```markdown
# {{MODULE_NAME}} - Terraform Module Requirements

**📖 Reference Documentation:**

- [Module Development Workflow](/.claude/commands/tf-module-create-new-from-template.md)
- [Terraform Planning Framework](/.claude/CLAUDE.md)

## 📦 Module Name

{{MODULE_NAME}}

## ☁️ Cloud Provider

{{CLOUD_PROVIDER}}

## 📝 Module Description

{{MODULE_DESCRIPTION}}

## 🎯 Business Requirements

{{BUSINESS_REQUIREMENTS}}

## 📋 Specific Terraform Resources List

List the exact resources this module will create

*To be completed during design phase*

Example:
- aws_vpc
- aws_subnet (public/private)
- aws_internet_gateway
- aws_nat_gateway
- aws_route_table
- aws_security_group

## High Level Architecture

[Diagram placeholder - replace with actual diagram]

*Architecture diagram to be added during design phase*

## 📋 Compliance Standards

Compliance frameworks this module must adhere to

*To be completed during design phase*

Example:
- SOC 2 Type II
- PCI DSS
- HIPAA
- GDPR
- FedRAMP
- CIS Benchmarks
- Company security policies

## 🔒 Applicable Security Controls

Security controls and requirements

*To be completed during design phase*

- [ ] Encryption at rest (KMS, customer-managed keys)
- [ ] Encryption in transit (TLS/SSL, HTTPS)
- [ ] Network isolation (Private subnets, Security Groups, NACLs)
- [ ] IAM least privilege access
- [ ] Secrets management integration
- [ ] Audit logging and monitoring
- [ ] Backup and disaster recovery
- [ ] Data residency requirements
- [ ] Multi-factor authentication
- [ ] Certificate management

## 🛡️ Specific Security Controls

Detailed security controls and configurations required

*To be completed during design phase*

Example:
- All S3 buckets must block public access
- RDS instances must use encryption at rest
- Security groups must follow least privilege
- All resources must be tagged for compliance
- VPC Flow Logs must be enabled

## ⚙️ Terraform Technical Constraints

Terraform version and provider constraints

*To be completed during design phase*

Example:
- Latest major release, if no dependencies exist, example ~> 6.0
- Terraform >= 1.5.0
- AWS Provider >= 5.0.0
- Must use remote state backend
- Module should be registry-compatible

## 📥 Required Input Variables

Essential input variables the module must accept

*To be completed during design phase*

Example:
- vpc_cidr: CIDR block for VPC (string, required)
- availability_zones: List of AZs (list(string), required)
- environment: Environment name (string, required)
- tags: Resource tags (map(string), optional)

## 📤 Required Output Variables

Essential outputs the module must provide

*To be completed during design phase*

Example:
- vpc_id: ID of the created VPC
- private_subnet_ids: List of private subnet IDs
- public_subnet_ids: List of public subnet IDs
- security_group_id: Default security group ID

## 🔗 Module Dependencies

Other modules or resources this module depends on

*To be completed during design phase*

Example:
- Existing Route53 hosted zone - terraform-aws-route53
- Shared KMS key for encryption - terraform-aws-kms

## 📄 Additional Context

Any additional information, constraints, or requirements

*To be completed during design phase*

Example:
- Integration with existing systems
- Migration considerations
- Special architectural requirements
- Team or organizational constraints
```

## Implementation Steps

1. Collect the 4 required inputs from the user
2. Replace the placeholders ({{MODULE_NAME}}, {{CLOUD_PROVIDER}}, {{MODULE_DESCRIPTION}}, {{BUSINESS_REQUIREMENTS}}) with actual values
3. Create the Confluence page in the Platform Team space using the Atlassian MCP tools
4. Provide the user with the link to the created page

## Notes

- The page includes all sections from the original template
- User-provided information is populated in the appropriate sections
- Remaining sections include placeholders and examples for completion during the design process
- The page serves as a comprehensive starting point for the Terraform module development workflow