package git_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"kato/internal/git"
)

func requireGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available in PATH")
	}
}

func initRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("command %v failed: %s", args, out)
		}
	}
	run("git", "init")
	run("git", "config", "user.email", "test@test.com")
	run("git", "config", "user.name", "Test User")
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# test\n"), 0644); err != nil {
		t.Fatal(err)
	}
	run("git", "add", ".")
	run("git", "commit", "-m", "init")
	return dir
}

func inDir(t *testing.T, dir string) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(orig) })
}

func TestListBranches(t *testing.T) {
	requireGit(t)
	dir := initRepo(t)
	inDir(t, dir)

	branches, err := git.ListBranches()
	if err != nil {
		t.Fatalf("ListBranches: %v", err)
	}
	if len(branches) == 0 {
		t.Fatal("expected at least one branch")
	}
	hasCurrent := false
	for _, b := range branches {
		if b.IsCurrent {
			hasCurrent = true
		}
		if b.Name == "" {
			t.Error("branch has empty name")
		}
	}
	if !hasCurrent {
		t.Error("expected exactly one branch to be marked as current")
	}
}

func TestListBranchesNotRepo(t *testing.T) {
	requireGit(t)
	inDir(t, t.TempDir())

	_, err := git.ListBranches()
	if err == nil {
		t.Error("expected error in non-git directory, got nil")
	}
}

func TestListGraph(t *testing.T) {
	requireGit(t)
	dir := initRepo(t)
	inDir(t, dir)

	lines, err := git.ListGraph(10, false)
	if err != nil {
		t.Fatalf("ListGraph: %v", err)
	}
	if len(lines) == 0 {
		t.Fatal("expected at least one graph line")
	}
	hasCommit := false
	for _, l := range lines {
		if l.IsCommit() {
			hasCommit = true
			if l.Hash == "" {
				t.Error("commit line has empty Hash")
			}
			if l.Short == "" {
				t.Error("commit line has empty Short")
			}
			break
		}
	}
	if !hasCommit {
		t.Error("expected at least one commit line in graph output")
	}
}

func TestListGraphNotRepo(t *testing.T) {
	requireGit(t)
	inDir(t, t.TempDir())

	_, err := git.ListGraph(10, false)
	if err == nil {
		t.Error("expected error in non-git directory, got nil")
	}
}

func TestRenameAndDelete(t *testing.T) {
	requireGit(t)
	dir := initRepo(t)
	inDir(t, dir)

	// Create a new non-current branch.
	cmd := exec.Command("git", "branch", "feature-test")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("create branch: %s", out)
	}

	if err := git.Rename("feature-test", "feature-renamed"); err != nil {
		t.Fatalf("Rename: %v", err)
	}

	branches, err := git.ListBranches()
	if err != nil {
		t.Fatalf("ListBranches after rename: %v", err)
	}
	found := false
	for _, b := range branches {
		if b.Name == "feature-renamed" {
			found = true
		}
	}
	if !found {
		t.Error("renamed branch 'feature-renamed' not found in branch list")
	}

	if err := git.Delete("feature-renamed"); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	branches, _ = git.ListBranches()
	for _, b := range branches {
		if b.Name == "feature-renamed" {
			t.Error("deleted branch 'feature-renamed' still appears in branch list")
		}
	}
}
