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
	sb.WriteString(g.WriteEnvWrappedFns(methods, processedMethods))
	if !g.conf.IsEnvTS() {
		sb.WriteString(g.WriteEnvExports(classes, methods, processedMethods))
	}
}

func (g *PackageGenerator) WriteEnvExports(classes map[string]*CPPClass, methods map[string]*CPPMethod, processedMethods map[string]*CPPMethod) string {
	sb := new(strings.Builder)
	sb.WriteString("module.exports = {\n")
	used := []string{}
	for name, c := range classes {
		if c.Decl != nil {
			used = append(used, name)
		}
	}
	used_len := len(used)
	for i, name := range used {
		g.writeIndent(sb, 1)
		sb.WriteString(name)
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
		sb.WriteString(name)
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
		sb.WriteString(name)
		if i < used_len-1 {
			sb.WriteString(",\n")
		}
	}
	sb.WriteString("\n}\n")
	return sb.String()
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
		if name == "var" {
			sb.WriteString(fmt.Sprintf("__%s", name))
		} else {
			sb.WriteString(fmt.Sprintf("_%s", name))
		}
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
	sb.WriteString(fmt.Sprintf("\n} = require(%q)\n\n", g.conf.ResolvedBindingsImportPath(g.conf.Path)))
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
		sb.WriteString("#_native_self\n")
	}
	sb.WriteByte('\n')
	g.writeIndent(sb, 1)
	sb.WriteString("constructor(t) {\n")
	g.writeIndent(sb, 2)
	sb.WriteString(fmt.Sprintf("this.#_native_self = new _%s(t);\n", className))
	g.writeIndent(sb, 1)
	sb.WriteString("}\n\n")

	for _, m := range methods {
		if g.conf.IsMethodWrapped(className, *m.Ident) {
			g.writeIndent(sb, 1)
			sb.WriteString(fmt.Sprintf("%s(", *m.Ident))
			for i, p := range *m.Overloads[0] {
				if i == 0 {
					continue
				}
				if i > 1 && i < len(*m.Overloads[0]) {
					sb.WriteString(", ")
				}
				sb.WriteString(*p.Ident)
				if g.conf.IsEnvTS() {
					tsType, _ := CPPTypeToTS(*p.Type)
					sb.WriteString(fmt.Sprintf(": %s", tsType))
				}
			}
			if g.conf.IsEnvTS() {
				tsType, _ := CPPTypeToTS(*m.Returns)
				sb.WriteString(fmt.Sprintf("): %s {\n", tsType))
			} else {
				sb.WriteString(") {\n")
			}
			g.writeIndent(sb, 2)
			sb.WriteString(fmt.Sprintf("return this.#_native_self.%s(this.#_native_self,", *m.Ident))
			for i, p := range *m.Overloads[0] {
				if i == 0 {
					continue
				}
				if i > 1 && i < len(*m.Overloads[0]) {
					sb.WriteString(", ")
				}
				sb.WriteString(*p.Ident)
			}
			sb.WriteString(");\n")
			g.writeIndent(sb, 1)
			sb.WriteString("}\n\n")
		}
	}

	sb.WriteString("}\n\n")
	return sb.String()
}

func (g *PackageGenerator) WriteEnvWrappedFns(methods map[string]*CPPMethod, processedMethods map[string]*CPPMethod) string {
	sb := new(strings.Builder)
	for _, m := range methods {
		if !g.conf.IsMethodIgnored(*m.Ident) {
			if *m.Ident == "var" {
				sb.WriteString(fmt.Sprintf("const %s = (", "_"+*m.Ident))
			} else {
				sb.WriteString(fmt.Sprintf("const %s = (", *m.Ident))
			}
			for i, p := range *m.Overloads[0] {
				if i > 0 && i < len(*m.Overloads[0]) {
					sb.WriteString(", ")
				}
				sb.WriteString(*p.Ident)
			}
			sb.WriteString(") => {\n")
			g.writeIndent(sb, 1)
			if *m.Ident == "var" {
				sb.WriteString(fmt.Sprintf("return __%s(", *m.Ident))
			} else {
				sb.WriteString(fmt.Sprintf("return _%s(", *m.Ident))
			}
			for i, p := range *m.Overloads[0] {
				if i > 0 && i < len(*m.Overloads[0]) {
					sb.WriteString(", ")
				}
				sb.WriteString(*p.Ident)
			}
			sb.WriteString(");\n")
			sb.WriteString("}\n\n")
		}
	}

	for _, m := range processedMethods {
		if !g.conf.IsMethodIgnored(*m.Ident) {
			if !g.conf.IsMethodIgnored(*m.Ident) {
				if *m.Ident == "var" {
					sb.WriteString(fmt.Sprintf("const %s = (", "_"+*m.Ident))
				} else {
					sb.WriteString(fmt.Sprintf("const %s = (", *m.Ident))
				}
				for i, p := range *m.Overloads[0] {
					if i > 0 && i < len(*m.Overloads[0]) {
						sb.WriteString(", ")
					}
					sb.WriteString(*p.Ident)
				}
				sb.WriteString(") => {\n")
				g.writeIndent(sb, 1)
				if *m.Ident == "var" {
					sb.WriteString(fmt.Sprintf("return __%s(", *m.Ident))
				} else {
					sb.WriteString(fmt.Sprintf("return _%s(", *m.Ident))
				}
				for i, p := range *m.Overloads[0] {
					if i > 0 && i < len(*m.Overloads[0]) {
						sb.WriteString(", ")
					}
					sb.WriteString(*p.Ident)
				}
				sb.WriteString(");\n")
				sb.WriteString("}\n\n")
			}
		}
	}
	return sb.String()
}
