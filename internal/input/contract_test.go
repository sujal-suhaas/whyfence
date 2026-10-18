package input

import (
	"io/fs"
	"reflect"
	"strings"
	"testing"
)

func TestParseRange(t *testing.T) {
	readFile := func(path string) ([]byte, error) {
		if path == "missing.go" {
			return nil, fs.ErrNotExist
		}
		return []byte("one\ntwo\nthree\n"), nil
	}

	tests := []struct {
		name string
		spec string
		want Request
		err  string
	}{
		{name: "inclusive range", spec: "worker.go:2-3", want: Request{Mode: RangeMode, Path: "worker.go", Start: 2, End: 3}},
		{name: "bad separator", spec: "worker.go:2", err: "input: range must use path:start-end"},
		{name: "bad bounds", spec: "worker.go:0-2", err: "input: range must use positive start-end lines"},
		{name: "missing file", spec: "missing.go:1-1", err: `input: file "missing.go" not found`},
		{name: "past end", spec: "worker.go:1-4", err: `input: range 1-4 exceeds file "worker.go" with 3 lines`},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := ParseRange(testCase.spec, readFile)
			if testCase.err != "" {
				if err == nil || err.Error() != testCase.err {
					t.Fatalf("error = %v, want %q", err, testCase.err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, testCase.want) {
				t.Fatalf("request = %#v, want %#v", got, testCase.want)
			}
		})
	}
}

func TestParseDiffSource(t *testing.T) {
	readFile := func(path string) ([]byte, error) {
		return []byte("diff --git a/one.go b/one.go\n--- a/one.go\n+++ b/one.go\n@@ -1 +1 @@\n-old\n+new\n"), nil
	}

	tests := []struct {
		name   string
		source string
		stdin  string
		want   Request
		err    string
	}{
		{name: "git headers", source: "change.patch", want: Request{Mode: DiffMode, Path: "change.patch", Files: 1, Hunks: 1}},
		{name: "stdin fallback headers", source: "-", stdin: "--- a/two.go\n+++ b/two.go\n@@ -1 +1 @@\n-old\n+new\n", want: Request{Mode: DiffMode, Path: "-", Files: 1, Hunks: 1}},
		{name: "empty", source: "-", err: "input: diff is empty"},
		{name: "malformed", source: "-", stdin: "not a diff\n", err: "input: malformed unified diff"},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := ParseDiffSource(testCase.source, readFile, strings.NewReader(testCase.stdin))
			if testCase.err != "" {
				if err == nil || err.Error() != testCase.err {
					t.Fatalf("error = %v, want %q", err, testCase.err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, testCase.want) {
				t.Fatalf("request = %#v, want %#v", got, testCase.want)
			}
		})
	}
}
