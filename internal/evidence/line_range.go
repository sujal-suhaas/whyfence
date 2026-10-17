// Package evidence recovers deterministic citations from local Git history.
package evidence

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// ReferenceKind identifies the kind of tracker object named by a commit subject.
type ReferenceKind string

const (
	// IssueReference identifies an issue reference.
	IssueReference ReferenceKind = "issue"
	// PullRequestReference identifies a pull request reference.
	PullRequestReference ReferenceKind = "pr"
)

// Reference identifies one explicit issue or pull request token in a subject.
type Reference struct {
	Kind   ReferenceKind
	Number int
	Raw    string
}

// Citation identifies one local Git commit relevant to a requested line range.
type Citation struct {
	Hash       string
	Subject    string
	References []Reference
	TestPaths  []string
}

type referencePattern struct {
	kind    ReferenceKind
	pattern *regexp.Regexp
}

var referencePatterns = []referencePattern{
	{kind: IssueReference, pattern: regexp.MustCompile(`(?i)\bissue\s*(?::\s*|\s+)#([0-9]+)\b`)},
	{kind: IssueReference, pattern: regexp.MustCompile(`(?i)\bhttps://github\.com/[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+/issues/([0-9]+)\b`)},
	{kind: PullRequestReference, pattern: regexp.MustCompile(`(?i)\b(?:pr|pull request)\s*(?::\s*|\s+)#([0-9]+)\b`)},
	{kind: PullRequestReference, pattern: regexp.MustCompile(`(?i)\bmerge\s+pull\s+request\s+#([0-9]+)\b`)},
	{kind: PullRequestReference, pattern: regexp.MustCompile(`(?i)\bhttps://github\.com/[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+/pull/([0-9]+)\b`)},
}

type referenceMatch struct {
	reference Reference
	start     int
}

// ExtractReferences parses supported explicit issue and pull request forms from
// one local commit subject without querying a remote tracker.
func ExtractReferences(subject string) []Reference {
	matches := make([]referenceMatch, 0)
	for _, referencePattern := range referencePatterns {
		for _, indexes := range referencePattern.pattern.FindAllStringSubmatchIndex(subject, -1) {
			if len(indexes) < 4 {
				continue
			}
			number, err := strconv.Atoi(subject[indexes[2]:indexes[3]])
			if err != nil || number < 1 {
				continue
			}
			matches = append(matches, referenceMatch{
				reference: Reference{
					Kind:   referencePattern.kind,
					Number: number,
					Raw:    strings.TrimSpace(subject[indexes[0]:indexes[1]]),
				},
				start: indexes[0],
			})
		}
	}
	sort.SliceStable(matches, func(i, j int) bool {
		return matches[i].start < matches[j].start
	})
	if len(matches) == 0 {
		return nil
	}

	references := make([]Reference, 0, len(matches))
	seen := make(map[string]struct{})
	for _, match := range matches {
		key := fmt.Sprintf("%s:%d", match.reference.Kind, match.reference.Number)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		references = append(references, match.reference)
	}
	return references
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
	for index := range citations {
		citations[index].TestPaths, err = changedTestPaths(repoRoot, citations[index].Hash)
		if err != nil {
			return nil, err
		}
	}
	return citations, nil
}

func changedTestPaths(repoRoot, hash string) ([]string, error) {
	command := exec.Command("git", "--no-pager", "-C", repoRoot, "diff-tree", "--root", "--no-commit-id", "--name-only", "--no-renames", "-r", "-z", hash)
	command.Env = append(os.Environ(),
		"GIT_CONFIG_NOSYSTEM=1",
		"GIT_TERMINAL_PROMPT=0",
		"LC_ALL=C",
		"LANG=C",
	)
	output, err := command.CombinedOutput()
	if err != nil {
		message := strings.TrimSpace(string(output))
		if message == "" {
			message = err.Error()
		}
		return nil, fmt.Errorf("evidence: Git changed-path recovery failed for %s: %s", hash, message)
	}

	var testPaths []string
	for _, changedPath := range strings.Split(strings.TrimSuffix(string(output), "\x00"), "\x00") {
		if isTestPath(changedPath) {
			testPaths = append(testPaths, changedPath)
		}
	}
	sort.Strings(testPaths)
	return testPaths, nil
}

func isTestPath(repoPath string) bool {
	return repoPath != "" && strings.HasSuffix(path.Base(repoPath), "_test.go")
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
		citations = append(citations, Citation{
			Hash:       fields[0],
			Subject:    fields[1],
			References: ExtractReferences(fields[1]),
		})
	}
	return citations
}
