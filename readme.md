![](https://d.pr/i/8uXzLx+ "screenshot")

# fzf-bibtex

A BibTeX source for fzf.

- Blazingly fast, even with extremely large BibTeX files
- Caches results intelligently (hence the speed)
- Uses a well-understood framework to parse BibTeX ([bibtool](https://ctan.org/pkg/bibtool))
- vim and neovim integration (with [fzf.vim](https://github.com/junegunn/fzf.vim) or [fzf-lua](https://github.com/ibhagwan/fzf-lua))
- Supports multiple BibTeX files
- Supports cross references (thanks to [\@cao](https://github.com/cao))
- Supports multiple citation formats
- BibLaTeX support (thanks to [\@ashwinvis](https://github.com/ashwinvis))

## Example use

To select items using fzf from a `.bib` file (as in image above):

``` {.bash}
bibtex-ls references.bib | fzf --multi --ansi
```

To cite items (using the pandoc '@' format) from a `.bib` file:

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

If you want vim/neovim integration:

- [fzf.vim](https://github.com/junegunn/fzf.vim)

or, if you like lua (for neovim only):

- [fzf-lua](https://github.com/ibhagwan/fzf-lua)

NB.  You only need one of the other of these (see mappings below).  You
can install both if you really want.

### Installation

``` {.bash}
go install github.com/msprev/fzf-bibtex/cmd/bibtex-ls@latest
go install github.com/msprev/fzf-bibtex/cmd/bibtex-markdown@latest
go install github.com/msprev/fzf-bibtex/cmd/bibtex-cite@latest
```

### Why these dependencies?

Parsing BibTeX is a non-trivial task.  It is best to do it in a
well-understood and reliable way.  fzf-bibtex uses an extremely stable,
reliable, and widely used parser, `bibtool`, which is the benchmark for
parsing BibTeX.  The goal of fzf-bibtex is to have no noticable delay
when searching, even for extremely large BibTeX files.  Writing it with
Go allows the desired responsiveness to be achieved.

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

Citation format may be customised with `-prefix`, `-postfix`, and `-separator` options.

Default values (suitable for pandoc '@' format):

- `-prefix="@"` `-postfix=""` `-separator="; @"`

Legacy `-mode` option provides presets for pandoc and LaTeX style
citations.  `-mode` options:

- `-mode=pandoc` = `-prefix="@"      -postfix=""  -separator="; @"`
- `-mode=latex`  = `-prefix="\cite{" -postfix="}" -separator=", "`

### bibtex-markdown

``` {.bash}
bibtex-markdown [-cache=...] [file1.bib file2.bib ...]
```

Pretty print items (in markdown) for selected `.bib` entries passed over stdin.

Cache directory may be set using the same environment variable as bibtex-ls.

## fzf.vim integration

Assuming the executables installed above are available to Vim in your file path, add the following code to your `vimrc` file (or, for neovim, your `init.vim`):

<details><summary>fzf-vim integration (normal mode)</summary>
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
</details>

- `<leader>c` will bring up fzf to cite selected items
- `<leader>m` will bring up fzf to markdown pretty print cite selected items


<details><summary>fzf-vim integration (insert mode)</summary>
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
</details>

- `@@` will bring up fzf to cite selected items

Alternative insert mode mapping (`@@`) that detects .bib files in parent, current or child directories (thanks to [\@ashwinvis](https://github.com/ashwinvis)):

<details><summary>fzf-vim integration (alternative insert mapping -- automatically reads from nearby .bib files)</summary>
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
</details>

## fzf-lua integration

If you use [fzf-lua](https://github.com/ibhagwan/fzf-lua) in neovim, you can add the following
code inside to your `init.lua` or similar config file.

<details><summary>fzf-lua integration</summary>

``` {.lua}
-- default list of bibfiles
-- can be overriden by changing vim.b.bibfiles inside buffer
local default_bibfiles = {
    }
-- default behaviour: doing nothing, unless in right filetype
vim.keymap.set("n", "[fzf]c", "<nop>")

-- default cache directory
local cachedir = vim.fn.stdpath("state") .. "/fzf-bibtex/"

-- actions
local pandoc = function(selected, opts)
    local result = vim.fn.system('bibtex-cite', selected)
    vim.api.nvim_put({ result }, "c", false, true)
    if opts.fzf_bibtex.mode == "i" then
        vim.api.nvim_feedkeys("i", "n", true)
    end
end

local citet = function(selected, opts)
    local result = vim.fn.system('bibtex-cite -prefix="\\citet{" -postfix="}" -separator=","', selected)
    vim.api.nvim_put({ result }, "c", false, true)
    if opts.mode == "i" then
        vim.api.nvim_feedkeys("i", "n", true)
    end
end

local citep = function(selected, opts)
    local result = vim.fn.system('bibtex-cite -prefix="\\citep{" -postfix="}" -separator=","', selected)
    vim.api.nvim_put({ result }, "c", false, true)
    if opts.mode == "i" then
        vim.api.nvim_feedkeys("i", "n", true)
    end
end

local markdown_print = function(selected, opts)
    local result = vim.fn.system("bibtex-markdown -cache=" .. cachedir .. " " .. table.concat(vim.b.bibfiles, " "),
        selected)
    local result_lines = {}
    for line in result:gmatch('[^\n]+') do
        table.insert(result_lines, line)
    end
    vim.api.nvim_put(result_lines, "l", true, true)
    if opts.mode == "i" then
        vim.api.nvim_feedkeys("i", "n", true)
    end
end

local fzf_bibtex_menu = function(mode)
    return function()
        -- check cache directory hasn't mysteriously disappeared
        if vim.fn.isdirectory(cachedir) == 0 then
            vim.fn.mkdir(cachedir, "p")
        end

        require 'fzf-lua'.config.set_action_helpstr(pandoc, "@-pandoc")
        require 'fzf-lua'.config.set_action_helpstr(citet, "\\citet{}")
        require 'fzf-lua'.config.set_action_helpstr(citep, "\\citep{}")
        require 'fzf-lua'.config.set_action_helpstr(markdown_print, "markdown-pretty-print")

        -- header line: the bibtex filenames
        local filenames = {}
        for i, fullpath in ipairs(vim.b.bibfiles) do
            filenames[i] = vim.fn.fnamemodify(fullpath, ":t")
        end
        local header = table.concat(filenames, "\\ ")

        -- set default action
        local default_action = nil
        if vim.bo.ft == "markdown" then
            default_action = pandoc
        elseif
            vim.bo.ft == "tex" then
            default_action = citet
        end

        -- run fzf
        return require 'fzf-lua'.fzf_exec(
            "bibtex-ls "
            .. "-cache=" .. cachedir .. " "
            .. table.concat(vim.b.bibfiles, " "), {
                actions = {
                        ['default'] = default_action,
                        ['alt-a'] = pandoc,
                        ['alt-t'] = citet,
                        ['alt-p'] = citep,
                        ['alt-m'] = markdown_print,
                },
                fzf_bibtex = { ['mode'] = mode },
                fzf_opts = { ['--prompt'] = 'BibTeX> ',['--header'] = header }
            })
    end
end

-- Only enable mapping in tex or markdown
vim.api.nvim_create_autocmd("Filetype", {
    desc = "Set up keymaps for fzf-bibtex",
    group = vim.api.nvim_create_augroup("fzf-bibtex", { clear = true }),
    pattern = { "markdown", "tex" },
    callback = function()
        vim.b.bibfiles = default_bibfiles
        vim.keymap.set("n", "<leader>c", fzf_bibtex_menu("n"), { buffer = true, desc = "FZF: BibTeX [C]itations" })
        vim.keymap.set("i", "@@", fzf_bibtex_menu("i"), { buffer = true, desc = "FZF: BibTeX [C]itations" })
    end
})
```
</details>

- `<leader>c` will bring up fzf to cite selected items
    - `<cr>`: insert with default citation style
    - `<alt-a>`: insert citation with pandoc @ style
    - `<alt-t>`: insert citation with LaTeX \\citet{} style
    - `<alt-p>`: insert citation with LaTeX \\citep{} style
    - `<alt-m>`: pretty print selected items in markdown
- `@@` in insert mode brings up the same fzf menu.

Mappings will only be active for `tex` or `markdown` filetypes.


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
can be seen by running bibtool with the `rsc` file specified in [this string](https://github.com/msprev/fzf-bibtex/blob/ae9b939fb30448a85a6b18370bfdab4a451eeba4/bibtex/bibtex.go#L57).


## Release notes

- 1.1 (17 February 2020)
    - support arbitrary citation formats
- 1.0 (4 November 2018)
    - first version

## Similar

- [unite-bibtex](https://github.com/msprev/unite-bibtex) -- no longer maintained; this replaces it.
