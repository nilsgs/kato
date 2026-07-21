package alias

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Alias represents a named shell alias.
type Alias struct {
	Name    string
	Command string
}

// storePath returns the path to ~/.kato/aliases.
func storePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("finding home directory: %w", err)
	}
	return filepath.Join(home, ".kato", "aliases"), nil
}

// Load reads all aliases from the store. Returns an empty slice if the file
// does not exist yet.
func Load() ([]Alias, error) {
	path, err := storePath()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("opening alias store: %w", err)
	}
	defer f.Close()

	return parse(f), nil
}

// parse reads name=command lines from r, skipping blank lines and comments.
func parse(r interface{ Read([]byte) (int, error) }) []Alias {
	var aliases []Alias
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.IndexByte(line, '=')
		if idx < 1 {
			continue
		}
		aliases = append(aliases, Alias{
			Name:    strings.TrimSpace(line[:idx]),
			Command: strings.TrimSpace(line[idx+1:]),
		})
	}
	return aliases
}

// Save writes aliases to the store, overwriting any existing content.
func Save(aliases []Alias) error {
	path, err := storePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return fmt.Errorf("creating kato directory: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("writing alias store: %w", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, a := range aliases {
		fmt.Fprintf(w, "%s=%s\n", a.Name, a.Command)
	}
	return w.Flush()
}

// Add appends a new alias. Returns an error if the name already exists.
func Add(name, command string) error {
	aliases, err := Load()
	if err != nil {
		return err
	}
	for _, a := range aliases {
		if a.Name == name {
			return fmt.Errorf("alias %q already exists", name)
		}
	}
	aliases = append(aliases, Alias{Name: name, Command: command})
	return Save(aliases)
}

// Update replaces the command for an existing alias. Returns an error if the
// name is not found.
func Update(name, newName, newCommand string) error {
	aliases, err := Load()
	if err != nil {
		return err
	}
	found := false
	for i, a := range aliases {
		if a.Name == name {
			aliases[i] = Alias{Name: newName, Command: newCommand}
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("alias %q not found", name)
	}
	return Save(aliases)
}

// Delete removes an alias by name. Returns an error if not found.
func Delete(name string) error {
	aliases, err := Load()
	if err != nil {
		return err
	}
	n := len(aliases)
	aliases = filterOut(aliases, name)
	if len(aliases) == n {
		return fmt.Errorf("alias %q not found", name)
	}
	return Save(aliases)
}

func filterOut(aliases []Alias, name string) []Alias {
	out := aliases[:0]
	for _, a := range aliases {
		if a.Name != name {
			out = append(out, a)
		}
	}
	return out
}
