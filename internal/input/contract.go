// Package input defines and validates whyfence's deterministic CLI input contract.
package input

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strconv"
	"strings"
)

// Mode identifies the accepted input form.
type Mode string

const (
	// RangeMode accepts a path and inclusive 1-based line range.
	RangeMode Mode = "range"
	// DiffMode accepts a unified diff from a file or standard input.
	DiffMode Mode = "diff"
)

// Request is a validated CLI input summary.
type Request struct {
	Mode  Mode
	Path  string
	Start int
	End   int
	Files int
	Hunks int
}

// ParseRange validates path:start-end and the referenced file's line bounds.
func ParseRange(spec string, readFile func(string) ([]byte, error)) (Request, error) {
	separator := strings.LastIndex(spec, ":")
	if separator <= 0 || separator == len(spec)-1 {
		return Request{}, errors.New("input: range must use path:start-end")
	}
	path := spec[:separator]
	span := strings.Split(spec[separator+1:], "-")
	if len(span) != 2 || span[0] == "" || span[1] == "" {
		return Request{}, errors.New("input: range must use path:start-end")
	}
	start, startErr := strconv.Atoi(span[0])
	end, endErr := strconv.Atoi(span[1])
	if startErr != nil || endErr != nil || start < 1 || end < start {
		return Request{}, errors.New("input: range must use positive start-end lines")
	}
	contents, err := readFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Request{}, fmt.Errorf("input: file %q not found", path)
		}
		return Request{}, fmt.Errorf("input: cannot read file %q", path)
	}
	lineCount := countLines(contents)
	if lineCount == 0 {
		return Request{}, fmt.Errorf("input: file %q is empty", path)
	}
	if end > lineCount {
		return Request{}, fmt.Errorf("input: range %d-%d exceeds file %q with %d lines", start, end, path, lineCount)
	}
	return Request{Mode: RangeMode, Path: path, Start: start, End: end}, nil
}

// ParseDiffSource reads and validates a diff source. A source of - means stdin.
func ParseDiffSource(source string, readFile func(string) ([]byte, error), stdin io.Reader) (Request, error) {
	var contents []byte
	var err error
	if source == "-" {
		contents, err = io.ReadAll(stdin)
	} else {
		contents, err = readFile(source)
	}
	if err != nil {
		if source != "-" && errors.Is(err, fs.ErrNotExist) {
			return Request{}, fmt.Errorf("input: diff %q not found", source)
		}
		if source == "-" {
			return Request{}, errors.New("input: cannot read diff from stdin")
		}
		return Request{}, fmt.Errorf("input: cannot read diff %q", source)
	}
	files, hunks, err := parseDiff(string(contents))
	if err != nil {
		return Request{}, err
	}
	return Request{Mode: DiffMode, Path: source, Files: files, Hunks: hunks}, nil
}

// Summary returns stable terminal output for a validated request.
func (r Request) Summary() string {
	if r.Mode == RangeMode {
		return fmt.Sprintf("input: range\nfile: %s\nlines: %d-%d\nstatus: accepted", r.Path, r.Start, r.End)
	}
	return fmt.Sprintf("input: diff\nsource: %s\nfiles: %d\nhunks: %d\nstatus: accepted", r.Path, r.Files, r.Hunks)
}

func countLines(contents []byte) int {
	if len(contents) == 0 {
		return 0
	}
	text := string(contents)
	lines := strings.Count(text, "\n")
	if !strings.HasSuffix(text, "\n") {
		lines++
	}
	return lines
}

func parseDiff(contents string) (int, int, error) {
	if strings.TrimSpace(contents) == "" {
		return 0, 0, errors.New("input: diff is empty")
	}
	scanner := bufio.NewScanner(strings.NewReader(contents))
	diffHeaders, oldFileHeaders, hunks := 0, 0, 0
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "diff --git "):
			diffHeaders++
		case strings.HasPrefix(line, "--- "):
			oldFileHeaders++
		case strings.HasPrefix(line, "@@ ") || strings.HasPrefix(line, "@@@ "):
			hunks++
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, 0, errors.New("input: cannot read diff")
	}
	files := diffHeaders
	if files == 0 {
		files = oldFileHeaders
	}
	if files == 0 || hunks == 0 {
		return 0, 0, errors.New("input: malformed unified diff")
	}
	return files, hunks, nil
}
