package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/benaskins/axon-mind"
)

func main() {
	os.Exit(run())
}

func run() int {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: mind query [-f file.pl]... <goal>")
		return 2
	}

	switch os.Args[1] {
	case "query":
		return runQuery(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		fmt.Fprintln(os.Stderr, "usage: mind query [-f file.pl]... <goal>")
		return 2
	}
}

type fileFlags []string

func (f *fileFlags) String() string { return fmt.Sprintf("%v", *f) }
func (f *fileFlags) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func runQuery(args []string) int {
	fs := flag.NewFlagSet("query", flag.ContinueOnError)
	var files fileFlags
	fs.Var(&files, "f", "Prolog file to load (can be repeated)")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	remaining := fs.Args()
	if len(remaining) != 1 {
		fmt.Fprintln(os.Stderr, "usage: mind query [-f file.pl]... <goal>")
		return 2
	}
	goal := remaining[0]

	// Build engine options from file flags
	opts := make([]mind.Option, 0, len(files))
	for _, f := range files {
		opts = append(opts, mind.WithFile(f))
	}

	e := mind.NewEngine(opts...)

	results, err := e.Query(goal)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}

	if len(results) == 0 {
		fmt.Println("[]")
		return 1
	}

	data, err := mind.SolutionsJSON(results)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 2
	}

	fmt.Println(string(data))
	return 0
}
