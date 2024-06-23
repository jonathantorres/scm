package main

import (
	"fmt"
	"os"
)

func main() {
	var filename string

	if len(os.Args) <= 1 {
		fmt.Fprintf(os.Stderr, "the file to run was not specified\n")
		os.Exit(1)
	}

	// open the file to parse
	filename = os.Args[1]

	tree, err := openAndParse(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("error opening and parsing: %s\n", err))
		os.Exit(1)
	}

	// start evaluation
	startEval(make_begin(tree))
}
