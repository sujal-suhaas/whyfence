// Package evidence recovers deterministic citations from local Git history.
package evidence

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Citation identifies one local Git commit relevant to a requested line range.
type Citation struct {
	Hash    string
	Subject string
}

// RecoverLineRange traces path's inclusive 1-based line range through local Git
// history and returns citations in Git's newest-to-oldest deterministic order.
func RecoverLineRange(path string, start, end int) ([]Citation, error) {
	if start < 1 || end < start {
		return nil, errors.New("evidence: line range must use positive start-end lines")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("evidence: cannot resolve file %q", path)
	}
	info, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("evidence: cannot inspect file %q", path)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("evidence: path %q is a directory", path)
	}

	repoRoot, err := findRepositoryRoot(filepath.Dir(absPath))
	if err != nil {
		return nil, err
	}
	relPath, err := filepath.Rel(repoRoot, absPath)
	if err != nil || relPath == "." || relPath == ".." || strings.HasPrefix(relPath, ".."+string(filepath.Separator)) {
		return nil, fmt.Errorf("evidence: file %q is outside the Git repository", path)
	}
	relPath = filepath.ToSlash(relPath)

	lineSpec := fmt.Sprintf("-L%d,%d:%s", start, end, relPath)
	command := exec.Command("git", "--no-pager", "-C", repoRoot, "log", "--no-ext-diff", "--no-patch", "--format=%H%x09%s", lineSpec)
	command.Env = append(os.Environ(),
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_TERMINAL_PROMPT=0",
		"GIT_PAGER=cat",
		"LC_ALL=C",
		"LANG=C",
	)
	output, err := command.CombinedOutput()
	if err != nil {
		message := strings.TrimSpace(string(output))
		if message == "" {
			message = err.Error()
		}
		return nil, fmt.Errorf("evidence: Git line-range recovery failed: %s", message)
	}

	citations := parseCitations(string(output))
	if len(citations) == 0 {
		return nil, fmt.Errorf("evidence: no Git commits found for %s:%d-%d", relPath, start, end)
	}
	return citations, nil
}

func findRepositoryRoot(directory string) (string, error) {
	command := exec.Command("git", "--no-pager", "-C", directory, "rev-parse", "--show-toplevel")
	command.Env = append(os.Environ(),
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_TERMINAL_PROMPT=0",
		"LC_ALL=C",
		"LANG=C",
	)
	output, err := command.CombinedOutput()
	if err != nil {
		return "", errors.New("evidence: file is not inside a Git repository")
	}
	root := strings.TrimSpace(string(output))
	if root == "" {
		return "", errors.New("evidence: Git repository root is empty")
	}
	return filepath.Clean(root), nil
}

func parseCitations(output string) []Citation {
	lines := strings.Split(output, "\n")
	citations := make([]Citation, 0, len(lines))
	seen := make(map[string]struct{})
	for _, line := range lines {
		if line == "" {
			continue
		}
		fields := strings.SplitN(line, "\t", 2)
		if len(fields) != 2 || fields[0] == "" {
			continue
		}
		if _, ok := seen[fields[0]]; ok {
			continue
		}
		seen[fields[0]] = struct{}{}
		citations = append(citations, Citation{Hash: fields[0], Subject: fields[1]})
	}
	return citations
}
