## ADDED Requirements

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

### Requirement: Install at specified path
The `git-syn install` command SHALL accept a `--path` flag to initialize a repository at a path other than the current working directory.

#### Scenario: Successful installation at explicit path
- **WHEN** the user runs `git-syn install --path /some/repo` and `/some/repo` is a valid git repository
- **THEN** the `pre-push` hook SHALL be written to `/some/repo/.git/hooks/pre-push`
- **THEN** the hook file SHALL have permissions set to `0700`
- **THEN** `.gitremotes` SHALL be written to `/some/repo/.gitremotes`
- **THEN** each remote from `.gitremotes` SHALL be appended to `/some/repo/.git/config` (without duplicating existing entries)
- **THEN** the command SHALL print `Updated git hooks. Git SYN initialized.`

#### Scenario: Explicit path is not a git repository
- **WHEN** the user runs `git-syn install --path /not/a/repo` and the target is not a git repository
- **THEN** the command SHALL exit with a non-zero status and print a clear error message

### Requirement: Pre-push hook content
The `pre-push` hook SHALL be embedded in the binary as a Go string constant using a `{{Command}}` placeholder, following the Git LFS pattern.

#### Scenario: git-syn is available in PATH
- **WHEN** the `pre-push` hook is invoked by git
- **AND** `git-syn` is found in PATH
- **THEN** the hook SHALL delegate to `git syn pre-push "$@"`

#### Scenario: git-syn is not available in PATH
- **WHEN** the `pre-push` hook is invoked by git
- **AND** `git-syn` is NOT found in PATH
- **THEN** the hook SHALL print a warning directing the user to remove the hook if they no longer need Git SYN
- **THEN** the hook SHALL exit with code 2

### Requirement: Remote URL type support
The `git-syn install` command SHALL write `.gitremotes` entries only for supported remote URL types.

#### Scenario: HTTPS remote
- **WHEN** a configured remote has an HTTPS URL (e.g., `https://gitlab.com/user/repo.git`)
- **THEN** the URL SHALL be written to `.gitremotes` as-is in gitconfig format

#### Scenario: OpenSSH remote
- **WHEN** a configured remote has an OpenSSH URL (e.g., `git@gitlab.com:user/repo.git`)
- **THEN** the URL SHALL be written to `.gitremotes` as-is in gitconfig format

#### Scenario: Unsupported URL scheme
- **WHEN** a configured remote has a URL that is neither HTTPS nor OpenSSH
- **THEN** the command SHALL skip that remote and print a warning identifying the remote name and unsupported URL

### Requirement: Append .gitremotes entries to .git/config
After writing `.gitremotes`, the `git-syn install` command SHALL read the file back and append each remote entry to `.git/config` so that git is aware of them for synchronization.

#### Scenario: Remotes appended to .git/config
- **WHEN** `.gitremotes` has been written with one or more remote entries
- **THEN** each remote's `[remote "name"]` section and `url` SHALL be appended to `.git/config`
- **THEN** remotes already present in `.git/config` SHALL NOT be duplicated

#### Scenario: .gitremotes is empty
- **WHEN** `.gitremotes` contains no remote entries
- **THEN** `.git/config` SHALL not be modified
- **THEN** the command SHALL still complete successfully

### Requirement: Accurate help text
The `git-syn install` command SHALL have descriptive help text consistent with the README.

#### Scenario: User requests help
- **WHEN** the user runs `git-syn install --help`
- **THEN** the output SHALL describe that the command installs Git SYN into a repository for remote synchronization
- **THEN** the output SHALL document the `--path` flag
