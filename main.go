package main

import "os"

var Name = "rexdep"
var Version = "v0.0.0"
var Description = "Roughly extract dependency from source code"
var Author = "itchyny"

func main() {
	os.Exit(run(os.Args))
}
