# Retrieve AWS Foundational Security Best Practices 

```bash
aws securityhub list-security-control-definitions \
  --standards-arn "arn:aws:securityhub:ap-southeast-2::standards/aws-foundational-security-best-practices/v/1.0.0" \
  --region ap-southeast-2 \
  --max-items 1000 \
  --output json > aws-foundational-security-all-controls.json
```