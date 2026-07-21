# kato nav

Interactively navigate the filesystem tree and jump to the selected directory.

## Usage

```
kato nav
```

The command opens a full-screen picker starting at your current working directory. Navigate with arrow keys and press `enter` to confirm — your shell will change to the selected directory.

> **Note:** Shell integration is installed automatically by `install.sh` / `install.ps1`. After restarting your terminal, `kato nav` changes your directory directly.

> **Icons:** The picker uses [Nerd Font](https://www.nerdfonts.com/) glyphs to show folder-type icons. A Nerd Font must be set as your terminal font for icons to render correctly. Without one, the UI still works — icons appear as placeholder characters.

## Key Bindings

| Key | Action |
|-----|--------|
| `↑` / `↓` | Move selection up / down |
| `→` | Descend into the selected directory |
| `←` | Go up to the parent directory |
| `enter` | Confirm — jump to the currently-browsed directory |
| `h` | Toggle visibility of hidden directories (dot-prefixed) |
| `q` / `esc` | Cancel — stay in the current directory |

## How Shell Integration Works

A Go process cannot change the parent shell's directory directly. The install scripts add a `kato` shell wrapper to your profile that intercepts `kato nav`, runs the picker, and `cd`s to the chosen path. All other `kato` subcommands pass through to the real binary unchanged.

The wrapper is added once and is idempotent — re-running the installer will not duplicate it.

### Manual setup (if needed)

If you installed kato manually without using the install script, add the appropriate snippet to your shell profile:

**bash / zsh** (`~/.bashrc` or `~/.zshrc`):

```sh
# kato shell integration — enables kato nav to change directory
kato() {
  if [ "$1" = "nav" ]; then
    local dir
    dir=$(command kato "$@") && cd "$dir"
  else
    command kato "$@"
  fi
}
```

**fish** (`~/.config/fish/functions/kato.fish`):

```fish
# kato shell integration — enables kato nav to change directory
function kato
  if test "$argv[1]" = "nav"
    set dir (command kato $argv)
    and cd $dir
  else
    command kato $argv
  end
end
```

**PowerShell** (`$PROFILE`):

```powershell
# kato shell integration — enables kato nav to change directory
function kato {
  if ($args[0] -eq 'nav') {
    $dir = & "$env:USERPROFILE\.kato\bin\kato.exe" nav
    if ($dir) { Set-Location $dir }
  } else {
    & "$env:USERPROFILE\.kato\bin\kato.exe" @args
  }
}
```

## Notes

- Navigation starts at the current working directory (`$PWD`).
- Only directories are shown; files are filtered out.
- Hidden directories (names starting with `.`) are hidden by default. Press `h` to toggle them.
- Pressing `esc` or `q` cancels without changing directory.

