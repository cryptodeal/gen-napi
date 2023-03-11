package napi

import (
	"fmt"
	"os"
	"path/filepath"
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

func (g *PackageGenerator) WriteClassImports(sb *strings.Builder) {
	for name, c := range g.conf.ClassOpts {
		if g.conf.IsEnvTS() {
			sb.WriteString("import ")
		} else {
			sb.WriteString("const ")
		}
		sb.WriteString(fmt.Sprintf(" { %s } ", name))
		if g.conf.IsEnvTS() {
			sb.WriteString("from ")
		} else {
			sb.WriteString("= require(")
		}
		sb.WriteString(fmt.Sprintf("'%s'", g.conf.ResolvedImportPath(c.PathToImpl)))
		if !g.conf.IsEnvTS() {
			sb.WriteString(")")
		}
		sb.WriteString(";\n")
	}
}

func (g *PackageGenerator) WriteEnumImports(sb *strings.Builder, imports map[string]bool) {
	if len(g.ParsedData.Enums) == 0 {
		return
	}

	if g.conf.IsEnvTS() {
		sb.WriteString("import ")
	} else {
		sb.WriteString("const ")
	}
	sb.WriteString("{ ")
	count := 0
	import_count := len(imports)
	for name := range imports {
		if count > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(name)
		count++
		if count == import_count {
			sb.WriteByte(' ')
		}
	}
	sb.WriteString("} ")
	if g.conf.IsEnvTS() {
		sb.WriteString("from ")
	} else {
		sb.WriteString("= require(")
	}
	base_import_path := strings.Split(g.conf.ResolvedWrappedEnumOutPath(filepath.Dir(g.conf.Path)), "/")
	used_import_path := base_import_path[len(base_import_path)-1]
	if g.conf.IsEnvTS() {
		used_import_path = used_import_path[:strings.IndexByte(used_import_path, '.')]
	}
	sb.WriteString(fmt.Sprintf("'./%s'", used_import_path))

	if !g.conf.IsEnvTS() {
		sb.WriteString(")")
	}
	sb.WriteString(";\n")

}

func (g *PackageGenerator) WriteEnvWrapper(sb *strings.Builder) {
	temp_sb := new(strings.Builder)
	temp_sb.WriteString(g.conf.JSWrapperOpts.FrontMatter)
	g.WriteClassImports(temp_sb)
	g.WriteEnvImports(temp_sb)
	enum_imports := map[string]bool{}
	g.WriteEnvWrappedFns(temp_sb, enum_imports)
	if !g.conf.IsEnvTS() {
		g.WriteEnvExports(temp_sb)
	}
	g.WriteEnumImports(sb, enum_imports)
	sb.WriteString(temp_sb.String())
}

func (g *PackageGenerator) WriteEnvExports(sb *strings.Builder) {
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
}

func (g *PackageGenerator) WriteEnvImports(sb *strings.Builder) {
	sb.WriteString("const addon = ")
	if g.conf.IsEnvTS() {
		sb.WriteString("import.meta.require(")
	} else {
		sb.WriteString("require(")
	}
	sb.WriteString(fmt.Sprintf("'%s');\n\n", g.conf.ResolvedImportPath(g.conf.JSWrapperOpts.AddonPath)))
}

func (g *PackageGenerator) WriteEnums() error {
	sb := new(strings.Builder)
	g.writeFileCodegenHeader(sb)
	g.writeFileSourceHeader(sb, *g.Path)
	enum_count := len(g.ParsedData.Enums)
	for i, e := range g.ParsedData.Enums {
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
			sb.WriteString("}\n")
		} else {
			// write as if it's a compiled typescript enum if JS out type
			sb.WriteString(fmt.Sprintf("var %s;\n", e.Name))
			sb.WriteString(fmt.Sprintf("(function (%s) {\n", e.Name))
			for _, v := range e.Values {
				g.writeIndent(sb, 1)
				sb.WriteString(fmt.Sprintf("%s[%s[%q] = %d] = %q;\n", e.Name, e.Name, *v.Name, v.Value, *v.Name))
			}
			sb.WriteString(fmt.Sprintf("})(%s || (%s = {}));\n", e.Name, e.Name))
		}
		if i < enum_count-1 {
			sb.WriteByte('\n')
		}
	}
	outPath := g.conf.ResolvedWrappedEnumOutPath(filepath.Dir(g.conf.Path))
	err := os.MkdirAll(filepath.Dir(outPath), os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(outPath, []byte(sb.String()), os.ModePerm)
	if err != nil {
		panic(err)
	}

	return nil
}

func (g *PackageGenerator) WriteDefaultArgVal(sb *strings.Builder, p *CPPArg) {
	if p.DefaultValue != nil && p.DefaultValue.Val != nil {
		val := *p.DefaultValue.Val
		val = strings.ReplaceAll(val, "{", "[")
		val = strings.ReplaceAll(val, "}", "]")
		sb.WriteString(fmt.Sprintf(" = %s", val))
	}
}

func IsBigIntTypedArray(ts_type string) bool {
	return strings.Contains(ts_type, "BigInt64Array") || strings.Contains(ts_type, "BigUint64Array")
}

func (g *PackageGenerator) WriteWrappedFn(sb *strings.Builder, method_name string, args []GenArgData, returns GenReturnData, is_void bool) {
	name := method_name
	if isInvalidName(method_name) {
		name = fmt.Sprintf("_%s", method_name)
	}

	if g.conf.IsEnvTS() {
		sb.WriteString("export ")
	}

	sb.WriteString(fmt.Sprintf("const %s = (", name))

	arg_count := len(args)
	for i, arg := range args {
		if i > 0 && i < arg_count {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%s: %s", arg.Name, arg.JSType))
		if arg.DefaultValue != nil {
			sb.WriteString(fmt.Sprintf(" = %s", *arg.DefaultValue))
		}
	}
	sb.WriteString(")")
	if g.conf.IsEnvTS() && !is_void {
		sb.WriteString(fmt.Sprintf(": %s", stripNameSpace(returns.JSType)))
	}
	sb.WriteString(" => {\n")
	g.writeIndent(sb, 1)
	sb.WriteString("return ")
	if returns.NapiType == External {
		sb.WriteString(fmt.Sprintf("new %s(", stripNameSpace(returns.JSType)))
	}
	sb.WriteString(fmt.Sprintf("addon._%s(", method_name))

	for i, arg := range args {
		if i > 0 && i < arg_count {
			sb.WriteString(", ")
		}
		sb.WriteString(arg.Name)
		if arg.TypedArrayInfo != nil {
			instanceof_name := arg.NapiType.TypedArrayType()
			sb.WriteString(fmt.Sprintf(" instanceof %s ? %s : new %s(%s", instanceof_name, arg.Name, instanceof_name, arg.Name))
			// if bigint, need to convert to bigint
			if IsBigIntTypedArray(arg.JSType) {
				sb.WriteString(".map((v) => typeof v === 'number' ? BigInt(v) : v)")
			}
			sb.WriteString(")")
		} else if arg.NapiType == External {
			sb.WriteString("._native_ref")
		}
	}
	if returns.NapiType == External {
		sb.WriteByte(')')
	}
	sb.WriteString(");\n")
	sb.WriteString("}\n\n")
}

func (g *PackageGenerator) WriteExternalTypeHelpers(sb *strings.Builder, type_name string) {
	sb.WriteString(fmt.Sprintf("export type _Native_%s = unknown & Record<string, never>;\n\n", type_name))
	sb.WriteString(fmt.Sprintf("export abstract class _Base_%s {\n", type_name))
	g.writeIndent(sb, 1)
	sb.WriteString(fmt.Sprintf("protected _native_%s: _Native_%s;\n\n", type_name, type_name))
	g.writeIndent(sb, 1)
	sb.WriteString(fmt.Sprintf("get _native_ref(): _Native_%s {\n", type_name))
	g.writeIndent(sb, 2)
	sb.WriteString(fmt.Sprintf("return this._native_%s;\n", type_name))
	g.writeIndent(sb, 1)
	sb.WriteString("}\n")
	sb.WriteString("}\n\n")
}

func (g *PackageGenerator) WriteForcedMethod(sb *strings.Builder, method_name string, args []FnArg, returns string, is_void bool) {
	name := method_name
	if isInvalidName(method_name) {
		name = fmt.Sprintf("_%s", method_name)
	}

	if g.conf.IsEnvTS() {
		sb.WriteString("export ")
	}

	sb.WriteString(fmt.Sprintf("const %s = (", name))

	arg_count := len(args)
	for i, arg := range args {
		if i > 0 && i < arg_count {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%s: %s", arg.Name, arg.TSType))
		if arg.Default != "" {
			sb.WriteString(fmt.Sprintf(" = %s", arg.Default))
		}
	}
	sb.WriteString(")")
	if g.conf.IsEnvTS() && !is_void && returns != "" {
		sb.WriteString(fmt.Sprintf(": %s", returns))
	}
	sb.WriteString(" => {\n")
	g.writeIndent(sb, 1)
	sb.WriteString("return ")
	if g.isClass(returns) {
		sb.WriteString(fmt.Sprintf("new %s(", returns))
	}
	sb.WriteString(fmt.Sprintf("addon._%s(", method_name))

	for i, arg := range args {
		if i > 0 && i < arg_count {
			sb.WriteString(", ")
		}
		sb.WriteString(arg.Name)
		if g.isClass(arg.TSType) {
			sb.WriteString("._native_ref")
		}
	}
	if g.isClass(returns) {
		sb.WriteByte(')')
	}
	sb.WriteString(");\n")
	sb.WriteString("}\n\n")
}

func ParseEnumImports(arg_helpers []GenArgData, imports map[string]bool) {
	for _, arg := range arg_helpers {
		if arg.NapiType == NumberEnum {
			imports[arg.JSType] = true
		}
	}
}

func (g *PackageGenerator) WriteEnvWrappedFns(sb *strings.Builder, imports map[string]bool) {
	for _, m := range g.ParsedData.Methods {
		if !g.conf.IsMethodIgnored(*m.Name) {
			arg_helpers := g.GetArgData(m.Overloads[0])
			ParseEnumImports(*arg_helpers, imports)
			g.WriteWrappedFn(sb, *m.Name, *arg_helpers, m.Returns.ParseReturnData(g), m.IsVoid())
		}
	}

	for _, m := range g.ParsedData.Lits {
		if !g.conf.IsMethodIgnored(*m.Name) {
			arg_helpers := g.GetArgData(m.Overloads[0])
			ParseEnumImports(*arg_helpers, imports)
			g.WriteWrappedFn(sb, *m.Name, *arg_helpers, m.Returns.ParseReturnData(g), m.IsVoid())
		}
	}

	for name, c := range g.ParsedData.Classes {
		if c.Decl != nil {
			if c.FieldDecl != nil {
				for _, f := range *c.FieldDecl {
					if f.Name != nil && !g.conf.IsFieldIgnored(name, *f.Name) {
						arg_helpers := g.GetArgData(&f.Args)
						ParseEnumImports(*arg_helpers, imports)
						g.WriteWrappedFn(sb, *f.Name, *arg_helpers, f.Returns.ParseReturnData(g), f.IsVoid())
					}
				}
			}

			if v, ok := g.conf.ClassOpts[name]; ok && len(v.ForcedMethods) > 0 {
				for _, m := range v.ForcedMethods {
					g.WriteForcedMethod(sb, m.Name, m.Args, m.TSReturnType, m.IsVoid)
				}
			}
		}
	}

	for _, m := range g.conf.GlobalForcedMethods {
		for _, arg := range m.Args {
			isEnum, _ := g.IsTypeEnum(arg.TSType)
			if isEnum {
				imports[arg.TSType] = true
			}
		}
		g.WriteForcedMethod(sb, m.Name, m.Args, m.TSReturnType, m.IsVoid)
	}
}

func (g *PackageGenerator) WriteShimWrappedFn(sb *strings.Builder, method_name string, class_name string, args []GenArgData, returns GenReturnData, is_void bool, imports map[string]bool) {
	g.writeIndent(sb, 2)
	sb.WriteString(fmt.Sprintf("%s(", method_name))
	for i, arg := range args {
		if i == 0 {
			continue
		}
		if i > 1 {
			sb.WriteString(", ")
		}
		if arg.NapiType == NumberEnum {
			imports[stripNameSpace(arg.JSType)] = true
		}
		sb.WriteString(fmt.Sprintf("%s: %s", arg.Name, arg.JSType))
		if arg.DefaultValue != nil {
			sb.WriteString(fmt.Sprintf(" = %s", *arg.DefaultValue))
		}
	}
	sb.WriteString(")")
	if g.conf.IsEnvTS() && !is_void {
		sb.WriteString(fmt.Sprintf(": %s", stripNameSpace(returns.JSType)))
	}
	sb.WriteString(" {\n")
	g.writeIndent(sb, 3)
	sb.WriteString("return ")
	if g.isClass(returns.RawType.Name) {
		sb.WriteString(fmt.Sprintf("new %s(", class_name))
	}
	sb.WriteString(fmt.Sprintf("addon._%s(", method_name))

	for i, arg := range args {
		if i == 0 {
			sb.WriteString("this._native_ref")
			continue
		}
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(arg.Name)
		if arg.TypedArrayInfo != nil {
			instanceof_name := arg.NapiType.TypedArrayType()
			sb.WriteString(fmt.Sprintf(" instanceof %s ? %s : new %s(%s", instanceof_name, arg.Name, instanceof_name, arg.Name))
			// if bigint, need to convert to bigint
			if IsBigIntTypedArray(arg.JSType) {
				sb.WriteString(".map((v) => typeof v === 'number' ? BigInt(v) : v)")
			}
			sb.WriteString(")")
		} else if arg.NapiType == External {
			sb.WriteString("._native_ref")
		}
	}
	if g.isClass(returns.RawType.Name) {
		sb.WriteByte(')')
	}
	sb.WriteString(");\n")
	g.writeIndent(sb, 2)
	sb.WriteString("},\n\n")
}

func (g *PackageGenerator) WriteClassInstructions(sb *strings.Builder, var_name string) {
	sb.WriteString("/**\n")
	sb.WriteString(" * USAGE INSTRUCTIONS:\n")
	sb.WriteString(fmt.Sprintf(" * 1. Import `gen_%s_ops_shim` to file where corresponding js impl lives\n", var_name))
	sb.WriteString(fmt.Sprintf(" * 2. copy the following logic following class declaration into `%s`:\n", g.conf.ClassOpts[var_name].PathToImpl))
	sb.WriteString(" */\n")
	sb.WriteString("/*\n")
	if g.conf.IsEnvTS() {
		g.writeIndent(sb, 1)
		sb.WriteString(fmt.Sprintf("export interface %s extends ReturnType<typeof gen_%s_ops_shim> {} // eslint-disable-line\n", var_name, var_name))
	}
	g.writeIndent(sb, 1)
	sb.WriteString(fmt.Sprintf("for (const [method, closure] of Object.entries(gen_%s_ops_shim(%s))) {\n", var_name, var_name))
	g.writeIndent(sb, 2)
	sb.WriteString(fmt.Sprintf("%s.prototype[method] = closure;\n", var_name))
	g.writeIndent(sb, 1)
	sb.WriteString("}\n")
	sb.WriteString("*/\n")
}

func (g *PackageGenerator) WriteClassShims(name string, c *CPPClass) error {
	if c.Decl != nil && c.FieldDecl != nil {
		shimmed_methods := g.WriteObjectShims(name)
		outPath := g.conf.ResolvedShimPath(filepath.Dir(g.conf.Path), name)
		err := os.MkdirAll(filepath.Dir(outPath), os.ModePerm)
		if err != nil {
			return nil
		}

		err = os.WriteFile(outPath, []byte(shimmed_methods), os.ModePerm)
		if err != nil {
			return nil
		}
	}
	return nil
}

func (g *PackageGenerator) WriteObjectShims(name string) string {
	sb := new(strings.Builder)
	imports := map[string]bool{}
	if v, ok := g.ParsedData.Classes[name]; ok && v.FieldDecl != nil {
		usedName := fmt.Sprintf("_%s", name)
		sb.WriteString(fmt.Sprintf("import type { %s } from '%s';\n", name, g.conf.ResolvedImportPath(g.conf.ClassOpts[name].PathToImpl)))
		g.WriteEnvImports(sb)
		g.writeFileSourceHeader(sb, *g.Path)

		g.WriteExternalTypeHelpers(sb, name)

		g.WriteClassInstructions(sb, name)
		// write fn to shim the class/struct methods
		if g.conf.IsEnvTS() {
			sb.WriteString("\nexport ")
		}
		sb.WriteString(fmt.Sprintf("const gen_%s_ops_shim = (%s", name, usedName))

		if g.conf.IsEnvTS() {
			sb.WriteString(fmt.Sprintf(": new (...args: unknown[]) => %s", name))
		}
		sb.WriteString(") => {\n")
		g.writeIndent(sb, 1)
		sb.WriteString("return {\n")
		for _, m := range *v.FieldDecl {
			if m.Name != nil && !g.conf.IsFieldIgnored(name, *m.Name) {
				arg_helpers := g.GetArgData(&m.Args)
				res_helpers := m.Returns.ParseReturnData(g)
				if res_helpers.NapiType == NumberEnum {
					imports[stripNameSpace(res_helpers.JSType)] = true
				}
				g.WriteShimWrappedFn(sb, *m.Name, usedName, *arg_helpers, res_helpers, m.IsVoid(), imports)
			}
		}

		for _, m := range g.ParsedData.Lits {
			if m.Name != nil && !g.conf.IsFieldIgnored(name, *m.Name) && stripNameSpace(m.Returns.Name) == name {
				arg_helpers := g.GetArgData(m.Overloads[0])
				res_helpers := m.Returns.ParseReturnData(g)
				if res_helpers.NapiType == NumberEnum {
					imports[stripNameSpace(res_helpers.JSType)] = true
				}
				g.WriteShimWrappedFn(sb, *m.Name, usedName, *arg_helpers, res_helpers, m.IsVoid(), imports)
			}
		}

		for _, m := range g.ParsedData.Methods {
			if m.Name != nil && !g.conf.IsFieldIgnored(name, *m.Name) && stripNameSpace(m.Returns.Name) == name {
				arg_helpers := *g.GetArgData(m.Overloads[0])
				if arg_helpers[0].RawType.Name == name {
					res_helpers := m.Returns.ParseReturnData(g)
					if res_helpers.NapiType == NumberEnum {
						imports[stripNameSpace(res_helpers.JSType)] = true
					}
					g.WriteShimWrappedFn(sb, *m.Name, usedName, arg_helpers, res_helpers, m.IsVoid(), imports)
				}
			}
		}
		g.writeIndent(sb, 1)
		sb.WriteString("}\n")
		sb.WriteString("}\n\n")
	}

	res_sb := new(strings.Builder)
	g.writeFileCodegenHeader(res_sb)
	g.WriteEnumImports(res_sb, imports)
	res_sb.WriteString(sb.String())

	return res_sb.String()
}
