# rexdep [![CI Status](https://github.com/itchyny/rexdep/workflows/CI/badge.svg)](https://github.com/itchyny/rexdep/actions)
### Roughly extract dependency relation from source code
The rexdep command is a tool for extracting dependency relation from software.
The command enables us to see the dependency relation among files, modules or packages written in any programming languages.

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

So, how can we extract dependency relation from a code base of software?

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
For example, there are some UML class diagram generating tools for object-oriented languages.
They investigate source codes at the level of abstract syntax tree and therefore much powerful.
The rexdep command, on the other hand, simply checks the source code by a regular expression given by the user.
It may not as powerful as AST-level tools, but the idea is very simple and can be used for many languages.
Now you understand what rexdep stands for; *roughly extract dependency relation*.

## Installation
### Homebrew
```bash
brew install itchyny/tap/rexdep
```

### Build from source
```bash
go get github.com/itchyny/rexdep
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

Note that on Windows environment, use double quotes instead of single quotes (for example: `rexdep --pattern "import +(\S+)" test1`).

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

### Module name: --module
The rexdep command uses the filenames to identify each file.
This is useful for C language but not for many other languages.
For example, we write `module ModuleName where` in Haskell.
Consider the following two files.
```bash
 $ cat Foo.hs
module Foo where
import Bar
import System.Directory
 $ cat Bar.hs
module Bar where
import System.IO
```
Firstly, we try the following command.
```bash
 $ rexdep --pattern 'import (\S+)' Foo.hs Bar.hs
Foo.hs Bar
Foo.hs System.Directory
Bar.hs System.IO
```
The result looks bad and we will not be able to obtain the dependency graph we really want.

We have to tell rexdep to extract the module name as well.
The rexdep command enables us to specify the module name pattern.
We specify a regular expression to the argument of --module.
```bash
 $ rexdep --pattern 'import (\S+)' --module 'module (\S+)' Foo.hs Bar.hs
Foo Bar
Foo System.Directory
Bar System.IO
```
Now it looks fine.

A file sometimes contains multiple modules.
When rexdep finds a new module in the file, the following packages will be regarded as imported by the new module.
For example, let us consider the following sample file.
```bash
 $ cat sample1
module A
import B
import C

module B
import C

module C
import D
```
In this example, three modules exist in one file.
The following command works nice for this sample.
```bash
 $ rexdep --pattern 'import (\S+)' --module 'module (\S+)' sample1
A B
A C
B C
C D
```
Let me show another example.
```bash
 $ cat sample2
A depends on B, C and D.
B depends on C and D.
C depends on D.
D does not depend on other modules.
```
Can we extract the relations between them? With rexdep, the answer is yes.
```bash
 $ rexdep --pattern 'depends on ([A-Z]+)(?:, ([A-Z]+))?(?:, ([A-Z]+))?(?: and ([A-Z]+))?' --module '^([A-Z]+) depends on ' sample2
A B
A C
A D
B C
B D
C D
```
This example does not look like a practical use case.
However, consider the following Makefile.
```make
 $ cat Makefile
all: clean build

build: deps
	go build

install: deps
	go install

deps:
	go get -d -v .

clean:
	go clean

.PHONY: all build install deps clean
```
Now rexdep can extract the target dependency relation from this Makefile.
```sh
 $ rexdep --pattern '^[^.][^:]+: +(\S+)(?: +(\S+))?(?: +(\S+))?' --module '^([^.][^:]+):' Makefile
all build
all clean
build deps
install deps
```
The rexdep command enables the user to specify a regular expression and still seems to be useful for pretty complicated pattern like the above example.

### Extraction range: --start, --end
The rexdep command checks the regular expression of `--pattern` against each line of files.
However, it is sometimes difficult to extract from multiline syntax.
Here's an example.
```go
 $ cat sample.go
package main

// "foo"
import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	fmt.Printf("Hello")
	// ...
```
The rexdep enables you to specify the range to extract from.
```
 $ rexdep --pattern '"(\S+)"' --module '^package +(\S+)' --start '^import +\($' --end '^\)$' sample.go
main fmt
main github.com/urfave/cli
main os
```
When the argument of `--start` is specified, rexdep finds the line which matches to the regular expression.
After it hits the starting line, it turns the internal enabled flag on and starts the extraction with the regular expression of `--patern`.
Then it finds the ending line, it turns the enabled flag off and stops the extraction procedure.

Both starting and ending lines are inclusive.
For example, see the following example.
```scala
object Foo extends Bar with Baz with Qux {
  // random comment: X extends Y
  // random comment: Y with Z
}

object Qux
  extends Quux
     with Quuy
     with Quuz {
  // with Qyyy
}
```
Firstly, we try with only `--pattern` and `--module` to extract the inheritance relation.
```
 $ rexdep --pattern '(?:(?:extends|with) +([[:alnum:]]+))(?: +(?:extends|with) +([[:alnum:]]+))?(?: +(?:extends|with) +([[:alnum:]]+))?' --module '^object +(\S+)' sample.scala
Foo Bar
Foo Baz
Foo Qux
Foo Y
Foo Z
Qux Quux
Qux Quuy
Qux Quuz
Qux Qyyy
```
This result is a failure; it keeps the extraction procedure inside bodies of the objects.
For this example, `--start` and `--end` work well.
```sh
 $ rexdep --pattern '(?:(?:extends|with) +([[:alnum:]]+))(?: +(?:extends|with) +([[:alnum:]]+))?(?: +(?:extends|with) +([[:alnum:]]+))?' --module '^object +(\S+)' --start '^object' --end '{' sample.scala
Foo Bar
Foo Baz
Foo Qux
Qux Quux
Qux Quuy
Qux Quuz
```
The rexdep command stops the extraction when it finds the line matches against the regular expression of `--end`.

### Output: --format, --output
We can change the output format with --format option.
The rexdep command outputs the module names each line with a space by default.
In order to see the dependency graph, which is a very common case, you can specify `--format dot` and pipe the output to the dot command.
You can also use `--format json` and pipe to the jq command.

The rexdep uses the standard output to print the result.
When you specify a filename to the --output option, the rexdep command outputs the result to the specified file.

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
 $ rexdep --pattern '"github.com/(?:hashicorp/consul/(?:\S+/)*)?(\S+)"' --module '^package +(\S+)' --start '^import +["(]' --end '^\)$|^import +"' --format dot $(find ./consul/ -name '*.go' | grep -v '_test') | dot -Tpng -o consul.png
```
[![consul](https://raw.githubusercontent.com/wiki/itchyny/rexdep/image/consul-1.png)](https://raw.githubusercontent.com/wiki/itchyny/rexdep/image/consul.png)
It is difficult to extract dependency relation between files from Go source codes. We can use functions from the other files at the same directory without writing import. Instead, we can extract dependency relation between packages. The rexdep command extract the imported packages between the lines matched by the `start` and `end` arguments.

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
