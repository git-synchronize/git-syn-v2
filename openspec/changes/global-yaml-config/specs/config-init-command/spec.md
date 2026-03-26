## ADDED Requirements

### Requirement: config init subcommand
The `git-syn config init` command SHALL scaffold a default, fully-commented config file at the XDG config path.

#### Scenario: Successful init
- **WHEN** the user runs `git-syn config init`
- **AND** no config file exists at `$XDG_CONFIG_HOME/git-syn/config.yaml`
- **THEN** the directory `$XDG_CONFIG_HOME/git-syn/` SHALL be created if it does not exist
- **THEN** a YAML file with all supported keys, their default values, and inline comments SHALL be written
- **THEN** the command SHALL print `Created config at <path>`

#### Scenario: Config file already exists
- **WHEN** the user runs `git-syn config init`
- **AND** `$XDG_CONFIG_HOME/git-syn/config.yaml` already exists
- **THEN** the command SHALL print an error: `Config already exists at <path>. Remove it first to reinitialize.`
- **THEN** the existing file SHALL not be modified
- **THEN** the command SHALL exit non-zero
