# rexdep [![Travis Build Status](https://travis-ci.org/itchyny/rexdep.svg?branch=master)](https://travis-ci.org/itchyny/rexdep)
### Roughly extract dependency from source code
The rexdep command is a tool for extracting dependency relation from software.

[![rexdep](https://raw.githubusercontent.com/wiki/itchyny/rexdep/image/rexdep.png)](https://raw.githubusercontent.com/wiki/itchyny/rexdep/image/rexdep.png)

When we see a large project, we want to understand the code structure intuitively.
We sometimes join a project at work, the project has lots of source code.
We sometimes want to make pull requests to famous OSS software, but the software is too large.
How can we understand the file structure of a huge project without reading through the software?

It is a good solution to check the module dependency relations among the software.
A module in a software depends on other modules in that software.
Some module depends on many other modules or another module does not depend on others.
Extracting the module dependency enables us to understand the relationship among the modules.
We can use the dependency graph to read the software top down or bottom up.
We sometimes find the core modules in the project because such modules are depended by many other modules.

So, how can we extract dependency relationship from a code base of software?

The idea of rexdep is very simple; in many cases, we can extract the module names by a regular expression.
Let me explain by a simple example in C language, where we want to extract the dependency relations between the files, rather than modules.
Consider a situation that `test1.c` includes `test2.c` and `test3.c`.
```c
/* This is test1.c */
#include "test2.c"
#include "test3.c"

int main...
```
The file `test1.c` depends on `test2.c` and `test3.c`.
This relation can be easily extracted using a simple regular expression.
```bash
 $ grep '^#include ".*"' test1.c | sed 's/^#include *"\(.*\)"/\1/'
test2.c
test3.c
```
This way can be applied for many languages.
For example, `import` keyword is used in Python and Haskell and `require` in Ruby.

The rexdep command enables us to specify the `pattern`, the regular expression to extract the module dependency from source codes.
For the above example, we can use rexdep to extract the dependency between the C source codes.
```bash
 $ rexdep --pattern '^\s*#include\s*"(\S+)"' test1.c
test1.c test2.c
test1.c test3.c
```
The captured string is regarded as the file names included by the source code.
We can of course specify multiple files.
We can also specify directories and rexdep recursively investigate the source files under the subdirectories.
Allowing the user to specify by regular expression, it can be used for various languages; `^\s*#include\s*[<"](\S+)[>"]` for C language or `^\s*import +(?:qualified +)?([[:alnum:].]+)` for Haskell language.

There are some other tools targeting on specific languages.
They investigate source codes at the level of abstract syntax tree and therefore much powerful.
The rexdep command, on the other hand, simply checks the source code by a regular expression given by the user.
It may not as powerful as AST-level tools, but the idea is very simple and can be used for many languages.
Now you understand what rexdep stands for; *roughly extract dependency*.

## Installation
### Homebrew
```bash
 $ brew tap itchyny/rexdep
 $ brew install rexdep
```

### Download binary from GitHub Releases
[Releases・itchyny/rexdep - GitHub](https://github.com/itchyny/rexdep/releases)

### Build from source
```bash
 $ go get github.com/itchyny/rexdep
 $ go install github.com/itchyny/rexdep
```

## Usage
### Basic usage: --pattern, --format
Consider the following sample file.
```bash
 $ cat test1
import test2
import test3
```
We want to extract the dependency relation from this file `test1`.
Specify `--pattern` the regular expression to extract the module names it imports.
```bash
 $ rexdep --pattern 'import +(\S+)' test1
test1 test2
test1 test3
```
The output result shows that `test1` depends on `test2` and `test3`.
Each line contains the space separated module names.
The captured strings in the `--pattern` argument are interpreted as the module names imported by each file.
The regular expression is compiled to `Regexp` type of Go language, so refer to the [document](https://golang.org/s/re2syntax) for the regexp syntax or try `go doc regexp/syntax`.

We can use the rexdep command to output in the dot language
```bash
 $ rexdep --pattern 'import +(\S+)' --format dot test1
digraph "graph" {
  "test1" -> "test2";
  "test1" -> "test3";
}
```
and it works perfectly with graphviz (dot command).
```bash
 $ rexdep --pattern 'import +(\S+)' --format dot test1 | dot -Tpng -o test.png
```
![example](https://raw.githubusercontent.com/wiki/itchyny/rexdep/image/example-1.png)

Also, it works for multiple files.
```bash
 $ rexdep --pattern 'import +(\S+)' --format dot test{1,2,3,4}
digraph "graph" {
  "test1" -> "test2";
  "test1" -> "test3";
  "test2" -> "test4";
  "test3" -> "test4";
}
 $ rexdep --pattern 'import +(\S+)' --format dot test{1,2,3,4} | dot -Tpng -o test.png
```
![example](https://raw.githubusercontent.com/wiki/itchyny/rexdep/image/example-2.png)

This is a basic example of rexdep.

You can also change the output format to JSON,
```bash
 $ rexdep --pattern 'import +(\S+)' --format json test{1,2,3,4}
{
  "test1": [
    "test2",
    "test3"
  ],
  "test2": [
    "test4"
  ],
  "test3": [
    "test4"
  ]
}
```
which may be piped to [jq](https://stedolan.github.io/jq/) command.

## Examples
### [Git](https://github.com/git/git)
```sh
 $ git clone --depth 1 https://github.com/git/git
 $ rexdep --pattern '^\s*#include\s*[<"](\S+)[>"]' --format dot ./git/*.h | dot -Tpng -o git.png
```
[![git](https://raw.githubusercontent.com/wiki/itchyny/rexdep/image/git-1.png)](https://raw.githubusercontent.com/wiki/itchyny/rexdep/image/git.png)
The code base of Git is large and the above example checks only header files. The above image (click to see the full image) reveals the relationship between the header files.

### [Vim](https://github.com/vim/vim)
```sh
 $ git clone --depth 1 https://github.com/vim/vim
 $ rexdep --pattern '^\s*#include\s*[<"](\S+)[>"]' --format dot ./vim/src/*.{c,h} | dot -Tpng -o vim.png
```
[![vim](https://raw.githubusercontent.com/wiki/itchyny/rexdep/image/vim-1.png)](https://raw.githubusercontent.com/wiki/itchyny/rexdep/image/vim.png)
We notice that the structure is flat and many files include `vim.h`.

### [consul](https://github.com/hashicorp/consul)
```sh
 $ git clone --depth 1 https://github.com/hashicorp/consul
 $ rexdep --pattern '"github.com/(?:hashicorp/consul/(?:\S+/)*)?(\S+)"' --start '^import +["(]' --end '^\)$|^import +"' --format dot --trimext $(find ./consul/ -name '*.go' | grep -v '_test') | dot -Tpng -o consul.png
```
[![consul](https://raw.githubusercontent.com/wiki/itchyny/rexdep/image/consul-1.png)](https://raw.githubusercontent.com/wiki/itchyny/rexdep/image/consul.png)
It is difficult to extract dependency relation from source codes written in Go because we can use functions from the other codes at the same directory without writing import. We can skip the directories by `(?:\S/)*`. The rexdep command extract imports between the lines matched by the `start` and `end` arguments.

### [pandoc](https://github.com/jgm/pandoc)
```sh
 $ git clone --depth 1 https://github.com/jgm/pandoc
 $ rexdep --pattern '^\s*import +(?:qualified +)?([[:alnum:].]+Pandoc[[:alnum:].]*)' --module '^module +([[:alnum:].]+Pandoc[[:alnum:].]*)' --format dot --recursive ./pandoc/src/ | dot -Tpng -o pandoc.png
```
[![pandoc](https://raw.githubusercontent.com/wiki/itchyny/rexdep/image/pandoc-1.png)](https://raw.githubusercontent.com/wiki/itchyny/rexdep/image/pandoc.png)
We can flexibly limit the modules by specific words.

### [lens](https://github.com/ekmett/lens)
```sh
 $ git clone --depth 1 https://github.com/ekmett/lens
 $ rexdep --pattern '^\s*import +(?:qualified +)?(\S+Lens\S*)' --module '^module +(\S+Lens\S*)' --format dot --recursive ./lens/src/ | dot -Tpng -o lens.png
```
[![lens](https://raw.githubusercontent.com/wiki/itchyny/rexdep/image/lens-1.png)](https://raw.githubusercontent.com/wiki/itchyny/rexdep/image/lens.png)
It is very fun to see the dependency graph of cool libraries like lens, isn't it?

## Bug Tracker
Report bug at [Issues・itchyny/rexdep - GitHub](https://github.com/itchyny/rexdep/issues).

## Author
itchyny (https://github.com/itchyny)

## License
This software is released under the MIT License, see LICENSE.
