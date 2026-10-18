# Architecture

## Purpose

whyfence explains why a local Git-tracked code range exists before it is
removed or rewritten. It also accepts a unified diff and reports its validated
file/hunk summary. The core is deterministic and offline.

## Current system map

- `cmd/whyfence/main.go` is the CLI entry point. It parses `--version`, exactly
  one of `--range` or `--diff`, and the optional named `--report` output, then
  orchestrates validation, range evidence, risk evaluation, and formatting.
- `cmd/whyfence/main_test.go` protects the version/help, range, diff, error, and
  report contracts without external dependencies.
- `internal/input/` validates the inclusive file-range and unified-diff input
  contracts. Its direct parser tests cover boundaries; CLI tests cover the
  user-visible contract.
- `internal/fixtures/` builds deterministic temporary Git repositories at test
  time for evidence and generated-fixture regression checks, including matching
  and non-matching test changes.
- `internal/accuracy/` is a generated-fixture regression harness. It counts
  exact expected-reason matches from test-time resolvers and production
  line-range recovery; it is not a benchmark of accuracy on arbitrary public
  repositories.
- `internal/evidence/` invokes local Git's line-range history traversal, parses
  full-hash/subject citations, extracts supported local issue/PR references,
  and identifies same-commit changed test paths. Evidence is range-only; diff
  input does not query history.
- `internal/risk/` reduces recovered range citations to a deterministic,
  evidence-backed low, medium, or high removal-risk level without mutating the
  evidence. It is not a safety, runtime-impact, ownership, or business-
  criticality judgment.
- `internal/report/` renders validated analysis as deterministic standalone
  HTML with inline CSS and no network or external asset dependencies. Range
  reports include evidence and risk; diff reports include only the accepted
  summary and explicitly omit evidence/risk claims.
- Without `--report`, the CLI writes stable terminal output. With a named
  `--report PATH`, it atomically writes the HTML artifact and prints only a
  stable confirmation.
- `site/` is a separate static presentation boundary. It contains the plain
  landing page, terminal demo, provenance notes, and checked-in report
  snapshots. `.github/workflows/pages.yml` checks out the repository and
  uploads `site/` directly to GitHub Pages; it does not build or run the CLI.
- `Dockerfile` is the reproducible core-test boundary. It pins Go, installs Git,
  copies only `go.mod`, `cmd/`, and `internal/`, and does not include `site/`.
  `.github/workflows/docker-test.yml` repeats the core checks with container
  networking disabled.

## System boundaries

- **Range input:** validate the path and lines, recover cited local Git history,
  collect test/reference signals, evaluate evidence-backed risk, and format
  terminal output or a standalone report.
- **Diff input:** validate and summarize the source, file count, and hunk count;
  history evidence and risk are not implemented for this mode.
- **Optional report output:** named HTML only; no server, timestamp, JavaScript,
  external assets, or network access is required.
- **Optional LLM summaries:** outside the deterministic core and disabled in
  tests; no LLM is required for the current flow.

Update this map whenever files, components, or boundaries change.
