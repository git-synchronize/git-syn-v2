## Why

`.gitremotes` must currently be edited by hand and users must manually run `git config` to keep `.git/config` in sync. There is no programmatic way to add or remove a sync remote, making the workflow error-prone.

## What Changes

- Add `git-syn remote` subcommand group with three subcommands:
  - `git-syn remote add <name> <url>` — validate URL scheme, append to `.gitremotes`, register in `.git/config`
  - `git-syn remote remove <name>` — remove from `.gitremotes` and deregister from `.git/config`
  - `git-syn remote list` — list remotes from `.gitremotes`, flagging any missing from `.git/config`
- All subcommands accept `--path <dir>` (default: current directory), consistent with `install`/`uninstall`
- `add` rejects non-HTTPS/non-OpenSSH URLs with a clear error (same validation as `install`)
- `add` is idempotent: re-adding an existing remote name is a no-op with a notice

## Capabilities

### New Capabilities

- `remote-management`: The `git-syn remote` subcommand group (`add`, `remove`, `list`) for programmatically managing `.gitremotes` and the corresponding `.git/config` entries

### Modified Capabilities

<!-- none -->

## Impact

- `cmd/remote.go`: new file implementing the `remote` command group and its three subcommands
- `cmd/root.go`: register the `remote` parent command in `init()`
- Shares URL validation logic with `cmd/init.go` — consider extracting to `cmd/util.go`
