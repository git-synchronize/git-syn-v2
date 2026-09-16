# Tasks

## 1. Extract shared sync logic
- [x] 1.1 Extract `syncAll(repo *git.Repository, entries []remoteEntry) int` from `prePushCmd.Run` in `cmd/pre-push.go`
- [x] 1.2 Update `prePushCmd.Run` to call `syncAll` and verify `pre-push` behavior is unchanged

## 2. Update configuration to support intervals
- [x] 2.1 Add `SyncInterval time.Duration` field with `yaml:"sync_interval"` to the `Config` struct in `cmd/config.go`
- [x] 2.2 Implement custom YAML unmarshalling or post-parse step to accept a duration string (e.g., `5m`) and convert it to `time.Duration`
- [x] 2.3 Return an error from `loadConfig` when `sync_interval` is present but not a valid duration string
- [x] 2.4 Add a commented-out `sync_interval` entry to `configTemplate` in `cmd/config.go`

## 3. Implement the `daemon` command
- [x] 3.1 Create `cmd/daemon.go` with `daemonCmd` registered to root via `init()`
- [x] 3.2 Add `--path` and `--interval` flags mirroring the `pre-push` flag conventions
- [x] 3.3 Implement interval resolution: `--interval` flag > `ActiveConfig.SyncInterval`; error if both are zero
- [x] 3.4 Implement the ticker loop using `time.NewTicker` and `context` cancellation via `signal.NotifyContext`
- [x] 3.5 Implement `--once` mode: call `syncAll` once and exit with the same exit code semantics as `pre-push`
- [x] 3.6 Ensure graceful shutdown: wait for any in-progress `syncAll` call to complete before returning from the command

## 4. Verification and testing
- [x] 4.1 Add a unit test for the duration-string parsing in `loadConfig` (valid, zero, invalid inputs)
- [x] 4.2 Add a unit test for interval resolution logic (flag priority over config)

## 5. Documentation
- [x] 5.1 Update `README.md` with a section describing the `daemon` command and example invocations
- [ ] 5.2 Add a man page source entry for `git-syn-daemon(1)` in `doc/man/` if a man page generator is configured
