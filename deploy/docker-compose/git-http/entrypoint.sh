#!/bin/sh
set -eu

mkdir -p /var/lib/git
chown nginx:nginx /var/lib/git

# git-syn serve (and the git-http-backend processes it spawns) run as the
# same unprivileged user as nginx's workers, so repositories they create
# and repositories nginx serves have consistent ownership.
gosu nginx git-syn serve --path /var/lib/git --listen 127.0.0.1:8080 &

exec nginx -g "daemon off;"
