# Flow

## Current flow

The CLI entry point sends parsed arguments to `run`; `--version` writes the
pinned development version to stdout and returns success. The smoke test calls
that same function directly. `internal/fixtures` creates a temporary repository,
initializes Git, writes known history with fixed metadata, and removes it after
the test. `internal/accuracy` runs injected resolvers plus the production
line-range resolver, compares each result with its known reason, and reports the
calculated recovered count. Docker builds the module, runs `go test ./...`, and
CI repeats the checks with networking disabled.

For analysis input, `runWithInput` delegates `--range` to `input.ParseRange` or
`--diff` to `input.ParseDiffSource`. A validated range then calls
`evidence.RecoverLineRange`, which resolves the repository root, runs local
`git log -L`, parses and de-duplicates full-hash/subject citations, and formats
the citations. A validated diff still becomes a stable terminal summary. Input
and evidence failures are deterministic non-zero errors.

## Planned core flow

CLI entry point → input/range or diff parser → local Git evidence collector →
commit/issue/PR/test citation builder → deterministic risk evaluator → terminal
report formatter.

Update this map whenever execution crosses a new file or function.
