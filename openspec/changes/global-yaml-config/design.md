## Context

All git-syn behavior is currently hardcoded. The `pre-push-command` change introduces push failure handling and push strategy decisions that need user-configurable defaults. A user-level config file is the standard mechanism for CLI tools; the XDG Base Directory Specification (`~/.config/<app>/`) is the modern convention for Linux/macOS tools.

## Goals / Non-Goals

**Goals:**
- Load a YAML config file from the XDG config directory at startup
- Expose `on_failure`, `push_strategy`, and `allowed_schemes` as configurable fields
- Provide `git-syn config init` to scaffold the file with documented defaults
- Apply defaults silently when the file is absent

**Non-Goals:**
- Per-repo config file (future, not this change)
- Config editing subcommands (`set`, `get`) — the YAML file is edited by hand or via `config init`
- Windows `%APPDATA%` path support (can be added later)

## Decisions

**`gopkg.in/yaml.v3` for YAML parsing.**
Already common in the Go ecosystem, zero transitive deps, well-maintained. Alternative: `github.com/ghodss/yaml` (converts YAML to JSON then uses `encoding/json`). Rejected — unnecessary indirection for a simple config struct.

**Config loaded once at process startup into a package-level singleton.**
Each subcommand imports the config package and reads the singleton. Alternative: pass config as a parameter to every command. Rejected — the cobra command pattern doesn't naturally thread context through; a singleton is idiomatic for CLI config.

**`config init` writes a fully-commented YAML template, not a minimal one.**
Users are more likely to discover and tune settings when all options are visible with their defaults and descriptions. A minimal file makes the schema opaque.

**`$XDG_CONFIG_HOME` respected; falls back to `~/.config`.**
Standard XDG behavior. `os.UserConfigDir()` in Go stdlib implements this correctly on Linux/macOS.

## Risks / Trade-offs

[Config struct changes break existing files] → Use `yaml:"omitempty"` and default-on-zero-value semantics. Adding new fields is always backward-compatible; removing fields is a no-op (unknown keys are ignored by `yaml.v3` with `yaml:",inline"` or simply ignored by default).

[Package-level singleton is harder to test] → Acceptable for a CLI tool. Integration tests can set `$XDG_CONFIG_HOME` to a temp dir to control config state.
