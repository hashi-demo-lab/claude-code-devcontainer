# Terraform module design planning

## add to claudes memory

- planning_workflow @/workspace/.claude/templates/workflow-planning.mmd
- github_issues_module_template "📋 Terraform Module Requirements"
- tfsec rules @/workspace/.claude/planning/security_rules/tfsec_rules.md
- terraform checkov security rules @/workspace/.claude/planning/security_rules/terraform_graph_checks
- AWS foundational security best practices @//workspace/.claude/planning/security_rules/aws-foundational-security-all-controls.json
- Terraform style guide and best practices @/workspace/.claude/planning/best_practices/terraform_best_practices.md
- Module structure "The module structure, file and directory layout is already preexisting in the module repository it should have been cloned into a subfolder, but check to confirm, and set the working directory to the module"
- design outputs "All planning an design outputs should be writtern to a subfolder of the module called "design and planning".
- Module status "The module should already to cloned into the sub directory of the workspace and should follow the formation terraform-<provider name>-<provider resource>"
- Terraform mcp "For Terraform when choosing MCP servers using the terraform hashicorp/terraform-mcp-server "prioritize tool lookups using the for getting provider documentation vs using awslabs."
- GitHub rules "when working in GitHub always work under a feature branch, never commit directly to main"
- GitHub issues "For planning GitHub issues should always be labelled documentation"
- Module patterns "When looking at module for AWS they are published by aws-ia use 'SearchSpecificAwsIaModules', for Azure modules are published by Azure, IBM modules are published by terraform-ibm-modules"
- Mermaid diagrams "For all mermaid diagrams, apply neo theme and a hierarchial layout"
- Plannig outputs "create all module design planning outputs in the design and planning subfolder of the module, terraform-<provider name>-<provider resource>/design and planning, architecture-design.md, architecture-design.mmd, requirements-analysis.md"

## Role Assignments & Collaboration Model

This workflow follows a structured human-AI collaboration model:

- 👤 **Human Only**: Tasks requiring human judgment, business context, or final approval
- 🤖 **AI Only**: Research, analysis, and documentation tasks that benefit from AI capabilities
- 👥 **Human-AI Pair**: Collaborative tasks combining human insight with AI assistance
- ❓ **Decision Points**: Critical approval gates requiring human decision-making

## Before starting

- ensure ide is connected by running /ide in claude.
- ensure that the current working directory is /workspace, this wokrfllow is intended to be run from a standardized devcontainer and this should start at /workspace.

## planning steps

This task is focused on planning and design for a Terraform modules.
The intenton of this prompt is for planning and design only we are not writing any Terraform code.

Steps should be performed in the following order:

1. **Understand Planning Workflow** 👥 (Human-AI Pair)

   - Review the planning_workflow from the mermaid diagram workflow-planning.mmd @/workspace/.claude/templates/workflow-planning.mmd
   - Confirm understanding of the complete workflow phases and role assignments

2. **Assess the module repository**

   - The module should already to cloned into the sub directory of the workspace and should follow the formation terraform-<provider name>-<provider resource>"
   - ultrathink about potential designs patterns and considerations.
   - get the latest provider versions for the targeted provider via MCP

3. **Read the GitHub issue template**
example:

```bash
   gh issue view-template .github/ISSUE_TEMPLATE/terraform-module-requirements.yml
```

4. **Create GitHub Issue** 👥 (Human-AI Pair)

   **Sub-steps:**

   a. **Template Population** 🤖 (AI-Only)

   - AI creates populated GitHub issue template "📋 Terraform Module Requirements"
   - Include provider requirements, basic functionality, security needs
   - Add initial scope and objectives based on user input

   b. **Issue Creation** 🤖 (AI-Only)

   - AI uses GitHub CLI (`gh issue create`) to create the issue in the repository. example below:

   ```bash
         gh issue create \
      --template ".github/ISSUE_TEMPLATE/terraform-module-requirements.yml" \
      --title "<module description>" \
      --assignee username \
      --label "documentation"
   ```

   - Use populated template as the issue body
   - Apply appropriate labels (e.g., "documentation")
   - Assign to appropriate milestone if exists

   c. **Verification** 👥 (Human-AI Pair)

   - Verify issue creation was successful
   - Confirm issue URL is accessible
   - Update issue with any additional context from user
   - Proceed only after successful issue creation

5. **AI-Assisted Planning Phase** 👥 (Human-AI Pair)

   - **Requirements Analysis** 👥 (Technical & Security Requirements)

     - Analyze functional and non-functional requirements
     - Review security requirements using tfsec rules and terraform checkov security rules
     - Reference Terraform Style Guide and best practices
     - Document compliance and governance requirements
     - Understand existing modules patterns from the public module registry. Use MCP to get module patterns.

   - **Resource Research** 🤖 (AI-Only Task)
     - Use MCP servers to research AWS provider documentation
     - Search for existing AWS-IA modules that could be leveraged
     - Research best practices for the specific resource types
     - Identify security scanning requirements and tools
     - Document findings and recommendations

6. **Architecture Design Creation** 👥 (Human-AI Pair)

   - Create comprehensive architecture design with AI assistance
   - Generate architecture diagrams and documentation using Mermaid 
   - Define module structure, inputs, outputs, and dependencies
   - Document security controls and compliance measures
   - Create cost estimation framework

7. **Create GitHub PR for all planniong and design artifacts** 👥 (Human-AI Pair)
   - Create pull request with architecture documentation
   - Include all design artifacts and diagrams
   - Add architectural decision records (ADRs) if applicable
   - Ensure documentation follows project standards

8. **Planning Phase Completion** 🎯
   - Confirm all planning artifacts are complete and approved
   - Transition planning issue to "Ready for Development" status
   - Create development phase GitHub issue or milestone
   - Document lessons learned and process improvements

9. **Review & Approval Cycle**

   - **Architecture Review** 👤 (Human-Only Task)

     - Technical review of proposed architecture
     - Security and compliance validation
     - Cost and operational impact assessment

   - **Decision Point** ❓ (Approval Gateway)

     - ✅ **If Approved**: Proceed to development phase
     - ❌ **If Changes Requested**: Return to step 4 for design updates

   - **Design Updates** 👥 (Human-AI Pair - if needed)
     - Address review feedback
     - Update architecture documentation
     - Push updates and request re-review

