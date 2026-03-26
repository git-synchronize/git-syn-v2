## Why

All git-syn behavior is currently hardcoded. Users who want different push failure handling or push ordering have no way to configure it without rebuilding the binary. A user-level config file enables personalization without per-repo changes.

## What Changes

- Load `~/.config/git-syn/config.yaml` (respecting `$XDG_CONFIG_HOME`) at startup
- Define three top-level config keys:
  - `on_failure: warn | abort` — whether a failing remote aborts the push or just prints a warning (default: `warn`)
  - `push_strategy: parallel | sequential` — how remotes are pushed to (default: `parallel`)
  - `allowed_schemes: [https, ssh]` — URL schemes accepted by `install` and `remote add` (default: `["https", "ssh"]`)
- If the file does not exist, all defaults apply silently — no error
- Add `git-syn config init` subcommand to write a commented default config file to the XDG path

## Capabilities

### New Capabilities

- `global-config`: User-level YAML configuration loaded from `~/.config/git-syn/config.yaml`; controls failure behavior, push strategy, and allowed URL schemes
- `config-init-command`: `git-syn config init` subcommand that writes a default commented config file

### Modified Capabilities

<!-- none -->

## Impact

- `cmd/config.go`: new file for config struct, loader, and `config init` subcommand
- `cmd/pre-push.go`: reads `on_failure` and `push_strategy` from loaded config
- `cmd/init.go`: reads `allowed_schemes` from loaded config
- New dependency: `gopkg.in/yaml.v3` (or `github.com/ghodss/yaml`)
