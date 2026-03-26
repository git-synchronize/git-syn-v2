/*
Copyright © 2024-2026 Lucas Ramage <lucas.ramage@infinite-omicron.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const defaultOnFailure = "warn"
const defaultPushStrategy = "parallel"
const defaultSyncInterval = 5 * time.Minute

var defaultAllowedSchemes = []string{"https", "ssh"}

// Config holds the global git-syn configuration loaded from the XDG config file.
type Config struct {
	OnFailure      string        `yaml:"on_failure"`
	PushStrategy   string        `yaml:"push_strategy"`
	AllowedSchemes []string      `yaml:"allowed_schemes"`
	SyncInterval   time.Duration `yaml:"sync_interval"`
}

// ActiveConfig is the package-level singleton loaded at process startup.
var ActiveConfig *Config

func defaultConfig() *Config {
	return &Config{
		OnFailure:      defaultOnFailure,
		PushStrategy:   defaultPushStrategy,
		AllowedSchemes: defaultAllowedSchemes,
		SyncInterval:   defaultSyncInterval,
	}
}

func configFilePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "git-syn", "config.yaml"), nil
}

func loadConfig() (*Config, error) {
	path, err := configFilePath()
	if err != nil {
		return defaultConfig(), nil
	}

	return loadConfigFromPath(path)
}

func loadConfigFromPath(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return defaultConfig(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read config %s: %w", path, err)
	}

	cfg := defaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config %s: %w", path, err)
	}

	if cfg.OnFailure == "" {
		cfg.OnFailure = defaultOnFailure
	}
	if cfg.PushStrategy == "" {
		cfg.PushStrategy = defaultPushStrategy
	}
	if len(cfg.AllowedSchemes) == 0 {
		cfg.AllowedSchemes = defaultAllowedSchemes
	}
	if cfg.SyncInterval == 0 {
		cfg.SyncInterval = defaultSyncInterval
	}

	return cfg, nil
}

const configTemplate = `# git-syn configuration file

# on_failure: controls behavior when a remote push fails
# Options: warn (default), abort
# - warn: report failed remotes as warnings, exit 0 if at least one remote succeeded
# - abort: exit non-zero immediately if any remote push fails
on_failure: warn

# push_strategy: controls whether remotes are pushed concurrently or sequentially
# Options: parallel (default), sequential
# - parallel: push to all remotes concurrently
# - sequential: push to remotes one at a time in the order they appear in .gitremotes
push_strategy: parallel

# sync_interval: default interval for the 'daemon' command (e.g., 5m, 1h)
# sync_interval: 5m

# allowed_schemes: URL schemes accepted by install and remote add
# Options: https, ssh
allowed_schemes:
  - https
  - ssh
`

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "create a default config file",
	Long:  `Create a fully-commented default config file at the XDG config path (~/.config/git-syn/config.yaml).`,
	Run: func(cmd *cobra.Command, args []string) {
		path, err := configFilePath()
		if err != nil {
			log.Fatalf("failed to resolve config path: %v", err)
		}

		if _, err := os.Stat(path); err == nil {
			fmt.Fprintf(os.Stderr, "Config already exists at %s. Remove it first to reinitialize.\n", path)
			os.Exit(1)
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			log.Fatalf("failed to create config directory: %v", err)
		}

		if err := os.WriteFile(path, []byte(configTemplate), 0644); err != nil {
			log.Fatalf("failed to write config file: %v", err)
		}

		fmt.Printf("Created config at %s\n", path)
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "manage git-syn configuration",
}

func init() {
	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	ActiveConfig = cfg

	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configInitCmd)
}
