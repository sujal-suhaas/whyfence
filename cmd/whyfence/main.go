package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/sujal-suhaas/whyfence/internal/evidence"
	"github.com/sujal-suhaas/whyfence/internal/input"
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
		fmt.Fprintf(stdout, "input: range\nfile: %s\nlines: %d-%d\n", request.Path, request.Start, request.End)
		for _, citation := range citations {
			fmt.Fprintf(stdout, "commit: %s %s\n", citation.Hash, citation.Subject)
			for _, reference := range citation.References {
				fmt.Fprintf(stdout, "reference: %s #%d\n", reference.Kind, reference.Number)
			}
		}
		fmt.Fprintln(stdout, "status: accepted")
		return 0
	}
	fmt.Fprintln(stdout, request.Summary())
	return 0
}
