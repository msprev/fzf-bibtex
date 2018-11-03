package main

import (
	"bibby/cache"
	"bibby/startup"
	"bufio"
	"fmt"
	"os"
	"sync"
)

const usage string = `bibby-ls [-cache=...] [file1.bib file2.bib ...]
  lists to stdout the content of .bib files, one record per line
`

const debug = false

func main() {
	cacheDir, bibFiles := startup.ReadArgs(usage)
	if debug {
		fmt.Println("cachedir:", cacheDir)
		fmt.Println("bib files: ", bibFiles)
	}
	var wg sync.WaitGroup
	for _, bibFile := range bibFiles {
		wg.Add(1)
		go ls(cacheDir, bibFile, &wg)
	}
	wg.Wait()
}

func ls(cacheDir string, bibFile string, wg *sync.WaitGroup) {
	if debug {
		fmt.Println("go ls " + bibFile)
	}
	cache.ReadAndDo(cacheDir, bibFile, "ls", printLine)
	wg.Done()
}

func printLine(s string) {
	fmt.Println(s)
}

func readStdin(doSomething func(string)) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		doSomething(scanner.Text()) // Println will add back the final '\n'
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
