## ADDED Requirements

### Requirement: pre-push subcommand entry point
The `git-syn pre-push` command SHALL act as a git pre-push hook delegate, reading remotes from `.gitremotes` and mirroring the exact refs being pushed to each other remote.

#### Scenario: Successful push to all mirror remotes
- **WHEN** the user's git push triggers the pre-push hook
- **AND** `.gitremotes` contains one or more valid remotes
- **THEN** `git-syn pre-push` SHALL push the same refs to every remote in `.gitremotes` except the one git is already pushing to
- **THEN** the command SHALL print a per-remote result line: `✔ <name>` on success or `✘ <name>: <error>` on failure
- **THEN** the command SHALL exit 0

#### Scenario: Partial failure — at least one remote succeeds
- **WHEN** pushing to one or more remotes fails but at least one succeeds
- **THEN** the command SHALL print failure lines for each failed remote
- **THEN** the command SHALL exit 0

#### Scenario: All remotes fail
- **WHEN** pushing to every mirror remote in `.gitremotes` fails
- **THEN** the command SHALL print a failure line for each remote
- **THEN** the command SHALL exit non-zero

### Requirement: Refspec forwarding from hook stdin
The `git-syn pre-push` command SHALL read the refs being pushed from stdin and mirror exactly those refs to each mirror remote.

#### Scenario: Single branch push
- **WHEN** git invokes the hook for `git push gitlab feat/test`
- **AND** stdin contains `refs/heads/feat/test <local-sha> refs/heads/feat/test <remote-sha>`
- **THEN** the command SHALL push `refs/heads/feat/test` to each mirror remote
- **THEN** the command SHALL NOT push any other refs

#### Scenario: New branch (remote sha is all zeros)
- **WHEN** stdin contains a zero remote sha (`0000000...`) indicating a new branch
- **THEN** the command SHALL push the local ref to create the branch on each mirror remote

#### Scenario: Branch deletion (local sha is all zeros)
- **WHEN** stdin contains a zero local sha (`0000000...`) indicating a deletion
- **THEN** the command SHALL delete the corresponding ref on each mirror remote

#### Scenario: Multiple refs pushed at once
- **WHEN** stdin contains multiple ref lines (e.g. `git push --all`)
- **THEN** the command SHALL push all listed refs to each mirror remote in a single push operation per remote

### Requirement: Skip the primary remote
The `git-syn pre-push` command SHALL not push to the remote that git is already pushing to, as that push is handled by git itself.

#### Scenario: Primary remote is in .gitremotes
- **WHEN** git invokes the hook with remote name `gitlab` as the first positional argument
- **AND** `.gitremotes` contains an entry named `gitlab`
- **THEN** the command SHALL skip `gitlab` and push only to the other remotes
- **THEN** the command SHALL NOT report a result line for the skipped remote

### Requirement: Mirror pushes use force
Mirror remotes are followers, not arbiters of history. The `git-syn pre-push` command SHALL use force when pushing to mirror remotes to ensure they stay in sync even after rebases or amends on the primary.

#### Scenario: Non-fast-forward push
- **WHEN** the local branch has been rebased or amended since the last sync
- **THEN** the push to mirror remotes SHALL succeed using force semantics
- **THEN** the push SHALL NOT be blocked by non-fast-forward rejection from the mirror

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
