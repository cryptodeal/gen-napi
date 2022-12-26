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

func (g *PackageGenerator) writeMethod(sb *strings.Builder, m *CPPMethod, classes map[string]*CPPClass) {
	lower_caser := cases.Lower(language.AmericanEnglish)
	sb.WriteString(fmt.Sprintf("static Napi::Value %s(const Napi::CallbackInfo& info) {\n", *m.Ident))
	g.writeIndent(sb, 1)
	sb.WriteString("Napi::Env env = info.Env();\n")
	hasObject := false
	if len(m.Overloads) == 1 {
		expected_count := len(*m.Overloads[0])
		// single overload, parse args
		g.writeIndent(sb, 1)
		sb.WriteString(fmt.Sprintf("if (info.Length() != %d) {\n", expected_count))
		g.writeIndent(sb, 2)
		sb.WriteString(fmt.Sprintf("Napi::TypeError::New(info.Env(), %q).ThrowAsJavaScriptException();\n", fmt.Sprintf("`%s` expects exactly %d args", *m.Ident, expected_count)))
		g.writeIndent(sb, 2)
		sb.WriteString("return env.Null();\n")
		g.writeIndent(sb, 1)
		sb.WriteString("}\n")
		if expected_count > 0 {
			for i, arg := range *m.Overloads[0] {
				if arg.IsPrimitive {
					napiTypeHandler := ""
					jsTypeEquivalent := ""
					valGetter := "Value"
					var needsCast *string
					switch strings.ReplaceAll(*arg.Type, "const ", "") {
					case "float":
						napiTypeHandler = "IsNumber"
						jsTypeEquivalent = "number"
						valGetter = "FloatValue"

					case "double":
						napiTypeHandler = "IsNumber"
						jsTypeEquivalent = "number"
						valGetter = "DoubleValue"

					case "int32_t":
						napiTypeHandler = "IsNumber"
						jsTypeEquivalent = "number"
						valGetter = "Int32Value"

					case "int64_t":
						napiTypeHandler = "IsNumber"
						jsTypeEquivalent = "number"
						valGetter = "Int64Value"

					case "int", "signed", "unsigned", "unsigned int", "char", "unsigned char", "long long", "unsigned long long", "short", "unsigned short", "int8_t", "uint8_t", "int16_t", "uint16_t", "uint32_t", "uint64_t", "size_t":
						napiTypeHandler = "IsNumber"
						jsTypeEquivalent = "number"
						needsCast = arg.Type
					case "bool":
						napiTypeHandler = "IsBoolean"
						jsTypeEquivalent = "boolean"
					}
					g.writeIndent(sb, 1)
					sb.WriteString(fmt.Sprintf("if (!info[%d].%s()) {\n", i, napiTypeHandler))
					g.writeIndent(sb, 2)
					sb.WriteString(fmt.Sprintf("Napi::TypeError::New(info.Env(), %q).ThrowAsJavaScriptException();\n", fmt.Sprintf("`%s` expects args[%d] to be typeof `%s`", *m.Ident, i, jsTypeEquivalent)))
					g.writeIndent(sb, 2)
					sb.WriteString("return env.Null();\n")
					g.writeIndent(sb, 1)
					sb.WriteString("}\n")
					g.writeIndent(sb, 1)
					if needsCast != nil {
						sb.WriteString(fmt.Sprintf("static_cast<%s>(", *needsCast))
					}
					sb.WriteString(fmt.Sprintf("%s %s = info[%d].As<Napi::%s>().%s()", *arg.Type, *arg.Ident, i, strings.ReplaceAll(napiTypeHandler, "Is", ""), valGetter))
					if needsCast != nil {
						sb.WriteByte(')')
					}
					sb.WriteString(";\n")
				} else if isClass(*arg.Type, classes) {
					hasObject = true
					g.writeIndent(sb, 1)
					sb.WriteString(fmt.Sprintf("if (!info[%d].IsObject()) {\n", i))
					g.writeIndent(sb, 2)
					sb.WriteString(fmt.Sprintf("Napi::TypeError::New(info.Env(), %q).ThrowAsJavaScriptException();\n", fmt.Sprintf("`%s` expects args[%d] to be instanceof `%s`", *m.Ident, i, *arg.Type)))
					g.writeIndent(sb, 2)
					sb.WriteString("return env.Null();\n")
					g.writeIndent(sb, 1)
					sb.WriteString("}\n")
					g.writeIndent(sb, 1)
					sb.WriteString(fmt.Sprintf("Napi::Object %s_obj = info[0].As<Napi::Object>();\n", *arg.Ident))
				} else if strings.Contains(*arg.Type, "std::vector") {
					g.writeIndent(sb, 1)
					sb.WriteString(fmt.Sprintf("if (!info[%d].IsArray()) {\n", i))
					g.writeIndent(sb, 2)
					sb.WriteString(fmt.Sprintf("Napi::TypeError::New(info.Env(), %q).ThrowAsJavaScriptException();\n", fmt.Sprintf("`%s` expects args[%d] to be typeof `number[]`", *m.Ident, i)))
					g.writeIndent(sb, 2)
					sb.WriteString("return env.Null();\n")
					g.writeIndent(sb, 1)
					sb.WriteString("}\n")
				}
			}
		}

		// handle casting objects to `Napi::ObjectWrap<Class>`
		if hasObject {
			g.writeIndent(sb, 1)
			sb.WriteString("if (")
			for i, arg := range *m.Overloads[0] {
				if isClass(*arg.Type, classes) {
					if i > 0 {
						sb.WriteString(" && ")
					}
					sb.WriteString(fmt.Sprintf("%s_obj.InstanceOf(%s::constructor->Value())", *arg.Ident, *arg.Type))
				}
			}
			sb.WriteString(") {\n")
			obj_name := ""
			obj_type := ""
			namespace := ""
			for i, arg := range *m.Overloads[0] {
				if isClass(*arg.Type, classes) {
					obj_name = *arg.Ident
					obj_type = *arg.Type
					namespace = *classes[*arg.Type].NameSpace
					g.writeIndent(sb, 2)
					sb.WriteString(fmt.Sprintf("%s* %s = Napi::ObjectWrap<%s>::Unwrap(%s_obj);\n", *arg.Type, *arg.Ident, *arg.Type, *arg.Ident))
				} else if strings.Contains(*arg.Type, "std::vector") {
					g.writeIndent(sb, 2)
					tmpType := *arg.Type
					sb.WriteString(fmt.Sprintf("auto axes = jsArrayArg<%s>(info[%d].As<Napi::Array>(), g_row_major, %s->_%s->ndim(), env);\n", tmpType[strings.Index(*arg.Type, "<")+1:strings.Index(*arg.Type, ">")], i, obj_name, lower_caser.String(obj_type)))
				}
			}
			g.writeIndent(sb, 2)
			sb.WriteString(fmt.Sprintf("%s::%s res;\n", namespace, obj_type))
			g.writeIndent(sb, 2)
			sb.WriteString(fmt.Sprintf("res = %s::%s(", namespace, *m.Ident))
			for i, arg := range *m.Overloads[0] {
				if i > 0 {
					sb.WriteString(", ")
				}
				if isClass(*arg.Type, classes) {
					sb.WriteString(fmt.Sprintf("*(%s->_%s)", *arg.Ident, lower_caser.String(*arg.Type)))
				} else {
					sb.WriteString(*arg.Ident)
				}
			}
			sb.WriteString(");\n")
			g.writeIndent(sb, 2)
			sb.WriteString("g_bytes_used += res.bytes();\n")
			g.writeIndent(sb, 2)
			sb.WriteString(fmt.Sprintf("auto* tensor = new %s::%s(res);\n", namespace, obj_type))
			g.writeIndent(sb, 2)
			sb.WriteString(fmt.Sprintf("auto wrapped = Napi::External<%s::%s>::New(env, tensor);\n", namespace, obj_type))
			g.writeIndent(sb, 2)
			sb.WriteString(fmt.Sprintf("Napi::Value wrapped_out = %s::constructor->New({wrapped});\n", obj_type))
			g.writeIndent(sb, 2)
			sb.WriteString("return wrapped_out;\n")
			g.writeIndent(sb, 1)
			sb.WriteString("}\n")
			g.writeIndent(sb, 1)
			sb.WriteString("return env.Null();\n")
		}
	}
	// TODO: handle cases w multiple overloads
	sb.WriteString("}\n\n")
}

func (g *PackageGenerator) writeBindings(sb *strings.Builder, classes map[string]*CPPClass, methods map[string]*CPPMethod) {
	sb.WriteString(fmt.Sprintf("#include %q\n", filepath.Base(g.conf.ResolvedHeaderOutPath(filepath.Dir(*g.Path)))))
	g.writeBindingsFrontmatter(sb)
	sb.WriteString("using namespace Napi;\n")
	g.writeFileSourceHeader(sb, *g.Path)
	g.writeGlobalVars(sb)
	g.writeHelpers(sb)

	sb.WriteString("// exported functions\n")
	for _, f := range methods {
		g.writeMethod(sb, f, classes)
	}

	sb.WriteString("// NAPI exports\n")
	sb.WriteString("Napi::Object Init(Napi::Env env, Napi::Object exports) {\n")
	for _, f := range methods {
		g.writeIndent(sb, 1)
		sb.WriteString(fmt.Sprintf("exports.Set(Napi::String::New(env, %q), Napi::Function::New(env, %s));\n", ("_" + *f.Ident), *f.Ident))
	}
	sb.WriteString("}\n\n")
	sb.WriteString("NODE_API_MODULE(addon, Init)\n")
}
