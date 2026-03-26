## ADDED Requirements

### Requirement: Verbose flag
The CLI SHALL support a `--verbose` persistent flag on the root command that enables detailed operational output for all subcommands.

#### Scenario: Verbose flag present
- **WHEN** the user runs any git-syn subcommand with `--verbose`
- **THEN** verbose-level log messages are written to stderr in addition to normal output

#### Scenario: Verbose flag absent
- **WHEN** the user runs any git-syn subcommand without `--verbose` or `--debug`
- **THEN** verbose-level log messages are suppressed

#### Scenario: Verbose appears in help
- **WHEN** the user runs `git-syn --help`
- **THEN** `--verbose` is listed with a description of its effect

### Requirement: Debug flag
The CLI SHALL support a `--debug` persistent flag on the root command that enables low-level diagnostic output for all subcommands.

#### Scenario: Debug flag present
- **WHEN** the user runs any git-syn subcommand with `--debug`
- **THEN** debug-level log messages are written to stderr in addition to verbose-level messages

#### Scenario: Debug flag absent
- **WHEN** the user runs any git-syn subcommand without `--debug`
- **THEN** debug-level log messages are suppressed

#### Scenario: Debug appears in help
- **WHEN** the user runs `git-syn --help`
- **THEN** `--debug` is listed with a description of its effect

### Requirement: Debug implies verbose
The CLI SHALL treat `--debug` as a superset of `--verbose`.

#### Scenario: Debug enables verbose output
- **WHEN** the user runs a subcommand with `--debug` but without `--verbose`
- **THEN** both verbose-level and debug-level messages are written to stderr

### Requirement: Log output target
All verbose and debug log output SHALL be written to stderr.

#### Scenario: Stdout unaffected by flags
- **WHEN** the user runs any subcommand with `--verbose` or `--debug`
- **THEN** existing stdout output is unchanged and stderr receives the additional log lines

### Requirement: Flag availability across subcommands
Both `--verbose` and `--debug` SHALL be available as persistent flags, inherited by all subcommands without per-command registration.

#### Scenario: Flag inherited by subcommand
- **WHEN** the user runs `git-syn pre-push --verbose`
- **THEN** verbose logging is active for the pre-push execution

#### Scenario: Flag inherited by any future subcommand
- **WHEN** a new subcommand is added to the CLI
- **THEN** it inherits `--verbose` and `--debug` without any additional flag registration
