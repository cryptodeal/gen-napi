package napi

import (
	"fmt"
	"strings"
)

func (g *PackageGenerator) Generate() (string, string, error) {
	bindings := new(strings.Builder)
	env_wrapper := new(strings.Builder)

	g.ParsedData.Methods = g.parseMethods(g.RootNode, *g.Input)
	g.ParsedData.Classes = parseClasses(g.RootNode, *g.Input)
	g.ParsedData.Lits = g.parseLiterals(g.RootNode, *g.Input)
	g.ParsedData.Enums = g.parseEnums(g.RootNode, *g.Input, true)
	for _, enum := range g.ParsedData.Enums {
		fmt.Printf("Name: %q; Values:\n", *enum.Ident)
		for _, v := range enum.Values {
			fmt.Printf("\t%s = %d\n", *v.Ident, v.Value)
		}
	}

	g.writeFileCodegenHeader(bindings)
	g.writeFileCodegenHeader(env_wrapper)

	g.writeBindings(bindings)
	g.WriteEnvWrapper(env_wrapper)

	return bindings.String(), env_wrapper.String(), nil
}
