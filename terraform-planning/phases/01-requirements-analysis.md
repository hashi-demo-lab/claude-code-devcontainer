# Phase 1: Requirements Analysis

## Prompt Structure for Requirements Gathering

### Primary Analysis Questions

**Infrastructure Requirements:**
- What AWS/Azure/GCP services need to be provisioned?
- What are the networking requirements (VPC, subnets, security groups)?
- What are the compute requirements (instances, containers, serverless)?
- What storage requirements exist (S3, databases, file systems)?

**Business Context:**
- What is the business purpose of this infrastructure?
- What are the compliance requirements (SOC2, HIPAA, PCI-DSS)?
- What are the expected traffic patterns and scaling needs?
- What environments are needed (dev, staging, prod)?

**Technical Constraints:**
- What existing infrastructure must be integrated with?
- What naming conventions and tagging standards must be followed?
- What security policies and access controls are required?
- What monitoring and logging requirements exist?

**Resource Dependencies:**
- What external services or APIs will be consumed?
- What data sources need to be accessed?
- What secrets management is required?
- What backup and disaster recovery needs exist?

### Output Format

Create a structured requirements document covering:

1. **Functional Requirements**
   - Primary infrastructure components
   - Service configurations
   - Integration points

2. **Non-Functional Requirements**
   - Performance expectations
   - Security requirements
   - Compliance needs
   - Cost constraints

3. **Technical Requirements**
   - Provider-specific constraints
   - Version requirements
   - Dependency mappings

4. **Operational Requirements**
   - Deployment procedures
   - Monitoring needs
   - Maintenance windows