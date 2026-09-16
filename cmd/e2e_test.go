/*
Copyright © 2024-2026 Lucas Ramage <lucas.ramage@infinite-omicron.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	gitSynBinaryOnce sync.Once
	gitSynBinaryPath string
	gitSynBinaryErr  error
)

// buildGitSynBinary compiles the git-syn binary once per test run and
// returns its path. The binary is named exactly "git-syn" so that, once its
// directory is on PATH, both the pre-push hook script's "command -v
// git-syn" check and git's own "git syn <verb>" -> "git-syn <verb>"
// external-command dispatch resolve it correctly.
func buildGitSynBinary(t *testing.T) string {
	t.Helper()

	gitSynBinaryOnce.Do(func() {
		dir, err := os.MkdirTemp("", "git-syn-bin-*")
		if err != nil {
			gitSynBinaryErr = err
			return
		}
		binPath := filepath.Join(dir, "git-syn")
		buildCmd := exec.Command("go", "build", "-o", binPath, "gitlab.com/git-syn/git-syn")
		out, err := buildCmd.CombinedOutput()
		if err != nil {
			gitSynBinaryErr = fmt.Errorf("go build git-syn failed: %w\n%s", err, out)
			return
		}
		gitSynBinaryPath = binPath
	})

	require.NoError(t, gitSynBinaryErr)
	return gitSynBinaryPath
}

// TestPrePushHookDelegation verifies that git-syn install writes a working
// pre-push hook, and that a real `git push` on the installed repository
// delegates to `git-syn pre-push`, which pushes to every entry in
// .gitremotes — including remotes other than the one the push was aimed at,
// matching the tool's real mirroring use case (push to your primary, the
// hook fans out to the mirrors).
func TestPrePushHookDelegation(t *testing.T) {
	binPath := buildGitSynBinary(t)
	binDir := filepath.Dir(binPath)
	pathWithBin := binDir + string(os.PathListSeparator) + os.Getenv("PATH")

	_, repoDir, mirrorName, mirrorPath := newTestRepoPair(t)

	// A second bare repo the outer `git push` targets directly, distinct
	// from the mirror git-syn tracks in .gitremotes. Pushing to the same
	// remote both directly and via the hook races two independent pushes
	// against one ref, which isn't the tool's intended usage.
	primaryPath := t.TempDir()
	_, err := git.PlainInit(primaryPath, true)
	require.NoError(t, err)

	install := exec.Command(binPath, "install", "--path", repoDir)
	install.Env = append(os.Environ(), "PATH="+pathWithBin)
	out, err := install.CombinedOutput()
	require.NoError(t, err, "git-syn install failed: %s", out)

	hookPath := filepath.Join(repoDir, ".git", "hooks", "pre-push")
	info, err := os.Stat(hookPath)
	require.NoError(t, err, "expected install to write a pre-push hook")
	assert.NotZero(t, info.Mode()&0100, "hook should be executable")

	// install's own remote-discovery skips the bare-repo path URL (it isn't
	// an accepted https/ssh scheme), so .gitremotes comes out empty. Seed it
	// directly with the fixture mirror, exactly as it would look had the
	// user tracked it via a scheme install/remote-add would accept.
	gitremotesPath := filepath.Join(repoDir, ".gitremotes")
	require.NoError(t, appendToGitremotes(gitremotesPath, mirrorName, mirrorPath))

	addPrimary := exec.Command("git", "remote", "add", "primary", primaryPath)
	addPrimary.Dir = repoDir
	out, err = addPrimary.CombinedOutput()
	require.NoError(t, err, "git remote add primary failed: %s", out)

	push := exec.Command("git", "push", "primary", "master")
	push.Dir = repoDir
	push.Env = append(os.Environ(), "PATH="+pathWithBin)
	out, err = push.CombinedOutput()
	require.NoError(t, err, "git push failed: %s", out)
	assert.Contains(t, string(out), "✔ "+mirrorName, "pre-push hook output should show the delegated push to the mirror succeeding")

	primaryRepo, err := git.PlainOpen(primaryPath)
	require.NoError(t, err)
	primaryHead, err := primaryRepo.Head()
	require.NoError(t, err)
	assert.False(t, primaryHead.Hash().IsZero(), "the direct git push should have landed on primary")

	mirrorRepo, err := git.PlainOpen(mirrorPath)
	require.NoError(t, err)
	mirrorHead, err := mirrorRepo.Head()
	require.NoError(t, err)
	assert.False(t, mirrorHead.Hash().IsZero(), "the pre-push hook should have delegated the push to the mirror too")
	assert.Equal(t, primaryHead.Hash(), mirrorHead.Hash(), "primary and mirror should end up at the same commit")
}

// TestInvalidYAMLConfigExitsNonZero verifies that a malformed config.yaml
// causes git-syn to print an error and exit non-zero at startup, before any
// command runs.
func TestInvalidYAMLConfigExitsNonZero(t *testing.T) {
	binPath := buildGitSynBinary(t)

	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tmpHome, ".config"))

	configDir, err := os.UserConfigDir()
	require.NoError(t, err)
	gitSynConfigDir := filepath.Join(configDir, "git-syn")
	require.NoError(t, os.MkdirAll(gitSynConfigDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(gitSynConfigDir, "config.yaml"), []byte("on_failure: [\n"), 0644))

	run := exec.Command(binPath, "--help")
	run.Env = os.Environ()
	out, err := run.CombinedOutput()

	require.Error(t, err, "expected git-syn to exit non-zero on invalid config YAML")
	var exitErr *exec.ExitError
	require.True(t, errors.As(err, &exitErr))
	assert.Equal(t, 1, exitErr.ExitCode())
	assert.Contains(t, string(out), "error:")
}
