# Terraform module design planning

## add to memory

- planning_workflow @/workspace/.claude/templates/workflow-planning.mmd
- github_issues_module_template "📋 Terraform Module Requirements"
- tfsec rules @/workspace/.claude/planning/security_rules/tfsec_rules.md
- terraform checkov security rules @/workspace/.claude/planning/security_rules/terraform_graph_checks
- Terraform Style Guide and best practices @/workspace/.claude/planning/best_practices/terraform_best_practices.md

## Role Assignments & Collaboration Model

This workflow follows a structured human-AI collaboration model:

- 👤 **Human Only**: Tasks requiring human judgment, business context, or final approval
- 🤖 **AI Only**: Research, analysis, and documentation tasks that benefit from AI capabilities  
- 👥 **Human-AI Pair**: Collaborative tasks combining human insight with AI assistance
- ❓ **Decision Points**: Critical approval gates requiring human decision-making

## Before starting

- ensure ide is connected by running /IDE in claude.

## planning steps

This task is focused on planning and design for a Terraform modules.
The intenton of this prompt is for planning and design only we are not writing any Terraform code.

Steps should be performed in the following order:

1. **Understand Planning Workflow** 👥 (Human-AI Pair)
   - Review the planning_workflow from the mermaid diagram workflow-planning.mmd @/workspace/.claude/templates/workflow-planning.mmd
   - Confirm understanding of the complete workflow phases and role assignments

2. **Create GitHub Issue** 👥 (Human-AI Pair)
   - Prompt the user to create a Github issue from existing template "📋 Terraform Module Requirements"
   - Populate known inputs (providers, basic requirements, etc.)
   - Ensure issue captures initial scope and objectives

3. **AI-Assisted Planning Phase** 👥 (Human-AI Pair)
   - **Requirements Analysis** 👥 (Technical & Security Requirements)
     - Analyze functional and non-functional requirements
     - Review security requirements using tfsec rules and terraform checkov security rules
     - Reference Terraform Style Guide and best practices
     - Document compliance and governance requirements
   
   - **Resource Research** 🤖 (AI-Only Task)
     - Use MCP servers to research AWS provider documentation
     - Search for existing AWS-IA modules that could be leveraged
     - Research best practices for the specific resource types
     - Identify security scanning requirements and tools
     - Document findings and recommendations

4. **Architecture Design Creation** 👥 (Human-AI Pair)
   - Create comprehensive architecture design with AI assistance
   - Generate architecture diagrams and documentation
   - Define module structure, inputs, outputs, and dependencies
   - Document security controls and compliance measures
   - Create cost estimation framework

5. **Design PR Creation** 👥 (Human-AI Pair)
   - Create pull request with architecture documentation
   - Include all design artifacts and diagrams
   - Add architectural decision records (ADRs) if applicable
   - Ensure documentation follows project standards

6. **Review & Approval Cycle** 
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

7. **Planning Phase Completion** 🎯
   - Confirm all planning artifacts are complete and approved
   - Transition planning issue to "Ready for Development" status
   - Create development phase GitHub issue or milestone
   - Document lessons learned and process improvements
