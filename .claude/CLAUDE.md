# CLAUDE.md - Terraform Development Focus

- MCP Resources and tools to consult

    - Resources
        - *terraform_aws_best_practices* for AWS best practices about security, code base structure and organization, AWS Provider version management, and usage of community modules

        - terraform_awscc_provider_resources_listing for available AWS Cloud Control API resources
        - terraform_aws_provider_resources_listing for available AWS resources

    - Tools

        - SearchSpecificAwsIaModules tool to check for specialized AWS-IA modules first (Bedrock, OpenSearch Serverless, SageMaker, Streamlit)
        - SearchUserProvidedModule tool to analyze any Terraform Registry module provided by the user
        - SearchAwsccProviderDocs tool to look up specific Cloud Control API resources
        - SearchAwsProviderDocs tool to look up specific resource documentation.