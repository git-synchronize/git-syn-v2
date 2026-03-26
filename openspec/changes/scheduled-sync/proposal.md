## Why

Users of git-syn currently rely exclusively on the pre-push hook to synchronize remotes, meaning mirrors only update when a developer actively pushes. For hosted or server-side mirrors, this leaves remotes stale between pushes and provides no resilience if the hook is bypassed or not installed.

## What Changes

- Add a new `git syn daemon` command that runs as a long-lived process and periodically synchronizes all configured remotes on a user-defined schedule.
- Add a `sync_interval` field to the global config (`~/.config/git-syn/config.yaml`) accepting a Go duration string (e.g., `5m`, `1h`) or cron expression.
- The daemon watches the repository and triggers the same push logic used by `pre-push` at each interval tick.
- Graceful shutdown on `SIGTERM`/`SIGINT` — in-flight push operations are allowed to complete before exit.
- Optional: a `--once` flag to run a single sync cycle and exit (useful for external cron invocation).

## Capabilities

### New Capabilities

- `daemon`: Long-lived daemon process that synchronizes remotes on a configurable schedule, with start/stop lifecycle, signal handling, and interval-based triggering.

### Modified Capabilities

- `config`: Add `sync_interval` configuration field to the global config schema.

## Impact

- **New file**: `cmd/daemon.go` — cobra subcommand and ticker/scheduler loop.
- **Modified**: `cmd/config.go` — extend `Config` struct with `SyncInterval` field.
- **Modified**: `cmd/root.go` — register `daemonCmd`.
- **Shared**: existing push logic in `cmd/pre-push.go` must be extractable (or already callable) by the daemon loop.
- **Dependencies**: no new third-party dependencies required; uses Go standard library (`time`, `os/signal`, `context`).
- **Deployment**: daemon is intended to be managed by systemd, launchd, or an external supervisor — no bundled service file in this change.
