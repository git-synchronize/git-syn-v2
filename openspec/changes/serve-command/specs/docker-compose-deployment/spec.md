## ADDED Requirements

### Requirement: docker-compose deployment runs git-syn serve
The `deploy/docker-compose/git-http/` image SHALL run `git-syn serve` to host repositories, rather than building and running a separate bridge program.

#### Scenario: Image no longer builds a separate bridge program
- **WHEN** `deploy/docker-compose/git-http/` is built
- **THEN** the resulting image SHALL contain the `git-syn` binary and the `git` package, and SHALL NOT require a separate Go module or multi-stage build step for a CGI bridge program

#### Scenario: nginx continues to front the service
- **WHEN** the docker-compose stack is running
- **THEN** nginx SHALL continue to reverse-proxy to the git-hosting process (now `git-syn serve` instead of `cgi-bridge`), preserving the existing TLS/reverse-proxy documentation in the deployment's README
