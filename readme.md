![](https://d.pr/i/8uXzLx+ "screenshot")

# fzf-bibtex

A BibTeX source for fzf.

- Blazingly fast, even with large BibTeX files
- Caches results intelligently (hence the speed)
- Uses a well-understood framework to parse BibTeX ([bibtool](https://ctan.org/pkg/bibtool))
- Vim integration (via [fzf.vim](https://github.com/junegunn/fzf.vim))
- Supports multiple BibTeX files
- Supports cross references (thanks to [\@cao](https://github.com/cao))
- Supports multiple citation formats
- Basic BibLaTeX support (thanks to [\@ashwinvis](https://github.com/ashwinvis))

## Example use

To select items using fzf from a `.bib` file (as in image above):

``` {.bash}
bibtex-ls references.bib | fzf --multi --ansi
```

To cite items (using the pandoc '@' format or LaTeX \cite command) from a `.bib` file:

``` {.bash}
bibtex-ls references.bib | fzf --multi --ansi | bibtex-cite
```

To pretty print items (in markdown) from a `.bib` file:

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

If you want vim integration:

- [fzf.vim](https://github.com/junegunn/fzf.vim)

### Installation

``` {.bash}
go get github.com/msprev/fzf-bibtex/cmd/bibtex-ls
go install github.com/msprev/fzf-bibtex/cmd/bibtex-ls
go install github.com/msprev/fzf-bibtex/cmd/bibtex-markdown
go install github.com/msprev/fzf-bibtex/cmd/bibtex-cite
```

## Command line use

### bibtex-ls

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

(NB. If you are tinkering with fzf-bibtex's codebase, beware of outdated caches.
Cache is *only* updated if the underlying BibTeX file has been changed.
If you change the fzf-bibtex codebase, make sure to flush the cache by `touch`ing the BibTeX files, or deleting the cache, before you run new code on them).

### bibtex-cite

``` {.bash}
bibtex-cite [-mode=pandoc|latex] [-prefix=...] [-postfix=...] [-separator=...]
```

Pretty print citations for selected entries passed over stdin.

Citation format may be customised with `-prefix`, `-postfix`, and `-separator` options. `-mode` option provides presets for pandoc and LaTeX citations.

Default values (suitable for pandoc citations):

- `-prefix="@"` `-postfix=""` `-separator="; @"`

`-mode` options:

- `-mode=pandoc` => `-prefix="@"      -postfix=""  -separator="; @"`
- `-mode=latex`  => `-prefix="\cite{" -postfix="}" -separator=", "`

### bibtex-markdown

``` {.bash}
bibtex-markdown [-cache=...] [file1.bib file2.bib ...]
```

Pretty print items (in markdown) for selected `.bib` entries passed over stdin.

Cache directory may be set using the same environment variable as bibtex-ls.

## Vim integration

Assuming the executables installed above are available to Vim in your file path, add this to your `vimrc` file:

``` {.vim}
let $FZF_BIBTEX_CACHEDIR = 'PATH-TO-CACHE-DIR'
let $FZF_BIBTEX_SOURCES = 'PATH-TO-BIBTEX-FILE'

function! s:bibtex_cite_sink(lines)
    let r=system("bibtex-cite ", a:lines)
    execute ':normal! a' . r
endfunction

function! s:bibtex_markdown_sink(lines)
    let r=system("bibtex-markdown ", a:lines)
    execute ':normal! a' . r
endfunction

nnoremap <silent> <leader>c :call fzf#run({
                        \ 'source': 'bibtex-ls',
                        \ 'sink*': function('<sid>bibtex_cite_sink'),
                        \ 'up': '40%',
                        \ 'options': '--ansi --layout=reverse-list --multi --prompt "Cite> "'})<CR>

nnoremap <silent> <leader>m :call fzf#run({
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

An alternative insert mode mapping that detects .bib files in parent, current or child directories (thanks to [\@ashwinvis](https://github.com/ashwinvis)):

```
function! Bibtex_ls()
  let bibfiles = (
      \ globpath('.', '*.bib', v:true, v:true) +
      \ globpath('..', '*.bib', v:true, v:true) +
      \ globpath('*/', '*.bib', v:true, v:true)
      \ )
  let bibfiles = join(bibfiles, ' ')
  let source_cmd = 'bibtex-ls '.bibfiles
  return source_cmd
endfunction

function! s:bibtex_cite_sink_insert(lines)
    let r=system("bibtex-cite ", a:lines)
    execute ':normal! a' . r
    call feedkeys('a', 'n')
endfunction

inoremap <silent> @@ <c-g>u<c-o>:call fzf#run({
                        \ 'source': Bibtex_ls(),
                        \ 'sink*': function('<sid>bibtex_cite_sink_insert'),
                        \ 'up': '40%',
                        \ 'options': '--ansi --layout=reverse-list --multi --prompt "Cite> "'})<CR>
```


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
expand.macros = On
expand.crossref = On
preserve.keys = On
preserve.key.case = On
print.line.length { 1000 }
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
```

## Release notes

- 1.1 (17 February 2020)
    - support arbitrary citation formats
- 1.0 (4 November 2018)
    - first version

## Similar

- [unite-bibtex](https://github.com/msprev/unite-bibtex) -- no longer maintained; this replaces it.
