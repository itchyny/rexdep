package main

import "os"

var name = "rexdep"
var version = "v0.0.0"
var description = "Roughly extract dependency from source code"
var author = "itchyny"

func main() {
	os.Exit(run(os.Args))
}
