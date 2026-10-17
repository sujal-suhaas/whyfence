# Handoff

## Done

- Added issue #2 on `feat/2-pin-go-toolchain`: pinned Go 1.23.0, a minimal
  `--version` CLI, one smoke test, a Git-installing Dockerfile, and offline
  Docker test workflow.
- Verified host Go tests, a clean no-cache Docker rebuild, Git in the image,
  offline container tests, and the offline version flag.
- Updated the architecture, flow, and feature records for issue #2.

## In progress

- Issue #2 is ready to commit and push; no PR has been opened yet.
- Issue #3 is open for deterministic temporary Git fixtures.

## Next

- Commit and push issue #2, then open its PR only after reviewing the diff.
- Implement issue #3 on `feat/3-fixture-builder` after explicit approval.

## Gotchas

- Docker builds need network access for `apt-get install git`; test containers
  run with `--network=none`.
- The sandbox Docker daemon uses the `vfs` storage driver.
- `docs/Bug.md` and `docs/Decisions.md` remain unrelated unstaged documentation
  from the prior skeleton setup; do not include them in the issue #2 commit.
- The README install section remains intentionally unchanged.

## Verification

- `go test ./...`: passed.
- `docker build --no-cache -t whyfence:test .`: passed.
- `docker run --rm --network=none whyfence:test git --version`: passed.
- `docker run --rm --network=none whyfence:test go test ./...`: passed.
- `docker run --rm --network=none whyfence:test go run ./cmd/whyfence --version`:
  passed; output `whyfence 0.1.0-dev`.
