package main

import (
	"bufio"
	"errors"
	"os"
	"path"
	"path/filepath"
)

func extract(name string, config *Config) ([]*Dependency, []error) {
	fileinfo, err := os.Lstat(name)
	if err != nil {
		return nil, []error{err}
	}
	if fileinfo.IsDir() {
		if config.Recursive {
			var dependencies []*Dependency
			var errs []error
			filepath.Walk(name, func(name string, info os.FileInfo, err error) error {
				if err == nil && !info.IsDir() {
					deps, err := extractFile(name, config)
					errs = append(errs, err...)
					dependencies = append(dependencies, deps...)
				}
				return nil
			})
			return dependencies, errs
		} else {
			err := errors.New(name + " is a directory. Specify source code files. Or you mean --recursive (-r)?")
			return nil, []error{err}
		}
	} else {
		return extractFile(name, config)
	}
}

func extractFile(name string, config *Config) ([]*Dependency, []error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, []error{err}
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	module := path.Base(name)
	return extractCore(module, scanner, config), nil
}

func extractCore(module string, scanner *bufio.Scanner, config *Config) []*Dependency {
	dependency := &Dependency{From: module}
	dependencies := []*Dependency{}
	appended := make(map[string]bool)
	enable := config.Start == nil
	for scanner.Scan() {
		line := scanner.Text()
		if config.Start != nil && config.Start.MatchString(line) {
			enable = true
		}
		if config.Module != nil {
			if matches := config.Module.FindStringSubmatch(line); matches != nil {
				dependency = &Dependency{From: matches[len(matches)-1]}
				appended = make(map[string]bool)
			}
		}
		if enable {
			if matches := config.Pattern.FindStringSubmatch(line); len(matches) >= 1 {
				for _, name := range matches[1:] {
					if name != "" && !appended[name] {
						if dependency.To == nil {
							dependencies = append(dependencies, dependency)
						}
						dependency.To = append(dependency.To, name)
						appended[name] = true
					}
				}
			}
		}
		if enable && config.End != nil && config.End.MatchString(line) {
			enable = false
		}
	}
	return dependencies
}
