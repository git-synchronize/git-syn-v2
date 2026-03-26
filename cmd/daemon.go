/*
Copyright © 2024-2026 Lucas Ramage <lucas.ramage@infinite-omicron.com>
*/
package cmd

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

var (
	daemonPath     string
	daemonInterval time.Duration
	daemonOnce     bool
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "run synchronization on a schedule",
	Long: `Run the synchronization process periodically.

The daemon will push to all remotes in .gitremotes at the specified interval.
If --once is specified, it will run exactly once and exit.`,
	Run: func(cmd *cobra.Command, args []string) {
		path := daemonPath
		if path == "" {
			var err error
			path, err = os.Getwd()
			if err != nil {
				log.Fatalf("failed to get current directory: %v", err)
			}
		}

		interval := daemonInterval
		if interval == 0 {
			interval = ActiveConfig.SyncInterval
		}

		if interval == 0 && !daemonOnce {
			log.Fatal("sync interval must be greater than zero")
		}

		gitremotesPath := filepath.Join(path, ".gitremotes")
		repo, err := git.PlainOpen(path)
		if err != nil {
			log.Fatalf("%s is not a git repository", path)
		}

		if daemonOnce {
			entries, err := parseGitremotes(gitremotesPath)
			if err != nil {
				log.Fatalf("failed to read .gitremotes: %v", err)
			}
			os.Exit(syncAll(repo, entries))
		}

		slog.Info("starting daemon", "path", path, "interval", interval)

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			entries, err := parseGitremotes(gitremotesPath)
			if err != nil {
				slog.Error("failed to read .gitremotes", "error", err)
			} else {
				syncAll(repo, entries)
			}

			select {
			case <-ticker.C:
				continue
			case <-ctx.Done():
				slog.Info("shutting down daemon")
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)
	daemonCmd.Flags().StringVar(&daemonPath, "path", "", "path to the git repository (default: current directory)")
	daemonCmd.Flags().DurationVar(&daemonInterval, "interval", 0, "sync interval (e.g., 5m, 1h) (default: from config)")
	daemonCmd.Flags().BoolVar(&daemonOnce, "once", false, "run once and exit")
}
