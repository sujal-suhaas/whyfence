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
			Name:     "commit subject reference",
			Expected: "PR #42",
			Resolve: func() (string, error) {
				subject, err := commitSubject(history.Repository, history.Fix.Hash)
				if err != nil {
					return "", err
				}
				if !strings.Contains(subject, "PR #42") {
					return "", fmt.Errorf("PR reference missing from %q", subject)
				}
				return "PR #42", nil
			},
		},
		{
			Name:     "test changed with target",
			Expected: history.TestPath,
			Resolve: func() (string, error) {
				citations, err := lineRangeCitations(history)
				if err != nil {
					return "", err
				}
				return citationTestPath(citations, history.Fix.Hash)
			},
		},
		{
			Name: "line-range evidence",
			Expected: citationSummary([]evidence.Citation{
				{
					Hash:    history.CodeOnly.Hash,
					Subject: history.CodeOnly.Subject,
				},
				{
					Hash:       history.Fix.Hash,
					Subject:    history.Fix.Subject,
					References: []evidence.Reference{{Kind: evidence.PullRequestReference, Number: 42, Raw: "PR #42"}},
					TestPaths:  []string{history.TestPath},
				},
				{
					Hash:       history.Baseline.Hash,
					Subject:    history.Baseline.Subject,
					References: []evidence.Reference{{Kind: evidence.IssueReference, Number: 41, Raw: "issue #41"}},
					TestPaths:  []string{history.TestPath},
				},
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
		{
			Name:     "code changed without test",
			Expected: "none",
			Resolve: func() (string, error) {
				citations, err := lineRangeCitations(history)
				if err != nil {
					return "", err
				}
				for _, citation := range citations {
					if citation.Hash == history.CodeOnly.Hash && len(citation.TestPaths) != 0 {
						return "", fmt.Errorf("unexpected test paths on code-only commit: %v", citation.TestPaths)
					}
				}
				return "none", nil
			},
		},
		{
			Name:     "production issue reference",
			Expected: "issue #41",
			Resolve: func() (string, error) {
				citations, err := lineRangeCitations(history)
				if err != nil {
					return "", err
				}
				return citationReference(citations, history.Baseline.Hash, evidence.IssueReference, 41)
			},
		},
		{
			Name:     "production PR reference",
			Expected: "pr #42",
			Resolve: func() (string, error) {
				citations, err := lineRangeCitations(history)
				if err != nil {
					return "", err
				}
				return citationReference(citations, history.Fix.Hash, evidence.PullRequestReference, 42)
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

func citationSummary(citations []evidence.Citation) string {
	parts := make([]string, 0, len(citations))
	for _, citation := range citations {
		part := citation.Hash + " " + citation.Subject
		for _, testPath := range citation.TestPaths {
			part += " [test " + testPath + "]"
		}
		for _, reference := range citation.References {
			part += fmt.Sprintf(" [%s #%d]", reference.Kind, reference.Number)
		}
		parts = append(parts, part)
	}
	return strings.Join(parts, "\n")
}

func lineRangeCitations(history *fixtures.KnownHistory) ([]evidence.Citation, error) {
	path := filepath.Join(history.Repository.Dir, filepath.FromSlash(history.TargetPath))
	return evidence.RecoverLineRange(path, 3, 3)
}

func citationTestPath(citations []evidence.Citation, hash string) (string, error) {
	for _, citation := range citations {
		if citation.Hash != hash {
			continue
		}
		if len(citation.TestPaths) != 1 {
			return "", fmt.Errorf("commit %s has test paths %v", hash, citation.TestPaths)
		}
		return citation.TestPaths[0], nil
	}
	return "", fmt.Errorf("commit %s not found", hash)
}

func citationReference(citations []evidence.Citation, hash string, kind evidence.ReferenceKind, number int) (string, error) {
	for _, citation := range citations {
		if citation.Hash != hash {
			continue
		}
		for _, reference := range citation.References {
			if reference.Kind == kind && reference.Number == number {
				return fmt.Sprintf("%s #%d", reference.Kind, reference.Number), nil
			}
		}
	}
	return "", fmt.Errorf("reference %s #%d not found on commit %s", kind, number, hash)
}
