package risk

import (
	"reflect"
	"testing"

	"github.com/sujal-suhaas/whyfence/internal/evidence"
)

func TestEvaluateRiskLevels(t *testing.T) {
	tests := []struct {
		name      string
		citations []evidence.Citation
		want      Level
	}{
		{name: "empty evidence is conservative", want: High},
		{name: "no supporting signals", citations: []evidence.Citation{{Hash: "low"}}, want: Low},
		{name: "tests only", citations: []evidence.Citation{{TestPaths: []string{"internal/worker_test.go"}}}, want: Medium},
		{name: "references only", citations: []evidence.Citation{{References: []evidence.Reference{{Kind: evidence.IssueReference, Number: 41}}}}, want: Medium},
		{name: "unclassified references only", citations: []evidence.Citation{{References: []evidence.Reference{{Kind: evidence.UnclassifiedReference, Number: 18}}}}, want: Medium},
		{name: "tests and references", citations: []evidence.Citation{{TestPaths: []string{"internal/worker_test.go"}, References: []evidence.Reference{{Kind: evidence.PullRequestReference, Number: 42}}}}, want: High},
		{name: "signals across citations", citations: []evidence.Citation{{TestPaths: []string{"internal/worker_test.go"}}, {References: []evidence.Reference{{Kind: evidence.IssueReference, Number: 41}}}}, want: High},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			if got := Evaluate(testCase.citations); got != testCase.want {
				t.Fatalf("risk = %q, want %q", got, testCase.want)
			}
		})
	}
}

func TestEvaluateDoesNotMutateCitations(t *testing.T) {
	citations := []evidence.Citation{{
		Hash:       "abc",
		TestPaths:  []string{"internal/worker_test.go"},
		References: []evidence.Reference{{Kind: evidence.IssueReference, Number: 41}},
	}}
	before := append([]evidence.Citation(nil), citations...)
	if got := Evaluate(citations); got != High {
		t.Fatalf("risk = %q, want %q", got, High)
	}
	if !reflect.DeepEqual(citations, before) {
		t.Fatalf("citations changed from %#v to %#v", before, citations)
	}
}
