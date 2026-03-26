## ADDED Requirements

### Requirement: Config supports a sync_interval field
The system SHALL accept an optional `sync_interval` field in `~/.config/git-syn/config.yaml` whose value is a Go duration string (e.g., `5m`, `1h`, `30s`). When present and non-zero, this value is used as the default interval for `git syn daemon`.

#### Scenario: Valid duration string is accepted
- **WHEN** `sync_interval: 5m` is present in the config file
- **THEN** `loadConfig()` parses it as `5 * time.Minute` with no error

#### Scenario: Zero or absent field disables default interval
- **WHEN** `sync_interval` is absent from the config file or set to `0`
- **THEN** `ActiveConfig.SyncInterval` is zero and the daemon requires an explicit `--interval` flag

#### Scenario: Invalid duration string is rejected
- **WHEN** `sync_interval: bogus` is present in the config file
- **THEN** `loadConfig()` returns an error and git-syn exits with a non-zero status on startup

#### Scenario: config init includes sync_interval documentation
- **WHEN** user runs `git syn config init`
- **THEN** the generated config file includes a commented-out `sync_interval` entry with a description of accepted values
