package fixtures

import "time"

// KnownHistory describes a small history used to exercise recovery assertions.
type KnownHistory struct {
	Repository Repository
	TargetPath string
	TestPath   string
	Baseline   Commit
	Fix        Commit
}

// NewKnownHistory creates target code and a regression test in separate commits.
func NewKnownHistory() (*KnownHistory, error) {
	repo, err := New()
	if err != nil {
		return nil, err
	}
	history := &KnownHistory{
		Repository: *repo,
		TargetPath: "internal/worker.go",
		TestPath:   "internal/worker_test.go",
	}
	if err := repo.WriteFile(history.TargetPath, "package internal\n\nfunc timeoutSeconds() int { return 30 }\n"); err != nil {
		_ = repo.Cleanup()
		return nil, err
	}
	if err := repo.WriteFile(history.TestPath, "package internal\n\nfunc TestTimeoutSeconds(t any) {}\n"); err != nil {
		_ = repo.Cleanup()
		return nil, err
	}
	history.Baseline, err = repo.Commit("Add worker timeout guard for issue #41", time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC))
	if err != nil {
		_ = repo.Cleanup()
		return nil, err
	}
	if err := repo.WriteFile(history.TargetPath, "package internal\n\nfunc timeoutSeconds() int { return 60 }\n"); err != nil {
		_ = repo.Cleanup()
		return nil, err
	}
	if err := repo.WriteFile(history.TestPath, "package internal\n\nfunc TestTimeoutSeconds(t any) { want := 60 }\n"); err != nil {
		_ = repo.Cleanup()
		return nil, err
	}
	history.Fix, err = repo.Commit("Fix timeout handling (PR #42)", time.Date(2024, 2, 3, 4, 5, 6, 0, time.UTC))
	if err != nil {
		_ = repo.Cleanup()
		return nil, err
	}
	return history, nil
}
