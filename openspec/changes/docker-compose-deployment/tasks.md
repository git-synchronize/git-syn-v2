## 1. Directory structure

- [x] 1.1 Create `deploy/docker-compose/` directory

## 2. docker-compose.yml

- [x] 2.1 Write `deploy/docker-compose/docker-compose.yml` with a `git-serve` service using a pinned git-http-backend image
- [x] 2.2 Configure the service to listen on port 80
- [x] 2.3 Mount a named volume for repository storage
- [x] 2.4 Pass `GIT_HTTP_EXPORT_ALL` and repository path via environment variables from `.env`

## 3. Environment variables

- [x] 3.1 Write `deploy/docker-compose/.env.example` documenting all required and optional variables with description comments

## 4. Documentation

- [x] 4.1 Write `deploy/docker-compose/README.md` with a quickstart: clone, copy `.env.example` to `.env`, `docker compose up`, add as git-syn remote
- [x] 4.2 Document that TLS is not handled by this stack and recommend fronting with a reverse proxy

## 5. Smoke test

- [ ] 5.1 Run `docker compose up` locally and verify `git clone` and `git push` work against the running server
