## ADDED Requirements

### Requirement: pre-push subcommand entry point
The `git-syn pre-push` command SHALL act as a git pre-push hook delegate, reading remotes from `.gitremotes` and pushing to each one.

#### Scenario: Successful push to all remotes
- **WHEN** the user's git push triggers the pre-push hook
- **AND** `.gitremotes` contains one or more valid remotes
- **THEN** `git-syn pre-push` SHALL push to every remote listed in `.gitremotes`
- **THEN** the command SHALL print a per-remote result line: `✔ <name>` on success or `✘ <name>: <error>` on failure
- **THEN** the command SHALL exit 0

#### Scenario: Partial failure — at least one remote succeeds
- **WHEN** pushing to one or more remotes fails but at least one succeeds
- **THEN** the command SHALL print failure lines for each failed remote
- **THEN** the command SHALL exit 0

#### Scenario: All remotes fail
- **WHEN** pushing to every remote in `.gitremotes` fails
- **THEN** the command SHALL print a failure line for each remote
- **THEN** the command SHALL exit non-zero

### Requirement: .gitremotes file required
The `git-syn pre-push` command SHALL fail gracefully when `.gitremotes` does not exist.

#### Scenario: .gitremotes is absent
- **WHEN** `git-syn pre-push` is invoked and `.gitremotes` does not exist in the repository root
- **THEN** the command SHALL print an error directing the user to run `git-syn install`
- **THEN** the command SHALL exit non-zero

### Requirement: Concurrent remote pushes
The `git-syn pre-push` command SHALL push to all remotes concurrently to minimize total push time.

#### Scenario: Multiple remotes pushed in parallel
- **WHEN** `.gitremotes` contains two or more remotes
- **THEN** pushes to all remotes SHALL be initiated concurrently
- **THEN** the command SHALL wait for all pushes to complete before printing results and exiting

### Requirement: Accurate help text
The `git-syn pre-push` command SHALL have descriptive help text.

#### Scenario: User requests help
- **WHEN** the user runs `git-syn pre-push --help`
- **THEN** the output SHALL describe that the command pushes to all remotes in `.gitremotes`
- **THEN** the output SHALL note it is intended to be invoked by the git pre-push hook
