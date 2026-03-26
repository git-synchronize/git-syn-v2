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

func TestMain(m *testing.M) {
	// Initialize ActiveConfig with defaults for tests
	ActiveConfig = defaultConfig()
	os.Exit(m.Run())
}

func TestValidateRemoteURL(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		allowedSchemes []string
		wantErr        bool
	}{
		{
			name:           "HTTPS allowed",
			url:            "https://github.com/user/repo.git",
			allowedSchemes: []string{"https"},
			wantErr:        false,
		},
		{
			name:           "HTTPS rejected when only SSH allowed",
			url:            "https://github.com/user/repo.git",
			allowedSchemes: []string{"ssh"},
			wantErr:        true,
		},
		{
			name:           "git@ SSH allowed",
			url:            "git@github.com:user/repo.git",
			allowedSchemes: []string{"ssh"},
			wantErr:        false,
		},
		{
			name:           "ssh:// SSH allowed",
			url:            "ssh://git@github.com/user/repo.git",
			allowedSchemes: []string{"ssh"},
			wantErr:        false,
		},
		{
			name:           "HTTP rejected by default",
			url:            "http://github.com/user/repo.git",
			allowedSchemes: []string{"https", "ssh"},
			wantErr:        true,
		},
		{
			name:           "invalid scheme",
			url:            "ftp://github.com/user/repo.git",
			allowedSchemes: []string{"https", "ssh"},
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRemoteURL(tt.url, tt.allowedSchemes)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseGitremotes(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "git-syn-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name     string
		content  string
		expected []remoteEntry
		wantErr  bool
	}{
		{
			name:    "single remote",
			content: "[remote \"origin\"]\n\turl = https://github.com/user/repo.git\n",
			expected: []remoteEntry{
				{name: "origin", url: "https://github.com/user/repo.git"},
			},
			wantErr: false,
		},
		{
			name:    "multiple remotes",
			content: "[remote \"origin\"]\n\turl = https://github.com/user/repo.git\n[remote \"backup\"]\n\turl = git@github.com:user/repo-backup.git\n",
			expected: []remoteEntry{
				{name: "origin", url: "https://github.com/user/repo.git"},
				{name: "backup", url: "git@github.com:user/repo-backup.git"},
			},
			wantErr: false,
		},
		{
			name:     "empty file",
			content:  "",
			expected: nil,
			wantErr:  false,
		},
		{
			name:     "incomplete stanza",
			content:  "[remote \"incomplete\"]\n",
			expected: nil,
			wantErr:  false,
		},
		{
			name:     "missing url",
			content:  "[remote \"missing\"]\n\tother = value\n",
			expected: nil,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, tt.name)
			err := os.WriteFile(path, []byte(tt.content), 0644)
			require.NoError(t, err)

			entries, err := parseGitremotes(path)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, entries)
			}
		})
	}

	t.Run("missing file", func(t *testing.T) {
		_, err := parseGitremotes(filepath.Join(tmpDir, "does-not-exist"))
		assert.Error(t, err)
	})
}
