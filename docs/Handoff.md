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
- Merged issue #6 through PR #14 as squash commit `b1abebc8`; its remote branch
  was deleted.
- Pulled `main` and reran the full Docker, line-range, and accuracy checks after
  the merge.
- Started issue #7 on `feat/7-reference-evidence` with explicit local issue/PR
  extraction attached to originating line-range citations.
- Merged issue #7 through PR #15 as squash commit `6bf4cc84`; its remote branch
  was deleted.
- Pulled `main` and reran the full Docker, reference, and accuracy checks after
  the merge.
- Started issue #8 on `feat/8-changed-tests` with root-safe same-commit test-path
  recovery, matching/non-matching fixtures, and CLI test citations.
- Merged issue #8 through PR #16 as squash commit `aba23165`; its remote branch
  was deleted.
- Pulled `main` and reran the full Docker, changed-test, and accuracy checks after
  the merge.
- Started issue #9 on `feat/9-risk-flag` with evidence-backed low, medium, and
  high risk evaluation and generated risk fixtures.
- Merged issue #9 through PR #17 as squash commit `80ac53dc`; its remote branch
  was deleted.
- Pulled `main` and reran the post-merge Docker build, full suite, risk tests,
  accuracy, and CLI contract checks.
- Started issue #18 on `feat/18-static-html-report` for deterministic standalone
  HTML report output.
- Merged issue #18 through PR #19 as squash commit `bad17659`; its remote branch
  was deleted.
- Pulled `main` and reran the full Docker, report, accuracy, version, and help
  checks after the merge.

## In progress

- Issues #2 and #3 are merged and verified on `main`.
- Issue #4 is merged and verified on `main`.
- Issue #5 is merged and verified on `main`.
- Issue #6 is merged and verified on `main`.
- Issue #7 is merged and verified on `main`.
- Issue #8 is merged and verified on `main`.
- Issue #9 is merged and verified on `main`.
- Issue #18 is merged and verified on `main`.
- The worktree was clean at the start of issue #18; no unrelated changes are
  included in this feature.

## Next

- After issue #18 is merged and verified, review the next approved deliverable
  before requesting or using another issue-level GO.

## Gotchas

- Docker builds need network access for `apt-get install git`; test containers
  run with `--network=none`.
- Docker Desktop may need to be started before local Docker verification.
- The post-merge Docker checks run with `--network=none`; the image build needs
  network access for `apt-get install git`.
- README now documents the tested Go/Docker setup and report commands.

- Docker daemon was restarted locally during issue #9 verification because the
  socket was unavailable; subsequent post-merge build, tests, version, and help
  checks passed.

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
- `git pull --ff-only origin main`: passed after PR #14 merge.
- `docker build --no-cache -t whyfence:test .` after PR #14 merge: passed.
- `docker run --rm --network=none whyfence:test go test ./... -count=1` after
  PR #14 merge: passed.
- `docker run --rm --network=none whyfence:test go test ./... -run
  TestLineRange -count=1 -v` after PR #14 merge: passed.
- `docker run --rm --network=none whyfence:test go test ./... -run
  TestRecoveryAccuracy -count=1 -v` after PR #14 merge: passed; output
  `recovered 4 of 4`.
- `docker run --rm --network=none whyfence:test go run ./cmd/whyfence
  --version` after PR #14 merge: passed; output `whyfence 0.1.0-dev`.
- `docker run --rm --network=none whyfence:test go run ./cmd/whyfence --help`
  after PR #14 merge: passed.
- Docker-mounted Go 1.23.0 `go test ./... -run
  'Test(Issue|PR)Reference' -count=1 -v`: passed; issue and PR references were
  recovered from their originating fixture commits.
- Docker-mounted Go 1.23.0 `go test ./... -run TestLineRange -count=1 -v` on
  issue #7: passed with citation references attached.
- Docker-mounted Go 1.23.0 `go test ./... -run TestRecoveryAccuracy -count=1
  -v` on issue #7: passed; output `recovered 6 of 6`.
- `git pull --ff-only origin main`: passed after PR #15 merge.
- `docker build --no-cache -t whyfence:test .` after PR #15 merge: passed.
- `docker run --rm --network=none whyfence:test go test ./... -count=1` after
  PR #15 merge: passed.
- `docker run --rm --network=none whyfence:test go test ./... -run
  TestIssueReference -count=1 -v` after PR #15 merge: passed.
- `docker run --rm --network=none whyfence:test go test ./... -run
  TestPRReference -count=1 -v` after PR #15 merge: passed.
- `docker run --rm --network=none whyfence:test go test ./... -run
  TestRecoveryAccuracy -count=1 -v` after PR #15 merge: passed; output
  `recovered 6 of 6`.
- `docker run --rm --network=none whyfence:test go run ./cmd/whyfence
  --version` after PR #15 merge: passed; output `whyfence 0.1.0-dev`.
- `docker run --rm --network=none whyfence:test go run ./cmd/whyfence --help`
  after PR #15 merge: passed.
- Docker-mounted Go 1.23.0 `go test ./... -run TestChangedTests -count=1 -v`:
  passed; matching `_test.go` paths, root-commit paths, and code-only changes
  were classified deterministically.
- Docker-mounted Go 1.23.0 `go test ./... -run TestLineRange -count=1 -v` on
  issue #8: passed with same-commit test paths attached to citations.
- Docker-mounted Go 1.23.0 `go test ./... -run TestRecoveryAccuracy -count=1
  -v` on issue #8: passed; output `recovered 7 of 7`.
- `git pull --ff-only origin main`: passed after PR #16 merge.
- `docker build --no-cache -t whyfence:test .` after PR #16 merge: passed.
- `docker run --rm --network=none whyfence:test go test ./... -count=1` after
  PR #16 merge: passed.
- `docker run --rm --network=none whyfence:test go test ./... -run
  TestChangedTests -count=1 -v` after PR #16 merge: passed.
- `docker run --rm --network=none whyfence:test go test ./... -run
  TestRecoveryAccuracy -count=1 -v` after PR #16 merge: passed; output
  `recovered 7 of 7`.
- `docker run --rm --network=none whyfence:test go run ./cmd/whyfence
  --version` after PR #16 merge: passed; output `whyfence 0.1.0-dev`.
- `docker run --rm --network=none whyfence:test go run ./cmd/whyfence --help`
  after PR #16 merge: passed.
- Go 1.23.0 `go test ./... -run
  'TestEvaluateRiskLevels|TestEvaluateDoesNotMutateCitations|TestRiskFixtureBuilder'
  -count=1 -v`: passed.
- Go 1.23.0 `go test ./... -run TestRecoveryAccuracy -count=1 -v`: passed;
  output `recovered 12 of 12`.
- Go 1.23.0 `go test ./... -count=1`: passed on issue #9.
- Docker `docker build --no-cache -t whyfence:test .`: passed on issue #9.
- Docker `docker run --rm --network=none whyfence:test git --version`: passed;
  output `git version 2.39.5`.
- Docker `docker run --rm --network=none whyfence:test go test ./... -count=1`:
  passed on issue #9.
- Docker `docker run --rm --network=none whyfence:test go test ./... -run
  TestRecoveryAccuracy -count=1 -v`: passed; output `recovered 12 of 12`.
- Docker `docker run --rm --network=none whyfence:test go test ./... -run
  TestCLIContract -count=1 -v`: passed.
- `git pull --ff-only origin main`: passed after PR #17 merge.
- Docker `docker build --no-cache -t whyfence:test .` after PR #17 merge:
  passed.
- Docker `docker run --rm --network=none whyfence:test go test ./... -count=1`
  after PR #17 merge: passed.
- Docker `docker run --rm --network=none whyfence:test go test ./... -run
  'TestEvaluateRiskLevels|TestEvaluateDoesNotMutateCitations|TestRiskFixtureBuilder'
  -count=1 -v` after PR #17 merge: passed.
- Docker `docker run --rm --network=none whyfence:test go test ./... -run
  TestRecoveryAccuracy -count=1 -v` after PR #17 merge: passed; output
  `recovered 12 of 12`.
- Docker `docker run --rm --network=none whyfence:test go test ./... -run
  TestCLIContract -count=1 -v` after PR #17 merge: passed.
- Post-merge Docker `go run ./cmd/whyfence --version` and `--help` retries:
  blocked by SIGBUS in the local container runtime, followed by sandbox I/O
  failure; not claimed as passed.
- Docker `docker run --rm --network=none whyfence:test go run ./cmd/whyfence
  --version`: passed; output `whyfence 0.1.0-dev`.
- Docker `docker run --rm --network=none whyfence:test go run ./cmd/whyfence
  --help`: passed; input flags were listed.
- `git pull --ff-only origin main`: passed after PR #19 merge.
- Docker `docker build --no-cache -t whyfence:test .` after PR #19 merge:
  passed.
- Docker `docker run --rm --network=none whyfence:test go test ./... -count=1`
  after PR #19 merge: passed.
- Docker `docker run --rm --network=none whyfence:test go test ./... -run
  'TestHTMLReport|TestDiffReportOmitsRisk|TestWriteAtomicOverwritesReport'
  -count=1 -v` after PR #19 merge: passed.
- Docker `docker run --rm --network=none whyfence:test go test ./... -run
  TestCLIReport -count=1 -v` after PR #19 merge: passed.
- Docker `docker run --rm --network=none whyfence:test go test ./... -run
  TestRecoveryAccuracy -count=1 -v` after PR #19 merge: passed; output
  `recovered 12 of 12`.
- Docker `docker run --rm --network=none whyfence:test go run ./cmd/whyfence
  --version` after PR #19 merge: passed; output `whyfence 0.1.0-dev`.
- Docker `docker run --rm --network=none whyfence:test go run ./cmd/whyfence
  --help` after PR #19 merge: passed; `--report` was listed.
