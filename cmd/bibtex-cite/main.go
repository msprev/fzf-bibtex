package main

import (
	"flag"
	"fmt"
	"fzf-bibtex/startup"
	"os"
	"strings"
)

const usage string = `bibtex-cite 
	pretty print citations (in pandoc '@' format) for selected .bib entries passed over stdin.
`

const debug = false

func main() {
	readArgs(usage)
	keys := startup.ReadKeysFromStdin()
	if len(keys) == 0 {
		os.Exit(-1)
	}
	fmt.Print(strings.Join(keys, "; "))
}

func readArgs(usage string) {
	flag.Usage = func() {
		fmt.Println(usage)
		flag.PrintDefaults()
	}
	flag.Parse()
}
