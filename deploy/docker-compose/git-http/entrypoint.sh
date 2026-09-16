#!/bin/sh
set -eu

mkdir -p /var/lib/git
chown nginx:nginx /var/lib/git

# cgi-bridge (and the git-http-backend processes it spawns) run as the same
# unprivileged user as nginx's workers, so repositories they create and
# repositories nginx serves have consistent ownership.
GIT_PROJECT_ROOT=/var/lib/git gosu nginx /usr/local/bin/cgi-bridge &

exec nginx -g "daemon off;"
