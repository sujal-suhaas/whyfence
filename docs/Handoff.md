# Handoff

## Done

- Merged issue #2 through PR #10 as squash commit `73a2f614`.
- Pulled `main` and verified a clean no-cache Docker rebuild, Git in the image,
  offline container tests, and the offline version flag after the merge.
- Issue #3 is implemented and pushed on `feat/3-fixture-builder` as `1b5a766`.

## In progress

- Issue #3 has no PR yet; its branch is still based on the pre-merge issue #2
  history and must be rebased onto current `main` before PR review.
- Untracked `docs/Bug.md` and `docs/Decisions.md` remain from the documentation
  skeleton setup and are intentionally untouched.

## Next

- Rebase issue #3 onto current `main` and verify its exclusive diff.
- Show the issue #3 diff and wait for explicit approval before opening its PR.
- After PR checks pass and a second approval, squash-merge issue #3, pull
  `main`, rerun Docker tests, and update this handoff again.

## Gotchas

- Docker builds need network access for `apt-get install git`; test containers
  run with `--network=none`.
- Docker Desktop may need to be started before local Docker verification.
- The post-merge Docker checks run with `--network=none`; the image build needs
  network access for `apt-get install git`.
- The README install section remains intentionally unchanged.

## Verification

- `git pull --ff-only origin main`: passed after PR #10 merge.
- `docker build --no-cache -t whyfence:test .`: passed on merged `main`.
- `docker run --rm --network=none whyfence:test git --version`: passed.
- `docker run --rm --network=none whyfence:test go test ./...`: passed.
- `docker run --rm --network=none whyfence:test go run ./cmd/whyfence --version`:
  passed; output `whyfence 0.1.0-dev`.
