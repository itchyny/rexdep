package main

import "os"

var name = "rexdep"
var version = "0.3.2"
var description = "Roughly extract dependency from source code"
var author = "itchyny"

func main() {
	os.Exit(run(os.Args))
}
