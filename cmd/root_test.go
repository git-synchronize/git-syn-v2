/*
Copyright © 2024-2026 Lucas Ramage <lucas.ramage@infinite-omicron.com>
*/
package cmd

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggingFlagsRegistration(t *testing.T) {
	v, err := rootCmd.PersistentFlags().GetBool("verbose")
	assert.NoError(t, err)
	assert.False(t, v)

	d, err := rootCmd.PersistentFlags().GetBool("debug")
	assert.NoError(t, err)
	assert.False(t, d)
}

func TestSlogConfiguration(t *testing.T) {
	var buf bytes.Buffer

	t.Run("verbose enabled", func(t *testing.T) {
		buf.Reset()
		handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
		logger := slog.New(handler)

		logger.Info("info message")
		logger.Debug("debug message")

		output := buf.String()
		assert.Contains(t, output, "level=INFO")
		assert.Contains(t, output, "msg=\"info message\"")
		assert.NotContains(t, output, "level=DEBUG")
	})

	t.Run("debug enabled", func(t *testing.T) {
		buf.Reset()
		handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
		logger := slog.New(handler)

		logger.Info("info message")
		logger.Debug("debug message")

		output := buf.String()
		assert.Contains(t, output, "level=INFO")
		assert.Contains(t, output, "level=DEBUG")
		assert.Contains(t, output, "msg=\"debug message\"")
	})
}
