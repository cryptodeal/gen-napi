package napi

import (
	"fmt"
	"strings"
)

func isInvalidName(name string) bool {
	invalid := []string{"var", "eval"}
	for _, n := range invalid {
		if strings.EqualFold(name, n) {
			return true
		}
	}
	return false
}

func isTypedArray(typeName string) bool {
	ta := []string{"Int8Array", "Uint8Array", "Uint8ClampedArray", "Int16Array", "Uint16Array", "Int32Array", "Uint32Array", "Float32Array", "Float64Array", "BigInt64Array", "BigUint64Array"}
	for _, t := range ta {
		if strings.EqualFold(typeName, t) {
			return true
		}
	}
	return false
}

func (g *PackageGenerator) WriteEnvWrapper(sb *strings.Builder, classes map[string]*CPPClass, methods map[string]*CPPMethod, processedMethods map[string]*CPPMethod) {
	sb.WriteString(g.conf.JSWrapperOpts.FrontMatter)
	sb.WriteString(g.WriteEnvImports(classes, methods, processedMethods))
	sb.WriteString(g.WriteEnvWrappedFns(methods, processedMethods, classes))
	if !g.conf.IsEnvTS() {
		sb.WriteString(g.WriteEnvExports(classes, methods, processedMethods))
	}
}

func (g *PackageGenerator) WriteEnvExports(classes map[string]*CPPClass, methods map[string]*CPPMethod, processedMethods map[string]*CPPMethod) string {
	sb := new(strings.Builder)
	sb.WriteString("module.exports = {\n")
	used := []string{}
	for name, m := range methods {
		if !g.conf.IsMethodIgnored(*m.Ident) {
			if isInvalidName(name) {
				used = append(used, "_"+name)
			} else {
				used = append(used, name)
			}
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
	for name, m := range processedMethods {
		if !g.conf.IsMethodIgnored(*m.Ident) {
			if isInvalidName(name) {
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
		if isInvalidName(m.Name) {
			sb.WriteString(fmt.Sprintf("_%s", m.Name))
		} else {
			sb.WriteString(m.Name)
		}
		if i < used_len-1 {
			sb.WriteString(",\n")
		}
	}

	for name, c := range classes {
		if c.Decl != nil {
			if v, ok := g.conf.ClassOpts[name]; ok && len(v.ForcedMethods) > 0 {
				used_len = len(v.ForcedMethods)
				for i, m := range v.ForcedMethods {
					if i == 0 {
						sb.WriteString(",\n")
					}
					g.writeIndent(sb, 1)
					if isInvalidName(m.Name) {
						sb.WriteString(fmt.Sprintf("_%s", m.Name))
					} else {
						sb.WriteString(m.Name)
					}
					if i < used_len-1 {
						sb.WriteString(",\n")
					}
				}
			}
		}
	}
	sb.WriteString("\n}\n")
	return sb.String()
}

func (g *PackageGenerator) WriteEnvImports(classes map[string]*CPPClass, methods map[string]*CPPMethod, processedMethods map[string]*CPPMethod) string {
	hasClassImports := false
	sb := new(strings.Builder)
	sb.WriteString("const {\n")
	for name, c := range classes {
		if c.Decl != nil {
			if v, ok := g.conf.ClassOpts[name]; ok && len(v.ForcedMethods) > 0 {
				for i, m := range v.ForcedMethods {
					hasClassImports = true
					if i > 0 {
						sb.WriteString(",\n")
					}
					g.writeIndent(sb, 1)
					if isInvalidName(m.Name) {
						sb.WriteString(fmt.Sprintf("_%s: __%s", m.Name, m.Name))
					} else {
						sb.WriteString(fmt.Sprintf("_%s", m.Name))
					}
				}
			}
		}
	}
	used := []string{}
	for name, m := range methods {
		if !g.conf.IsMethodIgnored(*m.Ident) {
			used = append(used, name)
		}
	}

	used_len := len(used)
	for i, name := range used {
		if i == 0 && hasClassImports {
			sb.WriteString(",\n")
		}
		g.writeIndent(sb, 1)
		if isInvalidName(name) {
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
		if isInvalidName(name) {
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
		if isInvalidName(m.Name) {
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

func (g *PackageGenerator) WriteEnvWrappedFns(methods map[string]*CPPMethod, processedMethods map[string]*CPPMethod, classes map[string]*CPPClass) string {
	sb := new(strings.Builder)
	for _, m := range methods {
		if !g.conf.IsMethodIgnored(*m.Ident) {
			if g.conf.IsEnvTS() {
				sb.WriteString("export ")
			}
			if isInvalidName(*m.Ident) {
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
			if isInvalidName(*m.Ident) {
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

	for name, c := range classes {
		if c.Decl != nil {
			if v, ok := g.conf.ClassOpts[name]; ok && len(v.ForcedMethods) > 0 {
				for _, m := range v.ForcedMethods {
					if g.conf.IsEnvTS() {
						sb.WriteString("export ")
					}
					if isInvalidName(m.Name) {
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
					if isInvalidName(m.Name) {
						sb.WriteString(fmt.Sprintf("__%s(", m.Name))
					} else {
						sb.WriteString(fmt.Sprintf("_%s(", m.Name))
					}
					for i, p := range m.Args {
						if i > 0 && i < len(m.Args) {
							sb.WriteString(", ")
						}
						sb.WriteString(p.Name)
						if isClass(p.TSType, classes) {
							sb.WriteString("._native_self")
						} else if isTypedArray(m.Args[i].TSType) {
							sb.WriteString(".buffer")
						}
					}
					if isClass(m.TSReturnType, classes) {
						sb.WriteByte(')')
					}
					sb.WriteString(");\n")
					sb.WriteString("}\n\n")
				}
			}
		}
	}

	for _, m := range processedMethods {
		if !g.conf.IsMethodIgnored(*m.Ident) {
			if !g.conf.IsMethodIgnored(*m.Ident) {
				if g.conf.IsEnvTS() {
					sb.WriteString("export ")
				}
				if isInvalidName(*m.Ident) {
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
				if isInvalidName(*m.Ident) {
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
		if isInvalidName(m.Name) {
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
		if isInvalidName(m.Name) {
			sb.WriteString(fmt.Sprintf("__%s(", m.Name))
		} else {
			sb.WriteString(fmt.Sprintf("_%s(", m.Name))
		}
		for i, p := range m.Args {
			if i > 0 && i < len(m.Args) {
				sb.WriteString(", ")
			}
			sb.WriteString(p.Name)
			if isTypedArray(m.Args[i].TSType) {
				sb.WriteString(".buffer")
			}
		}
		if isClass(m.TSReturnType, classes) {
			sb.WriteByte(')')
		}
		sb.WriteString(");\n")
		sb.WriteString("}\n\n")
	}

	return sb.String()
}
