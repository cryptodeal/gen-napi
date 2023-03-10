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
	sb.WriteString(g.WriteEnvWrappedFns(methods, processedMethods, classes))
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
			if name == "var" {
				used = append(used, "_"+name)
			} else {
				used = append(used, name)
			}
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
			if name == "var" {
				used = append(used, "_"+name)
			} else {
				used = append(used, name)
			}
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
	used_len = len(g.conf.GlobalForcedMethods)
	for i, m := range g.conf.GlobalForcedMethods {
		if i == 0 {
			sb.WriteString(",\n")
		}
		g.writeIndent(sb, 1)
		sb.WriteString(m.Name)
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
			sb.WriteString(fmt.Sprintf("_%s: __%s", name, name))
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
		if name == "var" {
			sb.WriteString(fmt.Sprintf("_%s: __%s", name, name))
		} else {
			sb.WriteString(fmt.Sprintf("_%s", name))
		}
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
		if name == "var" {
			sb.WriteString(fmt.Sprintf("_%s: __%s", name, name))
		} else {
			sb.WriteString(fmt.Sprintf("_%s", name))
		}
		if i < used_len-1 {
			sb.WriteString(",\n")
		}
	}
	used_len = len(g.conf.GlobalForcedMethods)
	for i, m := range g.conf.GlobalForcedMethods {
		if i == 0 {
			sb.WriteString(",\n")
		}
		g.writeIndent(sb, 1)
		if m.Name == "var" {
			sb.WriteString(fmt.Sprintf("_%s: __%s", m.Name, m.Name))
		} else {
			sb.WriteString(fmt.Sprintf("_%s", m.Name))
		}
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
	sb.WriteString("#_native_self")
	if g.conf.IsEnvTS() {
		sb.WriteString(": any")
	}
	sb.WriteString(";\n\n")
	g.writeIndent(sb, 1)
	sb.WriteString("constructor(t) {\n")
	g.writeIndent(sb, 2)
	sb.WriteString(fmt.Sprintf("this.#_native_self = new _%s(t);\n", className))
	g.writeIndent(sb, 1)
	sb.WriteString("}\n\n")

	g.writeIndent(sb, 1)
	sb.WriteString("get _native_self() {\n")
	g.writeIndent(sb, 2)
	sb.WriteString("return this.#_native_self;\n")
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
				if i >= m.ExpectedArgs {
					break
				}
				if i > 1 && i < len(*m.Overloads[0]) {
					sb.WriteString(", ")
				}
				sb.WriteString(*p.Ident)
				if g.conf.IsEnvTS() {
					tsType, _ := CPPTypeToTS(*p.Type)
					if v, ok := g.conf.TypeMappings[tsType]; ok {
						sb.WriteString(fmt.Sprintf(": %s", v.TSType))
					} else {
						if strings.Contains(tsType, "std::vector") {
							vectorType := tsType[strings.Index(tsType, "<")+1 : strings.Index(tsType, ">")]
							tsType, _ = CPPTypeToTS(vectorType)
							tsType = tsType + "[]"
						}
						sb.WriteString(fmt.Sprintf(": %s", tsType))
					}
				}
			}
			tsType, _ := CPPTypeToTS(*m.Returns)
			if g.conf.IsEnvTS() {
				if v, ok := g.conf.TypeMappings[tsType]; ok {
					sb.WriteString(fmt.Sprintf("): %s {\n", v.TSType))
				} else {
					sb.WriteString(fmt.Sprintf("): %s {\n", tsType))
				}
			} else {
				sb.WriteString(") {\n")
			}
			g.writeIndent(sb, 2)
			sb.WriteString("return ")
			if tsType == className {
				sb.WriteString(fmt.Sprintf("new %s(", tsType))
			}
			sb.WriteString(fmt.Sprintf("this.#_native_self.%s(", *m.Ident))
			for i, p := range *m.Overloads[0] {
				if i == 0 {
					continue
				}
				if i >= m.ExpectedArgs {
					break
				}
				if i > 1 && i < len(*m.Overloads[0]) {
					sb.WriteString(", ")
				}
				sb.WriteString(*p.Ident)
				if *p.Type == className {
					sb.WriteString(".#_native_self")
				}
			}
			if tsType == className {
				sb.WriteByte(')')
			}
			sb.WriteString(");\n")
			g.writeIndent(sb, 1)
			sb.WriteString("}\n\n")
		}
	}

	for _, m := range processedMethods {
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
					if v, ok := g.conf.TypeMappings[tsType]; ok {
						sb.WriteString(fmt.Sprintf(": %s", v.TSType))
					} else {
						if strings.Contains(tsType, "std::vector") {
							vectorType := tsType[strings.Index(tsType, "<")+1 : strings.Index(tsType, ">")]
							tsType, _ = CPPTypeToTS(vectorType)
							tsType = tsType + "[]"
						}
						sb.WriteString(fmt.Sprintf(": %s", tsType))
					}
				}
			}
			tsType, _ := CPPTypeToTS(*m.Returns)
			if g.conf.IsEnvTS() {
				if v, ok := g.conf.TypeMappings[tsType]; ok {
					sb.WriteString(fmt.Sprintf("): %s {\n", v.TSType))
				} else {
					if strings.Contains(tsType, "std::vector") {
						vectorType := tsType[strings.Index(tsType, "<")+1 : strings.Index(tsType, ">")]
						tsType, _ = CPPTypeToTS(vectorType)
						tsType = tsType + "[]"
					}
					sb.WriteString(fmt.Sprintf("): %s {\n", tsType))
				}
			} else {
				sb.WriteString(") {\n")
			}
			g.writeIndent(sb, 2)
			sb.WriteString("return ")
			if tsType == className {
				sb.WriteString(fmt.Sprintf("new %s(", tsType))
			}
			sb.WriteString(fmt.Sprintf("this.#_native_self.%s(", *m.Ident))
			for i, p := range *m.Overloads[0] {
				if i == 0 {
					continue
				}
				if i > 1 && i < len(*m.Overloads[0]) {
					sb.WriteString(", ")
				}
				sb.WriteString(*p.Ident)
				if *p.Type == className {
					sb.WriteString(".#_native_self")
				}
			}
			if tsType == className {
				sb.WriteByte(')')
			}
			sb.WriteString(");\n")
			g.writeIndent(sb, 1)
			sb.WriteString("}\n\n")
		}
	}

	for _, m := range *class.FieldDecl {
		if m.Ident != nil {
			if g.conf.IsFieldWrapped(className, *m.Ident) {
				g.writeIndent(sb, 1)
				sb.WriteString(fmt.Sprintf("%s(", *m.Ident))
				if m.Args != nil {
					for i, p := range *m.Args {
						if i == 0 {
							continue
						}
						if i > 1 && i < len(*m.Args) {
							sb.WriteString(", ")
						}
						sb.WriteString(*p.Ident)
						if g.conf.IsEnvTS() {
							tsType, _ := CPPTypeToTS(*p.Type)
							if v, ok := g.conf.TypeMappings[tsType]; ok {
								sb.WriteString(fmt.Sprintf(": %s", v.TSType))
							} else {
								if strings.Contains(tsType, "std::vector") {
									vectorType := tsType[strings.Index(tsType, "<")+1 : strings.Index(tsType, ">")]
									tsType, _ = CPPTypeToTS(vectorType)
									tsType = tsType + "[]"
								}
								sb.WriteString(fmt.Sprintf(": %s", tsType))
							}
						}
					}
				}
				tsType, _ := CPPTypeToTS(*m.Returns.FullType)
				if g.conf.IsEnvTS() {
					if m.Returns.FullType == nil {
						sb.WriteString(") {\n")
					} else {
						if v, ok := g.conf.TypeMappings[tsType]; ok {
							sb.WriteString(fmt.Sprintf("): %s {\n", v.TSType))
						} else {
							sb.WriteString(fmt.Sprintf("): %s {\n", tsType))
						}
					}
				} else {
					sb.WriteString(") {\n")
				}
				g.writeIndent(sb, 2)
				sb.WriteString("return ")
				if tsType == className {
					sb.WriteString(fmt.Sprintf("new %s(", tsType))
				}
				sb.WriteString(fmt.Sprintf("this.#_native_self.%s(", *m.Ident))
				if m.Args != nil {
					for i, p := range *m.Args {
						if i == 0 {
							continue
						}
						if i > 1 && i < len(*m.Args) {
							sb.WriteString(", ")
						}
						sb.WriteString(*p.Ident)
						if *p.Type == className {
							sb.WriteString(".#_native_self")
						}
					}
				}
				if tsType == className {
					sb.WriteByte(')')
				}
				sb.WriteString(");\n")
				g.writeIndent(sb, 1)
				sb.WriteString("}\n\n")
			}
		}
	}

	if v, ok := g.conf.ClassOpts[className]; ok {
		for _, m := range v.ForcedMethods {
			g.writeIndent(sb, 1)
			sb.WriteString(fmt.Sprintf("%s(", m.Name))
			for i, p := range m.Args {
				if i > 0 && i < len(m.Args) {
					sb.WriteString(", ")
				}
				sb.WriteString(p.Name)
				if g.conf.IsEnvTS() {
					sb.WriteString(fmt.Sprintf(": %s", p.TSType))
				}
			}
			if g.conf.IsEnvTS() {
				if m.IsVoid {
					sb.WriteString(") {\n")
				} else {
					sb.WriteString(fmt.Sprintf("): %s {\n", m.TSReturnType))
				}
			} else {
				sb.WriteString(") {\n")
			}
			g.writeIndent(sb, 2)
			sb.WriteString("return ")
			if m.TSReturnType == className {
				sb.WriteString(fmt.Sprintf("new %s(", m.TSReturnType))
			}
			sb.WriteString(fmt.Sprintf("this.#_native_self.%s(", m.Name))
			for i, p := range m.Args {
				if i > 0 && i < len(m.Args) {
					sb.WriteString(", ")
				}
				if p.TSType == className {
					sb.WriteString(".#_native_self")
				}
				sb.WriteString(p.Name)
			}
			if m.TSReturnType == className {
				sb.WriteByte(')')
			}
			sb.WriteString(");\n")
			g.writeIndent(sb, 1)
			sb.WriteString("}\n\n")
		}
	}

	sb.WriteString("}\n\n")
	return sb.String()
}

func (g *PackageGenerator) WriteEnvWrappedFns(methods map[string]*CPPMethod, processedMethods map[string]*CPPMethod, classes map[string]*CPPClass) string {
	sb := new(strings.Builder)
	for _, m := range methods {
		if !g.conf.IsMethodIgnored(*m.Ident) {
			if g.conf.IsEnvTS() {
				sb.WriteString("export ")
			}
			if *m.Ident == "var" {
				sb.WriteString(fmt.Sprintf("const %s = (", "_"+*m.Ident))
			} else {
				sb.WriteString(fmt.Sprintf("const %s = (", *m.Ident))
			}
			for i, p := range *m.Overloads[0] {
				if i >= m.ExpectedArgs {
					break
				}
				if i > 0 && i < len(*m.Overloads[0]) {
					sb.WriteString(", ")
				}
				sb.WriteString(*p.Ident)
				if g.conf.IsEnvTS() {
					tsType, _ := CPPTypeToTS(*p.Type)
					if v, ok := g.conf.TypeMappings[tsType]; ok {
						sb.WriteString(fmt.Sprintf(": %s", v.TSType))
					} else {
						if strings.Contains(tsType, "std::vector") {
							vectorType := tsType[strings.Index(tsType, "<")+1 : strings.Index(tsType, ">")]
							tsType, _ = CPPTypeToTS(vectorType)
							tsType = tsType + "[]"
						}
						sb.WriteString(fmt.Sprintf(": %s", tsType))
					}
				}
			}
			tsType, _ := CPPTypeToTS(*m.Returns)
			if g.conf.IsEnvTS() {
				if v, ok := g.conf.TypeMappings[tsType]; ok {
					sb.WriteString(fmt.Sprintf("): %s {\n", v.TSType))
				} else {
					sb.WriteString(fmt.Sprintf("): %s => {\n", tsType))
				}
			} else {
				sb.WriteString(") => {\n")
			}
			g.writeIndent(sb, 1)
			sb.WriteString("return ")
			if isClass(tsType, classes) {
				sb.WriteString(fmt.Sprintf("new %s(", tsType))
			}
			if *m.Ident == "var" {
				sb.WriteString(fmt.Sprintf("__%s(", *m.Ident))
			} else {
				sb.WriteString(fmt.Sprintf("_%s(", *m.Ident))
			}
			for i, p := range *m.Overloads[0] {
				if i >= m.ExpectedArgs {
					break
				}
				if i > 0 && i < len(*m.Overloads[0]) {
					sb.WriteString(", ")
				}
				sb.WriteString(*p.Ident)
				if isClass(*p.Type, classes) {
					sb.WriteString("._native_self")
				}
			}
			if isClass(tsType, classes) {
				sb.WriteByte(')')
			}
			sb.WriteString(");\n")
			sb.WriteString("}\n\n")
		}
	}

	for _, m := range processedMethods {
		if !g.conf.IsMethodIgnored(*m.Ident) {
			if !g.conf.IsMethodIgnored(*m.Ident) {
				if g.conf.IsEnvTS() {
					sb.WriteString("export ")
				}
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
					if g.conf.IsEnvTS() {
						tsType, _ := CPPTypeToTS(*p.Type)
						if v, ok := g.conf.TypeMappings[tsType]; ok {
							sb.WriteString(fmt.Sprintf(": %s", v.TSType))
						} else {
							sb.WriteString(fmt.Sprintf(": %s", tsType))
						}
					}
				}
				tsType, _ := CPPTypeToTS(*m.Returns)
				if g.conf.IsEnvTS() {
					if v, ok := g.conf.TypeMappings[tsType]; ok {
						sb.WriteString(fmt.Sprintf("): %s {\n", v.TSType))
					} else {
						sb.WriteString(fmt.Sprintf("): %s => {\n", tsType))
					}
				} else {
					sb.WriteString(") => {\n")
				}
				g.writeIndent(sb, 1)
				sb.WriteString("return ")
				if isClass(tsType, classes) {
					sb.WriteString(fmt.Sprintf("new %s(", tsType))
				}
				if *m.Ident == "var" {
					sb.WriteString(fmt.Sprintf("__%s(", *m.Ident))
				} else {
					sb.WriteString(fmt.Sprintf("_%s(", *m.Ident))
				}
				for i, p := range *m.Overloads[0] {
					if i > 0 && i < len(*m.Overloads[0]) {
						sb.WriteString(", ")
					}
					sb.WriteString(*p.Ident)
					if isClass(*p.Type, classes) {
						sb.WriteString("._native_self")
					}
				}
				if isClass(tsType, classes) {
					sb.WriteByte(')')
				}
				sb.WriteString(");\n")
				sb.WriteString("}\n\n")
			}
		}
	}

	for _, m := range g.conf.GlobalForcedMethods {
		if g.conf.IsEnvTS() {
			sb.WriteString("export ")
		}
		if m.Name == "var" {
			sb.WriteString(fmt.Sprintf("const %s = (", "_"+m.Name))
		} else {
			sb.WriteString(fmt.Sprintf("const %s = (", m.Name))
		}
		for i, p := range m.Args {
			if i > 0 && i < len(m.Args) {
				sb.WriteString(", ")
			}
			sb.WriteString(p.Name)
			if g.conf.IsEnvTS() {
				sb.WriteString(fmt.Sprintf(": %s", p.TSType))
			}
		}
		if g.conf.IsEnvTS() {
			if m.IsVoid {
				sb.WriteString(") => {\n")
			} else {
				sb.WriteString(fmt.Sprintf("): %s => {\n", m.TSReturnType))
			}
		} else {
			sb.WriteString(") => {\n")
		}
		g.writeIndent(sb, 1)
		sb.WriteString("return ")
		if isClass(m.TSReturnType, classes) {
			sb.WriteString(fmt.Sprintf("new %s(", m.TSReturnType))
		}
		if m.Name == "var" {
			sb.WriteString(fmt.Sprintf("__%s(", m.Name))
		} else {
			sb.WriteString(fmt.Sprintf("_%s(", m.Name))
		}
		for i, p := range m.Args {
			if i > 0 && i < len(m.Args) {
				sb.WriteString(", ")
			}
			sb.WriteString(p.Name)
		}
		if isClass(m.TSReturnType, classes) {
			sb.WriteByte(')')
		}
		sb.WriteString(");\n")
		sb.WriteString("}\n\n")
	}

	return sb.String()
}
