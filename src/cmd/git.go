package cmd

import "github.com/spf13/cobra"

func newGitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "git",
		Short: "Enhance common Git workflows",
		Long:  "Enhance common Git workflows with interactive pickers and safer defaults. Kato shells out to the installed git binary and respects your existing configuration.",
	}

	cmd.AddCommand(newBranchCmd())
	cmd.AddCommand(newLogCmd())

	return cmd
}
