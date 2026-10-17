package accuracy

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sujal-suhaas/whyfence/internal/evidence"
	"github.com/sujal-suhaas/whyfence/internal/fixtures"
)

func TestRecoveryAccuracy(t *testing.T) {
	history, err := fixtures.NewKnownHistory()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := history.Repository.Cleanup(); err != nil {
			t.Errorf("cleanup accuracy fixture: %v", err)
		}
	})

	cases := []Case{
		{
			Name:     "baseline target reason",
			Expected: history.Baseline.Subject,
			Resolve: func() (string, error) {
				return commitSubject(history.Repository, history.Baseline.Hash)
			},
		},
		{
			Name:     "linked issue reference",
			Expected: "#42",
			Resolve: func() (string, error) {
				subject, err := commitSubject(history.Repository, history.Fix.Hash)
				if err != nil {
					return "", err
				}
				if !strings.Contains(subject, "#42") {
					return "", fmt.Errorf("issue reference missing from %q", subject)
				}
				return "#42", nil
			},
		},
		{
			Name:     "test changed with target",
			Expected: history.TestPath,
			Resolve: func() (string, error) {
				return changedTestPath(history.Repository, history.Fix.Hash, history.TestPath)
			},
		},
		{
			Name: "line-range evidence",
			Expected: citationSummary([]evidence.Citation{
				{Hash: history.Fix.Hash, Subject: history.Fix.Subject},
				{Hash: history.Baseline.Hash, Subject: history.Baseline.Subject},
			}),
			Resolve: func() (string, error) {
				path := filepath.Join(history.Repository.Dir, filepath.FromSlash(history.TargetPath))
				citations, err := evidence.RecoverLineRange(path, 3, 3)
				if err != nil {
					return "", err
				}
				return citationSummary(citations), nil
			},
		},
	}

	report := Evaluate(cases)
	t.Log(report.Summary())
	if !report.Passed() {
		t.Fatalf("accuracy failures: %s", report.FailureSummary())
	}
}

func TestEvaluateCalculatesPartialRecovery(t *testing.T) {
	report := Evaluate([]Case{
		{Name: "recovered", Expected: "known", Resolve: func() (string, error) { return "known", nil }},
		{Name: "missed", Expected: "known", Resolve: func() (string, error) { return "unknown", nil }},
	})
	if report.Summary() != "recovered 1 of 2" {
		t.Fatalf("summary = %q, want %q", report.Summary(), "recovered 1 of 2")
	}
	if report.Passed() {
		t.Fatal("partial recovery should not pass")
	}
}

func commitSubject(repo fixtures.Repository, hash string) (string, error) {
	output, err := repo.Git("show", "-s", "--format=%s", hash)
	return strings.TrimSpace(output), err
}

func changedTestPath(repo fixtures.Repository, hash, expectedPath string) (string, error) {
	output, err := repo.Git("diff-tree", "--no-commit-id", "--name-only", "-r", hash+"^", hash)
	if err != nil {
		return "", err
	}
	for _, path := range strings.Split(strings.TrimSpace(output), "\n") {
		if path == expectedPath {
			return path, nil
		}
	}
	return "", fmt.Errorf("changed test %q not found in %q", expectedPath, output)
}

func citationSummary(citations []evidence.Citation) string {
	parts := make([]string, 0, len(citations))
	for _, citation := range citations {
		parts = append(parts, citation.Hash+" "+citation.Subject)
	}
	return strings.Join(parts, "\n")
}
