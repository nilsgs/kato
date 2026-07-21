package ui

import (
	"fmt"
	"strings"

	"kato/internal/alias"
	"kato/internal/shell"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// aliasViewState tracks which UI mode is currently active.
type aliasViewState int

const (
	aliasViewBrowse aliasViewState = iota
	aliasViewAdd
	aliasViewEdit
	aliasViewDeleteConfirm
)

// aliasItem wraps an alias.Alias to satisfy the list.Item interface.
type aliasItem struct{ a alias.Alias }

func (i aliasItem) Title() string {
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6")).Render(i.a.Name) +
		"  " + helpStyle.Render(i.a.Command)
}
func (i aliasItem) Description() string { return "" }
func (i aliasItem) FilterValue() string { return i.a.Name + " " + i.a.Command }

// AliasResult describes the outcome of an alias manager session.
type AliasResult struct {
	Message string
}

// AliasModel is the Bubble Tea model for the interactive alias manager.
type AliasModel struct {
	list      list.Model
	nameInput textinput.Model
	cmdInput  textinput.Model
	state     aliasViewState
	editName  string // original name when editing
	result    AliasResult
	statusMsg string
	err       error
	quitting  bool
}

// NewAliasModel constructs an AliasModel from a slice of aliases.
func NewAliasModel(aliases []alias.Alias) AliasModel {
	items := make([]list.Item, len(aliases))
	for i, a := range aliases {
		items[i] = aliasItem{a}
	}

	height := len(aliases)
	if height > 15 {
		height = 15
	}
	if height < 3 {
		height = 3
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetHeight(1)
	delegate.SetSpacing(0)

	l := list.New(items, delegate, listWidth, height+2)
	l.SetShowTitle(false)
	l.SetShowFilter(false)
	l.SetShowStatusBar(false)
	l.SetShowPagination(true)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)

	nameIn := textinput.New()
	nameIn.Prompt = ""
	nameIn.CharLimit = 64
	nameIn.Placeholder = "name"

	cmdIn := textinput.New()
	cmdIn.Prompt = ""
	cmdIn.CharLimit = 512
	cmdIn.Placeholder = "command"

	return AliasModel{list: l, nameInput: nameIn, cmdInput: cmdIn}
}

// Err returns any error that occurred during the session.
func (m AliasModel) Err() error { return m.err }

// Result returns the outcome of the session.
func (m AliasModel) Result() AliasResult { return m.result }

// Init implements tea.Model.
func (m AliasModel) Init() tea.Cmd { return nil }

// Update implements tea.Model.
func (m AliasModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case aliasViewBrowse:
		return m.updateBrowse(msg)
	case aliasViewAdd, aliasViewEdit:
		return m.updateForm(msg)
	case aliasViewDeleteConfirm:
		return m.updateDeleteConfirm(msg)
	}
	return m, nil
}

func (m AliasModel) updateBrowse(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, m.list.Height())
		return m, nil

	case tea.KeyMsg:
		if m.list.FilterState() != list.Filtering {
			switch msg.String() {
			case "q", "esc":
				m.quitting = true
				return m, tea.Quit

			case "a":
				m.nameInput.SetValue("")
				m.cmdInput.SetValue("")
				m.state = aliasViewAdd
				m.statusMsg = ""
				return m, m.nameInput.Focus()

			case "enter", "e":
				item, ok := m.list.SelectedItem().(aliasItem)
				if !ok {
					return m, nil
				}
				m.editName = item.a.Name
				m.nameInput.SetValue(item.a.Name)
				m.cmdInput.SetValue(item.a.Command)
				m.state = aliasViewEdit
				m.statusMsg = ""
				return m, m.nameInput.Focus()

			case "d":
				if _, ok := m.list.SelectedItem().(aliasItem); !ok {
					return m, nil
				}
				m.state = aliasViewDeleteConfirm
				m.statusMsg = ""
				return m, nil

			case "c":
				item, ok := m.list.SelectedItem().(aliasItem)
				if !ok {
					return m, nil
				}
				if err := clipboard.WriteAll(item.a.Command); err != nil {
					m.statusMsg = errorStyle.Render("clipboard: " + err.Error())
				} else {
					m.statusMsg = helpStyle.Render("Copied: " + item.a.Command)
				}
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m AliasModel) updateForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.state = aliasViewBrowse
			m.nameInput.Blur()
			m.cmdInput.Blur()
			return m, nil

		case "tab", "shift+tab":
			// Toggle focus between name and command fields.
			if m.nameInput.Focused() {
				m.nameInput.Blur()
				return m, m.cmdInput.Focus()
			}
			m.cmdInput.Blur()
			return m, m.nameInput.Focus()

		case "enter":
			name := strings.TrimSpace(m.nameInput.Value())
			cmd := strings.TrimSpace(m.cmdInput.Value())
			if name == "" || cmd == "" {
				m.statusMsg = errorStyle.Render("Name and command are required")
				return m, nil
			}
			var err error
			if m.state == aliasViewAdd {
				err = alias.Add(name, cmd)
			} else {
				err = alias.Update(m.editName, name, cmd)
			}
			if err != nil {
				m.statusMsg = errorStyle.Render(err.Error())
				return m, nil
			}
			action := "Added"
			if m.state == aliasViewEdit {
				action = "Updated"
			}
			m.result = AliasResult{Message: fmt.Sprintf("%s alias '%s'", action, name)}
			m.quitting = true
			return m, tea.Quit
		}
	}

	var nameCmd, cmdCmd tea.Cmd
	m.nameInput, nameCmd = m.nameInput.Update(msg)
	m.cmdInput, cmdCmd = m.cmdInput.Update(msg)
	return m, tea.Batch(nameCmd, cmdCmd)
}

func (m AliasModel) updateDeleteConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "y", "Y":
			item, ok := m.list.SelectedItem().(aliasItem)
			if !ok {
				m.state = aliasViewBrowse
				return m, nil
			}
			if err := alias.Delete(item.a.Name); err != nil {
				m.state = aliasViewBrowse
				m.statusMsg = errorStyle.Render(err.Error())
				return m, nil
			}
			m.result = AliasResult{Message: fmt.Sprintf("Deleted alias '%s'", item.a.Name)}
			m.quitting = true
			return m, tea.Quit
		case "n", "N", "esc", "q":
			m.state = aliasViewBrowse
			return m, nil
		}
	}
	return m, nil
}

// View implements tea.Model.
func (m AliasModel) View() string {
	if m.quitting {
		return ""
	}

	switch m.state {
	case aliasViewBrowse:
		view := m.list.View()
		hint := helpStyle.Render("a: add  enter/e: edit  d: delete  c: copy  /: filter  q: quit")
		if m.list.FilterState() == list.Filtering {
			view = helpStyle.Render("filter: "+m.list.FilterInput.Value()+"_") + "\n" + view
		} else if m.statusMsg != "" {
			view = view + "\n" + m.statusMsg
		}
		return view + "\n" + hint

	case aliasViewAdd, aliasViewEdit:
		label := "new alias"
		if m.state == aliasViewEdit {
			label = "edit '" + m.editName + "'"
		}
		header := subtitleStyle.Render(label)

		nameLabel := "  name:    "
		cmdLabel := "  command: "
		if m.nameInput.Focused() {
			nameLabel = lipgloss.NewStyle().Bold(true).Render(nameLabel)
		}
		if m.cmdInput.Focused() {
			cmdLabel = lipgloss.NewStyle().Bold(true).Render(cmdLabel)
		}

		form := header + "\n" +
			nameLabel + m.nameInput.View() + "\n" +
			cmdLabel + m.cmdInput.View()
		hint := helpStyle.Render("\ntab: switch field  enter: save  esc: cancel")
		if m.statusMsg != "" {
			hint = "\n" + m.statusMsg + hint
		}
		return form + hint

	case aliasViewDeleteConfirm:
		item, _ := m.list.SelectedItem().(aliasItem)
		return errorStyle.Render("delete '"+item.a.Name+"'?") + helpStyle.Render("  y/n")
	}

	return ""
}

// AliasSetupHint returns a setup hint shown when no aliases exist yet.
func AliasSetupHint(k shell.Kind) string {
	initPath, _ := shell.InitScriptPath(k)
	sourceLine, err := shell.SourceLine(k, initPath)
	if err != nil {
		return ""
	}
	return helpStyle.Render(
		"No aliases yet. Press 'a' to add one.\n\n" +
			"To load aliases in your shell, add this line to your shell profile:\n  " +
			sourceLine,
	)
}
