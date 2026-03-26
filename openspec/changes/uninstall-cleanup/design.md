## Context

`git-syn uninstall` currently calls `os.Remove` on the pre-push hook and prints a success message. It has no knowledge of `.gitremotes` or `.git/config`. The cleanup logic needed by `--clean` is the inverse of what `install` does: read `.gitremotes`, remove each listed remote from `.git/config`, then delete `.gitremotes`.

## Goals / Non-Goals

**Goals:**
- Add `--clean` flag to `uninstall` that removes `.gitremotes` and deregisters its remotes from `.git/config`
- Update default (no-flag) output to warn about leftover artifacts
- Keep the default behavior non-destructive

**Non-Goals:**
- Automatically inferring which `.git/config` remotes were added by git-syn without `.gitremotes` present
- A `--force` flag that removes remotes even if `.gitremotes` is missing

## Decisions

**`--clean` reads `.gitremotes` to know which remotes to remove.**
We only remove remotes that git-syn explicitly manages. Alternative: remove all remotes not in the original git config. Rejected — there is no baseline to compare against; this could destroy user-added remotes.

**If `.gitremotes` is absent and `--clean` is passed, proceed silently.**
Idempotent cleanup is important for scripted use. An error when the file is already gone would be surprising.

**Remove remotes from `.git/config` before deleting `.gitremotes`.**
If the process is interrupted, `.gitremotes` still exists and re-running `--clean` will re-attempt the config cleanup. Reverse order would leave orphaned `.git/config` entries with no way to identify them.

## Risks / Trade-offs

[User manually added a remote with the same name as a `.gitremotes` entry] → The cleanup correctly removes it from `.git/config`, which may be unexpected. Document that `--clean` removes all remotes listed in `.gitremotes` regardless of how they were added.

[Small change surface] → This is a contained modification to `uninstallCmd`. No new files, no new dependencies. Risk of regression is low.
