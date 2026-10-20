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

The repository uses only the Go standard library. `go.mod` is the complete
dependency manifest and lock state: there are no third-party modules, so Go does
not need a `go.sum` checksum file. The exact Go runtime is pinned by `go 1.23.0`
and the Docker image tag.

## Project status

The offline CLI core, self-contained static HTML report, and plain GitHub Pages
landing page are implemented and tested. No package is published; use the Go
toolchain or Docker commands below.

The completed deliverables were intentionally ordered:

1. CLI core.
2. Offline fixture-repository tests and Docker rebuild checks.
3. Self-contained static HTML reports.
4. A plain static landing page in `site/`.

## Setup from a clean Linux checkout

Requirements:

- Go 1.23.0
- Git
- Docker, for the reproducible offline build/run checks

Clone the repository and enter its root:

```sh
git clone https://github.com/sujal-suhaas/whyfence.git
cd whyfence
```

Verify the pinned toolchain and run the one project test command:

```sh
go version  # use Go 1.23.0
git --version
go test ./... -count=1
```

No `.env` file is required. `.env.example` documents that the CLI has no
runtime settings; never put credentials in project files.

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

## Build and run a local binary

Build the CLI from the repository root:

```sh
mkdir -p bin
go build -trimpath -o bin/whyfence ./cmd/whyfence
./bin/whyfence --help
```

Run a range analysis with the built binary:

```sh
./bin/whyfence --range cmd/whyfence/main.go:26-32
```

## Test and Docker build/run

The one test command is:

```sh
go test ./... -count=1
```

Build the pinned Docker image. The image runs `whyfence --help` by default and
contains the compiled CLI binary:

```sh
docker build --no-cache -t whyfence:local .
docker run --rm --network=none whyfence:local
docker run --rm --network=none whyfence:local whyfence --version
docker run --rm --network=none whyfence:local whyfence --range cmd/whyfence/main.go:26-32
```

The Docker build itself runs `go test ./... -count=1` before producing the
binary. To rerun the same test command explicitly inside the image:

```sh
docker run --rm --network=none whyfence:local go test ./... -count=1
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
