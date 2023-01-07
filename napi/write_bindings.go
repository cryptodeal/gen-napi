package napi

import (
	"fmt"
	"strings"
)

func (g *PackageGenerator) writeArgChecker(sb *strings.Builder, name string, checker string, idx int, msg string) {
	g.writeIndent(sb, 1)
	sb.WriteString(fmt.Sprintf("if (!info[%d].%s()) {\n", idx, checker))
	g.writeIndent(sb, 2)
	sb.WriteString(fmt.Sprintf("Napi::TypeError::New(env, %q).ThrowAsJavaScriptException();\n", fmt.Sprintf("`%s` expects args[%d] to be %s", name, idx, msg)))
	g.writeIndent(sb, 2)
	sb.WriteString("return env.Null();\n")
	g.writeIndent(sb, 1)
	sb.WriteString("}\n")

}

func (g *PackageGenerator) writeArgChecks(sb *strings.Builder, name string, args *[]*CPPArg, expected_arg_count int, classes map[string]*CPPClass) {
	if expected_arg_count == 0 {
		return
	}

	g.writeIndent(sb, 1)
	sb.WriteString(fmt.Sprintf("if (info.Length() != %d) {\n", expected_arg_count))
	g.writeIndent(sb, 2)
	errMsg := fmt.Sprintf("`%s` expects exactly %d arg", name, expected_arg_count)
	if expected_arg_count > 1 {
		errMsg += "s"
	}
	sb.WriteString(fmt.Sprintf("Napi::TypeError::New(env, %q).ThrowAsJavaScriptException();\n", errMsg))
	g.writeIndent(sb, 2)
	sb.WriteString("return env.Null();\n")
	g.writeIndent(sb, 1)
	sb.WriteString("}\n")

	// write arg checks, transforms, and declare arg variable as possible
	if v, ok := g.conf.MethodTransforms[name]; ok && v.ArgCheckTransforms != "" {
		sb.WriteString(v.ArgCheckTransforms)
		for i, arg := range *args {
			if i > expected_arg_count {
				break
			}
			// write arg transform only if it doesn't rely on other args
			if v2, ok2 := v.ArgTransforms[*arg.Ident]; ok2 && !strings.Contains(v2, "/arg_") {
				g.writeIndent(sb, 1)
				// `/arg/` is template for argument's index in the Napi callback info
				parsedTransform := strings.ReplaceAll(v2, "/arg/", fmt.Sprintf("info[%d]", i))
				sb.WriteString(parsedTransform)
			}
		}
		return
	}

	// TODO: clean up logic/handle more cases
	for i, arg := range *args {
		if i > expected_arg_count {
			break
		}
		if arg.IsPrimitive {
			napiTypeHandler := ""
			jsTypeEquivalent := ""
			valGetter := "Value"
			var needsCast *string
			switch *arg.Type {
			case "float":
				napiTypeHandler = "IsNumber"
				jsTypeEquivalent = "number"
				valGetter = "FloatValue"

			case "double":
				napiTypeHandler = "IsNumber"
				jsTypeEquivalent = "number"
				valGetter = "DoubleValue"

			case "long long", "char", "signed", "int8_t", "int32_t", "int16_t", "short":
				napiTypeHandler = "IsNumber"
				jsTypeEquivalent = "number"
				valGetter = "Int32Value"
				needsCast = arg.Type

			case "int", "int64_t":
				napiTypeHandler = "IsNumber"
				jsTypeEquivalent = "number"
				valGetter = "Int64Value"
				needsCast = arg.Type

			case "unsigned long long", "unsigned char", "unsigned", "uint8_t", "uint16_t", "unsigned short", "uint32_t":
				napiTypeHandler = "IsNumber"
				jsTypeEquivalent = "number"
				valGetter = "Uint32Value"
				needsCast = arg.Type

			case "unsigned int", "uint64_t", "size_t", "uintptr_t":
				napiTypeHandler = "IsNumber"
				jsTypeEquivalent = "number"
				valGetter = "Uint64Value"
				needsCast = arg.Type
			case "bool":
				napiTypeHandler = "IsBoolean"
				jsTypeEquivalent = "boolean"
			}
			g.writeArgChecker(sb, name, napiTypeHandler, i, fmt.Sprintf("typeof `%s`)", jsTypeEquivalent))
			if _, ok := g.conf.MethodTransforms[name].ArgTransforms[*arg.Ident]; !ok {
				g.writeIndent(sb, 1)
				sb.WriteString(fmt.Sprintf("%s %s = ", *arg.Type, *arg.Ident))
				if needsCast != nil {
					sb.WriteString(fmt.Sprintf("static_cast<%s>(", *needsCast))
				}
				sb.WriteString(fmt.Sprintf("info[%d].As<Napi::%s>().%s()", i, strings.ReplaceAll(napiTypeHandler, "Is", ""), valGetter))
				if needsCast != nil {
					sb.WriteByte(')')
				}
				sb.WriteString(";\n")
			}
		} else if isClass(*arg.Type, classes) {
			g.writeArgChecker(sb, name, "IsExternal", i, fmt.Sprintf("native `%s` (typeof `Napi::External<%s::%s>`)", *arg.Type, *g.NameSpace, *arg.Type))
			if _, ok := g.conf.MethodTransforms[name].ArgTransforms[*arg.Ident]; !ok {
				g.writeIndent(sb, 1)
				sb.WriteString(fmt.Sprintf("%s::%s* %s = UnExternalize<%s::%s>(info[%d]);\n", *g.NameSpace, *arg.Type, *arg.Ident, *g.NameSpace, *arg.Type, i))
			}
		} else if strings.Contains(*arg.Type, "std::vector") {
			argType := *arg.Type
			type_test := argType[strings.Index(argType, "<")+1 : strings.Index(argType, ">")]
			tsType, isObject := CPPTypeToTS(type_test)
			g.writeArgChecker(sb, name, "IsArray", i, fmt.Sprintf("typeof `%s[]`)", tsType))
			g.writeIndent(sb, 1)
			sb.WriteString(fmt.Sprintf("int len_%s = info[%d].As<Napi::Array>().Length();\n", *arg.Ident, i))
			g.writeIndent(sb, 1)
			sb.WriteString(fmt.Sprintf("for (auto i = 0; i < len_%s; ++i) {\n", *arg.Ident))
			g.writeIndent(sb, 2)
			if isObject {
				sb.WriteString(fmt.Sprintf("if (!info[%d].As<Napi::Array>().Get(i).IsExternal()) {\n", i))
			} else {
				sb.WriteString(fmt.Sprintf("if (!info[%d].As<Napi::Array>().Get(i).Is%s()) {\n", i, g.casers.upper.String(tsType[0:1])+tsType[1:]))
			}
			g.writeIndent(sb, 3)
			sb.WriteString(fmt.Sprintf("Napi::TypeError::New(env, (%q + std::to_string(i) + %q)).ThrowAsJavaScriptException();\n", fmt.Sprintf("`%s` expects args[%d][", name, i), fmt.Sprintf("] to be typeof `%s`", tsType)))
			g.writeIndent(sb, 3)
			sb.WriteString("return env.Null();\n")
			g.writeIndent(sb, 2)
			sb.WriteString("}\n")
			g.writeIndent(sb, 1)
			sb.WriteString("}\n")

		} else if v, ok := g.conf.TypeMappings[*arg.Type]; ok {
			g.writeIndent(sb, 1)
			if strings.Contains(v.TSType, "Array") || strings.Contains(v.TSType, "[]") {
				sb.WriteString(fmt.Sprintf("if (!info[%d].IsArray()) {\n", i))
			} else if strings.Contains(v.TSType, "any") || strings.Contains(v.TSType, "object") || strings.Contains(v.TSType, "Record<") || strings.Contains(v.TSType, "Map<") {
				sb.WriteString(fmt.Sprintf("if (!info[%d].IsExternal()) {\n", i))
			} else if strings.Contains(v.TSType, "string") {
				sb.WriteString(fmt.Sprintf("if (!info[%d].IsString()) {\n", i))
			} else if strings.Contains(v.TSType, "number") {
				sb.WriteString(fmt.Sprintf("if (!info[%d].IsNumber()) {\n", i))
			}
			g.writeIndent(sb, 2)
			sb.WriteString(fmt.Sprintf("Napi::TypeError::New(env, %q).ThrowAsJavaScriptException();\n", fmt.Sprintf("`%s` expects args[%d] to be typeof `%s`", name, i, v.TSType)))
			g.writeIndent(sb, 2)
			sb.WriteString("return env.Null();\n")
			g.writeIndent(sb, 1)
			sb.WriteString("}\n")
		}

		if v, ok := g.conf.MethodTransforms[name].ArgTransforms[*arg.Ident]; ok {
			if !strings.Contains(v, "/arg_") {
				g.writeIndent(sb, 1)
				sb.WriteString(strings.ReplaceAll(v, "/arg/", fmt.Sprintf("info[%d]", i)))
			}
		}
	}

	// TODO: all logic for arg type checks needs to live here
}

func (g *PackageGenerator) writeMethod(sb *strings.Builder, m *CPPMethod, classes map[string]*CPPClass) {
	parsedName := "_" + *m.Ident
	sb.WriteString(fmt.Sprintf("static Napi::Value %s(const Napi::CallbackInfo& info) {\n", parsedName))
	g.writeIndent(sb, 1)
	sb.WriteString("Napi::Env env = info.Env();\n")
	// if len(m.Overloads) == 1 {
	var arg_count int
	if v, ok := g.conf.MethodTransforms[*m.Ident]; ok {
		arg_count = v.ArgCount
		m.ExpectedArgs = arg_count
	} else {
		arg_count = len(*m.Overloads[0])
		m.ExpectedArgs = arg_count
	}
	// single overload, parse args
	g.writeArgChecks(sb, *m.Ident, m.Overloads[0], arg_count, classes)

	obj_name := ""
	outType := *m.Returns
	for i, arg := range *m.Overloads[0] {
		if i > arg_count {
			break
		}
		tmpType := *arg.Type
		if v, ok := g.conf.MethodTransforms[*m.Ident].ArgTransforms[*arg.Ident]; ok && strings.Contains(v, "/arg_") {
			g.writeIndent(sb, 2)
			if strings.Contains(v, "/arg_") {
				for j, val := range *m.Overloads[0] {
					v = strings.ReplaceAll(v, fmt.Sprintf("/arg_%d/", j), *val.Ident)
				}
			}
			// TODO: this might need better handling for class wrappers
			sb.WriteString(strings.ReplaceAll(v, "/arg/", fmt.Sprintf("info[%d]", i)))
		} else if isClass(*arg.Type, classes) {
			obj_name = *arg.Ident
		} else if strings.Contains(*arg.Type, "std::vector") && !strings.EqualFold(tmpType[strings.Index(*arg.Type, "<")+1:strings.Index(*arg.Type, ">")], *m.Returns) {
			g.writeIndent(sb, 2)
			sb.WriteString(fmt.Sprintf("auto axes = jsArrayArg<%s>(info[%d].As<Napi::Array>(), g_row_major, %s->ndim(), env);\n", tmpType[strings.Index(*arg.Type, "<")+1:strings.Index(*arg.Type, ">")], i, obj_name))
		} else {
			fmt.Println("TODO: handle type ", *arg.Type)
		}
	}
	g.writeIndent(sb, 2)
	sb.WriteString(fmt.Sprintf("%s::%s _res;\n", *g.NameSpace, outType))
	if v, ok := g.conf.MethodTransforms[*m.Ident]; ok && v.ReturnTransforms != "" {
		parsed_transform := strings.ReplaceAll(v.ReturnTransforms, "/return/", "_res")
		for i, arg := range *m.Overloads[0] {
			fmtd_arg := ""
			if isClass(*arg.Type, classes) {
				fmtd_arg = fmt.Sprintf("*(%s)", *arg.Ident)
			} else {
				fmtd_arg = *arg.Ident
			}
			parsed_transform = strings.ReplaceAll(parsed_transform, fmt.Sprintf("/arg_%d/", i), fmtd_arg)
		}
		transformed_lines := strings.Split(parsed_transform, "\n")
		length := len(transformed_lines)
		for i, line := range transformed_lines {
			g.writeIndent(sb, 2)
			if i == length-1 {
				sb.WriteString(line)
			} else {
				sb.WriteString(fmt.Sprintf("%s\n", line))
			}
		}
	} else {
		g.writeIndent(sb, 2)
		sb.WriteString(fmt.Sprintf("_res = %s::%s(", *g.NameSpace, *m.Ident))
		for i, arg := range *m.Overloads[0] {
			if i > 0 {
				sb.WriteString(", ")
			}
			if _, ok := g.conf.TypeMappings[*arg.Type]; ok {
				sb.WriteString(fmt.Sprintf("%s::%s(%s)", *g.NameSpace, *arg.Type, *arg.Ident))
			} else if isClass(*arg.Type, classes) {
				sb.WriteString(fmt.Sprintf("*(%s)", *arg.Ident))
			} else {
				sb.WriteString(*arg.Ident)
			}
		}
		sb.WriteString(");\n")
	}
	if v, ok := g.conf.GlobalTypeOutTransforms[*m.Returns]; ok {
		g.writeIndent(sb, 2)
		sb.WriteString(strings.ReplaceAll(v, "/return/", "_res"))
	}
	g.writeIndent(sb, 2)
	sb.WriteString(fmt.Sprintf("auto* out = new %s::%s(_res);\n", *g.NameSpace, outType))
	g.writeIndent(sb, 2)
	sb.WriteString(fmt.Sprintf("Napi::External<%s::%s> _external_out = Externalize%s(env, out);\n", *g.NameSpace, outType, outType))
	g.writeIndent(sb, 2)
	sb.WriteString("return _external_out;\n")

	/* TODO: Handle cases w multiple overloads
	} else {
		// TODO: handle cases w multiple overloads
		g.writeIndent(sb, 1)
		sb.WriteString("return env.Null();\n")
	}
	*/
	sb.WriteString("}\n\n")
}

func (g *PackageGenerator) writeClassField(sb *strings.Builder, f *CPPFieldDecl, className string, classes map[string]*CPPClass) {
	if f.Ident != nil && g.conf.IsFieldWrapped(className, *f.Ident) {
		var returnType string
		isVoid := false
		if f.Returns != nil && *f.Returns.FullType != "void" {
			returnType = *f.Returns.FullType
		} else {
			isVoid = true
		}
		sb.WriteString("static ")
		if isVoid {
			sb.WriteString("void ")
		} else {
			sb.WriteString("Napi::Value ")
		}
		sb.WriteString(fmt.Sprintf("_%s(const Napi::CallbackInfo& info) {\n", *f.Ident))
		g.writeIndent(sb, 1)
		sb.WriteString("Napi::Env env = info.Env();\n")
		argCount := 1
		if f.Args != nil {
			argCount += len(*f.Args)
		}
		g.writeIndent(sb, 1)
		sb.WriteString(fmt.Sprintf("if (info.Length() != %d) {\n", argCount))
		g.writeIndent(sb, 2)
		sb.WriteString(fmt.Sprintf("Napi::TypeError::New(env, %q).ThrowAsJavaScriptException();\n", fmt.Sprintf("`%s` expects exactly %d args", *f.Ident, argCount)))
		g.writeIndent(sb, 2)
		sb.WriteString("return env.Null();\n")
		g.writeIndent(sb, 1)
		sb.WriteString("}\n")

		g.writeIndent(sb, 1)
		sb.WriteString(fmt.Sprintf("if (!info[%d].IsExternal()) {\n", 0))
		g.writeIndent(sb, 2)
		sb.WriteString(fmt.Sprintf("Napi::TypeError::New(env, %q).ThrowAsJavaScriptException();\n", fmt.Sprintf("`%s` expects args[%d] to be typeof `%s`", *f.Ident, 0, className)))
		g.writeIndent(sb, 2)
		sb.WriteString("return env.Null();\n")
		g.writeIndent(sb, 1)
		sb.WriteString("}\n")
		g.writeIndent(sb, 1)
		sb.WriteString(fmt.Sprintf("%s::%s* _tmp_external = UnExternalize<%s::%s>(info[%d]);\n", *g.NameSpace, className, *g.NameSpace, className, 0))
		if f.Args != nil {
			for i, arg := range *f.Args {
				typeHandler, isObject := CPPTypeToTS(*arg.Type)
				if v, ok := g.conf.TypeMappings[*arg.Type]; ok {
					g.writeIndent(sb, 1)
					sb.WriteString(fmt.Sprintf("if (!info[%d].Is%s()) {\n", i+1, v.NapiType))
					g.writeIndent(sb, 2)
					sb.WriteString(fmt.Sprintf("Napi::TypeError::New(env, %q).ThrowAsJavaScriptException();\n", fmt.Sprintf("`%s` expects args[%d] to be typeof `%s`", *f.Ident, i+1, typeHandler)))
					g.writeIndent(sb, 2)
					sb.WriteString("return env.Null();\n")
					g.writeIndent(sb, 1)
					sb.WriteString("}\n")
					g.writeIndent(sb, 1)
					sb.WriteString(fmt.Sprintf("auto %s = static_cast<%s>(info[%d].As<Napi::%s>().%s());\n", *arg.Ident, v.CastsTo, i, v.NapiType, v.CastNapi))
				} else if isObject {
					// TODO: handle
				} else {
					// TODO: handle
				}
			}
		}
		g.writeIndent(sb, 1)
		if f.Returns != nil && *f.Returns.FullType != "void" {
			sb.WriteString("auto _res = ")
			returnType = *f.Returns.FullType
		}
		sb.WriteString(fmt.Sprintf("_tmp_external->%s(", *f.Ident))
		if f.Args != nil {
			for i, arg := range *f.Args {
				if i > 0 {
					sb.WriteString(", ")
				}
				if _, ok := g.conf.TypeMappings[*arg.Type]; ok {
					sb.WriteString(fmt.Sprintf("%s::%s(%s)", *g.NameSpace, *arg.Type, *arg.Ident))
				} else if isClass(*arg.Type, classes) {
					sb.WriteString(fmt.Sprintf("*(%s)", *arg.Ident))
				} else {
					sb.WriteString(*arg.Ident)
				}
			}
		}
		sb.WriteString(");\n")

		if f.Returns != nil && *f.Returns.FullType != "void" {
			jsType, isObject := CPPTypeToTS(returnType)
			if g.conf.TypeHasHandler(returnType) != nil {
				t := g.conf.TypeHasHandler(returnType)
				g.writeIndent(sb, 1)
				sb.WriteString(strings.ReplaceAll(t.Handler, "/val/", "_res"))
				g.writeIndent(sb, 1)
				sb.WriteString(fmt.Sprintf("return %s;\n", t.OutVar))
			} else if isObject && isClass(returnType, classes) {
				if v, ok := g.conf.GlobalTypeOutTransforms[returnType]; ok {
					g.writeIndent(sb, 1)
					sb.WriteString(strings.ReplaceAll(v, "/return/", "_res"))
				}
				g.writeIndent(sb, 1)
				sb.WriteString(fmt.Sprintf("auto* out = new %s::%s(_res);\n", *g.NameSpace, returnType))
				g.writeIndent(sb, 1)
				sb.WriteString(fmt.Sprintf("Napi::External<%s::%s> _external_out = Externalize%s(env, out);", *g.NameSpace, returnType, returnType))
				g.writeIndent(sb, 1)
				g.writeIndent(sb, 1)
				sb.WriteString("return _external_out;\n")
			} else {
				napiHandler := g.casers.upper.String(jsType[0:1]) + jsType[1:]
				sb.WriteString(fmt.Sprintf("return Napi::%s::New(env, %s);\n", napiHandler, "_res"))
			}
		}
		sb.WriteString("}\n\n")
	}
}

func (g *PackageGenerator) writeAddonExport(sb *strings.Builder, name string) {
	g.writeIndent(sb, 1)
	parsedName := ("_" + name)
	sb.WriteString(fmt.Sprintf("exports.Set(Napi::String::New(env, %q), Napi::Function::New(env, %s));\n", parsedName, parsedName))
}

// makes calls to functions that write bindings
func (g *PackageGenerator) writeBindings(sb *strings.Builder, classes map[string]*CPPClass, methods map[string]*CPPMethod, processedMethods map[string]*CPPMethod) {
	sb.WriteString("#include <napi.h>\n")
	g.writeHeaderFrontmatter(sb)
	g.writeBindingsFrontmatter(sb)
	g.writeFileSourceHeader(sb, *g.Path)
	g.writeGlobalVars(sb)
	// write any helpers functions (non-exported; specified in config)
	g.writeHelpers(sb, classes)

	sb.WriteString("// exported functions\n\n")
	for _, f := range methods {
		g.writeMethod(sb, f, classes)
	}

	for _, f := range processedMethods {
		g.writeMethod(sb, f, classes)
	}

	// write any forced methods (specified in config)
	for _, f := range g.conf.GlobalForcedMethods {
		sb.WriteString(fmt.Sprintf("%s\n", strings.Replace(f.FnBody, f.Name, "_"+f.Name, 1)))
	}

	// writes NAPI `Init` function (init NAPI exports)
	sb.WriteString("// NAPI exports\n\n")
	sb.WriteString("Napi::Object Init(Napi::Env env, Napi::Object exports) {\n")
	for name, c := range classes {
		// check if header contained class constructor declaration(s)
		if c.Decl != nil {
			// write exports for wrapped class fields (specified in config)
			if c.FieldDecl != nil {
				for _, f := range *c.FieldDecl {
					if f.Ident != nil && g.conf.IsFieldWrapped(name, *f.Ident) {
						g.writeAddonExport(sb, *f.Ident)
					}
				}
			}
			// write exports for any optionally forced class methods (specified in config)
			if v, ok := g.conf.ClassOpts[name]; ok && len(v.ForcedMethods) > 0 {
				for _, f := range v.ForcedMethods {
					g.writeAddonExport(sb, f.Name)
				}
			}
		}
	}

	// write exports for methods defined in header
	for _, f := range methods {
		g.writeAddonExport(sb, *f.Ident)
	}

	// write exports for methods requiring pre-processing
	for _, f := range processedMethods {
		g.writeAddonExport(sb, *f.Ident)
	}

	// write any optionally forced global methods (specified in config)
	for _, f := range g.conf.GlobalForcedMethods {
		g.writeAddonExport(sb, f.Name)
	}

	g.writeIndent(sb, 1)
	sb.WriteString("return exports;\n")
	sb.WriteString("}\n\n")
	sb.WriteString("NODE_API_MODULE(addon, Init)\n")
}
