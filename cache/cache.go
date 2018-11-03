package cache

import (
	"bibby/bibtex"
	"bufio"
	"encoding/base32"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

const debug = false

func IsFresh(cacheDir string, bibFile string) bool {
	cacheFile := cacheName(bibFile)
	if debug {
		fmt.Println(cacheFile)
	}
	// wait while locked
	for islocked(cacheDir, cacheFile) {
		if debug {
			fmt.Println("waiting...")
		}
		time.Sleep(100 * time.Millisecond)
	}
	// lock!
	lock(cacheDir, cacheFile)
	defer unlock(cacheDir, cacheFile)
	// read .timestamp file
	toRead := filepath.Join(cacheDir, cacheFile+".timestamp")
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

func Refresh(cacheDir string, bibFile string) {
	RefreshAndDo(cacheDir, bibFile, "", nil)
}

func RefreshAndDo(cacheDir string, bibFile string, subcache string, doSomething func(string)) {
	// wait while locked
	cacheFile := cacheName(bibFile)
	for islocked(cacheDir, cacheFile) {
		if debug {
			fmt.Println("waiting...")
		}
		time.Sleep(100 * time.Millisecond)
	}
	lock(cacheDir, cacheFile)
	defer unlock(cacheDir, cacheFile)
	data := make(map[string]string)
	bibtex.Parse(&data, bibFile, subcache, doSomething)
	var wg sync.WaitGroup
	for k, _ := range data {
		wg.Add(1)
		go write(filepath.Join(cacheDir, cacheFile), k, &data, &wg)
	}
	wg.Wait()
	// update timestamp
	timestamp := time.Now().UnixNano()
	f, err := os.Create(filepath.Join(cacheDir, cacheFile) + ".timestamp")
	check(err)
	defer f.Close()
	f.WriteString(strconv.FormatInt(timestamp, 10))
}

func ReadAndDo(cacheDir string, bibFile string, subcache string, doSomething func(string)) {
	if !IsFresh(cacheDir, bibFile) {
		RefreshAndDo(cacheDir, bibFile, subcache, doSomething)
		return
	}
	cacheFile := cacheName(bibFile)
	// wait while locked
	for islocked(cacheDir, cacheFile) {
		if debug {
			fmt.Println("waiting...")
		}
		time.Sleep(100 * time.Millisecond)
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

func write(cacheFile string, key string, data *map[string]string, wg *sync.WaitGroup) {
	if debug {
		fmt.Println("writing " + cacheFile + "." + key)
		fmt.Println((*data)[key])
	}
	f, err := os.Create(cacheFile + "." + key)
	check(err)
	defer f.Close()
	f.WriteString((*data)[key])
	wg.Done()
}
