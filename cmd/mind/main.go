package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: mind query [-f file.pl]... <goal>")
		os.Exit(2)
	}
	fmt.Fprintln(os.Stderr, "not yet implemented")
	os.Exit(2)
}
