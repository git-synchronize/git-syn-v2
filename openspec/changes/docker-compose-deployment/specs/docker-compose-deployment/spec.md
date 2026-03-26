## ADDED Requirements

### Requirement: Docker Compose stack for git mirror
The `deploy/docker-compose/` directory SHALL contain a `docker-compose.yml` that runs a self-contained git HTTP server suitable for use as a git-syn sync target.

#### Scenario: Stack starts successfully
- **WHEN** the user copies `.env.example` to `.env`, sets required variables, and runs `docker compose up`
- **THEN** the git-serve service SHALL start and listen on port 80
- **THEN** the service SHALL serve git repositories from the configured repository path

#### Scenario: Repository data persists across restarts
- **WHEN** the stack is stopped and restarted
- **THEN** repository data stored in the named volume SHALL persist

### Requirement: Environment variable configuration
The stack SHALL be configurable via environment variables documented in `.env.example`.

#### Scenario: .env.example present
- **WHEN** the user inspects `deploy/docker-compose/.env.example`
- **THEN** all required environment variables SHALL be listed with a description comment
- **THEN** default values SHALL be provided where applicable

### Requirement: Quickstart documentation
The `deploy/docker-compose/` directory SHALL include a `README.md` with setup instructions.

#### Scenario: User follows quickstart
- **WHEN** the user follows the steps in `README.md`
- **THEN** the steps SHALL be sufficient to start the stack and push a repository to it as a git-syn remote
