/*
Copyright © 2024-2026 Lucas Ramage <lucas.ramage@infinite-omicron.com>
*/
package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveGitHTTPBackend(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		execPath := t.TempDir()
		backendPath := filepath.Join(execPath, "git-http-backend")
		require.NoError(t, os.WriteFile(backendPath, []byte("#!/bin/sh\n"), 0755))

		got, err := resolveGitHTTPBackend(execPath)
		require.NoError(t, err)
		assert.Equal(t, backendPath, got)
	})

	t.Run("not found", func(t *testing.T) {
		execPath := t.TempDir()

		_, err := resolveGitHTTPBackend(execPath)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "git-http-backend not found")
		assert.Contains(t, err.Error(), execPath)
	})
}

func TestGitExecPath(t *testing.T) {
	path, err := gitExecPath()
	require.NoError(t, err)
	assert.NotEmpty(t, path)
}
