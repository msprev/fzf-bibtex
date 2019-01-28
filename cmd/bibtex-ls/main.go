package main

import (
	"fmt"
	"strings"
	"github.com/msprev/fzf-bibtex/cache"
	"github.com/msprev/fzf-bibtex/format"
	"github.com/msprev/fzf-bibtex/startup"
)

const usage string = `bibtex-ls [-cache=...] [file1.bib file2.bib ...]
  lists to stdout the content of .bib files, one record per line
`

const debug = false

func main() {
	cacheDir, bibFiles := startup.ReadArgs(usage)
	if debug {
		fmt.Println("cachedir: ", cacheDir)
		fmt.Println("bib files: ", bibFiles)
	}
	ls(cacheDir, bibFiles)
}

func ls(cacheDir string, bibFiles []string) {
	if debug {
		fmt.Println("ls " + strings.Join(bibFiles, " "))
	}
	cache.ReadAndDo(cacheDir, bibFiles, "fzf", format.EntryToFZF, printLine)
}

func printLine(s string) {
	fmt.Println(s)
}
