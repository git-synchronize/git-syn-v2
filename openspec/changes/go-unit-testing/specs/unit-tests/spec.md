## ADDED Requirements

### Requirement: parseGitremotes parses well-formed .gitremotes files
The system SHALL correctly parse a `.gitremotes` file containing one or more `[remote "name"]` / `url = <url>` stanzas and return the corresponding `remoteEntry` slice.

#### Scenario: Single remote
- **WHEN** the file contains one remote stanza with a valid URL
- **THEN** `parseGitremotes` returns a slice with exactly one `remoteEntry` matching the name and URL

#### Scenario: Multiple remotes
- **WHEN** the file contains multiple remote stanzas
- **THEN** `parseGitremotes` returns all entries in order of appearance

#### Scenario: Empty file
- **WHEN** the file exists but is empty
- **THEN** `parseGitremotes` returns an empty slice with no error

#### Scenario: Missing file
- **WHEN** the file path does not exist
- **THEN** `parseGitremotes` returns an error satisfying `os.IsNotExist`

#### Scenario: Stanza with no url line is skipped
- **WHEN** a remote stanza header is present but no `url = ` line follows before the next stanza
- **THEN** the incomplete stanza is silently skipped and does not appear in the returned slice

### Requirement: validateRemoteURL accepts and rejects URLs by allowed schemes
The system SHALL accept URLs matching the configured allowed schemes and reject all others, returning a descriptive error.

#### Scenario: HTTPS URL allowed
- **WHEN** `allowedSchemes` includes `"https"` and the URL starts with `https://`
- **THEN** `validateRemoteURL` returns nil

#### Scenario: SSH URL with git@ prefix allowed
- **WHEN** `allowedSchemes` includes `"ssh"` and the URL starts with `git@`
- **THEN** `validateRemoteURL` returns nil

#### Scenario: SSH URL with ssh:// prefix allowed
- **WHEN** `allowedSchemes` includes `"ssh"` and the URL starts with `ssh://`
- **THEN** `validateRemoteURL` returns nil

#### Scenario: HTTP URL rejected
- **WHEN** `allowedSchemes` contains only `"https"` and `"ssh"` and the URL starts with `http://`
- **THEN** `validateRemoteURL` returns a non-nil error mentioning the unsupported scheme

#### Scenario: HTTPS URL rejected when only SSH allowed
- **WHEN** `allowedSchemes` contains only `"ssh"` and the URL starts with `https://`
- **THEN** `validateRemoteURL` returns a non-nil error

### Requirement: loadConfig returns defaults when config file is absent
The system SHALL return a `Config` populated with default values when no config file exists at the XDG path, without returning an error.

#### Scenario: Missing config file yields defaults
- **WHEN** no config file exists at the resolved path
- **THEN** `loadConfigFromPath` returns a `Config` with `OnFailure: "warn"`, `PushStrategy: "parallel"`, and `AllowedSchemes: ["https", "ssh"]`, and a nil error

### Requirement: loadConfig parses valid YAML correctly
The system SHALL parse a well-formed config YAML file and return a `Config` reflecting its values.

#### Scenario: All fields present
- **WHEN** the config file contains valid YAML for `on_failure`, `push_strategy`, and `allowed_schemes`
- **THEN** `loadConfigFromPath` returns a `Config` with fields matching the file values and a nil error

#### Scenario: Partial fields use defaults for missing keys
- **WHEN** the config file omits `push_strategy`
- **THEN** `loadConfigFromPath` returns `PushStrategy: "parallel"` (the default) and a nil error

### Requirement: loadConfig rejects invalid YAML
The system SHALL return a non-nil error when the config file contains malformed YAML that cannot be parsed.

#### Scenario: Malformed YAML returns error
- **WHEN** the config file contains invalid YAML (e.g., unclosed brackets, wrong indentation type)
- **THEN** `loadConfigFromPath` returns a non-nil error

### Requirement: Test suite runs via go test
The system SHALL provide a test suite that passes when executed with `go test ./...` in the repository root, with no manual setup beyond a standard Go toolchain.

#### Scenario: go test ./... exits zero
- **WHEN** `go test ./...` is run in the repository root with no config file present
- **THEN** all tests pass and the command exits with status 0

#### Scenario: mise run test runs the Go tests
- **WHEN** `mise run test` is invoked
- **THEN** it runs `go test ./...` and exits with the same status
