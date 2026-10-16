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
- **Status:** Committed on `feat/2-pin-go-toolchain`; PR and CI review pending.
