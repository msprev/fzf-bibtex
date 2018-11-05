package startup

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const debug = false

// ReadArgs returns cacheDir, bibFiles read from environment variables and commandline
// - returns temp dir from OS, if no cacheDir specified
// - commandline argument supercedes value set by environment variables
// - exits with error, if no bibFiles specified
func ReadArgs(usage string) (string, []string) {
	// read environment variables
	cacheDir := os.Getenv("FZF_BIBTEX_CACHEDIR")
	bibFiles := make([]string, 0)
	if os.Getenv("FZF_BIBTEX_SOURCES") != "" {
		bibFiles = strings.Split(os.Getenv("FZF_BIBTEX_SOURCES"), ":")
	}
	// read commandline arguments
	wordPtr := flag.String("cache", "", "path to cache directory for list")
	flag.Usage = func() {
		fmt.Println(usage)
		flag.PrintDefaults()
		fmt.Println("\nSet default arguments with environment variables:")
		fmt.Println("")
		fmt.Println("  $FZF_BIBTEX_CACHEDIR \t path to cache directory")
		fmt.Println("  $FZF_BIBTEX_SOURCES \t paths to .bib files, separated by ':'")
	}
	flag.Parse()
	// commandline arguments supercede environment variables
	if *wordPtr != "" {
		cacheDir = *wordPtr
	}
	if len(flag.Args()) > 0 {
		bibFiles = flag.Args()
	}
	// no bib files? exit with error
	if len(bibFiles) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	// no cache dir? request a temp dir from OS
	if cacheDir == "" {
		cacheDir = os.TempDir()
	}
	// expand bibfiles to their absolute path, remove duplicates & non-existents
	cleanedBfs := make([]string, 0)
	for _, bf := range bibFiles {
		abspath, _ := filepath.Abs(bf)
		// bib file a duplicate? ignore it
		if stringInSlice(abspath, cleanedBfs) {
			continue
		}
		// bib file does not exist? ignore it
		_, err := os.Stat(abspath)
		if err != nil {
			fmt.Println("Error: file not found " + abspath)
			continue
		}
		cleanedBfs = append(cleanedBfs, abspath)
	}
	return cacheDir, cleanedBfs
}

// helper function for ReadArgs
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func ReadKeysFromStdin() []string {
	var keys []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()
		sl := strings.Split(s, "@")
		k := "@" + stripansi(sl[len(sl)-1])
		keys = append(keys, k)
	}
	if scanner.Err() != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", scanner.Err())
	}
	if debug {
		for _, k := range keys {
			fmt.Println(k, len(k), "characters long")
		}
	}
	return keys
}

// helper function for readKeysFromStdin
// (code taken from https://github.com/acarl005/stripansi)
const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var re = regexp.MustCompile(ansi)

func stripansi(str string) string {
	return re.ReplaceAllString(str, "")
}
