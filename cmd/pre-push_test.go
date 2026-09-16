/*
Copyright © 2024-2026 Lucas Ramage <lucas.ramage@infinite-omicron.com>
*/
package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestRepoPair creates a local repository with one commit and a bare
// repository registered as its remote, entirely on the local filesystem
// (go-git's "file" transport needs no network access). It returns the local
// repository, its filesystem path, and the name under which the bare repo
// is registered as a remote.
func newTestRepoPair(t *testing.T) (repo *git.Repository, repoPath, remoteName, remotePath string) {
	t.Helper()

	remoteName = "origin"
	remotePath = t.TempDir()
	_, err := git.PlainInit(remotePath, true)
	require.NoError(t, err)

	repoPath = t.TempDir()
	repo, err = git.PlainInit(repoPath, false)
	require.NoError(t, err)

	wt, err := repo.Worktree()
	require.NoError(t, err)

	require.NoError(t, os.WriteFile(filepath.Join(repoPath, "README.md"), []byte("hello\n"), 0644))
	_, err = wt.Add("README.md")
	require.NoError(t, err)

	_, err = wt.Commit("initial commit", &git.CommitOptions{
		Author: &object.Signature{Name: "Test", Email: "test@example.com", When: time.Now()},
	})
	require.NoError(t, err)

	_, err = repo.CreateRemote(&gitconfig.RemoteConfig{
		Name: remoteName,
		URLs: []string{remotePath},
	})
	require.NoError(t, err)

	return repo, repoPath, remoteName, remotePath
}

// captureStdout runs fn with os.Stdout redirected to a pipe and returns
// whatever fn wrote to it.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	old := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)
	os.Stdout = w

	fn()

	require.NoError(t, w.Close())
	os.Stdout = old

	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	require.NoError(t, err)
	return buf.String()
}

// withConfig temporarily replaces ActiveConfig for the duration of the test.
func withConfig(t *testing.T, cfg *Config) {
	t.Helper()
	orig := ActiveConfig
	ActiveConfig = cfg
	t.Cleanup(func() { ActiveConfig = orig })
}

func TestPushRemote(t *testing.T) {
	t.Run("successful push", func(t *testing.T) {
		repo, _, remoteName, remotePath := newTestRepoPair(t)

		err := pushRemote(repo, remoteEntry{name: remoteName})
		require.NoError(t, err)

		bareRepo, err := git.PlainOpen(remotePath)
		require.NoError(t, err)
		head, err := bareRepo.Head()
		require.NoError(t, err)
		assert.False(t, head.Hash().IsZero())
	})

	t.Run("already up to date", func(t *testing.T) {
		repo, _, remoteName, _ := newTestRepoPair(t)
		require.NoError(t, pushRemote(repo, remoteEntry{name: remoteName}))

		err := pushRemote(repo, remoteEntry{name: remoteName})
		assert.NoError(t, err)
	})

	t.Run("remote not registered", func(t *testing.T) {
		repo, _, _, _ := newTestRepoPair(t)

		err := pushRemote(repo, remoteEntry{name: "does-not-exist"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "does-not-exist")
		assert.Contains(t, err.Error(), "git-syn install")
	})
}

func TestRunParallel(t *testing.T) {
	withConfig(t, &Config{OnFailure: "warn", PushStrategy: "parallel", AllowedSchemes: defaultAllowedSchemes})

	repo, _, goodRemote, _ := newTestRepoPair(t)
	entries := []remoteEntry{
		{name: "missing-remote"},
		{name: goodRemote},
	}

	var code int
	output := captureStdout(t, func() {
		code = runParallel(repo, entries)
	})

	assert.Contains(t, output, "✔ "+goodRemote)
	assert.Contains(t, output, "✘ missing-remote")
	assert.Equal(t, 0, code)
}

func TestRunSequential(t *testing.T) {
	t.Run("abort stops before later remotes", func(t *testing.T) {
		withConfig(t, &Config{OnFailure: "abort", PushStrategy: "sequential", AllowedSchemes: defaultAllowedSchemes})

		repo, _, goodRemote, _ := newTestRepoPair(t)
		entries := []remoteEntry{
			{name: "missing-remote"},
			{name: goodRemote},
		}

		var code int
		output := captureStdout(t, func() {
			code = runSequential(repo, entries)
		})

		assert.Contains(t, output, "✘ missing-remote")
		assert.NotContains(t, output, goodRemote)
		assert.Equal(t, 1, code)
	})

	t.Run("warn continues through all remotes", func(t *testing.T) {
		withConfig(t, &Config{OnFailure: "warn", PushStrategy: "sequential", AllowedSchemes: defaultAllowedSchemes})

		repo, _, goodRemote, _ := newTestRepoPair(t)
		entries := []remoteEntry{
			{name: "missing-remote"},
			{name: goodRemote},
		}

		var code int
		output := captureStdout(t, func() {
			code = runSequential(repo, entries)
		})

		assert.Contains(t, output, "✘ missing-remote")
		assert.Contains(t, output, "✔ "+goodRemote)
		assert.Equal(t, 0, code)
	})
}

func TestSyncAll(t *testing.T) {
	t.Run("parallel strategy attempts every remote even under abort", func(t *testing.T) {
		withConfig(t, &Config{OnFailure: "abort", PushStrategy: "parallel", AllowedSchemes: defaultAllowedSchemes})

		repo, _, goodRemote, _ := newTestRepoPair(t)
		entries := []remoteEntry{
			{name: "missing-remote"},
			{name: goodRemote},
		}

		output := captureStdout(t, func() {
			syncAll(repo, entries)
		})

		assert.Contains(t, output, "✔ "+goodRemote)
	})

	t.Run("sequential strategy stops before later remotes under abort", func(t *testing.T) {
		withConfig(t, &Config{OnFailure: "abort", PushStrategy: "sequential", AllowedSchemes: defaultAllowedSchemes})

		repo, _, goodRemote, _ := newTestRepoPair(t)
		entries := []remoteEntry{
			{name: "missing-remote"},
			{name: goodRemote},
		}

		output := captureStdout(t, func() {
			syncAll(repo, entries)
		})

		assert.NotContains(t, output, goodRemote)
	})
}
