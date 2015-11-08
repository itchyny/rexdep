package main

import (
	"github.com/codegangsta/cli"
)

var Flags = []cli.Flag{
	cli.StringFlag{
		Name:  "module, m",
		Value: "",
		Usage: "Pattern to extract module names (example: '^module\\s+(\\S+)')",
	},
	cli.StringFlag{
		Name:  "pattern, p",
		Value: "",
		Usage: "Pattern to extract imports (example: '^import\\s+(\\S+)')",
	},
	cli.StringFlag{
		Name:  "start, s",
		Value: "",
		Usage: "Pattern to start matching",
	},
	cli.StringFlag{
		Name:  "end, e",
		Value: "",
		Usage: "Pattern to end matching",
	},
	cli.BoolFlag{
		Name:  "recursive, r",
		Usage: "Recursively inspect files in subdirectories",
	},
	cli.StringFlag{
		Name:  "digraph, g",
		Value: "",
		Usage: "Specify the name of directed graph for graphviz. If omited, the command does not print digraph line.",
	},
	cli.BoolFlag{
		Name:  "trimext, t",
		Usage: "Trim the file extension.",
	},
	cli.BoolFlag{
		Name:  "help, h",
		Usage: "Shows the help of the command",
	},
}
