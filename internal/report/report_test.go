package report

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHTMLReport(t *testing.T) {
	risk := "high"
	data := Data{
		Status: "accepted",
		Input: Input{
			Mode:  "range",
			Path:  `src/<worker>.go`,
			Start: 3,
			End:   4,
		},
		Citations: []Citation{{
			Hash:    "abc123",
			Subject: `Fix <script>alert("x")</script>`,
			TestPaths: []string{
				`internal/worker_test.go`,
			},
			References: []Reference{{Kind: "issue", Number: 41, Raw: "issue #41"}},
		}},
		Risk: &risk,
	}

	first, err := Render(data)
	if err != nil {
		t.Fatal(err)
	}
	second, err := Render(data)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first, second) {
		t.Fatal("identical report data produced different HTML")
	}
	for _, want := range []string{
		"<!doctype html>",
		"<style>",
		"&lt;script&gt;alert(&#34;x&#34;)&lt;/script&gt;",
		"src/&lt;worker&gt;.go",
		"risk-high",
		"internal/worker_test.go",
		"issue #41",
	} {
		if !strings.Contains(string(first), want) {
			t.Fatalf("report omitted %q", want)
		}
	}
	for _, forbidden := range []string{"<script", "<link", "http://", "https://"} {
		if strings.Contains(string(first), forbidden) {
			t.Fatalf("report contains forbidden external/runtime content %q", forbidden)
		}
	}
}

func TestDiffReportOmitsRisk(t *testing.T) {
	data := Data{
		Status: "accepted",
		Input:  Input{Mode: "diff", Path: "change.patch", Files: 2, Hunks: 3},
	}
	output, err := Render(data)
	if err != nil {
		t.Fatal(err)
	}
	text := string(output)
	for _, want := range []string{"change.patch", "2", "3", "Diff input is accepted", "risk are unavailable"} {
		if !strings.Contains(text, want) {
			t.Fatalf("diff report omitted %q", want)
		}
	}
	if strings.Contains(text, "<dt>Risk</dt>") {
		t.Fatal("diff report unexpectedly contains a risk value")
	}
}

func TestRenderOmitsTrailingWhitespace(t *testing.T) {
	output, err := Render(Data{Status: "accepted", Input: Input{Mode: "diff", Path: "change.patch"}})
	if err != nil {
		t.Fatal(err)
	}
	for lineNumber, line := range strings.Split(string(output), "\n") {
		if strings.TrimRight(line, " \t") != line {
			t.Fatalf("line %d has trailing whitespace: %q", lineNumber+1, line)
		}
	}
}

func TestWriteAtomicOverwritesReport(t *testing.T) {
	path := filepath.Join(t.TempDir(), "report.html")
	if err := WriteAtomic(path, []byte("first")); err != nil {
		t.Fatal(err)
	}
	if err := WriteAtomic(path, []byte("second")); err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != "second" {
		t.Fatalf("report contents = %q, want %q", contents, "second")
	}
}
