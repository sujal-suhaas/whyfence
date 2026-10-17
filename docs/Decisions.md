# Decisions

Record the reason for each durable choice and the alternatives rejected.

## 2026-10-17 — choose Go for the implementation

**Context:** Deliverables 1 and 2 must rebuild cleanly in Docker and run offline;
the core should have as few third-party dependencies as possible.

**Options considered:**

- **Go:** standard-library Git/process, parsing, and testing support; one
  distributable binary; no third-party dependency is needed for the planned
  deterministic core.
- **Python:** strong text-processing ergonomics and a familiar test ecosystem,
  but it needs a pinned interpreter plus environment/dependency handling and
  would add packaging/runtime choices even if kept dependency-free.

**Decision:** Use Go, initially with only the standard library. This minimizes
network-sensitive dependency resolution, makes the Docker runtime small and
predictable, and keeps the CLI portable as one binary.

**Rejected alternative:** Python was not chosen because its runtime and
packaging surface is larger for this offline, git/process-oriented CLI. It can
be reconsidered only with explicit evidence that Go blocks a required feature.

**Reproducibility note:** The implementation step must pin an exact Go toolchain
in the Dockerfile and `go.mod`. With no third-party modules, `go.mod` is the
module/version manifest and no `go.sum` is needed; if a dependency is added,
commit the generated `go.sum` as the lockfile.

## 2026-10-17 — build the verification foundation before evidence slices

**Context:** The Docker rebuild is a judging criterion, evidence recovery needs
repositories with known answers, and a late accuracy suite would hide failures
until the end of the schedule. The core evidence/risk work is large enough to
need independently reviewable changes.

**Decision:** Deliver the first two deliverables in this order:

1. Pin the Go toolchain, add a minimal Docker rebuild, and make one test pass.
2. Build temporary Git fixtures at test time with fixed authors and dates.
3. Start the `recovered N of M` harness with 2–3 cases.
4. Implement and test the CLI input/output contract.
5. Add evidence slices separately: line-range recovery, issue/PR reference
   extraction, and tests changed alongside the target code.
6. Add the deterministic risk flag after the evidence slices are stable.

**Why:** This keeps Docker green from the first implementation PR, makes every
evidence result testable against known history, and lets the accuracy suite grow
with each feature instead of becoming a rushed final check. Splitting the core
also preserves the one-change-per-request rule and makes rollback boundaries
clear.

**Rejected alternatives:** A late Docker gate, fixtures added after evidence
work, and one combined evidence/risk PR were rejected because they postpone
failures, make recovery unverifiable, or create an oversized review and rollback
surface.

## 2026-10-17 — use one GO for an issue lifecycle

**Context:** Requiring separate approvals to commit, push, open a PR, and merge
the same already-approved issue repeated the same decision and slowed a solo
buildathon workflow.

**Decision:** Treat an issue-level `GO` as authorization for that issue's full
lifecycle: implement, test, show the diff and evidence, commit, push, open the
PR, wait for checks, squash-merge, delete the branch, pull `main`, rerun Docker,
and update `Handoff.md`. Pause only for a blocker requiring manual user action
or a product decision.

**Why:** The safety gates remain—the diff, actual checks, rollback plan, CI
result, and post-merge Docker run are still required—but the user makes one
scope decision instead of repeating approval for mechanical Git transitions.

**Rejected alternative:** Separate approval at every lifecycle stage was
rejected as redundant once the issue scope and acceptance criteria were already
approved.

## 2026-10-17 — use explicit range and diff input modes

**Context:** The CLI needs a deterministic contract before evidence recovery is
implemented. Ambiguous positional input or automatic format detection would
make invalid-input behavior and tests harder to reproduce.

**Decision:** Accept exactly one explicit input mode: `--range path:start-end`
for an inclusive 1-based file range, or `--diff path` for a unified diff file;
`--diff -` reads the diff from stdin. Validate paths, line bounds, non-empty
input, and unified-diff headers/hunks before producing a stable summary.

**Why:** Explicit modes make the user contract visible in `--help`, give every
failure a deterministic `input:` error, and leave the production evidence
collector independent from input parsing.

**Rejected alternatives:** Positional arguments and auto-detecting whether a
single string is a path, range, or diff were rejected because they create
ambiguous parsing and platform-sensitive error behavior.

## 2026-10-17 — trace ranges with Git line history

**Context:** Issue #6 needs to recover local commits that explain a targeted
line range, with full-hash and subject citations that remain reproducible in
offline tests.

**Options considered:**

- **`git log -L`:** follows the selected lines through their local history and
  returns the commits in Git's stable newest-to-oldest order.
- **`git blame`:** identifies only the latest commit associated with the current
  lines and does not recover the range's evolution.
- **Patch parsing:** would duplicate Git's line-tracking behavior and create a
  larger, less trustworthy implementation surface.

**Decision:** Invoke `git log -L` with `--format=%H%x09%s`, parse the first tab
as the citation boundary, preserve Git's order, and de-duplicate hashes. Resolve
the repository from the requested file path so the CLI works with both absolute
and current-working-directory-relative paths. Keep the collector local-only and
return deterministic `evidence:` errors for unsupported paths or histories.

**Why:** Git already owns the difficult line-history traversal; using its full
object hash and subject gives each result a direct citation while avoiding
network, dependency, and heuristic variability.

**Rejected alternative:** `git blame` was rejected because a latest-owner-only
answer would miss the baseline commit in the known two-commit history and would
not satisfy the range-recovery goal.
