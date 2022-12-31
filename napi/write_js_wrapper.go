package napi

import (
	"fmt"
	"strings"
)

func (g *PackageGenerator) WriteEnvWrapper(sb *strings.Builder, classes map[string]*CPPClass, methods map[string]*CPPMethod, processedMethods map[string]*CPPMethod) {
	sb.WriteString(g.WriteEnvImports(classes, methods, processedMethods))
}

func (g *PackageGenerator) WriteEnvImports(classes map[string]*CPPClass, methods map[string]*CPPMethod, processedMethods map[string]*CPPMethod) string {
	sb := &strings.Builder{}
	used := []string{}
	for name, c := range classes {
		if c.Decl != nil {
			used = append(used, name)
		}
	}
	used_len := len(used)
	for i, name := range used {
		g.writeIndent(sb, 1)
		sb.WriteString(fmt.Sprintf("_%s", name))
		if i < used_len {
			sb.WriteString(",\n")
		}
	}
	used = []string{}
	for name, m := range methods {
		if !g.conf.IsMethodIgnored(*m.Ident) {
			used = append(used, name)
		}
	}
	used_len = len(used)
	for i, name := range used {
		if i == 0 {
			sb.WriteString(",\n")
		}
		g.writeIndent(sb, 1)
		sb.WriteString(fmt.Sprintf("_%s", name))
		if i < used_len {
			sb.WriteString(",\n")
		}
	}
	used = []string{}
	for name, m := range processedMethods {
		if !g.conf.IsMethodIgnored(*m.Ident) {
			used = append(used, name)
		}
	}
	used_len = len(used)
	for i, name := range used {
		if i == 0 {
			sb.WriteString(",\n")
		}
		g.writeIndent(sb, 1)
		sb.WriteString(fmt.Sprintf("_%s", name))
		if i < used_len {
			sb.WriteString(",\n")
		}
	}
	sb.WriteString(fmt.Sprintf("\n} = require('%s')\n\n", g.conf.JSWrapperOpts.AddonPath))
	return sb.String()
}
