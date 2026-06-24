package cmd

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"kato/internal/git"
	"kato/internal/ui"
)

func newBranchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "branch",
		Aliases: []string{"b"},
		Short:   "Interactively manage local branches",
		Long:    "Open an interactive branch picker to switch, rename, or delete local branches.",
		RunE:    runBranch,
	}
	cmd.Flags().IntP("page", "p", 3, "number of branches per page")
	return cmd
}

func runBranch(cmd *cobra.Command, _ []string) error {
	pageSize, _ := cmd.Flags().GetInt("page")
	if pageSize < 1 {
		return fmt.Errorf("page size must be at least 1")
	}

	branches, err := git.ListBranches()
	if err != nil {
		return fmt.Errorf("listing branches: %w", err)
	}
	if len(branches) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No local branches found.")
		return nil
	}

	m := ui.NewBranchModel(branches, pageSize)
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
