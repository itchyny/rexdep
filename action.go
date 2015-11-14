package main

import (
	"fmt"
	"io"
	"os"

	"github.com/codegangsta/cli"
)

func action(ctx *cli.Context) {
	config, errors := makeConfig(ctx)
	if errors != nil {
		for _, err := range errors {
			fmt.Fprintf(os.Stderr, err.Error())
		}
		cli.ShowAppHelp(ctx)
	} else {
		dependencies, errors := gatherDependency(config)
		output(ctx.App.Writer, config, dependencies, errors)
	}
}

func output(writer io.Writer, config *Config, dependencies []*Dependency, errors []error) {
	for _, err := range errors {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
	}
	switch config.Format {
	case "dot":
		outputDot(writer, dependencies)
	case "csv":
		outputCsv(writer, dependencies)
	case "json":
		outputJson(writer, dependencies)
	default:
		outputDefault(writer, dependencies)
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
