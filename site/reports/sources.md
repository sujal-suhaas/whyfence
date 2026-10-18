# Pre-generated report sources

These reports are checked-in snapshots. Each was generated with the whyfence
CLI from a full local clone at the exact commit listed below. The HTML files do
not fetch the source repository when opened.

## whyfence

- Repository: <https://github.com/sujal-suhaas/whyfence>
- Commit: `bad176596cf3639db37bfb5d12fcd5e53d054cb6`
- Input: `cmd/whyfence/main.go:26-32`
- Report: [whyfence.html](whyfence.html)
- License: MIT
- Generation: from the repository root, `whyfence --range cmd/whyfence/main.go:26-32 --report site/reports/whyfence.html`

## spf13/cobra

- Repository: <https://github.com/spf13/cobra>
- Commit: `adbc8813901bba65827259daa8e22ff94ec1f30e`
- Input: `cobra.go:15-18`
- Report: [cobra.html](cobra.html)
- License: Apache-2.0
- Generation: from the Cobra repository root, `whyfence --range cobra.go:15-18 --report report.html`; the resulting file was copied to this report gallery.

## urfave/cli

- Repository: <https://github.com/urfave/cli>
- Commit: `192ed0b94135f015fb95e23b06e09688fce9737f`
- Input: `cli.go:1-18`
- Report: [urfave-cli.html](urfave-cli.html)
- License: MIT
- Generation: from the urfave/cli repository root, `whyfence --range cli.go:1-18 --report report.html`; the resulting file was copied to this report gallery.

## BurntSushi/toml

- Repository: <https://github.com/BurntSushi/toml>
- Commit: `d733fc535e4a9f3c454e421a87b90ade5d49bf31`
- Input: `decode.go:3-18`
- Report: [burntsushi-toml.html](burntsushi-toml.html)
- License: MIT
- Generation: from the BurntSushi/toml repository root, `whyfence --range decode.go:3-18 --report report.html`; the resulting file was copied to this report gallery.

## sirupsen/logrus

- Repository: <https://github.com/sirupsen/logrus>
- Commit: `6ef0748d1177dadaeee9aa79ffef604570f487d7`
- Input: `logger.go:12-18`
- Report: [logrus.html](logrus.html)
- License: MIT
- Generation: from the logrus repository root, `whyfence --range logger.go:12-18 --report report.html`; the resulting file was copied to this report gallery.

## Attribution and scope

The repository names, source files, commit subjects, and license names above
remain the property of their respective projects and contributors. whyfence is
not affiliated with these projects. The reports are examples of local Git
history recovery, not endorsements or live project status. See each linked
repository for its complete license and attribution terms.
