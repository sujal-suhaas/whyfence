package evidence

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/sujal-suhaas/whyfence/internal/fixtures"
)

func TestLineRange(t *testing.T) {
	history, err := fixtures.NewKnownHistory()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := history.Repository.Cleanup(); err != nil {
			t.Errorf("cleanup line-range fixture: %v", err)
		}
	})

	path := filepath.Join(history.Repository.Dir, filepath.FromSlash(history.TargetPath))
	want := []Citation{
		{Hash: history.Fix.Hash, Subject: history.Fix.Subject},
		{Hash: history.Baseline.Hash, Subject: history.Baseline.Subject},
	}
	got, err := RecoverLineRange(path, 3, 3)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("citations = %#v, want %#v", got, want)
	}

	repeated, err := RecoverLineRange(path, 3, 3)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(repeated, got) {
		t.Fatalf("repeated citations = %#v, want %#v", repeated, got)
	}
}

func TestLineRangeRejectsNonRepositoryFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "worker.go")
	if err := writeTestFile(path); err != nil {
		t.Fatal(err)
	}
	if _, err := RecoverLineRange(path, 1, 1); err == nil || err.Error() != "evidence: file is not inside a Git repository" {
		t.Fatalf("error = %v, want deterministic non-repository error", err)
	}
}

func TestLineRangeRejectsInvalidBounds(t *testing.T) {
	if _, err := RecoverLineRange("worker.go", 0, 1); err == nil || err.Error() != "evidence: line range must use positive start-end lines" {
		t.Fatalf("error = %v, want deterministic bounds error", err)
	}
}

func writeTestFile(path string) error {
	return os.WriteFile(path, []byte("package worker\n"), 0o644)
}
