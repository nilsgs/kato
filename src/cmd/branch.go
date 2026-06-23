package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"kato/internal/git"
	"kato/internal/ui"
)

func newBranchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "branch",
		Short: "Interactively manage local branches",
		Long:  "Open an interactive branch picker to switch, rename, or delete local branches.",
		RunE:  runBranch,
	}
}

func runBranch(cmd *cobra.Command, _ []string) error {
	branches, err := git.ListBranches()
	if err != nil {
		return fmt.Errorf("listing branches: %w", err)
	}
	if len(branches) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No local branches found.")
		return nil
	}

	m := ui.NewBranchModel(branches)
	p := tea.NewProgram(m)
	final, err := p.Run()
	if err != nil {
		return fmt.Errorf("branch picker: %w", err)
	}

	model, ok := final.(ui.BranchModel)
	if !ok {
		return nil
	}

	if model.Err() != nil {
		return model.Err()
	}

	if msg := model.Result().Message; msg != "" {
		fmt.Fprintln(cmd.OutOrStdout(), msg)
	}

	return nil
}
