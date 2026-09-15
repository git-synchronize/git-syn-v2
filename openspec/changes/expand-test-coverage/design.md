## Context

`go-unit-testing` covered the pure, already-testable functions in `cmd` (`parseGitremotes`, `validateRemoteURL`, `loadConfigFromPath`) and explicitly deferred everything touching the filesystem, git operations, or `log.Fatalf`/`os.Exit`. That left `init.go`, `pre-push.go`'s push logic, and `remote.go`'s `removeFromGitremotes` at 0% coverage. Two testability obstacles remain from that change's open questions:

1. `init_repo` (`cmd/init.go`) calls `log.Fatalf` on every error path, which calls `os.Exit(1)` directly — any test that exercises an error branch kills the test binary. Its sibling `cleanGitremotes` already returns `error` instead, so the fix is to bring `init_repo` in line with an established pattern, not invent a new one.
2. The `remote add`/`remove`/`list` `Run` closures in `cmd/remote.go` mix cobra wiring with business logic (URL validation, `.gitremotes` I/O, git config registration) and call `log.Fatalf`/`os.Exit` inline, with no separable function to test.

`handleResults` (`cmd/pre-push.go`) turned out to already be testable — it returns `int` and only calls `fmt.Printf`, with no `os.Exit` — so the prior design doc's concern about it is stale; no changes are needed there.

Four other openspec changes (`uninstall-cleanup`, `scheduled-sync`, `remote-management-subcommands`, `pre-push-command`) are in progress and touch this same package. Per discussion with the user, this change proceeds against current code now rather than blocking on those landing.

## Goals / Non-Goals

**Goals:**
- Make `init_repo` and the `remote` subcommands' business logic testable by having them return `error` instead of calling `log.Fatalf`/`os.Exit` internally, with `Run` closures translating errors to the same exit behavior as today.
- Add unit tests for `removeFromGitremotes`, `pushRemote`, `syncAll`, `runParallel`, `runSequential`, `init_repo`, and `cleanGitremotes`.
- Use real local git repos (via `go-git`'s `PlainInit`) as test fixtures instead of mocking `git.Repository`.
- Add a CI job that runs the test suite (via `mise run check`) on every branch and MR, closing the gap where nothing currently runs `go test` in CI.

**Non-Goals:**
- Testing cobra command wiring itself (`Execute`, flag parsing, `Run` closures as black boxes) — consistent with `go-unit-testing`'s original scope.
- Testing `main.go` (a two-line entrypoint).
- Introducing a hard coverage percentage gate in CI — the package is small and still under active churn from the four in-progress changes; a threshold now is more likely to cause friction than catch real regressions. Coverage stays visible via `go test ./... -cover` output in CI logs.
- Changing any CLI-visible behavior, flags, exit codes, or output messages.
- Reworking the docker `build` or `pages` jobs already in `.gitlab-ci.yml`.

## Decisions

### 1. Refactor `init_repo` to return `error`
**Decision**: Change `init_repo(path string)` to `init_repo(path string) error`, replacing each `log.Fatalf` call with `return fmt.Errorf(...)`. `installCmd`'s `Run` closure calls `log.Fatalf` on a non-nil return, preserving today's exit-1-with-message behavior.
**Rationale**: Mirrors the pattern `cleanGitremotes` already uses in the same file. Minimal, mechanical, and makes every error branch in `init_repo` reachable from a test.
**Alternative considered**: Test `init_repo` out-of-process (invoke the compiled binary as a subprocess and assert on exit code/stderr). Rejected — much slower, harder to set up fixtures for, and inconsistent with how the rest of the suite tests business logic directly.

### 2. Extract `remote add`/`remove`/`list` business logic into error-returning functions
**Decision**: Pull the non-cobra logic out of each `Run` closure into a function (e.g. `addRemote(path, name, url string) error`, `removeRemote(path, name string) error`, `listRemotes(path string) ([]remoteStatus, error)` or similar shape) that returns an error/result instead of writing to stdout/stderr and calling `os.Exit` directly. `Run` closures call these and handle output/exit as before.
**Rationale**: Same motivation as `init_repo` — makes the logic unit-testable without a subprocess. Keeps the diff mechanical: move code, don't change behavior.
**Alternative considered**: Leave `remote.go` as-is and only test `removeFromGitremotes` (already a standalone function). Rejected — it would leave the higher-value integration paths (duplicate-remote detection, git config registration/removal, unregistered-remote listing) untested, which is most of what makes `remote.go` risky to change.

### 3. Real git repo fixtures via `go-git` `PlainInit`, not mocks
**Decision**: For tests needing a `*git.Repository`, create one with `git.PlainInit(t.TempDir(), false)` and, where a remote target is needed, a second bare repo (`git.PlainInit(otherDir, true)`) registered as a local-path remote.
**Rationale**: `go-git` doesn't expose a small interface for `*git.Repository` that's practical to mock, and testing against the real library exercises the actual code paths (including error types like `git.NoErrAlreadyUpToDate` and `git.ErrRemoteNotFound`) with no network access required.
**Alternative considered**: Define a narrow interface (`type gitRepo interface {...}`) and mock it. Rejected as premature abstraction for this change — it would require touching every call site that takes `*git.Repository` for a testing-only concern.
**Risk to verify early**: confirm `go-git` v5.11.0 pushes to a local filesystem path remote without requiring an explicit side-effect import (e.g. `github.com/go-git/go-git/v5/plumbing/transport/file`). Spike this in the first `pre-push_test.go` test before writing the rest.

### 4. CI runs `mise run check` via installed `mise`, not raw `go` commands
**Decision**: The new CI job installs `mise` and runs `mise run check`, rather than reimplementing fmt/vet/test steps directly with a `golang` base image.
**Rationale**: Keeps one source of truth for what "check" means — local dev and CI run the exact same task definitions in `.mise.toml`. Avoids drift between what a contributor runs locally and what CI enforces.
**Alternative considered**: Use an official `golang` image and run `gofmt -l .`, `go vet ./...`, `go test ./...` directly. Rejected — duplicates task logic in two places (`.mise.toml` and `.gitlab-ci.yml`), which will silently diverge over time.

### 5. No `only: [main]` restriction on the new job
**Decision**: The new `test` job has no branch restriction, so it runs on every push and MR pipeline by default.
**Rationale**: The entire point is to catch regressions before merge; gating it to `main` only (like today's `build`/`pages` jobs) would defeat the purpose.

## Risks / Trade-offs

- **Goroutine non-determinism in `runParallel`** → results arrive on a buffered channel in non-deterministic order. Mitigation: tests assert on the *set* of results (success/failure per remote name) via the channel drained into a map, not on ordering.
- **Refactor touches code that in-progress sibling changes also touch** (`remote-management-subcommands`, `uninstall-cleanup` both likely touch `remote.go`/`init.go`) → merge/rebase friction. Mitigation: keep the `log.Fatalf`-to-`error` refactor mechanical and minimal (no renaming, no behavior changes) to keep conflicts easy to resolve.
- **`go-git` local-path transport quirks** (file locking on the bare repo, `.git` vs bare directory conventions) could make fixtures flakier than expected. Mitigation: spike the fixture pattern first (see Decision 3) before writing the full test matrix.
- **No coverage gate means the new tests aren't self-enforcing** — a future PR could remove business logic without a test catching it via a coverage drop. Accepted trade-off per the non-goal above; revisit once the four in-progress changes land and the package stabilizes.

## Migration Plan

1. Refactor `init_repo` to return `error`; update `installCmd.Run`.
2. Extract error-returning helpers from `remote add`/`remove`/`list` `Run` closures; update closures to call them.
3. Spike a `go-git` local bare-repo fixture pattern; confirm push/remote semantics work without network access.
4. Write `cmd/pre-push_test.go` (`pushRemote`, `syncAll`, `runParallel`, `runSequential`).
5. Write `cmd/init_test.go` (`init_repo`, `cleanGitremotes`).
6. Add tests for `removeFromGitremotes` (new `cmd/remote_test.go` or extend `cmd/util_test.go`).
7. Add the `test` job to `.gitlab-ci.yml` running `mise run check`.
8. Run `mise run check` locally; confirm `go test ./...` and `go vet ./...` are clean.

## Open Questions

- Does `go-git` v5.11.0 need an explicit `plumbing/transport/file` import for local-path pushes, or is it registered by default via the `go-git` root package? → Resolve during the spike in step 3.
- Should the CI job pin a specific Go version instead of following `.mise.toml`'s `go = "latest"`? → Lean: no, keep CI and local dev on the same floating pin `.mise.toml` already uses; revisit if CI flakiness from toolchain drift shows up.
