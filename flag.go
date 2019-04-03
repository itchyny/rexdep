package main

import "github.com/urfave/cli"

var flags = []cli.Flag{
	cli.StringFlag{
		Name:  "pattern, p",
		Value: "",
		Usage: "Pattern to extract imports (example: '^import\\s+(\\S+)')",
	},
	cli.StringFlag{
		Name:  "module, m",
		Value: "",
		Usage: "Pattern to extract module names (example: '^module\\s+(\\S+)')",
	},
	cli.BoolFlag{
		Name:  "reverse",
		Usage: "Assume reversed order; module names appear after dependant patterns",
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
		Name:  "root",
		Usage: "Regard the specified directory as the root directory of the modules",
	},
	cli.StringFlag{
		Name:  "format",
		Value: "",
		Usage: "Specify the format of the output; 'default', 'dot', 'csv', 'tsv' and 'json' are available.",
	},
	cli.StringFlag{
		Name:  "output, o",
		Value: "",
		Usage: "Output file. If omitted, stdout is used.",
	},
	cli.BoolFlag{
		Name:  "help, h",
		Usage: "Shows the help of the command",
	},
}
