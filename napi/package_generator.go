package napi

import (
	"strings"
)

func (g *PackageGenerator) Generate() (string, string, error) {
	bindings := new(strings.Builder)
	env_wrapper := new(strings.Builder)

	methods := g.parseMethods(g.RootNode, *g.Input)
	classes := parseClasses(g.RootNode, *g.Input)
	lits := g.parseLiterals(g.RootNode, *g.Input)

	g.writeFileCodegenHeader(bindings)
	g.writeFileCodegenHeader(env_wrapper)

	g.writeBindings(bindings, classes, methods, lits)
	g.WriteEnvWrapper(env_wrapper, classes, methods, lits)

	return bindings.String(), env_wrapper.String(), nil
}
