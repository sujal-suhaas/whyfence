package evidence

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/sujal-suhaas/whyfence/internal/fixtures"
)

func TestLineRange(t *testing.T) {
	history := knownHistory(t)

	path := filepath.Join(history.Repository.Dir, filepath.FromSlash(history.TargetPath))
	want := []Citation{
		{
			Hash:      history.CodeOnly.Hash,
			Subject:   history.CodeOnly.Subject,
			TestPaths: nil,
		},
		{
			Hash:       history.Fix.Hash,
			Subject:    history.Fix.Subject,
			References: []Reference{{Kind: PullRequestReference, Number: 42, Raw: "PR #42"}},
			TestPaths:  []string{history.TestPath},
		},
		{
			Hash:       history.Baseline.Hash,
			Subject:    history.Baseline.Subject,
			References: []Reference{{Kind: IssueReference, Number: 41, Raw: "issue #41"}},
			TestPaths:  []string{history.TestPath},
		},
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

func TestChangedTests(t *testing.T) {
	history := knownHistory(t)
	path := filepath.Join(history.Repository.Dir, filepath.FromSlash(history.TargetPath))
	citations, err := RecoverLineRange(path, 3, 3)
	if err != nil {
		t.Fatal(err)
	}

	if got := citationByHash(citations, history.Fix.Hash).TestPaths; !reflect.DeepEqual(got, []string{history.TestPath}) {
		t.Fatalf("matching test paths = %#v, want %#v", got, []string{history.TestPath})
	}
	if got := citationByHash(citations, history.Baseline.Hash).TestPaths; !reflect.DeepEqual(got, []string{history.TestPath}) {
		t.Fatalf("root-commit test paths = %#v, want %#v", got, []string{history.TestPath})
	}
	if got := citationByHash(citations, history.CodeOnly.Hash).TestPaths; len(got) != 0 {
		t.Fatalf("non-matching test paths = %#v, want none", got)
	}

	repeated, err := RecoverLineRange(path, 3, 3)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(repeated, citations) {
		t.Fatalf("repeated changed-test evidence = %#v, want %#v", repeated, citations)
	}
}

func TestIssueReference(t *testing.T) {
	history := knownHistory(t)
	path := filepath.Join(history.Repository.Dir, filepath.FromSlash(history.TargetPath))
	citations, err := RecoverLineRange(path, 3, 3)
	if err != nil {
		t.Fatal(err)
	}
	citation := citationByHash(citations, history.Baseline.Hash)
	want := []Reference{{Kind: IssueReference, Number: 41, Raw: "issue #41"}}
	if !reflect.DeepEqual(citation.References, want) {
		t.Fatalf("issue references = %#v, want %#v", citation.References, want)
	}
}

func TestPRReference(t *testing.T) {
	history := knownHistory(t)
	path := filepath.Join(history.Repository.Dir, filepath.FromSlash(history.TargetPath))
	citations, err := RecoverLineRange(path, 3, 3)
	if err != nil {
		t.Fatal(err)
	}
	citation := citationByHash(citations, history.Fix.Hash)
	want := []Reference{{Kind: PullRequestReference, Number: 42, Raw: "PR #42"}}
	if !reflect.DeepEqual(citation.References, want) {
		t.Fatalf("PR references = %#v, want %#v", citation.References, want)
	}
}

func TestExtractReferences(t *testing.T) {
	message := "issue #41, PR: #42, Merge pull request #44 from acme/topic, https://github.com/acme/tool/issues/43, https://github.com/acme/tool/pull/45, issue #41, #99"
	want := []Reference{
		{Kind: IssueReference, Number: 41, Raw: "issue #41"},
		{Kind: PullRequestReference, Number: 42, Raw: "PR: #42"},
		{Kind: PullRequestReference, Number: 44, Raw: "Merge pull request #44"},
		{Kind: IssueReference, Number: 43, Raw: "https://github.com/acme/tool/issues/43"},
		{Kind: PullRequestReference, Number: 45, Raw: "https://github.com/acme/tool/pull/45"},
	}
	if got := ExtractReferences(message); !reflect.DeepEqual(got, want) {
		t.Fatalf("references = %#v, want %#v", got, want)
	}

	if got := ExtractReferences("unresolved #99 and acme/tool#100"); len(got) != 0 {
		t.Fatalf("ambiguous references = %#v, want none", got)
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

func knownHistory(t *testing.T) *fixtures.KnownHistory {
	t.Helper()
	history, err := fixtures.NewKnownHistory()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := history.Repository.Cleanup(); err != nil {
			t.Errorf("cleanup evidence fixture: %v", err)
		}
	})
	return history
}

func citationByHash(citations []Citation, hash string) Citation {
	for _, citation := range citations {
		if citation.Hash == hash {
			return citation
		}
	}
	return Citation{}
}
