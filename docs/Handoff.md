# Handoff

## Current status

- `main` is clean and synchronized with `origin/main` after the issue #25
  handoff.
- Deliverables 1–4 are implemented: deterministic CLI core, offline fixture
  regression checks, standalone HTML reports, and the static Pages site.
- The published site is <https://sujal-suhaas.github.io/whyfence/>.
- Issues #24 and #25 are merged and verified on `main`.
- Issue #28 is resolved and verified on `main`.

## Completed lifecycle

| Scope | Result |
| --- | --- |
| Issues #2–#9 | Core, fixtures, CLI contract, range evidence, references, changed tests, and risk merged and verified on `main`. |
| Issue #18 / PR #19 | Deterministic standalone HTML reports; merged as `bad17659`. |
| Issue #20 / PR #21 | Static landing page, five report snapshots, and Pages workflow; merged as `789559f`. |
| Issue #22 / PR #23 | Pages enablement fix; merged as `7075929`. |
| Issue #24 / PR #26 | Documentation alignment and direct input parser tests; merged as `106585b`. |
| Issue #25 / PR #27 | Honest `unclassified` GitHub squash references; merged as `457d1af`. |
| Issue #28 / PR #29 | Subject-end boundary fix for unclassified references; merged as `8f86cc8`. |
| Pages deployment | Workflow run `39126827123` passed; site published successfully. |

## Next

1. Prepare the final demo/submission handoff. Do not start stretch work without
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

- Latest merged-main Docker rebuild after issue #28:
  `docker build --no-cache -t whyfence:test .` — passed.
- Latest merged-main offline core suite:
  `docker run --rm --network=none whyfence:test go test ./... -count=1` — passed.
- Latest merged-main Git check:
  `docker run --rm --network=none whyfence:test git --version` — passed with
  Git 2.39.5.
- Latest merged-main offline CLI smoke checks:
  `docker run --rm --network=none whyfence:test go run ./cmd/whyfence --version`
  and `--help` — passed; help listed `--report`.
- Mounted self-demo after issue #28:
  `docker run --rm --network=none -v "$PWD:/src" -w /src whyfence:test go run
  ./cmd/whyfence --range cmd/whyfence/main.go:26-32` — passed; output included
  `unclassified #18` and `risk: high`.
- Focused evidence/risk/report tests after issue #25 — passed, including
  unclassified-reference parsing and risk behavior.
- Focused suffix-boundary test after issue #28 — passed; `Title (#46) trailing`
  remains unclassified-reference free.
- Static site smoke check — passed with five standalone reports and twenty local
  landing-page links.
- Local HTTP smoke check — passed for the landing page, report files, and demo.
- Pages workflow run `39126827123` — passed Configure Pages, upload, and deploy.
- Docker workflow run `39126827120` — passed image build, Git check, offline
  tests, and version check.
- PR #29 Docker CI and post-merge Docker checks — passed.
- `git diff --check` — passed after issue #28.

Update this handoff at the end of each issue lifecycle with only the current
branch, next action, gotchas, and actual checks.
