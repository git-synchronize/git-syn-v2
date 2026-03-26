## 1. Logger Setup

- [ ] 1.1 In `cmd/root.go` `PersistentPreRun`, configure `slog.SetDefault` with a `slog.NewTextHandler` writing to `os.Stderr` at the appropriate level (`LevelInfo` for `--verbose`, `LevelDebug` for `--debug`); discard all output when neither flag is set
- [ ] 1.2 Confirm debug implies verbose by setting the handler level to `LevelDebug` (which also passes `LevelInfo` records through)

## 2. Root Command Flags

- [ ] 2.1 Register `--verbose` as a persistent bool flag on `rootCmd` in `cmd/root.go` (no short alias — `-v` is taken by `--version`)
- [ ] 2.2 Register `--debug` as a persistent bool flag on `rootCmd` in `cmd/root.go`
- [ ] 2.3 Add a `PersistentPreRun` hook on `rootCmd` that reads the flag values and sets the package-level log level
- [ ] 2.4 Update the custom help output in `cmd/root.go` to list `--verbose` and `--debug` with descriptions

## 3. Instrumentation

- [ ] 3.1 Add `Verbosef` calls in `pre-push.go` for per-remote push start events (e.g. "pushing to remote %q")
- [ ] 3.2 Add `Debugf` calls in `pre-push.go` for internal state (e.g. strategy chosen, goroutine count)
- [ ] 3.3 Add `Debugf` calls in `util.go` for remote parsing events (e.g. entries read from `.gitremotes`)

## 4. Tests

- [ ] 4.1 Add a test that verifies `--verbose` and `--debug` are registered as persistent flags on the root command
- [ ] 4.2 Add a test that verifies `Debugf` output appears when debug level is active and is suppressed otherwise
