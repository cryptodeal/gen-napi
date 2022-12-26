package napi

import (
	"strings"
)

func (g *PackageGenerator) Generate() (string, string, error) {
	bindings := new(strings.Builder)
	headers := new(strings.Builder)

	// write headers for generated file for specific package
	g.writeFileCodegenHeader(bindings)
	g.writeFileCodegenHeader(headers)

	return bindings.String(), headers.String(), nil
}
