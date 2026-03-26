## Context

The `git-syn install` command writes a pre-push hook at `.git/hooks/pre-push` that calls `git syn pre-push "$@"`. Git passes push information via stdin (one line per ref: `<local-ref> <local-sha> <remote-ref> <remote-sha>`) and passes the remote name and URL as positional arguments. The subcommand must be a valid git hook delegate: read that input, push to all remotes in `.gitremotes`, and exit with the appropriate code.

## Goals / Non-Goals

**Goals:**
- Implement `git-syn pre-push` as a working git hook delegate
- Push to all remotes in `.gitremotes` and report per-remote results
- Exit non-zero only when all pushes fail (partial success is still success)

**Non-Goals:**
- Configurable push strategy (parallel vs. sequential) — deferred to global-yaml-config change
- Authentication handling beyond what `go-git` provides natively
- Modifying refspecs before forwarding to mirrors

## Decisions

**Use `go-git` for push operations (same library as `install`).**
Keeps the dependency surface flat. Alternative: shell out to `git push` directly. Rejected because it requires `git` on PATH and introduces subprocess management complexity; `go-git` is already vendored.

**Parse `.gitremotes` directly rather than reading `.git/config`.**
`.gitremotes` is the source of truth for git-syn managed remotes. `.git/config` may contain additional remotes not intended to be mirrored. Parsing the same INI-like format used by `install` keeps behavior consistent.

**Exit 0 if at least one remote succeeds.**
The hook is a safety net, not a blocker. Requiring all remotes to succeed would make a transient mirror outage block legitimate pushes. Users who want strict behavior can use `on_failure: abort` once global config is implemented.

**Concurrency: goroutine per remote, collect results.**
Pushing to N remotes sequentially penalizes the common case. A goroutine per remote with a `sync.WaitGroup` and a result channel is straightforward and sufficient at this scale.

## Risks / Trade-offs

[go-git push may have auth limitations] → For HTTPS remotes requiring credentials, go-git relies on the system credential helper or explicit credential config. Document this limitation; users can fall back to SSH for mirror remotes.

[Goroutine-per-remote doesn't bound parallelism] → Acceptable for the typical case (2-5 remotes). If users configure dozens of remotes this could spike. A bounded worker pool is a future optimization.

## Open Questions

- Should we forward stdin to each push, or is stdin only relevant for the hook's own validation? (Likely: stdin only matters if we want to skip no-op pushes — can be addressed in a follow-up.)
