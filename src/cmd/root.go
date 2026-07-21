package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
)

// NewRootCmd constructs the root kato command. Callers may use this directly in
// tests to get a fresh, isolated command tree.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "kato",
		Short:         "Small, unobtrusive superpowers for your terminal",
		Long:          "Kato enhances command-line workflows with interactive pickers and safer defaults.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.Version = version + "+" + commit
	root.SetVersionTemplate("{{.Version}}\n")

	root.AddCommand(newGitCmd())
	root.AddCommand(newNavCmd())
	root.AddCommand(newAliasCmd())

	return root
}

// Execute runs the root command and exits on error.
func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
