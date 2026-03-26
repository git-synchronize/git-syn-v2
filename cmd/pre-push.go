/*
Copyright © 2024 Lucas Ramage <lucas.ramage@infinite-omicron.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

type pushResult struct {
	name string
	err  error
}

func pushRemote(repo *git.Repository, r remoteEntry) error {
	remote, err := repo.Remote(r.name)
	if err != nil {
		return fmt.Errorf("remote %q not found in .git/config (run git-syn install): %w", r.name, err)
	}

	err = remote.Push(&git.PushOptions{
		RemoteName: r.name,
	})
	if err == git.NoErrAlreadyUpToDate {
		return nil
	}
	return err
}

var prePushPath string

var prePushCmd = &cobra.Command{
	Use:   "pre-push",
	Short: "push to all remotes in .gitremotes",
	Long: `Push to all remotes tracked by git-syn in .gitremotes.

This command is intended to be invoked by the git pre-push hook installed by
git-syn install. It reads remotes from .gitremotes and pushes to each one,
reporting per-remote success or failure.`,
	Run: func(cmd *cobra.Command, args []string) {
		path := prePushPath
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
			fmt.Fprintln(os.Stderr, "error: .gitremotes not found — run git-syn install to initialize this repository")
			os.Exit(1)
		}
		if err != nil {
			log.Fatalf("failed to read .gitremotes: %v", err)
		}

		repo, err := git.PlainOpen(path)
		if err != nil {
			log.Fatalf("%s is not a git repository", path)
		}

		if ActiveConfig.PushStrategy == "sequential" {
			runSequential(repo, entries)
		} else {
			runParallel(repo, entries)
		}
	},
}

func runParallel(repo *git.Repository, entries []remoteEntry) {
	results := make(chan pushResult, len(entries))
	var wg sync.WaitGroup

	for _, entry := range entries {
		wg.Add(1)
		go func(r remoteEntry) {
			defer wg.Done()
			err := pushRemote(repo, r)
			results <- pushResult{name: r.name, err: err}
		}(entry)
	}

	wg.Wait()
	close(results)

	handleResults(results)
}

func runSequential(repo *git.Repository, entries []remoteEntry) {
	results := make(chan pushResult, len(entries))

	for _, entry := range entries {
		err := pushRemote(repo, entry)
		results <- pushResult{name: entry.name, err: err}
		if err != nil && ActiveConfig.OnFailure == "abort" {
			close(results)
			handleResults(results)
			return
		}
	}

	close(results)
	handleResults(results)
}

func handleResults(results <-chan pushResult) {
	var succeeded, failed int
	for r := range results {
		if r.err != nil {
			fmt.Printf("✘ %s: %v\n", r.name, r.err)
			failed++
		} else {
			fmt.Printf("✔ %s\n", r.name)
			succeeded++
		}
	}

	if ActiveConfig.OnFailure == "abort" && failed > 0 {
		os.Exit(1)
	}
	if succeeded == 0 && failed > 0 {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(prePushCmd)
	prePushCmd.Flags().StringVar(&prePushPath, "path", "", "path to the git repository (default: current directory)")
}
