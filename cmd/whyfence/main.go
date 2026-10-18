package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/sujal-suhaas/whyfence/internal/evidence"
	"github.com/sujal-suhaas/whyfence/internal/input"
	"github.com/sujal-suhaas/whyfence/internal/report"
	"github.com/sujal-suhaas/whyfence/internal/risk"
)

const version = "0.1.0-dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	return runWithInput(args, os.Stdin, stdout, stderr)
}

func runWithInput(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("whyfence", flag.ContinueOnError)
	flags.SetOutput(stderr)
	showVersion := flags.Bool("version", false, "print the whyfence version")
	rangeSpec := flags.String("range", "", "an inclusive file range: path:start-end")
	diffSource := flags.String("diff", "", "a unified diff file, or - for stdin")
	reportPath := flags.String("report", "", "write a self-contained HTML report to a named file")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if *showVersion {
		fmt.Fprintf(stdout, "whyfence %s\n", version)
		return 0
	}
	if (*rangeSpec == "") == (*diffSource == "") {
		fmt.Fprintln(stderr, "input: exactly one of --range or --diff is required")
		return 2
	}
	if *reportPath == "-" {
		fmt.Fprintln(stderr, "output: report path must be a named file")
		return 2
	}
	var request input.Request
	var err error
	if *rangeSpec != "" {
		request, err = input.ParseRange(*rangeSpec, os.ReadFile)
	} else {
		request, err = input.ParseDiffSource(*diffSource, os.ReadFile, stdin)
	}
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if request.Mode == input.RangeMode {
		citations, err := evidence.RecoverLineRange(request.Path, request.Start, request.End)
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 2
		}
		if *reportPath != "" {
			level := string(risk.Evaluate(citations))
			return writeReport(*reportPath, report.Data{
				Status:    "accepted",
				Input:     reportInput(request),
				Citations: reportCitations(citations),
				Risk:      &level,
			}, stdout, stderr)
		}
		fmt.Fprintf(stdout, "input: range\nfile: %s\nlines: %d-%d\n", request.Path, request.Start, request.End)
		for _, citation := range citations {
			fmt.Fprintf(stdout, "commit: %s %s\n", citation.Hash, citation.Subject)
			for _, testPath := range citation.TestPaths {
				fmt.Fprintf(stdout, "test: %s\n", testPath)
			}
			for _, reference := range citation.References {
				fmt.Fprintf(stdout, "reference: %s #%d\n", reference.Kind, reference.Number)
			}
		}
		fmt.Fprintf(stdout, "risk: %s\n", risk.Evaluate(citations))
		fmt.Fprintln(stdout, "status: accepted")
		return 0
	}
	if *reportPath != "" {
		return writeReport(*reportPath, report.Data{
			Status: "accepted",
			Input:  reportInput(request),
		}, stdout, stderr)
	}
	fmt.Fprintln(stdout, request.Summary())
	return 0
}

func writeReport(path string, data report.Data, stdout, stderr io.Writer) int {
	contents, err := report.Render(data)
	if err != nil {
		fmt.Fprintf(stderr, "output: cannot render report: %v\n", err)
		return 2
	}
	if err := report.WriteAtomic(path, contents); err != nil {
		fmt.Fprintf(stderr, "output: cannot write report %q: %v\n", path, err)
		return 2
	}
	fmt.Fprintf(stdout, "report: %s\nstatus: accepted\n", path)
	return 0
}

func reportInput(request input.Request) report.Input {
	return report.Input{
		Mode:  string(request.Mode),
		Path:  request.Path,
		Start: request.Start,
		End:   request.End,
		Files: request.Files,
		Hunks: request.Hunks,
	}
}

func reportCitations(citations []evidence.Citation) []report.Citation {
	converted := make([]report.Citation, 0, len(citations))
	for _, citation := range citations {
		references := make([]report.Reference, 0, len(citation.References))
		for _, reference := range citation.References {
			references = append(references, report.Reference{
				Kind:   string(reference.Kind),
				Number: reference.Number,
				Raw:    reference.Raw,
			})
		}
		converted = append(converted, report.Citation{
			Hash:       citation.Hash,
			Subject:    citation.Subject,
			TestPaths:  append([]string(nil), citation.TestPaths...),
			References: references,
		})
	}
	return converted
}
