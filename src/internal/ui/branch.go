package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"kato/internal/git"
)

// viewState tracks which UI mode is currently active.
type viewState int

const (
	viewBrowse viewState = iota
	viewRename
	viewDeleteConfirm
)

// Result describes the outcome of a branch picker session.
type Result struct {
	Message string
}

// branchItem wraps a git.Branch to satisfy the list.Item interface.
type branchItem struct{ branch git.Branch }

const maxSubject = 50

func (b branchItem) Title() string {
	prefix := "  "
	nameStyle := lipgloss.NewStyle()
	if b.branch.IsCurrent {
		prefix = "* "
		nameStyle = nameStyle.Bold(true)
	}
	name := nameStyle.Render(prefix + b.branch.Name)
	if b.branch.Hash == "" {
		return name
	}
	subject := []rune(b.branch.Subject)
	if len(subject) > maxSubject {
		subject = append(subject[:maxSubject], '…')
	}
	return name + "  " + hashStyle.Render(b.branch.Hash) + "  " + subjectStyle.Render(string(subject))
}

func (b branchItem) Description() string { return "" }
func (b branchItem) FilterValue() string { return b.branch.Name }

var (
	helpStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	subtitleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
	errorStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	hashStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	subjectStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("246"))
	branchNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
)

// BranchModel is the Bubble Tea model for the interactive branch picker.
type BranchModel struct {
	list        list.Model
	renameInput textinput.Model
	state       viewState
	result      Result
	statusMsg   string
	err         error
	quitting    bool
}

const (
	listWidth      = 100
	listMaxVisible = 8
)

// NewBranchModel constructs a BranchModel from a slice of local branches.
func NewBranchModel(branches []git.Branch) BranchModel {
	items := make([]list.Item, len(branches))
	for i, b := range branches {
		items[i] = branchItem{b}
	}

	visible := len(branches)
	if visible > listMaxVisible {
		visible = listMaxVisible
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false

	l := list.New(items, delegate, listWidth, visible)
	l.SetShowTitle(false)
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)

	ri := textinput.New()
	ri.Prompt = ""
	ri.CharLimit = 200

	return BranchModel{list: l, renameInput: ri}
}

// Err returns any error that occurred during the session.
func (m BranchModel) Err() error { return m.err }

// Result returns the outcome of the branch picker session.
func (m BranchModel) Result() Result { return m.result }

// Init implements tea.Model.
func (m BranchModel) Init() tea.Cmd { return nil }

// Update implements tea.Model.
func (m BranchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case viewBrowse:
		return m.updateBrowse(msg)
	case viewRename:
		return m.updateRename(msg)
	case viewDeleteConfirm:
		return m.updateDeleteConfirm(msg)
	}
	return m, nil
}

func (m BranchModel) updateBrowse(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Only update width; height stays compact (set at construction time).
		m.list.SetSize(msg.Width, m.list.Height())
		return m, nil

	case tea.KeyMsg:
		// Only intercept custom keys when the list is not in filter-entry mode.
		if m.list.FilterState() != list.Filtering {
			switch msg.String() {
			case "q", "esc":
				m.quitting = true
				return m, tea.Quit

			case "enter":
				item, ok := m.list.SelectedItem().(branchItem)
				if !ok {
					return m, nil
				}
				if item.branch.IsCurrent {
					m.quitting = true
					return m, tea.Quit
				}
				if err := git.Switch(item.branch.Name); err != nil {
					m.err = err
					m.quitting = true
					return m, tea.Quit
				}
				m.result = Result{
					Message: fmt.Sprintf("Switched to branch '%s'", item.branch.Name),
				}
				m.quitting = true
				return m, tea.Quit

			case "r":
				item, ok := m.list.SelectedItem().(branchItem)
				if !ok {
					return m, nil
				}
				m.renameInput.Placeholder = "new name for '" + item.branch.Name + "'"
				m.renameInput.SetValue("")
				m.state = viewRename
				return m, m.renameInput.Focus()

			case "d":
				item, ok := m.list.SelectedItem().(branchItem)
				if !ok {
					return m, nil
				}
				if item.branch.IsCurrent {
					m.statusMsg = helpStyle.Render("Cannot delete the current branch")
					return m, nil
				}
				m.state = viewDeleteConfirm
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m BranchModel) updateRename(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.state = viewBrowse
			m.renameInput.Blur()
			return m, nil
		case "enter":
			newName := strings.TrimSpace(m.renameInput.Value())
			if newName == "" {
				m.state = viewBrowse
				m.renameInput.Blur()
				return m, nil
			}
			item, ok := m.list.SelectedItem().(branchItem)
			if !ok {
				m.state = viewBrowse
				return m, nil
			}
			if err := git.Rename(item.branch.Name, newName); err != nil {
				m.err = err
				m.quitting = true
				return m, tea.Quit
			}
			m.result = Result{
				Message: fmt.Sprintf("Renamed '%s' to '%s'", item.branch.Name, newName),
			}
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.renameInput, cmd = m.renameInput.Update(msg)
	return m, cmd
}

func (m BranchModel) updateDeleteConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "y", "Y":
			item, ok := m.list.SelectedItem().(branchItem)
			if !ok {
				m.state = viewBrowse
				return m, nil
			}
			if err := git.Delete(item.branch.Name); err != nil {
				m.state = viewBrowse
				m.statusMsg = errorStyle.Render(err.Error())
				return m, nil
			}
			m.result = Result{
				Message: fmt.Sprintf("Deleted branch '%s'", item.branch.Name),
			}
			m.quitting = true
			return m, tea.Quit
		case "n", "N", "esc", "q":
			m.state = viewBrowse
			return m, nil
		}
	}
	return m, nil
}

// View implements tea.Model.
func (m BranchModel) View() string {
	if m.quitting {
		return ""
	}

	switch m.state {
	case viewBrowse:
		if m.statusMsg != "" {
			return m.list.View() + "\n" + m.statusMsg
		}
		return m.list.View()

	case viewRename:
		item, _ := m.list.SelectedItem().(branchItem)
		return subtitleStyle.Render("rename ") + branchNameStyle.Render(item.branch.Name) + subtitleStyle.Render(" → ") + m.renameInput.View()

	case viewDeleteConfirm:
		item, _ := m.list.SelectedItem().(branchItem)
		return errorStyle.Render("delete '"+item.branch.Name+"'?") + helpStyle.Render("  y/n")
	}

	return ""
}
