# Terraform Module Template

## Module Planning Template

### Basic Information

- **Module Name**: [module-name]
- **Provider**: [aws/azure/gcp/kubernetes]
- **Version**: [semantic-version]
- **Purpose**: [brief-description]

### Requirements Analysis

```
Business Context:
- [ ] Purpose and scope defined
- [ ] Compliance requirements identified
- [ ] Integration points documented

Technical Requirements:
- [ ] Resource dependencies mapped
- [ ] Security requirements specified
- [ ] Performance requirements defined
```

### Architecture Design

```
Module Structure:
- [ ] Resource hierarchy planned
- [ ] Module boundaries defined
- [ ] Integration patterns specified

Security Design:
- [ ] Access control patterns defined
- [ ] Network security planned
- [ ] Data protection specified
```

### Module Specification

```
Interface Design:
- [ ] Required variables defined
- [ ] Optional variables specified
- [ ] Output values planned
- [ ] Validation rules created

Resource Planning:
- [ ] Primary resources identified
- [ ] Supporting resources planned
- [ ] Data sources specified
- [ ] Lifecycle rules defined
```

### Implementation Planning

```
Development Strategy:
- [ ] Development order planned
- [ ] Testing approach defined
- [ ] Quality gates established
- [ ] Documentation requirements set

Risk Management:
- [ ] Technical risks identified
- [ ] Mitigation strategies planned
- [ ] Rollback procedures defined
- [ ] Monitoring approach specified
```

## Module Structure Template

```
modules/
├── [module-name]/
│   ├── main.tf            # Primary resources
│   ├── variables.tf       # Input variables
│   ├── outputs.tf         # Output values
│   ├── versions.tf        # Provider versions
│   ├── locals.tf          # Local values
│   ├── data.tf            # Data sources
│   ├── README.md          # Module documentation
│   ├── examples/          # Usage examples
│   │   ├── basic/
│   │   └── advanced/
│   └── tests/            # Module tests
│       ├── unit/
│       └── integration/
```

## Documentation Template

### README Structure

1. Module description and purpose
2. Usage examples (basic and advanced)
3. Input variable reference
4. Output value reference
5. Dependencies and requirements
6. Known issues and limitations
