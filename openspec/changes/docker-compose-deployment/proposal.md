## Why

Git SYN's purpose is disaster recovery and censorship resistance. Users need a simple, self-hosted mirror to push to — but standing one up currently requires manually composing several tools. A ready-made Docker Compose file removes that friction.

## What Changes

- Add `deploy/docker-compose/docker-compose.yml` with a single `git-serve` service:
  - Builds a small custom image (`deploy/docker-compose/git-http/`) running Apache httpd with `mod_cgid`, serving `git-http-backend` (git's own documented way to serve repos over HTTP)
  - Exposes port 80 (intended to sit behind a reverse proxy for TLS)
  - Mounts a named volume for repository storage
- Add `deploy/docker-compose/README.md` with a quickstart (build and start the stack, create a bare repo as `www-data`, enable `http.receivepack`, add it as a git-syn remote)

Out of scope: TLS termination in compose (use a reverse proxy), web UI, authentication beyond git protocol basics.

**Revision note**: the original design used a public `cirocosta/gitserver-http` image (nginx + fcgiwrap). Verifying it (`docker compose up`, clone, push) surfaced two config bugs (a nonexistent image tag, and env vars the image silently ignored) and, after fixing those, a structural one: `git push` reproducibly hung and failed with an nginx 504, fcgiwrap never responding to `git-receive-pack`. No actively maintained alternative HTTP git server image was found either; the few that exist share the same nginx+fcgiwrap+CGI architecture and the same risk. This change builds a custom image instead, on Apache + `mod_cgid` rather than nginx + fcgiwrap, avoiding that architecture entirely. Verified end to end: clone, push (both the system `git` client and `git-syn`'s own `go-git` HTTP transport), and a restart to confirm repository data persists.

## Capabilities

### New Capabilities

- `docker-compose-deployment`: A `deploy/docker-compose/` directory providing a self-contained compose stack for hosting a git mirror

### Modified Capabilities

<!-- none -->

## Impact

- `deploy/docker-compose/`: `docker-compose.yml`, `README.md`, `git-http/` (Dockerfile, `httpd.conf`, `entrypoint.sh`)
- No Go code changes
- No new Go dependencies
