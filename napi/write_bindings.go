package napi

import (
	"fmt"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func isClass(argType string, classes map[string]*CPPClass) bool {
	_, ok := classes[argType]
	return ok
}

func CPPTypeToTS(t string) (string, bool) {
	switch t {
	case "int", "int8_t", "uint8_t", "signed", "unsigned", "short", "long", "long int", "size_t", "signed long", "signed long int", "unsigned long", "unsigned long int", "long long", "long long int", "signed long long", "signed long long int", "unsigned long long", "unsigned long long int", "long double", "signed char", "unsigned char", "short int", "signed short", "unsigned_short", "signed int", "unsigned int", "unsigned short int", "signed short int", "uint16_t", "uint32_t", "uint64_t", "int16_t", "int32_t", "int64_t", "float", "double":
		return "number", false
	case "string", "std::string", "char":
		return "string", false
	case "bool":
		return "boolean", false
	default:
		return t, true
	}
}

func (g *PackageGenerator) writeMethod(sb *strings.Builder, m *CPPMethod, classes map[string]*CPPClass, wrappedClass *string) {
	lower_caser := cases.Lower(language.AmericanEnglish)
	upper_caser := cases.Upper(language.AmericanEnglish)
	parsedName := "_" + *m.Ident
	if wrappedClass == nil {
		sb.WriteString(fmt.Sprintf("static Napi::Value %s(const Napi::CallbackInfo& info) {\n", parsedName))
	} else {
		sb.WriteString(fmt.Sprintf("Napi::Value %s::%s(const Napi::CallbackInfo& info) {\n", *wrappedClass, *m.Ident))
	}
	g.writeIndent(sb, 1)
	sb.WriteString("Napi::Env env = info.Env();\n")
	hasObject := false
	// if len(m.Overloads) == 1 {
	var expected_count int
	if v, ok := g.conf.MethodReturnTransforms[*m.Ident]; ok {
		expected_count = (strings.Count(strings.Split(v, "\n")[0], ",") + 1)
	} else {
		expected_count = len(*m.Overloads[0])
	}
	// single overload, parse args
	g.writeIndent(sb, 1)
	argCount := expected_count
	if wrappedClass != nil {
		argCount--
	}
	if argCount > 0 {
		sb.WriteString(fmt.Sprintf("if (info.Length() != %d) {\n", argCount))
		g.writeIndent(sb, 2)
		sb.WriteString(fmt.Sprintf("Napi::TypeError::New(info.Env(), %q).ThrowAsJavaScriptException();\n", fmt.Sprintf("`%s` expects exactly %d args", *m.Ident, argCount)))
		g.writeIndent(sb, 2)
		sb.WriteString("return env.Null();\n")
		g.writeIndent(sb, 1)
		sb.WriteString("}\n")
	}
	if expected_count > 0 {
		if v, ok := g.conf.MethodArgCheckTransforms[*m.Ident]; ok {
			lines := strings.Split(v, "\n")
			for _, line := range lines {
				g.writeIndent(sb, 1)
				sb.WriteString(fmt.Sprintf("%s\n", line))
			}
			for i, arg := range *m.Overloads[0] {
				if wrappedClass != nil && i == 0 {
					continue
				}
				if i > expected_count {
					break
				}
				if v, ok := g.conf.MethodArgTransforms[*m.Ident][*arg.Ident]; ok && !strings.Contains(v, "/arg_") {
					g.writeIndent(sb, 1)
					sb.WriteString(strings.ReplaceAll(v, "/arg/", fmt.Sprintf("info[%d]", i)))
				}
			}
		} else {
			for i, arg := range *m.Overloads[0] {
				if wrappedClass != nil && i == 0 {
					continue
				}
				if i > expected_count {
					break
				}
				argIdx := i
				if wrappedClass != nil {
					argIdx--
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
					g.writeIndent(sb, 1)
					sb.WriteString(fmt.Sprintf("if (!info[%d].%s()) {\n", argIdx, napiTypeHandler))
					g.writeIndent(sb, 2)
					sb.WriteString(fmt.Sprintf("Napi::TypeError::New(info.Env(), %q).ThrowAsJavaScriptException();\n", fmt.Sprintf("`%s` expects args[%d] to be typeof `%s`", *m.Ident, argIdx, jsTypeEquivalent)))
					g.writeIndent(sb, 2)
					sb.WriteString("return env.Null();\n")
					g.writeIndent(sb, 1)
					sb.WriteString("}\n")
					if _, ok := g.conf.MethodArgTransforms[*m.Ident][*arg.Ident]; !ok {
						g.writeIndent(sb, 1)
						sb.WriteString(fmt.Sprintf("%s %s = ", *arg.Type, *arg.Ident))
						if needsCast != nil {
							sb.WriteString(fmt.Sprintf("static_cast<%s>(", *needsCast))
						}
						sb.WriteString(fmt.Sprintf("info[%d].As<Napi::%s>().%s()", argIdx, strings.ReplaceAll(napiTypeHandler, "Is", ""), valGetter))
						if needsCast != nil {
							sb.WriteByte(')')
						}
						sb.WriteString(";\n")
					}
				} else if isClass(*arg.Type, classes) {
					hasObject = true
					g.writeIndent(sb, 1)
					sb.WriteString(fmt.Sprintf("if (!info[%d].IsObject()) {\n", argIdx))
					g.writeIndent(sb, 2)
					sb.WriteString(fmt.Sprintf("Napi::TypeError::New(info.Env(), %q).ThrowAsJavaScriptException();\n", fmt.Sprintf("`%s` expects args[%d] to be instanceof `%s`", *m.Ident, argIdx, *arg.Type)))
					g.writeIndent(sb, 2)
					sb.WriteString("return env.Null();\n")
					g.writeIndent(sb, 1)
					sb.WriteString("}\n")
					if _, ok := g.conf.MethodArgTransforms[*m.Ident][*arg.Ident]; !ok {
						g.writeIndent(sb, 1)
						sb.WriteString(fmt.Sprintf("Napi::Object %s_obj = info[%d].As<Napi::Object>();\n", *arg.Ident, argIdx))
					}
				} else if strings.Contains(*arg.Type, "std::vector") {
					argType := *arg.Type
					type_test := argType[strings.Index(argType, "<")+1 : strings.Index(argType, ">")]
					g.writeIndent(sb, 1)
					sb.WriteString(fmt.Sprintf("if (!info[%d].IsArray()) {\n", argIdx))
					g.writeIndent(sb, 2)
					tsType, isObject := CPPTypeToTS(type_test)
					sb.WriteString(fmt.Sprintf("Napi::TypeError::New(info.Env(), %q).ThrowAsJavaScriptException();\n", fmt.Sprintf("`%s` expects args[%d] to be typeof `%s[]`", *m.Ident, argIdx, tsType)))
					g.writeIndent(sb, 2)
					sb.WriteString("return env.Null();\n")
					g.writeIndent(sb, 1)
					sb.WriteString("}\n")
					g.writeIndent(sb, 1)
					sb.WriteString(fmt.Sprintf("int len_%s = info[%d].As<Napi::Array>().Length();\n", *arg.Ident, argIdx))
					g.writeIndent(sb, 1)
					sb.WriteString(fmt.Sprintf("for (auto i = 0; i < len_%s; ++i) {\n", *arg.Ident))
					g.writeIndent(sb, 2)
					if isObject {
						sb.WriteString(fmt.Sprintf("if (!info[%d].As<Napi::Array>().Get(i).IsObject()) {\n", argIdx))
					} else {
						sb.WriteString(fmt.Sprintf("if (!info[%d].As<Napi::Array>().Get(i).Is%s()) {\n", argIdx, upper_caser.String(tsType[0:1])+tsType[1:]))
					}
					g.writeIndent(sb, 3)
					sb.WriteString(fmt.Sprintf("Napi::TypeError::New(info.Env(), (%q + std::to_string(i) + %q)).ThrowAsJavaScriptException();\n", fmt.Sprintf("`%s` expects args[%d][", *m.Ident, argIdx), fmt.Sprintf("] to be typeof `%s`", tsType)))
					g.writeIndent(sb, 3)
					sb.WriteString("return env.Null();\n")
					g.writeIndent(sb, 2)
					sb.WriteString("}\n")
					g.writeIndent(sb, 1)
					sb.WriteString("}\n")
					if wrappedClass != nil {
						g.writeIndent(sb, 2)
						tmpType := *arg.Type
						sb.WriteString(fmt.Sprintf("auto axes = jsArrayArg<%s>(info[%d].As<Napi::Array>(), g_row_major, this->_%s->ndim(), env);\n", tmpType[strings.Index(*arg.Type, "<")+1:strings.Index(*arg.Type, ">")], argIdx, lower_caser.String(*wrappedClass)))
					}
				} else if v, ok := g.conf.TypeMappings[*arg.Type]; ok {
					g.writeIndent(sb, 1)
					if strings.Contains(v.TSType, "Array") || strings.Contains(v.TSType, "[]") {
						sb.WriteString(fmt.Sprintf("if (!info[%d].IsArray()) {\n", argIdx))
					} else if strings.Contains(v.TSType, "any") || strings.Contains(v.TSType, "object") || strings.Contains(v.TSType, "Record<") || strings.Contains(v.TSType, "Map<") {
						sb.WriteString(fmt.Sprintf("if (!info[%d].IsObject()) {\n", argIdx))
					} else if strings.Contains(v.TSType, "string") {
						sb.WriteString(fmt.Sprintf("if (!info[%d].IsString()) {\n", argIdx))
					} else if strings.Contains(v.TSType, "number") {
						sb.WriteString(fmt.Sprintf("if (!info[%d].IsNumber()) {\n", argIdx))
					}
					g.writeIndent(sb, 2)
					sb.WriteString(fmt.Sprintf("Napi::TypeError::New(info.Env(), %q).ThrowAsJavaScriptException();\n", fmt.Sprintf("`%s` expects args[%d] to be typeof `%s`", *m.Ident, argIdx, v.TSType)))
					g.writeIndent(sb, 2)
					sb.WriteString("return env.Null();\n")
					g.writeIndent(sb, 1)
					sb.WriteString("}\n")
				}
				if v, ok := g.conf.MethodArgTransforms[*m.Ident][*arg.Ident]; ok {
					if wrappedClass != nil && strings.Contains(v, "/arg_0/") {
						v = strings.ReplaceAll(v, "/arg_0/", "this")
					}
					if !strings.Contains(v, "/arg_") {
						g.writeIndent(sb, 1)
						sb.WriteString(strings.ReplaceAll(v, "/arg/", fmt.Sprintf("info[%d]", argIdx)))
					}
				}
			}
		}
	}

	// handle casting objects to `Napi::ObjectWrap<Class>`
	if hasObject {
		g.writeIndent(sb, 1)
		sb.WriteString("if (")
		for i, arg := range *m.Overloads[0] {
			if wrappedClass != nil && i == 0 {
				continue
			}
			if isClass(*arg.Type, classes) {
				if (wrappedClass == nil && i > 0) || i > 1 {
					sb.WriteString(" && ")
				}
				sb.WriteString(fmt.Sprintf("%s_obj.InstanceOf(%s::constructor->Value())", *arg.Ident, *arg.Type))
			}
		}
		sb.WriteString(") {\n")
		obj_name := ""
		obj_type := ""
		for i, arg := range *m.Overloads[0] {
			argIdx := i
			if wrappedClass != nil {
				argIdx--
			}
			if i > expected_count {
				break
			}
			if v, ok := g.conf.MethodArgTransforms[*m.Ident][*arg.Ident]; ok && strings.Contains(v, "/arg_") {
				g.writeIndent(sb, 2)
				if strings.Contains(v, "/arg_") {
					for j, val := range *m.Overloads[0] {
						if wrappedClass != nil && j == 0 {
							v = strings.ReplaceAll(v, "/arg_%d/", "this")
						} else {
							v = strings.ReplaceAll(v, fmt.Sprintf("/arg_%d/", j), *val.Ident)
						}
					}
				}
				// TODO: this might need better handling for class wrappers
				sb.WriteString(strings.ReplaceAll(v, "/arg/", fmt.Sprintf("info[%d]", argIdx)))
			} else if isClass(*arg.Type, classes) {
				if wrappedClass != nil && i == 0 {
					continue
				}
				obj_name = *arg.Ident
				obj_type = *arg.Type
				g.writeIndent(sb, 2)
				sb.WriteString(fmt.Sprintf("%s* %s = Napi::ObjectWrap<%s>::Unwrap(%s_obj);\n", *arg.Type, *arg.Ident, *arg.Type, *arg.Ident))
			} else if strings.Contains(*arg.Type, "std::vector") {
				g.writeIndent(sb, 2)
				tmpType := *arg.Type
				sb.WriteString(fmt.Sprintf("auto axes = jsArrayArg<%s>(info[%d].As<Napi::Array>(), g_row_major, %s->_%s->ndim(), env);\n", tmpType[strings.Index(*arg.Type, "<")+1:strings.Index(*arg.Type, ">")], i, obj_name, lower_caser.String(obj_type)))
			} else {
				fmt.Println("TODO: handle type ", *arg.Type)
			}
		}
		g.writeIndent(sb, 2)
		sb.WriteString(fmt.Sprintf("%s::%s _res;\n", *g.NameSpace, obj_type))
		if v, ok := g.conf.MethodReturnTransforms[*m.Ident]; ok {
			parsed_transform := strings.ReplaceAll(v, "/return/", "_res")
			for i, arg := range *m.Overloads[0] {
				fmtd_arg := ""
				if isClass(*arg.Type, classes) {
					if wrappedClass != nil && i == 0 {
						fmtd_arg = fmt.Sprintf("*(this->_%s)", lower_caser.String(*arg.Type))
					} else {
						fmtd_arg = fmt.Sprintf("*(%s->_%s)", *arg.Ident, lower_caser.String(*arg.Type))
					}
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
					if wrappedClass != nil && i == 0 {
						sb.WriteString(fmt.Sprintf("*(this->_%s)", lower_caser.String(*arg.Type)))

					} else {
						sb.WriteString(fmt.Sprintf("*(%s->_%s)", *arg.Ident, lower_caser.String(*arg.Type)))

					}
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
		sb.WriteString(fmt.Sprintf("auto* out = new %s::%s(_res);\n", *g.NameSpace, obj_type))
		g.writeIndent(sb, 2)
		sb.WriteString(fmt.Sprintf("auto wrapped = Napi::External<%s::%s>::New(env, out);\n", *g.NameSpace, obj_type))
		g.writeIndent(sb, 2)
		sb.WriteString(fmt.Sprintf("Napi::Value wrapped_out = %s::constructor->New({wrapped});\n", obj_type))
		g.writeIndent(sb, 2)
		sb.WriteString("return wrapped_out;\n")
		g.writeIndent(sb, 1)
		sb.WriteString("}\n")
		g.writeIndent(sb, 1)
		sb.WriteString("return env.Null();\n")
	} else {
		g.writeIndent(sb, 1)
		sb.WriteString(fmt.Sprintf("%s::%s _res;\n", *g.NameSpace, *m.Returns))
		if v, ok := g.conf.MethodReturnTransforms[*m.Ident]; ok {
			parsed_transform := strings.ReplaceAll(v, "/return/", "_res")
			for i, arg := range *m.Overloads[0] {
				fmtd_arg := ""
				if isClass(*arg.Type, classes) {
					if wrappedClass != nil && i == 0 {
						fmtd_arg = fmt.Sprintf("*(this->_%s)", lower_caser.String(*arg.Type))
					} else {
						fmtd_arg = fmt.Sprintf("*(%s->_%s)", *arg.Ident, lower_caser.String(*arg.Type))
					}
				} else {
					fmtd_arg = *arg.Ident
				}
				parsed_transform = strings.ReplaceAll(parsed_transform, fmt.Sprintf("/arg_%d/", i), fmtd_arg)
			}
			transformed_lines := strings.Split(parsed_transform, "\n")
			length := len(transformed_lines)
			for i, line := range transformed_lines {
				g.writeIndent(sb, 1)
				if i == length-1 {
					sb.WriteString(line)
				} else {
					sb.WriteString(fmt.Sprintf("%s\n", line))
				}
			}
		} else {
			g.writeIndent(sb, 1)
			sb.WriteString(fmt.Sprintf("_res = %s::%s(", *g.NameSpace, *m.Ident))
			for i, arg := range *m.Overloads[0] {
				if i > 0 {
					sb.WriteString(", ")
				}
				if _, ok := g.conf.TypeMappings[*arg.Type]; ok {
					sb.WriteString(fmt.Sprintf("%s::%s(%s)", *g.NameSpace, *arg.Type, *arg.Ident))
				} else if isClass(*arg.Type, classes) {
					if wrappedClass != nil && i == 0 {
						sb.WriteString(fmt.Sprintf("*(this->_%s)", lower_caser.String(*arg.Type)))
					} else {
						sb.WriteString(fmt.Sprintf("*(%s->_%s)", *arg.Ident, lower_caser.String(*arg.Type)))
					}
				} else {
					sb.WriteString(*arg.Ident)
				}
			}
			sb.WriteString(");\n")
		}
		if v, ok := g.conf.GlobalTypeOutTransforms[*m.Returns]; ok {
			g.writeIndent(sb, 1)
			sb.WriteString(strings.ReplaceAll(v, "/return/", "_res"))
		}
		g.writeIndent(sb, 1)
		sb.WriteString(fmt.Sprintf("auto* out = new %s::%s(_res);\n", *g.NameSpace, *m.Returns))
		g.writeIndent(sb, 1)
		sb.WriteString(fmt.Sprintf("auto wrapped = Napi::External<%s::%s>::New(env, out);\n", *g.NameSpace, *m.Returns))
		g.writeIndent(sb, 1)
		sb.WriteString(fmt.Sprintf("Napi::Value wrapped_out = %s::constructor->New({wrapped});\n", *m.Returns))
		g.writeIndent(sb, 1)
		sb.WriteString("return wrapped_out;\n")
	}
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
	lower_caser := cases.Lower(language.AmericanEnglish)
	upper_caser := cases.Upper(language.AmericanEnglish)

	var returnType string
	isVoid := false
	if f.Returns != nil && *f.Returns.FullType != "void" {
		returnType = *f.Returns.FullType
	} else {
		isVoid = true
	}
	if isVoid {
		sb.WriteString("void ")
	} else {
		sb.WriteString("Napi::Value ")
	}
	sb.WriteString(fmt.Sprintf("%s::%s(const Napi::CallbackInfo& info) {\n", className, *f.Ident))
	g.writeIndent(sb, 1)
	sb.WriteString("Napi::Env env = info.Env();\n")
	if f.Args != nil {
		argCount := len(*f.Args)
		if argCount > 0 {
			g.writeIndent(sb, 1)
			sb.WriteString(fmt.Sprintf("if (info.Length() != %d) {\n", argCount))
			g.writeIndent(sb, 2)
			sb.WriteString(fmt.Sprintf("Napi::TypeError::New(info.Env(), %q).ThrowAsJavaScriptException();\n", fmt.Sprintf("`%s` expects exactly %d args", *f.Ident, argCount)))
			g.writeIndent(sb, 2)
			sb.WriteString("return env.Null();\n")
			g.writeIndent(sb, 1)
			sb.WriteString("}\n")
		}
	}

	if f.Args != nil {
		for i, arg := range *f.Args {
			typeHandler, isObject := CPPTypeToTS(*arg.Type)
			if v, ok := g.conf.TypeMappings[*arg.Type]; ok {
				g.writeIndent(sb, 1)
				sb.WriteString(fmt.Sprintf("if (!info[%d].Is%s()) {\n", i, v.NapiType))
				g.writeIndent(sb, 2)
				sb.WriteString(fmt.Sprintf("Napi::TypeError::New(info.Env(), %q).ThrowAsJavaScriptException();\n", fmt.Sprintf("`%s` expects args[%d] to be typeof `%s`", *f.Ident, i, typeHandler)))
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
	sb.WriteString(fmt.Sprintf("this->_%s->%s(", lower_caser.String(className), *f.Ident))
	if f.Args != nil {
		for i, arg := range *f.Args {
			if i > 0 {
				sb.WriteString(", ")
			}
			if _, ok := g.conf.TypeMappings[*arg.Type]; ok {
				sb.WriteString(fmt.Sprintf("%s::%s(%s)", *g.NameSpace, *arg.Type, *arg.Ident))
			} else if isClass(*arg.Type, classes) {
				sb.WriteString(fmt.Sprintf("*(%s->_%s)", *arg.Ident, lower_caser.String(*arg.Type)))
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
			sb.WriteString(fmt.Sprintf("auto _wrapped = Napi::External<%s::%s>::New(env, out);\n", *g.NameSpace, returnType))
			g.writeIndent(sb, 1)
			sb.WriteString(fmt.Sprintf("Napi::Value wrapped_out = %s::constructor->New({_wrapped});\n", returnType))
			g.writeIndent(sb, 1)
			sb.WriteString("return wrapped_out;\n")
		} else {
			napiHandler := upper_caser.String(jsType[0:1]) + jsType[1:]
			sb.WriteString(fmt.Sprintf("return Napi::%s::New(env, %s);\n", napiHandler, "_res"))
		}
	}
	sb.WriteString("}\n\n")
}

func (g *PackageGenerator) writeClass(sb *strings.Builder, class *CPPClass, classes map[string]*CPPClass, name string, methods map[string]*CPPMethod, processedMethods map[string]*CPPMethod) {

	if class.FieldDecl != nil {
		for _, f := range *class.FieldDecl {
			if f.Ident != nil {
				fmt.Println(*f.Ident)
			}
		}
	}
	// write class constructor (passed in as config option)
	sb.WriteString(fmt.Sprintf("// %q class constructor\n", name))
	sb.WriteString(g.conf.ClassOpts[name].Constructor)
	sb.WriteString(fmt.Sprintf("// exported %q class methods\n", name))
	for _, f := range methods {
		if g.conf.IsMethodWrapped(name, *f.Ident) {
			g.writeMethod(sb, f, classes, &name)
		}
	}
	for _, f := range processedMethods {
		if g.conf.IsMethodWrapped(name, *f.Ident) {
			g.writeMethod(sb, f, classes, &name)
		}
	}

	if class.FieldDecl != nil {
		for _, f := range *class.FieldDecl {
			if f.Ident != nil && g.conf.IsFieldWrapped(name, *f.Ident) {
				g.writeClassField(sb, f, name, classes)
			}
		}
	}

	if v, ok := g.conf.ClassOpts[name]; ok {
		for _, f := range v.ForcedMethods {
			sb.WriteString(fmt.Sprintf("%s\n", f.FnBody))
		}
	}

	sb.WriteString(fmt.Sprintf("Napi::FunctionReference* %s::constructor;\n", name))
	sb.WriteString(fmt.Sprintf("Napi::Function %s::GetClass(Napi::Env env) {\n", name))
	g.writeIndent(sb, 1)
	sb.WriteString(fmt.Sprintf("Napi::Function func = DefineClass(env, %q, {\n", name))
	for _, f := range methods {
		if g.conf.IsMethodWrapped(name, *f.Ident) {
			g.writeIndent(sb, 2)
			sb.WriteString(fmt.Sprintf("%s::InstanceMethod(%q, &%s::%s),\n", name, *f.Ident, name, *f.Ident))
		}
	}
	for _, f := range processedMethods {
		if g.conf.IsMethodWrapped(name, *f.Ident) {
			g.writeIndent(sb, 2)
			sb.WriteString(fmt.Sprintf("%s::InstanceMethod(%q, &%s::%s),\n", name, *f.Ident, name, *f.Ident))
		}
	}
	if class.FieldDecl != nil {
		for _, f := range *class.FieldDecl {
			if f.Ident != nil && g.conf.IsFieldWrapped(name, *f.Ident) {
				g.writeIndent(sb, 2)
				sb.WriteString(fmt.Sprintf("%s::InstanceMethod(%q, &%s::%s),\n", name, *f.Ident, name, *f.Ident))
			}
		}
	}
	if v, ok := g.conf.ClassOpts[name]; ok {
		for _, f := range v.ForcedMethods {
			g.writeIndent(sb, 2)
			sb.WriteString(fmt.Sprintf("%s::InstanceMethod(%q, &%s::%s),\n", name, f.Name, name, f.Name))
		}
	}

	g.writeIndent(sb, 1)
	sb.WriteString("});\n")
	g.writeIndent(sb, 1)
	sb.WriteString("constructor = new Napi::FunctionReference();\n")
	g.writeIndent(sb, 1)
	sb.WriteString("*constructor = Napi::Persistent(func);\n")
	g.writeIndent(sb, 1)
	sb.WriteString("return func;\n")
	sb.WriteString("}\n\n")
}

// makes calls to functions that write bindings
func (g *PackageGenerator) writeBindings(sb *strings.Builder, classes map[string]*CPPClass, methods map[string]*CPPMethod, processedMethods map[string]*CPPMethod) {
	sb.WriteString(fmt.Sprintf("#include %q\n", filepath.Base(g.conf.ResolvedHeaderOutPath(filepath.Dir(*g.Path)))))
	g.writeBindingsFrontmatter(sb)
	sb.WriteString("using namespace Napi;\n")
	g.writeFileSourceHeader(sb, *g.Path)
	g.writeGlobalVars(sb)
	g.writeHelpers(sb)

	sb.WriteString("// exported functions\n")
	for _, f := range methods {
		g.writeMethod(sb, f, classes, nil)
	}

	for _, f := range processedMethods {
		g.writeMethod(sb, f, classes, nil)
	}

	for name, c := range classes {
		if c.Decl != nil {
			g.writeClass(sb, c, classes, name, methods, processedMethods)
		}
	}

	for _, f := range g.conf.GlobalForcedMethods {
		sb.WriteString(fmt.Sprintf("%s\n", strings.Replace(f.FnBody, f.Name, "_"+f.Name, 1)))
	}
	// writes NAPI `Init` function (init NAPI exports)
	sb.WriteString("// NAPI exports\n")
	sb.WriteString("Napi::Object Init(Napi::Env env, Napi::Object exports) {\n")
	for name, c := range classes {
		if c.Decl != nil {
			g.writeIndent(sb, 1)
			parsedName := ("_" + name)
			sb.WriteString(fmt.Sprintf("exports.Set(Napi::String::New(env, %q), %s::GetClass(env));\n", parsedName, name))
		}
	}
	for _, f := range methods {
		g.writeIndent(sb, 1)
		parsedName := ("_" + *f.Ident)
		sb.WriteString(fmt.Sprintf("exports.Set(Napi::String::New(env, %q), Napi::Function::New(env, %s));\n", parsedName, parsedName))
	}
	for _, f := range processedMethods {
		g.writeIndent(sb, 1)
		parsedName := ("_" + *f.Ident)
		sb.WriteString(fmt.Sprintf("exports.Set(Napi::String::New(env, %q), Napi::Function::New(env, %s));\n", parsedName, parsedName))
	}
	for _, f := range g.conf.GlobalForcedMethods {
		g.writeIndent(sb, 1)
		sb.WriteString(fmt.Sprintf("exports.Set(Napi::String::New(env, %q), Napi::Function::New(env, %s));\n", "_"+f.Name, f.Name))
	}
	g.writeIndent(sb, 1)
	sb.WriteString("return exports;\n")
	sb.WriteString("}\n\n")
	sb.WriteString("NODE_API_MODULE(addon, Init)\n")
}
