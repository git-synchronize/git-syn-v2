## ADDED Requirements

### Requirement: git-syn serve hosts repositories over HTTP
The system SHALL provide a `serve` command that runs an HTTP server hosting bare git repositories under a configured path, supporting both clone and push.

#### Scenario: Clone succeeds
- **WHEN** `git-syn serve --path <dir> --listen <addr>` is running and a bare repository exists under `<dir>`
- **THEN** a `git clone` against `http://<addr>/<repo-name>.git` SHALL succeed

#### Scenario: Push succeeds
- **WHEN** `git-syn serve` is running, a bare repository under its path has `http.receivepack` enabled, and a client pushes to it
- **THEN** the push SHALL succeed and the commit SHALL be present in the server-side repository

#### Scenario: Default path is the current directory
- **WHEN** `git-syn serve` is run without `--path`
- **THEN** it SHALL serve repositories from the current working directory, consistent with other git-syn commands' `--path` default

#### Scenario: Default listen address
- **WHEN** `git-syn serve` is run without `--listen`
- **THEN** it SHALL bind to `:8080`

### Requirement: git-http-backend is located portably
The system SHALL locate the `git-http-backend` binary via `git --exec-path` rather than a hardcoded filesystem path, and SHALL fail with a clear error at startup if it cannot be found.

#### Scenario: git-http-backend found via git --exec-path
- **WHEN** `git-syn serve` starts and `git --exec-path` reports a directory containing `git-http-backend`
- **THEN** `serve` SHALL use that binary to handle requests

#### Scenario: git-http-backend missing
- **WHEN** `git --exec-path`'s directory does not contain a `git-http-backend` binary
- **THEN** `git-syn serve` SHALL exit with a non-zero status and an error naming the missing binary and the path checked, rather than starting a server that would fail on first request

### Requirement: Graceful shutdown
The `serve` command SHALL shut down cleanly on `SIGINT`/`SIGTERM`, consistent with the `daemon` command's shutdown behavior.

#### Scenario: Signal triggers clean shutdown
- **WHEN** `git-syn serve` receives `SIGINT` or `SIGTERM`
- **THEN** it SHALL stop accepting new connections, allow in-flight requests to complete, and exit zero
