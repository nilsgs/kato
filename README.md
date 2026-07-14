![Banner](img/banner.png)


# Kato

> Small, unobtrusive Git superpowers in your terminal.

`kato` is a native CLI tool that enhances common Git workflows with interactive
pickers and safer defaults. It shells out to the installed `git` binary, so it
works with your existing Git configuration.

![Demo](img/branch_and_log_demo.gif)

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

| Command | Alias | Purpose |
|---|---|---|
| `kato git branch` | `kato git b` | Interactive local branch picker: switch, rename, delete |
| `kato git log` | `kato git l` | Interactive commit graph: visualise topology, copy hash |

More Git commands are planned (`kato git add`, `kato git cherry-pick`, `kato git tag`), and the top-level namespace is available for future non-Git functionality.

See [docs/branch.md](docs/branch.md) and [docs/log.md](docs/log.md) for full keybindings and examples.

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
kato --version
# 0.1.0+abc1234
```
