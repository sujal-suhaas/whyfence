# Handoff

## Done

- Issue #2 is committed and pushed on `feat/2-pin-go-toolchain` as `a5c0d1a`.
- Started issue #3 on `feat/3-fixture-builder`.
- Added a standard-library fixture builder that creates temporary Git repos with
  fixed authors, dates, known target/test history, and cleanup.
- Updated Docker to copy and test the new internal fixture package.

## In progress

- Issue #3 is implemented and verified; it is awaiting diff approval before
  commit and push.
- No PR has been opened for issue #2 or #3.

## Next

- Show the issue #3 diff and wait for explicit approval before committing.
- After approval, commit and push `feat/3-fixture-builder`.

## Gotchas

- Docker builds need network access for `apt-get install git`; test containers
  run with `--network=none`.
- The sandbox Docker daemon uses the `vfs` storage driver.
- `docs/Bug.md` and `docs/Decisions.md` remain unrelated unstaged documentation
  from the prior skeleton setup; do not include them in the issue #3 commit.
- The README install section remains intentionally unchanged.

## Verification

- `go test ./... -run TestFixtureBuilder -count=1`: passed.
- `go test ./... -count=1`: passed.
- `docker build --no-cache -t whyfence:test .`: passed.
- `docker run --rm --network=none whyfence:test git --version`: passed.
- `docker run --rm --network=none whyfence:test go test ./... -count=1`: passed.
- `docker run --rm --network=none whyfence:test go test ./... -run
  TestFixtureBuilder -count=1`: passed.
