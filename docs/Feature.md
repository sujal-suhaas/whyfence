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

- **Goal:** Start a calculated generated-fixture regression harness with three
  known cases and stable `recovered N of M` output.
- **Approach:** Keep the harness resolver-agnostic; use test-only Git resolvers
  over generated history until production evidence recovery lands.
- **Files touched:** `internal/accuracy/`, `docs/Architecture.md`,
  `docs/Flow.md`, `docs/Feature.md`, `docs/Handoff.md`.
- **Verification:** Docker-mounted Go 1.23.0 ran `go test ./... -run
  TestRecoveryAccuracy -count=1 -v` with all 3 generated fixture cases recovered,
  and `go test ./... -count=1` passed.
- **Status:** Merged via PR #12 as squash commit `bddf6615`.

## Issue #5 — define deterministic CLI contract

- **Goal:** Accept file ranges or unified diffs and emit stable terminal output
  with deterministic non-zero errors for invalid input.
- **Approach:** Use explicit `--range path:start-end` and `--diff path`/`-`
  modes; validate the input in `internal/input` before formatting a summary.
- **Files touched:** `cmd/whyfence/`, `internal/input/`,
  `docs/Architecture.md`, `docs/Flow.md`, `docs/Decisions.md`,
  `docs/Feature.md`, `docs/Handoff.md`.
- **Verification:** Docker-mounted Go 1.23.0 passed `TestCLIContract`, its
  deterministic error cases, and the full suite; clean Docker build, Git check,
  offline full tests, contract tests, `--version`, and `--help` all passed.
- **Status:** Merged via PR #13 as squash commit `32ae790d`.

## Issue #6 — recover line-range evidence

- **Goal:** Trace a targeted file range through local Git history and emit
  citations containing commit hashes and subjects.
- **Approach:** Resolve the file's Git repository, run deterministic
  `git log -L`, parse and de-duplicate full-hash/subject records, and wire the
  range path into the CLI while leaving diff evidence for its own issue slice.
- **Files touched:** `cmd/whyfence/`, `internal/evidence/`,
  `internal/accuracy/`, `docs/Architecture.md`, `docs/Flow.md`,
  `docs/Decisions.md`, `docs/Feature.md`, `docs/Handoff.md`.
- **Verification:** Docker-mounted Go 1.23.0 passed `TestLineRange` and the
  generated-fixture regression harness with all 4 known cases recovered;
  full-suite and offline Docker checks passed before and after merge.
- **Status:** Merged via PR #14 as squash commit `b1abebc8`.

## Issue #7 — extract local issue and PR references

- **Goal:** Surface deterministic issue and pull request references from local
  commit subjects without using the GitHub API.
- **Approach:** Extend line-range citations with provenance-preserving
  references; parse explicit issue, PR, merge-subject, and GitHub URL forms while
  intentionally leaving ambiguous bare `#number` tokens unresolved.
- **Files touched:** `cmd/whyfence/`, `internal/evidence/`,
  `internal/fixtures/`, `internal/accuracy/`, `docs/Architecture.md`,
  `docs/Flow.md`, `docs/Decisions.md`, `docs/Feature.md`, `docs/Handoff.md`.
- **Verification:** Docker-mounted Go 1.23.0 passed `TestIssueReference`,
  `TestPRReference`, `TestLineRange`, and the generated-fixture regression
  harness with all 6 known cases recovered; full-suite and offline Docker checks
  passed before and after merge.
- **Status:** Merged via PR #15 as squash commit `6bf4cc84`.

## Issue #8 — detect tests changed alongside code

- **Goal:** Identify test files changed in the same commits as recovered target
  code and preserve those test citations.
- **Approach:** Query each recovered commit with root-safe local `git diff-tree`,
  classify only repository-relative basenames ending in `_test.go`, sort paths,
  attach them to the originating citation, and render `test:` lines in range
  output. The fixture includes matching test changes and a code-only
  non-matching change.
- **Files touched:** `cmd/whyfence/`, `internal/evidence/`,
  `internal/fixtures/`, `internal/accuracy/`, `docs/Architecture.md`,
  `docs/Flow.md`, `docs/Decisions.md`, `docs/Feature.md`, `docs/Handoff.md`.
- **Verification:** Docker-mounted Go 1.23.0 passed `TestChangedTests`,
  `TestLineRange`, and the generated-fixture regression harness with all 7 known
  cases recovered; full-suite and offline Docker checks passed before and after
  merge.
- **Status:** Merged via PR #16 as squash commit `aba23165`.

## Issue #9 — add deterministic risk flag

- **Goal:** Classify recovered range evidence as low, medium, or high removal
  risk and finish the generated-fixture regression suite.
- **Approach:** Add an order-independent `internal/risk` evaluator using changed
  test paths and explicit references as the two evidence signals; generate
  isolated low, medium, and high fixture histories plus an aggregate high case;
  and render `risk: <level>` for range output only.
- **Files touched:** `cmd/whyfence/`, `internal/risk/`, `internal/fixtures/`,
  `internal/accuracy/`, `docs/Architecture.md`, `docs/Flow.md`,
  `docs/Decisions.md`, `docs/Feature.md`, `docs/Handoff.md`.
- **Verification:** Go 1.23.0 and Docker passed risk unit tests, generated risk
  fixtures, and the generated-fixture regression suite with all 12 known cases
  recovered, plus the full suite, CLI contract checks, no-cache image rebuild,
  Git availability, offline version, and offline help checks before merge. The
  first post-merge version/help retry hit a local Docker SIGBUS followed by a
  sandbox I/O failure; a later retry passed both checks.
- **Status:** Merged via PR #17 as squash commit `80ac53dc`; the initial local
  post-merge smoke retry was blocked by the container runtime, but a later retry
  passed and is recorded in `docs/Handoff.md`.

## Issue #18 — generate self-contained static HTML reports

- **Goal:** Write deterministic, standalone HTML reports from the existing range
  and diff analysis results without a server or network access.
- **Approach:** Add `--report PATH`, map validated input/evidence into an
  `internal/report` data model, render inline-CSS HTML with `html/template`, and
  atomically write a named file. Range reports include risk/evidence; diff
  reports remain summary-only with no risk claim.
- **Files touched:** `cmd/whyfence/`, `internal/report/`, `README.md`,
  `docs/Architecture.md`, `docs/Flow.md`, `docs/Decisions.md`,
  `docs/Feature.md`, `docs/Handoff.md`.
- **Verification:** Docker-mounted Go 1.23.0 passed `TestHTMLReport`,
  `TestCLIReport`, `TestCLIContract`, `TestRecoveryAccuracy`, the full suite,
  no-cache image rebuild, Git availability, offline version/help, and
  `git diff --check`; the same report, full-suite, accuracy, version, and help
  checks passed after merge.
- **Status:** Merged via PR #19 as squash commit `bad17659`.

## Issue #20 — build the static GitHub Pages landing page

- **Goal:** Publish a plain static project landing page with an accurate CLI
  explanation, terminal demo, tested setup section, and offline report gallery.
- **Approach:** Keep HTML and CSS local in `site/`, avoid runtime JavaScript or
  CDN dependencies, add five CLI-generated standalone report snapshots, record
  each source repository/commit/range/license in `site/reports/sources.md`, and
  upload `site/` directly through a no-build GitHub Pages workflow.
- **Files touched:** `site/index.html`, `site/styles.css`, `site/demo/terminal.txt`,
  `site/reports/*.html`, `site/reports/sources.md`, `.github/workflows/pages.yml`,
  `internal/report/`, `README.md`, `docs/Decisions.md`, `docs/Feature.md`,
  `docs/Handoff.md`.
- **Verification:** Static Python smoke check passed with five standalone reports
  and twenty local landing-page links; Docker-mounted Go 1.23.0 passed
  `go test ./... -count=1` and the documented `go run ... --report` demo;
  `git diff --check` passed. The same Docker/site checks passed after merge, and
  the GitHub Pages workflow deployed the site successfully after issue #22
  enabled the repository Pages configuration.
- **Rollback:** Revert the issue #20 feature commit; the CLI, report renderer,
  tests, and Docker boundary remain intact.
- **Status:** Merged via PR #21 as squash commit `789559f`; Pages deployment
  fix completed via PR #23 as squash commit `7075929`.

## Issue #24 — align docs and add direct input parser coverage

- **Goal:** Make project claims match the current range/diff behavior, remove
  stale lifecycle notes, and protect the input parser with package-level tests.
- **Approach:** Rewrite Architecture and Flow as current maps, distinguish
  range evidence/risk from diff summary-only behavior in README and the landing
  page, condense Handoff to current state, correct the fixture-regression
  wording and historical verification note, remove the stale deadline estimate,
  and add direct range/diff boundary tests under `internal/input/`.
- **Files touched:** `AGENTS.md`, `README.md`, `site/index.html`,
  `docs/Architecture.md`, `docs/Flow.md`, `docs/Feature.md`,
  `docs/Handoff.md`, `docs/Decisions.md`, `internal/input/contract_test.go`.
- **Verification:** Docker Go 1.23.0 passed the direct `internal/input` tests;
  documentation smoke checks and `git diff --check` passed. Full Docker and
  static-site checks remain part of the issue lifecycle.
- **Rollback:** Revert the issue #24 commit; runtime behavior and the report
  renderer remain unchanged.
- **Status:** Implementation ready for review; PR lifecycle pending.

## Issue #25 — preserve unclassified GitHub squash references

- **Goal:** Keep common GitHub-generated `Title (#N)` context without falsely
  classifying the number as an issue or pull request.
- **Approach:** Parse only a parenthesized numeric subject suffix as an
  `unclassified` reference, keep bare and cross-repository numeric forms
  ignored, preserve citation provenance/order/de-duplication, and let the
  existing reference signal raise risk without changing its weighting.
- **Files touched:** `internal/evidence/line_range.go`,
  `internal/evidence/line_range_test.go`, `internal/risk/risk_test.go`,
  `internal/report/`, `site/reports/*.html`, `site/index.html`, `README.md`,
  `docs/Architecture.md`, `docs/Flow.md`, `docs/Decisions.md`,
  `docs/Feature.md`, `docs/Handoff.md`.
- **Verification:** Docker Go 1.23.0 passed focused evidence/risk/report tests,
  the full suite, no-cache rebuild, Git/version/help checks, and the mounted
  self-demo. The self-demo shows `unclassified #18` and high risk; regenerated
  reports and static/HTTP site checks passed while diff behavior remains
  unchanged.
- **Rollback:** Revert the issue #25 commit; explicit issue/PR parsing and the
  existing risk matrix remain intact.
- **Status:** Implementation ready for review; PR lifecycle pending.
