package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

const usage string = `bibtex-cite 
	pretty print citations (in pandoc '@' format) for selected .bib entries passed over stdin.
`

const debug = false

func main() {
	readArgs(usage)
	// read keys from stdin into a slice
	keys := readKeysFromStdin()
	if len(keys) == 0 {
		os.Exit(-1)
	}
	fmt.Print(strings.Join(keys, "; "))
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
		fmt.Fprintln(os.Stderr, "reading standard input:", scanner.Err())
	}
	if debug {
		fmt.Println(keys)
	}
	return keys
}

func readArgs(usage string) {
	flag.Usage = func() {
		fmt.Println(usage)
		flag.PrintDefaults()
	}
	flag.Parse()
}
