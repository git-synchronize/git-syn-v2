/*
Copyright © 2024 Lucas Ramage <lucas.ramage@infinite-omicron.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

const hookContent = "#!/bin/sh\ncommand -v git-syn >/dev/null 2>&1 || { printf >&2 \"\\n%s\\n\\n\" \"This repository is configured for Git SYN but 'git-syn' was not found on your path. If you no longer wish to use Git SYN, remove this hook by deleting the '{{Command}}' file in the hooks directory (usually '.git/hooks').\"; exit 2; }\ngit syn {{Command}} \"$@\""

func write_hook(path, content string) error {
	return os.WriteFile(path, []byte(content+"\n"), 0600)
}

func set_file_permissions(path string, mode os.FileMode) error {
	return os.Chmod(path, mode)
}

func install_hook(hookType, hooksDir string) error {
	content := strings.ReplaceAll(hookContent, "{{Command}}", hookType)
	dst := filepath.Join(hooksDir, hookType)
	if err := write_hook(dst, content); err != nil {
		return err
	}
	return set_file_permissions(dst, 0700)
}

func init_repo(path string) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		log.Fatalf("%s is not a git repository", path)
	}

	remotes, err := repo.Remotes()
	if err != nil {
		log.Fatalf("failed to list remotes: %v", err)
	}

	gitremotesPath := filepath.Join(path, ".gitremotes")
	f, err := os.Create(gitremotesPath)
	if err != nil {
		log.Fatalf("failed to create .gitremotes: %v", err)
	}
	f.Close()

	for _, remote := range remotes {
		rc := remote.Config()
		if len(rc.URLs) == 0 {
			continue
		}
		url := rc.URLs[0]
		if err := validateRemoteURL(url); err != nil {
			fmt.Fprintf(os.Stderr, "warning: skipping remote %q with unsupported URL scheme: %s\n", rc.Name, url)
			continue
		}
		if err := appendToGitremotes(gitremotesPath, rc.Name, url); err != nil {
			log.Fatalf("failed to write .gitremotes: %v", err)
		}
	}

	entries, err := parseGitremotes(gitremotesPath)
	if err != nil {
		log.Fatalf("failed to read .gitremotes: %v", err)
	}

	for _, entry := range entries {
		if err := registerRemoteInConfig(repo, entry.name, entry.url); err != nil {
			log.Fatalf("failed to register remote %q in .git/config: %v", entry.name, err)
		}
	}
}

func cleanGitremotes(repoPath string) error {
	gitremotesPath := filepath.Join(repoPath, ".gitremotes")

	entries, err := parseGitremotes(gitremotesPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to read .gitremotes: %w", err)
	}

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return fmt.Errorf("%s is not a git repository", repoPath)
	}

	for _, entry := range entries {
		if err := repo.DeleteRemote(entry.name); err != nil && err != git.ErrRemoteNotFound {
			return fmt.Errorf("failed to remove remote %q from .git/config: %w", entry.name, err)
		}
	}

	if err := os.Remove(gitremotesPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove .gitremotes: %w", err)
	}

	return nil
}

var installPath string

var installCmd = &cobra.Command{
	Use:   "install [--path <dir>]",
	Short: "install extension to repository",
	Long: `Install git-syn into a git repository.

Installs a pre-push hook into .git/hooks/ and writes a .gitremotes file
containing all configured HTTPS and OpenSSH remotes. The remotes are also
registered in .git/config for synchronization.

Git SYN enables remote repository synchronization across multiple git forges
for disaster recovery and censorship resistance.`,
	Run: func(cmd *cobra.Command, args []string) {
		path := installPath
		if path == "" {
			var err error
			path, err = os.Getwd()
			if err != nil {
				log.Fatalf("failed to get current directory: %v", err)
			}
		}

		hooksDir := filepath.Join(path, ".git", "hooks")
		if err := install_hook("pre-push", hooksDir); err != nil {
			log.Fatalf("failed to install pre-push hook: %v", err)
		}

		init_repo(path)

		fmt.Println("Updated git hooks. Git SYN initialized.")
	},
}

var uninstallPath string
var uninstallClean bool

var uninstallCmd = &cobra.Command{
	Use:   "uninstall [--path <dir>]",
	Short: "remove extension from repository",
	Long: `Remove git-syn from a git repository.

Removes the pre-push hook from .git/hooks/ that was installed by git-syn.`,
	Run: func(cmd *cobra.Command, args []string) {
		path := uninstallPath
		if path == "" {
			var err error
			path, err = os.Getwd()
			if err != nil {
				log.Fatalf("failed to get current directory: %v", err)
			}
		}

		hookPath := filepath.Join(path, ".git", "hooks", "pre-push")
		if err := os.Remove(hookPath); err != nil && !os.IsNotExist(err) {
			log.Fatalf("failed to remove pre-push hook: %v", err)
		}

		if uninstallClean {
			if err := cleanGitremotes(path); err != nil {
				log.Fatalf("failed to clean up remotes: %v", err)
			}
			fmt.Println("Removed git hooks and remotes. Git SYN fully uninstalled.")
		} else {
			fmt.Println("Removed git hooks. Git SYN uninstalled. Run with --clean to also remove .gitremotes and remote config entries.")
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().StringVar(&installPath, "path", "", "path to the git repository (default: current directory)")

	rootCmd.AddCommand(uninstallCmd)
	uninstallCmd.Flags().StringVar(&uninstallPath, "path", "", "path to the git repository (default: current directory)")
	uninstallCmd.Flags().BoolVar(&uninstallClean, "clean", false, "also remove .gitremotes and git config entries")
}
