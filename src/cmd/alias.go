package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"kato/internal/alias"
	"kato/internal/shell"
	"kato/internal/ui"
)

func newAliasCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "alias",
		Short: "Interactively manage shell aliases",
		Long: "Open an interactive alias manager to browse, add, edit, delete, and search " +
			"shell aliases. Aliases are stored in ~/.kato/aliases and loaded by your shell at startup.",
		RunE: runAlias,
	}
}

func runAlias(_ *cobra.Command, _ []string) error {
	sh := shell.Detect()

	if _, err := shell.WriteInitScript(sh); err != nil {
		// Non-fatal: warn and continue.
		fmt.Fprintln(os.Stderr, "warning: could not write kato init script:", err)
	}

	aliases, err := alias.Load()
	if err != nil {
		return fmt.Errorf("loading aliases: %w", err)
	}

	if len(aliases) == 0 {
		// Show setup hint and launch TUI so the user can add their first alias.
		fmt.Fprintln(os.Stderr, ui.AliasSetupHint(sh))
	}

	m := ui.NewAliasModel(aliases)
	p := tea.NewProgram(m)
	final, err := p.Run()
	if err != nil {
		return fmt.Errorf("alias manager: %w", err)
	}

	model, ok := final.(ui.AliasModel)
	if !ok {
		return nil
	}

	if model.Err() != nil {
		return model.Err()
	}

	if msg := model.Result().Message; msg != "" {
		fmt.Fprintln(os.Stdout, msg)
	}

	return nil
}
