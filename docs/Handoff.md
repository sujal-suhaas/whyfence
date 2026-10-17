# Handoff

## Done

- Merged issue #2 through PR #10 as squash commit `73a2f614`.
- Pulled `main` and verified a clean no-cache Docker rebuild, Git in the image,
  offline container tests, and the offline version flag after the merge.
- Merged issue #3 through PR #11 as squash commit `0e83c7e0`; its remote branch
  was deleted.
- Pulled `main` and reran the full Docker and offline fixture checks after the
  merge.
- Started issue #4 on `feat/4-accuracy-harness` with a calculated harness and
  three generated-fixture cases.

## In progress

- Issues #2 and #3 are merged and verified on `main`.
- Issue #4 is implemented and verified in Docker; it is awaiting diff approval
  before commit and push.
- Untracked `docs/Bug.md` and `docs/Decisions.md` remain from the documentation
  skeleton setup and are intentionally untouched.

## Next

- Show the issue #4 diff and wait for explicit approval before committing.
- After approval, commit and push `feat/4-accuracy-harness` and open its PR only
  after reviewing the exclusive diff.

## Gotchas

- Docker builds need network access for `apt-get install git`; test containers
  run with `--network=none`.
- Docker Desktop may need to be started before local Docker verification.
- The post-merge Docker checks run with `--network=none`; the image build needs
  network access for `apt-get install git`.
- `docs/Bug.md` and `docs/Decisions.md` remain unrelated unstaged documentation
  from the prior skeleton setup and are not part of issue #3.
- The README install section remains intentionally unchanged.

## Verification

- `git pull --ff-only origin main`: passed after PR #10 merge.
- `docker build --no-cache -t whyfence:test .`: passed on merged `main`.
- `docker run --rm --network=none whyfence:test git --version`: passed.
- `docker run --rm --network=none whyfence:test go test ./...`: passed.
- `docker run --rm --network=none whyfence:test go run ./cmd/whyfence --version`:
  passed; output `whyfence 0.1.0-dev`.
- `docker build --no-cache -t whyfence:test .` after PR #11 merge: passed.
- `docker run --rm --network=none whyfence:test go test ./... -count=1` after
  PR #11 merge: passed.
- `docker run --rm --network=none whyfence:test go test ./... -run
  TestFixtureBuilder -count=1` after PR #11 merge: passed.
- Docker-mounted Go 1.23.0 `go test ./... -run TestRecoveryAccuracy -count=1
  -v`: passed; output `recovered 3 of 3`.
- Docker-mounted Go 1.23.0 `go test ./... -count=1`: passed.
