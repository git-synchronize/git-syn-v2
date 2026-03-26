## 1. Shared utilities

- [x] 1.1 Create `cmd/util.go` and extract `validateRemoteURL(url string) error` from `cmd/init.go`
- [x] 1.2 Extract `appendToGitremotes(path, name, url string) error` into `cmd/util.go`
- [x] 1.3 Extract `registerRemoteInConfig(repo *git.Repository, name, url string) error` into `cmd/util.go`
- [x] 1.4 Create `parseGitremotes(path string) ([]remoteEntry, error)` in `cmd/util.go` (shared with pre-push-command if both are implemented together)

## 2. remote parent command

- [x] 2.1 Create `cmd/remote.go` with a `remoteCmd` cobra.Command (no Run — parent only)
- [x] 2.2 Register `remoteCmd` in `cmd/root.go` `init()`
- [x] 2.3 Add `--path` persistent flag on `remoteCmd` (inherited by all subcommands)

## 3. remote add

- [x] 3.1 Implement `remoteAddCmd` that validates the URL scheme via `validateRemoteURL`
- [x] 3.2 Check for duplicate name in `.gitremotes`; print notice and exit 0 if duplicate
- [x] 3.3 Call `appendToGitremotes` and `registerRemoteInConfig` on success
- [x] 3.4 Print `Added remote '<name>'.` on success

## 4. remote remove

- [x] 4.1 Implement `remoteRemoveCmd` that reads `.gitremotes` and errors if name not found
- [x] 4.2 Remove the matching stanza from `.gitremotes` (line-by-line rewrite)
- [x] 4.3 Deregister the remote from `.git/config` using `go-git`
- [x] 4.4 Print `Removed remote '<name>'.` on success

## 5. remote list

- [x] 5.1 Implement `remoteListCmd` that reads `.gitremotes` and prints each name and URL
- [x] 5.2 Cross-reference `.git/config`; flag remotes missing from config with `[unregistered]`
- [x] 5.3 Print `No remotes tracked by git-syn.` when `.gitremotes` is empty or absent

## 6. Manual verification

- [ ] 6.1 Verify `remote add` → `remote list` → `remote remove` round-trip leaves no artifacts
- [ ] 6.2 Verify `remote add` with an invalid URL scheme prints a useful error
