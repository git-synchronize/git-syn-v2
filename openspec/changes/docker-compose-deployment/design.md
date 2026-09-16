## Context

Users need a git server to push to as a sync target. Docker Compose is the lowest-friction way to run a self-hosted server for users who aren't on Kubernetes. The stack only needs to serve git repositories over HTTP; TLS is delegated to a reverse proxy (nginx, Caddy, Traefik) that the user provides.

The original plan (below, kept for history) was to use a third-party image bundling nginx + git-http-backend, needing zero custom configuration. Verifying it end to end (a real `docker compose up`, then clone and push) found its `git push` path structurally broken: fcgiwrap never responded to `git-receive-pack`, reproducibly, across multiple attempts and configurations. A follow-up search for an actively maintained alternative found none; the few HTTP git server images that exist use the same nginx+fcgiwrap+CGI architecture and carry the same risk. This revision builds a small custom image instead, using Apache + `mod_cgid` rather than nginx + fcgiwrap.

## Goals / Non-Goals

**Goals:**
- A minimal `docker-compose.yml` that serves git repos over HTTP via `git-http-backend`, with working push
- A small, auditable custom image (`git-http/`) rather than an unmaintained third-party one
- A brief README with quickstart instructions

**Non-Goals:**
- TLS termination inside the compose stack
- Authentication beyond git protocol basics (fronting with nginx basic auth or an OAuth proxy is the user's concern)
- Web UI

## Decisions

**Build a small custom image: the official Debian-based `httpd:2.4` plus the `git` package, `mod_cgid` serving `git-http-backend`.**
This is git's own documented pattern for serving repositories over HTTP (`git help http-backend` describes exactly this Apache + `mod_cgid`/`mod_cgi` setup). Alpine's `git` package was tried first and rejected: it doesn't ship `git-http-backend` at all, only client-side HTTP tools (`git-http-fetch`, `git-http-push`, `git-remote-http(s)`). Debian's `git` package includes it at `/usr/lib/git-core/git-http-backend`, a long-standing, well-documented path.
Alternative: search harder for an existing public image. Rejected. The HTTP git server image ecosystem is thin and unmaintained, and shares the exact nginx+fcgiwrap architecture that caused the original bug; a from-scratch Apache image is both more reliable and small enough to own directly (three files: `Dockerfile`, `httpd.conf`, `entrypoint.sh`).

**`mod_cgid`, not `mod_cgi`.**
`httpd:2.4`'s default MPM is `event` (threaded). `mod_cgi` isn't safe to use with a threaded MPM; `mod_cgid` runs each CGI invocation through a separate daemon process instead and is the correct choice here.

**A custom, minimal `httpd.conf`, not the stock one with additions.**
The stock config carries directives unrelated to this single purpose (virtual hosts, status pages, autoindex, etc.). A from-scratch config listing exactly the modules and directives needed (`ScriptAlias`, `SetEnv GIT_PROJECT_ROOT`/`GIT_HTTP_EXPORT_ALL`, the `<Directory>` block granting `+ExecCGI`) is easier to audit than diffing against upstream's.

**Document `docker compose exec -u www-data` for creating bare repos, not plain `exec`.**
`docker compose exec` defaults to root. A repo directory created that way is root-owned, and git-http-backend (running as `www-data` via the Apache worker) later refuses it with a `safe.directory` ownership error, the same failure mode found and fixed in an earlier (SSH-based, since abandoned) attempt at this same problem. The fix is `-u www-data` on the `exec` call, not relaxing `safe.directory`, which is a real protection, not container-specific noise.

**Document `git config http.receivepack true` per repository.**
`git-http-backend` refuses `git-receive-pack` (push) by default regardless of the web server in front of it or `GIT_HTTP_EXPORT_ALL` (which only affects read access). This is a git-level safety default, not something to work around; each bare repo needs it set explicitly once, which the README's quickstart does.

**Port 80 exposed on the host.**
The expected deployment pattern is a reverse proxy on the same host mapping a domain to this port. Exposing 8080 would require an extra proxy config step for no benefit.

## Risks / Trade-offs

[No authentication by default] → The stack is intended for a trusted network or behind a proxy with auth. Documented clearly in the README. Users who need auth should front it with nginx basic auth or an OAuth proxy.

[Verified manually, not covered by CI] → This stack was verified with a real `docker compose up`, `git init --bare`, enabling `http.receivepack`, clone, push, and a restart to confirm repository persistence, using both the system `git` client and `git-syn`'s own `go-git`-based HTTP transport. None of this runs in CI (the `test` job doesn't have Docker available). A regression here would only surface on the next manual verification pass.

[Image freshness] → `httpd:2.4` and the `git` package version aren't pinned to exact versions, only to the `2.4` Apache line. Acceptable for a self-built image where a rebuild picks up patches; revisit if reproducibility becomes a concern.

---

## Superseded: original third-party HTTP image design

**Use `nginxinc/nginx-unprivileged` + `git-http-backend` via CGI (or a purpose-built image like `jkarlos/git-server-docker`).**
A purpose-built image bundles nginx + git-http-backend and requires zero custom configuration. Alternative: a bare `git daemon` image. Rejected as `git daemon` uses the git protocol (port 9418), not HTTP; users can't push to it over HTTPS through a standard proxy. Superseded: no purpose-built image was found that actually works; the custom Apache image replaces this entirely.

**Single service, single named volume.**
Keep complexity minimal. Multiple repos are served from subdirectories of the mounted volume. Alternative: one container per repo. Rejected as operationally heavy for a DR mirror. This decision carried over unchanged into the current design.

**[Image freshness]** → A specific third-party image version should be pinned in `docker-compose.yml` to avoid surprise breaking changes. Superseded: not applicable to a self-built image.
