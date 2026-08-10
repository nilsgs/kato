package ui

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	navPathStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
	navHintStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	navHiddenStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

// folderMeta holds the Nerd Font glyph and hex color for a folder type.
// Requires a Nerd Font installed in the terminal to render correctly.
type folderMeta struct {
	icon  string
	color string
}

// wellKnownFolders maps lowercase folder names to icon+color metadata.
// Colors and icons are inspired by the Terminal-Icons PowerShell module.
var wellKnownFolders = map[string]folderMeta{
	// Version control
	".git":              {"\ue702", "#FF4500"},
	".github":           {"\uf09a", "#C0C0C0"},
	"github":            {"\uf09a", "#C0C0C0"},

	// Development
	"src":               {"\uf121", "#00FF7F"},
	"source":            {"\uf121", "#00FF7F"},
	"dev":               {"\uf121", "#00FF7F"},
	"development":       {"\uf121", "#00FF7F"},
	"projects":          {"\uf121", "#00FF7F"},
	"bin":               {"\uf188", "#00FFF7"},
	"output":            {"\uf019", "#00FF7F"},
	"tests":             {"\uf0c3", "#87CEEB"},
	"test":              {"\uf0c3", "#87CEEB"},
	"specs":             {"\uf0c3", "#87CEEB"},

	// Documentation
	"docs":              {"\uf02d", "#00BFFF"},
	"doc":               {"\uf02d", "#00BFFF"},
	"documents":         {"\uf02d", "#00BFFF"},

	// Media
	"images":            {"\uf03e", "#9ACD32"},
	"image":             {"\uf03e", "#9ACD32"},
	"photos":            {"\uf03e", "#9ACD32"},
	"pictures":          {"\uf03e", "#9ACD32"},
	"videos":            {"\uf008", "#FFA500"},
	"video":             {"\uf008", "#FFA500"},
	"movies":            {"\uf008", "#FFA500"},
	"media":             {"\uf144", "#D3D3D3"},
	"music":             {"\uf001", "#DB7093"},
	"songs":             {"\uf001", "#DB7093"},
	"audio":             {"\uf001", "#DB7093"},

	// Common user folders
	"downloads":         {"\uf019", "#D3D3D3"},
	"download":          {"\uf019", "#D3D3D3"},
	"desktop":           {"\uf108", "#00FBFF"},
	"favorites":         {"\uf005", "#F7D72C"},
	"fonts":             {"\uf031", "#DC143C"},
	"users":             {"\uf0c0", "#F4F4F4"},
	"apps":              {"\uf009", "#FF143C"},
	"applications":      {"\uf009", "#FF143C"},

	// Configuration
	".config":           {"\uf013", "#87CEAF"},
	".cache":            {"\uf013", "#87ECAF"},
	".local":            {"\uf07b", "#D3D3D3"},

	// Editors / IDEs
	".vscode":           {"\ue70c", "#87CEFA"},
	".vscode-insiders":  {"\ue70c", "#24BFA5"},
	".idea":             {"\uf0eb", "#FF6B6B"},

	// Package managers / runtimes
	"node_modules":      {"\ue71e", "#6B8E23"},

	// Cloud / infrastructure
	".docker":           {"\uf308", "#2391E6"},
	".kube":             {"\uf1b2", "#326DE6"},
	".terraform":        {"\uf1b2", "#948EEC"},
	".aws":              {"\uf270", "#EC912D"},
	".azure":            {"\uf17b", "#00BFFF"},
}

const defaultFolderIcon = "\uf07b"
const defaultFolderColor = "#89AAD4"

// folderIcon returns the Nerd Font icon and hex color for a directory name.
func folderIcon(name string) (icon, color string) {
	if m, ok := wellKnownFolders[strings.ToLower(name)]; ok {
		return m.icon, m.color
	}
	return defaultFolderIcon, defaultFolderColor
}

// dirItem wraps a directory name to satisfy the list.Item interface.
type dirItem struct {
	name   string
	icon   string
	color  string
	hidden bool
}

func newDirItem(name string) dirItem {
	icon, color := folderIcon(name)
	return dirItem{
		name:   name,
		icon:   icon,
		color:  color,
		hidden: strings.HasPrefix(name, "."),
	}
}

func (d dirItem) Title() string {
	iconStr := lipgloss.NewStyle().Foreground(lipgloss.Color(d.color)).Render(d.icon)
	name := d.name
	if d.hidden {
		name = navHiddenStyle.Render(name)
	}
	return iconStr + " " + name
}

func (d dirItem) Description() string { return "" }
func (d dirItem) FilterValue() string { return d.name }

// NavModel is the Bubble Tea model for the interactive directory navigator.
type NavModel struct {
	list       list.Model
	currentDir string
	showHidden bool
	chosen     string
	statusMsg  string
	err        error
	quitting   bool
}

// NewNavModel constructs a NavModel starting at startDir.
func NewNavModel(startDir string) (NavModel, error) {
	abs, err := filepath.Abs(startDir)
	if err != nil {
		return NavModel{}, err
	}
	dirs, err := listDirs(abs, false)
	if err != nil {
		return NavModel{}, err
	}
	return NavModel{
		list:       newDirList(dirs),
		currentDir: abs,
	}, nil
}

// Chosen returns the directory path selected by the user, or empty string if cancelled.
func (m NavModel) Chosen() string { return m.chosen }

// Err returns any error that occurred during the session.
func (m NavModel) Err() error { return m.err }

// Init implements tea.Model.
func (m NavModel) Init() tea.Cmd { return nil }

// Update implements tea.Model.
func (m NavModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, m.list.Height())
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			m.chosen = m.currentDir
			if item, ok := m.list.SelectedItem().(dirItem); ok {
				m.chosen = filepath.Join(m.currentDir, item.name)
			}
			m.quitting = true
			return m, tea.Quit

		case "right":
			item, ok := m.list.SelectedItem().(dirItem)
			if !ok {
				return m, nil
			}
			next := filepath.Join(m.currentDir, item.name)
			dirs, err := listDirs(next, m.showHidden)
			if err != nil {
				m.statusMsg = helpStyle.Render(err.Error())
				return m, nil
			}
			m.currentDir = next
			m.list = newDirList(dirs)
			m.statusMsg = ""
			return m, nil

		case "left":
			parent := filepath.Dir(m.currentDir)
			if parent == m.currentDir {
				return m, nil // already at filesystem root
			}
			childName := filepath.Base(m.currentDir)
			dirs, err := listDirs(parent, m.showHidden)
			if err != nil {
				m.statusMsg = helpStyle.Render(err.Error())
				return m, nil
			}
			m.list = newDirList(dirs)
			m.currentDir = parent
			for i, d := range dirs {
				if d == childName {
					m.list.Select(i)
					break
				}
			}
			m.statusMsg = ""
			return m, nil

		case "h":
			m.showHidden = !m.showHidden
			dirs, err := listDirs(m.currentDir, m.showHidden)
			if err != nil {
				m.statusMsg = helpStyle.Render(err.Error())
				return m, nil
			}
			m.list = newDirList(dirs)
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View implements tea.Model.
func (m NavModel) View() string {
	if m.quitting {
		return ""
	}

	header := navPathStyle.Render("\uf07c " + m.currentDir)

	var body string
	if len(m.list.Items()) == 0 {
		body = navHintStyle.Render("  (no subdirectories)")
	} else {
		body = m.list.View()
	}

	hint := navHintStyle.Render("↑/↓: move  →: open  ←: up  enter: select highlighted  h: toggle hidden  q: quit")

	view := header + "\n" + body + "\n" + hint
	if m.statusMsg != "" {
		view += "\n" + m.statusMsg
	}
	return view
}

// newDirList builds a list.Model from a slice of directory names.
func newDirList(dirs []string) list.Model {
	items := make([]list.Item, len(dirs))
	for i, d := range dirs {
		items[i] = newDirItem(d)
	}

	height := len(dirs)
	if height > 15 {
		height = 15
	}
	if height < 1 {
		height = 1
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
	l.SetFilteringEnabled(false)

	return l
}

// listDirs returns subdirectory names inside path.
func listDirs(path string, showHidden bool) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var dirs []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if !showHidden && strings.HasPrefix(e.Name(), ".") {
			continue
		}
		dirs = append(dirs, e.Name())
	}
	return dirs, nil
}
