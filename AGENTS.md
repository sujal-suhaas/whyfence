# whyfence project instructions

whyfence is a solo, local, git-only CLI for the Developer Tools & Open Source
track. The core is deterministic, cited, and fully offline; LLM summaries are
optional and disabled in tests. Deadline: 20 Oct, 11:59 PM IST.

## Deliverables — do not start a later item until the earlier one works and is tested

1. CLI core: file range in for cited local-Git evidence and risk; diff in for a
   validated summary because diff evidence/risk is not implemented yet.
2. Offline tests, including fixture git repos with known reasons, reporting
   `recovered N of M`, plus Docker rebuild checks.
3. Self-contained static HTML report output from the CLI; no server.
4. `site/` GitHub Pages landing page: plain HTML/CSS/JS, no build step, real
   sample report, short terminal demo, tested install section, and 3–5 public-
   repository pre-generated reports.

Stretch, including any “paste a GitHub link” feature, requires explicit approval.

## Hard constraints

- Judges rebuild on clean Linux and run offline. Keep git, a pinned language
  version, a lockfile, Dockerfile, one offline test command, README with exact
  setup/run/test/deploy steps, and `.env.example`.
- Keep all required files in the project; never use machine-specific paths or
  commit secrets. Keep `site/` outside tests and the Docker rebuild.
- Every documented install command must work today from a clean environment;
  never document an unpublished package. README and landing-page claims must
  match the code. All development happens in OmniRush sessions.

## Working rules

1. Ask “why” before “what”: state reasoning and assumptions before coding and
   challenge flawed logic. One change per request; if bundled, propose splits.
2. At the start of each session, read `docs/Architecture.md` fully. Before
   each change, read the relevant sections of `docs/Architecture.md` and
   `docs/Flow.md`. Keep Architecture as the system map and Flow as the
   file/function execution map; update the relevant document when structure or
   execution flow changes.
3. In `docs/Decisions.md`, record why a choice was made and alternatives
   rejected, not merely what changed. Use sub-agents for independent tasks when
   helpful, and record that decision there.
4. Comments explain intent and reasoning, never merely what code does.
5. Keep append-only `docs/Bug.md` and `docs/Feature.md`, one entry per item;
   include problem/goal, cause/approach, files touched, verification, status.
6. For every change, record a test checklist with the exact command and actual
   result, plus a rollback plan: which commit to revert and other impact.
7. For each issue approved with `GO`, show/report every diff and actual checks;
   commit/push, open PR, wait for checks, squash-merge/delete the branch, pull
   `main`, rerun Docker, and update `Handoff.md`. Pause only if blocked.
8. At session end, write `docs/Handoff.md`: done, in progress, next, gotchas.
9. If behind at the 19 Oct checkpoint, cut in this order: landing page, then
   HTML report; never cut the core or its tests.

## Git workflow
- Create a public GitHub repo named `whyfence` on the connected account, with
  README, license, and `.gitignore`; use `main` as the default branch.
- `main` must always pass the offline test command inside Docker.
- Open one GitHub issue per feature or bug, mirroring `Feature.md` or `Bug.md`,
  with a one-line goal and done criteria.
- Use one branch per issue: `feat/<n>-short-name` or `fix/<n>-short-name`.
- Keep commits small and descriptive, and reference the issue number.
- Open one PR per issue, not per commit. Its body must include why, what changed,
  a test checklist with real evidence, a rollback plan, and `Closes #n`.
- A `GO` for an issue authorizes its full lifecycle; preserve the diff/evidence
  review, then open and squash-merge its PR after required checks pass.
- After each merge, pull `main`, rerun the Docker tests, and update `Handoff.md`.
- Never force-push `main` or rewrite history.
