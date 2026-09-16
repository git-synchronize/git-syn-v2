## ADDED Requirements

### Requirement: Docker Compose stack for git mirror
The `deploy/docker-compose/` directory SHALL contain a `docker-compose.yml` that builds and runs a self-contained git HTTP server suitable for use as a git-syn sync target, with working push.

#### Scenario: Stack starts successfully
- **WHEN** the user runs `docker compose up -d --build`
- **THEN** the git-serve service SHALL build its image, start, and listen on port 80

#### Scenario: Push succeeds
- **WHEN** a bare repository has been created on the server with `http.receivepack` enabled, and the user pushes to it over HTTP
- **THEN** the push SHALL complete successfully and the commit SHALL be present in the server-side repository

#### Scenario: Push succeeds at realistic sizes
- **WHEN** the user pushes a commit containing tens of megabytes of new data
- **THEN** the push SHALL complete successfully, not hang or fail with an error

#### Scenario: Repository data persists across restarts
- **WHEN** the stack is stopped and restarted
- **THEN** repository data stored in the named volume SHALL persist

### Requirement: Quickstart documentation
The `deploy/docker-compose/` directory SHALL include a `README.md` with setup instructions.

#### Scenario: User follows quickstart
- **WHEN** the user follows the steps in `README.md`
- **THEN** the steps SHALL be sufficient to start the stack, create a bare repository with push enabled, and push to it as a git-syn remote
