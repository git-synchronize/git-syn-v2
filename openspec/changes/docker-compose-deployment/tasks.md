## 1. Directory structure

- [x] 1.1 Create `deploy/docker-compose/` directory

## 2. Custom Apache + git-http-backend image

- [x] 2.1 Write `deploy/docker-compose/git-http/Dockerfile`: Debian-based `httpd:2.4` plus the `git` package (Alpine's `git` package doesn't ship `git-http-backend`, only client-side HTTP tools)
- [x] 2.2 Write `deploy/docker-compose/git-http/httpd.conf`: a minimal config loading `mod_cgid` (not `mod_cgi`, unsafe with the default threaded `event` MPM), `ScriptAlias`-ing to `/usr/lib/git-core/git-http-backend/`, and setting `GIT_PROJECT_ROOT`/`GIT_HTTP_EXPORT_ALL`
- [x] 2.3 Write `deploy/docker-compose/git-http/entrypoint.sh`: ensures `/var/lib/git` exists and is owned by `www-data`, then execs `httpd-foreground`

## 3. docker-compose.yml

- [x] 3.1 Write `deploy/docker-compose/docker-compose.yml` with a `git-serve` service that builds `git-http/` locally
- [x] 3.2 Expose port 80 on the host
- [x] 3.3 Mount a named volume for repository storage (`/var/lib/git`)

## 4. Documentation

- [x] 4.1 Write `deploy/docker-compose/README.md` with a quickstart: build and start the stack, create a bare repo as `www-data`, enable `http.receivepack`, add it as a git-syn remote
- [x] 4.2 Document that TLS is not handled by this stack and recommend fronting with a reverse proxy

## 5. Smoke test

- [x] 5.1 Run `docker compose up --build` locally and verify `git clone` and `git push` work against the running server. Verified: clone, push (both the system `git` client and `git-syn`'s own `go-git`-based HTTP transport), and a restart to confirm repository data persists across it. Two prior approaches were tried and abandoned before this one: the original third-party `cirocosta/gitserver-http` image (nginx + fcgiwrap), whose push path was structurally broken, and a custom SSH-based image, built when a miscommunication led to switching protocols entirely before the actual requirement (stay on HTTP, just fix the image) was confirmed. See `design.md`'s "Superseded" section and `proposal.md`'s revision note for that history.
