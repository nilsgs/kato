# kato git branch

`kato git branch` (alias: `kato git b`) opens an inline interactive picker showing local branches with their short commit hash and subject line.

```sh
kato git branch          # default: 10 branches per page
kato git branch -p 5     # show 5 branches per page
kato git b               # shorthand alias
kato git b -p 20
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
kato git branch
# → navigate to main, press enter
```

### Rename current branch

```sh
kato git branch
# → press r
# → type new name, press enter
```

### Delete a merged branch

```sh
kato git branch
# → navigate to the branch, press d
# → press y to confirm
```

### Show more branches at once

```sh
kato git branch -p 20
```
