# Tasks

## 1. Implement the logging infrastructure
- [x] 1.1 In `cmd/root.go` `PersistentPreRun`, configure `slog.SetDefault` with a `slog.NewTextHandler` writing to `os.Stderr` at the appropriate level (`LevelInfo` for `--verbose`, `LevelDebug` for `--debug`); discard all output when neither flag is set
- [x] 1.2 Confirm debug implies verbose by setting the handler level to `LevelDebug` (which also passes `LevelInfo` records through)

## 2. Register command-line flags
- [x] 2.1 Register `--verbose` as a persistent bool flag on `rootCmd` in `cmd/root.go` (no short alias — `-v` is taken by `--version`)
- [x] 2.2 Register `--debug` as a persistent bool flag on `rootCmd` in `cmd/root.go`
- [x] 2.3 Add a `PersistentPreRun` hook on `rootCmd` that reads the flag values and sets the package-level log level
- [x] 2.4 Update the custom help output in `cmd/root.go` to list `--verbose` and `--debug` with descriptions

## 3. Instrument the codebase with logs
- [x] 3.1 Add `slog.Info` calls in `pre-push.go` for per-remote push start events (e.g. "pushing to remote %q")
- [x] 3.2 Add `slog.Debug` calls in `pre-push.go` for internal state (e.g. strategy chosen, goroutine count)
- [x] 3.3 Add `slog.Debug` calls in `util.go` for remote parsing events (e.g. entries read from `.gitremotes`)

## 4. Verification and testing
- [x] 4.1 Add a test that verifies `--verbose` and `--debug` are registered as persistent flags on the root command
- [x] 4.2 Add a test that verifies `Debugf` output appears when debug level is active and is suppressed otherwise
