package main

import (
	"flag"
	"fmt"
	"github.com/msprev/fzf-bibtex/startup"
	"os"
	"strings"
)

const usage string = `bibtex-cite [-mode=pandoc|latex]
	Pretty print citations in LaTeX \cite command or pandoc @ format for selected entries passed over stdin.
`

const debug = false

func main() {
	citeMode := readArgs(usage)

	keys := startup.ReadKeysFromStdin()
	if len(keys) == 0 {
		os.Exit(-1)
	}

	if citeMode == "pandoc" {
		fmt.Print(strings.Join(keys, ", "))
	} else if citeMode == "latex" {
		for k, v := range keys {
			keys[k] = v[1:]
		}
		fmt.Print("\\cite{" + strings.Join(keys, ", ") + "}")
	}
}

func readArgs(usage string) (string) {
	citeMode := os.Getenv("FZF_BIBTEX_MODE")

	citeModeFlag := flag.String("mode", "pandoc", "citation mode")
	flag.Usage = func() {
		fmt.Println(usage)
		flag.PrintDefaults()
		fmt.Println("\nSet default arguments with environment variables:")
		fmt.Println("")
		fmt.Println("  $FZF_BIBTEX_MODE \t cite style, latex or pandoc")
	}
	flag.Parse()

	if *citeModeFlag == "pandoc" || *citeModeFlag == "latex" {
		citeMode = *citeModeFlag
	} else {
		panic("invalid cite mode: " + *citeModeFlag)
	}

	return citeMode
}
