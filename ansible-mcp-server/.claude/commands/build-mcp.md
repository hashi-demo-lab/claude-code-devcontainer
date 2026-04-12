Run `go build ./...` followed by `go vet ./...` from the workspace root.

Report results clearly:
- If build succeeds and vet passes: print "BUILD OK"
- If build fails: print the compiler error with file and line number, then stop
- If vet reports issues: list each issue with file, line, and description, then stop

Do not proceed to any other step if either command fails.
