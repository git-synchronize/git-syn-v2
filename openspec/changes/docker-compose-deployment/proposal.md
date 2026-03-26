## Why

Git SYN's purpose is disaster recovery and censorship resistance. Users need a simple, self-hosted mirror to push to — but standing one up currently requires manually composing several tools. A ready-made Docker Compose file removes that friction.

## What Changes

- Add `deploy/docker-compose/docker-compose.yml` with a single `git-serve` service:
  - Based on a minimal image that serves git repositories over HTTP via `git-http-backend`
  - Exposes port 80 (intended to sit behind a reverse proxy for TLS)
  - Mounts a named volume for repository storage
- Add `deploy/docker-compose/.env.example` documenting required variables (`GIT_REPO_PATH`, `GIT_HTTP_EXPORT_ALL`)
- Add `deploy/docker-compose/README.md` with a quickstart (clone, copy `.env.example` to `.env`, `docker compose up`)

Out of scope: TLS termination in compose (use a reverse proxy), web UI, authentication beyond git protocol basics.

## Capabilities

### New Capabilities

- `docker-compose-deployment`: A `deploy/docker-compose/` directory providing a self-contained compose stack for hosting a git mirror

### Modified Capabilities

<!-- none -->

## Impact

- `deploy/docker-compose/`: new directory with `docker-compose.yml`, `.env.example`, `README.md`
- No Go code changes
- No new Go dependencies
