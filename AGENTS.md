# kg Agent Notes

## Layout

- `src/main.go`: CLI entry point. Calls `cmd.Execute()`.
- `src/cmd/`: Cobra root command, subcommands, flags, and output formatting.
- `src/cmd/root.go`: `NewRootCmd()`, `Execute()`, version/commit vars stamped at build time.
- `src/cmd/branch.go`: `kg branch` – lists branches, launches the Bubble Tea picker.
- `src/internal/git/`: Low-level git shell-out helpers (ListBranches, Switch, Rename, Delete).
- `src/internal/ui/`: Bubble Tea models for interactive commands.
- `src/internal/ui/branch.go`: Interactive branch picker model (browse, rename, delete-confirm states).
- `install.ps1`: Builds and installs `kg` on Windows; adds `~/.kato/bin` to user PATH.
- `install.sh`: Builds and installs `kg` on Unix; adds `~/.kato/bin` to shell profile.
- `specs/`: Smoko smoke specs.
- `.smokorc`: Smoko image, timeout, and Docker build command.

## Build And Test

Use Task targets for normal validation:

```sh
task test
task build
task smoke
task ci
```

Raw Go fallback for focused debugging:

```sh
cd src
go test ./... -v -count=1
```

## Smoke Tests

Smoke specs live under `specs/` and are run with Smoko inside Docker:

```sh
task smoke
```

`.smokorc` owns the image build. Do not duplicate Docker build commands in
Task or docs.

## Constraints

- `kg` shells out to the installed `git` binary. It does not implement Git
  behavior through a library.
- Do not add pass-through `git` command support. Each subcommand must be
  explicit and intentional.
- Default branch delete uses `-d` (safe delete). Force delete is out of scope
  for the initial version.
- Confirm all destructive actions before executing.
- Do not delete the currently checked-out branch.

## Documentation

- Keep `README.md` as the concise user and developer front door.
- See `docs/usage.md` for extended usage details.

## Versioning

The version comes from `VERSION` and is stamped into builds by the build
scripts via ldflags: `-X kato/cmd.version=<version> -X kato/cmd.commit=<commit>`.
