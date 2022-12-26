package napi

import (
	"strings"
)

func (g *PackageGenerator) Generate() (string, string, error) {
	bindings := new(strings.Builder)
	headers := new(strings.Builder)

	methods := parseMethods(g.RootNode, *g.Input)
	classes := parseClasses(g.RootNode, *g.Input)

	// write headers for generated file for specific package
	g.writeFileCodegenHeader(bindings)
	g.writeFileCodegenHeader(headers)

	// write headers
	g.writeHeader(headers, classes, methods)
	g.writeBindings(bindings, classes, methods)

	return bindings.String(), headers.String(), nil
}
