# Terraform Module Requirements

Use this template to capture comprehensive technical and security requirements for new Terraform modules.

**📖 Reference Documentation:**

- [Module Development Workflow](/.claude/commands/tf-module-create-new-from-template.md)
- [Terraform Planning Framework](/.claude/CLAUDE.md)

## 📦 Module Name

**Module name following terraform-PROVIDER-NAME convention**

e.g., terraform-aws-vpc, terraform-azure-storage

## ☁️ Cloud Provider

**Primary cloud provider for this module**

- [ ] AWS
- [ ] Azure
- [ ] GCP
- [ ] Multi-cloud
- [ ] Other

## 📝 Module Description

**Brief description of what this module will accomplish**

Describe the infrastructure this module will manage...

## 🎯 Business Requirements

**Business use case and requirements driving this module**

- Business problem this module solves
- Expected outcomes
- Success criteria

## 📋 Specific Resources List

**List the exact resources this module will create**

Example:

- aws_vpc
- aws_subnet (public/private)
- aws_internet_gateway
- aws_nat_gateway
- aws_route_table
- aws_security_group

## 🔒 Security Requirements

**Security controls and requirements**

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

## 📋 Compliance Standards
*Compliance frameworks this module must adhere to*

Example:
- SOC 2 Type II
- PCI DSS
- HIPAA
- GDPR
- FedRAMP
- CIS Benchmarks
- Company security policies

## 🛡️ Specific Security Controls
*Detailed security controls and configurations required*

Example:
- All S3 buckets must block public access
- RDS instances must use encryption at rest
- Security groups must follow least privilege
- All resources must be tagged for compliance
- VPC Flow Logs must be enabled

## ⚙️ Terraform Technical Constraints
*Terraform version and provider constraints*

Example:
- Latest is major release, if no dependencies exist, example ~> 6.0
- Terraform >= 1.5.0
- AWS Provider >= 5.0.0
- Must use remote state backend
- Module should be registry-compatible

## 📥 Required Input Variables
*Essential input variables the module must accept*

Example:
- vpc_cidr: CIDR block for VPC (string, required)
- availability_zones: List of AZs (list(string), required)
- environment: Environment name (string, required)
- tags: Resource tags (map(string), optional)

## 📤 Required Output Variables
*Essential outputs the module must provide*

Example:
- vpc_id: ID of the created VPC
- private_subnet_ids: List of private subnet IDs
- public_subnet_ids: List of public subnet IDs
- security_group_id: Default security group ID

## 🔗 Module Dependencies
*Other modules or resources this module depends on*

Example:
- Existing Route53 hosted zone
- Shared KMS key for encryption
- Corporate DNS resolver
- Existing IAM roles

## 📅 Target Completion Date
*When this module needs to be available*

YYYY-MM-DD

## 📄 Additional Context
*Any additional information, constraints, or requirements*

Example:
- Integration with existing systems
- Migration considerations
- Special architectural requirements
- Team or organizational constraints

## ✅ Acceptance Criteria
*Criteria that must be met for this module to be considered complete*

- [ ] Module follows terraform-<PROVIDER>-<NAME> naming convention
- [ ] All security requirements implemented and verified
- [ ] Module passes all quality gates
- [ ] Documentation is complete and accurate
- [ ] Examples demonstrate common use cases
- [ ] Module is published to internal/public registry
- [ ] Team training completed if required
