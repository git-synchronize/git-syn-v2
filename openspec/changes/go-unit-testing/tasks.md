# Tasks

## 1. Refactor `validateRemoteURL` to accept config as a parameter
- [x] 1.1 Change `validateRemoteURL(rawURL string)` signature to `validateRemoteURL(rawURL string, allowedSchemes []string) error` in `cmd/util.go`
- [x] 1.2 Update all call sites of `validateRemoteURL` to pass `ActiveConfig.AllowedSchemes`
- [x] 1.3 Extract `loadConfigFromPath(path string) (*Config, error)` from `loadConfig` in `cmd/config.go`, and update `loadConfig` to call it

## 2. Install unit testing framework
- [x] 2.1 Run `go get github.com/stretchr/testify` to add testify to `go.mod`
- [x] 2.2 Run `go mod tidy` to update `go.sum`

## 3. Implement unit tests for `cmd/util.go`
- [x] 3.1 Create `cmd/util_test.go` with `package cmd`
- [x] 3.2 Add `TestMain` to set `ActiveConfig` to `defaultConfig()` before all tests run
- [x] 3.3 Add table-driven tests for `parseGitremotes`: single remote, multiple remotes, empty file, missing file, incomplete stanza
- [x] 3.4 Add table-driven tests for `validateRemoteURL`: HTTPS allowed, git@ SSH allowed, ssh:// SSH allowed, HTTP rejected, HTTPS rejected when SSH-only

## 4. Implement unit tests for `cmd/config.go`
- [x] 4.1 Create `cmd/config_test.go` with `package cmd`
- [x] 4.2 Add test for `loadConfigFromPath` with a missing path → returns defaults, nil error
- [x] 4.3 Add test for `loadConfigFromPath` with all fields set → returns correct values
- [x] 4.4 Add test for `loadConfigFromPath` with partial fields → missing fields use defaults
- [x] 4.5 Add test for `loadConfigFromPath` with invalid YAML → returns non-nil error

## 5. Verification and documentation
- [x] 5.1 Run `go test ./...` and confirm all tests pass with zero failures
- [x] 5.2 Run `mise run check` (fmt + lint + test) and confirm clean output
- [x] 5.3 Add a "Testing" section to `README.md` documenting `go test ./...` and `mise run test`
