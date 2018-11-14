package main

import (
	"fmt"
	"github.com/msprev/fzf-bibtex/cache"
	"github.com/msprev/fzf-bibtex/format"
	"github.com/msprev/fzf-bibtex/startup"
	"os"
	"strings"
)

const usage string = `bibtex-markdown [-cache=...] [file1.bib file2.bib ...]
  pretty print items (in markdown) for selected .bib entries passed over stdin.
`

const debug = false

func main() {
	cacheDir, bibFiles := startup.ReadArgs(usage)
	if debug {
		fmt.Println("cachedir: ", cacheDir)
		fmt.Println("bib files: ", bibFiles)
	}
	keys := startup.ReadKeysFromStdin()
	if len(keys) == 0 {
		os.Exit(-1)
	}
	// create a key printer function for read keys
	printIfKeyMatches := makePrinter(keys)
	// pass it to cache-backed markdown outputing function
	for _, bibFile := range bibFiles {
		markdown(cacheDir, bibFile, printIfKeyMatches)
	}
}

func markdown(cacheDir string, bibFile string, printIfKeyMatches func(string)) {
	if debug {
		fmt.Println("markdown " + bibFile)
	}
	cache.ReadAndDo(cacheDir, bibFile, "markdown", format.EntryToMarkdown, printIfKeyMatches)
}

func makePrinter(keys []string) func(string) {
	return func(s string) {
		for _, k := range keys {
			if strings.HasPrefix(s, k+" ") {
				sl := strings.SplitN(s, k+" ", 2)
				fmt.Println(sl[1])
			}
		}
	}
}
