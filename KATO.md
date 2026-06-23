# Kato

Kato is a native CLI tool that gives Git small, unobtrusive superpowers in the
terminal.

The goal is not to replace Git. Kato should enhance places where Git is
cumbersome: choosing from lists, remembering exact names, doing small multi-step
chores, making destructive actions safer, and adding useful context before the
user acts.

## Naming

- Tool/family name: `Kato`
- Git command: `kg`
- Git-specific product name when needed: `Kato Git`

The command name `kg` is short for `kato git`. This keeps room for Kato to grow
into other domains later without making the overall product name Git-specific.

## Product Principles

- Kato is a helper, not a Git replacement.
- Kato should be unobtrusive and terminal-native.
- Kato should require explicit subcommands.
- Running `kg` by itself should show help and available commands.
- Kato should support explicit behavior only, not arbitrary pass-through to Git.
- Kato should shell out to the installed `git` binary instead of implementing Git
  behavior through a library.
- Kato should align with the user's existing Git configuration and repository
  behavior.
- Kato should make destructive operations deliberate and clear.

## Tech Stack

- Language: Go
- Command framework: Cobra
- Terminal UI stack: Charm
  - Bubble Tea for interactive flows
  - Bubbles for reusable UI components
  - Lip Gloss for styling
- Distribution target: single native `kg` binary

This matches the existing Go CLI tools in this workspace, which use Cobra, while
Charm is the preferred stack for rich terminal UX.

## MVP Focus

The first implementation focus is:

```sh
kg branch
```

Other commands should be mentioned and reserved in the design, but branch should
be implemented and refined first.

## `kg branch`

`kg branch` helps users switch between local branches without remembering or
typing exact branch names.

### Scope

Initial scope:

- List local branches only.
- Show the current branch clearly.
- Let the user select a branch and switch to it.
- Support branch rename.
- Support branch delete.
- Protect against deleting the current branch.
- Confirm destructive actions before running them.
- Keep remote branch handling out of scope for now.

### Suggested Interaction

When the user runs:

```sh
kg branch
```

Kato opens an interactive local branch picker.

Expected controls:

- `up` / `down` or `j` / `k`: move selection
- type or `/`: filter branches
- `enter`: switch to the selected branch
- `r`: rename the selected branch
- `d`: delete the selected branch after confirmation
- `q` / `esc`: quit without changing anything

### Git Commands

Kato should shell out to Git for behavior.

Likely commands:

```sh
git branch --format=...
git switch <branch>
git branch -m <old> <new>
git branch -d <branch>
```

Use safe structured output where possible. Avoid parsing Git's decorative human
output when Git offers a stable format.

### Delete Behavior

Default delete should be the safe Git delete:

```sh
git branch -d <branch>
```

If Git refuses because the branch is not fully merged, Kato should report that
clearly. Force delete can be considered later, but it should not be the default.

### Rename Behavior

Rename should prompt for the new branch name and then run:

```sh
git branch -m <old> <new>
```

If renaming the current branch, Git supports that, but the UI should make it
clear which branch is being renamed.

## Later Commands

These commands are good candidates after `kg branch`, but they should not drive
the first implementation.

### `kg log`

Purpose: provide a rich terminal history browser good enough that the user does
not need to open a GUI just to understand recent history.

Current intended shape:

- Show a visually clear commit log.
- Support navigation through commits.
- Allow selecting a commit and copying its hash.
- Do not perform Git actions by default.
- Possible later actions: expand details, search/filter, copy hash.

### `kg add`

Purpose: provide a file-level staging picker.

Current intended shape:

- Show changed files from the working tree.
- Include clear status markers for modified, added, deleted, renamed, and
  untracked files.
- Let the user select files.
- Stage selected files with `git add -- <files>`.
- Hunk-level staging is out of scope for the initial version.

### `kg cherry-pick`

Purpose: choose one commit from a searchable commit list and cherry-pick it.

Current intended shape:

- Open an interactive commit picker.
- Show enough commit context to avoid selecting the wrong commit.
- Select exactly one commit.
- Run `git cherry-pick <hash>`.
- If conflicts occur, leave the repository in Git's normal conflict state and
  explain what happened.

### `kg tag`

Purpose: browse and manage tags.

Current intended shape:

- Browse/search tags.
- Checkout selected tag.
- Create a new tag.
- Delete selected tag with confirmation.
- Show useful tag context such as name, target commit, date, and annotation when
  available.

### `kg commit`

Deferred.

Longer-term direction: integrate with an LLM or agent harness, such as Codex, to
write good comprehensive commit messages from staged changes. The user should be
able to review and edit the generated message before committing.

### `kg rebase`

Deferred.

Potential future direction: guided interactive rebase or safer history-editing
workflows. This is high value but higher risk, so it should wait until the basic
interaction model is proven.

## Open Questions

- Exact visual design for the branch picker.
- Whether filtering should start immediately on typing or require `/`.
- Whether branch rows should show last commit summary/date.
- Whether `kg branch` should include a create-new-branch action.
- Whether force delete should exist behind an explicit keybinding or separate
  confirmation flow.
- How clipboard support should be handled cross-platform for later `kg log`.
