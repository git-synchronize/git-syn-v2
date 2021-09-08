# Git Synchronize

_Remote git repository syncing._

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

Remote git repository syncing.

  -h, --help       display this help and exit
  -v, --version    output version information and exit
  install          install extension to repository
  uninstall        remove extension from repository
  monitor          view or enable/disable repository sync

Git SYN online help: <https://gitlab.com/git-syn/git-syn>
Full documentation <https://git-syn.gitlab.io/git-syn>
or available locally via: man git-syn
```

## License

SPDX-License-Identifier: [GPL-2.0-or-later](COPYING)

## Reference

- [How to integrate new subcommands](https://git.kernel.org/pub/scm/git/git.git/plain/Documentation/howto/new-command.txt)
- [/srv : Data for services provided by this system](https://refspecs.linuxfoundation.org/FHS_3.0/fhs/ch03s17.html)

## See Also

- [Git Large File Storage (LFS)](https://git-lfs.github.com)

- [GitLab Repository Mirroring](https://docs.gitlab.com/ee/user/project/repository/repository_mirroring.html)

