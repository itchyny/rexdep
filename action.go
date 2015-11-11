package main

import (
	"fmt"
	"io"

	"github.com/codegangsta/cli"
)

func action(ctx *cli.Context) {
	config, errors := makeConfig(ctx)
	if errors != nil {
		for _, err := range errors {
			fmt.Fprintf(ctx.App.Writer, err.Error())
		}
		cli.ShowAppHelp(ctx)
	} else {
		dependencies, errors := gatherDependency(config)
		output(ctx.App.Writer, config, dependencies, errors)
	}
}

func output(writer io.Writer, config *Config, dependencies []*Dependency, errors []error) {
	for _, err := range errors {
		fmt.Fprintf(os.Stderr, "# "+err.Error()+"\n")
	}
	indent := ""
	if config.Digraph != "" {
		fmt.Fprintf(writer, "digraph \"%s\" {\n", config.Digraph)
		indent = "  "
	}
	for _, dependency := range dependencies {
		for _, to := range dependency.To {
			fmt.Fprintf(writer, "%s\"%s\" -> \"%s\";\n", indent, dependency.From, to)
		}
	}
	if config.Digraph != "" {
		fmt.Fprintf(writer, "}\n")
	}
}

func gatherDependency(config *Config) ([]*Dependency, []error) {
	var errors []error
	var dependencies []*Dependency
	for _, path := range config.Paths {
		deps, err := extract(path, config)
		if err != nil {
			errors = append(errors, err...)
		} else {
			dependencies = append(dependencies, deps...)
		}
	}
	return dependencies, errors
}
