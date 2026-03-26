## Context

`.gitremotes` is an INI-format file (same syntax as `.git/config`) written by `git-syn install`. It stores the set of remotes git-syn manages. Currently there is no programmatic interface for modifying it — users must hand-edit both `.gitremotes` and `.git/config`. The `install` command already contains URL validation and remote registration logic that the new subcommands should reuse.

## Goals / Non-Goals

**Goals:**
- Provide `add`, `remove`, and `list` subcommands under a `remote` parent command
- Keep `.gitremotes` and `.git/config` in sync on every mutation
- Reuse existing URL validation and remote registration logic from `cmd/init.go`

**Non-Goals:**
- Renaming remotes
- Changing the URL of an existing remote (remove + add is sufficient)
- Global remote management outside a specific repo

## Decisions

**Extract shared logic into `cmd/util.go`.**
Both `install` and `remote add` need URL scheme validation and remote registration. Duplicating the logic would cause drift. A small `util.go` with `validateRemoteURL`, `appendToGitremotes`, and `registerRemote` functions is cleaner than embedding logic in either command file.

**`remote` is a Cobra parent command with `add`/`remove`/`list` as children.**
This matches git's own `git remote` UX, which users already know. Alternative: flat `remote-add`, `remote-remove`, `remote-list` commands. Rejected — harder to discover and inconsistent with git conventions.

**`add` is idempotent: re-adding an existing name is a no-op with a notice, not an error.**
Consistent with `git-syn install` which skips duplicate remotes silently. Returning a non-zero exit code for a re-add would break scripts that call `remote add` defensively.

**`remove` only removes remotes that appear in `.gitremotes`.**
It must not touch remotes that were added to `.git/config` by other means. Limiting scope to `.gitremotes`-managed remotes prevents accidental removal of the user's primary `origin`.

## Risks / Trade-offs

[`.gitremotes` parse/write round-trip] → Writing the INI file after removing an entry requires careful handling to avoid stripping comments or reformatting user edits. Use a line-by-line approach: read all lines, filter out the target stanza, rewrite. Simple and predictable.

[Concurrent modification of `.gitremotes`] → Not addressed; single-user CLI tools rarely need file locking. Document that concurrent use is unsupported.
