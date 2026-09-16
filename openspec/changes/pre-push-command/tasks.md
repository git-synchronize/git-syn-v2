## 1. Core command scaffolding

- [x] 1.1 Create `cmd/pre-push.go` with a `prePushCmd` cobra.Command registered in `init()`
- [x] 1.2 Add `--path` flag (default: current directory) consistent with other commands
- [x] 1.3 Register `prePushCmd` in `cmd/root.go` `init()`

## 2. .gitremotes parsing

- [x] 2.1 Implement `parseGitremotes(path string) ([]remoteEntry, error)` that reads `.gitremotes` and returns name/URL pairs
- [x] 2.2 Return a clear error when `.gitremotes` is absent (direct user to run `git-syn install`)

## 3. Concurrent push logic

- [x] 3.1 Implement `pushRemote(repo *git.Repository, r remoteEntry) error` using `go-git`
- [x] 3.2 Fan out a goroutine per remote using `sync.WaitGroup` and a results channel
- [x] 3.3 Collect results and print per-remote `✔ <name>` or `✘ <name>: <error>` lines

## 4. Exit code logic

- [x] 4.1 Exit 0 if at least one remote push succeeded
- [x] 4.2 Exit non-zero if all remote pushes failed
- [x] 4.3 Exit non-zero if `.gitremotes` is absent or unreadable

## 5. Help text and manual verification

- [x] 5.1 Write accurate `Use`, `Short`, and `Long` descriptions for the command
- [x] 5.2 Manually verify: `git-syn install` followed by a git push delegates to `git-syn pre-push` and pushes to all `.gitremotes` entries
