## 1. Dependency

- [x] 1.1 Add `gopkg.in/yaml.v3` to `go.mod` and `go.sum`

## 2. Config struct and loader

- [x] 2.1 Create `cmd/config.go` with a `Config` struct (`OnFailure`, `PushStrategy`, `AllowedSchemes` fields)
- [x] 2.2 Define default values as constants or a `defaultConfig()` function
- [x] 2.3 Implement `loadConfig() (*Config, error)` that resolves the XDG path, reads the file, unmarshals YAML, and applies defaults for missing fields
- [x] 2.4 Expose a package-level `ActiveConfig` singleton initialized in a `cobra.PersistentPreRun` or `init()`

## 3. Integrate config into commands

- [x] 3.1 Read `ActiveConfig.AllowedSchemes` in `cmd/init.go` `install` URL validation (replace hardcoded `["https", "ssh"]`)
- [x] 3.2 Read `ActiveConfig.OnFailure` in `cmd/pre-push.go` push result handling
- [x] 3.3 Read `ActiveConfig.PushStrategy` in `cmd/pre-push.go` to switch between concurrent and sequential push

## 4. config init subcommand

- [x] 4.1 Implement `configInitCmd` that resolves the XDG config path
- [x] 4.2 Create the config directory if it does not exist
- [x] 4.3 Write a fully-commented default YAML template to the config file
- [x] 4.4 Error and exit non-zero if the config file already exists
- [x] 4.5 Print `Created config at <path>` on success

## 5. Register and document

- [x] 5.1 Register a `configCmd` parent and `configInitCmd` child in `cmd/root.go`
- [x] 5.2 Manually verify: missing config file applies defaults silently; invalid YAML prints an error and exits non-zero
