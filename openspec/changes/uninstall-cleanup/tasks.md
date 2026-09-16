# Tasks

## 1. Register the `uninstall` command and implement remote cleanup
- [x] 1.1 Register `uninstallCmd` to root via `init()` in `cmd/init.go`
- [x] 1.2 Implement `cleanGitremotes(repoPath string)` in `cmd/init.go` that parses `.gitremotes` and deletes each entry from the git configuration using `repo.DeleteRemote(name)`

## 2. Implement the `--clean` flag for `uninstall`
- [x] 2.1 Add the `clean` bool flag to `uninstallCmd` in `cmd/init.go`

## 3. Implement the `uninstall` command logic
- [x] 3.1 Update the `Run` function of `uninstallCmd` to remove the `pre-push` hook, and conditionally call `cleanGitremotes` if the `--clean` flag is set

## 4. Verification and manual testing
- [x] 4.1 Verify `git-syn uninstall` (no flags) leaves `.gitremotes` and `.git/config` entries intact and prints the warning message
- [x] 4.2 Verify `git-syn uninstall --clean` removes `.gitremotes`, removes `.git/config` entries, and prints the clean message
- [x] 4.3 Verify `git-syn uninstall --clean` succeeds when `.gitremotes` is already absent
