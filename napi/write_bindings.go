package napi

import (
	"fmt"
	"strings"
)

func (g *PackageGenerator) writeArgCountChecker(sb *strings.Builder, name string, expected_arg_count int, optional int) {
	if expected_arg_count == 0 && optional == 0 {
		return
	}
	g.writeIndent(sb, 1)
	if optional == 0 {
		sb.WriteString(fmt.Sprintf("if (info.Length() != %d) {\n", expected_arg_count))
	} else {
		g.writeIndent(sb, 1)
		sb.WriteString("auto _arg_count = info.Length();\n")
		g.writeIndent(sb, 1)
		sb.WriteString(fmt.Sprintf("if (_arg_count < %d || _arg_count > %d) {\n", expected_arg_count, expected_arg_count+optional))
	}
	g.writeIndent(sb, 2)
	var errMsg string
	if optional == 0 {
		errMsg = fmt.Sprintf("`%s` expects exactly %d arg", name, expected_arg_count)
		if expected_arg_count > 1 {
			errMsg += "s"
		}
	} else {
		errMsg = fmt.Sprintf("`%s` expects between %d to %d args", name, expected_arg_count, expected_arg_count+optional)
	}
	sb.WriteString(fmt.Sprintf("Napi::TypeError::New(env, %q).ThrowAsJavaScriptException();\n", errMsg))
	g.writeIndent(sb, 2)
	sb.WriteString("return env.Undefined();\n")
	g.writeIndent(sb, 1)
	sb.WriteString("}\n")
}

func (g *PackageGenerator) writeArgTypeChecker(sb *strings.Builder, name string, checker string, idx int, msg string, indents int, arrName *string, arg *CPPArg) {
	isArrayItem := arrName != nil
	hasDefault := arg != nil && arg.DefaultValue != nil
	if hasDefault {
		// TODO: handle default values for non-enum types
		isEnum, name := g.IsTypeEnum(*arg.Type)
		if isEnum {
			g.writeIndent(sb, indents)
			sb.WriteString(fmt.Sprintf("%s %s;\n", *name, *arg.Ident))
			g.writeIndent(sb, indents)
			sb.WriteString(fmt.Sprintf("if (!info[%d].IsUndefined()) {\n", idx))
			indents++
		}
	}
	g.writeIndent(sb, indents)
	if isArrayItem {
		sb.WriteString(fmt.Sprintf("Napi::Value arrayItem = %s[i];\n", *arrName))
	}
	sb.WriteString("if (!")
	// required to handle checking array index items (i.e. info[0][i])
	if isArrayItem {
		sb.WriteString(fmt.Sprintf("arrayItem.%s", checker))
	} else {
		sb.WriteString(fmt.Sprintf("info[%d].%s", idx, checker))
	}
	sb.WriteString("()) {\n")
	g.writeIndent(sb, indents+1)
	sb.WriteString("Napi::TypeError::New(env, ")
	// customize error msg when checking indexes in array
	if isArrayItem {
		sb.WriteString(fmt.Sprintf("(%q + std::to_string(i) + %q)", fmt.Sprintf("`%s` expects args[%d][", name, idx), fmt.Sprintf("] to be %s", msg)))
	} else {
		sb.WriteString(fmt.Sprintf("%q", fmt.Sprintf("`%s` expects args[%d] to be %s", name, idx, msg)))
	}
	sb.WriteString(").ThrowAsJavaScriptException();\n")
	g.writeIndent(sb, indents+1)
	sb.WriteString("return env.Undefined();\n")
	g.writeIndent(sb, indents)
	sb.WriteString("}\n")
	if hasDefault {
		// TODO: handle default values for non-enum types
		isEnum, enumName := g.IsTypeEnum(*arg.Type)
		if isEnum {
			g.writeIndent(sb, indents)
			sb.WriteString(fmt.Sprintf("%s = static_cast<%s>(info[%d].As<Napi::Number>().Int32Value());\n", *arg.Ident, *enumName, idx))
			indents--
			g.writeIndent(sb, indents)
			sb.WriteString("} else {\n")
			indents++
			g.writeIndent(sb, indents)
			sb.WriteString(fmt.Sprintf("%s = %s::%s;\n", *arg.Ident, *enumName, *arg.DefaultValue.Val))
			indents--
			g.writeIndent(sb, indents)
			sb.WriteString("}\n")
		}
	}
}

func (g *PackageGenerator) writeArgChecks(sb *strings.Builder, name string, args *[]*CPPArg, expected_arg_count int, optionalArgs int) {
	if expected_arg_count == 0 {
		return
	}
	g.writeArgCountChecker(sb, name, expected_arg_count, optionalArgs)
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
		if i >= expected_arg_count {
			break
		}
		if arg.Ident == nil {
			fmt.Printf("WARNING: arg.Ident is nil for %q", name)
		}
		isArgTransform, argTransformVal := g.conf.IsArgTransform(name, *arg.Ident)
		isEnum, _ := g.IsTypeEnum(*arg.Type)
		if isEnum {
			g.writeArgTypeChecker(sb, name, "IsNumber", i, "typeof `number`", 1, nil, arg)
		} else if arg.IsPrimitive {
			if !arg.IsPointer {
				napiTypeHandler := "IsNumber"
				jsTypeEquivalent := "number"
				valGetter := "Value"
				var needsCast *string
				switch *arg.Type {
				case "float":
					valGetter = "FloatValue"
				case "std::string":
					valGetter = "Utf8Value"
					napiTypeHandler = "IsString"
					jsTypeEquivalent = "string"
				case "string":
					valGetter = "Utf8Value"
					napiTypeHandler = "IsString"
					jsTypeEquivalent = "string"
				case "double":
					valGetter = "DoubleValue"
				case "long long", "char", "signed", "int8_t", "int32_t", "int16_t", "short":
					valGetter = "Int32Value"
					needsCast = arg.Type
				case "int", "int64_t":
					valGetter = "Int64Value"
					needsCast = arg.Type
				case "unsigned long long", "unsigned char", "unsigned", "uint8_t", "uint16_t", "unsigned short", "uint32_t":
					valGetter = "Uint32Value"
					needsCast = arg.Type
				case "unsigned int", "uint64_t", "size_t", "uintptr_t":
					valGetter = "Uint64Value"
					needsCast = arg.Type
				case "bool":
					napiTypeHandler = "IsBoolean"
					jsTypeEquivalent = "boolean"
				}
				g.writeArgTypeChecker(sb, name, napiTypeHandler, i, fmt.Sprintf("typeof `%s`)", jsTypeEquivalent), 1, nil, arg)
				// get val from arg if no transform is specified
				if !isArgTransform {
					g.writeIndent(sb, 1)
					sb.WriteString(fmt.Sprintf("%s %s = ", *arg.Type, *arg.Ident))
					// only cast when necessary
					if needsCast != nil {
						sb.WriteString(fmt.Sprintf("static_cast<%s>(", *needsCast))
					}
					sb.WriteString(fmt.Sprintf("info[%d].As<Napi::%s>().%s()", i, strings.ReplaceAll(napiTypeHandler, "Is", ""), valGetter))
					if needsCast != nil {
						sb.WriteByte(')')
					}
					sb.WriteString(";\n")
				}
			} else {
				jsTypeEquivalent, arrayType, needsCast, _ := PrimitivePtrToTS(*arg.Type)
				g.writeArgTypeChecker(sb, name, "IsTypedArray", i, fmt.Sprintf("typeof `%s`)", jsTypeEquivalent), 1, nil, arg)
				// get val from arg if no transform is specified
				if !isArgTransform {
					g.writeIndent(sb, 1)
					sb.WriteString(fmt.Sprintf("%s *%s = ", *arg.Type, *arg.Ident))
					if needsCast != nil {
						sb.WriteString(fmt.Sprintf("reinterpret_cast<%s *>(", *needsCast))
					}
					sb.WriteString(fmt.Sprintf("info[%d].As<Napi::TypedArrayOf<%s>>().Data()", i, arrayType))
					if needsCast != nil {
						sb.WriteByte(')')
					}
					sb.WriteString(";\n")
				}
			}
		} else if isClass(*arg.Type, g.ParsedData.Classes) {
			g.writeArgTypeChecker(sb, name, "IsExternal", i, fmt.Sprintf("native `%s` (typeof `Napi::External<%s::%s>`)", *arg.Type, *g.NameSpace, *arg.Type), 1, nil, arg)
			if !isArgTransform {
				g.writeIndent(sb, 1)
				sb.WriteString(fmt.Sprintf("%s::%s* %s = UnExternalize<%s::%s>(info[%d]);\n", *g.NameSpace, *arg.Type, *arg.Ident, *g.NameSpace, *arg.Type, i))
			}
		} else if strings.Contains(*arg.Type, "std::vector") {
			argType := *arg.Type
			type_test := argType[strings.Index(argType, "<")+1 : strings.Index(argType, ">")]
			tsType, isObject := CPPTypeToTS(type_test, false)
			g.writeArgTypeChecker(sb, name, "IsArray", i, fmt.Sprintf("typeof `%s[]`)", tsType), 1, nil, arg)
			g.writeIndent(sb, 1)
			arrName := fmt.Sprintf("_tmp_parsed_%s", *arg.Ident)
			sb.WriteString(fmt.Sprintf("Napi::Array %s = info[%d].As<Napi::Array>();\n", arrName, i))
			g.writeIndent(sb, 1)
			sb.WriteString(fmt.Sprintf("size_t len_%s = %s.Length();\n", *arg.Ident, arrName))
			g.writeIndent(sb, 1)
			sb.WriteString(fmt.Sprintf("for (size_t i = 0; i < len_%s; ++i) {\n", *arg.Ident))
			g.writeIndent(sb, 2)
			if isObject {
				g.writeArgTypeChecker(sb, name, "IsExternal", i, fmt.Sprintf("native `%s` (typeof `Napi::External<%s::%s>`)", tsType, *g.NameSpace, tsType), 2, &arrName, nil)
			} else {
				g.writeArgTypeChecker(sb, name, fmt.Sprintf("Is%s", g.casers.upper.String(tsType[0:1])+tsType[1:]), i, fmt.Sprintf("typeof `%s`", tsType), 2, &arrName, nil)
			}
			g.writeIndent(sb, 1)
			sb.WriteString("}\n")
		} else if v, ok := g.conf.TypeMappings[*arg.Type]; ok {
			g.writeIndent(sb, 1)
			errMsg := fmt.Sprintf("typeof `%s`)", v.TSType)
			if strings.Contains(v.TSType, "Array") || strings.Contains(v.TSType, "[]") {
				g.writeArgTypeChecker(sb, name, "IsArray", i, errMsg, 1, nil, arg)
			} else if strings.Contains(v.TSType, "any") || strings.Contains(v.TSType, "object") || strings.Contains(v.TSType, "Record<") || strings.Contains(v.TSType, "Map<") {
				g.writeArgTypeChecker(sb, name, "IsObject", i, errMsg, 1, nil, arg)
			} else if strings.Contains(v.TSType, "string") {
				g.writeArgTypeChecker(sb, name, "IsString", i, errMsg, 1, nil, arg)
			} else if strings.Contains(v.TSType, "number") {
				g.writeArgTypeChecker(sb, name, "IsNumber", i, errMsg, 1, nil, arg)
			}
		}
		if isArgTransform && !strings.Contains(*argTransformVal, "/arg_") {
			g.writeIndent(sb, 1)
			sb.WriteString(strings.ReplaceAll(*argTransformVal, "/arg/", fmt.Sprintf("info[%d]", i)))
		}
	}
}

func (g *PackageGenerator) writeMethod(sb *strings.Builder, m *CPPMethod) {
	parsedName := "_" + *m.Ident
	sb.WriteString(fmt.Sprintf("static Napi::Value %s(const Napi::CallbackInfo& info) {\n", parsedName))
	g.writeIndent(sb, 1)
	sb.WriteString("Napi::Env env = info.Env();\n")
	// if len(m.Overloads) == 1 {
	arg_count := 0
	optional_args := 0
	if v, ok := g.conf.MethodTransforms[*m.Ident]; ok {
		arg_count = v.ArgCount
		m.ExpectedArgs = arg_count
	} else {
		for _, arg := range *m.Overloads[0] {
			isEnum, _ := g.IsTypeEnum(*arg.Type)
			if arg.DefaultValue != nil && isEnum {
				optional_args++
			}
			arg_count++
		}
		m.ExpectedArgs = arg_count
	}
	// single overload, parse args
	g.writeArgChecks(sb, *m.Ident, m.Overloads[0], arg_count, optional_args)

	obj_name := ""
	outType := *m.Returns
	for i, arg := range *m.Overloads[0] {
		if i > arg_count {
			break
		}
		isArgTransform, argTransformVal := g.conf.IsArgTransform(*m.Ident, *arg.Ident)

		tmpType := *arg.Type
		if isArgTransform && strings.Contains(*argTransformVal, "/arg_") {
			g.writeIndent(sb, 2)
			for j, val := range *m.Overloads[0] {
				*argTransformVal = strings.ReplaceAll(*argTransformVal, fmt.Sprintf("/arg_%d/", j), *val.Ident)
			}
			sb.WriteString(strings.ReplaceAll(*argTransformVal, "/arg/", fmt.Sprintf("info[%d]", i)))
			// TODO: this might need better handling for class wrappers
		} else if isClass(*arg.Type, g.ParsedData.Classes) {
			obj_name = *arg.Ident
		} else if strings.Contains(*arg.Type, "std::vector") && !strings.EqualFold(tmpType[strings.Index(*arg.Type, "<")+1:strings.Index(*arg.Type, ">")], *m.Returns) {
			g.writeIndent(sb, 2)
			invertVal := "false"
			if g.conf.VectorOpts.DimAccessor != "" {
				invertVal = fmt.Sprintf("%s->%s", obj_name, g.conf.VectorOpts.DimAccessor)
			}
			sb.WriteString(fmt.Sprintf("auto %s = jsArrayToVector<%s>(info[%d].As<Napi::Array>(), g_row_major, %s);\n", *arg.Ident, tmpType[strings.Index(*arg.Type, "<")+1:strings.Index(*arg.Type, ">")], i, invertVal))
		} else if !arg.IsPrimitive {
			// TODO: all arg coercions should probably be written in a single place to prevent duplication
			var ptrType string
			if arg.IsPointer {
				ptrType = "*"
			}
			fmt.Printf("TODO: Method %q has unhandled argument: `%s %s%s`\n", *m.Ident, *arg.Type, ptrType, *arg.Ident)
		}
	}
	g.writeIndent(sb, 2)
	_, arrayType, needsCast, _ := PrimitivePtrToTS(*m.Returns)

	if !m.ReturnsPrimitive {
		sb.WriteString(fmt.Sprintf("%s::%s ", *g.NameSpace, outType))
	} else if m.ReturnsPointer && needsCast != nil && arrayType != "" {
		sb.WriteString(fmt.Sprintf("%s ", arrayType))
	} else {
		sb.WriteString(fmt.Sprintf("%s ", *m.Returns))
	}

	if m.ReturnsPointer {
		sb.WriteByte('*')
	}
	sb.WriteString("_res;\n")

	isReturnTransform, isGrouped, transform := g.conf.IsReturnTransform(m)
	if isReturnTransform {
		if isGrouped {
			g.writeIndent(sb, 2)
			sb.WriteString(fmt.Sprintf("_res = %s::%s(", *g.NameSpace, *m.Ident))
			for i, arg := range *m.Overloads[0] {
				if i > 0 {
					sb.WriteString(", ")
				}
				isEnum, _ := g.IsTypeEnum(*arg.Type)
				if isEnum {
					sb.WriteString(*arg.Ident)
				} else if _, ok := g.conf.TypeMappings[*arg.Type]; ok {
					sb.WriteString(fmt.Sprintf("%s::%s(%s)", *g.NameSpace, *arg.Type, *arg.Ident))
				} else if isClass(*arg.Type, g.ParsedData.Classes) {
					sb.WriteString(fmt.Sprintf("*(%s)", *arg.Ident))
				} else {
					sb.WriteString(*arg.Ident)
				}
			}
			sb.WriteString(");\n")
		}
		parsed_transform := strings.ReplaceAll(*transform, "/return/", "_res")
		for i, arg := range *m.Overloads[0] {
			fmtd_arg := ""
			if isClass(*arg.Type, g.ParsedData.Classes) {
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
		// handle w/o any transformations
	} else {
		g.writeIndent(sb, 2)
		sb.WriteString("_res = ")
		_, arrayType, needsCast, _ := PrimitivePtrToTS(*m.Returns)
		if m.ReturnsPointer && needsCast != nil && arrayType != "" {
			sb.WriteString(fmt.Sprintf("reinterpret_cast<%s *>(", arrayType))
		}
		sb.WriteString(fmt.Sprintf("%s::%s(", *g.NameSpace, *m.Ident))
		for i, arg := range *m.Overloads[0] {
			if i > arg_count {
				break
			}
			if i > 0 {
				sb.WriteString(", ")
			}
			if _, ok := g.conf.TypeMappings[*arg.Type]; ok {
				sb.WriteString(fmt.Sprintf("%s::%s(%s)", *g.NameSpace, *arg.Type, *arg.Ident))
			} else if isClass(*arg.Type, g.ParsedData.Classes) {
				sb.WriteString(fmt.Sprintf("*(%s)", *arg.Ident))
			} else {
				sb.WriteString(*arg.Ident)
			}
		}
		sb.WriteByte(')')
		if m.ReturnsPointer && needsCast != nil && arrayType != "" {
			sb.WriteByte(')')
		}
		sb.WriteString(";\n")
	}
	if *m.Returns != "void" || isReturnTransform {
		g.writeIndent(sb, 2)
		returnType := *m.Returns
		if m.ReturnsPrimitive && m.ReturnsPointer {
			_, arrayType, _, napi_short_type := PrimitivePtrToTS(returnType)
			sb.WriteString("size_t _res_byte_len = sizeof(_res);\n")
			g.writeIndent(sb, 2)
			sb.WriteString("size_t _res_elem_len = _res_byte_len / sizeof(*_res);\n")
			g.writeIndent(sb, 2)
			sb.WriteString(fmt.Sprintf("std::unique_ptr<std::vector<%s>> _res_native_array = std::make_unique<std::vector<%s>>(_res, _res + _res_elem_len);\n", arrayType, arrayType))
			g.writeIndent(sb, 2)
			sb.WriteString(fmt.Sprintf("Napi::ArrayBuffer _res_arraybuffer = Napi::ArrayBuffer::New(env, _res_native_array->data(), _res_byte_len, DeleteArrayBuffer<%s>, _res_native_array.get());\n", arrayType))
			g.writeIndent(sb, 2)
			sb.WriteString("_res_native_array.release();\n")
			g.writeIndent(sb, 2)
			sb.WriteString("Napi::MemoryManagement::AdjustExternalMemory(env, _res_byte_len);\n")
			g.writeIndent(sb, 2)
			sb.WriteString(fmt.Sprintf("return Napi::TypedArrayOf<%s>::New(env, _res_elem_len, _res_arraybuffer, 0, napi_%s_array);\n", arrayType, napi_short_type))
		} else {
			jsType, isObject := CPPTypeToTS(returnType, false)
			if g.conf.TypeHasHandler(returnType) != nil {
				t := g.conf.TypeHasHandler(returnType)
				g.writeIndent(sb, 1)
				sb.WriteString(strings.ReplaceAll(t.Handler, "/val/", "_res"))
				g.writeIndent(sb, 1)
				sb.WriteString(fmt.Sprintf("return %s;\n", t.OutVar))
			} else if isObject && isClass(returnType, g.ParsedData.Classes) {
				if v, ok := g.conf.GlobalTypeOutTransforms[returnType]; ok {
					g.writeIndent(sb, 1)
					sb.WriteString(strings.ReplaceAll(v, "/return/", "_res"))
				}
				g.writeIndent(sb, 1)
				sb.WriteString(fmt.Sprintf("auto* out = new %s::%s(_res);\n", *g.NameSpace, returnType))
				g.writeIndent(sb, 1)
				sb.WriteString(fmt.Sprintf("return Externalize%s(env, out);", returnType))
			} else {
				napiHandler := g.casers.upper.String(jsType[0:1]) + jsType[1:]
				if napiHandler == "Bigint" {
					napiHandler = "BigInt"
				}
				sb.WriteString(fmt.Sprintf("return Napi::%s::New(env, %s);\n", napiHandler, "_res"))
			}
		}
	}
	/* TODO: Handle cases w multiple overloads
	} else {
		// TODO: handle cases w multiple overloads
		g.writeIndent(sb, 1)
		sb.WriteString("return env.Undefined();\n")
	}
	*/
	sb.WriteString("}\n\n")
}

func (g *PackageGenerator) writeClassField(sb *strings.Builder, f *CPPFieldDecl, className string) {
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
		g.writeArgCountChecker(sb, *f.Ident, argCount, 0)
		g.writeArgTypeChecker(sb, *f.Ident, "IsExternal", 0, fmt.Sprintf("native `%s` (typeof `Napi::External<%s::%s>`)", className, *g.NameSpace, className), 1, nil, nil)
		g.writeIndent(sb, 1)
		sb.WriteString(fmt.Sprintf("%s::%s* _tmp_external = UnExternalize<%s::%s>(info[%d]);\n", *g.NameSpace, className, *g.NameSpace, className, 0))
		if f.Args != nil {
			for i, arg := range *f.Args {
				typeHandler, _ := CPPTypeToTS(*arg.Type, false)
				if v, ok := g.conf.TypeMappings[*arg.Type]; ok {
					g.writeArgTypeChecker(sb, *f.Ident, fmt.Sprintf("Is%s", v.NapiType), i+1, fmt.Sprintf("typeof `%s`)", typeHandler), 1, nil, nil)
					g.writeIndent(sb, 1)
					sb.WriteString(fmt.Sprintf("auto %s = static_cast<%s>(info[%d].As<Napi::%s>().%s());\n", *arg.Ident, v.CastsTo, i, v.NapiType, v.CastNapi))
				}
				// TODO: handle unmapped types
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
				} else if isClass(*arg.Type, g.ParsedData.Classes) {
					sb.WriteString(fmt.Sprintf("*(%s)", *arg.Ident))
				} else {
					sb.WriteString(*arg.Ident)
				}
			}
		}
		sb.WriteString(");\n")

		if f.Returns != nil && *f.Returns.FullType != "void" {
			jsType, isObject := CPPTypeToTS(returnType, false)
			if g.conf.TypeHasHandler(returnType) != nil {
				t := g.conf.TypeHasHandler(returnType)
				g.writeIndent(sb, 1)
				sb.WriteString(strings.ReplaceAll(t.Handler, "/val/", "_res"))
				g.writeIndent(sb, 1)
				sb.WriteString(fmt.Sprintf("return %s;\n", t.OutVar))
			} else if isObject && isClass(returnType, g.ParsedData.Classes) {
				if v, ok := g.conf.GlobalTypeOutTransforms[returnType]; ok {
					g.writeIndent(sb, 1)
					sb.WriteString(strings.ReplaceAll(v, "/return/", "_res"))
				}
				g.writeIndent(sb, 1)
				sb.WriteString(fmt.Sprintf("auto* out = new %s::%s(_res);\n", *g.NameSpace, returnType))
				g.writeIndent(sb, 1)
				sb.WriteString(fmt.Sprintf("return Externalize%s(env, out);", returnType))
			} else {
				napiHandler := g.casers.upper.String(jsType[0:1]) + jsType[1:]
				usedVar := "_res"
				if napiHandler == "Bigint" {
					napiHandler = "BigInt"
					usedVar = "(int64_t)_res"
				}
				sb.WriteString(fmt.Sprintf("return Napi::%s::New(env, %s);\n", napiHandler, usedVar))
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
func (g *PackageGenerator) writeBindings(sb *strings.Builder) {
	g.writeRequiredIncludes(sb)

	// g.writeHeaderFrontmatter(sb)
	g.writeBindingsFrontmatter(sb)
	g.writeFileSourceHeader(sb, *g.Path)
	g.writeGlobalVars(sb)
	// write any helpers functions (non-exported; specified in config)
	g.writeHelpers(sb)

	sb.WriteString("// exported functions\n\n")
	// write methods (not requiring preprocessing)
	for _, f := range g.ParsedData.Methods {
		g.writeMethod(sb, f)
	}

	// write methods that required preprocessing
	for _, f := range g.ParsedData.Lits {
		g.writeMethod(sb, f)
	}

	// write any forced methods (specified in config)
	for _, f := range g.conf.GlobalForcedMethods {
		sb.WriteString(fmt.Sprintf("%s\n\n", strings.Replace(f.FnBody, f.Name, "_"+f.Name, 1)))
	}

	// writes NAPI `Init` function (init NAPI exports)
	sb.WriteString("// NAPI exports\n\n")
	sb.WriteString("Napi::Object Init(Napi::Env env, Napi::Object exports) {\n")
	for name, c := range g.ParsedData.Classes {
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
	for _, f := range g.ParsedData.Methods {
		g.writeAddonExport(sb, *f.Ident)
	}

	// write exports for methods requiring pre-processing
	for _, f := range g.ParsedData.Lits {
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
