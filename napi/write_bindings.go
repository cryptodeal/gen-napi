package napi

import (
	"fmt"
	"strings"
)

func (g *PackageGenerator) writeArgCountChecker(sb *strings.Builder, name string, expected_arg_count int, optional int, is_void bool) {
	if expected_arg_count == 0 && optional == 0 {
		return
	}
	g.writeIndent(sb, 1)
	sb.WriteString("const auto _arg_count = info.Length();\n")
	if optional == 0 {
		g.writeIndent(sb, 1)
		sb.WriteString(fmt.Sprintf("if (_arg_count != %d) {\n", expected_arg_count))
	} else {
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
		errMsg += ", but received "
	} else {
		errMsg = fmt.Sprintf("`%s` expects %d to %d args, but received ", name, expected_arg_count, expected_arg_count+optional)
	}
	sb.WriteString(fmt.Sprintf("Napi::TypeError::New(env, %q + std::to_string(_arg_count)).ThrowAsJavaScriptException();\n", errMsg))
	g.writeIndent(sb, 2)
	if !is_void {
		sb.WriteString("return env.Undefined();\n")
	} else {
		sb.WriteString("return;\n")
	}
	g.writeIndent(sb, 1)
	sb.WriteString("}\n")
}

func (g *PackageGenerator) writeArgTypeChecker(sb *strings.Builder, name string, checker string, idx int, msg string, indents int, arrName *string, arg *CPPArg, is_void bool) {
	hasDefault := arg != nil && arg.DefaultValue != nil
	/*
		if `hasDefault` is true, the generated JS wrapper has
		generated the corresponding default values for the arg;
		any type checks would be redundant, so skip for performance
	*/
	if arg == nil || hasDefault {
		return
	}

	isArrayItem := arrName != nil

	if isArrayItem {
		sb.WriteString(fmt.Sprintf("Napi::Array %s = info[%d].As<Napi::Array>();\n", *arrName, idx))
		g.writeIndent(sb, 1)
		sb.WriteString(fmt.Sprintf("size_t len_%s = %s.Length();\n", *arg.Name, *arrName))
		g.writeIndent(sb, 1)
		sb.WriteString(fmt.Sprintf("for (size_t i = 0; i < len_%s; ++i) {\n", *arg.Name))
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
	if !is_void {
		sb.WriteString("return env.Undefined();\n")
	} else {
		sb.WriteString("return;\n")
	}

	g.writeIndent(sb, indents)
	sb.WriteString("}\n")
	if isArrayItem {
		g.writeIndent(sb, 1)
		sb.WriteString("}\n")
	}
}

func (g *PackageGenerator) writeArgChecks(sb *strings.Builder, name string, args *[]*CPPArg, expected_arg_count int, optionalArgs int, is_void bool) {
	if expected_arg_count == 0 {
		return
	}
	g.writeArgCountChecker(sb, name, expected_arg_count, optionalArgs, is_void)
	// write arg checks, transforms, and declare arg variable as possible
	if v, ok := g.conf.MethodTransforms[name]; ok && v.ArgCheckTransforms != "" {
		sb.WriteString(v.ArgCheckTransforms)
		for i, arg := range *args {
			if arg == nil {
				continue
			}
			if i > expected_arg_count {
				break
			}
			// write arg transform only if it doesn't rely on other args
			if v2, ok2 := v.ArgTransforms[*arg.Name]; ok2 && !strings.Contains(v2, "/arg_") {
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
		if i >= expected_arg_count && optionalArgs == 0 {
			break
		}
		if arg == nil {
			continue
		}
		if arg.Name == nil {
			fmt.Printf("WARNING: arg.Ident is nil for %q", name)
		}
		checks := arg.GetArgChecks(g)
		g.WriteArgCheck(sb, checks, name, i, arg, is_void)
		g.WriteArgGetter(sb, checks, name, arg, i)
	}
}

func (g *PackageGenerator) writeMethod(sb *strings.Builder, m *CPPMethod) {
	parsedName := "_" + *m.Name
	return_type := "Napi::Value"
	is_void := m.IsVoid()
	if is_void {
		return_type = "void"
	}
	sb.WriteString(fmt.Sprintf("static %s %s(const Napi::CallbackInfo& info) {\n", return_type, parsedName))
	g.writeIndent(sb, 1)
	sb.WriteString("Napi::Env env = info.Env();\n")
	// if len(m.Overloads) == 1 {
	arg_count := 0
	optional_args := 0
	if v, ok := g.conf.MethodTransforms[*m.Name]; ok && v.ArgCount != nil {
		arg_count = *v.ArgCount
		m.ExpectedArgs = arg_count
	} else {
		for _, arg := range *m.Overloads[0] {
			isEnum, _ := g.IsTypeEnum(*arg.Type)
			if arg.DefaultValue != nil && isEnum {
				optional_args++
			} else {
				arg_count++
			}
		}
		m.ExpectedArgs = arg_count
		m.OptionalArgs = optional_args
	}
	// single overload, parse args
	g.writeArgChecks(sb, *m.Name, m.Overloads[0], arg_count, optional_args, is_void)

	obj_name := ""
	outType := m.Returns.Name
	for i, arg := range *m.Overloads[0] {
		if i > arg_count {
			break
		}
		isArgTransform, argTransformVal := g.conf.IsArgTransform(*m.Name, *arg.Name)
		isEnum, _ := g.IsTypeEnum(*arg.Type)

		if isArgTransform && strings.Contains(*argTransformVal, "/arg_") {
			g.writeIndent(sb, 2)
			for j, val := range *m.Overloads[0] {
				*argTransformVal = strings.ReplaceAll(*argTransformVal, fmt.Sprintf("/arg_%d/", j), *val.Name)
			}
			sb.WriteString(strings.ReplaceAll(*argTransformVal, "/arg/", fmt.Sprintf("info[%d]", i)))
			// TODO: this might need better handling for class wrappers
		} else if g.isClass(*arg.Type) {
			obj_name = *arg.Name
			// rework below conditional
		} else if IsArgTemplate(arg) && !strings.EqualFold(*arg.Template.Args[0].Name, m.Returns.Name) {
			g.writeIndent(sb, 2)
			invertVal := "false"
			if g.conf.VectorOpts.DimAccessor != "" {
				invertVal = fmt.Sprintf("%s->%s", obj_name, g.conf.VectorOpts.DimAccessor)
			}
			sb.WriteString(fmt.Sprintf("auto %s = jsArrayToVector<%s>(info[%d].As<Napi::Array>(), g_row_major, %s);\n", *arg.Name, *arg.Template.Args[0].Name, i, invertVal))
		} else if !arg.IsPrimitive && !isEnum {
			// TODO: all arg coercions should probably be written in a single place to prevent duplication
			var ptrType string
			if arg.IsPointer {
				ptrType = "*"
			}
			fmt.Printf("TODO: Method %q has unhandled argument: `%s %s%s`\n", *m.Name, *arg.Type, ptrType, *arg.Name)
		}
	}
	if !is_void {
		g.writeIndent(sb, 2)
		_, arrayType, needsCast, _ := PrimitivePtrToTS(m.Returns.Name)

		if !m.Returns.IsPrimitive {
			sb.WriteString(fmt.Sprintf("%s::%s ", *g.NameSpace, outType))
		} else if m.Returns.IsPointer && needsCast != nil && arrayType != "" {
			sb.WriteString(fmt.Sprintf("%s ", arrayType))
		} else {
			sb.WriteString(fmt.Sprintf("%s ", m.Returns.Name))
		}

		if m.Returns.IsPointer {
			sb.WriteByte('*')
		}
		sb.WriteString("_res;\n")
	}

	isReturnTransform, isGrouped, transform := g.conf.IsReturnTransform(m)
	if isReturnTransform {
		if isGrouped {
			g.writeIndent(sb, 2)
			if !is_void {
				sb.WriteString("_res = ")
			}
			sb.WriteString(fmt.Sprintf("%s::%s(", *m.NameSpace, *m.Name))
			for i, arg := range *m.Overloads[0] {
				if i > 0 {
					sb.WriteString(", ")
				}
				isEnum, _ := g.IsTypeEnum(*arg.Type)
				if isEnum {
					sb.WriteString(*arg.Name)
				} else if _, ok := g.conf.TypeMappings[*arg.Type]; ok {
					sb.WriteString(fmt.Sprintf("%s::%s(%s)", *g.NameSpace, *arg.Type, *arg.Name))
				} else if g.isClass(*arg.Type) {
					sb.WriteString(fmt.Sprintf("*(%s)", *arg.Name))
				} else {
					sb.WriteString(*arg.Name)
				}
			}
			sb.WriteString(");\n")
		}
		parsed_transform := strings.ReplaceAll(*transform, "/return/", "_res")
		for i, arg := range *m.Overloads[0] {
			fmtd_arg := ""
			if g.isClass(*arg.Type) {
				fmtd_arg = fmt.Sprintf("*(%s)", *arg.Name)
			} else {
				fmtd_arg = *arg.Name
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
		if !is_void {
			sb.WriteString("_res = ")
		}
		_, arrayType, needsCast, _ := PrimitivePtrToTS(m.Returns.Name)
		if m.Returns.IsPointer && needsCast != nil && arrayType != "" {
			sb.WriteString(fmt.Sprintf("reinterpret_cast<%s *>(", arrayType))
		}
		sb.WriteString(fmt.Sprintf("%s::%s(", *m.NameSpace, *m.Name))
		for i, arg := range *m.Overloads[0] {
			if i > arg_count {
				break
			}
			if i > 0 {
				sb.WriteString(", ")
			}
			if _, ok := g.conf.TypeMappings[*arg.Type]; ok {
				sb.WriteString(fmt.Sprintf("%s::%s(%s)", *g.NameSpace, *arg.Type, *arg.Name))
			} else if g.isClass(*arg.Type) {
				sb.WriteString(fmt.Sprintf("*(%s)", *arg.Name))
			} else {
				sb.WriteString(*arg.Name)
			}
		}
		sb.WriteByte(')')
		if m.Returns.IsPointer && needsCast != nil && arrayType != "" {
			sb.WriteByte(')')
		}
		sb.WriteString(";\n")
	}
	if m.Returns.Name != "void" || isReturnTransform {
		g.writeIndent(sb, 2)
		returnType := m.Returns.Name
		if m.Returns.IsPrimitive && m.Returns.IsPointer {
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
			jsType, isObject := g.CPPTypeToTS(returnType, false)
			if g.conf.TypeHasHandler(returnType) != nil {
				t := g.conf.TypeHasHandler(returnType)
				g.writeIndent(sb, 1)
				sb.WriteString(strings.ReplaceAll(t.Handler, "/val/", "_res"))
				g.writeIndent(sb, 1)
				sb.WriteString(fmt.Sprintf("return %s;\n", t.OutVar))
			} else if isObject && g.isClass(returnType) {
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
	classData := g.getClass(className)
	if f.Name != nil && g.conf.IsFieldWrapped(className, *f.Name) {
		var returnType string
		isVoid := false
		if f.Returns != nil && f.Returns.Name != "void" {
			returnType = f.Returns.Name
		} else {
			isVoid = true
		}
		sb.WriteString("static ")
		if isVoid {
			sb.WriteString("void ")
		} else {
			sb.WriteString("Napi::Value ")
		}
		sb.WriteString(fmt.Sprintf("_%s(const Napi::CallbackInfo& info) {\n", *f.Name))
		g.writeIndent(sb, 1)
		sb.WriteString("Napi::Env env = info.Env();\n")
		argCount := 1
		if f.Args != nil {
			argCount += len(f.Args)
		}
		g.writeArgTypeChecker(sb, *f.Name, "IsExternal", 0, fmt.Sprintf("native `%s` (typeof `Napi::External<%s::%s>`)", stripNameSpace(className), *classData.NameSpace, stripNameSpace(className)), 1, nil, nil, isVoid)
		g.writeIndent(sb, 1)
		sb.WriteString(fmt.Sprintf("%s::%s* _tmp_external = UnExternalize<%s::%s>(info[%d]);\n", *classData.NameSpace, stripNameSpace(className), *classData.NameSpace, stripNameSpace(className), 0))

		if f.Args != nil {
			f.Args = append([]*CPPArg{nil}, f.Args...)
			g.writeArgChecks(sb, *f.Name, &f.Args, argCount, 0, isVoid)
		}

		g.writeIndent(sb, 1)
		if f.Returns != nil && f.Returns.Name != "void" {
			sb.WriteString("auto _res = ")
			returnType = f.Returns.Name
		}
		sb.WriteString(fmt.Sprintf("_tmp_external->%s(", *f.Name))
		if f.Args != nil {
			for i, arg := range f.Args {
				if arg == nil {
					continue
				}
				if i > 1 {
					sb.WriteString(", ")
				}
				if _, ok := g.conf.TypeMappings[*arg.Type]; ok {
					sb.WriteString(fmt.Sprintf("%s::%s(%s)", *g.NameSpace, *arg.Type, *arg.Name))
				} else if g.isClass(*arg.Type) {
					// TODO: check whether expects ptr to Class
					sb.WriteString(fmt.Sprintf("*(%s)", *arg.Name))
				} else {
					sb.WriteString(*arg.Name)
				}
			}
		}
		sb.WriteString(");\n")

		if f.Returns != nil && f.Returns.Name != "void" {
			jsType, isObject := g.CPPTypeToTS(returnType, false)
			if g.conf.TypeHasHandler(returnType) != nil {
				t := g.conf.TypeHasHandler(returnType)
				g.writeIndent(sb, 1)
				sb.WriteString(strings.ReplaceAll(t.Handler, "/val/", "_res"))
				g.writeIndent(sb, 1)
				sb.WriteString(fmt.Sprintf("return %s;\n", t.OutVar))
			} else if isObject && g.isClass(returnType) {
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
					if f.Name != nil && g.conf.IsFieldWrapped(name, *f.Name) {
						g.writeAddonExport(sb, *f.Name)
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
		g.writeAddonExport(sb, *f.Name)
	}

	// write exports for methods requiring pre-processing
	for _, f := range g.ParsedData.Lits {
		g.writeAddonExport(sb, *f.Name)
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
