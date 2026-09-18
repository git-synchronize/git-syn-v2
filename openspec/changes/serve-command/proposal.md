## Why

git-syn is currently client-only: it pushes to remotes configured in `.gitremotes`, but hosting a remote to push to means standing up separate infrastructure by hand. `deploy/docker-compose/` already does this for HTTP (nginx reverse-proxying to `git-http-backend` via a small standalone Go program, `cgi-bridge/`, that lives outside git-syn's own module and is only used by that one deployment example). Folding that capability into git-syn itself as a `serve` command means one binary does both ends of the sync loop, and the docker-compose deployment (and any future deployment) gets it for free instead of maintaining a second, separate Go program.

## What Changes

- Add a `git-syn serve` command (`cmd/serve.go`) that runs an HTTP server hosting bare git repositories under a configured root, for both clone and push:
  - `--path` (default: current directory): repository root, equivalent to `GIT_PROJECT_ROOT`
  - `--listen` (default: `:8080`): bind address
  - Locates the external `git-http-backend` binary via `git --exec-path` rather than a hardcoded filesystem path, since that path differs across distributions (the deploy-only `cgi-bridge` hardcoded the Debian path, `/usr/lib/git-core/git-http-backend`, which doesn't exist on every system `git-syn serve` might run on)
  - Serves it via the standard library's `net/http/cgi` package, the same approach `cgi-bridge` already uses and that was verified working (clone, push up to 100MB, restart persistence) in the change that built it
- No built-in authentication, TLS, or access control. Matches git-syn's other commands (`pre-push`, `daemon`) and the current docker-compose deployment, all of which leave auth/TLS to the operator (a reverse proxy, in the documented deployment).
- Update `deploy/docker-compose/git-http/` to run `git-syn serve` instead of building and running the separate `cgi-bridge` program: the image becomes git-syn plus the `git` package, no second Go module or multi-stage build for a bridge program. nginx stays in front (unchanged reverse-proxy/TLS story in the README) rather than exposing `git-syn serve` directly.
- Document the command in `doc/man/git-syn.1.md` and the root `README.md`'s command list, matching the existing entries for `daemon` and other commands.

**Explicitly out of scope**: a native, `git-http-backend`-free implementation using `go-git`'s server-side protocol support. `go-git` v5 (git-syn's current dependency) has only low-level protocol plumbing with no HTTP wrapper; `go-git` v6 adds exactly that (`backend/http.Backend`) but is alpha-only with no stable release and unconfirmed push-path test depth. Revisit once v6 stabilizes; see `design.md`.

## Capabilities

### New Capabilities

- `serve-command`: a `git-syn serve` command hosting git repositories over HTTP for clone and push, replacing the deploy-only `cgi-bridge` program.

### Modified Capabilities

- `docker-compose-deployment`: the `deploy/docker-compose/git-http/` image runs `git-syn serve` instead of building and running its own `cgi-bridge` program.

## Impact

- **New**: `cmd/serve.go`, `cmd/serve_test.go`
- **Modified**: `cmd/root.go` (command registration, if not self-registering via `init()` alone), `doc/man/git-syn.1.md`, `README.md`
- **Modified**: `deploy/docker-compose/git-http/Dockerfile`, `entrypoint.sh`, `README.md` (deploy-specific quickstart)
- **Removed**: `deploy/docker-compose/git-http/cgi-bridge/` (its logic moves into `cmd/serve.go`)
- No new Go dependencies (uses only the standard library, same as `cgi-bridge` did)
