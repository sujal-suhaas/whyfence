package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sujal-suhaas/whyfence/internal/fixtures"
)

func TestVersionFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if exitCode := run([]string{"--version"}, &stdout, &stderr); exitCode != 0 {
		t.Fatalf("--version exit code = %d, want 0", exitCode)
	}
	if got, want := stdout.String(), "whyfence "+version+"\n"; got != want {
		t.Fatalf("--version output = %q, want %q", got, want)
	}
	if stderr.Len() != 0 {
		t.Fatalf("--version wrote unexpected stderr: %q", stderr.String())
	}
}

func TestHelpFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if exitCode := run([]string{"--help"}, &stdout, &stderr); exitCode != 0 {
		t.Fatalf("--help exit code = %d, want 0", exitCode)
	}
	if !strings.Contains(stderr.String(), "-range string") || !strings.Contains(stderr.String(), "-diff string") {
		t.Fatalf("help output omitted input flags: %q", stderr.String())
	}
}

func TestCLIContract(t *testing.T) {
	history, err := fixtures.NewKnownHistory()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := history.Repository.Cleanup(); err != nil {
			t.Errorf("cleanup CLI fixture: %v", err)
		}
	})
	file := filepath.Join(history.Repository.Dir, filepath.FromSlash(history.TargetPath))

	t.Run("accepts file range", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		if exitCode := run([]string{"--range", file + ":1-3"}, &stdout, &stderr); exitCode != 0 {
			t.Fatalf("range exit code = %d, stderr = %q", exitCode, stderr.String())
		}
		want := fmt.Sprintf("input: range\nfile: %s\nlines: 1-3\ncommit: %s %s\ncommit: %s %s\ntest: %s\nreference: pr #42\ncommit: %s %s\ntest: %s\nreference: issue #41\nrisk: high\nstatus: accepted\n", file, history.CodeOnly.Hash, history.CodeOnly.Subject, history.Fix.Hash, history.Fix.Subject, history.TestPath, history.Baseline.Hash, history.Baseline.Subject, history.TestPath)
		if stdout.String() != want {
			t.Fatalf("range output = %q, want %q", stdout.String(), want)
		}
	})

	t.Run("accepts diff file", func(t *testing.T) {
		diffFile := filepath.Join(t.TempDir(), "change.patch")
		diff := "diff --git a/worker.go b/worker.go\n--- a/worker.go\n+++ b/worker.go\n@@ -1,1 +1,1 @@\n-old\n+new\n"
		if err := os.WriteFile(diffFile, []byte(diff), 0o644); err != nil {
			t.Fatal(err)
		}
		var stdout, stderr bytes.Buffer
		if exitCode := run([]string{"--diff", diffFile}, &stdout, &stderr); exitCode != 0 {
			t.Fatalf("diff exit code = %d, stderr = %q", exitCode, stderr.String())
		}
		want := "input: diff\nsource: " + diffFile + "\nfiles: 1\nhunks: 1\nstatus: accepted\n"
		if stdout.String() != want {
			t.Fatalf("diff output = %q, want %q", stdout.String(), want)
		}
	})

	t.Run("accepts diff stdin", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		stdin := bytes.NewBufferString("--- a/worker.go\n+++ b/worker.go\n@@ -1 +1 @@\n-old\n+new\n")
		if exitCode := runWithInput([]string{"--diff", "-"}, stdin, &stdout, &stderr); exitCode != 0 {
			t.Fatalf("stdin diff exit code = %d, stderr = %q", exitCode, stderr.String())
		}
		if stdout.String() != "input: diff\nsource: -\nfiles: 1\nhunks: 1\nstatus: accepted\n" {
			t.Fatalf("stdin diff output = %q", stdout.String())
		}
	})
}

func TestCLIContractErrors(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.go")
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "no input", want: "input: exactly one of --range or --diff is required\n"},
		{name: "both inputs", args: []string{"--range", "worker.go:1-1", "--diff", "-"}, want: "input: exactly one of --range or --diff is required\n"},
		{name: "malformed range", args: []string{"--range", "worker.go:bad"}, want: "input: range must use path:start-end\n"},
		{name: "missing range file", args: []string{"--range", missing + ":1-1"}, want: "input: file \"" + missing + "\" not found\n"},
		{name: "empty diff", args: []string{"--diff", "-"}, want: "input: diff is empty\n"},
		{name: "malformed diff", args: []string{"--diff", "-"}, want: "input: malformed unified diff\n"},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			stdin := bytes.NewBufferString("")
			if testCase.name == "malformed diff" {
				stdin = bytes.NewBufferString("not a diff\n")
			}
			if exitCode := runWithInput(testCase.args, stdin, &stdout, &stderr); exitCode != 2 {
				t.Fatalf("exit code = %d, want 2; stderr = %q", exitCode, stderr.String())
			}
			if stderr.String() != testCase.want {
				t.Fatalf("stderr = %q, want %q", stderr.String(), testCase.want)
			}
			if stdout.Len() != 0 {
				t.Fatalf("unexpected stdout = %q", stdout.String())
			}
		})
	}
}
