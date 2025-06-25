# Phase 2: Architecture Design

## Prompt Structure for Technical Architecture

### High-Level Design Questions

**System Architecture:**
- How should components be logically grouped into modules?
- What are the data flow patterns between components?
- How will different environments be differentiated?
- What are the scaling and availability patterns?

**Module Structure:**
- What should be the module boundaries and responsibilities?
- How will modules communicate and share data?
- What common modules can be reused across projects?
- How will module versioning and updates be managed?

**Security Architecture:**
- How will network segmentation be implemented?
- What access control patterns will be used?
- How will secrets and credentials be managed?
- What encryption requirements exist for data at rest and in transit?

**Infrastructure Patterns:**
- What deployment patterns will be used (blue-green, canary, rolling)?
- How will state management and locking be handled?
- What backup and disaster recovery patterns are needed?
- How will cost optimization be implemented?

### Design Deliverables

1. **Architecture Diagrams**
   - High-level system architecture
   - Network topology diagrams
   - Data flow diagrams
   - Security architecture diagrams

2. **Module Design**
   - Module dependency graph
   - Interface definitions
   - Resource hierarchy
   - State management strategy

3. **Technical Specifications**
   - Resource sizing and configuration
   - Security group rules and access patterns
   - Monitoring and alerting specifications
   - Backup and recovery procedures

4. **Implementation Strategy**
   - Deployment sequence and dependencies
   - Testing and validation approach
   - Rollback procedures
   - Performance benchmarks