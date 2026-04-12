Scaffold a new MCP tool implementation. The argument is the tool name (e.g. `get_module_docs`).

Steps:

1. Read spec.md and find the section for this tool. Extract:
   - Input fields (name, type, required, description)
   - Output fields
   - Source (which API or URL pattern)
   - Fallback logic if any

2. Create `internal/tools/<tool_name>.go` with:
   - Input struct with json tags
   - Handler function signature: `func handle<ToolName>(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error)`
   - Stub implementation that returns a "not implemented" error
   - Registration call added to `internal/tools/registry.go`

3. Create `internal/tools/<tool_name>_test.go` with:
   - At least three table-driven test cases: happy path, missing required input, API error response
   - Tests must fail (stub not implemented) — this is expected and correct at this stage

4. Print the files created and confirm the test file compiles even though tests fail.

Do not implement the handler logic yet — that is the next step after the test file is confirmed.
