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

| Command | Alias | Purpose |
|---|---|---|
| `kg branch` | `kg b` | Interactive local branch picker: switch, rename, delete |
| `kg log` | `kg l` | Interactive commit graph: visualise topology, copy hash |

More commands are planned (`kg add`, `kg cherry-pick`, `kg tag`).

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
kg --version
# 0.1.0+abc1234
```
