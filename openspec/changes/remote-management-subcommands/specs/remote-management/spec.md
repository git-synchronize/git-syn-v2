## ADDED Requirements

### Requirement: remote add subcommand
The `git-syn remote add <name> <url>` command SHALL append a new remote to `.gitremotes` and register it in `.git/config`.

#### Scenario: Successful add
- **WHEN** the user runs `git-syn remote add backup https://gitlab.com/user/repo.git`
- **AND** the remote name does not already exist in `.gitremotes`
- **THEN** the remote SHALL be appended to `.gitremotes` in gitconfig format
- **THEN** the remote SHALL be registered in `.git/config`
- **THEN** the command SHALL print `Added remote 'backup'.`

#### Scenario: Duplicate name is a no-op
- **WHEN** the user runs `git-syn remote add backup <url>`
- **AND** a remote named `backup` already exists in `.gitremotes`
- **THEN** neither `.gitremotes` nor `.git/config` SHALL be modified
- **THEN** the command SHALL print `Remote 'backup' already exists. No changes made.`
- **THEN** the command SHALL exit 0

#### Scenario: Unsupported URL scheme rejected
- **WHEN** the user runs `git-syn remote add mirror file:///local/path`
- **THEN** the command SHALL print an error indicating the URL scheme is not supported
- **THEN** `.gitremotes` and `.git/config` SHALL not be modified
- **THEN** the command SHALL exit non-zero

### Requirement: remote remove subcommand
The `git-syn remote remove <name>` command SHALL remove a remote from `.gitremotes` and deregister it from `.git/config`.

#### Scenario: Successful remove
- **WHEN** the user runs `git-syn remote remove backup`
- **AND** `backup` exists in `.gitremotes`
- **THEN** the `backup` stanza SHALL be removed from `.gitremotes`
- **THEN** the `backup` entry SHALL be removed from `.git/config`
- **THEN** the command SHALL print `Removed remote 'backup'.`

#### Scenario: Remote not in .gitremotes
- **WHEN** the user runs `git-syn remote remove unknown`
- **AND** `unknown` does not exist in `.gitremotes`
- **THEN** the command SHALL print an error: `Remote 'unknown' not found in .gitremotes`
- **THEN** `.git/config` SHALL not be modified
- **THEN** the command SHALL exit non-zero

### Requirement: remote list subcommand
The `git-syn remote list` command SHALL display all remotes tracked in `.gitremotes`.

#### Scenario: Remotes present
- **WHEN** `.gitremotes` contains one or more remotes
- **THEN** the command SHALL print each remote name and URL, one per line
- **THEN** any remote present in `.gitremotes` but missing from `.git/config` SHALL be flagged with `[unregistered]`

#### Scenario: .gitremotes is empty or absent
- **WHEN** `.gitremotes` does not exist or contains no remotes
- **THEN** the command SHALL print `No remotes tracked by git-syn.`
- **THEN** the command SHALL exit 0

### Requirement: --path flag support
All `git-syn remote` subcommands SHALL accept a `--path <dir>` flag defaulting to the current directory, consistent with `install` and `uninstall`.

#### Scenario: Explicit path used
- **WHEN** the user runs `git-syn remote list --path /some/repo`
- **THEN** the command SHALL read `.gitremotes` from `/some/repo/.gitremotes`
- **THEN** the command SHALL read `.git/config` from `/some/repo/.git/config`
