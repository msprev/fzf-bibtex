package main

import (
	"fmt"
	"fzf-bibtex/cache"
	"fzf-bibtex/format"
	"fzf-bibtex/startup"
)

const usage string = `bibtex-ls [-cache=...] [file1.bib file2.bib ...]
  lists to stdout the content of .bib files, one record per line
`

const debug = false

func main() {
	cacheDir, bibFiles := startup.ReadArgs(usage)
	if debug {
		fmt.Println("cachedir:", cacheDir)
		fmt.Println("bib files: ", bibFiles)
	}
	for _, bibFile := range bibFiles {
		ls(cacheDir, bibFile)
	}
}

func ls(cacheDir string, bibFile string) {
	if debug {
		fmt.Println("ls " + bibFile)
	}
	cache.ReadAndDo(cacheDir, bibFile, "fzf", format.EntryToFZF, printLine)
}

func printLine(s string) {
	fmt.Println(s)
}
