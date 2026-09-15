## ADDED Requirements

### Requirement: init_repo and remote subcommand logic return errors instead of exiting
The system SHALL implement `init_repo` and the business logic behind the `remote add`, `remote remove`, and `remote list` commands as functions that return an `error` (or an error alongside a result) rather than calling `log.Fatalf` or `os.Exit` internally, so that every code path is reachable from a unit test. The cobra `Run` closures SHALL translate a returned error into the same exit-code and message behavior the CLI exhibits today.

#### Scenario: init_repo returns an error instead of exiting on a git failure
- **WHEN** `init_repo` is called with a path that is not a git repository
- **THEN** it returns a non-nil error describing the failure instead of calling `log.Fatalf`

#### Scenario: install command preserves today's exit behavior
- **WHEN** the `install` command's `Run` closure receives a non-nil error from `init_repo`
- **THEN** it calls `log.Fatalf` with that error, producing the same exit code and message the CLI produced before the refactor

#### Scenario: remote add business logic returns an error instead of exiting
- **WHEN** the extracted `remote add` logic is called against a path that is not a git repository
- **THEN** it returns a non-nil error instead of calling `log.Fatalf` directly

#### Scenario: remote add business logic reports an already-tracked remote without erroring
- **WHEN** the extracted `remote add` logic is called with a name already present in `.gitremotes`
- **THEN** it reports the remote as already tracked and returns a nil error, without appending a duplicate entry

#### Scenario: remote remove business logic distinguishes "not tracked" from other errors
- **WHEN** the extracted `remote remove` logic is called with a name not present in `.gitremotes`
- **THEN** it returns a distinguishable "not tracked" error, allowing the `Run` closure to print today's `error: remote '<name>' not found in .gitremotes` message rather than a generic `log.Fatalf`

#### Scenario: remote list business logic reports registration status
- **WHEN** the extracted `remote list` logic is called on a repository where one tracked remote is registered in `.git/config` and another is not
- **THEN** it returns a result marking the unregistered remote accordingly, without writing directly to stdout or calling `os.Exit`

### Requirement: removeFromGitremotes correctly rewrites the .gitremotes file
The system SHALL remove exactly the named remote's stanza from a `.gitremotes` file, leaving all other stanzas and their contents unchanged and in order.

#### Scenario: Removing a remote among several
- **WHEN** `.gitremotes` contains stanzas for `origin`, `mirror`, and `backup`, and `removeFromGitremotes` is called for `mirror`
- **THEN** the resulting file contains the `origin` and `backup` stanzas unchanged, in their original order, with no `mirror` stanza

#### Scenario: Removing the only remote
- **WHEN** `.gitremotes` contains a single stanza for `origin`, and `removeFromGitremotes` is called for `origin`
- **THEN** the resulting file contains no remote stanzas

#### Scenario: Removing a name that is not present
- **WHEN** `.gitremotes` contains stanzas that do not include the requested name
- **THEN** `removeFromGitremotes` returns nil and the file is unchanged

#### Scenario: Missing file returns an error
- **WHEN** `removeFromGitremotes` is called with a path that does not exist
- **THEN** it returns a non-nil error

### Requirement: pushRemote reports success, no-op, and failure correctly
The system SHALL push to the named remote and distinguish a successful push, a no-op push (already up to date), and a failed push.

#### Scenario: Successful push
- **WHEN** `pushRemote` is called for a remote registered in the repository with new commits to push
- **THEN** it returns nil

#### Scenario: Push when already up to date
- **WHEN** `pushRemote` is called for a remote whose ref is already up to date
- **THEN** it returns nil (the `git.NoErrAlreadyUpToDate` case is treated as success)

#### Scenario: Remote not registered
- **WHEN** `pushRemote` is called with a `remoteEntry` whose name is not registered in the repository's git config
- **THEN** it returns an error that mentions the remote name and suggests running `git-syn install`

### Requirement: syncAll dispatches to the configured push strategy
The system SHALL run pushes in parallel when `ActiveConfig.PushStrategy` is `"parallel"` (the default) and sequentially when it is `"sequential"`.

#### Scenario: Parallel strategy pushes all remotes concurrently
- **WHEN** `ActiveConfig.PushStrategy` is `"parallel"` and `syncAll` is called with multiple remote entries
- **THEN** `runParallel` pushes to all of them and `syncAll` returns its result

#### Scenario: Sequential strategy pushes remotes one at a time
- **WHEN** `ActiveConfig.PushStrategy` is `"sequential"` and `syncAll` is called with multiple remote entries
- **THEN** `runSequential` pushes to them in order and `syncAll` returns its result

### Requirement: runSequential respects on_failure policy
The system SHALL stop pushing to remaining remotes on the first failure when `ActiveConfig.OnFailure` is `"abort"`, and continue pushing to all remotes otherwise.

#### Scenario: Abort on first failure
- **WHEN** `ActiveConfig.OnFailure` is `"abort"` and a push to the first of three remotes fails
- **THEN** `runSequential` does not attempt the remaining two remotes and returns a non-zero result

#### Scenario: Continue past failures when not aborting
- **WHEN** `ActiveConfig.OnFailure` is `"warn"` and a push to the first of three remotes fails
- **THEN** `runSequential` still attempts the remaining two remotes

### Requirement: runParallel pushes to all remotes regardless of individual failures
The system SHALL attempt a push to every remote concurrently and collect all results, regardless of whether earlier pushes failed.

#### Scenario: One remote fails, others still attempted
- **WHEN** `runParallel` is called with three remote entries and one of them fails to push
- **THEN** the other two pushes are still attempted and their results are included in the outcome

### Requirement: cleanGitremotes removes registered remotes and the .gitremotes file
The system SHALL remove every remote listed in `.gitremotes` from the repository's git config and delete the `.gitremotes` file itself, tolerating remotes that are already absent from git config.

#### Scenario: Successful cleanup
- **WHEN** `cleanGitremotes` is called on a repository with a `.gitremotes` file listing remotes that are registered in git config
- **THEN** each listed remote is removed from git config and the `.gitremotes` file no longer exists, and nil is returned

#### Scenario: Remote already missing from git config
- **WHEN** a remote listed in `.gitremotes` is not present in the repository's git config
- **THEN** `cleanGitremotes` does not return an error for that remote and continues cleaning the rest

#### Scenario: Missing .gitremotes file is a no-op
- **WHEN** `cleanGitremotes` is called on a repository with no `.gitremotes` file
- **THEN** it returns nil without error

### Requirement: CI runs the test suite on every branch and merge request
The system SHALL run `mise run check` (format, lint, and test) in CI for every branch push and merge request pipeline, not only on `main`.

#### Scenario: Test job runs on a feature branch
- **WHEN** a commit is pushed to a non-main branch
- **THEN** the CI pipeline includes a job that runs `mise run check` and the pipeline fails if any test fails

#### Scenario: Test job runs on a merge request
- **WHEN** a merge request pipeline is triggered
- **THEN** the CI pipeline includes a job that runs `mise run check` and the pipeline fails if any test fails
