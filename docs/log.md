# kg log

`kg log` (alias: `kg l`) opens an interactive commit graph browser for the current branch, showing topology with branch relationships.

```sh
kg log              # current branch, last 100 commits, 10 lines per page
kg log --all        # include all branches
kg log -n 200       # load more history
kg log -p 20        # show 20 lines per page
kg l                # shorthand alias
```

## Columns

Each commit row shows: `graph  short-hash  (refs)  author  relative-date  subject`

Connector lines (`|`, `/`, `\`) show branch topology between commits. Branch and tag labels are highlighted in bright green.

## Keybindings

| Key | Action |
|---|---|
| `↑` / `↓` or `j` / `k` | Move selection |
| `space` | Expand selected commit (full message + diff stat) |
| `enter` / `c` | Copy selected commit hash to clipboard |
| `q` / `esc` | Quit without copying |

### In the detail panel

| Key | Action |
|---|---|
| `↑` / `↓` | Scroll the detail panel |
| `space` / `esc` | Collapse the detail panel |

## Examples

### Browse current branch history

```sh
kg log
```

### Include all branches

```sh
kg log --all
```

### Copy a commit hash

```sh
kg log
# → navigate to the commit, press enter or c
# → "Copied abc1234 to clipboard" is printed
```
