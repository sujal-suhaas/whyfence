# whyfence

whyfence is a local, git-only CLI that helps explain why code exists before it
is deleted or rewritten. Given a file range, it produces cited evidence from
local Git history, including related commits, supported issue/PR references,
unclassified GitHub squash suffixes, and tests changed alongside the code, plus
an evidence-backed risk flag. Given a unified diff, it validates and summarizes
the source, changed-file count, and hunk count; diff history evidence and risk
are not implemented yet.

The core is designed to be deterministic and fully offline. Any LLM summary is
optional and disabled in tests.

## Project status

The offline CLI core, self-contained static HTML report, and plain GitHub Pages
landing page are implemented and tested. No package is published; use the Go
toolchain or Docker commands below.

The completed deliverables were intentionally ordered:

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

Analyze an inclusive file range from this or another local Git repository:

```sh
go run ./cmd/whyfence --range cmd/whyfence/main.go:26-32
```

Analyze a unified diff:

```sh
go run ./cmd/whyfence --diff change.patch
cat change.patch | go run ./cmd/whyfence --diff -
```

## Static HTML reports

Write a deterministic, standalone report to a named file:

```sh
go run ./cmd/whyfence --range cmd/whyfence/main.go:26-32 --report report.html
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

The generated `.html` report is a standalone artifact that can be opened locally
or uploaded to a static host. The checked-in `site/` landing page is published
by `.github/workflows/pages.yml` without a build step.

To enable GitHub Pages for a fork or another copy of the repository:

1. Push `main` and open **Settings → Pages**.
2. Select **GitHub Actions** as the build and deployment source.
3. Push a subsequent change to `main`, or run **Deploy static site** from the
   repository's Actions tab.

The workflow uploads only `site/`; the Go CLI and Docker test image do not run
as part of the Pages deployment.

## License

MIT. See [LICENSE](LICENSE).
