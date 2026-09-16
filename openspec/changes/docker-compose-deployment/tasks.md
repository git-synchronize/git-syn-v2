## 1. Directory structure

- [x] 1.1 Create `deploy/docker-compose/` directory

## 2. Custom nginx + git-http-backend image

- [x] 2.1 Write `deploy/docker-compose/git-http/cgi-bridge/main.go` (and its own `go.mod`): a standalone Go program using `net/http/cgi` to run `git-http-backend` as a CGI process and expose it over plain HTTP
- [x] 2.2 Write `deploy/docker-compose/git-http/Dockerfile`: multi-stage build, `golang:alpine` to build `cgi-bridge`, Debian-based `nginx:1.27` (not `nginx:alpine`; Alpine's `git` package doesn't ship `git-http-backend`) plus `git` and `gosu` for the runtime
- [x] 2.3 Write `deploy/docker-compose/git-http/nginx.conf`: a minimal config reverse-proxying to `cgi-bridge` over HTTP, with default (not disabled) request buffering
- [x] 2.4 Write `deploy/docker-compose/git-http/entrypoint.sh`: ensures `/var/lib/git` exists and is owned by `nginx`, starts `cgi-bridge` as that same user via `gosu`, then execs nginx in the foreground

## 3. docker-compose.yml

- [x] 3.1 Write `deploy/docker-compose/docker-compose.yml` with a `git-serve` service that builds `git-http/` locally
- [x] 3.2 Expose port 80 on the host
- [x] 3.3 Mount a named volume for repository storage (`/var/lib/git`)

## 4. Documentation

- [x] 4.1 Write `deploy/docker-compose/README.md` with a quickstart: build and start the stack, create a bare repo as the `nginx` user, enable `http.receivepack`, add it as a git-syn remote
- [x] 4.2 Document that TLS is not handled by this stack and recommend fronting with a reverse proxy

## 5. Smoke test

- [x] 5.1 Run `docker compose up --build` locally and verify `git clone` and `git push` work against the running server, including at realistic push sizes. Verified: clone, push at several sizes (a small commit, 20MB, 100MB, both the system `git` client and `git-syn`'s own `go-git`-based HTTP transport), and a restart to confirm repository data persists across it. An initial version with `proxy_request_buffering off` failed pushes above roughly 20MB with an HTTP 400; removing that directive (nginx's default buffering) fixed it. Two earlier approaches were tried and abandoned before this one: the original third-party `cirocosta/gitserver-http` image (nginx + fcgiwrap), whose push path was structurally broken, and a custom Apache + `mod_cgid` image, which worked but switched away from nginx before the actual requirement (keep nginx, fix the image) was confirmed. See `design.md`'s "Superseded" sections and `proposal.md`'s revision note for that history.
