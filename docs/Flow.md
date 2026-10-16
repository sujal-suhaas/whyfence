# Flow

## Current flow

The CLI entry point sends parsed arguments to `run`; `--version` writes the
pinned development version to stdout and returns success. The smoke test calls
that same function directly. Docker builds the module, runs `go test ./...`,
and CI repeats the test and version checks with networking disabled.

## Planned core flow

CLI entry point → input/range or diff parser → local Git evidence collector →
commit/issue/PR/test citation builder → deterministic risk evaluator → terminal
report formatter.

Update this map whenever execution crosses a new file or function.
