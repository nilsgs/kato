package ui

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"kato/internal/git"
)

const detailHeight = 12

var (
	graphCharStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	refLabelStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	authorStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	dateStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	detailBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("241")).
				Padding(0, 1)
)

// graphItem wraps a git.GraphLine to satisfy the list.Item interface.
type graphItem struct{ line git.GraphLine }

func (g graphItem) Title() string {
	gl := g.line
	graph := graphCharStyle.Render(gl.Graph)
	if !gl.IsCommit() {
		return graph
	}
	hash := hashStyle.Render(gl.Short)
	refs := ""
	if gl.Refs != "" {
		refs = " " + refLabelStyle.Render("("+gl.Refs+")")
	}
	author := authorStyle.Render(truncate(gl.Author, 15))
	date := dateStyle.Render(truncate(gl.Date, 13))
	subject := subjectStyle.Render(truncate(gl.Subject, 55))
	return fmt.Sprintf("%s%s%s  %s  %s  %s", graph, hash, refs, author, date, subject)
}

func (g graphItem) Description() string { return "" }
func (g graphItem) FilterValue() string {
	if !g.line.IsCommit() {
		return ""
	}
	return g.line.Short + " " + g.line.Subject + " " + g.line.Author
}

// LogResult describes the outcome of a log session.
type LogResult struct {
	CopiedHash string
	Message    string
}

// LogModel is the Bubble Tea model for the interactive commit log browser.
type LogModel struct {
	list      list.Model
	detail    viewport.Model
	expanded  bool
	width     int
	result    LogResult
	statusMsg string
	err       error
	quitting  bool
}

// NewLogModel constructs a LogModel from graph lines.
// pageSize controls how many lines are visible at once.
func NewLogModel(lines []git.GraphLine, pageSize int) LogModel {
	items := make([]list.Item, len(lines))
	for i, l := range lines {
		items[i] = graphItem{l}
	}

	visible := len(lines)
	if visible > pageSize {
		visible = pageSize
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = false
	delegate.SetHeight(1)
	delegate.SetSpacing(0)

	l := list.New(items, delegate, listWidth, visible+2)
	l.SetShowTitle(false)
	l.SetShowFilter(false)
	l.SetShowStatusBar(false)
	l.SetShowPagination(true)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)

	return LogModel{list: l}
}

// Err returns any error that occurred during the session.
func (m LogModel) Err() error { return m.err }

// Result returns the outcome of the log session.
func (m LogModel) Result() LogResult { return m.result }

// Init implements tea.Model.
func (m LogModel) Init() tea.Cmd { return nil }

// Update implements tea.Model.
func (m LogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.list.SetSize(msg.Width, m.list.Height())
		m.detail.Width = detailViewportWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		if m.expanded {
			return m.updateExpanded(msg)
		}
		return m.updateBrowse(msg)
	}

	if !m.expanded {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m LogModel) updateBrowse(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		m.quitting = true
		return m, tea.Quit

	case " ":
		item, ok := m.list.SelectedItem().(graphItem)
		if !ok || !item.line.IsCommit() {
			return m, nil
		}
		detail, err := git.CommitDetail(item.line.Hash)
		if err != nil {
			m.statusMsg = errorStyle.Render(err.Error())
			return m, nil
		}
		vp := viewport.New(detailViewportWidth(m.width), detailHeight)
		vp.SetContent(detail)
		m.detail = vp
		m.expanded = true
		m.statusMsg = ""
		return m, nil

	case "enter", "c":
		item, ok := m.list.SelectedItem().(graphItem)
		if !ok || !item.line.IsCommit() {
			return m, nil
		}
		short := item.line.Short
		if err := clipboard.WriteAll(short); err != nil {
			m.statusMsg = errorStyle.Render("clipboard error: " + err.Error())
			return m, nil
		}
		m.result = LogResult{
			CopiedHash: short,
			Message:    fmt.Sprintf("Copied %s to clipboard", short),
		}
		m.quitting = true
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m LogModel) updateExpanded(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc", " ":
		m.expanded = false
		return m, nil
	}
	var cmd tea.Cmd
	m.detail, cmd = m.detail.Update(msg)
	return m, cmd
}

// View implements tea.Model.
func (m LogModel) View() string {
	if m.quitting {
		return ""
	}

	view := m.list.View()

	if m.expanded {
		hint := helpStyle.Render("space/esc: close  ↑/↓: scroll")
		panel := detailBorderStyle.Render(m.detail.View())
		view = view + "\n" + hint + "\n" + panel
	} else if m.statusMsg != "" {
		view = view + "\n" + m.statusMsg
	}

	return view
}

func detailViewportWidth(w int) int {
	if w > 4 {
		return w - 4 // account for border + padding
	}
	return 80
}

func truncate(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max]) + "…"
}
