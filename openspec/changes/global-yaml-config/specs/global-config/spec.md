## ADDED Requirements

### Requirement: Load config from XDG config directory
The git-syn binary SHALL load a YAML config file from `$XDG_CONFIG_HOME/git-syn/config.yaml` (defaulting to `~/.config/git-syn/config.yaml`) at startup.

#### Scenario: Config file present and valid
- **WHEN** `~/.config/git-syn/config.yaml` exists and contains valid YAML
- **THEN** the values SHALL override the built-in defaults for all subsequent command operations

#### Scenario: Config file absent
- **WHEN** `~/.config/git-syn/config.yaml` does not exist
- **THEN** all built-in defaults SHALL apply silently
- **THEN** the command SHALL not print any warning or error about the missing file

#### Scenario: Config file present but invalid YAML
- **WHEN** `~/.config/git-syn/config.yaml` exists but contains invalid YAML
- **THEN** the command SHALL print an error identifying the file and the parse failure
- **THEN** the command SHALL exit non-zero

### Requirement: on_failure config key
The `on_failure` key SHALL control behavior when a remote push fails during `git-syn pre-push`.

#### Scenario: on_failure set to warn (default)
- **WHEN** `on_failure: warn` is set (or the key is absent)
- **AND** one or more remote pushes fail
- **THEN** failed remotes SHALL be reported as warnings
- **THEN** the overall command SHALL exit 0 if at least one remote succeeded

#### Scenario: on_failure set to abort
- **WHEN** `on_failure: abort` is set
- **AND** any remote push fails
- **THEN** the command SHALL print the failure and exit non-zero immediately

### Requirement: push_strategy config key
The `push_strategy` key SHALL control whether remotes are pushed concurrently or sequentially.

#### Scenario: push_strategy set to parallel (default)
- **WHEN** `push_strategy: parallel` is set (or the key is absent)
- **THEN** `git-syn pre-push` SHALL push to all remotes concurrently

#### Scenario: push_strategy set to sequential
- **WHEN** `push_strategy: sequential` is set
- **THEN** `git-syn pre-push` SHALL push to remotes one at a time in the order they appear in `.gitremotes`

### Requirement: allowed_schemes config key
The `allowed_schemes` key SHALL control which URL schemes are accepted by `install` and `remote add`.

#### Scenario: Custom allowed schemes
- **WHEN** `allowed_schemes: [https]` is set
- **AND** the user runs `git-syn remote add mirror git@github.com:user/repo.git`
- **THEN** the command SHALL reject the SSH URL and print an error

### Requirement: config init subcommand
The `git-syn config init` command SHALL write a commented default config file to the XDG config path.

#### Scenario: Config file does not exist
- **WHEN** the user runs `git-syn config init`
- **AND** no config file exists at the XDG path
- **THEN** the command SHALL create the directory if needed and write a fully-commented YAML file with all default values
- **THEN** the command SHALL print the path of the created file

#### Scenario: Config file already exists
- **WHEN** the user runs `git-syn config init`
- **AND** a config file already exists
- **THEN** the command SHALL print an error and exit non-zero without overwriting the file
