package main

import (
	"flag"
	"fmt"
	"github.com/msprev/fzf-bibtex/startup"
	"os"
	"strings"
)

const usage string = `bibtex-cite [-mode=pandoc|latex] [-prefix=...] [-postfix=...] [-separator=...]
	Pretty print citations for selected entries passed over stdin.

    Citation format may be customised with -prefix, -postfix, and -separator options.
    -mode provides presets for pandoc and LaTeX citations.
`

const debug = false

func main() {
	citeMode, prefix, postfix, separator := readArgs(usage)

	if citeMode == "pandoc" {
		prefix="@"
        postfix=""
        separator="; @"
	} else if citeMode == "latex" {
		prefix="\\cite{"
        postfix="}"
        separator=", "
	}

	keys := startup.ReadKeysFromStdin()
	if len(keys) == 0 {
		os.Exit(-1)
	}

    for k, v := range keys {
        keys[k] = v[1:]
    }
    fmt.Print(prefix + strings.Join(keys, separator) + postfix)
}

func readArgs(usage string) (string, string, string, string) {
	citeMode := os.Getenv("FZF_BIBTEX_MODE")

	citeModeFlag := flag.String("mode", "", ` -mode=pandoc => -prefix="@"      -postfix=""  -separator="; @"
 -mode=latex  => -prefix="\cite{" -postfix="}" -separator=", "
`)
	prefix := flag.String("prefix", "@", "string before citation key")
	postfix := flag.String("postfix", "", "string after citation key")
	separator := flag.String("separator", "; @", "string separating citations keys")
	flag.Usage = func() {
		fmt.Println(usage)
		flag.PrintDefaults()
		fmt.Println("\nSet default arguments with environment variables:")
		fmt.Println("")
		fmt.Println("  $FZF_BIBTEX_MODE \t cite style: 'latex' or 'pandoc'")
	}
	flag.Parse()

	if *citeModeFlag == "pandoc" || *citeModeFlag == "latex" {
		citeMode = *citeModeFlag
	} else {
		citeMode = ""
	}

	return citeMode, *prefix, *postfix, *separator
}
