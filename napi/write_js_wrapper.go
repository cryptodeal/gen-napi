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

func (g *PackageGenerator) WriteEnvWrapper(sb *strings.Builder) {
	sb.WriteString(g.conf.JSWrapperOpts.FrontMatter)
	sb.WriteString(g.WriteEnvImports())
	sb.WriteString(g.WriteEnums())
	sb.WriteString(g.WriteEnvWrappedFns())
	if !g.conf.IsEnvTS() {
		sb.WriteString(g.WriteEnvExports())
	}
}

func (g *PackageGenerator) WriteEnvExports() string {
	sb := new(strings.Builder)
	sb.WriteString("module.exports = {\n")
	used := []string{}
	for name, m := range g.ParsedData.Methods {
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
	for name, m := range g.ParsedData.Lits {
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

	for name, c := range g.ParsedData.Classes {
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

	for i, e := range g.ParsedData.Enums {
		if i == 0 {
			sb.WriteString(",\n")
		}
		g.writeIndent(sb, 1)
		sb.WriteString(*e.Ident)
		if i < used_len-1 {
			sb.WriteString(",\n")
		}
	}
	sb.WriteString("\n}\n")
	return sb.String()
}

func (g *PackageGenerator) WriteEnvImports() string {
	hasClassImports := false
	sb := new(strings.Builder)
	sb.WriteString("const {\n")
	for name, c := range g.ParsedData.Classes {
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
	for name, m := range g.ParsedData.Methods {
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
	for name, m := range g.ParsedData.Lits {
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
	sb.WriteString("\n} = ")
	if g.conf.IsEnvTS() {
		sb.WriteString("import.meta.require(")
	} else {
		sb.WriteString("require(")
	}
	sb.WriteString(fmt.Sprintf("%q)\n\n", g.conf.ResolvedBindingsImportPath(g.conf.Path)))
	return sb.String()
}

func (g *PackageGenerator) WriteEnums() string {
	sb := new(strings.Builder)
	for _, e := range g.ParsedData.Enums {
		if g.conf.IsEnvTS() {
			sb.WriteString(fmt.Sprintf("export enum %s {\n", *e.Ident))
			count := len(e.Values)
			for i, v := range e.Values {
				g.writeIndent(sb, 1)
				sb.WriteString(fmt.Sprintf("%s = %d", *v.Ident, v.Value))
				if i < count-1 {
					sb.WriteByte(',')
				}
				sb.WriteByte('\n')
			}
			sb.WriteString("}\n\n")
		} else {
			// write as if it's a compiled typescript enum if JS out type
			sb.WriteString(fmt.Sprintf("var %s;\n", *e.Ident))
			sb.WriteString(fmt.Sprintf("(function (%s) {\n", *e.Ident))
			for _, v := range e.Values {
				g.writeIndent(sb, 1)
				sb.WriteString(fmt.Sprintf("%s[%s[%q] = %d] = %q;\n", *e.Ident, *e.Ident, *v.Ident, v.Value, *v.Ident))
			}
			sb.WriteString(fmt.Sprintf("})(%s || (%s = {}));\n\n", *e.Ident, *e.Ident))
		}
	}
	return sb.String()
}

func (g *PackageGenerator) WriteEnvWrappedFns() string {
	sb := new(strings.Builder)
	for _, m := range g.ParsedData.Methods {
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
				if i >= m.ExpectedArgs && m.OptionalArgs == 0 {
					break
				}
				if i > 0 && i < len(*m.Overloads[0]) {
					sb.WriteString(", ")
				}
				sb.WriteString(*p.Ident)
				if g.conf.IsEnvTS() {
					tsType, _ := g.CPPTypeToTS(*p.Type, p.IsPointer)
					if v, ok := g.conf.TypeMappings[tsType]; ok && v.TSType != "" {
						sb.WriteString(fmt.Sprintf(": %s", v.TSType))
					} else {
						if strings.Contains(tsType, "std::vector") {
							vectorType := tsType[strings.Index(tsType, "<")+1 : strings.Index(tsType, ">")]
							tsType, _ = g.CPPTypeToTS(vectorType, p.IsPointer)
							tsType = tsType + "[]"
							if p.DefaultValue != nil && p.DefaultValue.Val != nil {
								val := *p.DefaultValue.Val
								val = strings.ReplaceAll(val, "{", "[")
								val = strings.ReplaceAll(val, "]", "}")
								tsType += fmt.Sprintf(" = %s", val)
							}
						}
						isEnum, _ := g.IsTypeEnum(*p.Type)
						if isEnum && p.DefaultValue.Val != nil {
							tsType += fmt.Sprintf(" = %s.%s", *p.Type, *p.DefaultValue.Val)
						}
						sb.WriteString(fmt.Sprintf(": %s", tsType))
					}
				}
			}
			tsType, _ := g.CPPTypeToTS(*m.Returns, m.ReturnsPointer)
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
			if isClass(tsType, g.ParsedData.Classes) {
				sb.WriteString(fmt.Sprintf("new %s(", tsType))
			}
			if isInvalidName(*m.Ident) {
				sb.WriteString(fmt.Sprintf("__%s(", *m.Ident))
			} else {
				sb.WriteString(fmt.Sprintf("_%s(", *m.Ident))
			}
			for i, p := range *m.Overloads[0] {
				if i >= m.ExpectedArgs && m.OptionalArgs == 0 {
					break
				}
				if i > 0 && i < len(*m.Overloads[0]) {
					sb.WriteString(", ")
				}
				sb.WriteString(*p.Ident)
				if isClass(*p.Type, g.ParsedData.Classes) {
					sb.WriteString("._native_self")
				}
			}
			if isClass(tsType, g.ParsedData.Classes) {
				sb.WriteByte(')')
			}
			sb.WriteString(");\n")
			sb.WriteString("}\n\n")
		}
	}

	for name, c := range g.ParsedData.Classes {
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
							sb.WriteString("): void => {\n")
						} else if m.TSReturnType == "" {
							sb.WriteString(") => {\n")
						} else {
							sb.WriteString(fmt.Sprintf("): %s => {\n", m.TSReturnType))
						}
					} else {
						sb.WriteString(") => {\n")
					}
					g.writeIndent(sb, 1)
					sb.WriteString("return ")
					if isClass(m.TSReturnType, g.ParsedData.Classes) {
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
						if isClass(p.TSType, g.ParsedData.Classes) {
							sb.WriteString("._native_self")
						} else if isTypedArray(m.Args[i].TSType) {
							sb.WriteString(".buffer")
						}
					}
					if isClass(m.TSReturnType, g.ParsedData.Classes) {
						sb.WriteByte(')')
					}
					sb.WriteString(");\n")
					sb.WriteString("}\n\n")
				}
			}
		}
	}

	for _, m := range g.ParsedData.Lits {
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
						tsType, _ := g.CPPTypeToTS(*p.Type, p.IsPointer)
						if v, ok := g.conf.TypeMappings[tsType]; ok {
							sb.WriteString(fmt.Sprintf(": %s", v.TSType))
						} else {
							sb.WriteString(fmt.Sprintf(": %s", tsType))
						}
					}
				}
				tsType, _ := g.CPPTypeToTS(*m.Returns, m.ReturnsPointer)
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
				if isClass(tsType, g.ParsedData.Classes) {
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
					if isClass(*p.Type, g.ParsedData.Classes) {
						sb.WriteString("._native_self")
					}
				}
				if isClass(tsType, g.ParsedData.Classes) {
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
		if isClass(m.TSReturnType, g.ParsedData.Classes) {
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
		if isClass(m.TSReturnType, g.ParsedData.Classes) {
			sb.WriteByte(')')
		}
		sb.WriteString(");\n")
		sb.WriteString("}\n\n")
	}

	return sb.String()
}
