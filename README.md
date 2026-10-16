# whyfence

whyfence is a local, git-only CLI that helps explain why code exists before it
is deleted or rewritten. Given a file range or diff, it will produce cited
evidence from local Git history, including related commits, issues, pull
requests, and tests changed alongside the code, plus a risk flag.

The core is designed to be deterministic and fully offline. Any LLM summary is
optional and disabled in tests.

## Project status

This repository currently contains the project bootstrap only. The CLI, runtime,
offline test command, Docker rebuild, static report, and GitHub Pages site have
not been implemented yet. No package-install command is published or documented
until a real implementation exists.

The planned deliverables are intentionally ordered:

1. CLI core.
2. Offline fixture-repository tests and Docker rebuild checks.
3. Self-contained static HTML reports.
4. A plain static landing page in `site/`.

## Development commands

There is no setup, run, test, or deploy command yet. These sections will be
updated with tested commands as each implementation step is completed. Do not
install an unpublished package or claim that the CLI is runnable at this stage.

## License

MIT. See [LICENSE](LICENSE).
