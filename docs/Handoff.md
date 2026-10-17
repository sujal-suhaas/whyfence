# Handoff

## Done

- Merged issue #2 through PR #10 as squash commit `73a2f614`.
- Pulled `main` and verified a clean no-cache Docker rebuild, Git in the image,
  offline container tests, and the offline version flag after the merge.
- Merged issue #3 through PR #11 as squash commit `0e83c7e0`; its remote branch
  was deleted.
- Pulled `main` and reran the full Docker and offline fixture checks after the
  merge.
- Merged issue #4 through PR #12 as squash commit `bddf6615`; its remote branch
  was deleted.
- Pulled `main` and reran the full Docker, fixture, and accuracy checks after the
  merge.
- Merged issue #5 through PR #13 as squash commit `32ae790d`; its remote branch
  was deleted.
- Pulled `main` and reran the full Docker and CLI contract checks after the
  merge.
- Started issue #6 on `feat/6-line-range-evidence` with production Git
  line-range recovery, cited commit output, and an accuracy case.

## In progress

- Issues #2 and #3 are merged and verified on `main`.
- Issue #4 is merged and verified on `main`.
- Issue #5 is merged and verified on `main`.
- Issue #6 is implemented and verified in its branch; its commit, PR, merge, and
  post-merge lifecycle remain in progress under the issue-level GO workflow.
- Untracked `docs/Bug.md` remains from the documentation skeleton setup and is
  intentionally untouched.

## Next

- After issue #6 is merged and verified, review issue #7's plan before requesting
  or using its issue-level GO.

## Gotchas

- Docker builds need network access for `apt-get install git`; test containers
  run with `--network=none`.
- Docker Desktop may need to be started before local Docker verification.
- The post-merge Docker checks run with `--network=none`; the image build needs
  network access for `apt-get install git`.
- `docs/Bug.md` remains unrelated untracked documentation from the prior
  skeleton setup and is not part of issue #6.
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
- `docker build --no-cache -t whyfence:test .` after PR #12 merge: passed.
- `docker run --rm --network=none whyfence:test go test ./... -count=1` after
  PR #12 merge: passed.
- `docker run --rm --network=none whyfence:test go test ./... -run
  TestRecoveryAccuracy -count=1 -v` after PR #12 merge: passed; output
  `recovered 3 of 3`.
- `docker run --rm --network=none whyfence:test go run ./cmd/whyfence
  --version` after PR #12 merge: passed; output `whyfence 0.1.0-dev`.
- Docker-mounted Go 1.23.0 `go test ./... -run TestCLIContract -count=1 -v`:
  passed, including range, diff-file, diff-stdin, and deterministic errors.
- `docker build --no-cache -t whyfence:test .` on issue #5: passed.
- `docker run --rm --network=none whyfence:test go test ./... -count=1`: passed.
- `docker run --rm --network=none whyfence:test go test ./... -run
  TestCLIContract -count=1 -v`: passed.
- `docker run --rm --network=none whyfence:test go run ./cmd/whyfence --help`:
  passed; input flags were listed.
- `git pull --ff-only origin main`: passed after PR #13 merge.
- `docker build --no-cache -t whyfence:test .` after PR #13 merge: passed.
- `docker run --rm --network=none whyfence:test go test ./... -count=1` after
  PR #13 merge: passed.
- `docker run --rm --network=none whyfence:test go test ./... -run
  TestCLIContract -count=1 -v` after PR #13 merge: passed.
- `docker run --rm --network=none whyfence:test go run ./cmd/whyfence --help`
  after PR #13 merge: passed.
- Docker-mounted Go 1.23.0 `go test ./... -run TestLineRange -count=1 -v`:
  passed; the production collector recovered both known commits and rejected
  non-repository and invalid-boundary inputs deterministically.
- Docker-mounted Go 1.23.0 `go test ./... -run TestRecoveryAccuracy -count=1
  -v`: passed; output `recovered 4 of 4`.
