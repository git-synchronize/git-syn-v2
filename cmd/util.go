/*
Copyright © 2024-2026 Lucas Ramage <lucas.ramage@infinite-omicron.com>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
)

type remoteEntry struct {
	name string
	url  string
}

// validateRemoteURL returns an error if the URL scheme is not in allowedSchemes.
func validateRemoteURL(rawURL string, allowedSchemes []string) error {
	for _, scheme := range allowedSchemes {
		switch scheme {
		case "https":
			if strings.HasPrefix(rawURL, "https://") {
				return nil
			}
		case "ssh":
			if strings.HasPrefix(rawURL, "git@") || strings.HasPrefix(rawURL, "ssh://") {
				return nil
			}
		}
	}
	return fmt.Errorf("unsupported URL scheme: %s", rawURL)
}

// appendToGitremotes appends a remote stanza to the .gitremotes file.
func appendToGitremotes(path, name, url string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "[remote %q]\n\turl = %s\n", name, url)
	return err
}

// registerRemoteInConfig adds a remote to .git/config if not already present.
func registerRemoteInConfig(repo *git.Repository, name, url string) error {
	cfg, err := repo.Config()
	if err != nil {
		return err
	}
	if cfg.Remotes == nil {
		cfg.Remotes = make(map[string]*gitconfig.RemoteConfig)
	}
	if _, exists := cfg.Remotes[name]; !exists {
		cfg.Remotes[name] = &gitconfig.RemoteConfig{
			Name: name,
			URLs: []string{url},
		}
	}
	return repo.SetConfig(cfg)
}

// parseGitremotes reads a .gitremotes file and returns the list of remote entries.
func parseGitremotes(path string) ([]remoteEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var entries []remoteEntry
	var currentName string

	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "[remote \"") && strings.HasSuffix(trimmed, "\"]") {
			currentName = trimmed[len(`[remote "`) : len(trimmed)-len(`"]`)]
		} else if strings.HasPrefix(trimmed, "url = ") && currentName != "" {
			url := strings.TrimPrefix(trimmed, "url = ")
			slog.Debug("parsed remote entry", "name", currentName, "url", url)
			entries = append(entries, remoteEntry{name: currentName, url: url})
			currentName = ""
		}
	}

	return entries, nil
}
