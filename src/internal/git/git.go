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

// GraphLine represents one output line from git log --graph.
type GraphLine struct {
	Graph   string // leading graph decoration chars (e.g. "* ", "| * ", "| ")
	Hash    string // full 40-char SHA; empty for graph-only connector lines
	Short   string // abbreviated hash
	Refs    string // raw ref decoration (%D), e.g. "HEAD -> main, origin/main"
	Author  string
	Date    string // relative date
	Subject string
}

// IsCommit reports whether the line carries commit data.
func (g GraphLine) IsCommit() bool { return g.Hash != "" }

// ListGraph returns up to n lines of git log --graph output.
// When all is true, --all is passed to include every branch.
func ListGraph(n int, all bool) ([]GraphLine, error) {
	if n <= 0 {
		n = 100
	}
	args := []string{
		"log",
		"--graph",
		fmt.Sprintf("--max-count=%d", n),
		"--pretty=format:%H\t%h\t%D\t%an\t%ar\t%s",
	}
	if all {
		args = append(args, "--all")
	}
	out, err := runGit(args...)
	if err != nil {
		return nil, err
	}
	var lines []GraphLine
	for _, raw := range strings.Split(strings.TrimRight(out, "\n"), "\n") {
		if raw == "" {
			continue
		}
		lines = append(lines, parseGraphLine(raw))
	}
	return lines, nil
}

func parseGraphLine(raw string) GraphLine {
	tabIdx := strings.IndexByte(raw, '\t')
	if tabIdx == -1 {
		// Connector-only line (|, /, \, spaces — no commit data).
		return GraphLine{Graph: raw}
	}
	graphAndHash := raw[:tabIdx]
	rest := raw[tabIdx+1:]
	gl := GraphLine{}
	// The full hash is the last 40 hex chars of the graph+hash segment.
	if len(graphAndHash) >= 40 && isHex(graphAndHash[len(graphAndHash)-40:]) {
		gl.Hash = graphAndHash[len(graphAndHash)-40:]
		gl.Graph = graphAndHash[:len(graphAndHash)-40]
	} else {
		gl.Graph = graphAndHash
	}
	parts := strings.SplitN(rest, "\t", 5)
	if len(parts) > 0 {
		gl.Short = parts[0]
	}
	if len(parts) > 1 {
		gl.Refs = parts[1]
	}
	if len(parts) > 2 {
		gl.Author = parts[2]
	}
	if len(parts) > 3 {
		gl.Date = parts[3]
	}
	if len(parts) > 4 {
		gl.Subject = parts[4]
	}
	return gl
}

func isHex(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return len(s) > 0
}

// CommitDetail returns the full output of git show for the given hash,
// including author, date, full commit message, and diff stat.
func CommitDetail(hash string) (string, error) {
	return runGit("show", "--stat", hash)
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
