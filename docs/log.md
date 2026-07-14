# kato git log

`kato git log` (alias: `kato git l`) opens an interactive commit graph browser for the current branch, showing topology with branch relationships.

```sh
kato git log              # current branch, last 100 commits, 10 lines per page
kato git log --all        # include all branches
kato git log -n 200       # load more history
kato git log -p 20        # show 20 lines per page
kato git l                # shorthand alias
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
kato git log
```

### Include all branches

```sh
kato git log --all
```

### Copy a commit hash

```sh
kato git log
# → navigate to the commit, press enter or c
# → "Copied abc1234 to clipboard" is printed
```
