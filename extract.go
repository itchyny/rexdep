package main

import (
	"bufio"
	"errors"
	"os"
	"path"
	"path/filepath"
)

func extract(name string, config *Config) (*Dependency, []error) {
	fileinfo, err := os.Lstat(name)
	if err != nil {
		return nil, []error{err}
	}
	if fileinfo.IsDir() {
		if config.Recursive {
			dependency := newDependency()
			var errs []error
			filepath.Walk(name, func(name string, info os.FileInfo, err error) error {
				if err == nil && !info.IsDir() {
					deps, err := extractFile(name, config)
					errs = append(errs, err...)
					if deps != nil {
						dependency.concat(deps)
					}
				}
				return nil
			})
			return dependency, errs
		}
		err := errors.New(name + " is a directory. Specify source code files. Or you mean --recursive (-r)?")
		return nil, []error{err}
	}
	return extractFile(name, config)
}

func extractFile(name string, config *Config) (*Dependency, []error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, []error{err}
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	module := path.Base(name)
	if config.Reldir != "" {
		name, err = filepath.Abs(name)
		if err != nil {
			return nil, []error{err}
		}
		module, err = filepath.Rel(config.Reldir, name)
		if err != nil {
			return nil, []error{err}
		}
	}
	return extractCore(module, scanner, config), nil
}

func extractCore(module string, scanner *bufio.Scanner, config *Config) *Dependency {
	dependency := newDependency()
	enable := config.Start == nil
	for scanner.Scan() {
		line := scanner.Text()
		if config.Start != nil && config.Start.MatchString(line) {
			enable = true
		}
		if enable {
			if config.Module != nil {
				if matches := config.Module.FindStringSubmatch(line); matches != nil {
					module = matches[len(matches)-1]
				}
			}
			if matches := config.Pattern.FindStringSubmatch(line); len(matches) >= 1 {
				for _, name := range matches[1:] {
					if name != "" {
						dependency.add(module, name)
					}
				}
			}
		}
		if enable && config.End != nil && config.End.MatchString(line) {
			enable = false
		}
	}
	return dependency
}
