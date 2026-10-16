package main

import (
	"bytes"
	"testing"
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
