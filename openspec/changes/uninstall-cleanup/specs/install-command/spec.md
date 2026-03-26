## MODIFIED Requirements

### Requirement: Install command entry point
The `git-syn install` command SHALL initialize Git SYN in the current repository by writing a `pre-push` git hook into `.git/hooks/` and writing a `.gitremotes` file.

#### Scenario: Successful installation
- **WHEN** the user runs `git-syn install` inside a git repository
- **THEN** the `pre-push` hook SHALL be written to `.git/hooks/pre-push`
- **THEN** the hook file SHALL have permissions set to `0700`
- **THEN** a `.gitremotes` file SHALL be created in the current directory with each remote's name and URL in gitconfig format
- **THEN** each remote from `.gitremotes` SHALL be appended to `.git/config` (without duplicating existing entries)
- **THEN** the command SHALL print `Updated git hooks. Git SYN initialized.`

#### Scenario: Directory is not a git repository
- **WHEN** the user runs `git-syn install` in a directory that is not a git repository
- **THEN** the command SHALL exit with a non-zero status and print a clear error message identifying the directory as not a git repository

### Requirement: Uninstall removes hook and warns about artifacts
The `git-syn uninstall` command SHALL remove the pre-push hook and warn the user about remaining artifacts (`.gitremotes` and `.git/config` entries).

#### Scenario: Successful uninstall without --clean
- **WHEN** the user runs `git-syn uninstall` (no flags)
- **THEN** `.git/hooks/pre-push` SHALL be removed
- **THEN** `.gitremotes` SHALL remain on disk
- **THEN** `.git/config` entries added by git-syn SHALL remain
- **THEN** the command SHALL print `Removed git hooks. Git SYN uninstalled. Run with --clean to also remove .gitremotes and remote config entries.`

#### Scenario: Successful uninstall with --clean
- **WHEN** the user runs `git-syn uninstall --clean`
- **THEN** `.git/hooks/pre-push` SHALL be removed
- **THEN** `.gitremotes` SHALL be deleted
- **THEN** all remotes listed in `.gitremotes` SHALL be removed from `.git/config`
- **THEN** the command SHALL print `Removed git hooks and remotes. Git SYN fully uninstalled.`

#### Scenario: --clean with missing .gitremotes
- **WHEN** the user runs `git-syn uninstall --clean`
- **AND** `.gitremotes` does not exist
- **THEN** the command SHALL proceed silently without error
- **THEN** the hook SHALL still be removed if present
- **THEN** the command SHALL print `Removed git hooks and remotes. Git SYN fully uninstalled.`

#### Scenario: Hook already absent
- **WHEN** the user runs `git-syn uninstall`
- **AND** `.git/hooks/pre-push` does not exist
- **THEN** the command SHALL succeed silently (idempotent)
