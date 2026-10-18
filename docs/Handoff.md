# Handoff

## Current status

- `main` was clean and synchronized at `a71ca47` when issue #24 started.
- Deliverables 1–4 are implemented: deterministic CLI core, offline fixture
  regression checks, standalone HTML reports, and the static Pages site.
- The published site is <https://sujal-suhaas.github.io/whyfence/>.
- Work in progress: issue #24 on `fix/24-docs-and-input-tests` corrects stale
  claims and adds direct `internal/input` contract tests.
- Work in progress: issue #25 on `feat/25-unclassified-squash-references`
  preserves narrow GitHub squash suffixes as honest `unclassified` references.

## Completed lifecycle

| Scope | Result |
| --- | --- |
| Issues #2–#9 | Core, fixtures, CLI contract, range evidence, references, changed tests, and risk merged and verified on `main`. |
| Issue #18 / PR #19 | Deterministic standalone HTML reports; merged as `bad17659`. |
| Issue #20 / PR #21 | Static landing page, five report snapshots, and Pages workflow; merged as `789559f`. |
| Issue #22 / PR #23 | Pages enablement fix; merged as `7075929`. |
| Pages deployment | Workflow run `39123292582` passed; site published successfully. |

## Next

1. Finish issue #24 with focused parser tests, documentation checks, and the
   issue PR lifecycle.
2. Finish the issue #25 PR lifecycle and review the regenerated self-report.
3. Prepare the final demo/submission handoff. Do not start stretch work without
   explicit approval.

## Truth boundaries

- Range input recovers cited local Git history, supported references, same-commit
  changed Go test paths, and an evidence-backed removal-risk signal.
- Diff input validates and summarizes source, files, and hunks only. It does not
  recover history or emit risk.
- The `recovered N of M` output is generated-fixture regression coverage with
  answers defined by the test fixtures. It is not a real-repository accuracy
  percentage or an external benchmark.
- Risk is a conservative evidence classification, not proof that removal is
  safe or unsafe and not a measure of ownership, runtime impact, dependencies,
  or business criticality.
- Reference parsing never queries GitHub. Explicit issue/PR forms are
  classified; parenthesized squash suffixes are labeled `unclassified`.

## Gotchas

- Docker image builds need network access for `apt-get install git`; test
  containers run with `--network=none`.
- Docker Desktop may need to be started before local verification.
- The Pages repository configuration needed one-time authorized enablement. The
  workflow now uses `enablement: true`; forks may still need Settings → Pages.
- `site/` is outside the Docker `COPY` list and is deployed separately by
  `.github/workflows/pages.yml`.
- No package is published. Use Go 1.23.0, Git, or Docker as documented in the
  README.

## Verification evidence

- Latest merged-main Docker rebuild before issue #24:
  `docker build --no-cache -t whyfence:test .` — passed.
- Latest merged-main offline core suite:
  `docker run --rm --network=none whyfence:test go test ./... -count=1` — passed.
- Latest merged-main Git check:
  `docker run --rm --network=none whyfence:test git --version` — passed with
  Git 2.39.5.
- Latest merged-main offline CLI smoke checks:
  `docker run --rm --network=none whyfence:test go run ./cmd/whyfence --version`
  and `--help` — passed; help listed `--report`.
- Static site smoke check — passed with five standalone reports and twenty local
  landing-page links.
- Local HTTP smoke check — passed for the landing page, report files, and demo.
- Pages workflow run `39123292582` — passed Configure Pages, upload, and deploy.
- `git diff --check` — passed on the merged landing-page and Pages-fix changes.

Update this handoff at the end of each issue lifecycle with only the current
branch, next action, gotchas, and actual checks.
