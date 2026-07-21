package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"kato/internal/ui"
)

func newNavCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "nav",
		Short: "Interactively navigate to a directory",
		Long:  "Open an interactive directory picker. Use arrow keys to browse the filesystem tree, then press enter to select a directory. Prints the chosen path to stdout for shell integration (e.g. cd $(kato nav)).",
		RunE:  runNav,
	}
}

func runNav(cmd *cobra.Command, _ []string) error {
	startDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	m, err := ui.NewNavModel(startDir)
	if err != nil {
		return fmt.Errorf("initialising navigator: %w", err)
	}

	p := tea.NewProgram(m, tea.WithOutput(os.Stderr))
	final, err := p.Run()
	if err != nil {
		return fmt.Errorf("nav: %w", err)
	}

	model, ok := final.(ui.NavModel)
	if !ok {
		return nil
	}

	if model.Err() != nil {
		return model.Err()
	}

	if path := model.Chosen(); path != "" {
		fmt.Fprintln(cmd.OutOrStdout(), path)
	}

	return nil
}
