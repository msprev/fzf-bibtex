![](https://d.pr/i/8uXzLx+ "screenshot")

# fzf-bibtex

A BibTeX source for fzf. Example use:

To select items using fzf from a `.bib` file (as in image above):

``` {.bash}
bibtex-ls references.bib | fzf --multi --ansi
```

To select, then cite items (using pandoc '@' citation format) from a `.bib` file:

``` {.bash}
bibtex-ls references.bib | fzf --multi --ansi | bibtex-cite
```

To select, and then pretty print items (in markdown) from a `.bib` file:

``` {.bash}
bibtex-ls references.bib | fzf --multi --ansi | bibtex-markdown references.bib
```

## Installation

### Requirements

- [fzf](https://github.com/junegunn/fzf)
- [bibtool](https://ctan.org/pkg/bibtool)
- [go](https://golang.org/)

On the Mac, these can be installed by [homebrew](https://brew.sh/):

``` {.bash}
brew install fzf
brew install bib-tool
brew install go
```

- [fzf.vim](https://github.com/junegunn/fzf.vim) (if you want vim integration)

### Installation

``` {.bash}
git clone https://github.com/msprev/fzf-bibtex
cd fzf-bibtex
go install ./...
```

### Vim integration

Add this to your `vimrc` file:

``` {.vim}
let $FZF_BIBTEX_CACHEDIR = 'PATH-TO-CACHE-DIR'
let $FZF_BIBTEX_SOURCES = 'PATH-TO-BIBTEX-FILE'

function! s:bibtex_cite_sink(lines)
    let r=system("bibby-cite ", a:lines)
    execute ':normal! i' . r
endfunction

function! s:bibtex_markdown_sink(lines)
    let r=system("bibtex-markdown ", a:lines)
    execute ':normal! i' . r
endfunction

nnoremap <leader>c :call fzf#run({
                        \ 'source': 'bibtex-ls',
                        \ 'sink*': function('<sid>bibtex_cite_sink'),
                        \ 'up': '40%',
                        \ 'options': '--ansi --layout=reverse-list --multi --prompt "Cite> "'})<CR>

nnoremap <leader>m :call fzf#run({
                        \ 'source': 'bibtex-ls',
                        \ 'sink*': function('<sid>bibtex_markdown_sink'),
                        \ 'up': '40%',
                        \ 'options': '--ansi --layout=reverse-list --multi --prompt "Markdown> "'})<CR>
```

`<leader>c` will bring up fzf to insert citation to selected items.

`<leader>m` will bring up fzf to insert pretty markdown versions of selected items.

An insert mode mapping, typing '@@' brings up fzf to insert a citation:

``` {.vim}
function! s:bibtex_cite_sink_insert(lines)
    let r=system("bibtex-cite ", a:lines)
    execute ':normal! a' . r
    call feedkeys('a', 'n')
endfunction

inoremap <silent> @@ <c-g>u<c-o>:call fzf#run({
                        \ 'source': 'bibtex-ls',
                        \ 'sink*': function('<sid>bibtex_cite_sink_insert'),
                        \ 'up': '40%',
                        \ 'options': '--ansi --layout=reverse-list --multi --prompt "Cite> "'})<CR>
```


## Command line use

``` {.bash}
bibtex-ls [-cache=...] [file1.bib file2.bib ...]
```

Lists to stdout the content of .bib files, one record per line.

If the following environment variables are set, then these command line arguments can be omitted.

- `FZF_BIBTEX_CACHEDIR`: path to a cache directory
- `FZF_BIBTEX_SOURCES`: path to bibtex file; multiple items separated by a '`:`'

The cache directory should be a suitable directory for bibtex-ls temporary files.
Parsing BibTeX databases is computationally intensive, so the command caches the results.
    The cache is updated if the underlying BibTeX file has been changed.
    If no cache directory is specified, the operating system's directory for temporary files is used.

``` {.bash}
bibtex-cite
```

Pretty print citations (in pandoc '@' format) for selected `.bib` entries passed over stdin.

``` {.bash}
bibtex-markdown [-cache=...] [file1.bib file2.bib ...]
```

Pretty print items (in markdown) for selected `.bib` entries passed over stdin.

Cache directory may be set using the same environment variable as bibtex-ls.

## Errors?

fzf-bibtex uses [bibtool](https://ctan.org/pkg/bibtool) to parse BibTeX
files.  If there is an error, it is likely that your BibTeX file is not
being parsed correctly.  You can locate the cause, and correct it, by
running bibtool directly on your BibTeX file from the command line.  Look
at any errors reported from:

``` {.bash}
bibtool references.bib -o parsed.bib
```

The BibTeX fields that fzf-bibtex asks bibtool to extract from your file
can be seen by running bibtool with this `rsc` file:

```
preserve.keys = On
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
```

## Release notes

- 1.0 (4 November 2018)
    - first version

## Similar

- [unite-bibtex](https://github.com/msprev/unite-bibtex) -- no longer maintained; this replaces it.
