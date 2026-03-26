## Why

Git SYN is modeled after the Git LFS extension. `git syn install` follows the same pattern as `git lfs install` and has two core goals:

1. **Configure the `pre-push` hook** — write an embedded hook script into `.git/hooks/pre-push` with `0700` permissions
2. **Populate `.git/config` from `.gitremotes`** — write remotes to `.gitremotes`, then register each entry in `.git/config`

Git SYN's purpose is remote repository synchronization for disaster recovery and censorship resistance. The `.gitremotes` file captures all configured remotes so they can be registered in `.git/config` and kept in sync.

The scaffolded `git-syn init` command partially implements goal 2 but has placeholder help text, no hook installation, and no user feedback. This change completes both goals and renames the command to `install` to match the Git LFS convention.

## What Changes

- Rename the command from `init` to `install`
- Replace placeholder `Use`, `Short`, and `Long` descriptions with accurate help text
- Add a `--path` flag (defaults to current directory)
- Embed `pre-push` hook script as a Go string constant with a `{{Command}}` placeholder (Git LFS `hookBaseContent` pattern); validates `git-syn` is in PATH, exits 2 with a warning if not, delegates to `git syn pre-push "$@"` if found
- Write the hook into `.git/hooks/pre-push` and set permissions to `0700`
- Filter remotes to HTTPS and OpenSSH only; skip others with a warning
- Write `.gitremotes`, then append each remote to `.git/config` (skipping duplicates)
- Print `Updated git hooks. Git SYN initialized.` on success

## Capabilities

### New Capabilities

- `install-command`: The `git-syn install` command that installs a `pre-push` hook and writes a `.gitremotes` config file, registering all remotes in `.git/config` for synchronization

### Modified Capabilities

<!-- none -->

## Impact

- `cmd/init.go`: primary file changed (command renamed to `install`; hook content embedded as constant)
- `.gitremotes` file: config artifact written by the command
- No API or dependency changes
