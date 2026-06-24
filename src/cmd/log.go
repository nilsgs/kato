package cmd

import (
	"fmt"

	"kato/internal/git"
	"kato/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func newLogCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "log",
		Aliases: []string{"l"},
		Short:   "Interactively browse the commit graph",
		Long:    "Open an interactive commit graph showing branch topology. Navigate commits and copy a hash to the clipboard.",
		RunE:    runLog,
	}
	cmd.Flags().IntP("count", "n", 100, "number of commits to load")
	cmd.Flags().IntP("page", "p", 10, "number of lines visible at once")
	cmd.Flags().Bool("all", false, "include all branches")
	return cmd
}

func runLog(cmd *cobra.Command, _ []string) error {
	count, _ := cmd.Flags().GetInt("count")
	if count < 1 {
		return fmt.Errorf("count must be at least 1")
	}
	pageSize, _ := cmd.Flags().GetInt("page")
	if pageSize < 1 {
		return fmt.Errorf("page size must be at least 1")
	}
	all, _ := cmd.Flags().GetBool("all")

	lines, err := git.ListGraph(count, all)
	if err != nil {
		return fmt.Errorf("building log: %w", err)
	}
	if len(lines) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No commits found.")
		return nil
	}

	m := ui.NewLogModel(lines, pageSize)
	p := tea.NewProgram(m)
	final, err := p.Run()
	if err != nil {
		return fmt.Errorf("log browser: %w", err)
	}

	model, ok := final.(ui.LogModel)
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
