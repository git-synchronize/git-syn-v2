# Tasks

## 1. Extract common utilities to `cmd/util.go`
- [x] 1.1 Create `cmd/util.go` and extract `validateRemoteURL(url string) error` from `cmd/init.go`
- [x] 1.2 Extract `appendToGitremotes(path, name, url string) error`
- [x] 1.3 Extract `registerRemoteInConfig(repo *git.Repository, name, url string) error`

## 2. Register the `remote` subcommand
- [x] 2.1 Add `remoteCmd` to `rootCmd` via `init()` in `cmd/remote.go`
- [x] 2.2 Add a persistent `--path` flag to `remoteCmd`

## 3. Implement the `remote add` subcommand
- [x] 3.1 Implement `remoteAddCmd` that validates the URL scheme via `validateRemoteURL`
- [x] 3.2 Implement `remoteAddCmd` logic to update both `.gitremotes` and `.git/config`

## 4. Implement the `remote remove` subcommand
- [x] 4.1 Implement `removeFromGitremotes(path, name string) error` to rewrite the `.gitremotes` file without the target remote stanza
- [x] 4.2 Implement `remoteRemoveCmd` that calls `removeFromGitremotes` and `repo.DeleteRemote(name)`

## 5. Implement the `remote list` subcommand
- [x] 5.1 Implement `remoteListCmd` that parses `.gitremotes` and prints each entry name and URL, with an `[unregistered]` suffix if missing from `.git/config`

## 6. Verification and manual testing
- [ ] 6.1 Verify `remote add` → `remote list` → `remote remove` round-trip leaves no artifacts
- [ ] 6.2 Verify `remote add` with an invalid URL scheme prints a useful error
