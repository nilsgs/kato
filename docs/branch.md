# kg branch

`kg branch` (alias: `kg b`) opens an inline interactive picker showing local branches with their short commit hash and subject line.

```sh
kg branch          # default: 3 branches per page
kg branch -p 5     # show 5 branches per page
kg b               # shorthand alias
kg b -p 10
```

## Keybindings

| Key | Action |
|---|---|
| `↑` / `↓` or `j` / `k` | Move selection |
| Type or `/` | Filter branches |
| `enter` | Switch to selected branch |
| `r` | Rename selected branch |
| `d` | Delete selected branch (with confirmation) |
| `q` / `esc` | Quit without changes |

## Examples

### Switch to main

```sh
kg branch
# → navigate to main, press enter
```

### Rename current branch

```sh
kg branch
# → press r
# → type new name, press enter
```

### Delete a merged branch

```sh
kg branch
# → navigate to the branch, press d
# → press y to confirm
```

### Show more branches at once

```sh
kg branch -p 10
```
