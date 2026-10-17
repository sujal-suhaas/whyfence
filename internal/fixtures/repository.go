// Package fixtures builds deterministic Git repositories for offline tests.
package fixtures

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	// FixedAuthorName prevents a developer's global Git identity from changing fixture history.
	FixedAuthorName = "whyfence fixture"
	// FixedAuthorEmail keeps generated commit metadata stable without using a real address.
	FixedAuthorEmail = "fixture@whyfence.local"
)

// Commit records the stable identity and timestamp of a generated commit.
type Commit struct {
	Hash    string
	Subject string
	When    time.Time
}

// Repository is a temporary Git repository created for one test run.
type Repository struct {
	Dir string
}

// New creates a Git repository outside the workspace so fixtures cannot become nested repos.
func New() (*Repository, error) {
	dir, err := os.MkdirTemp("", "whyfence-fixture-")
	if err != nil {
		return nil, fmt.Errorf("create fixture directory: %w", err)
	}
	repo := &Repository{Dir: dir}
	if _, err := repo.Git("init", "--quiet"); err != nil {
		_ = repo.Cleanup()
		return nil, err
	}
	for _, setting := range [][2]string{
		{"user.name", FixedAuthorName},
		{"user.email", FixedAuthorEmail},
	} {
		if _, err := repo.Git("config", setting[0], setting[1]); err != nil {
			_ = repo.Cleanup()
			return nil, err
		}
	}
	return repo, nil
}

// Cleanup removes the temporary repository and is safe to call more than once.
func (r *Repository) Cleanup() error {
	if r == nil || r.Dir == "" {
		return nil
	}
	if err := os.RemoveAll(r.Dir); err != nil {
		return fmt.Errorf("remove fixture repository %q: %w", r.Dir, err)
	}
	return nil
}

// WriteFile writes a fixture file while rejecting paths that escape the repository.
func (r *Repository) WriteFile(name, contents string) error {
	path, err := r.filePath(name)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create parent for %q: %w", name, err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		return fmt.Errorf("write fixture file %q: %w", name, err)
	}
	return nil
}

// Git runs a Git command in the temporary repository and returns combined output.
func (r *Repository) Git(args ...string) (string, error) {
	if r == nil || r.Dir == "" {
		return "", fmt.Errorf("fixture repository is nil")
	}
	command := exec.Command("git", args...)
	command.Dir = r.Dir
	command.Env = append(os.Environ(),
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_TERMINAL_PROMPT=0",
	)
	output, err := command.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(string(output)))
	}
	return string(output), nil
}

// Commit stages all fixture changes and records them with fixed identity and time metadata.
func (r *Repository) Commit(subject string, when time.Time) (Commit, error) {
	if strings.TrimSpace(subject) == "" {
		return Commit{}, fmt.Errorf("commit subject cannot be empty")
	}
	if when.IsZero() {
		return Commit{}, fmt.Errorf("commit time cannot be zero")
	}
	if _, err := r.Git("add", "--all"); err != nil {
		return Commit{}, err
	}
	command := exec.Command("git", "commit", "--quiet", "--message", subject)
	command.Dir = r.Dir
	date := when.UTC().Format(time.RFC3339)
	command.Env = append(os.Environ(),
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_TERMINAL_PROMPT=0",
		"GIT_AUTHOR_NAME="+FixedAuthorName,
		"GIT_AUTHOR_EMAIL="+FixedAuthorEmail,
		"GIT_AUTHOR_DATE="+date,
		"GIT_COMMITTER_NAME="+FixedAuthorName,
		"GIT_COMMITTER_EMAIL="+FixedAuthorEmail,
		"GIT_COMMITTER_DATE="+date,
	)
	output, err := command.CombinedOutput()
	if err != nil {
		return Commit{}, fmt.Errorf("git commit %q: %w: %s", subject, err, strings.TrimSpace(string(output)))
	}
	hashOutput, err := r.Git("rev-parse", "HEAD")
	if err != nil {
		return Commit{}, err
	}
	return Commit{Hash: strings.TrimSpace(hashOutput), Subject: subject, When: when.UTC()}, nil
}

func (r *Repository) filePath(name string) (string, error) {
	if name == "" || filepath.IsAbs(name) {
		return "", fmt.Errorf("fixture file path must be relative: %q", name)
	}
	clean := filepath.Clean(filepath.FromSlash(name))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("fixture file path escapes repository: %q", name)
	}
	return filepath.Join(r.Dir, clean), nil
}
