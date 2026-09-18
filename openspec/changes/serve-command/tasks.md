## 1. Core `serve` command

- [x] 1.1 Create `cmd/serve.go` with a `serveCmd` cobra command registered via `init()`, `--path` (default cwd, matching other commands) and `--listen` (default `:8080`) flags
- [x] 1.2 Implement `resolveGitHTTPBackend() (string, error)`: runs `git --exec-path`, joins `git-http-backend` onto the result, `os.Stat`s it, and returns a clear error naming the checked path if missing (split into `gitExecPath()` + `resolveGitHTTPBackend(execPath string)` for testability)
- [x] 1.3 Implement the HTTP handler using `net/http/cgi.Handler{Path: <resolved path>, Env: []string{"GIT_PROJECT_ROOT=" + path, "GIT_HTTP_EXPORT_ALL="}}`, mirroring `cgi-bridge/main.go`
- [x] 1.4 Implement graceful shutdown: `signal.NotifyContext` for `SIGINT`/`SIGTERM` (matching `daemon`'s pattern), calling `http.Server.Shutdown(ctx)`

## 2. Tests

- [x] 2.1 Unit test `resolveGitHTTPBackend`: found case (real `git --exec-path` on the test machine, which has git installed) and a not-found case (inject a fake `PATH`/exec-path pointing somewhere without `git-http-backend`, or extract the join+stat logic to accept the exec-path as a parameter so it's testable without mocking `exec.Command`)
- [x] 2.2 e2e test (`cmd/e2e_test.go`, reusing `buildGitSynBinary`): start `git-syn serve --path <tmpdir> --listen 127.0.0.1:<free port>` as a background subprocess, create a bare repo, real `git clone` and `git push` against it, then send `SIGTERM` and confirm the process exits zero within a reasonable timeout

## 3. Documentation

- [x] 3.1 Add a `serve` entry to `doc/man/git-syn.1.md`'s COMMANDS section, matching the detail level of the `daemon` entry (mentions `--path` and `--listen`)
- [x] 3.2 Add `serve` to the command list in the root `README.md`'s usage output (confirmed it matches live `git-syn -h` output exactly)

## 4. Update the docker-compose deployment

- [x] 4.1 Rewrite `deploy/docker-compose/git-http/Dockerfile`: multi-stage build compiling `git-syn` itself (build context must be the repo root, not `git-http/`, since git-syn's `go.mod` lives there — updated `docker-compose.yml`'s `build:` block accordingly), runtime stage installs `git` and copies in the `git-syn` binary
- [x] 4.2 Update `deploy/docker-compose/git-http/entrypoint.sh` to run `git-syn serve --path /var/lib/git --listen 127.0.0.1:8080` instead of starting `cgi-bridge`
- [x] 4.3 Confirm `deploy/docker-compose/git-http/nginx.conf` needs no changes (still proxies to `127.0.0.1:8080`)
- [x] 4.4 Delete `deploy/docker-compose/git-http/cgi-bridge/`
- [x] 4.5 Update `deploy/docker-compose/README.md`: replace references to `cgi-bridge` with `git-syn serve`, update the architecture description at the top

## 5. Verification

- [x] 5.1 Run `go test ./...` and `mise run check`; confirm clean
- [x] 5.2 Run `docker compose up --build` against the updated deployment and re-verify clone, push (including at a larger size, e.g. 20MB, given a prior regression was caught exactly this way), and restart persistence, exactly as done for the `cgi-bridge` version. Caught and fixed a real issue along the way: switching the Docker build context to the repo root (needed to build git-syn itself) broke the `COPY nginx.conf`/`COPY entrypoint.sh` lines, which were still relative to the old `git-http/` context; fixed by pointing them at `deploy/docker-compose/git-http/...` instead.
- [x] 5.3 Confirm `git-syn serve --path <dir>` run directly (outside Docker) also works, as a sanity check that the command isn't accidentally coupled to the container environment
