# AWS Security standards
```
aws securityhub describe-standards --output json > aws-security-standards.json
```

# Retrieve AWS Foundational Security Best Practices 

## confirm logged into AWS cli

## get aws-foundational-security-best-practices

```bash
aws securityhub list-security-control-definitions \
  --standards-arn "arn:aws:securityhub:ap-southeast-2::standards/aws-foundational-security-best-practices/v/1.0.0" \
  --region ap-southeast-2 \
  --max-items 1000 \
  --output json > aws-foundational-security-all-controls.json
```


# get aws trusted advisor checks

```bash
aws support describe-trusted-advisor-checks --language "en" --output json > describe-trusted-advisor-checks.json
```


# get AWS Well Architected Lense Summary

```bash
q
```


# security hub confirmance packs
```
https://github.com/awslabs/aws-config-rules/tree/master/aws-config-conformance-packs
```