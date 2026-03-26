## Context

git-syn-v2's Go code (~830 LOC across 7 files in `cmd/`) has no test coverage. All logic lives in the `cmd` package, which currently relies on two package-level globals: `ActiveConfig *Config` (loaded at `init()` time) and `prePushPath string`. Most interesting functions — `parseGitremotes`, `validateRemoteURL`, `loadConfig` — are already pure or near-pure functions that are straightforward to test.

The `ActiveConfig` global is the primary testability obstacle: `validateRemoteURL` reads from it directly, making isolated unit tests require either a global mutation or a refactor.

## Goals / Non-Goals

**Goals:**
- Unit tests for `parseGitremotes` (table-driven, covering well-formed input, malformed input, empty file).
- Unit tests for `validateRemoteURL` (all allowed/disallowed scheme combinations).
- Unit tests for `loadConfig` (valid YAML, missing file → defaults, invalid YAML → error; extended to cover `sync_interval` parsing once that change lands).
- Tests run via `go test ./...` / `mise run test` with no additional setup.
- `testify/assert` used for readable assertions; `testify/require` for fatal assertions.

**Non-Goals:**
- Integration or end-to-end tests (pushing to a real git remote).
- Testing cobra command wiring or CLI flag parsing.
- Deleting or migrating the legacy C tests in `test/`.
- 100% code coverage — focus is on the highest-value, lowest-friction units.

## Decisions

### 1. Refactor `validateRemoteURL` to accept config as a parameter
**Decision**: Change the signature to `validateRemoteURL(rawURL string, allowedSchemes []string) error` and update all call sites to pass `ActiveConfig.AllowedSchemes`.
**Rationale**: Removes the global dependency, making the function trivially testable without touching `ActiveConfig`. The call sites are few (only `remote add`).
**Alternative considered**: Mutating `ActiveConfig` in test setup via `TestMain` — works but is fragile and prevents parallel tests.

### 2. Test files colocated in `cmd` package (`cmd/*_test.go`)
**Decision**: Use `package cmd` (white-box) rather than `package cmd_test` (black-box).
**Rationale**: Most functions under test are unexported (`parseGitremotes`, `validateRemoteURL`, `loadConfig`). White-box testing is the pragmatic choice given the current package structure.

### 3. Use `testify/assert` + `testify/require`
**Decision**: Add `github.com/stretchr/testify` as a test dependency.
**Rationale**: Table-driven tests with `assert.Equal` / `assert.ErrorContains` are significantly more readable than raw `t.Errorf` comparisons, especially for multi-case tables. `testify` is the de-facto standard in the Go ecosystem with zero runtime footprint.
**Alternative considered**: stdlib only — viable but results in more boilerplate for negligible benefit.

### 4. `loadConfig` testability via temp files
**Decision**: Test `loadConfig` by writing YAML to `os.CreateTemp` and temporarily overriding the config path using an environment variable or by extracting a `loadConfigFromPath(path string)` helper.
**Rationale**: `loadConfig` derives its path from `os.UserConfigDir()`, which is not easily injectable. Extracting a `loadConfigFromPath` helper is a minimal, non-breaking refactor that makes the function fully testable without OS-level mocking.

## Risks / Trade-offs

- **Global `ActiveConfig` mutation in tests**: If any test exercises a code path that reads `ActiveConfig` before `validateRemoteURL` is refactored, it could panic (nil pointer). → Mitigation: refactor `validateRemoteURL` first; ensure tests initialize `ActiveConfig` in `TestMain` as a safety net.
- **`cmd` package init() side effects**: `init()` in `config.go` calls `loadConfig()` and sets `ActiveConfig`, which runs when the test binary starts. If the config file on the developer's machine has invalid content, tests will fail. → Mitigation: document that tests assume a valid or absent config file; the `loadConfigFromPath` refactor reduces this risk.

## Migration Plan

1. Refactor `validateRemoteURL` to accept `allowedSchemes []string`.
2. Extract `loadConfigFromPath(path string) (*Config, error)` from `loadConfig`.
3. Write `cmd/util_test.go` (parseGitremotes, validateRemoteURL).
4. Write `cmd/config_test.go` (loadConfigFromPath, defaults, invalid YAML).
5. Run `go mod tidy` to add testify.
6. Verify `mise run test` passes.

## Open Questions

- Should `handleResults` be tested? It currently writes directly to stdout and calls `os.Exit`, making it harder to test without capturing output. → Lean: defer to a later change; focus this one on pure functions.
- Should we introduce a `TestMain` to set `ActiveConfig` to a safe default before all tests run? → Lean: yes, as a one-liner safety net even after the refactor.
