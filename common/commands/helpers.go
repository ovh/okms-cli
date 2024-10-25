package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func NewMarkdownCmd(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:    "markdown TARGET",
		Hidden: true,
		Args:   cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			if err := os.MkdirAll(args[0], 0o755); err != nil {
				fmt.Fprintln(os.Stderr, "ERROR:", err)
				os.Exit(1)
			}
			if err := doc.GenMarkdownTree(rootCmd, args[0]); err != nil {
				fmt.Fprintln(os.Stderr, "ERROR:", err)
				os.Exit(1)
			}
		},
	}
}

func NewVersionCmd(version, commit, date *string) *cobra.Command {
	if version == nil {
		s := ""
		version = &s
	}
	return &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Print the version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s\n", *version)
			fmt.Printf("Commit: %s\n", *commit)
			fmt.Printf("Date: %s\n", *date)
		},
	}
}
