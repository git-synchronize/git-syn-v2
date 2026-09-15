## Why

The `go-unit-testing` change intentionally scoped itself to pure functions (`parseGitremotes`, `validateRemoteURL`, `loadConfigFromPath`), leaving `init.go`, `pre-push.go`'s push logic, and `remote.go`'s `removeFromGitremotes` untested. The `cmd` package now sits at 32.8% statement coverage, and there is no CI job that runs `go test` or `mise run check` at all — `.gitlab-ci.yml` only builds a Docker image and publishes docs, both gated to `main`. A broken test suite currently blocks nothing. Regressions in hook installation, remote registration, and parallel/sequential push logic go undetected, and refactors in those areas (several are in progress: `uninstall-cleanup`, `scheduled-sync`, `remote-management-subcommands`, `pre-push-command`) carry more risk than they need to.

## What Changes

- Refactor `init_repo` (`cmd/init.go`) to return `error` instead of calling `log.Fatalf` internally, matching the pattern already used by its sibling `cleanGitremotes`. The `install` command's `Run` closure translates the returned error to `log.Fatalf`/`os.Exit(1)` at the top, so CLI-visible behavior is unchanged.
- Apply the same treatment to the business logic inside the `remote add`/`remove`/`list` `Run` closures in `cmd/remote.go`: extract the non-cobra logic (URL validation + `.gitremotes` update + git config registration/removal/listing) into functions that return `error` rather than calling `log.Fatalf` or `os.Exit` inline. `Run` closures keep the same exit-code behavior.
- Add unit tests for `removeFromGitremotes` (`cmd/remote.go`) — pure string/file manipulation, no fixtures needed.
- Add unit tests for `pushRemote`, `syncAll`, `runParallel`, `runSequential` (`cmd/pre-push.go`) and `init_repo`, `cleanGitremotes` (`cmd/init.go`), using `go-git`'s `PlainInit` to create local bare repos as fixture "remotes" in `t.TempDir()` — no network access required.
- Add a CI job to `.gitlab-ci.yml` that runs `mise run check` (fmt + lint + test) on all branches/MRs, not just `main`.
- `go test ./... -cover` output remains visible in CI logs for coverage visibility; no hard coverage percentage gate is introduced in this change.

**Out of scope**: testing cobra command wiring itself (`Execute`, flag registration, `Run` closures as black boxes) — consistent with the `go-unit-testing` change's original non-goals. Only the extracted/refactored business logic underneath gets unit tests. `main.go` (a two-line entrypoint) is also out of scope.

## Capabilities

### New Capabilities

<!-- None — this change extends test coverage rather than introducing a new capability area. -->

### Modified Capabilities

- `unit-tests`: Extends the suite added by `go-unit-testing` to cover `removeFromGitremotes`, `pushRemote`, `syncAll`, `runParallel`, `runSequential`, `init_repo`, and `cleanGitremotes`; adds a requirement that `init_repo` and the `remote` subcommands' business logic return errors instead of calling `log.Fatalf`, to keep them testable; adds a requirement that CI runs the test suite on every branch/MR.

## Impact

- **Modified**: `cmd/init.go` — `init_repo` signature changes from `func(path string)` to `func(path string) error`; its caller (`installCmd.Run`) is updated to handle the returned error.
- **Modified**: `cmd/remote.go` — business logic extracted from `remoteAddCmd`/`remoteRemoveCmd`/`remoteListCmd` `Run` closures into error-returning helper functions; closures updated to call them and translate errors to `log.Fatalf`/`os.Exit(1)`.
- **New files**: `cmd/init_test.go`, `cmd/pre-push_test.go`, additions to `cmd/util_test.go` or a new `cmd/remote_test.go` for `removeFromGitremotes`.
- **Modified**: `.gitlab-ci.yml` — add a `test` job running `mise run check`, not gated to `main`.
- **No changes** to CLI-visible behavior, flags, exit codes, or output messages.
