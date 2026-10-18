# Bug log

Append-only. Add one entry per bug with the problem, cause, approach, files,
verification, rollback, and status.

## Issue #22 — enable GitHub Pages before static-site deployment

- **Problem:** The first post-merge run of `.github/workflows/pages.yml` failed
  in `actions/configure-pages@v5` with `Get Pages site failed ... Not Found`.
- **Cause:** The repository had no Pages site configured and the workflow used
  the action's default `enablement: false`.
- **Approach:** Set `enablement: true` so the supported Pages action can create
  the workflow-backed Pages configuration before uploading the existing static
  `site/` directory.
- **Files touched:** `.github/workflows/pages.yml`, `docs/Bug.md`.
- **Verification:** The pre-fix failure is recorded at
  <https://github.com/sujal-suhaas/whyfence/actions/runs/39122917936>. PR #23
  passed Docker CI, merged as `7075929`, and the post-fix Pages run
  <https://github.com/sujal-suhaas/whyfence/actions/runs/39123292582> passed
  Configure Pages, upload, and deploy. The published site is
  <https://sujal-suhaas.github.io/whyfence/>. Post-merge Docker rebuild and
  full tests also passed.
- **Rollback:** Revert the issue #22 fix commit; the site can still be enabled
  manually in repository Settings → Pages.
- **Status:** Resolved via PR #23 as squash commit `7075929`.

## Issue #28 — restrict unclassified references to subject suffixes

- **Problem:** The new `unclassified` parser boundary accepted `(#N)` followed
  by additional subject text, although the contract is specifically the
  GitHub-style end-of-subject suffix.
- **Cause:** The initial regular expression allowed whitespace after the
  parenthesized number instead of requiring end-of-subject.
- **Approach:** Require `(#N)` to end the subject and add a regression test for
  `Title (#46) trailing`.
- **Files touched:** `internal/evidence/line_range.go`,
  `internal/evidence/line_range_test.go`, `docs/Bug.md`.
- **Verification:** Focused `TestExtractReferences` passed; full offline/Docker
  checks and the issue PR lifecycle remain required.
- **Rollback:** Revert the issue #28 fix commit; issue #25's broader boundary
  behavior is restored.
- **Status:** Fix in progress.
