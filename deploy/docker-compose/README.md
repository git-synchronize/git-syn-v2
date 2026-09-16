# git-syn Docker Compose Deployment

Runs a self-hosted git HTTP server using `git-http-backend` for use as a git-syn mirror target.

## Quickstart

```sh
# 1. Clone and enter this directory
cd deploy/docker-compose

# 2. Start the server
docker compose up -d

# 3. Create a bare repository on the server
docker compose exec git-serve git init --bare /var/lib/git/myrepo.git

# 4. Add the server as a git-syn remote
git-syn remote add mirror http://localhost/myrepo.git
```

## Notes

**Repository storage**: the `cirocosta/gitserver-http` image hardcodes its
repository root to `/var/lib/git` inside the container and always exports
every repository — it does not read `GIT_PROJECT_ROOT` or
`GIT_HTTP_EXPORT_ALL` from the environment despite documenting them. The
compose file's named volume is mounted at `/var/lib/git` to match; there is
nothing to configure here.

**TLS**: This stack serves plain HTTP on port 80. TLS is not handled here.
Front the service with a reverse proxy (nginx, Caddy, Traefik) to add HTTPS.
Example nginx snippet:

```nginx
server {
    listen 443 ssl;
    server_name git.example.com;
    location / {
        proxy_pass http://localhost:80;
    }
}
```

**Authentication**: No authentication is configured by default. The stack is
intended for use on a trusted network or behind a proxy that enforces auth.
For basic auth, add an `auth_basic` block to your nginx configuration.
