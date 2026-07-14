package cmd_test

import (
	"bytes"
	"strings"
	"testing"

	"kato/cmd"
)

func TestRootHelp(t *testing.T) {
	root := cmd.NewRootCmd()
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"--help"})
	_ = root.Execute()
	got := buf.String()
	if !strings.Contains(got, "kato") {
		t.Errorf("expected help to mention 'kato', got:\n%s", got)
	}
	if !strings.Contains(got, "git") {
		t.Errorf("expected help to mention 'git' subcommand, got:\n%s", got)
	}
	if strings.Contains(got, "  branch") || strings.Contains(got, "  log") {
		t.Errorf("expected Git commands to be absent from root help, got:\n%s", got)
	}
}

func TestGitHelp(t *testing.T) {
	root := cmd.NewRootCmd()
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"git", "--help"})
	_ = root.Execute()
	got := buf.String()
	if !strings.Contains(got, "branch") {
		t.Errorf("expected Git help to mention 'branch', got:\n%s", got)
	}
	if !strings.Contains(got, "log") {
		t.Errorf("expected Git help to mention 'log', got:\n%s", got)
	}
}

func TestGitCommandAliases(t *testing.T) {
	root := cmd.NewRootCmd()
	for alias, want := range map[string]string{"b": "branch", "l": "log"} {
		found, remaining, err := root.Find([]string{"git", alias})
		if err != nil {
			t.Fatalf("find alias %q: %v", alias, err)
		}
		if found.Name() != want || len(remaining) != 0 {
			t.Errorf("alias %q resolved to %q with remaining args %v; want %q", alias, found.Name(), remaining, want)
		}
	}
}

func TestBranchDefaultPageSize(t *testing.T) {
	root := cmd.NewRootCmd()
	branch, _, err := root.Find([]string{"git", "branch"})
	if err != nil {
		t.Fatalf("find branch command: %v", err)
	}

	page := branch.Flags().Lookup("page")
	if page == nil {
		t.Fatal("expected branch command to define the page flag")
	}
	if page.DefValue != "10" {
		t.Errorf("default page size = %q; want %q", page.DefValue, "10")
	}
}

func TestLegacyRootCommandsAreRejected(t *testing.T) {
	for _, legacy := range []string{"branch", "log"} {
		root := cmd.NewRootCmd()
		root.SetArgs([]string{legacy})
		if err := root.Execute(); err == nil {
			t.Errorf("expected legacy root command %q to be rejected", legacy)
		}
	}
}

func TestRootVersion(t *testing.T) {
	root := cmd.NewRootCmd()
	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"--version"})
	_ = root.Execute()
	got := strings.TrimSpace(buf.String())
	if !strings.Contains(got, "dev") {
		t.Errorf("expected version to contain 'dev', got: %s", got)
	}
}

func TestRootUnknownCommand(t *testing.T) {
	root := cmd.NewRootCmd()
	var errBuf bytes.Buffer
	root.SetErr(&errBuf)
	root.SetArgs([]string{"notacommand"})
	err := root.Execute()
	if err == nil {
		t.Error("expected error for unknown command, got nil")
	}
}
