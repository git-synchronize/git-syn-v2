## 1. Add --clean flag

- [x] 1.1 Add `--clean` boolean flag to `uninstallCmd` in `cmd/init.go`

## 2. Cleanup logic

- [x] 2.1 Implement `cleanGitremotes(repoPath string) error` that reads `.gitremotes`, removes each listed remote from `.git/config`, then deletes `.gitremotes`
- [x] 2.2 Handle missing `.gitremotes` silently (idempotent)
- [x] 2.3 Call `cleanGitremotes` in `uninstallCmd.Run` when `--clean` is set

## 3. Update output messages

- [x] 3.1 Change default uninstall success message to: `Removed git hooks. Git SYN uninstalled. Run with --clean to also remove .gitremotes and remote config entries.`
- [x] 3.2 Add `--clean` success message: `Removed git hooks and remotes. Git SYN fully uninstalled.`

## 4. Manual verification

- [ ] 4.1 Verify `git-syn uninstall` (no flags) leaves `.gitremotes` and `.git/config` entries intact and prints the warning message
- [ ] 4.2 Verify `git-syn uninstall --clean` removes `.gitremotes`, removes `.git/config` entries, and prints the clean message
- [ ] 4.3 Verify `git-syn uninstall --clean` succeeds when `.gitremotes` is already absent
