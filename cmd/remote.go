/*
Copyright © 2024-2026 Lucas Ramage <lucas.ramage@infinite-omicron.com>
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

// removeFromGitremotes rewrites .gitremotes without the named remote stanza.
func removeFromGitremotes(path, name string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	header := fmt.Sprintf("[remote %q]", name)
	lines := strings.Split(string(data), "\n")
	var result []string
	skip := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == header {
			skip = true
			continue
		}
		if skip && strings.HasPrefix(trimmed, "[") {
			skip = false
		}
		if !skip {
			result = append(result, line)
		}
	}

	return os.WriteFile(path, []byte(strings.Join(result, "\n")), 0644)
}

var remotePath string

var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "manage tracked remotes",
}

var remoteAddCmd = &cobra.Command{
	Use:   "add <name> <url>",
	Short: "add a remote to .gitremotes",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		name, url := args[0], args[1]

		if err := validateRemoteURL(url, ActiveConfig.AllowedSchemes); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		path := remotePath
		if path == "" {
			var err error
			path, err = os.Getwd()
			if err != nil {
				log.Fatalf("failed to get current directory: %v", err)
			}
		}

		gitremotesPath := filepath.Join(path, ".gitremotes")
		existing, err := parseGitremotes(gitremotesPath)
		if err != nil && !os.IsNotExist(err) {
			log.Fatalf("failed to read .gitremotes: %v", err)
		}

		for _, e := range existing {
			if e.name == name {
				fmt.Printf("Remote '%s' is already tracked by git-syn.\n", name)
				return
			}
		}

		if err := appendToGitremotes(gitremotesPath, name, url); err != nil {
			log.Fatalf("failed to write .gitremotes: %v", err)
		}

		repo, err := git.PlainOpen(path)
		if err != nil {
			log.Fatalf("%s is not a git repository", path)
		}

		if err := registerRemoteInConfig(repo, name, url); err != nil {
			log.Fatalf("failed to register remote in .git/config: %v", err)
		}

		fmt.Printf("Added remote '%s'.\n", name)
	},
}

var remoteRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "remove a remote from .gitremotes",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		path := remotePath
		if path == "" {
			var err error
			path, err = os.Getwd()
			if err != nil {
				log.Fatalf("failed to get current directory: %v", err)
			}
		}

		gitremotesPath := filepath.Join(path, ".gitremotes")
		entries, err := parseGitremotes(gitremotesPath)
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "error: remote '%s' not found in .gitremotes\n", name)
			os.Exit(1)
		}
		if err != nil {
			log.Fatalf("failed to read .gitremotes: %v", err)
		}

		found := false
		for _, e := range entries {
			if e.name == name {
				found = true
				break
			}
		}
		if !found {
			fmt.Fprintf(os.Stderr, "error: remote '%s' not found in .gitremotes\n", name)
			os.Exit(1)
		}

		if err := removeFromGitremotes(gitremotesPath, name); err != nil {
			log.Fatalf("failed to update .gitremotes: %v", err)
		}

		repo, err := git.PlainOpen(path)
		if err != nil {
			log.Fatalf("%s is not a git repository", path)
		}

		if err := repo.DeleteRemote(name); err != nil && err != git.ErrRemoteNotFound {
			log.Fatalf("failed to remove remote from .git/config: %v", err)
		}

		fmt.Printf("Removed remote '%s'.\n", name)
	},
}

var remoteListCmd = &cobra.Command{
	Use:   "list",
	Short: "list tracked remotes",
	Run: func(cmd *cobra.Command, args []string) {
		path := remotePath
		if path == "" {
			var err error
			path, err = os.Getwd()
			if err != nil {
				log.Fatalf("failed to get current directory: %v", err)
			}
		}

		gitremotesPath := filepath.Join(path, ".gitremotes")
		entries, err := parseGitremotes(gitremotesPath)
		if os.IsNotExist(err) {
			fmt.Println("No remotes tracked by git-syn.")
			return
		}
		if err != nil {
			log.Fatalf("failed to read .gitremotes: %v", err)
		}

		if len(entries) == 0 {
			fmt.Println("No remotes tracked by git-syn.")
			return
		}

		repo, err := git.PlainOpen(path)
		if err != nil {
			log.Fatalf("%s is not a git repository", path)
		}

		cfg, err := repo.Config()
		if err != nil {
			log.Fatalf("failed to read .git/config: %v", err)
		}

		for _, entry := range entries {
			if _, registered := cfg.Remotes[entry.name]; !registered {
				fmt.Printf("%s %s [unregistered]\n", entry.name, entry.url)
			} else {
				fmt.Printf("%s %s\n", entry.name, entry.url)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(remoteCmd)
	remoteCmd.PersistentFlags().StringVar(&remotePath, "path", "", "path to the git repository (default: current directory)")
	remoteCmd.AddCommand(remoteAddCmd)
	remoteCmd.AddCommand(remoteRemoveCmd)
	remoteCmd.AddCommand(remoteListCmd)
}
