/*
Copyright © 2024-2026 Lucas Ramage <lucas.ramage@infinite-omicron.com>
*/
package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemoveFromGitremotes(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		remove   string
		expected string
	}{
		{
			name:     "removing a remote among several preserves order",
			content:  "[remote \"origin\"]\n\turl = https://github.com/user/repo.git\n[remote \"mirror\"]\n\turl = git@github.com:user/mirror.git\n[remote \"backup\"]\n\turl = https://github.com/user/backup.git\n",
			remove:   "mirror",
			expected: "[remote \"origin\"]\n\turl = https://github.com/user/repo.git\n[remote \"backup\"]\n\turl = https://github.com/user/backup.git\n",
		},
		{
			name:     "removing the only remote leaves an empty file",
			content:  "[remote \"origin\"]\n\turl = https://github.com/user/repo.git\n",
			remove:   "origin",
			expected: "",
		},
		{
			name:     "removing a name not present leaves the file unchanged",
			content:  "[remote \"origin\"]\n\turl = https://github.com/user/repo.git\n",
			remove:   "nonexistent",
			expected: "[remote \"origin\"]\n\turl = https://github.com/user/repo.git\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), ".gitremotes")
			require.NoError(t, os.WriteFile(path, []byte(tt.content), 0644))

			require.NoError(t, removeFromGitremotes(path, tt.remove))

			data, err := os.ReadFile(path)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(data))
		})
	}

	t.Run("missing file returns an error", func(t *testing.T) {
		err := removeFromGitremotes(filepath.Join(t.TempDir(), "does-not-exist"), "origin")
		assert.Error(t, err)
	})
}

func TestAddRemote(t *testing.T) {
	t.Run("adds a new remote to .gitremotes and git config", func(t *testing.T) {
		path := t.TempDir()
		_, err := git.PlainInit(path, false)
		require.NoError(t, err)

		alreadyTracked, err := addRemote(path, "origin", "https://example.com/user/repo.git")
		require.NoError(t, err)
		assert.False(t, alreadyTracked)

		entries, err := parseGitremotes(filepath.Join(path, ".gitremotes"))
		require.NoError(t, err)
		require.Len(t, entries, 1)
		assert.Equal(t, "origin", entries[0].name)
		assert.Equal(t, "https://example.com/user/repo.git", entries[0].url)

		repo, err := git.PlainOpen(path)
		require.NoError(t, err)
		cfg, err := repo.Config()
		require.NoError(t, err)
		_, registered := cfg.Remotes["origin"]
		assert.True(t, registered)
	})

	t.Run("reports an already-tracked remote without duplicating it", func(t *testing.T) {
		path := t.TempDir()
		_, err := git.PlainInit(path, false)
		require.NoError(t, err)

		gitremotesPath := filepath.Join(path, ".gitremotes")
		require.NoError(t, appendToGitremotes(gitremotesPath, "origin", "https://example.com/user/repo.git"))

		alreadyTracked, err := addRemote(path, "origin", "https://example.com/user/repo.git")
		require.NoError(t, err)
		assert.True(t, alreadyTracked)

		entries, err := parseGitremotes(gitremotesPath)
		require.NoError(t, err)
		assert.Len(t, entries, 1)
	})

	t.Run("returns an error for a non-git-repository path", func(t *testing.T) {
		path := t.TempDir()

		_, err := addRemote(path, "origin", "https://example.com/user/repo.git")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "is not a git repository")
	})
}

func TestRemoveRemote(t *testing.T) {
	t.Run("removes a tracked remote from .gitremotes and git config", func(t *testing.T) {
		path := t.TempDir()
		repo, err := git.PlainInit(path, false)
		require.NoError(t, err)

		_, err = repo.CreateRemote(&gitconfig.RemoteConfig{
			Name: "origin",
			URLs: []string{"https://example.com/user/repo.git"},
		})
		require.NoError(t, err)

		gitremotesPath := filepath.Join(path, ".gitremotes")
		require.NoError(t, appendToGitremotes(gitremotesPath, "origin", "https://example.com/user/repo.git"))

		require.NoError(t, removeRemote(path, "origin"))

		entries, err := parseGitremotes(gitremotesPath)
		require.NoError(t, err)
		assert.Empty(t, entries)

		cfg, err := repo.Config()
		require.NoError(t, err)
		_, stillRegistered := cfg.Remotes["origin"]
		assert.False(t, stillRegistered)
	})

	t.Run("returns errRemoteNotTracked for a name not present", func(t *testing.T) {
		path := t.TempDir()
		_, err := git.PlainInit(path, false)
		require.NoError(t, err)

		err = removeRemote(path, "nonexistent")
		assert.ErrorIs(t, err, errRemoteNotTracked)
	})

	t.Run("returns errRemoteNotTracked when .gitremotes is missing", func(t *testing.T) {
		path := t.TempDir()
		_, err := git.PlainInit(path, false)
		require.NoError(t, err)

		err = removeRemote(path, "origin")
		assert.ErrorIs(t, err, errRemoteNotTracked)
	})
}

func TestListRemotes(t *testing.T) {
	t.Run("reports registered vs unregistered remotes", func(t *testing.T) {
		path := t.TempDir()
		repo, err := git.PlainInit(path, false)
		require.NoError(t, err)

		_, err = repo.CreateRemote(&gitconfig.RemoteConfig{
			Name: "origin",
			URLs: []string{"https://example.com/user/repo.git"},
		})
		require.NoError(t, err)

		gitremotesPath := filepath.Join(path, ".gitremotes")
		require.NoError(t, appendToGitremotes(gitremotesPath, "origin", "https://example.com/user/repo.git"))
		require.NoError(t, appendToGitremotes(gitremotesPath, "mirror", "git@example.com:user/mirror.git"))

		statuses, err := listRemotes(path)
		require.NoError(t, err)
		require.Len(t, statuses, 2)

		byName := make(map[string]remoteStatus, len(statuses))
		for _, s := range statuses {
			byName[s.name] = s
		}

		assert.False(t, byName["origin"].unregistered)
		assert.True(t, byName["mirror"].unregistered)
	})

	t.Run("returns nil when .gitremotes is missing", func(t *testing.T) {
		path := t.TempDir()
		_, err := git.PlainInit(path, false)
		require.NoError(t, err)

		statuses, err := listRemotes(path)
		require.NoError(t, err)
		assert.Nil(t, statuses)
	})

	t.Run("returns nil when .gitremotes is empty", func(t *testing.T) {
		path := t.TempDir()
		_, err := git.PlainInit(path, false)
		require.NoError(t, err)

		require.NoError(t, os.WriteFile(filepath.Join(path, ".gitremotes"), nil, 0644))

		statuses, err := listRemotes(path)
		require.NoError(t, err)
		assert.Nil(t, statuses)
	})
}

// TestRemoteAddListRemoveRoundTrip verifies that add -> list -> remove leaves
// no trace of the remote in either .gitremotes or .git/config.
func TestRemoteAddListRemoveRoundTrip(t *testing.T) {
	path := t.TempDir()
	_, err := git.PlainInit(path, false)
	require.NoError(t, err)

	const name = "mirror"
	const url = "https://example.com/user/mirror.git"

	alreadyTracked, err := addRemote(path, name, url)
	require.NoError(t, err)
	require.False(t, alreadyTracked)

	statuses, err := listRemotes(path)
	require.NoError(t, err)
	require.Len(t, statuses, 1)
	assert.Equal(t, name, statuses[0].name)
	assert.Equal(t, url, statuses[0].url)
	assert.False(t, statuses[0].unregistered)

	require.NoError(t, removeRemote(path, name))

	statuses, err = listRemotes(path)
	require.NoError(t, err)
	assert.Empty(t, statuses)

	entries, err := parseGitremotes(filepath.Join(path, ".gitremotes"))
	require.NoError(t, err)
	assert.Empty(t, entries)

	repo, err := git.PlainOpen(path)
	require.NoError(t, err)
	cfg, err := repo.Config()
	require.NoError(t, err)
	_, stillRegistered := cfg.Remotes[name]
	assert.False(t, stillRegistered)
}
