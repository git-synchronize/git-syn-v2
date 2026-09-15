/*
Copyright © 2024-2026 Lucas Ramage <lucas.ramage@infinite-omicron.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

// errRemoteNotTracked indicates a remote name was not found in .gitremotes.
var errRemoteNotTracked = errors.New("remote not found in .gitremotes")

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

// addRemote records name/url in .gitremotes and registers it in .git/config.
// It reports alreadyTracked=true (and a nil error) if the remote is already
// tracked, without making any changes.
func addRemote(path, name, url string) (alreadyTracked bool, err error) {
	gitremotesPath := filepath.Join(path, ".gitremotes")
	existing, err := parseGitremotes(gitremotesPath)
	if err != nil && !os.IsNotExist(err) {
		return false, fmt.Errorf("failed to read .gitremotes: %w", err)
	}

	for _, e := range existing {
		if e.name == name {
			return true, nil
		}
	}

	if err := appendToGitremotes(gitremotesPath, name, url); err != nil {
		return false, fmt.Errorf("failed to write .gitremotes: %w", err)
	}

	repo, err := git.PlainOpen(path)
	if err != nil {
		return false, fmt.Errorf("%s is not a git repository", path)
	}

	if err := registerRemoteInConfig(repo, name, url); err != nil {
		return false, fmt.Errorf("failed to register remote in .git/config: %w", err)
	}

	return false, nil
}

// removeRemote drops name from .gitremotes and .git/config. It returns
// errRemoteNotTracked if name is not present in .gitremotes.
func removeRemote(path, name string) error {
	gitremotesPath := filepath.Join(path, ".gitremotes")
	entries, err := parseGitremotes(gitremotesPath)
	if os.IsNotExist(err) {
		return errRemoteNotTracked
	}
	if err != nil {
		return fmt.Errorf("failed to read .gitremotes: %w", err)
	}

	found := false
	for _, e := range entries {
		if e.name == name {
			found = true
			break
		}
	}
	if !found {
		return errRemoteNotTracked
	}

	if err := removeFromGitremotes(gitremotesPath, name); err != nil {
		return fmt.Errorf("failed to update .gitremotes: %w", err)
	}

	repo, err := git.PlainOpen(path)
	if err != nil {
		return fmt.Errorf("%s is not a git repository", path)
	}

	if err := repo.DeleteRemote(name); err != nil && err != git.ErrRemoteNotFound {
		return fmt.Errorf("failed to remove remote from .git/config: %w", err)
	}

	return nil
}

// remoteStatus describes a tracked remote and whether it is registered in
// .git/config.
type remoteStatus struct {
	name         string
	url          string
	unregistered bool
}

// listRemotes returns the remotes tracked in .gitremotes, annotated with
// whether each is currently registered in .git/config. A missing or empty
// .gitremotes file yields a nil slice and nil error.
func listRemotes(path string) ([]remoteStatus, error) {
	gitremotesPath := filepath.Join(path, ".gitremotes")
	entries, err := parseGitremotes(gitremotesPath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read .gitremotes: %w", err)
	}
	if len(entries) == 0 {
		return nil, nil
	}

	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("%s is not a git repository", path)
	}

	cfg, err := repo.Config()
	if err != nil {
		return nil, fmt.Errorf("failed to read .git/config: %w", err)
	}

	statuses := make([]remoteStatus, 0, len(entries))
	for _, entry := range entries {
		_, registered := cfg.Remotes[entry.name]
		statuses = append(statuses, remoteStatus{name: entry.name, url: entry.url, unregistered: !registered})
	}
	return statuses, nil
}

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

		alreadyTracked, err := addRemote(path, name, url)
		if err != nil {
			log.Fatalf("%v", err)
		}
		if alreadyTracked {
			fmt.Printf("Remote '%s' is already tracked by git-syn.\n", name)
			return
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

		if err := removeRemote(path, name); err != nil {
			if errors.Is(err, errRemoteNotTracked) {
				fmt.Fprintf(os.Stderr, "error: remote '%s' not found in .gitremotes\n", name)
				os.Exit(1)
			}
			log.Fatalf("%v", err)
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

		statuses, err := listRemotes(path)
		if err != nil {
			log.Fatalf("%v", err)
		}

		if len(statuses) == 0 {
			fmt.Println("No remotes tracked by git-syn.")
			return
		}

		for _, s := range statuses {
			if s.unregistered {
				fmt.Printf("%s %s [unregistered]\n", s.name, s.url)
			} else {
				fmt.Printf("%s %s\n", s.name, s.url)
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
