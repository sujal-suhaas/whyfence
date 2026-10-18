package fixtures

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFixtureBuilder(t *testing.T) {
	history, err := NewKnownHistory()
	if err != nil {
		t.Fatal(err)
	}
	cleaned := false
	t.Cleanup(func() {
		if cleaned {
			return
		}
		if err := history.Repository.Cleanup(); err != nil {
			t.Errorf("cleanup fixture repository: %v", err)
		}
	})

	workspace, err := filepath.Abs(".")
	if err != nil {
		t.Fatal(err)
	}
	tempRoot, err := filepath.Abs(os.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	fixtureDir, err := filepath.Abs(history.Repository.Dir)
	if err != nil {
		t.Fatal(err)
	}
	if relative, err := filepath.Rel(tempRoot, fixtureDir); err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		t.Fatalf("fixture %q is not under the system temp directory %q", fixtureDir, tempRoot)
	}
	if relative, err := filepath.Rel(workspace, fixtureDir); err == nil && relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		t.Fatalf("fixture %q is nested in workspace %q", fixtureDir, workspace)
	}

	if info, err := os.Stat(filepath.Join(fixtureDir, ".git")); err != nil || !info.IsDir() {
		t.Fatalf("fixture does not contain a Git directory: %v", err)
	}
	rootOutput, err := history.Repository.Git("rev-parse", "--show-toplevel")
	if err != nil {
		t.Fatal(err)
	}
	if got := filepath.Clean(strings.TrimSpace(rootOutput)); got != fixtureDir {
		t.Fatalf("Git root = %q, want %q", got, fixtureDir)
	}

	logOutput, err := history.Repository.Git("log", "--format=%H%x09%an%x09%ae%x09%aI%x09%cn%x09%ce%x09%cI%x09%s")
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(logOutput), "\n")
	if len(lines) != 3 {
		t.Fatalf("fixture commit count = %d, want 3; log = %q", len(lines), logOutput)
	}
	for _, line := range lines {
		fields := strings.Split(line, "\t")
		if len(fields) != 8 {
			t.Fatalf("unexpected commit metadata fields in %q", line)
		}
		if fields[1] != FixedAuthorName || fields[2] != FixedAuthorEmail || fields[4] != FixedAuthorName || fields[5] != FixedAuthorEmail {
			t.Fatalf("unexpected fixed identity in %q", line)
		}
	}
	if !strings.Contains(logOutput, "Tune timeout without test") || !strings.Contains(logOutput, "Fix timeout handling (PR #42)") || !strings.Contains(logOutput, "Add worker timeout guard for issue #41") {
		t.Fatalf("known history subjects missing from log: %q", logOutput)
	}
	if !strings.Contains(logOutput, history.Fix.When.Format("2006-01-02T15:04:05")) {
		t.Fatalf("fixed commit date missing from log: %q", logOutput)
	}

	statusOutput, err := history.Repository.Git("status", "--porcelain")
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(statusOutput) != "" {
		t.Fatalf("fixture has uncommitted changes: %q", statusOutput)
	}

	if err := history.Repository.Cleanup(); err != nil {
		t.Fatal(err)
	}
	cleaned = true
	if _, err := os.Stat(fixtureDir); !os.IsNotExist(err) {
		t.Fatalf("fixture directory still exists after cleanup: %v", err)
	}
}

func TestRiskFixtureBuilder(t *testing.T) {
	history, err := NewRiskHistory()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := history.Repository.Cleanup(); err != nil {
			t.Errorf("cleanup risk fixture repository: %v", err)
		}
	})

	logOutput, err := history.Repository.Git("log", "--format=%H%x09%s")
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(logOutput), "\n")
	if len(lines) != 4 {
		t.Fatalf("risk fixture commit count = %d, want 4; log = %q", len(lines), logOutput)
	}
	for _, subject := range []string{
		history.Low.Subject,
		history.TestOnly.Subject,
		history.Reference.Subject,
		history.High.Subject,
	} {
		if !strings.Contains(logOutput, subject) {
			t.Fatalf("risk fixture subject %q missing from log %q", subject, logOutput)
		}
	}
}
