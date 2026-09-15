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

func TestInitRepo(t *testing.T) {
	t.Run("registers remotes with supported schemes", func(t *testing.T) {
		path := t.TempDir()
		repo, err := git.PlainInit(path, false)
		require.NoError(t, err)

		_, err = repo.CreateRemote(&gitconfig.RemoteConfig{
			Name: "origin",
			URLs: []string{"https://example.com/user/repo.git"},
		})
		require.NoError(t, err)

		require.NoError(t, init_repo(path))

		entries, err := parseGitremotes(filepath.Join(path, ".gitremotes"))
		require.NoError(t, err)
		require.Len(t, entries, 1)
		assert.Equal(t, "origin", entries[0].name)
		assert.Equal(t, "https://example.com/user/repo.git", entries[0].url)
	})

	t.Run("skips remotes with unsupported schemes", func(t *testing.T) {
		path := t.TempDir()
		repo, err := git.PlainInit(path, false)
		require.NoError(t, err)

		_, err = repo.CreateRemote(&gitconfig.RemoteConfig{
			Name: "ftp-remote",
			URLs: []string{"ftp://example.com/user/repo.git"},
		})
		require.NoError(t, err)

		require.NoError(t, init_repo(path))

		entries, err := parseGitremotes(filepath.Join(path, ".gitremotes"))
		require.NoError(t, err)
		assert.Empty(t, entries)
	})

	t.Run("non-git-repository path returns an error", func(t *testing.T) {
		path := t.TempDir()

		err := init_repo(path)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "is not a git repository")
	})
}

func TestCleanGitremotes(t *testing.T) {
	t.Run("removes registered remotes and the .gitremotes file", func(t *testing.T) {
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

		require.NoError(t, cleanGitremotes(path))

		_, err = os.Stat(gitremotesPath)
		assert.True(t, os.IsNotExist(err))

		cfg, err := repo.Config()
		require.NoError(t, err)
		_, stillRegistered := cfg.Remotes["origin"]
		assert.False(t, stillRegistered)
	})

	t.Run("tolerates a remote already absent from git config", func(t *testing.T) {
		path := t.TempDir()
		_, err := git.PlainInit(path, false)
		require.NoError(t, err)

		gitremotesPath := filepath.Join(path, ".gitremotes")
		require.NoError(t, appendToGitremotes(gitremotesPath, "ghost", "https://example.com/user/repo.git"))

		assert.NoError(t, cleanGitremotes(path))
		_, err = os.Stat(gitremotesPath)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("missing .gitremotes is a no-op", func(t *testing.T) {
		path := t.TempDir()
		_, err := git.PlainInit(path, false)
		require.NoError(t, err)

		assert.NoError(t, cleanGitremotes(path))
	})
}
