package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/codegangsta/cli"
)

// Config is the command configuration
type Config struct {
	Pattern   *regexp.Regexp
	Module    *regexp.Regexp
	Reverse   bool
	Start     *regexp.Regexp
	End       *regexp.Regexp
	Format    string
	Paths     []string
	Recursive bool
	Root      string
	Output    io.Writer
}

func makeConfig(ctx *cli.Context) (*Config, []error) {
	var errs []error

	if ctx.GlobalBool("help") || ctx.NumFlags() == 0 {
		errs = append(errs, errors.New(""))
		return nil, errs
	}

	if ctx.GlobalString("pattern") == "" {
		errs = append(errs, errors.New("Specify --pattern (-p) to extract imports.\n\n"))
	}

	pattern, err := regexp.Compile(ctx.GlobalString("pattern"))
	if err != nil {
		errs = append(errs, errors.New(regexErrorMessage("--pattern (-p)")+err.Error()+"\n\n"))
	}

	module, err := regexp.Compile(ctx.GlobalString("module"))
	if err != nil {
		errs = append(errs, errors.New(regexErrorMessage("--module (-m)")+err.Error()+"\n\n"))
	}
	if ctx.GlobalString("module") == "" {
		module = nil
	}

	start, err := regexp.Compile(ctx.GlobalString("start"))
	if err != nil {
		errs = append(errs, errors.New(regexErrorMessage("--start (-s)")+err.Error()+"\n\n"))
	}
	if ctx.GlobalString("start") == "" {
		start = nil
	}

	end, err := regexp.Compile(ctx.GlobalString("end"))
	if err != nil {
		errs = append(errs, errors.New(regexErrorMessage("--end (-e)")+err.Error()+"\n\n"))
	}
	if ctx.GlobalString("end") == "" {
		end = nil
	}

	root := ctx.GlobalString("root")
	if root != "" {
		root, err = filepath.Abs(root)
		if err != nil {
			errs = append(errs, errors.New(regexErrorMessage("--root")+err.Error()+"\n\n"))
		}
	}

	output := ctx.App.Writer
	outfile := ctx.GlobalString("output")
	if outfile != "" {
		file, err := os.Create(outfile)
		if err != nil {
			errs = append(errs, errors.New("Cannot create the output file: "+outfile+"\n\n"))
		} else {
			output = file
		}
	}

	paths := ctx.Args()
	if len(paths) == 0 {
		errs = append(errs, errors.New("Specify source codes.\n\n"))
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return &Config{
		Pattern:   pattern,
		Module:    module,
		Reverse:   ctx.GlobalBool("reverse"),
		Start:     start,
		End:       end,
		Format:    ctx.GlobalString("format"),
		Paths:     paths,
		Recursive: ctx.GlobalBool("recursive"),
		Root:      root,
		Output:    output,
	}, nil
}

func regexErrorMessage(flag string) string {
	return "The argument of " + flag + " is invalid. Specify a valid regular expression.\n"
}
