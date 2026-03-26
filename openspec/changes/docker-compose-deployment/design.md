## Context

Users need a git server to push to as a sync target. Docker Compose is the lowest-friction way to run a self-hosted server for users who aren't on Kubernetes. The stack only needs to serve git repositories over HTTP — TLS is delegated to a reverse proxy (nginx, Caddy, Traefik) that the user provides.

## Goals / Non-Goals

**Goals:**
- A minimal `docker-compose.yml` that serves git repos over HTTP via `git-http-backend`
- A `.env.example` with all required variables documented
- A brief README with quickstart instructions

**Non-Goals:**
- TLS termination inside the compose stack
- Authentication (git-http-backend can be fronted by basic auth via nginx, but that's the user's concern)
- Web UI

## Decisions

**Use `nginxinc/nginx-unprivileged` + `git-http-backend` via CGI (or a purpose-built image like `jkarlos/git-server-docker`).**
`git http-backend` is the canonical git HTTP server. A purpose-built image bundles nginx + git-http-backend and requires zero custom configuration. Alternative: a bare `git daemon` image. Rejected — `git daemon` uses the git protocol (port 9418), not HTTP; users can't push to it over HTTPS through a standard proxy.

**Single service, single named volume.**
Keep complexity minimal. Multiple repos are served from subdirectories of the mounted volume. Alternative: one container per repo. Rejected — operationally heavy for a DR mirror.

**Port 80 exposed on the host.**
The expected deployment pattern is a reverse proxy on the same host mapping a domain to this port. Exposing 8080 would require an extra proxy config step for no benefit.

## Risks / Trade-offs

[No authentication by default] → The stack is intended for a trusted network or behind a proxy with auth. Document this clearly in the README. Users who need auth should front it with nginx basic auth or an OAuth proxy.

[Image freshness] → A specific image version should be pinned in `docker-compose.yml` to avoid surprise breaking changes. Use a recent stable tag with a comment noting the upgrade path.
