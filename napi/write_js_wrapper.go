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
		if !g.conf.IsMethodIgnored(*m.Name) {
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
		if !g.conf.IsMethodIgnored(*m.Name) {
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
		sb.WriteString(e.Name)
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
	sb.WriteString("\nconst {\n")
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
		if !g.conf.IsMethodIgnored(*m.Name) {
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
		if !g.conf.IsMethodIgnored(*m.Name) {
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
			sb.WriteString(fmt.Sprintf("export enum %s {\n", e.Name))
			count := len(e.Values)
			for i, v := range e.Values {
				g.writeIndent(sb, 1)
				sb.WriteString(fmt.Sprintf("%s = %d", *v.Name, v.Value))
				if i < count-1 {
					sb.WriteByte(',')
				}
				sb.WriteByte('\n')
			}
			sb.WriteString("}\n\n")
		} else {
			// write as if it's a compiled typescript enum if JS out type
			sb.WriteString(fmt.Sprintf("var %s;\n", e.Name))
			sb.WriteString(fmt.Sprintf("(function (%s) {\n", e.Name))
			for _, v := range e.Values {
				g.writeIndent(sb, 1)
				sb.WriteString(fmt.Sprintf("%s[%s[%q] = %d] = %q;\n", e.Name, e.Name, *v.Name, v.Value, *v.Name))
			}
			sb.WriteString(fmt.Sprintf("})(%s || (%s = {}));\n\n", e.Name, e.Name))
		}
	}
	return sb.String()
}

func (g *PackageGenerator) WriteDefaultArgVal(sb *strings.Builder, p *CPPArg) {
	if p.DefaultValue != nil && p.DefaultValue.Val != nil {
		val := *p.DefaultValue.Val
		val = strings.ReplaceAll(val, "{", "[")
		val = strings.ReplaceAll(val, "}", "]")
		sb.WriteString(fmt.Sprintf(" = %s", val))
	}
}

func (g *PackageGenerator) WriteEnvWrappedFns() string {
	sb := new(strings.Builder)
	for _, m := range g.ParsedData.Methods {
		if !g.conf.IsMethodIgnored(*m.Name) {
			if g.conf.IsEnvTS() {
				sb.WriteString("export ")
			}
			if isInvalidName(*m.Name) {
				sb.WriteString(fmt.Sprintf("const %s = (", "_"+*m.Name))
			} else {
				sb.WriteString(fmt.Sprintf("const %s = (", *m.Name))
			}
			for i, p := range *m.Overloads[0] {
				if i >= m.ExpectedArgs && m.OptionalArgs == 0 {
					break
				}
				if i > 0 && i < len(*m.Overloads[0]) {
					sb.WriteString(", ")
				}
				sb.WriteString(*p.Name)
				usedName := p.Type.Name
				if v, ok := g.conf.TypeMappings[p.Type.Name]; (!ok || v.TSType != "") && p.Type.MappedType != nil {
					usedName = p.Type.MappedType.Name
				}
				tsType, isClass := g.CPPTypeToTS(usedName, p.Type.IsPointer)
				if v, ok := g.conf.TypeMappings[tsType]; ok && v.TSType != "" {
					if g.conf.IsEnvTS() {
						sb.WriteString(fmt.Sprintf(": %s", stripNameSpace(v.TSType)))
					}
					g.WriteDefaultArgVal(sb, p)
				} else if p.Type.MappedType != nil {
					type_helpers := p.Type.MappedType.GetTypeHandlers(g, true)
					if g.conf.IsEnvTS() {
						sb.WriteString(fmt.Sprintf(": %s", type_helpers.JSType))
					}
					g.WriteDefaultArgVal(sb, p)
				} else {
					// TODO: write types for `pair` aka `[T1, T2]`
					if IsArgTemplate(p) && *p.Type.Template.Name == "vector" {
						tsType, _ = g.CPPTypeToTS(*p.Type.Template.Args[0].Name, p.Type.IsPointer)
						if g.conf.IsEnvTS() {
							if tsType == "pair" {
								helpers_1 := g.GetTypeHelpers(*p.Type.Template.Args[0].Args[0].Name)
								helpers_2 := g.GetTypeHelpers(*p.Type.Template.Args[0].Args[1].Name)
								sb.WriteString(fmt.Sprintf(": Array<[%s, %s]>", helpers_1.JSType, helpers_2.JSType))
							} else {
								sb.WriteString(fmt.Sprintf(": %s[]", tsType))
							}
						}
						g.WriteDefaultArgVal(sb, p)
					} else {
						sb.WriteString(fmt.Sprintf(": %s", stripNameSpace(tsType)))
					}
					isEnum, _ := g.IsTypeEnum(p.Type.Name)
					if isEnum && p.DefaultValue != nil && p.DefaultValue.Val != nil {
						sb.WriteString(fmt.Sprintf(" = %s.%s", p.Type.Name, *p.DefaultValue.Val))
					} else if !isClass && p.DefaultValue != nil && p.DefaultValue.Val != nil {
						val := *p.DefaultValue.Val
						val = strings.ReplaceAll(val, "{", "[")
						val = strings.ReplaceAll(val, "}", "]")
						sb.WriteString(fmt.Sprintf(" = %s", val))
					}
				}
			}
			tsType, _ := g.CPPTypeToTS(m.Returns.Name, m.Returns.IsPointer)
			if g.conf.IsEnvTS() {
				if v, ok := g.conf.TypeMappings[tsType]; ok {
					sb.WriteString(fmt.Sprintf("): %s => {\n", v.TSType))
				} else {
					sb.WriteString(fmt.Sprintf("): %s => {\n", tsType))
				}
			} else {
				sb.WriteString(") => {\n")
			}
			g.writeIndent(sb, 1)
			sb.WriteString("return ")
			if g.isClass(tsType) {
				sb.WriteString(fmt.Sprintf("new %s(", stripNameSpace(tsType)))
			}
			if isInvalidName(*m.Name) {
				sb.WriteString(fmt.Sprintf("__%s(", *m.Name))
			} else {
				sb.WriteString(fmt.Sprintf("_%s(", *m.Name))
			}
			for i, p := range *m.Overloads[0] {
				if i >= m.ExpectedArgs && m.OptionalArgs == 0 {
					break
				}
				if i > 0 && i < len(*m.Overloads[0]) {
					sb.WriteString(", ")
				}
				sb.WriteString(*p.Name)
				if g.isClass(p.Type.Name) {
					sb.WriteString("._native_self")
				}
			}
			if g.isClass(tsType) {
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
					if g.isClass(m.TSReturnType) {
						sb.WriteString(fmt.Sprintf("new %s(", stripNameSpace(m.TSReturnType)))
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
						if g.isClass(p.TSType) {
							sb.WriteString("._native_self")
						} else if isTypedArray(m.Args[i].TSType) {
							sb.WriteString(".buffer")
						}
					}
					if g.isClass(m.TSReturnType) {
						sb.WriteByte(')')
					}
					sb.WriteString(");\n")
					sb.WriteString("}\n\n")
				}
			}
		}
	}

	for _, m := range g.ParsedData.Lits {
		if !g.conf.IsMethodIgnored(*m.Name) {
			if !g.conf.IsMethodIgnored(*m.Name) {
				if g.conf.IsEnvTS() {
					sb.WriteString("export ")
				}
				if isInvalidName(*m.Name) {
					sb.WriteString(fmt.Sprintf("const %s = (", "_"+*m.Name))
				} else {
					sb.WriteString(fmt.Sprintf("const %s = (", *m.Name))
				}
				for i, p := range *m.Overloads[0] {
					if i > 0 && i < len(*m.Overloads[0]) {
						sb.WriteString(", ")
					}
					sb.WriteString(*p.Name)
					if g.conf.IsEnvTS() {
						tsType, _ := g.CPPTypeToTS(p.Type.Name, p.Type.IsPointer)
						if v, ok := g.conf.TypeMappings[tsType]; ok {
							sb.WriteString(fmt.Sprintf(": %s", v.TSType))
						} else {
							sb.WriteString(fmt.Sprintf(": %s", tsType))
						}
					}
				}
				tsType, _ := g.CPPTypeToTS(m.Returns.Name, m.Returns.IsPointer)
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
				if g.isClass(tsType) {
					sb.WriteString(fmt.Sprintf("new %s(", stripNameSpace(tsType)))
				}
				if isInvalidName(*m.Name) {
					sb.WriteString(fmt.Sprintf("__%s(", *m.Name))
				} else {
					sb.WriteString(fmt.Sprintf("_%s(", *m.Name))
				}
				for i, p := range *m.Overloads[0] {
					if i > 0 && i < len(*m.Overloads[0]) {
						sb.WriteString(", ")
					}
					sb.WriteString(*p.Name)
					if g.isClass(p.Type.Name) {
						sb.WriteString("._native_self")
					}
				}
				if g.isClass(tsType) {
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
		if g.isClass(m.TSReturnType) {
			sb.WriteString(fmt.Sprintf("new %s(", stripNameSpace(m.TSReturnType)))
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
		if g.isClass(m.TSReturnType) {
			sb.WriteByte(')')
		}
		sb.WriteString(");\n")
		sb.WriteString("}\n\n")
	}

	return sb.String()
}
