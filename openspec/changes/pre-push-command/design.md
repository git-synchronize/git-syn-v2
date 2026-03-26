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
- Modifying or filtering refspecs before forwarding to mirrors

## Decisions

**Use `go-git` for push operations (same library as `install`).**
Keeps the dependency surface flat. Alternative: shell out to `git push` directly. Rejected because it requires `git` on PATH and introduces subprocess management complexity; `go-git` is already vendored.

**Parse `.gitremotes` directly rather than reading `.git/config`.**
`.gitremotes` is the source of truth for git-syn managed remotes. `.git/config` may contain additional remotes not intended to be mirrored. Parsing the same INI-like format used by `install` keeps behavior consistent.

**Exit 0 if at least one remote succeeds.**
The hook is a safety net, not a blocker. Requiring all remotes to succeed would make a transient mirror outage block legitimate pushes. Users who want strict behavior can use `on_failure: abort` once global config is implemented.

**Concurrency: goroutine per remote, collect results.**
Pushing to N remotes sequentially penalizes the common case. A goroutine per remote with a `sync.WaitGroup` and a result channel is straightforward and sufficient at this scale.

**Parse refspecs from stdin, forward exactly those refs to mirror remotes.**
Git passes the refs being pushed as lines on stdin: `<local-ref> <local-sha> <remote-ref> <remote-sha>`. The command reads all lines before pushing and builds a `[]config.RefSpec` to pass to `PushOptions`. This ensures mirrors receive exactly what the user pushed — not a full sync of all branches. A zero local sha (`0000000...`) signals deletion and maps to a delete refspec (`:refs/heads/<name>`).

**Skip the remote named in `args[0]`.**
Git passes the primary remote name as the first positional argument. That remote is already being pushed to by git itself; pushing to it again from the hook is redundant. The command filters it out of the `.gitremotes` entries before pushing.

**Use `Force: true` on all mirror pushes.**
Mirrors are followers. A rebase or amend on the primary branch would cause a non-fast-forward rejection on the mirror without force. Since the user has already authorized the push to the primary, force-syncing mirrors is the correct behavior.

## Risks / Trade-offs

[go-git push may have auth limitations] → For HTTPS remotes requiring credentials, go-git relies on the system credential helper or explicit credential config. Document this limitation; users can fall back to SSH for mirror remotes.

[Goroutine-per-remote doesn't bound parallelism] → Acceptable for the typical case (2-5 remotes). If users configure dozens of remotes this could spike. A bounded worker pool is a future optimization.

[Force push on mirrors could overwrite divergent mirror state] → Acceptable: mirrors are not the source of truth. If a mirror has diverged (e.g. someone pushed directly to it), force-syncing is the correct resolution. Users who need bidirectional sync are outside the current scope.

## Open Questions

None.
