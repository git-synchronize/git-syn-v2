## 1. Refactor `init_repo` to return `error`

- [x] 1.1 Change `init_repo(path string)` to `init_repo(path string) error` in `cmd/init.go`, replacing each `log.Fatalf` call with `return fmt.Errorf(...)`
- [x] 1.2 Update `installCmd`'s `Run` closure to call `log.Fatalf` on a non-nil error from `init_repo`, preserving today's exit code and message text

## 2. Extract testable business logic from `remote add`/`remove`/`list`

- [x] 2.1 Extract the `remote add` logic (URL validation, duplicate check, `.gitremotes` append, git config registration) from `remoteAddCmd.Run` into an error-returning function
- [x] 2.2 Extract the `remote remove` logic (`.gitremotes` lookup, `removeFromGitremotes`, git config removal) from `remoteRemoveCmd.Run` into an error-returning function
- [x] 2.3 Extract the `remote list` logic (`.gitremotes` read, git config cross-reference) from `remoteListCmd.Run` into a function returning a result and error
- [x] 2.4 Update `remoteAddCmd`, `remoteRemoveCmd`, and `remoteListCmd`'s `Run` closures to call the extracted functions and translate errors/results to the same stdout/stderr output and exit codes as today

## 3. Spike the go-git local fixture pattern

- [x] 3.1 Write a throwaway test that creates a repo with `git.PlainInit(t.TempDir(), false)` and a bare repo with `git.PlainInit(t.TempDir(), true)`, registers the bare repo as a remote, and pushes a commit to confirm the local-path transport works without a network
- [x] 3.2 Note in a code comment (or remove if unnecessary) whether an explicit `plumbing/transport/file` import is required for the push to succeed (confirmed not needed: `go-git`'s `client` package registers the `file` protocol by default, and plain filesystem paths parse to that protocol via `parseFile`)
- [x] 3.3 Fold the working pattern into a small test helper (e.g. `newTestRepoPair(t *testing.T) (repo *git.Repository, remoteName, remotePath string)`) usable by tasks 4 and 5

## 4. Unit tests for `cmd/pre-push.go`

- [x] 4.1 Create `cmd/pre-push_test.go` with `package cmd`
- [x] 4.2 Add tests for `pushRemote`: successful push, push when already up to date, push to an unregistered remote name
- [x] 4.3 Add tests for `runParallel`: all remotes succeed; one of several fails while others still complete (assert on the result set, not order)
- [x] 4.4 Add tests for `runSequential`: `OnFailure: "abort"` stops after the first failure; `OnFailure: "warn"` continues through all remotes
- [x] 4.5 Add tests for `syncAll`: dispatches to `runParallel` when `PushStrategy: "parallel"`, to `runSequential` when `PushStrategy: "sequential"`

## 5. Unit tests for `cmd/init.go`

- [x] 5.1 Create `cmd/init_test.go` with `package cmd`
- [x] 5.2 Add tests for `init_repo`: success path creates `.gitremotes` and registers each remote with a supported scheme; a remote with an unsupported scheme is skipped with a warning but does not fail the run; a non-git-repository path returns an error
- [x] 5.3 Add tests for `cleanGitremotes`: removes each listed remote from git config and deletes `.gitremotes`; tolerates a remote already absent from git config; is a no-op (nil error) when `.gitremotes` does not exist

## 6. Unit tests for `cmd/remote.go`

- [x] 6.1 Create `cmd/remote_test.go` with `package cmd`
- [x] 6.2 Add table-driven tests for `removeFromGitremotes`: removing one remote among several preserves the others in order; removing the only remote leaves an empty file; removing a name not present is a no-op returning nil; a missing file returns an error
- [x] 6.3 Add tests for `addRemote`: adds a new remote to `.gitremotes` and registers it in git config; reports `alreadyTracked` without duplicating an existing entry; returns an error for a non-git-repository path
- [x] 6.4 Add tests for `removeRemote`: removes a tracked remote from `.gitremotes` and git config; returns `errRemoteNotTracked` for a name not present
- [x] 6.5 Add tests for `listRemotes`: reports registered vs. unregistered remotes correctly; returns nil/nil when `.gitremotes` is missing or empty

## 7. CI test job

- [x] 7.1 Add a `test` job to `.gitlab-ci.yml` that installs `mise` and runs `mise run check`, with no branch restriction so it runs on every push and merge request pipeline
- [x] 7.2 Confirm the job's output includes `go test ./... -cover` coverage percentages in the log (via `mise run test` invoking `go test ./...`, or by adding `-cover` to the `test` task in `.mise.toml` if not already present)

## 8. Verification

- [x] 8.1 Run `go test ./...` and confirm all tests pass with zero failures
- [x] 8.2 Run `mise run check` (fmt + lint + test) and confirm clean output
- [x] 8.3 Run `go tool cover -func` (or equivalent) and confirm `init.go`, `pre-push.go`, and `remote.go` no longer show 0%-covered business-logic functions
- [ ] 8.4 Push a branch and confirm the new CI `test` job runs and passes on a non-`main` branch (requires pushing to the GitLab remote — left for the user to trigger and confirm)
