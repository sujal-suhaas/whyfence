# Architecture

## Purpose

whyfence explains why a local Git-tracked code range or diff exists before it is
removed or rewritten. The core must be deterministic and offline.

## Current system map

- `cmd/whyfence/main.go` is the minimal Go CLI entry point and currently owns
  the `--version` behavior.
- `cmd/whyfence/main_test.go` protects the version behavior without external
  dependencies.
- `Dockerfile` is the reproducible test boundary; it pins Go and installs Git
  explicitly for the future local-history collector.
- `.github/workflows/docker-test.yml` builds the image and runs the tests and
  version check with container networking disabled.

## Planned boundaries

- CLI input accepts a file range or diff.
- Local Git history is the only evidence source for the core.
- Evidence is cited, then a deterministic risk flag is produced.
- Optional LLM summaries remain outside the deterministic core and are disabled
  in tests.

Update this map when files, components, or boundaries change.
