% git-syn(1)
% Lucas Ramage
% September 7, 2021

# NAME

git-syn - remote git repository syncing

# SYNOPSIS

**git-syn** [option] ... [command] ...

# DESCRIPTION

Git SYN is a command line extension for synchronizing remote git repositories.

# OPTIONS

**-h, --help**

:   Prints the command usage instructions.

**-v, --version**

:   Prints the version number.

# COMMANDS

Like Git, Git SYN commands are separated into high level ("porcelain") commands
and low level ("plumbing") commands.

## High level commands (porcelain)

**git-syn-install(1)**

:   Install Git SYN configuration.

**git-syn-uninstall(1)**

:   Remove Git SYN configuration.

**git-syn-monitor(1)**

:   View or enable/disable Git SYN repository syncing.

## Low level commands (plumbing)

**git-syn-pre-push(1)**

:   Git pre-push hook implementation.

# EXIT STATUS

Returns zero on success, errno values on failure.

# NOTES

The git-syn project site, with more information and the source code
repository, can be found at <https://gitlab.com/git-syn/git-syn>.

This tool is currently under development, please report any bugs at
the project site or directly to the author.

