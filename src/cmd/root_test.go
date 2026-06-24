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
	if !strings.Contains(got, "kg") {
		t.Errorf("expected help to mention 'kg', got:\n%s", got)
	}
	if !strings.Contains(got, "branch") {
		t.Errorf("expected help to mention 'branch' subcommand, got:\n%s", got)
	}
	if !strings.Contains(got, "log") {
		t.Errorf("expected help to mention 'log' subcommand, got:\n%s", got)
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
