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
the citations. Each subject is passed to `evidence.ExtractReferences`, which
recognizes documented explicit issue/PR forms and keeps each reference attached
to its source commit before the CLI formats it. The collector then runs a
root-safe local `git diff-tree` query for each recovered commit, filters changed
paths whose basenames end in `_test.go`, sorts those repository-relative paths,
and attaches them to the same citation. The CLI passes the complete citation set
to `risk.Evaluate`, which emits the stable evidence-backed risk level before the
formatter writes either terminal output or a report. With `--report PATH`, the
CLI maps the validated request and evidence into `report.Data`, renders a
deterministic standalone HTML document, atomically writes it, and prints only a
stable report confirmation. A validated diff still becomes a stable summary
report without a risk claim because diff evidence is not yet implemented. Input,
evidence, and report-write failures are deterministic non-zero errors.

## Planned core flow

CLI entry point → input/range or diff parser → local Git evidence collector →
commit/issue/PR/test citation builder → deterministic risk evaluator → terminal
or static HTML report formatter.

Update this map whenever execution crosses a new file or function.
