#!/bin/sh
set -eu

mkdir -p /var/lib/git
chown www-data:www-data /var/lib/git

exec httpd-foreground
