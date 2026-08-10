package ui

import (
	"os"
	"path/filepath"
	"testing"

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
