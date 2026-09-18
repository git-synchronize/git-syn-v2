% git-syn(1)
% Lucas Ramage
% September 18, 2026

# NAME

git-syn - remote git repository synchronization

# SYNOPSIS

**git-syn** [option] ... [command] ...

# DESCRIPTION

Git SYN is a command line extension for synchronizing remote git repositories. It enables remote repository synchronization across multiple git forges for disaster recovery and censorship resistance.

# OPTIONS

**-h, --help**

:   Prints the command usage instructions.

**-v, --version**

:   Prints the version number.

**--verbose**

:   Enable verbose output.

**--debug**

:   Enable debug output (implies verbose).

# COMMANDS

**config**

:   Manage git-syn configuration.
    Use `git syn config init` to create a default configuration file.

**daemon**

:   Run synchronization on a schedule.
    Pushes to all remotes in .gitremotes at a fixed interval until interrupted.
    Use `--interval` (e.g. `5m`, `1h`) to set the sync interval, overriding
    the `sync_interval` value in the configuration file. Use `--once` to
    run a single synchronization pass and exit instead of looping.

**install**

:   Install extension to repository.
    Installs a pre-push hook into .git/hooks/ and writes a .gitremotes file.

**pre-push**

:   Push to all remotes in .gitremotes.
    Typically invoked by the git pre-push hook.

**remote**

:   Manage tracked remotes.
    Subcommands: `add`, `remove`, `list`.

**serve**

:   Host git repositories over HTTP for clone and push.
    Repositories must already exist under `--path` (create them with
    `git init --bare`) and, to accept pushes, have `http.receivepack` set
    to `true`. Use `--listen` (default `:8080`) to set the bind address.
    There is no authentication; front this with a reverse proxy for TLS
    or access control.

**uninstall**

:   Remove extension from repository.
    Removes the pre-push hook.

# EXIT STATUS

Returns zero on success, non-zero on failure.

# NOTES

The git-syn project site, with more information and the source code
repository, can be found at <https://gitlab.com/git-syn/git-syn-v2>.

This tool is currently under development, please report any bugs at
the project site or directly to the author.
