## Why

`git-syn uninstall` only removes the pre-push hook, leaving `.gitremotes` and the git-syn-registered remotes in `.git/config` behind. Users who want a clean removal must manually edit both files, which is inconsistent with the managed install experience.

## What Changes

- Add `--clean` flag to `uninstall`: when set, also delete `.gitremotes` and remove all remotes listed in `.gitremotes` from `.git/config`
- Default behavior (no `--clean`) remains hook-removal only — non-destructive — but now prints a warning listing leftover artifacts and the command to clean them
- Success messages:
  - Without `--clean`: `Removed git hooks. Git SYN uninstalled. Run with --clean to also remove .gitremotes and remote config entries.`
  - With `--clean`: `Removed git hooks and remotes. Git SYN fully uninstalled.`
- If `.gitremotes` does not exist when `--clean` is passed, proceed silently (idempotent)

## Capabilities

### New Capabilities

<!-- none -->

### Modified Capabilities

- `install-command`: the `uninstall` path gains a `--clean` flag and updated output messaging; the cleanup logic (reading `.gitremotes`, removing remotes from `.git/config`) is a behavioral requirement change

## Impact

- `cmd/init.go`: `uninstallCmd` updated — add `--clean` flag, add cleanup logic, update output strings
- No new dependencies
