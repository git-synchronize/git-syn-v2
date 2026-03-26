/*
Copyright © 2024-2026 Lucas Ramage <lucas.ramage@infinite-omicron.com>
*/
package cmd

import (
	"fmt"
	"log"
	"log/slog"
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
	slog.Info("pushing to remote", "name", r.name)
	remote, err := repo.Remote(r.name)
	if err != nil {
		return fmt.Errorf("remote %q not found in .git/config (run git-syn install): %w", r.name, err)
	}

	err = remote.Push(&git.PushOptions{
		RemoteName: r.name,
	})
	if err == git.NoErrAlreadyUpToDate {
		slog.Debug("remote already up to date", "name", r.name)
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
		slog.Debug("using repository path", "path", path)

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

		slog.Debug("starting sync", "strategy", ActiveConfig.PushStrategy, "remote_count", len(entries))
		os.Exit(syncAll(repo, entries))
	},
}

func syncAll(repo *git.Repository, entries []remoteEntry) int {
	if ActiveConfig.PushStrategy == "sequential" {
		return runSequential(repo, entries)
	}
	return runParallel(repo, entries)
}

func runParallel(repo *git.Repository, entries []remoteEntry) int {
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

	return handleResults(results)
}

func runSequential(repo *git.Repository, entries []remoteEntry) int {
	results := make(chan pushResult, len(entries))

	for _, entry := range entries {
		err := pushRemote(repo, entry)
		results <- pushResult{name: entry.name, err: err}
		if err != nil && ActiveConfig.OnFailure == "abort" {
			close(results)
			return handleResults(results)
		}
	}

	close(results)
	return handleResults(results)
}

func handleResults(results <-chan pushResult) int {
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
		return 1
	}
	if succeeded == 0 && failed > 0 {
		return 1
	}
	return 0
}

func init() {
	rootCmd.AddCommand(prePushCmd)
	prePushCmd.Flags().StringVar(&prePushPath, "path", "", "path to the git repository (default: current directory)")
}
