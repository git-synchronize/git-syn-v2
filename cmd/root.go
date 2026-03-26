/*
Copyright © 2024-2026 Lucas Ramage <lucas.ramage@infinite-omicron.com>
*/
package cmd

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var (
	verbose bool
	debug   bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "git-syn",
	Short: "Remote git repository synchronization.",
	Long:  "Remote git repository synchronization.",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var level slog.Level
		var out io.Writer

		if debug {
			level = slog.LevelDebug
			out = os.Stderr
		} else if verbose {
			level = slog.LevelInfo
			out = os.Stderr
		} else {
			level = slog.LevelError + 1
			out = io.Discard
		}

		handler := slog.NewTextHandler(out, &slog.HandlerOptions{
			Level: level,
		})
		slog.SetDefault(slog.New(handler))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("version", "v", false, "output version information and exit")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "enable verbose output")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug output (implies verbose)")

	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		fmt.Printf("Usage: git-syn [option] ... [command] ...\n\n")
		fmt.Printf("Remote git repository synchronization.\n\n")
		fmt.Printf("  %-17s%s\n", "-h, --help", "display this help and exit")
		fmt.Printf("  %-17s%s\n", "-v, --version", "output version information and exit")
		fmt.Printf("  %-17s%s\n", "--verbose", "enable verbose output")
		fmt.Printf("  %-17s%s\n", "--debug", "enable debug output (implies verbose)")
		for _, c := range cmd.Commands() {
			if !c.Hidden && c.Name() != "help" {
				fmt.Printf("  %-17s%s\n", c.Name(), c.Short)
			}
		}
		fmt.Printf("\nGit SYN online help: <https://gitlab.com/git-syn/git-syn>\n")
		fmt.Printf("Full documentation <https://git-syn.gitlab.io/git-syn>\n")
		fmt.Printf("or available locally via: man git-syn\n")
	})
}
