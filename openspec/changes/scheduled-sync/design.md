## Context

git-syn currently synchronizes remotes only when the pre-push hook fires. The push logic lives in `cmd/pre-push.go` across three functions: `pushRemote`, `runParallel`, and `runSequential`. The `handleResults` helper reports per-remote outcomes. The global `ActiveConfig` singleton (loaded in `cmd/config.go`) controls push strategy and failure behavior.

A daemon mode was sketched in `doc/daemon.md` for the C-era design but was never implemented in the Go rewrite. The original design used `libuv` for filesystem event watching; the Go rewrite targets a simpler ticker-based approach without a fs-event dependency.

## Goals / Non-Goals

**Goals:**
- `git syn daemon` command runs as a foreground process, pushing all remotes at a fixed interval.
- Interval configurable via `sync_interval` in `~/.config/git-syn/config.yaml` (Go duration string, e.g., `5m`, `1h`).
- `--interval` flag overrides config for one-shot invocations.
- `--once` flag performs a single sync cycle then exits (enables use from external cron).
- Graceful shutdown on `SIGTERM`/`SIGINT`: finish any in-flight push before exit.
- Re-uses existing `runParallel`/`runSequential` push logic without duplication.

**Non-Goals:**
- Filesystem event watching (inotify/kqueue) — ticker is sufficient for the DR/mirror use case.
- Service file generation (systemd unit, launchd plist) — out of scope for this change.
- Daemon double-fork / PID file management — users are expected to manage the process with their init system or `nohup`.
- Per-repository daemon instances beyond what the `--path` flag already supports.

## Decisions

### 1. Ticker-based scheduling over cron expressions
**Decision**: Accept a Go `time.Duration` string (`5m`, `1h30m`) rather than a cron expression.
**Rationale**: The use case is periodic mirror synchronization, not calendar-based scheduling. A duration string is simpler to parse (standard library `time.ParseDuration`), simpler to document, and avoids a third-party cron library dependency. Users who need calendar scheduling can invoke `git syn daemon --once` from their own cron or systemd timer.
**Alternative considered**: `robfig/cron` — adds a dependency and complexity without meaningful benefit for this use case.

### 2. Extract shared push logic into a package-level function
**Decision**: Extract a `syncAll(repo *git.Repository, entries []remoteEntry) int` function from the `pre-push` command's `Run` body. Both `prePushCmd` and `daemonCmd` call it.
**Rationale**: Avoids duplicating the open-repo / parse-gitremotes / dispatch-strategy flow. The extraction is a pure refactor with no behavior change to `pre-push`.

### 3. `--path` flag on `daemonCmd`
**Decision**: Mirror the `--path` flag from `pre-push` on `daemon` as well.
**Rationale**: Allows daemonizing a repo that is not the shell's current directory, consistent with existing CLI conventions.

### 4. Graceful shutdown via `context.Context` + `os/signal`
**Decision**: Use `signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)` to derive a cancellable context. The ticker loop selects on `<-ctx.Done()`.
**Rationale**: Standard Go idiom; no external dependency. In-flight goroutines for parallel push use a `sync.WaitGroup`, so the shutdown waits for the current cycle to drain before returning.

### 5. Default interval
**Decision**: Default `sync_interval` is `0` (disabled). If the field is absent or zero, `daemon` requires an explicit `--interval` flag or errors out with a clear message.
**Rationale**: Opt-in is safer than accidentally backgrounding a push loop after a config upgrade. Users who want always-on behavior set `sync_interval` explicitly.

## Risks / Trade-offs

- **Credential expiry during long-running daemon**: SSH keys or HTTPS tokens may expire. The daemon will report a push failure per the `on_failure` policy but will not attempt re-authentication. → Mitigation: document that credential refresh is the operator's responsibility; `on_failure: warn` (default) keeps the daemon alive on auth failures.
- **Repository churn during push**: If a push cycle overlaps with a large local repack, `go-git` may return transient errors. → Mitigation: errors are logged and the daemon continues to the next tick; no retry within a cycle.
- **Single-instance enforcement**: Nothing prevents multiple daemon processes for the same repo. → Mitigation: out of scope; defer to init system (systemd `Type=simple` with `Restart=on-failure` handles this idiomatically).
- **`--once` vs `pre-push` overlap**: `--once` and `pre-push` are functionally equivalent. Users may be confused about which to use in hook context. → Mitigation: document that `pre-push` is the hook-facing command; `daemon --once` is for scripting/external cron.

## Migration Plan

1. Extract `syncAll` helper from `prePushCmd.Run` — no behavior change.
2. Add `SyncInterval` (`time.Duration`) field to `Config` struct and update `configTemplate`.
3. Implement `cmd/daemon.go` with `daemonCmd`.
4. Register `daemonCmd` in `cmd/root.go` (or its own `init()`).
5. No config migration required — new field is additive and defaults to disabled.

## Open Questions

- Should `daemon` log to stderr only, or support a `--log-file` flag? (Lean: stderr only for now; redirect at the shell/systemd level.)
- Should the sync interval be per-repository (`.gitremotes` extension) or global-only? (Lean: global only for now.)
