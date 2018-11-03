package bibtex

import (
	// "fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func Parse(data *map[string]string, bibFile string, subcache string, doSomething func(string)) {
	(*data)["markdown"] = string("markdown -- world there")
	bibtexStr := *bibtool(bibFile)   // read data from .bibfile as string
	bibtexStr = *cleanup(&bibtexStr) // clean up the string from LaTeX crap
	// now parse that string into fields
	sl := strings.Split(bibtexStr, "@")[1:]
	entries := make([]map[string]string, len(sl))
	fzf := ""
	for i, e := range sl {
		x := strings.TrimSpace(e)
		entries[i] = parseEntry(x)
		s := entryToFZF(entries[i])
		doSomething(s)
		fzf += s + "\n"
	}
	(*data)["ls"] = fzf
}

func parseEntry(entry string) map[string]string {
	m := make(map[string]string)
	lines := strings.Split(entry, "\n")
	// read key and type
	firstLine := lines[0]
	sl := strings.Fields(firstLine)
	m["type"] = strings.ToLower(sl[0])
	m["key"] = sl[1][:len(sl[1])-1] // remove last character ','
	// read other fields
	for _, l := range lines[1:] {
		sl := strings.Split(l, "=")
		k, v := sl[0], sl[1]
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		if k == "author" || k == "editor" {
			v = abbrevAuthors(v)
		}
		m[k] = v
	}
	return m
}

func abbrevAuthors(authors string) string {
	sl := strings.Split(authors, " and ")
	if len(sl) == 1 {
		return authors
	}
	if len(sl) == 2 {
		return sl[0] + " & " + sl[1]
	}
	last := len(sl) - 1
	return strings.Join(sl[0:last-1], ", ") + " & " + sl[last]
}

func entryToFZF(entry map[string]string) string {
	s := ""
	switch entry["type"] {
	case "article":
		s += entry["author"]
		s += " "
		s += "(" + entry["year"] + ")"
		s += " "
		s += "'" + entry["title"] + "'"
		s += ", "
		s += "\033[3m"
		s += entry["journal"]
		s += "\033[0m"
		s += " "
		s += entry["volume"]
		s += ", "
		s += entry["pages"]
	case "book":
		if _, ok := entry["editor"]; ok {
			s += entry["editor"]
			if strings.Contains(entry["editor"], " & ") {
				s += " (Eds.)"
			} else {
				s += " (Ed.)"
			}
		} else {
			s += entry["author"]
		}
		s += " "
		s += "(" + entry["year"] + ")"
		s += " "
		s += "\033[3m"
		s += entry["title"]
		s += "\033[0m"
		s += ", "
		s += entry["address"]
		s += ": "
		s += entry["publisher"]
	case "incollection", "inproceedings", "inbook":
		s += entry["author"]
		s += " "
		s += "(" + entry["year"] + ")"
		s += " "
		s += "'" + entry["title"] + "'"
		s += " in "
		if _, ok := entry["editor"]; ok {
			s += entry["editor"]
			if strings.Contains(entry["editor"], " & ") {
				s += " (Eds.)"
			} else {
				s += " (Ed.)"
			}
			s += " "
		}
		s += "\033[3m"
		s += entry["booktitle"]
		s += "\033[0m"
		s += ", "
		s += entry["address"]
		s += ": "
		s += entry["publisher"]
		s += ", pp. "
		s += entry["pages"]
	case "unpublished":
		s += entry["author"]
		s += " "
		s += "(" + entry["year"] + ")"
		s += " "
		s += "'" + entry["title"] + "'"
		s += ", unpublished manuscript"
	case "phdthesis", "mastersthesis":
		s += entry["author"]
		s += " "
		s += "(" + entry["year"] + ")"
		s += " "
		s += "'" + entry["title"] + "'"
		s += ", " + entry["school"]
		s += entry["author"]
		s += " "
		s += "(" + entry["year"] + ")"
		s += " "
		s += "'" + entry["title"] + "'"
		s += ", " + entry["school"]
	default:
		if _, ok := entry["editor"]; ok {
			s += entry["editor"]
			if strings.Contains(entry["editor"], " & ") {
				s += " (Eds.)"
			} else {
				s += " (Ed.)"
			}
		} else {
			s += entry["author"]
		}
		s += " "
		s += "(" + entry["year"] + ")"
		s += " "
		s += "'" + entry["title"] + "'"
	}
	// add type and bibtex key
	s += " "
	s += "\033[34m"
	s += "@" + entry["key"]
	s += "\033[0m"
	s += " "
	s += "\033[31m"
	s += "[" + entry["type"] + "]"
	s += "\033[0m"
	return s
}

func bibtool(bibFile string) *string {
	const rsc string = `preserve.keys = On
preserve.key.case = On
print.line.length { 1000 }
keep.field { author }
keep.field { title }
keep.field { year }
keep.field { journal }
keep.field { booktitle }
keep.field { editor }
keep.field { publisher }
keep.field { address }
keep.field { pages }
keep.field { school }
keep.field { volume }
`
	// create rsc file
	tmpfile, err := ioutil.TempFile(os.TempDir(), "bibby.*.rsc")
	check(err)
	defer os.Remove(tmpfile.Name()) // clean up
	// fmt.Println(tmpfile.Name())
	content := []byte(rsc)
	_, err = tmpfile.Write(content)
	check(err)
	err = tmpfile.Close()
	check(err)
	// run bibtool
	extCmd := exec.Command("bibtool", "-r", tmpfile.Name(), bibFile)
	extOut, _ := extCmd.StdoutPipe()
	err = extCmd.Start()
	check(err) // should handle this one better!
	extBytes, _ := ioutil.ReadAll(extOut)
	extCmd.Wait()
	bibtex := string(extBytes)
	return &bibtex
}

func cleanup(bibtex *string) *string {
	r := strings.NewReplacer(
		"{\\`a}", "á",
		"{\\'a}", "à",
		"{\\^a}", "â",
		"{\\\"a}", "ä",
		"{\\c{c}}", "ç",
		"{\\`e}", "é",
		"{\\'e}", "è",
		"{\\^e}", "ê",
		"{\\\"e}", "ë",
		"{\\`i}", "í",
		"{\\'i}", "ì",
		"{\\^i}", "î",
		"{\\\"i}", "ï",
		"{\\~n}", "ñ",
		"{\\`o}", "ó",
		"{\\'o}", "ò",
		"{\\^o}", "ô",
		"{\\\"o}", "ö",
		"{\\`u}", "ú",
		"{\\'u}", "ù",
		"{\\^u}", "û",
		"{\\\"u}", "ü",
		"{\\\"y}", "ÿ",
		"{\\ss}", "ß",
		"{\\`A}", "Á",
		"{\\'A}", "À",
		"{\\^A}", "Â",
		"{\\\"A}", "Ä",
		"{\\c{C}}", "Ç",
		"{\\`E}", "É",
		"{\\'E}", "È",
		"{\\^E}", "Ê",
		"{\\\"E}", "Ë",
		"{\\`I}", "Í",
		"{\\'I}", "Ì",
		"{\\^I}", "Î",
		"{\\\"I}", "Ï",
		"{\\~N}", "Ñ",
		"{\\`O}", "Ó",
		"{\\'O}", "Ò",
		"{\\^O}", "Ô",
		"{\\\"O}", "Ö",
		"{\\`U}", "Ú",
		"{\\'U}", "Ù",
		"{\\^U}", "Û",
		"{\\\"U}", "Ü",
		"{\\\"Y}", "Ÿ",
		"\\o", "ø",
		"\\ldots\\", "...",
		"\\ldots", "...",
		"\\dots\\", "...",
		"\\dots", "...",
		"~", " ",
		"``", "\"",
		"''", "\"",
		"`", "'",
		"\\&", "&",
		"$\\lambda$", "λ",
		"\\emph{", "",
		"{", "",
		"},", "",
		"}", "",
		"\\", "")
	clean := r.Replace(*bibtex)
	return &clean
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
