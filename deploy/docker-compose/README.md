# git-syn Docker Compose Deployment

Runs a self-hosted git HTTP server for use as a git-syn mirror target. The
server is a small custom image (`git-http/`, built locally by this compose
file): nginx in front of `git-syn serve`, which hosts `git-http-backend`
directly and exposes it over plain HTTP for nginx to reverse-proxy to.

nginx has no native CGI support, so it needs a bridge to talk to
`git-http-backend`. The usual bridge is `fcgiwrap`, which is what an earlier
version of this stack used (via the third-party `cirocosta/gitserver-http`
image) and what made `git push` hang and fail with an nginx 504 in every
test (fcgiwrap never responded to `git-receive-pack`). `git-syn serve`
sidesteps the FastCGI protocol entirely: nginx speaks ordinary HTTP/1.1 to
it, and it uses Go's standard `net/http/cgi` package to talk to
`git-http-backend`.

## Quickstart

```sh
# 1. Clone and enter this directory
cd deploy/docker-compose

# 2. Build and start the server
docker compose up -d --build

# 3. Create a bare repository on the server, as the nginx user (not root,
# which "docker compose exec" defaults to -- a root-owned repo directory
# fails git's safe.directory check when git-http-backend later serves it
# as nginx)
docker compose exec -u nginx git-serve git init --bare /var/lib/git/myrepo.git

# 4. Allow pushes to it. git-http-backend refuses git-receive-pack by
# default regardless of the web server in front of it.
docker compose exec -u nginx git-serve git -C /var/lib/git/myrepo.git config http.receivepack true

# 5. Add the server as a git-syn remote
git-syn remote add mirror http://localhost/myrepo.git
```

## Notes

**Repository storage**: repos live under `/var/lib/git` inside the
container, backed by the `git-data` named volume, so they persist across
restarts. Create each bare repo once (steps 3 and 4 above) before pushing
to it.

**Access control**: there's no authentication in front of nginx by default.
Every repository under `/var/lib/git` is readable by anyone who can reach
port 80, and pushable by anyone once `http.receivepack` is set. This stack
is intended for a trusted network or behind a proxy that adds auth (see
below).

**TLS**: This stack serves plain HTTP on port 80. TLS is not handled here.
Front the service with a reverse proxy (a separate nginx, Caddy, Traefik) to
add HTTPS. Example nginx snippet:

```nginx
server {
    listen 443 ssl;
    server_name git.example.com;
    location / {
        proxy_pass http://localhost:80;
    }
}
```

For basic auth, add an `auth_basic` block directly to the `location /` block
in `git-http/nginx.conf`.
