/*
Copyright © 2024-2026 Lucas Ramage <lucas.ramage@infinite-omicron.com>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/http/cgi"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

var (
	servePath   string
	serveListen string
)

// gitExecPath returns git's own reported exec-path, the directory
// containing its core binaries (including git-http-backend). Using git's
// own answer, rather than a hardcoded path like
// /usr/lib/git-core/git-http-backend, keeps this portable across
// distributions that install git-http-backend somewhere else (or, on some
// distributions, not at all).
func gitExecPath() (string, error) {
	out, err := exec.Command("git", "--exec-path").Output()
	if err != nil {
		return "", fmt.Errorf("failed to run git --exec-path: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// resolveGitHTTPBackend joins "git-http-backend" onto execPath and
// confirms it exists, so a missing binary fails fast at startup with a
// clear error instead of a mysterious 500 on the first request.
func resolveGitHTTPBackend(execPath string) (string, error) {
	backend := filepath.Join(execPath, "git-http-backend")
	if _, err := os.Stat(backend); err != nil {
		return "", fmt.Errorf("git-http-backend not found at %s: %w", backend, err)
	}
	return backend, nil
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "host git repositories over HTTP",
	Long: `Serve bare git repositories over HTTP for clone and push.

Repositories must already exist under --path (create them with
git init --bare) and, to accept pushes, have http.receivepack set to
true. There is no authentication; front this with a reverse proxy if
you need TLS or access control.`,
	Run: func(cmd *cobra.Command, args []string) {
		path := servePath
		if path == "" {
			var err error
			path, err = os.Getwd()
			if err != nil {
				log.Fatalf("failed to get current directory: %v", err)
			}
		}

		execPath, err := gitExecPath()
		if err != nil {
			log.Fatalf("%v", err)
		}
		backend, err := resolveGitHTTPBackend(execPath)
		if err != nil {
			log.Fatalf("%v", err)
		}

		handler := &cgi.Handler{
			Path: backend,
			Env: []string{
				"GIT_PROJECT_ROOT=" + path,
				"GIT_HTTP_EXPORT_ALL=",
			},
		}

		server := &http.Server{
			Addr:    serveListen,
			Handler: handler,
		}

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		go func() {
			<-ctx.Done()
			slog.Info("shutting down serve")
			_ = server.Shutdown(context.Background())
		}()

		slog.Info("serving repositories", "path", path, "listen", serveListen)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("serve failed: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVar(&servePath, "path", "", "path to the repository root (default: current directory)")
	serveCmd.Flags().StringVar(&serveListen, "listen", ":8080", "address to listen on")
}
