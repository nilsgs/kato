package ui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func TestNavEnterSelectsHighlightedDirectory(t *testing.T) {
	root := t.TempDir()
	mkdirAll(t, filepath.Join(root, "alpha"), filepath.Join(root, "beta"))

	m := newNavModelForTest(t, root)
	m = updateNavForTest(t, m, tea.KeyMsg{Type: tea.KeyEnter})

	want := filepath.Join(root, "alpha")
	if got := m.Chosen(); got != want {
		t.Errorf("chosen directory = %q; want %q", got, want)
	}
}

func TestNavEnterSelectsMovedHighlight(t *testing.T) {
	root := t.TempDir()
	mkdirAll(t, filepath.Join(root, "alpha"), filepath.Join(root, "beta"))

	m := newNavModelForTest(t, root)
	m = updateNavForTest(t, m, tea.KeyMsg{Type: tea.KeyDown})
	m = updateNavForTest(t, m, tea.KeyMsg{Type: tea.KeyEnter})

	want := filepath.Join(root, "beta")
	if got := m.Chosen(); got != want {
		t.Errorf("chosen directory = %q; want %q", got, want)
	}
}

func TestNavEnterSelectsHighlightedDirectoryAfterDescending(t *testing.T) {
	root := t.TempDir()
	mkdirAll(t, filepath.Join(root, "parent", "child"))

	m := newNavModelForTest(t, root)
	m = updateNavForTest(t, m, tea.KeyMsg{Type: tea.KeyRight})
	m = updateNavForTest(t, m, tea.KeyMsg{Type: tea.KeyEnter})

	want := filepath.Join(root, "parent", "child")
	if got := m.Chosen(); got != want {
		t.Errorf("chosen directory = %q; want %q", got, want)
	}
}

func TestNavEnterSelectsCurrentDirectoryWhenEmpty(t *testing.T) {
	root := t.TempDir()
	m := newNavModelForTest(t, root)

	m = updateNavForTest(t, m, tea.KeyMsg{Type: tea.KeyEnter})

	want, err := filepath.Abs(root)
	if err != nil {
		t.Fatalf("resolve temporary directory: %v", err)
	}
	if got := m.Chosen(); got != want {
		t.Errorf("chosen directory = %q; want %q", got, want)
	}
}

func TestNavTypingStartsFilter(t *testing.T) {
	root := t.TempDir()
	mkdirAll(t, filepath.Join(root, "alpha"), filepath.Join(root, "beta"))
	m := newNavModelForTest(t, root)

	m = updateNavForTest(t, m, runeKey('b'))

	if got := m.list.FilterState(); got != list.FilterApplied {
		t.Errorf("filter state = %v; want %v so selection is rendered", got, list.FilterApplied)
	}
	if got := m.list.FilterInput.Value(); got != "b" {
		t.Errorf("filter value = %q; want %q", got, "b")
	}
	if item, ok := m.list.SelectedItem().(dirItem); !ok || item.name != "beta" {
		t.Errorf("selected item = %v; want beta", m.list.SelectedItem())
	}
}

func TestNavBackspaceRestoresMatchesForNavigation(t *testing.T) {
	root := t.TempDir()
	mkdirAll(t, filepath.Join(root, "alpha"), filepath.Join(root, "alpine"), filepath.Join(root, "beta"))
	m := newNavModelForTest(t, root)

	m = updateNavForTest(t, m, runeKey('a'))
	m = updateNavForTest(t, m, runeKey('l'))
	m = updateNavForTest(t, m, runeKey('ø'))
	if got := len(m.list.VisibleItems()); got != 0 {
		t.Fatalf("visible directories after typo = %d; want 0", got)
	}

	m = updateNavForTest(t, m, tea.KeyMsg{Type: tea.KeyBackspace})
	if got := m.list.FilterInput.Value(); got != "al" {
		t.Errorf("filter value = %q; want %q", got, "al")
	}
	if got := len(m.list.VisibleItems()); got != 2 {
		t.Errorf("visible directories after backspace = %d; want 2", got)
	}

	m = updateNavForTest(t, m, tea.KeyMsg{Type: tea.KeyDown})
	m = updateNavForTest(t, m, tea.KeyMsg{Type: tea.KeyEnter})
	if got, want := m.Chosen(), filepath.Join(root, "alpine"); got != want {
		t.Errorf("chosen directory = %q; want %q", got, want)
	}
}

func TestNavBackspaceClearsFinalFilterCharacter(t *testing.T) {
	root := t.TempDir()
	mkdirAll(t, filepath.Join(root, "alpha"), filepath.Join(root, "beta"))
	m := newNavModelForTest(t, root)

	m = updateNavForTest(t, m, runeKey('b'))
	m = updateNavForTest(t, m, tea.KeyMsg{Type: tea.KeyBackspace})
	assertNavFilterReset(t, m)
	if got := len(m.list.VisibleItems()); got != 2 {
		t.Errorf("visible directories after clearing the filter = %d; want 2", got)
	}

	m = updateNavForTest(t, m, tea.KeyMsg{Type: tea.KeyDown})
	m = updateNavForTest(t, m, tea.KeyMsg{Type: tea.KeyBackspace})
	if item, ok := m.list.SelectedItem().(dirItem); !ok || item.name != "beta" || m.quitting {
		t.Errorf("backspace with no filter changed the selection or quit, selected item = %v", m.list.SelectedItem())
	}
}

func TestNavActionLettersAreFilterable(t *testing.T) {
	root := t.TempDir()
	mkdirAll(t, filepath.Join(root, "query"))

	for _, key := range []rune{'h', 'q'} {
		t.Run(string(key), func(t *testing.T) {
			m := newNavModelForTest(t, root)
			m = updateNavForTest(t, m, runeKey(key))

			if got := m.list.FilterInput.Value(); got != string(key) {
				t.Errorf("filter value = %q; want %q", got, string(key))
			}
			if m.quitting {
				t.Error("model quit while entering a filter")
			}
		})
	}
}

func TestNavAltHTogglesHiddenDirectories(t *testing.T) {
	root := t.TempDir()
	mkdirAll(t, filepath.Join(root, ".hidden"), filepath.Join(root, "visible"))
	m := newNavModelForTest(t, root)

	m = updateNavForTest(t, m, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}, Alt: true})

	if !m.showHidden {
		t.Error("showHidden = false; want true")
	}
	if got := len(m.list.Items()); got != 2 {
		t.Errorf("directory count = %d; want 2", got)
	}
}

func TestNavEnterSelectsFilteredDirectory(t *testing.T) {
	root := t.TempDir()
	mkdirAll(t, filepath.Join(root, "alpha"), filepath.Join(root, "beta"))
	m := newNavModelForTest(t, root)
	m = updateNavForTest(t, m, runeKey('b'))

	m = updateNavForTest(t, m, tea.KeyMsg{Type: tea.KeyEnter})

	if got, want := m.Chosen(), filepath.Join(root, "beta"); got != want {
		t.Errorf("chosen directory = %q; want %q", got, want)
	}
}

func TestNavArrowMovesThroughFilteredDirectories(t *testing.T) {
	root := t.TempDir()
	mkdirAll(t, filepath.Join(root, "alpha"), filepath.Join(root, "alpine"), filepath.Join(root, "beta"))
	m := newNavModelForTest(t, root)
	m = updateNavForTest(t, m, runeKey('a'))
	m = updateNavForTest(t, m, runeKey('l'))

	m = updateNavForTest(t, m, tea.KeyMsg{Type: tea.KeyDown})

	item, ok := m.list.SelectedItem().(dirItem)
	if !ok {
		t.Fatal("filtered list has no selected directory")
	}
	if got, want := item.name, "alpine"; got != want {
		t.Errorf("selected directory = %q; want %q", got, want)
	}
}

func TestNavEscClearsFilterBeforeQuitting(t *testing.T) {
	root := t.TempDir()
	mkdirAll(t, filepath.Join(root, "alpha"), filepath.Join(root, "beta"))
	m := newNavModelForTest(t, root)
	m = updateNavForTest(t, m, runeKey('b'))

	m = updateNavForTest(t, m, tea.KeyMsg{Type: tea.KeyEsc})
	if m.quitting {
		t.Error("first esc quit with an applied filter")
	}
	assertNavFilterReset(t, m)
	if got := len(m.list.VisibleItems()); got != 2 {
		t.Errorf("visible directories after esc = %d; want 2", got)
	}

	m = updateNavForTest(t, m, tea.KeyMsg{Type: tea.KeyEsc})
	if !m.quitting {
		t.Error("second esc did not quit")
	}
}

func TestNavTraversalResetsFilter(t *testing.T) {
	root := t.TempDir()
	mkdirAll(t, filepath.Join(root, "parent", "child"))
	m := newNavModelForTest(t, root)
	m.list.SetFilterText("parent")

	m = updateNavForTest(t, m, tea.KeyMsg{Type: tea.KeyRight})
	assertNavFilterReset(t, m)

	m.list.SetFilterText("child")
	m = updateNavForTest(t, m, tea.KeyMsg{Type: tea.KeyLeft})
	assertNavFilterReset(t, m)
}

func assertNavFilterReset(t *testing.T, m NavModel) {
	t.Helper()
	if got := m.list.FilterState(); got != list.Unfiltered {
		t.Errorf("filter state = %v; want %v", got, list.Unfiltered)
	}
	if got := m.list.FilterInput.Value(); got != "" {
		t.Errorf("filter value = %q; want empty", got)
	}
}

func runeKey(r rune) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
}

func newNavModelForTest(t *testing.T, root string) NavModel {
	t.Helper()
	m, err := NewNavModel(root)
	if err != nil {
		t.Fatalf("create navigation model: %v", err)
	}
	return m
}

func updateNavForTest(t *testing.T, m NavModel, msg tea.Msg) NavModel {
	t.Helper()
	updated, _ := m.Update(msg)
	next, ok := updated.(NavModel)
	if !ok {
		t.Fatalf("updated model has type %T; want NavModel", updated)
	}
	return next
}

func mkdirAll(t *testing.T, paths ...string) {
	t.Helper()
	for _, path := range paths {
		if err := os.MkdirAll(path, 0o755); err != nil {
			t.Fatalf("create directory %q: %v", path, err)
		}
	}
}
