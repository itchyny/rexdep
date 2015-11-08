# rexdep [![Travis Build Status](https://travis-ci.org/itchyny/rexdep.svg?branch=master)](https://travis-ci.org/itchyny/rexdep) [![Latest Version](https://img.shields.io/github/release/itchyny/rexdep.svg)](https://github.com/itchyny/rexdep/releases)
### Roughly extract dependency from source code
The rexdep command is a tool for extracting dependency relation from a set of source codes.
The idea of rexdep is very simple.
For example, consider a situation that `test1.c` includes `test2.c` and `test3.c`.
```c
#include "test2.c"
#include "test3.c"

int main...
```
The file `test1.c` depends on `test2.c` and `test3.c`.
This relation can be easily extracted by matching a simple regular expression.
```
 $ grep '^#include ".*"' test1.c | sed 's/^#include *"\(.*\)"/\1/'
test2.c
test3.c
```
This simple way can be used for many languages.
For example, `import` keyword is used in Python and Haskell, `require` is used in Ruby.

The rexdep command enables to specify the `pattern` to extract the module dependency.
For the above example, you can also use rexdep to extract the dependency.
```
 $ rexdep --pattern '^\s*#include\s*"(\S+)"' test1.c
"test1.c" -> "test2.c";
"test1.c" -> "test3.c";
```
You can of course specify multiple files, and you can even specify directories and rexdep recursively investigate the source files under the subdirectories.
You may use `^\s*#include\s*[<"](\S+)[>"]` for C language or `^\s*import +(?:qualified +)?(\S+)` for Haskell language.
Allowing the user to specify by regular expression, it can be used for various languages.

There are some tools targeting on specific languages.
They investigate the source code at the level of abstract syntax tree and therefore are much powerful.
The rexdep command, on the other hand, simply checks the source code by a regular expression given by the user.
It may not as powerful as AST-level tools, but the idea of rexdep is very simple and can be used for many languages.

## Examples
### [Git](https://github.com/git/git)
```sh
 $ git clone --depth 1 https://github.com/git/git
 $ rexdep --pattern '^\s*#include\s*[<"](\S+)[>"]' --digraph git ./git/*.h | dot -Tpng -o git.png
```
[![git](https://raw.githubusercontent.com/wiki/itchyny/rexdep/image/git-1.png)](https://raw.githubusercontent.com/wiki/itchyny/rexdep/image/git.png)
The source code of Git is large and the above example checks only header files.

### [Vim](https://github.com/vim/vim)
```sh
 $ git clone --depth 1 https://github.com/vim/vim
 $ rexdep --pattern '^\s*#include\s*[<"](\S+)[>"]' --digraph vim ./vim/src/*.{c,h} | dot -Tpng -o vim.png
```
[![vim](https://raw.githubusercontent.com/wiki/itchyny/rexdep/image/vim-1.png)](https://raw.githubusercontent.com/wiki/itchyny/rexdep/image/vim.png)
The above image is a part of the output. We notice that the structure is very flat and may files include `vim.h`.

### [consul](https://github.com/hashicorp/consul)
```sh
 $ git clone --depth 1 https://github.com/hashicorp/consul
 $ rexdep --pattern '\s*"(?:\S+/)+(\S+)"' --start '^import \($' --end '^\)$' --digraph go --trimext $(find ./consul/ -name '*.go' | grep -v '_test') | dot -Tpng -o consul.png
```
[![consul](https://raw.githubusercontent.com/wiki/itchyny/rexdep/image/consul-1.png)](https://raw.githubusercontent.com/wiki/itchyny/rexdep/image/consul.png)
In order to extract imports from source codes written in Go, you can use `--start` and `--end`. The `rexdep` command extract imports between the lines matched by the start and end arguments.

### [pandoc](https://github.com/jgm/pandoc)
```sh
 $ git clone --depth 1 https://github.com/jgm/pandoc
 $ rexdep --pattern '^\s*import +(?:qualified +)?(\S+(?:Pandoc)\S+)' --module '^module +(\S+(?:Pandoc)\S+)' --digraph pandoc --recursive ./pandoc/src/ | dot -Tpng -o pandoc.png
```
[![pandoc](https://raw.githubusercontent.com/wiki/itchyny/rexdep/image/pandoc-1.png)](https://raw.githubusercontent.com/wiki/itchyny/rexdep/image/pandoc.png)
Since the arguments are regular expressions, you can flexibly limit the modules by specific words.

## Installation
### Download binary from GitHub Releases
[Releases ・ itchyny/rexdep - GitHub](https://github.com/itchyny/rexdep/releases)

### Build from source
```bash
 $ go get github.com/itchyny/rexdep
 $ go install github.com/itchyny/rexdep
```

## Usage
To be documented.

```sh
 $ rexdep --pattern 'import (\S+)' FILES
```

## Author
itchyny (https://github.com/itchyny)

## License
This software is released under the MIT License, see LICENSE.
