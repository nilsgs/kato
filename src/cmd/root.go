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

// NewRootCmd constructs the root kg command. Callers may use this directly in
// tests to get a fresh, isolated command tree.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:           "kg",
		Short:         "Kato Git – small, unobtrusive Git superpowers",
		Long:          "kg enhances common Git workflows with interactive pickers and safer defaults.\nIt shells out to the installed git binary and respects your existing configuration.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.Version = version + "+" + commit
	root.SetVersionTemplate("{{.Version}}\n")

	root.AddCommand(newBranchCmd())

	return root
}

// Execute runs the root command and exits on error.
func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
