/*
Copyright © 2024-2026 Lucas Ramage <lucas.ramage@infinite-omicron.com>
*/
package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigFromPath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "git-syn-config-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	t.Run("missing path returns defaults", func(t *testing.T) {
		path := filepath.Join(tmpDir, "non-existent.yaml")
		cfg, err := loadConfigFromPath(path)
		assert.NoError(t, err)
		assert.Equal(t, defaultConfig(), cfg)
	})

	t.Run("all fields set", func(t *testing.T) {
		path := filepath.Join(tmpDir, "full.yaml")
		content := `
on_failure: abort
push_strategy: sequential
allowed_schemes:
  - https
sync_interval: 10m
`
		err := os.WriteFile(path, []byte(content), 0644)
		require.NoError(t, err)

		cfg, err := loadConfigFromPath(path)
		assert.NoError(t, err)
		assert.Equal(t, "abort", cfg.OnFailure)
		assert.Equal(t, "sequential", cfg.PushStrategy)
		assert.Equal(t, []string{"https"}, cfg.AllowedSchemes)
		assert.Equal(t, 10*time.Minute, cfg.SyncInterval)
	})

	t.Run("partial fields use defaults", func(t *testing.T) {
		path := filepath.Join(tmpDir, "partial.yaml")
		content := `
on_failure: abort
`
		err := os.WriteFile(path, []byte(content), 0644)
		require.NoError(t, err)

		cfg, err := loadConfigFromPath(path)
		assert.NoError(t, err)
		assert.Equal(t, "abort", cfg.OnFailure)
		assert.Equal(t, defaultPushStrategy, cfg.PushStrategy)
		assert.Equal(t, defaultAllowedSchemes, cfg.AllowedSchemes)
	})

	t.Run("invalid YAML returns error", func(t *testing.T) {
		path := filepath.Join(tmpDir, "invalid.yaml")
		content := `
on_failure: [unclosed bracket
`
		err := os.WriteFile(path, []byte(content), 0644)
		require.NoError(t, err)

		_, err = loadConfigFromPath(path)
		assert.Error(t, err)
	})
}
