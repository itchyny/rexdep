package main

import (
	"fmt"
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
		output(config, dependencies, errors)
	}
}

func output(config *Config, dependencies []*Dependency, errors []error) {
	for _, err := range errors {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
	}
	switch config.Format {
	case "dot":
		outputDot(config.Output, dependencies)
	case "csv":
		outputCsv(config.Output, dependencies)
	case "json":
		outputJson(config.Output, dependencies)
	default:
		outputDefault(config.Output, dependencies)
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
