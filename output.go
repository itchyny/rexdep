package main

import (
	"fmt"
	"io"
)

func outputDefault(writer io.Writer, dependencies []*Dependency) {
	for _, dependency := range dependencies {
		for _, to := range dependency.To {
			fmt.Fprintf(writer, "%s %s\n", dependency.From, to)
		}
	}
}

func outputDot(writer io.Writer, dependencies []*Dependency) {
	fmt.Fprintf(writer, "digraph \"graph\" {\n")
	for _, dependency := range dependencies {
		for _, to := range dependency.To {
			fmt.Fprintf(writer, "  \"%s\" -> \"%s\";\n", dependency.From, to)
		}
	}
	fmt.Fprintf(writer, "}\n")
}

func outputCsv(writer io.Writer, dependencies []*Dependency) {
	for _, dependency := range dependencies {
		for _, to := range dependency.To {
			fmt.Fprintf(writer, "%s,%s\n", dependency.From, to)
		}
	}
}

func outputJson(writer io.Writer, dependencies []*Dependency) {
	fmt.Fprintf(writer, "{")
	for i, dependency := range dependencies {
		if i > 0 {
			fmt.Fprintf(writer, ",")
		}
		fmt.Fprintf(writer, "\n  \"%s\": [\n", dependency.From)
		for _, to := range dependency.To {
			fmt.Fprintf(writer, "  \"%s\",\n", to)
		}
		fmt.Fprintf(writer, "]")
	}
	fmt.Fprintf(writer, "\n}\n")
}
