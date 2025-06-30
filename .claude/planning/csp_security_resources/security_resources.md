# Programmatic Security Sources for Cloud Providers

## AWS Programmatic Sources

| Source Type | Command/API | Description | Output Format |
|------------|-------------|-------------|---------------|
| **CLI Commands** |
| Security Hub | `aws securityhub get-findings` | Security findings across services | JSON |
| Security Hub Standards | `aws securityhub describe-standards` | Available compliance standards (CIS, PCI-DSS, etc.) | JSON |
| Config Rules | `aws configservice describe-config-rules` | All Config rules and their details | JSON |
| Config Compliance | `aws configservice describe-compliance-by-config-rule` | Compliance status by rule | JSON |
| Access Analyzer | `aws accessanalyzer list-findings` | IAM access findings | JSON |
| GuardDuty | `aws guardduty list-findings` | Threat detection findings | JSON |
| Service Control Policies | `aws organizations list-policies --filter SERVICE_CONTROL_POLICY` | Organization-level policies | JSON |
| Well-Architected | `aws wellarchitected list-lenses` | Available assessment lenses | JSON |
| Well-Architected Review | `aws wellarchitected get-lens-review` | Security pillar assessment | JSON |
| Control Tower | `aws controltower list-enabled-controls` | Enabled guardrails | JSON |
| **API/SDK** |
| SecurityHub API | `GetEnabledStandards()` | Enabled compliance frameworks | JSON |
| Config API | `DescribeComplianceByResource()` | Resource-level compliance | JSON |
| IAM API | `GetAccountAuthorizationDetails()` | Complete IAM configuration | JSON |
| Organizations API | `DescribePolicy()` | Policy details and attachments | JSON |
| Trusted Advisor API | `DescribeTrustedAdvisorChecks()` | Security best practice checks | JSON |

## Azure Programmatic Sources

| Source Type | Command/API | Description | Output Format |
|------------|-------------|-------------|---------------|
| **CLI Commands** |
| Policy Definitions | `az policy definition list` | All available policy definitions | JSON |
| Policy Assignments | `az policy assignment list` | Active policy assignments | JSON |
| Policy Compliance | `az policy state list` | Current compliance state | JSON |
| Security Assessments | `az security assessment list` | Security posture assessments | JSON |
| Regulatory Compliance | `az security regulatory-compliance-controls list` | Compliance control status | JSON |
| Compliance Standards | `az security regulatory-compliance-standards list` | Available standards (ISO, PCI, etc.) | JSON |
| Blueprints | `az blueprint list` | Security blueprints | JSON |
| Secure Score | `az security secure-score-controls list` | Security score breakdown | JSON |
| Security Settings | `az security setting list` | Security configuration | JSON |
| Resource Graph | `az graph query -q 'SecurityResources'` | Query security resources | JSON |
| **API/SDK** |
| Policy Insights API | `PolicyStates.List()` | Policy compliance states | JSON |
| Security Center API | `Assessments.List()` | Security assessments | JSON |
| Resource Graph API | `Resources.Query()` | Complex security queries | JSON |
| Management Groups API | `PolicyDefinitions.List()` | Hierarchical policies | JSON |
| Compliance API | `RegulatoryComplianceStandards.List()` | Compliance framework data | JSON |

## GCP Programmatic Sources

| Source Type | Command/API | Description | Output Format |
|------------|-------------|-------------|---------------|
| **CLI Commands** |
| Security Command Center | `gcloud scc findings list` | Security findings and vulnerabilities | JSON |
| Security Policies | `gcloud compute security-policies list` | Network security policies | JSON |
| Organization Policies | `gcloud resource-manager org-policies list` | Organization-level constraints | JSON |
| Asset Inventory | `gcloud asset search-all-resources --query` | Resource security properties | JSON |
| Access Context | `gcloud access-context-manager policies list` | Access control policies | JSON |
| Binary Authorization | `gcloud container binauthz policy export` | Container security policies | YAML |
| IAM Policies | `gcloud projects get-iam-policy` | Project IAM configuration | JSON |
| Policy Details | `gcloud org-policies describe` | Detailed policy constraints | JSON |
| Recommendations | `gcloud recommender recommendations list` | Security recommendations | JSON |
| **API/SDK** |
| Security Command Center | `ListFindings()` | Programmatic access to findings | JSON |
| Organization Policy | `ListConstraints()` | Available policy constraints | JSON |
| Cloud Asset API | `AnalyzeIamPolicy()` | IAM analysis and insights | JSON |
| Policy Analyzer | `QueryActivity()` | Policy usage analysis | JSON |
| Recommender API | `ListRecommendations()` | Security insights and suggestions | JSON |

## Usage Examples

### AWS Example
```bash
# Get all security hub findings for critical severity
aws securityhub get-findings --filters '{"SeverityLabel":[{"Value":"CRITICAL","Comparison":"EQUALS"}]}'

# Check Config rule compliance
aws configservice describe-compliance-by-config-rule --config-rule-names required-tags

# List all SCPs in organization
aws organizations list-policies --filter SERVICE_CONTROL_POLICY --output json

# Get Well-Architected security pillar review
aws wellarchitected get-lens-review --workload-id <workload-id> --lens-alias wellarchitected
```

### Azure Example
```bash
# List all security assessments
az security assessment list --query "[?properties.status.code=='Unhealthy']"

# Get regulatory compliance status
az security regulatory-compliance-controls list --standard-name 'PCI-DSS-3.2.1'

# Query policy compliance state
az policy state list --filter "ComplianceState eq 'NonCompliant'"

# Export security score
az security secure-score-controls list --output table
```

### GCP Example
```bash
# List high severity findings
gcloud scc findings list --filter="severity='HIGH' AND state='ACTIVE'"

# Get organization policies
gcloud resource-manager org-policies list --organization=ORGANIZATION_ID

# Export IAM policy
gcloud projects get-iam-policy PROJECT_ID --format=json

# Get security recommendations
gcloud recommender recommendations list --recommender=google.iam.policy.Recommender --location=global
```