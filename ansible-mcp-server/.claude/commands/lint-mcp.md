Run two linters in sequence from the workspace root:

1. `golangci-lint run ./...`
2. `gosec ./...`

For each linter report results as:
- Rule or check ID
- File and line number
- Description of the issue
- Suggested fix if available

Summary at the end:
- "golangci-lint: X issues"
- "gosec: Y issues"
- Overall: PASS (zero issues from both) or FAIL

Do not skip or suppress any findings.
