#!/bin/sh
command -v git-syn >/dev/null 2>&1 || { echo >&2 "\nThis repository is configured for Git SYN but 'git-syn' was not found on your path. If you no longer wish to use Git SYN, remove this hook by deleting .git/hooks/pre-push.\n"; exit 2; }
git syn pre-push "$@"
