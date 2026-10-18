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

## 2026-10-17 — classify only explicit local issue and PR references

**Context:** Issue #7 needs tracker references from local commit messages while
the deterministic core remains offline and cannot ask GitHub whether a number is
an issue or a pull request.

**Supported forms:**

- Issues: `issue #41`, `issue: #41`, and GitHub `/issues/41` URLs.
- Pull requests: `PR #42`, `PR: #42`, `pull request #42`, generated `Merge pull
  request #42 from ...` subjects, and GitHub `/pull/42` URLs.

**Decision:** Parse commit subjects only because the line-range collector already
retrieves `%s`. Return references in first-appearance order, de-duplicate the
same kind/number within one citation, and preserve each reference on its source
commit. Ignore bare `#42` and `owner/repo#42` because they are ambiguous without
remote metadata.

**Why:** Explicit wording or a URL path supplies enough local information to
classify the reference without network access. Keeping provenance on `Citation`
prevents a repeated number from losing the commit that mentioned it.

**Rejected alternatives:** Treating every `#number` as an issue or PR was
rejected because it would make a false classification. Querying GitHub was
rejected because offline, local-git-only execution is a hard constraint.

## 2026-10-17 — identify changed tests by same-commit `_test.go` paths

**Context:** Issue #8 needs related test evidence for each commit recovered from
a target line range. The relationship must be local, deterministic, and
auditable without inferring imports or querying a code host.

**Decision:** For every recovered commit, run `git diff-tree --root --name-only
--no-renames -r -z`, keep repository-relative paths whose basename ends in
`_test.go`, sort them lexically, and attach them to that commit's `Citation` as
`TestPaths`. A root commit is treated the same as any other commit. An empty
list means that commit changed no recognized test path.

**Why:** The same-commit rule is directly observable in local Git history, the
Go `_test.go` suffix is unambiguous for this repository, `--root` covers the
fixture's initial commit, and sorted slash-separated paths make output stable.
Disabling rename detection avoids reporting heuristic rename pairs as new test
evidence in this first slice.

**Rejected alternatives:** Directory names such as `tests/`, test data files,
and language-neutral conventions were rejected because they would broaden an
underspecified policy beyond the current Go project. A naïve parent-to-commit
diff was rejected because it misses root-commit changes.

## 2026-10-17 — classify evidence-backed removal risk

**Context:** The evidence slices now recover commit citations, explicit issue/PR
references, and same-commit changed tests. Issue #9 needs a deterministic risk
flag without claiming knowledge that local Git cannot provide.

**Decision:** Aggregate the complete recovered citation set into two signals:
whether any citation has a `TestPaths` entry and whether any citation has an
explicit `Reference`. Apply this matrix:

| Evidence signals | Level |
| --- | --- |
| neither tests nor references | `low` |
| tests only, or references only | `medium` |
| both tests and references | `high` |
| empty citation set | `high` |

The evaluator is order-independent, does not mutate citations, and describes
evidence-backed removal risk only. It does not measure runtime impact, ownership,
dependency reachability, or business criticality. The initial flag is emitted
for range input only; diff input has no recovered evidence yet.

**Why:** A changed test signals an exercised behavior, while an explicit issue or
PR reference signals a documented change context. The combination warrants the
strongest caution. Treating empty evidence as high is conservative because no
local history was recovered to support safe removal.

**Rejected alternatives:** A numeric score and weighted thresholds were rejected
because they would imply precision unsupported by the current evidence. Low risk
for empty evidence was rejected because absence of recovered history is not
evidence of safe removal. Emitting a diff risk was rejected until diff history
recovery exists.

## 2026-10-18 — write deterministic standalone HTML reports

**Context:** Deliverable 3 needs a browser-openable artifact without adding a
server, network access, or a second output dependency. The existing range and
diff analysis results already contain the data needed for a report.

**Decision:** Add `--report PATH` as an optional named output path. Render a
complete HTML document with `html/template`, inline CSS, semantic sections, and
no JavaScript, timestamp, external URL, font, image, or stylesheet. Range reports
include citations, changed tests, references, and risk. Diff reports include
their validated source/file/hunk summary and explicitly state that evidence and
risk are unavailable. Without `--report`, terminal output remains unchanged.
Write through a same-directory temporary file and rename it into place; reject
`--report -` to keep stdout and HTML unambiguous. Report failures use an
`output:` prefix and occur before any success confirmation.

**Why:** A named, inline-only artifact is portable to a local browser and GitHub
Pages, deterministic for tests, and safe to use offline. Keeping rendering
separate from evidence collection prevents HTML formatting from performing Git
queries or changing analysis behavior.

**Rejected alternatives:** JSON-first output was rejected because the immediate
deliverable is a browser-readable report. A server, external assets, and
JavaScript were rejected because judges rebuild and run offline. Timestamps were
rejected because they make identical analyses produce different bytes. Diff risk
was rejected because diff history evidence is not yet available.

## 2026-10-18 — publish a plain static landing page from `site/`

**Context:** Deliverable 4 needs a judge-readable project introduction, a real
report sample, a terminal demo, and public-repository examples. The project must
remain build-free and the CLI/Docker test boundary must stay focused on the
offline Go core.

**Decision:** Keep the landing page as plain local HTML and CSS in `site/`,
check in five standalone report snapshots, and publish that directory with a
GitHub Pages workflow that performs no site build. Record each snapshot's exact
repository commit, input range, generation command, and license in
`site/reports/sources.md`.

**Why:** Relative local links and self-contained reports work both on GitHub
Pages and from a downloaded checkout. A checked-in snapshot gallery proves the
actual CLI output without requiring the page to fetch GitHub data or run a
server. Keeping `site/` outside Docker's explicit `COPY` list preserves the
reproducible test boundary.

**Rejected alternatives:** A JavaScript-rendered gallery, CDN styles, and live
GitHub API calls were rejected because they add runtime network dependencies and
make offline review weaker. A package-manager install command was rejected
because whyfence is not published. A live report generator on the page was
rejected because reports are intentionally local, pinned snapshots.

## 2026-10-18 — describe fixture counts as regression coverage

**Context:** The offline harness reports `recovered N of M`, but every expected
reason is authored by the generated fixtures. Presenting that count as a
real-repository accuracy percentage would overstate what the test proves.

**Decision:** Describe the harness as generated-fixture regression coverage in
the README, architecture, flow, feature log, handoff, and demo guidance. Keep
the stable `recovered N of M` output because it makes missed regressions visible
without implying an external benchmark.

**Rejected alternative:** Calling the count an accuracy score was rejected
because the current suite does not sample independent repositories or establish
a ground-truth corpus outside its own fixtures.
