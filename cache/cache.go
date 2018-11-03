package cache

import (
	"bibby/bibtex"
	"bufio"
	"encoding/base32"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const debug = false

func IsFresh(cacheDir string, subcache string, bibFile string) bool {
	cacheFile := cacheName(bibFile)
	if debug {
		fmt.Println(cacheFile)
	}
	// wait while locked
	for islocked(cacheDir, cacheFile) {
		if debug {
			fmt.Println("waiting...")
		}
		time.Sleep(50 * time.Millisecond)
	}
	// lock!
	lock(cacheDir, cacheFile)
	defer unlock(cacheDir, cacheFile)
	// read .timestamp file
	toRead := filepath.Join(cacheDir, cacheFile+"."+subcache+".timestamp")
	file, err := os.Open(toRead)
	if err != nil {
		if debug {
			fmt.Println("cache does not exist yet " + toRead)
		}
		return false // cache does not exist yet
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	timestamp, _ := strconv.ParseInt(string(scanner.Text()), 10, 64)
	// read mod time of bibFile
	fi, err := os.Stat(bibFile)
	check(err)
	timestamp2 := fi.ModTime().UnixNano()
	if timestamp < timestamp2 {
		if debug {
			fmt.Println("cache is out of date " + bibFile)
		}
		return false // cache is out of date
	}
	if debug {
		fmt.Println("cache is up to date " + bibFile)
	}
	return true // cache is up to date
}

func RefreshAndDo(cacheDir string, bibFile string, subcache string, formatter func(map[string]string) string, doSomething func(string)) {
	// wait while locked
	cacheFile := cacheName(bibFile)
	for islocked(cacheDir, cacheFile) {
		if debug {
			fmt.Println("waiting...")
		}
		time.Sleep(50 * time.Millisecond)
	}
	lock(cacheDir, cacheFile)
	defer unlock(cacheDir, cacheFile)
	data := ""
	bibtex.Parse(&data, bibFile, formatter, doSomething)
	write(filepath.Join(cacheDir, cacheFile+"."+subcache), &data)
	// update timestamp
	timestamp := time.Now().UnixNano()
	f, err := os.Create(filepath.Join(cacheDir, cacheFile+"."+subcache+".timestamp"))
	check(err)
	defer f.Close()
	f.WriteString(strconv.FormatInt(timestamp, 10))
}

func ReadAndDo(cacheDir string, bibFile string, subcache string, formatter func(map[string]string) string, doSomething func(string)) {
	if !IsFresh(cacheDir, subcache, bibFile) {
		RefreshAndDo(cacheDir, bibFile, subcache, formatter, doSomething)
		return
	}
	cacheFile := cacheName(bibFile)
	// wait while locked
	for islocked(cacheDir, cacheFile) {
		if debug {
			fmt.Println("waiting...")
		}
		time.Sleep(50 * time.Millisecond)
	}
	lock(cacheDir, cacheFile)
	defer unlock(cacheDir, cacheFile)
	toRead := filepath.Join(cacheDir, cacheFile+"."+subcache)
	if debug {
		fmt.Println("opening: " + toRead)
	}
	file, err := os.Open(toRead)
	check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		doSomething(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func cacheName(bibFile string) string {
	abspath, _ := filepath.Abs(bibFile)
	return base32.StdEncoding.EncodeToString([]byte(abspath))
}

func lock(cacheDir string, cacheFile string) {
	if islocked(cacheDir, cacheFile) == true {
		panic(fmt.Sprintf("Attempted to lock an already locked cache %s", cacheFile))
	}
	lockFile := filepath.Join(cacheDir, cacheFile+".lock")
	f, err := os.Create(lockFile)
	check(err)
	defer f.Close()
}

func unlock(cacheDir string, cacheFile string) {
	if islocked(cacheDir, cacheFile) == false {
		panic(fmt.Sprintf("Attempted to unlock an already unlocked cache %s", cacheFile))
	}
	lockFile := filepath.Join(cacheDir, cacheFile+".lock")
	err := os.Remove(lockFile)
	check(err)
}

func islocked(cacheDir string, cacheFile string) bool {
	lockFile := filepath.Join(cacheDir, cacheFile+".lock")
	_, err := os.Stat(lockFile)
	if err == nil {
		return true
	}
	return false
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func write(cacheFile string, data *string) {
	if debug {
		fmt.Println("writing " + cacheFile)
		fmt.Println(*data)
	}
	f, err := os.Create(cacheFile)
	check(err)
	defer f.Close()
	f.WriteString(*data)
}
