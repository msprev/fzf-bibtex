package format

import (
	"strings"
)

func EntryToFZF(entry map[string]string) string {
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
	s += "\033[31m"
	s += "[" + entry["type"] + "]"
	s += "\033[0m"
	s += " "
	s += "\033[34m"
	s += "@" + entry["key"]
	s += "\033[0m"
	return s
}

func EntryToMarkdown(entry map[string]string) string {
	s := ""
	s += "@" + entry["key"] + " "
	switch entry["type"] {
	case "article":
		s += entry["author"]
		s += " "
		s += "(" + entry["year"] + ")"
		s += " "
		s += "'" + entry["title"] + "'"
		s += ", "
		s += "*"
		s += entry["journal"]
		s += "*"
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
		s += "*"
		s += entry["title"]
		s += "*"
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
		s += "*"
		s += entry["booktitle"]
		s += "*"
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
	return s
}
