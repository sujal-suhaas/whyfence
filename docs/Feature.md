# Feature log

Append-only. Add one entry per feature with the goal, approach, files,
verification, and status.

## Entry template

- **Goal:**
- **Approach:**
- **Files touched:**
- **Verification:**
- **Status:**

## Issue #2 — pin Go toolchain and establish Docker gate

- **Goal:** Pin Go 1.23.0, provide a minimal `--version` CLI, and establish a
  Git-installed Docker test boundary.
- **Approach:** Use only the standard library, run the smoke test with
  `go test ./...`, and run CI test containers with `--network=none`.
- **Files touched:** `go.mod`, `cmd/whyfence/`, `Dockerfile`,
  `.github/workflows/docker-test.yml`, `docs/Architecture.md`, `docs/Flow.md`,
  `docs/Feature.md`.
- **Verification:** `go test ./...` passed; `docker build -t whyfence:test .`
  passed; container Git check, `docker run --rm --network=none whyfence:test go
  test ./...`, and offline `--version` check all passed.
- **Status:** Merged via PR #10 as squash commit `73a2f614`.

## Issue #3 — build temporary Git fixtures

- **Goal:** Generate deterministic temporary fixture repositories with fixed
  authors and dates at test time.
- **Approach:** Use a standard-library Go helper to run `git init` in a system
  temp directory, configure a fixed identity, create known target/test history,
  and clean the repository after each test.
- **Files touched:** `internal/fixtures/`, `Dockerfile`,
  `docs/Architecture.md`, `docs/Flow.md`, `docs/Feature.md`.
- **Verification:** `go test ./... -run TestFixtureBuilder -count=1` and
  `go test ./... -count=1` passed; clean `docker build --no-cache -t
  whyfence:test .` passed; Git, full tests, and the fixture test passed in
  `--network=none` containers.
- **Status:** Merged via PR #11 as squash commit `0e83c7e0`.

## Issue #4 — start recovered-N-of-M accuracy harness

- **Goal:** Start a calculated accuracy harness with three known generated-fixture
  cases and stable `recovered N of M` output.
- **Approach:** Keep the harness resolver-agnostic; use test-only Git resolvers
  over generated history until production evidence recovery lands.
- **Files touched:** `internal/accuracy/`, `docs/Architecture.md`,
  `docs/Flow.md`, `docs/Feature.md`, `docs/Handoff.md`.
- **Verification:** Docker-mounted Go 1.23.0 ran `go test ./... -run
  TestRecoveryAccuracy -count=1 -v` with `recovered 3 of 3`, and
  `go test ./... -count=1` passed.
- **Status:** Implementation ready for review; commit pending approval.
