## Why

git-syn-v2 is a Go rewrite of the original C tool but carries over no test coverage — the legacy C tests in `test/` target the old implementation and cannot exercise the Go code. Without tests, regressions in core parsing and push logic go undetected and refactors (like extracting shared push logic for the daemon) are risky.

## What Changes

- Add Go unit test files (`*_test.go`) in the `cmd` package covering the most critical, pure-function units: config loading/parsing, `.gitremotes` parsing, and URL validation.
- Add `github.com/stretchr/testify` as a test-only dependency for table-driven assertion clarity (optional but recommended).
- The existing `mise run test` task already runs `go test ./...` — no build system changes needed.
- Legacy C tests in `test/` are left in place as historical artifacts; they are not wired into `go test`.
- Add a brief testing section to `README.md` documenting how to run tests.

## Capabilities

### New Capabilities

- `unit-tests`: Go unit test suite covering config parsing, `.gitremotes` parsing, and URL validation, runnable via `go test ./...` / `mise run test`.

### Modified Capabilities

<!-- None — no existing spec-level behavior is changing. -->

## Impact

- **New files**: `cmd/config_test.go`, `cmd/util_test.go`
- **Modified**: `go.mod` / `go.sum` — add `github.com/stretchr/testify` (test-only)
- **Modified**: `README.md` — add testing instructions
- **No production code changes** in this proposal; any refactoring needed to make code testable (e.g., accepting `ActiveConfig` as a parameter rather than using the global) is in scope as a minimal, targeted change.
