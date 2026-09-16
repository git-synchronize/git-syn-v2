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

// setupInstalledRepo creates a repo with a registered remote, a matching
// .gitremotes entry, and an installed pre-push hook, simulating the state
// left behind by a prior `git-syn install`.
func setupInstalledRepo(t *testing.T) (path string, repo *git.Repository) {
	t.Helper()

	path = t.TempDir()
	var err error
	repo, err = git.PlainInit(path, false)
	require.NoError(t, err)

	_, err = repo.CreateRemote(&gitconfig.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://example.com/user/repo.git"},
	})
	require.NoError(t, err)

	require.NoError(t, appendToGitremotes(filepath.Join(path, ".gitremotes"), "origin", "https://example.com/user/repo.git"))

	hooksDir := filepath.Join(path, ".git", "hooks")
	require.NoError(t, os.MkdirAll(hooksDir, 0755))
	require.NoError(t, install_hook("pre-push", hooksDir))

	return path, repo
}

func TestUninstallRepo(t *testing.T) {
	t.Run("no flags leaves .gitremotes and git config intact", func(t *testing.T) {
		path, repo := setupInstalledRepo(t)

		message, err := uninstallRepo(path, false)
		require.NoError(t, err)
		assert.Equal(t, "Removed git hooks. Git SYN uninstalled. Run with --clean to also remove .gitremotes and remote config entries.", message)

		_, err = os.Stat(filepath.Join(path, ".git", "hooks", "pre-push"))
		assert.True(t, os.IsNotExist(err))

		entries, err := parseGitremotes(filepath.Join(path, ".gitremotes"))
		require.NoError(t, err)
		assert.Len(t, entries, 1)

		cfg, err := repo.Config()
		require.NoError(t, err)
		_, registered := cfg.Remotes["origin"]
		assert.True(t, registered)
	})

	t.Run("--clean removes .gitremotes and git config entries", func(t *testing.T) {
		path, repo := setupInstalledRepo(t)

		message, err := uninstallRepo(path, true)
		require.NoError(t, err)
		assert.Equal(t, "Removed git hooks and remotes. Git SYN fully uninstalled.", message)

		_, err = os.Stat(filepath.Join(path, ".git", "hooks", "pre-push"))
		assert.True(t, os.IsNotExist(err))

		_, err = os.Stat(filepath.Join(path, ".gitremotes"))
		assert.True(t, os.IsNotExist(err))

		cfg, err := repo.Config()
		require.NoError(t, err)
		_, registered := cfg.Remotes["origin"]
		assert.False(t, registered)
	})

	t.Run("--clean succeeds when .gitremotes is already absent", func(t *testing.T) {
		path := t.TempDir()
		_, err := git.PlainInit(path, false)
		require.NoError(t, err)

		hooksDir := filepath.Join(path, ".git", "hooks")
		require.NoError(t, os.MkdirAll(hooksDir, 0755))
		require.NoError(t, install_hook("pre-push", hooksDir))

		message, err := uninstallRepo(path, true)
		require.NoError(t, err)
		assert.Equal(t, "Removed git hooks and remotes. Git SYN fully uninstalled.", message)
	})
}
