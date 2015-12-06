package main

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
)

func outputDefault(writer io.Writer, dependency *Dependency) {
	for _, module := range dependency.modules {
		for _, to := range keys(dependency.relation[module]) {
			fmt.Fprintf(writer, "%s %s\n", module, to)
		}
	}
}

func outputDot(writer io.Writer, dependency *Dependency) {
	fmt.Fprintf(writer, "digraph \"graph\" {\n")
	for _, module := range dependency.modules {
		for _, to := range keys(dependency.relation[module]) {
			fmt.Fprintf(writer, "  %s -> %s;\n", strconv.Quote(module), strconv.Quote(to))
		}
	}
	fmt.Fprintf(writer, "}\n")
}

func outputCsv(writer io.Writer, dependency *Dependency) {
	for _, module := range dependency.modules {
		for _, to := range keys(dependency.relation[module]) {
			fmt.Fprintf(writer, "%s,%s\n", strconv.Quote(module), strconv.Quote(to))
		}
	}
}

func outputTsv(writer io.Writer, dependency *Dependency) {
	escape := func(str string) string {
		unescape := regexp.MustCompile(`\\([\"'\\])`)
		ret := unescape.ReplaceAllString(strconv.Quote(str), "$1")
		return ret[1 : len(ret)-1]
	}
	for _, module := range dependency.modules {
		for _, to := range keys(dependency.relation[module]) {
			fmt.Fprintf(writer, "%s\t%s\n", escape(module), escape(to))
		}
	}
}

func outputJSON(writer io.Writer, dependency *Dependency) {
	fmt.Fprintf(writer, "{")
	for i, module := range dependency.modules {
		if i > 0 {
			fmt.Fprintf(writer, ",")
		}
		fmt.Fprintf(writer, "\n  %s: [", strconv.Quote(module))
		for j, to := range keys(dependency.relation[module]) {
			if j > 0 {
				fmt.Fprintf(writer, ",")
			}
			fmt.Fprintf(writer, "\n    %s", strconv.Quote(to))
		}
		fmt.Fprintf(writer, "\n  ]")
	}
	fmt.Fprintf(writer, "\n}\n")
}

func keys(m map[string]bool) []string {
	var xs []string
	for x := range m {
		xs = append(xs, x)
	}
	sort.Strings(xs)
	return xs
}
