package fixtures

import "time"

// RiskHistory contains one generated target for each evidence-backed risk level.
// Each target is committed separately so line-range recovery isolates its case.
type RiskHistory struct {
	Repository    Repository
	LowPath       string
	TestOnlyPath  string
	ReferencePath string
	HighPath      string
	Low           Commit
	TestOnly      Commit
	Reference     Commit
	High          Commit
}

// NewRiskHistory creates low, medium, and high evidence scenarios using local
// Git history with fixed metadata.
func NewRiskHistory() (*RiskHistory, error) {
	repo, err := New()
	if err != nil {
		return nil, err
	}
	history := &RiskHistory{
		Repository:    *repo,
		LowPath:       "risk/low.go",
		TestOnlyPath:  "risk/test_only.go",
		ReferencePath: "risk/reference.go",
		HighPath:      "risk/high.go",
	}
	commits := []struct {
		path       string
		contents   string
		testPath   string
		testSource string
		subject    string
		when       time.Time
		result     *Commit
	}{
		{
			path:     history.LowPath,
			contents: "package risk\n\nfunc Low() int { return 1 }\n",
			subject:  "Tune low-risk path",
			when:     time.Date(2024, 4, 1, 1, 2, 3, 0, time.UTC),
			result:   &history.Low,
		},
		{
			path:       history.TestOnlyPath,
			contents:   "package risk\n\nfunc TestOnly() int { return 2 }\n",
			testPath:   "risk/test_only_test.go",
			testSource: "package risk\n\nfunc TestTestOnly(t any) {}\n",
			subject:    "Add regression coverage",
			when:       time.Date(2024, 5, 2, 2, 3, 4, 0, time.UTC),
			result:     &history.TestOnly,
		},
		{
			path:     history.ReferencePath,
			contents: "package risk\n\nfunc Reference() int { return 3 }\n",
			subject:  "Document behavior for issue #73",
			when:     time.Date(2024, 6, 3, 3, 4, 5, 0, time.UTC),
			result:   &history.Reference,
		},
		{
			path:       history.HighPath,
			contents:   "package risk\n\nfunc High() int { return 4 }\n",
			testPath:   "risk/high_test.go",
			testSource: "package risk\n\nfunc TestHigh(t any) {}\n",
			subject:    "Add guarded path in PR #74",
			when:       time.Date(2024, 7, 4, 4, 5, 6, 0, time.UTC),
			result:     &history.High,
		},
	}
	for _, commit := range commits {
		if err := repo.WriteFile(commit.path, commit.contents); err != nil {
			_ = repo.Cleanup()
			return nil, err
		}
		if commit.testPath != "" {
			if err := repo.WriteFile(commit.testPath, commit.testSource); err != nil {
				_ = repo.Cleanup()
				return nil, err
			}
		}
		created, err := repo.Commit(commit.subject, commit.when)
		if err != nil {
			_ = repo.Cleanup()
			return nil, err
		}
		*commit.result = created
	}
	return history, nil
}
