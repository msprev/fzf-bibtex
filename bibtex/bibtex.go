package bibtex

import (
	// "fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func Parse(output *string, bibFiles []string, formatter func(map[string]string) string, doSomething func(string)) {
	bibtexStr := *bibtool(bibFiles)  // read data from .bibfile as string
	bibtexStr = *cleanup(&bibtexStr) // clean up the string from LaTeX crap
	sl := strings.Split(bibtexStr, "\n@")[1:]
	for _, e := range sl {
		s := formatter(parseEntry(strings.TrimSpace(e)))
		doSomething(s)
		*output += s + "\n"
	}
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

func bibtool(bibFiles []string) *string {
	const rsc string = `expand.macros = On
expand.crossref = On
preserve.keys = On
preserve.key.case = On
print.line.length { 10000 }
keep.field { date }
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
rename.field { year = date if year = ".+" }
add.field { year = "%-4.1d(date)" }
new.entry.type{Article}
new.entry.type{Book}
new.entry.type{MVBook}
new.entry.type{InBook}
new.entry.type{BookInBook}
new.entry.type{SuppBook}
new.entry.type{Booklet}
new.entry.type{Collection}
new.entry.type{MVCollection}
new.entry.type{InCollection}
new.entry.type{SuppCollection}
new.entry.type{Manual}
new.entry.type{Misc}
new.entry.type{Online}
new.entry.type{Patent}
new.entry.type{Periodical}
new.entry.type{SuppPeriodical}
new.entry.type{Proceedings}
new.entry.type{MVProceedings}
new.entry.type{Reference}
new.entry.type{MVReference}
new.entry.type{Inreference}
new.entry.type{Report}
new.entry.type{Set}
new.entry.type{Thesis}
new.entry.type{Unpublished}
new.entry.type{Cdata}
new.entry.type{CustomA}
new.entry.type{CustomB}
new.entry.type{CustomC}
new.entry.type{CustomD}
new.entry.type{CustomE}
new.entry.type{CustomF}
new.entry.type{Conference}
new.entry.type{Electronic}
new.entry.type{MasterThesis}
new.entry.type{PhdThesis}
new.entry.type{TechReport}
new.entry.type{WWW}
new.entry.type{Artwork}
new.entry.type{Audio}
new.entry.type{BibNote}
new.entry.type{Commentary}
new.entry.type{Image}
new.entry.type{Jurisdiction}
new.entry.type{Legislation}
new.entry.type{Legal}
new.entry.type{Letter}
new.entry.type{Movie}
new.entry.type{Music}
new.entry.type{Performance}
new.entry.type{Review}
new.entry.type{Software}
new.entry.type{Standard}
new.entry.type{Video}
new.entry.type{XData}
`
	rscFile, err := ioutil.TempFile(os.TempDir(), "fzf-bibtex.*.rsc")
	check(err)
	defer os.Remove(rscFile.Name())

	_, err = rscFile.Write([]byte(rsc))
	check(err)

	err = rscFile.Close()
	check(err)

	args := append([]string{"-r", rscFile.Name()}, bibFiles...)
	extCmd := exec.Command("bibtool", args...)
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
