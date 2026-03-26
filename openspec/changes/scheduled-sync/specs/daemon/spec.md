## ADDED Requirements

### Requirement: Daemon runs as a long-lived synchronization process
The system SHALL provide a `git syn daemon` command that runs as a foreground process and periodically pushes all configured remotes at a user-defined interval using the same push strategy and failure policy as the `pre-push` command.

#### Scenario: Daemon starts and logs first tick
- **WHEN** user runs `git syn daemon --interval 5m` in a repository with a valid `.gitremotes`
- **THEN** the daemon performs an initial sync cycle immediately and then repeats every 5 minutes until terminated

#### Scenario: Daemon uses interval from config
- **WHEN** `sync_interval: 10m` is set in `~/.config/git-syn/config.yaml` and user runs `git syn daemon` without `--interval`
- **THEN** the daemon uses a 10-minute interval

#### Scenario: Daemon errors when no interval is configured
- **WHEN** `sync_interval` is absent or zero in config and no `--interval` flag is provided
- **THEN** the daemon SHALL exit with a non-zero status and a descriptive error message indicating that an interval must be specified

#### Scenario: --interval flag overrides config
- **WHEN** `sync_interval: 10m` is set in config and user runs `git syn daemon --interval 2m`
- **THEN** the daemon uses a 2-minute interval, ignoring the config value

### Requirement: Daemon supports a single-cycle mode
The system SHALL provide a `--once` flag on `git syn daemon` that performs exactly one sync cycle and exits, enabling use from external schedulers such as cron or systemd timers.

#### Scenario: --once exits after one cycle
- **WHEN** user runs `git syn daemon --once`
- **THEN** the daemon performs one full sync cycle (pushing all remotes) and exits with the same exit code semantics as `git syn pre-push`

### Requirement: Daemon shuts down gracefully on SIGTERM or SIGINT
The system SHALL handle `SIGTERM` and `SIGINT` by completing any in-progress sync cycle before exiting, rather than terminating mid-push.

#### Scenario: SIGINT received while idle
- **WHEN** the daemon receives `SIGINT` between sync ticks
- **THEN** it exits cleanly with status 0

#### Scenario: SIGTERM received during active push
- **WHEN** the daemon receives `SIGTERM` while a sync cycle is in progress
- **THEN** it waits for the current cycle to complete and then exits

### Requirement: Daemon accepts a --path flag
The system SHALL accept a `--path` flag on `git syn daemon` to specify the repository root, consistent with `git syn pre-push --path`.

#### Scenario: Daemon runs against non-CWD repository
- **WHEN** user runs `git syn daemon --path /srv/git/myrepo --interval 15m`
- **THEN** the daemon synchronizes remotes for the repository at `/srv/git/myrepo`
