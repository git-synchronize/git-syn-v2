## Context

git-syn is a Go CLI built with Cobra. Currently the tool uses `fmt.Printf` for operational output and `log.Fatalf` for fatal errors — there is no log-level concept. Adding `--verbose` and `--debug` flags requires a lightweight leveled logging layer and wiring it into the Cobra root command.

**Constraint:** `-v` / `--version` is already registered on `rootCmd` in `cmd/root.go:33`. A short flag `-v` for `--verbose` would collide.

## Goals / Non-Goals

**Goals:**
- Add `--verbose` (no short flag) and `--debug` (no short flag) as persistent flags on `rootCmd`
- Use stdlib `log/slog` for leveled, structured logging (available since Go 1.21; module is now on 1.24)
- Gate verbose-level output and debug-level output behind the respective flags
- Write all log output to stderr

**Non-Goals:**
- Log file output or log rotation
- Removing `-v` / `--version` (leave as-is to avoid a breaking change)

## Decisions

### Decision 1: Short flags for `--verbose` and `--debug`

`-v` is already taken by `--version`. Assigning `-v` to `--verbose` would require removing or rebinding the version flag, which is a breaking change.

**Choice:** Register both `--verbose` and `--debug` with **no short-flag aliases**. Users who want brevity can use `--verbose` or `--debug` in full. This is consistent with many modern CLIs (e.g., `kubectl`, `helm`).

**Alternatives considered:**
- Reassign `-v` to `--verbose` and move `--version` to `-V`: breaking for existing scripts; deferred.
- Use `-V` for verbose: non-standard and unintuitive.

### Decision 2: Logging implementation — `log/slog`

Now that the module is on Go 1.24, `log/slog` is available with no new dependencies. A `slog.Logger` writing to `os.Stderr` via `slog.NewTextHandler` will be configured at startup with the appropriate minimum level (`slog.LevelInfo` for verbose, `slog.LevelDebug` for debug, or discarded entirely when neither flag is set).

Call sites use `slog.Info` for verbose-level messages and `slog.Debug` for debug-level messages. The logger is set as the default via `slog.SetDefault` so no logger instance needs to be threaded through the call graph.

**Alternatives considered:**
- Hand-rolled leveled logger: was the plan under Go 1.20; no longer necessary.
- `github.com/sirupsen/logrus` / `go.uber.org/zap`: feature-rich but heavyweight for a simple two-level gate.

### Decision 3: Flag wiring via `PersistentPreRun`

Cobra evaluates persistent flags before any command runs. The flag values (`--verbose`, `--debug`) are read in a `PersistentPreRun` hook on `rootCmd` and used to configure the package-level logger. This ensures subcommands automatically inherit the log level without any per-command wiring.

### Decision 4: Debug implies verbose

When `--debug` is set, verbose output is also enabled. This is the standard expectation (debug is a superset of verbose).

## Risks / Trade-offs

- **Short-flag gap**: No `-v` for verbose may surprise users expecting GNU-style `-v`. This is documented in `--help` output. Mitigated by clear flag descriptions.
- **slog text format verbosity**: `slog.NewTextHandler` emits `time=... level=... msg=...` lines which may feel noisy. The handler can be swapped for a minimal writer later if needed without changing call sites.
- **Existing `fmt.Printf` call sites**: Operational output in `pre-push.go` (e.g. `✔ remote-name`) currently goes to stdout. Verbose logging supplements this — it does not replace it. No existing output is silenced by this change.
