# Phase 3: Module Specification

## Prompt Structure for Detailed Module Design

### Module Definition Framework

**Core Module Properties:**
- What is the single responsibility of this module?
- What resources will this module create and manage?
- What are the required vs optional input variables?
- What outputs should this module expose?

**Interface Design:**
- What are the minimum required inputs for basic functionality?
- What optional inputs provide advanced configuration?
- How should complex objects be structured in variables?
- What validation rules should be applied to inputs?

**Resource Organization:**
- How should resources be logically grouped within the module?
- What locals should be defined for computed values?
- How will conditional resource creation be handled?
- What data sources are needed for external references?

**Integration Patterns:**
- How will this module integrate with other modules?
- What dependencies exist on external resources?
- How will module outputs be consumed by other components?
- What lifecycle management considerations exist?

### Specification Template

For each module, define:

1. **Module Metadata**
   ```hcl
   # Module: [name]
   # Version: [semantic_version] 
   # Description: [purpose]
   # Dependencies: [list_of_dependencies]
   ```

2. **Input Variables**
   - Required variables with descriptions and types
   - Optional variables with defaults and constraints
   - Variable validation rules
   - Sensitive variable handling

3. **Resource Definitions**
   - Primary resources and their configurations
   - Supporting resources (IAM, security groups, etc.)
   - Data source requirements
   - Random/external resource needs

4. **Output Values**
   - Essential outputs for module consumers
   - Debug outputs for troubleshooting
   - Sensitive output handling
   - Output descriptions and usage examples

5. **Module Documentation**
   - Usage examples and common patterns
   - Integration guidelines
   - Troubleshooting guide
   - Version compatibility matrix