package format

import (
    "strings"
)

func EntryToFZF(entry map[string]string) string {
    if entry["year"] == "" {
        entry["year"] = "no year"
    }
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
        if entry["journal"] != "" {
            s += entry["journal"]
        } else if entry["journaltitle"] != "" {
            s += entry["journaltitle"]
        }
        s += "\033[0m"
        if entry["volume"] != "" {
            s += " "
            s += entry["volume"]
        }
        if entry["pages"] != "" {
            s += ", "
            s += entry["pages"]
        }
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
        if entry["address"] != ""  {
            s += ", "
            s += entry["address"]
            s += ": "
            s += entry["publisher"]
        } else if entry["location"] != ""  {
            s += ", "
            s += entry["location"]
            s += ": "
            s += entry["publisher"]
        }
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
        if entry["address"] != ""  {
            s += ", "
            s += entry["address"]
            s += ": "
            s += entry["publisher"]
        } else if entry["location"] != ""  {
            s += ", "
            s += entry["location"]
            s += ": "
            s += entry["publisher"]
        }
        if entry["pages"] != "" {
            s += ", pp. "
            s += entry["pages"]
        }
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
    case "online":
        s += entry["author"]
        s += " "
        s += "(" + entry["year"] + ")"
        s += " "
        s += "'" + entry["title"] + "'"
        s += ", " + entry["url"]
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
    s += "\033[32m"
    s += "[" + entry["type"] + "]"
    s += "\033[0m"
    s += " "
    s += "\033[35m"
    s += "@" + entry["key"]
    s += "\033[0m"
    return s
}

func EntryToMarkdown(entry map[string]string) string {
    if entry["year"] == "" {
        entry["year"] = "no year"
    }
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
        if entry["journal"] != "" {
            s += entry["journal"]
        } else if entry["journaltitle"] != "" {
            s += entry["journaltitle"]
        }
        s += "*"
        if entry["volume"] != "" {
            s += " "
            s += entry["volume"]
        }
        if entry["pages"] != "" {
            s += ", "
            s += entry["pages"]
        }
        if entry["doi"] != "" {
            s += ". "
            s += "<https://doi.org/" + entry["doi"] + ">"
        }
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
        if entry["address"] != ""  {
            s += ", "
            s += entry["address"]
            s += ": "
            s += entry["publisher"]
        } else if entry["location"] != ""  {
            s += ", "
            s += entry["location"]
            s += ": "
            s += entry["publisher"]
        }
        if entry["doi"] != "" {
            s += ". "
            s += "<https://doi.org/" + entry["doi"] + ">"
        }
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
        if entry["address"] != ""  {
            s += ", "
            s += entry["address"]
            s += ": "
            s += entry["publisher"]
        } else if entry["location"] != ""  {
            s += ", "
            s += entry["location"]
            s += ": "
            s += entry["publisher"]
        }
        if entry["pages"] != "" {
            s += ", pp. "
            s += entry["pages"]
        }
        if entry["doi"] != "" {
            s += ". "
            s += "<https://doi.org/" + entry["doi"] + ">"
        }
    case "unpublished":
        s += entry["author"]
        s += " "
        s += "(" + entry["year"] + ")"
        s += " "
        s += "'" + entry["title"] + "'"
        s += ", unpublished manuscript"
        if entry["doi"] != "" {
            s += ". "
            s += "<https://doi.org/" + entry["doi"] + ">"
        }
    case "phdthesis", "mastersthesis":
        s += entry["author"]
        s += " "
        s += "(" + entry["year"] + ")"
        s += " "
        s += "'" + entry["title"] + "'"
        s += ", " + entry["school"]
        if entry["doi"] != "" {
            s += ". "
            s += "<https://doi.org/" + entry["doi"] + ">"
        }
    case "online":
        s += entry["author"]
        s += " "
        s += "(" + entry["year"] + ")"
        s += " "
        s += "'" + entry["title"] + "'"
        s += ", <" + entry["url"] + ">"
        if entry["doi"] != "" {
            s += ". "
            s += "<https://doi.org/" + entry["doi"] + ">"
        }
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
        if entry["doi"] != "" {
            s += ". "
            s += "<https://doi.org/" + entry["doi"] + ">"
        }
    }
    return s
}

