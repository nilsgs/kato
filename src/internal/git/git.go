package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Branch represents a local Git branch.
type Branch struct {
	Name      string
	IsCurrent bool
	Hash      string // short commit hash of HEAD
	Subject   string // HEAD commit subject line
}

// ListBranches returns all local Git branches.
// The currently checked-out branch is marked with IsCurrent = true.
func ListBranches() ([]Branch, error) {
	out, err := runGit("branch", "--format=%(HEAD)%(refname:short)\t%(objectname:short)\t%(subject)")
	if err != nil {
		return nil, err
	}
	var branches []Branch
	for _, line := range strings.Split(strings.TrimRight(out, "\n"), "\n") {
		if len(line) == 0 {
			continue
		}
		current := line[0] == '*'
		parts := strings.SplitN(line[1:], "\t", 3)
		name := strings.TrimSpace(parts[0])
		if name == "" {
			continue
		}
		b := Branch{Name: name, IsCurrent: current}
		if len(parts) > 1 {
			b.Hash = strings.TrimSpace(parts[1])
		}
		if len(parts) > 2 {
			b.Subject = strings.TrimSpace(parts[2])
		}
		branches = append(branches, b)
	}
	return branches, nil
}

// Switch switches the working tree to the named branch using git switch.
func Switch(name string) error {
	_, err := runGit("switch", name)
	return err
}

// Rename renames oldName to newName using git branch -m.
func Rename(oldName, newName string) error {
	_, err := runGit("branch", "-m", oldName, newName)
	return err
}

// Delete deletes a branch using the safe delete flag (-d).
// Git will refuse to delete a branch that has not been fully merged.
func Delete(name string) error {
	_, err := runGit("branch", "-d", name)
	return err
}

// runGit runs git with the given arguments and returns combined stdout.
// On non-zero exit it returns an error containing the stderr message.
func runGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("%s", msg)
	}
	return stdout.String(), nil
}
