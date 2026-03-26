## Why

git-syn currently provides no way for users to control log verbosity, making it difficult to understand what the tool is doing during sync operations or to diagnose issues. Adding `--verbose` and `--debug` flags aligns with GNU CLI conventions and gives users the observability they need for day-to-day use and troubleshooting.

## What Changes

- Add a global `--verbose` (`-v`) flag that enables detailed operational output across all subcommands
- Add a global `--debug` (`-d`) flag that enables low-level diagnostic output across all subcommands
- Both flags are persistent and available to all subcommands (not scoped to a single command)
- Debug implies verbose (enabling `--debug` also enables verbose output)
- Log output is written to stderr to keep stdout clean for machine-readable output

## Capabilities

### New Capabilities

- `verbose-debug-logging`: Global CLI flags (`--verbose`, `--debug`) and a logging layer that gates output at appropriate levels throughout the application

### Modified Capabilities

<!-- No existing spec-level requirements are changing -->

## Impact

- `main.go` / root command: add persistent flags and wire up log level
- All subcommands: replace ad-hoc `fmt.Print*` diagnostic output with leveled log calls
- No breaking changes to existing command interfaces or output formats
