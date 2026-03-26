## Why

The `git-syn install` command writes a `pre-push` hook that delegates to `git syn pre-push "$@"`, but that subcommand does not exist. Every installed repository has a broken hook — any push silently skips synchronization.

## What Changes

- Add `git-syn pre-push` subcommand that reads `.gitremotes` and pushes to all listed remotes
- Accept standard git pre-push hook input: remote name and URL on stdin, refspecs as arguments
- Push to each remote concurrently; report per-remote success/failure
- Exit 0 if at least one remote succeeds; exit non-zero only if all remotes fail
- Print a summary line per remote: `✔ origin-mirror` or `✘ origin-backup: <error>`

## Capabilities

### New Capabilities

- `pre-push-command`: The `git-syn pre-push` subcommand that reads `.gitremotes` and fans out pushes to all configured remotes, reporting per-remote results

### Modified Capabilities

<!-- none -->

## Impact

- `cmd/pre-push.go`: new file implementing the subcommand
- `cmd/root.go`: register the new subcommand in `init()`
- No new dependencies — uses `go-git` (already in go.mod) for push operations
