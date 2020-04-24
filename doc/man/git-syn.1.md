% git-syn(1)
% Lucas Ramage
% April 24, 2020

# NAME

git-syn - event-driven git remote repository syncing

# SYNOPSIS

**git-syn** [option] ... [command] ...

# DESCRIPTION

Git SYN is a command line extension for synchronizing Git remote repositories.

# OPTIONS

**-h, --help**

:   Prints the command usage instructions.

**-v, --version**

:   Prints the version number.

# COMMANDS

Currently, Git SYN only contains high level ("porcelain") commands.

**git-syn-install(1)**

:   Install Git SYN configuration.

**git-syn-uninstall(1)**

:   Remove Git SYN configuration.

**git-syn-monitor(1)**

:   View or enable/disable Git SYN repository syncing.

# EXIT STATUS

Returns zero on success, errno values on failure.

# NOTES

The git-syn project site, with more information and the source code
repository, can be found at <https://gitlab.com/oxr463/git-syn>.

This tool is currently under development, please report any bugs at
the project site or directly to the author.

