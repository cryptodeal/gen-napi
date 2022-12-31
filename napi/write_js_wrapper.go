package napi

import (
	"fmt"
	"strings"
)

func (g *PackageGenerator) WriteEnvWrapper(sb *strings.Builder, classes map[string]*CPPClass, methods map[string]*CPPMethod, processedMethods map[string]*CPPMethod) {
	sb.WriteString(g.WriteEnvImports(classes, methods, processedMethods))
	for name, c := range classes {
		if c.Decl != nil {
			sb.WriteString(g.WriteEnvClassWrapper(name, c, methods, processedMethods))
		}
	}
}

func (g *PackageGenerator) WriteEnvImports(classes map[string]*CPPClass, methods map[string]*CPPMethod, processedMethods map[string]*CPPMethod) string {
	sb := new(strings.Builder)
	sb.WriteString("const {\n")
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
		if i < used_len-1 {
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
		if i < used_len-1 {
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
		if i < used_len-1 {
			sb.WriteString(",\n")
		}
	}
	sb.WriteString(fmt.Sprintf("\n} = require('%s')\n\n", g.conf.ResolvedBindingsImportPath(g.conf.Path)))
	return sb.String()
}

func (g *PackageGenerator) WriteEnvClassWrapper(className string, class *CPPClass, methods map[string]*CPPMethod, processedMethods map[string]*CPPMethod) string {
	sb := new(strings.Builder)
	if g.conf.IsEnvTS() {
		sb.WriteString("export ")
	}
	sb.WriteString(fmt.Sprintf("class %s {\n", className))

	g.writeIndent(sb, 1)
	if g.conf.IsEnvTS() {
		sb.WriteString("private #_native_self: any;\n")
	} else {
		sb.WriteString("#_native_self")
	}

	g.writeIndent(sb, 1)
	sb.WriteString("constructor(t) {\n")
	g.writeIndent(sb, 2)
	sb.WriteString(fmt.Sprintf("this.#_native_self = new _%s(t);\n", className))
	g.writeIndent(sb, 1)
	sb.WriteString("}\n")

	for _, m := range methods {
		if g.conf.IsMethodWrapped(className, *m.Ident) {
			g.writeIndent(sb, 1)
			sb.WriteString(fmt.Sprintf("%s(...args) {\n", *m.Ident))
			g.writeIndent(sb, 2)
			sb.WriteString(fmt.Sprintf("return this.#_native_self.%s(...args);\n", *m.Ident))
			g.writeIndent(sb, 1)
			sb.WriteString("}\n\n")
		}
	}

	sb.WriteString("}\n\n")
	return sb.String()
}
