Run all three validation stages in sequence. Stop at the first failure — do not continue to later stages if an earlier one fails.

Stage 1 — Build: run /build-mcp
Stage 2 — Test: run /test-mcp  
Stage 3 — Lint: run /lint-mcp

After all three pass, print a summary:
```
BUILD  ✓
TEST   ✓ (X tests across Y packages)
LINT   ✓ (golangci-lint clean, gosec clean)
```

If any stage fails, print which stage failed and stop. Do not mark implementation
of the current tool as complete until all three stages pass.
