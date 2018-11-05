package main

import (
	"bufio"
	"fmt"
	"fzf-bibtex/cache"
	"fzf-bibtex/format"
	"fzf-bibtex/startup"
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
	// read keys from stdin into a slice
	keys := readKeysFromStdin()
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

func readKeysFromStdin() []string {
	var keys []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()
		sl := strings.Split(s, "@")
		keys = append(keys, "@"+sl[len(sl)-1])
	}
	if scanner.Err() != nil {
		// handle error.
		fmt.Fprintln(os.Stderr, "reading standard input:", scanner.Err())
	}
	if debug {
		fmt.Println(keys)
	}
	return keys
}
