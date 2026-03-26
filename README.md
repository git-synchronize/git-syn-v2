# Git Synchronize v2

_Remote git repository synchronization._

Git SYN is a command line extension for synchronizing git remote repositories.

## Dependencies

- [check](https://libcheck.github.io/check) (optional)
- [libgit2](https://libgit2.org)
- [pandoc](https://pandoc.org) (optional)

## Compiling

```sh
make
```

## Installation

```sh
make install
```

## Usage

```sh
git syn -h
Usage: git-syn [option] ... [command] ...

Remote git repository synchronization.

  -h, --help       display this help and exit
  -v, --version    output version information and exit
  --verbose        enable verbose output
  --debug          enable debug output (implies verbose)
  config           manage git-syn configuration
  daemon           run synchronization on a schedule
  install          install extension to repository
  pre-push         push to all remotes in .gitremotes
  remote           manage tracked remotes
  uninstall        remove extension from repository

Git SYN online help: <https://gitlab.com/git-syn/git-syn>
Full documentation <https://git-syn.gitlab.io/git-syn>
or available locally via: man git-syn
```

### Daemon Mode

To run git-syn as a background process (daemon) with a 5-minute sync interval:

```sh
git syn daemon --interval 5m
```

The interval can also be configured in `~/.config/git-syn/config.yaml`.

## Testing

To run the unit tests:

```sh
go test ./...
```

Or using `mise`:

```sh
mise run test
```

To run all checks (formatting, linting, and tests):

```sh
mise run check
```

## License

SPDX-License-Identifier: [GPL-2.0-or-later](COPYING)

## Reference

- [How to integrate new subcommands](https://git.kernel.org/pub/scm/git/git.git/plain/Documentation/howto/new-command.txt)
- [/srv : Data for services provided by this system](https://refspecs.linuxfoundation.org/FHS_3.0/fhs/ch03s17.html)

## See Also

- [Git Large File Storage (LFS)](https://git-lfs.github.com)

- [GitLab Repository Mirroring](https://docs.gitlab.com/ee/user/project/repository/repository_mirroring.html)

- [git-sync](https://github.com/kubernetes/git-sync)
