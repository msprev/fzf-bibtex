package startup

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ReadArgs returns cacheDir, bibFiles read from environment variables and commandline
// - returns temp dir from OS, if no cacheDir specified
// - commandline argument supercedes value set by environment variables
// - exits with error, if no bibFiles specified
func ReadArgs(usage string) (string, []string) {
	// read environment variables
	cacheDir := os.Getenv("BIBBY_CACHEDIR")
	bibFiles := make([]string, 0)
	if os.Getenv("BIBBY_SOURCES") != "" {
		bibFiles = strings.Split(os.Getenv("BIBBY_SOURCES"), ":")
	}
	// read commandline arguments
	wordPtr := flag.String("cache", "", "path to cache directory for list")
	flag.Usage = func() {
		fmt.Println(usage)
		flag.PrintDefaults()
		fmt.Println("\nSet default arguments with environment variables:")
		fmt.Println("")
		fmt.Println("  $BIBBY_CACHEDIR \t path to cache directory")
		fmt.Println("  $BIBBY_SOURCES \t paths to .bib files, separated by ':'")
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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
