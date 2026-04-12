Run `go test ./... -v -race -count=1` from the workspace root.

Report results:
- List each package tested with PASS or FAIL
- For any FAIL: show the test function name, what it tested, and the failure output
- For any panic: show the full stack trace
- Summary line at the end: "X passed, Y failed across Z packages"

A test file that does not compile counts as a failure — show the compiler error.
