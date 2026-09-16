## Context

Users need a git server to push to as a sync target. Docker Compose is the lowest-friction way to run a self-hosted server for users who aren't on Kubernetes. The stack only needs to serve git repositories over HTTP; TLS is delegated to a reverse proxy (nginx, Caddy, Traefik) that the user provides.

The original plan (kept below under "Superseded") used a third-party image bundling nginx + git-http-backend, needing zero custom configuration. Verifying it end to end found its `git push` path structurally broken: fcgiwrap never responded to `git-receive-pack`, reproducibly, across multiple attempts and configurations. No actively maintained alternative was found either. A first replacement (also kept below) built a custom Apache + `mod_cgid` image, which worked, but switched away from nginx when the actual requirement was to keep it and just fix the broken image. This revision bridges nginx to `git-http-backend` with a small Go program instead of relying on any FastCGI-family gateway.

## Goals / Non-Goals

**Goals:**
- A minimal `docker-compose.yml` that serves git repos over HTTP via nginx and `git-http-backend`, with working push at realistic sizes
- A small, auditable custom image (`git-http/`) rather than an unmaintained third-party one
- A brief README with quickstart instructions

**Non-Goals:**
- TLS termination inside the compose stack
- Authentication beyond git protocol basics (fronting with nginx basic auth or an OAuth proxy is the user's concern)
- Web UI

## Decisions

**Bridge nginx to `git-http-backend` with a small Go program using `net/http/cgi`, not fcgiwrap or uwsgi's cgi plugin.**
nginx has no native CGI support; it needs some bridge to a CGI script. fcgiwrap is the usual choice and is what broke `git-receive-pack` in the original design. uwsgi's cgi plugin is the other common option, but it's a comparatively obscure, thinly documented configuration surface, and repeating a FastCGI-family bridge risked hitting the same class of bug in a different disguise. `net/http/cgi` is part of Go's standard library, and git-syn is a Go project already comfortable owning this: nginx talks ordinary HTTP/1.1 to `cgi-bridge` (`proxy_pass`), and `cgi-bridge` uses `net/http/cgi` to run `git-http-backend` as a ordinary CGI child process, connecting the HTTP request/response directly to its stdin/stdout. There's no FastCGI record framing anywhere in the chain.
Alternative: nginx + uwsgi's cgi plugin. Rejected as a second attempt at the same architecture class (a FastCGI-derived protocol bridge) that failed once already, for comparatively little benefit over owning ~20 lines of Go directly.

**Use nginx's default request buffering (don't set `proxy_request_buffering off`).**
An initial version disabled request buffering, reasoning that a streamed push shouldn't be held in memory or on disk first. This reproducibly broke pushes above roughly 20MB with an HTTP 400 (confirmed: a 20MB push failed, the same push succeeded immediately after removing the directive, and 100MB pushes were then verified working too). The exact interaction wasn't chased further since default buffering solved it outright, is nginx's well-tested normal mode of operation for reverse-proxying arbitrary HTTP backends, and buffering a push's pack data to a temp file is an entirely ordinary cost for a DR-mirror use case.

**Debian-based `nginx:1.27`, not `nginx:alpine`.**
Same reasoning as the Apache attempt: Alpine's `git` package doesn't ship `git-http-backend` at all, only client-side HTTP tools (`git-http-fetch`, `git-http-push`, `git-remote-http(s)`). Debian's includes it at `/usr/lib/git-core/git-http-backend`.

**`cgi-bridge` is its own Go module, built in a separate Docker build stage.**
It's a small, standalone program with no relationship to git-syn's own module graph; giving it its own `go.mod` under `git-http/cgi-bridge/` keeps it fully independent of the root module (no shared dependencies, no risk of the deploy tooling affecting `go build ./...` for git-syn itself) while still living in the deploy directory it belongs to.

**`cgi-bridge` and nginx's workers run as the same `nginx` user (via `gosu` in `entrypoint.sh`), not root.**
Keeps repository ownership consistent between the process that creates a bare repo (documented as `docker compose exec -u nginx ...`) and the processes that later read and write it.

**Document `docker compose exec -u nginx` for creating bare repos, not plain `exec`.**
`docker compose exec` defaults to root. A repo directory created that way is root-owned, and `git-http-backend` (running as `nginx`) later refuses it with a `safe.directory` ownership error, the same failure mode found in the Apache attempt (there as `www-data`) and, before that, in a since-abandoned SSH-based attempt. The fix is `-u nginx` on the `exec` call, not relaxing `safe.directory`, which is a real protection, not container-specific noise.

**Document `git config http.receivepack true` per repository.**
`git-http-backend` refuses `git-receive-pack` (push) by default regardless of the web server in front of it or `GIT_HTTP_EXPORT_ALL` (which only affects read access). This is a git-level safety default, not something to work around; each bare repo needs it set explicitly once, which the README's quickstart does.

**Port 80 exposed on the host.**
The expected deployment pattern is a reverse proxy on the same host mapping a domain to this port. Exposing 8080 would require an extra proxy config step for no benefit.

## Risks / Trade-offs

[No authentication by default] → The stack is intended for a trusted network or behind a proxy with auth. Documented clearly in the README. Users who need auth should front it with nginx basic auth or an OAuth proxy.

[Verified manually, not covered by CI] → This stack was verified with a real `docker compose up`, `git init --bare`, enabling `http.receivepack`, clone, push at several sizes (a small commit, 20MB, 100MB), and a restart to confirm repository persistence, using both the system `git` client and `git-syn`'s own `go-git`-based HTTP transport. None of this runs in CI (the `test` job doesn't have Docker available). A regression here would only surface on the next manual verification pass.

[A custom bridge program is more to maintain than a config file] → Accepted given the two prior FastCGI-family bridges both proved unreliable (one broken outright, one requiring an undiagnosed-but-real config fix). ~20 lines of Go using only the standard library is a small, fully-owned surface in exchange for not depending on either.

[Image freshness] → `nginx:1.27` and the `git` package version aren't pinned to exact versions, only to the `1.27` nginx line. Acceptable for a self-built image where a rebuild picks up patches; revisit if reproducibility becomes a concern.

---

## Superseded: Apache + mod_cgid design

**Build a small custom image: the official Debian-based `httpd:2.4` plus the `git` package, `mod_cgid` serving `git-http-backend`.**
This is git's own documented pattern for serving repositories over HTTP (`git help http-backend` describes exactly this Apache + `mod_cgid`/`mod_cgi` setup), and it worked, verified end to end including push. Superseded: the actual requirement was to keep nginx as the web server, which this design didn't.

**`mod_cgid`, not `mod_cgi`.**
`httpd:2.4`'s default MPM is `event` (threaded); `mod_cgi` isn't safe to use with a threaded MPM. This concern doesn't carry over to the nginx design, which has no Apache MPM at all.

**A custom, minimal `httpd.conf`, not the stock one with additions.**
Not applicable to nginx.

## Superseded: original third-party HTTP image design

**Use `nginxinc/nginx-unprivileged` + `git-http-backend` via CGI (or a purpose-built image like `jkarlos/git-server-docker`).**
A purpose-built image bundles nginx + git-http-backend and requires zero custom configuration. Alternative: a bare `git daemon` image. Rejected as `git daemon` uses the git protocol (port 9418), not HTTP; users can't push to it over HTTPS through a standard proxy. Superseded: no purpose-built image was found that actually works; the custom `cgi-bridge` image replaces this entirely.

**Single service, single named volume.**
Keep complexity minimal. Multiple repos are served from subdirectories of the mounted volume. Alternative: one container per repo. Rejected as operationally heavy for a DR mirror. This decision carried over unchanged into the current design.

**[Image freshness]** → A specific third-party image version should be pinned in `docker-compose.yml` to avoid surprise breaking changes. Superseded: not applicable to a self-built image.
