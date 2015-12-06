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
		dependency, errors := gatherDependency(config)
		output(config, dependency, errors)
	}
}

func output(config *Config, dependency *Dependency, errors []error) {
	for _, err := range errors {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
	}
	switch config.Format {
	case "dot":
		outputDot(config.Output, dependency)
	case "csv":
		outputCsv(config.Output, dependency)
	case "tsv":
		outputTsv(config.Output, dependency)
	case "json":
		outputJSON(config.Output, dependency)
	default:
		outputDefault(config.Output, dependency)
	}
}

func gatherDependency(config *Config) (*Dependency, []error) {
	var errors []error
	dependency := newDependency()
	for _, path := range config.Paths {
		deps, err := extract(path, config)
		if err != nil {
			errors = append(errors, err...)
		} else {
			dependency.concat(deps)
		}
	}
	return dependency, errors
}
