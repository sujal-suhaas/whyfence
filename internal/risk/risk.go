// Package risk classifies evidence-backed removal risk deterministically.
package risk

import "github.com/sujal-suhaas/whyfence/internal/evidence"

// Level is the stable risk label emitted by whyfence.
type Level string

const (
	// Low means the recovered evidence contains neither changed tests nor
	// explicit issue or pull request references.
	Low Level = "low"
	// Medium means exactly one of changed tests or explicit references exists.
	Medium Level = "medium"
	// High means both changed tests and explicit references exist, or evidence is
	// empty and therefore cannot establish a safe removal context.
	High Level = "high"
)

// Evaluate classifies the complete recovered range without mutating citations.
// The result reflects evidence-backed removal risk, not runtime or business
// criticality.
func Evaluate(citations []evidence.Citation) Level {
	hasTests, hasReferences := false, false
	for _, citation := range citations {
		if len(citation.TestPaths) > 0 {
			hasTests = true
		}
		if len(citation.References) > 0 {
			hasReferences = true
		}
	}
	if len(citations) == 0 || hasTests && hasReferences {
		return High
	}
	if hasTests || hasReferences {
		return Medium
	}
	return Low
}
