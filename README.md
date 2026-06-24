![Banner](img/banner.png)


# Kato

> Small, unobtrusive Git superpowers in your terminal.

`kato` is a native CLI tool that enhances common Git workflows with interactive
pickers and safer defaults. It shells out to the installed `git` binary, so it
works with your existing Git configuration.

## Installation

```sh
# Unix / macOS
sh install.sh

# Windows (PowerShell)
.\install.ps1
```

Both scripts build the binary, copy it to `~/.kato/bin`, and add that directory
to your PATH.

## Commands

| Command | Purpose |
|---|---|
| `kg branch` | Interactive local branch picker: switch, rename, delete |

More commands are planned (`kg log`, `kg add`, `kg cherry-pick`, `kg tag`).

### `kg branch`

Opens an inline interactive branch picker. Each branch shows its short commit hash and subject line.

```sh
kg branch          # default: 3 branches per page
kg branch -p 5     # show 5 branches per page
kg b               # shorthand alias
kg b -p 10
```

Controls:

| Key | Action |
|---|---|
| `↑` / `↓` or `j` / `k` | Move selection |
| Type or `/` | Filter branches |
| `enter` | Switch to selected branch |
| `r` | Rename selected branch |
| `d` | Delete selected branch (with confirmation) |
| `q` / `esc` | Quit without changes |

## Development

```sh
task test    # Run Go unit tests
task build   # Build binary into dist/
task smoke   # Run Smoko end-to-end specs (requires Docker)
task ci      # Run test + build + smoke
```

## Versioning

Builds stamp `VERSION` plus the short Git commit hash:

```sh
kg --version
# 0.1.0+abc1234
```
