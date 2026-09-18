## Context

git-syn currently only pushes to remotes; nothing in the CLI hosts one. `deploy/docker-compose/git-http/cgi-bridge/main.go` proves the approach (nginx reverse-proxies to a Go program that runs the external `git-http-backend` binary as a CGI child process via `net/http/cgi`), verified end to end including push at up to 100MB. That program is a standalone Go module, unrelated to git-syn's own, used only by the one deployment example. This change moves that same approach into git-syn's own command tree.

A fully native implementation (no external `git-http-backend` binary, using `go-git`'s own protocol logic) was investigated and rejected for now: `go-git` v5 (git-syn's current dependency) exposes only low-level server-side protocol plumbing (`plumbing/transport/server`), with no HTTP wrapper — building one means hand-rolling smart-HTTP glue (info/refs negotiation, capability advertisement, content-type handling) from scratch, with real risk of subtle protocol bugs. `go-git` v6 ships exactly the missing piece (`backend/http.Backend`, a complete `http.Handler`) but is alpha-only (through `v6.0.0-alpha.5`), with no stable release and unconfirmed push-path test depth. Pinning a new feature to an alpha dependency isn't worth it when the subprocess approach is already proven working.

## Goals / Non-Goals

**Goals:**
- `git-syn serve` hosts bare repositories under a configured path over HTTP, supporting both clone and push
- Locate `git-http-backend` portably (not a hardcoded, distribution-specific path)
- Reuse the exact serving approach already verified in the `cgi-bridge` program, so its testing (clone, push at multiple sizes, restart persistence) carries over rather than needing to be redone from a different implementation
- `deploy/docker-compose/` uses this command instead of maintaining a separate program

**Non-Goals:**
- A native, `go-git`-only implementation with no external `git-http-backend` dependency (tracked as a future follow-up once `go-git` v6 stabilizes, not attempted here)
- Authentication, TLS, or per-repository access control (consistent with every other git-syn command; the operator's reverse proxy is where this belongs, as already documented for the docker-compose deployment)
- A web UI or repository management commands (creating bare repos, enabling `http.receivepack`) — those stay manual (`git init --bare`, `git config http.receivepack true`), matching the current docker-compose quickstart

## Decisions

**Command: `git-syn serve --path <dir> --listen <addr>`.**
`--path` (default: current directory) mirrors every other command's convention (`install`, `uninstall`, `remote`, `pre-push`, `daemon` all take `--path`). `--listen` (default `:8080`) names the bind address; `serve` was chosen over alternatives like `http-backend` or `server` as the shortest name that reads naturally (`git-syn serve`) and doesn't overload "server" with connotations of a long-running daemon process distinct from what this does (it is one, but so is `daemon`, and disambiguating by name — `daemon` schedules pushes, `serve` hosts a target to push to — reads clearly).

**Locate `git-http-backend` via `git --exec-path`, not a hardcoded path.**
The deploy-only `cgi-bridge` hardcoded `/usr/lib/git-core/git-http-backend`, correct for Debian but not universal (Alpine and others differ, and some distributions don't ship `git-http-backend` in the `git` package at all — discovered directly while building the docker-compose deployment). `git --exec-path` is git's own authoritative answer to "where are my core binaries," portable across every git installation. `serve` runs this once at startup, joins `git-http-backend` onto the result, and fails fast with a clear error if the binary isn't there (rather than a mysterious 500 on the first request).

**Serve via `net/http/cgi`, unchanged from `cgi-bridge`.**
No FastCGI protocol, no third-party bridge library, standard library only. This is a straight move of already-verified logic, not a rewrite: the same `cgi.Handler{Path: ..., Env: ...}` pattern, now parameterized by the `--path` flag instead of a hardcoded `/var/lib/git`.

**No graceful-shutdown-on-SIGINT beyond what `http.ListenAndServe` already provides at process level.**
`daemon` already establishes the project's pattern for long-running commands: `signal.NotifyContext` for clean shutdown on `SIGTERM`/`SIGINT`. `serve` follows the same pattern (`http.Server` with `Shutdown(ctx)` called from a signal handler) rather than inventing a different one, and the e2e test confirms the process exits cleanly rather than being killed.

**`deploy/docker-compose/git-http/` runs `git-syn serve`; nginx stays in front.**
Removing nginx and exposing `git-syn serve` directly would be a bigger change to the deployment's TLS story (the README currently documents fronting port 80 with a reverse proxy for HTTPS) for no immediate benefit — `serve` has no TLS support of its own, so *something* still needs to terminate HTTPS in front of it if a deployment wants that. Keeping nginx as that reverse proxy (now proxying to `git-syn serve` instead of `cgi-bridge`) is the smaller, lower-risk change. The image build simplifies either way: no more multi-stage Go build for a separate bridge program, just install git-syn and the `git` package.

## Risks / Trade-offs

[Subprocess-per-request, not a native implementation] → `git-http-backend` is spawned fresh per CGI invocation, same as `cgi-bridge` did and same as git's own documented Apache/nginx patterns do. Acceptable; this is standard practice for git-http-backend serving, not a git-syn-specific compromise.

[No auth means anyone reaching the port can read and, once `http.receivepack` is set, write to a repo] → Unchanged from the current docker-compose deployment's risk profile; documented, not solved, consistent with every other git-syn command.

[Testing the HTTP server path in CI] → `serve` can be unit-tested for flag/path-resolution logic, but exercising the real server (clone/push against it) needs an actual git-http-backend binary and a real TCP listener, similar to the existing `TestPrePushHookDelegation` e2e test. That test already runs in CI (the `test` job's image has `git` installed), so this isn't a new CI capability gap — the docker-compose stack's own build/verification (Docker, not available in CI) remains the one thing that stays manual-only.

## Future Follow-up

Once `go-git` v6 reaches a stable release with confirmed push-path (`git-receive-pack`) test coverage in `backend/http.Backend`, revisit `serve` as a candidate to drop the external `git-http-backend` dependency entirely, using `go-git`'s own object storage and protocol handling instead of shelling out. Not scheduled; noted here so it isn't rediscovered from scratch.
