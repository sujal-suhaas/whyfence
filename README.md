# whyfence

whyfence is a local, git-only CLI that helps explain why code exists before it
is deleted or rewritten. Given a file range or diff, it will produce cited
evidence from local Git history, including related commits, issues, pull
requests, and tests changed alongside the code, plus a risk flag.

The core is designed to be deterministic and fully offline. Any LLM summary is
optional and disabled in tests.

## Project status

The offline CLI core and self-contained static HTML report are implemented and
tested. The GitHub Pages landing page is still a later deliverable. No package is
published; use the Go toolchain or Docker commands below.

The planned deliverables are intentionally ordered:

1. CLI core.
2. Offline fixture-repository tests and Docker rebuild checks.
3. Self-contained static HTML reports.
4. A plain static landing page in `site/`.

## Setup

Requirements:

- Go 1.23.0
- Git
- Docker, for the reproducible offline checks

The project uses only the Go standard library. From the repository root:

```sh
go version
go test ./... -count=1
```

## Run

Analyze an inclusive file range from a local Git repository:

```sh
go run ./cmd/whyfence --range internal/worker.go:3-3
```

Analyze a unified diff:

```sh
go run ./cmd/whyfence --diff change.patch
cat change.patch | go run ./cmd/whyfence --diff -
```

## Static HTML reports

Write a deterministic, standalone report to a named file:

```sh
go run ./cmd/whyfence --range internal/worker.go:3-3 --report report.html
go run ./cmd/whyfence --diff change.patch --report diff-report.html
```

Reports contain inline CSS and no network-dependent assets or scripts. Range
reports include local Git citations, changed tests, explicit issue/PR references,
and the evidence-backed risk level. Diff reports include only the accepted diff
summary; diff history evidence and risk are not available yet.

## Test and Docker rebuild

Run the complete offline test suite:

```sh
go test ./... -count=1
```

Rebuild and test through the pinned Go 1.23.0 Docker boundary:

```sh
docker build --no-cache -t whyfence:test .
docker run --rm --network=none whyfence:test go test ./... -count=1
```

## Deploy

There is no deployment target yet. The generated `.html` report is a standalone
artifact that can be opened locally or uploaded to a static host. The `site/`
landing page is planned as the next deliverable.

## License

MIT. See [LICENSE](LICENSE).
