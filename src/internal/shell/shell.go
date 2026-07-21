package shell

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

//go:embed scripts/init.sh
var initSh []byte

//go:embed scripts/init.fish
var initFish []byte

//go:embed scripts/init.ps1
var initPs1 []byte

// Kind represents a supported shell type.
type Kind int

const (
	Unknown   Kind = iota
	BashOrZsh      // bash or zsh — share the same init script
	Fish
	PowerShell
)

// Detect infers the current shell from environment variables.
func Detect() Kind {
	if runtime.GOOS == "windows" {
		if sh := os.Getenv("SHELL"); sh != "" {
			return classifyUnix(sh)
		}
		return PowerShell
	}
	return classifyUnix(os.Getenv("SHELL"))
}

func classifyUnix(shell string) Kind {
	switch strings.ToLower(filepath.Base(shell)) {
	case "bash", "zsh", "sh", "dash":
		return BashOrZsh
	case "fish":
		return Fish
	case "pwsh", "powershell":
		return PowerShell
	}
	return Unknown
}

// InitScriptPath returns the path to the kato init script for the given shell.
func InitScriptPath(k Kind) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("finding home directory: %w", err)
	}
	dir := filepath.Join(home, ".kato")
	switch k {
	case BashOrZsh:
		return filepath.Join(dir, "init.sh"), nil
	case Fish:
		return filepath.Join(dir, "init.fish"), nil
	case PowerShell:
		return filepath.Join(dir, "init.ps1"), nil
	default:
		return "", fmt.Errorf("unsupported shell")
	}
}

// WriteInitScript copies the embedded init script for k to ~/.kato/.
// It is safe to call multiple times — it always overwrites with the current content.
func WriteInitScript(k Kind) (string, error) {
	path, err := InitScriptPath(k)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return "", fmt.Errorf("creating kato directory: %w", err)
	}

	var content []byte
	switch k {
	case BashOrZsh:
		content = initSh
	case Fish:
		content = initFish
	case PowerShell:
		content = initPs1
	default:
		return "", fmt.Errorf("unsupported shell")
	}

	if err := os.WriteFile(path, content, 0o640); err != nil {
		return "", fmt.Errorf("writing init script: %w", err)
	}
	return path, nil
}

// SourceLine returns the shell-appropriate line to add to the profile.
func SourceLine(k Kind, initPath string) (string, error) {
	switch k {
	case BashOrZsh, Fish:
		return fmt.Sprintf(`source "%s"`, initPath), nil
	case PowerShell:
		return fmt.Sprintf(`. "%s"`, initPath), nil
	default:
		return "", fmt.Errorf("unsupported shell")
	}
}
