package main

import (
	"github.com/codegangsta/cli"
)

var Name string
var Description string
var Version string
var Author string

func run(args []string) int {
	app := newApp()
	if app.Run(args) != nil {
		return 1
	}
	return 0
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Name = Name
	app.HelpName = Name
	app.Usage = Description
	app.Version = Version
	app.Author = Author
	app.Flags = Flags
	app.HideHelp = true
	app.Action = action
	return app
}
