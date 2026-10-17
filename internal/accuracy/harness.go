// Package accuracy provides deterministic accounting for recovery cases.
package accuracy

import (
	"fmt"
	"strings"
)

// Case is one expected-reason recovery check supplied by a test-time resolver.
type Case struct {
	Name     string
	Expected string
	Resolve  func() (string, error)
}

// Failure records why one case did not recover its expected reason.
type Failure struct {
	Name     string
	Expected string
	Actual   string
	Err      error
}

// Report contains calculated accuracy results for a set of cases.
type Report struct {
	Recovered int
	Total     int
	Failures  []Failure
}

// Evaluate runs every resolver and counts exact expected-reason matches.
func Evaluate(cases []Case) Report {
	report := Report{Total: len(cases)}
	for _, testCase := range cases {
		failure := Failure{Name: testCase.Name, Expected: testCase.Expected}
		if testCase.Resolve == nil {
			failure.Err = fmt.Errorf("resolver is nil")
			report.Failures = append(report.Failures, failure)
			continue
		}
		actual, err := testCase.Resolve()
		if err != nil {
			failure.Err = err
			report.Failures = append(report.Failures, failure)
			continue
		}
		failure.Actual = actual
		if actual != testCase.Expected {
			report.Failures = append(report.Failures, failure)
			continue
		}
		report.Recovered++
	}
	return report
}

// Summary returns the stable output used by offline accuracy checks.
func (r Report) Summary() string {
	return fmt.Sprintf("recovered %d of %d", r.Recovered, r.Total)
}

// Passed reports whether every case recovered its expected reason.
func (r Report) Passed() bool {
	return r.Recovered == r.Total && len(r.Failures) == 0
}

// FailureSummary formats failures without making the successful count implicit.
func (r Report) FailureSummary() string {
	if len(r.Failures) == 0 {
		return ""
	}
	parts := make([]string, 0, len(r.Failures))
	for _, failure := range r.Failures {
		message := fmt.Sprintf("%s: expected %q, got %q", failure.Name, failure.Expected, failure.Actual)
		if failure.Err != nil {
			message += ": " + failure.Err.Error()
		}
		parts = append(parts, message)
	}
	return strings.Join(parts, "; ")
}
