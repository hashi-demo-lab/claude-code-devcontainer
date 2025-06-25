# Terraform Module Development Planning Framework

A structured approach to planning and developing Terraform modules with comprehensive prompt templates for each development phase.

## Directory Structure

```
terraform-planning/
├── phases/                          # Planning phase prompts
│   ├── 01-requirements-analysis.md  # Business and technical requirements
│   ├── 02-architecture-design.md    # System and module architecture  
│   ├── 03-module-specification.md   # Detailed module interfaces
│   └── 04-implementation-planning.md # Development strategy and risks
├── workflows/                       # Process documentation
│   └── module-development-workflow.md # Complete development process
├── templates/                       # Reusable templates
│   └── module-template.md          # Module planning checklist
└── README.md                       # This file
```

## Usage

1. **Start with Requirements Analysis** (`phases/01-requirements-analysis.md`)
   - Gather business context and technical constraints
   - Identify compliance and security requirements
   - Document integration needs and dependencies

2. **Design the Architecture** (`phases/02-architecture-design.md`)
   - Create high-level system architecture
   - Define module boundaries and responsibilities
   - Plan security and deployment strategies

3. **Specify Module Details** (`phases/03-module-specification.md`)
   - Define precise module interfaces
   - Document input variables and outputs
   - Plan resource organization and lifecycle

4. **Plan Implementation** (`phases/04-implementation-planning.md`)
   - Establish development order and testing strategy
   - Set up quality assurance processes
   - Define risk management approaches

## Key Features

- **Phase-based Planning**: Structured approach with clear deliverables
- **Comprehensive Prompts**: Detailed questions for each planning phase
- **Quality Focus**: Built-in quality gates and review processes
- **Risk Management**: Proactive identification and mitigation strategies
- **Template-driven**: Consistent approach across all modules

## Integration with Development Environment

This framework integrates with the DevContainer's pre-installed tools:
- **Terraform Tools**: TFLint, TFSec, Terrascan for quality and security
- **Documentation**: terraform-docs for automated documentation
- **Testing**: Checkov for security scanning
- **Cost Analysis**: Infracost for cost estimation