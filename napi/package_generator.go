package napi

import (
	"strings"
)

func (g *PackageGenerator) Generate() (string, string, error) {
	bindings := new(strings.Builder)
	env_wrapper := new(strings.Builder)

	g.ParsedData.Methods = g.parseMethods(g.RootNode, *g.Input)
	g.ParsedData.Classes = g.parseClasses(g.RootNode, *g.Input)
	g.ParsedData.Lits = g.parseLiterals(g.RootNode, *g.Input)
	g.ParsedData.Enums = g.parseEnumDecls(g.RootNode, *g.Input, true)
	g.writeFileCodegenHeader(bindings)
	g.writeFileCodegenHeader(env_wrapper)

	g.WriteBindings(bindings)
	g.WriteEnvWrapper(env_wrapper)

	return bindings.String(), env_wrapper.String(), nil
}
