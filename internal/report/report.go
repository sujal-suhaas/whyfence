// Package report renders deterministic, self-contained HTML analysis reports.
package report

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
)

// Data is the complete input to the static report renderer.
type Data struct {
	Status    string
	Input     Input
	Citations []Citation
	Risk      *string
}

// Input describes the validated CLI input represented by a report.
type Input struct {
	Mode  string
	Path  string
	Start int
	End   int
	Files int
	Hunks int
}

// Citation preserves the evidence provenance shown in a report.
type Citation struct {
	Hash       string
	Subject    string
	TestPaths  []string
	References []Reference
}

// Reference identifies an issue or pull request attached to a citation.
type Reference struct {
	Kind   string
	Number int
	Raw    string
}

var documentTemplate = template.Must(template.New("report").Parse(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>whyfence evidence report</title>
  <style>
    :root { color-scheme: light; font-family: system-ui, sans-serif; background: #f5f7fb; color: #172033; }
    body { margin: 0; padding: 2rem; }
    main { max-width: 72rem; margin: 0 auto; }
    header, section, footer { background: #fff; border: 1px solid #dbe2ee; border-radius: 0.75rem; padding: 1.25rem; margin-bottom: 1rem; box-shadow: 0 0.25rem 1rem rgba(23, 32, 51, 0.05); }
    header { display: flex; align-items: baseline; justify-content: space-between; gap: 1rem; }
    h1, h2, p { margin-top: 0; }
    h1 { margin-bottom: 0; font-size: 1.6rem; }
    h2 { font-size: 1.1rem; }
    .eyebrow { color: #52627a; font-size: 0.78rem; font-weight: 700; letter-spacing: 0.12em; text-transform: uppercase; }
    .status, .risk { border-radius: 999px; display: inline-block; font-size: 0.78rem; font-weight: 700; padding: 0.35rem 0.7rem; text-transform: uppercase; }
    .status { background: #e8f5ed; color: #17633a; }
    .risk-low { background: #e8f5ed; color: #17633a; }
    .risk-medium { background: #fff3d6; color: #815b00; }
    .risk-high { background: #fde8e8; color: #a32121; }
    dl { display: grid; gap: 0.7rem; grid-template-columns: 10rem 1fr; margin: 0; }
    dt { color: #52627a; font-weight: 700; }
    dd { margin: 0; overflow-wrap: anywhere; }
    code { background: #eef2f8; border-radius: 0.3rem; padding: 0.12rem 0.3rem; }
    ol { display: grid; gap: 1rem; padding-left: 1.4rem; }
    li { padding-left: 0.3rem; }
    .commit { font-weight: 700; overflow-wrap: anywhere; }
    .subject { font-weight: 400; }
    .evidence { color: #52627a; font-size: 0.92rem; margin: 0.5rem 0 0; }
    .evidence ul { margin: 0.35rem 0 0; }
    .notice { color: #52627a; margin-bottom: 0; }
    footer { color: #52627a; font-size: 0.85rem; }
  </style>
</head>
<body>
  <main>
    <header>
      <div>
        <p class="eyebrow">whyfence</p>
        <h1>Evidence report</h1>
      </div>
      <span class="status">{{.Status}}</span>
    </header>

    <section aria-labelledby="input-heading">
      <h2 id="input-heading">Input</h2>
      <dl>
        <dt>Mode</dt><dd><code>{{.Input.Mode}}</code></dd>
        {{if eq .Input.Mode "range"}}
        <dt>File</dt><dd><code>{{.Input.Path}}</code></dd>
        <dt>Lines</dt><dd><code>{{.Input.Start}}-{{.Input.End}}</code></dd>
        {{else}}
        <dt>Source</dt><dd><code>{{.Input.Path}}</code></dd>
        <dt>Files</dt><dd><code>{{.Input.Files}}</code></dd>
        <dt>Hunks</dt><dd><code>{{.Input.Hunks}}</code></dd>
        {{end}}
        {{if .Risk}}
        <dt>Risk</dt><dd><span class="risk risk-{{.Risk}}">{{.Risk}}</span></dd>
        {{end}}
      </dl>
    </section>

    {{if .Citations}}
    <section aria-labelledby="evidence-heading">
      <h2 id="evidence-heading">Recovered evidence</h2>
      <ol>
        {{range .Citations}}
        <li>
          <div class="commit"><code>{{.Hash}}</code> <span class="subject">{{.Subject}}</span></div>
          {{if .TestPaths}}<div class="evidence"><strong>Tests changed</strong><ul>{{range .TestPaths}}<li><code>{{.}}</code></li>{{end}}</ul></div>{{end}}
          {{if .References}}<div class="evidence"><strong>References</strong><ul>{{range .References}}<li><code>{{.Kind}} #{{.Number}}</code></li>{{end}}</ul></div>{{end}}
        </li>
        {{end}}
      </ol>
    </section>
    {{else if eq .Input.Mode "diff"}}
    <section aria-labelledby="availability-heading">
      <h2 id="availability-heading">Evidence availability</h2>
      <p class="notice">Diff input is accepted, but local commit evidence and risk are unavailable for this input mode.</p>
    </section>
    {{end}}

    <footer>Status: {{.Status}}</footer>
  </main>
</body>
</html>
`))

// Render returns deterministic standalone HTML for data.
func Render(data Data) ([]byte, error) {
	var output byteBuffer
	if err := documentTemplate.Execute(&output, data); err != nil {
		return nil, fmt.Errorf("render report: %w", err)
	}
	return trimTrailingWhitespace(output.Bytes()), nil
}

// trimTrailingWhitespace keeps generated snapshots readable and makes their
// byte-level diffs reflect report content instead of template indentation.
func trimTrailingWhitespace(contents []byte) []byte {
	lines := bytes.Split(contents, []byte{'\n'})
	for index := range lines {
		lines[index] = bytes.TrimRight(lines[index], " \t")
	}
	return bytes.Join(lines, []byte{'\n'})
}

// WriteAtomic writes contents to path through a same-directory temporary file.
func WriteAtomic(path string, contents []byte) error {
	if path == "" || path == "-" {
		return errors.New("report path must be a named file")
	}
	directory := filepath.Dir(path)
	temporary, err := os.CreateTemp(directory, ".whyfence-report-*")
	if err != nil {
		return err
	}
	temporaryName := temporary.Name()
	defer os.Remove(temporaryName)
	if _, err := temporary.Write(contents); err != nil {
		_ = temporary.Close()
		return err
	}
	if err := temporary.Close(); err != nil {
		return err
	}
	if err := os.Rename(temporaryName, path); err == nil {
		return nil
	}
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return os.Rename(temporaryName, path)
}

type byteBuffer struct {
	data []byte
}

func (b *byteBuffer) Write(value []byte) (int, error) {
	b.data = append(b.data, value...)
	return len(value), nil
}

func (b *byteBuffer) Bytes() []byte {
	return b.data
}

var _ io.Writer = (*byteBuffer)(nil)
