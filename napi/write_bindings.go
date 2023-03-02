package napi

import (
	"fmt"
	"strings"
)

func (g *PackageGenerator) WriteArgTypeCheck(sb *strings.Builder, name string, checker string, idx int, msg string, arrName *string, arg *CPPArg, is_void bool) {
	if arg == nil {
		return
	}
	isArrayItem := arrName != nil
	indents := 1

	if isArrayItem {
		sb.WriteString(fmt.Sprintf("Napi::Array %s = info[%d].As<Napi::Array>();\n", *arrName, idx))
		g.writeIndent(sb, indents)
		sb.WriteString(fmt.Sprintf("size_t len_%s = %s.Length();\n", *arg.Name, *arrName))
		g.writeIndent(sb, indents)
		sb.WriteString(fmt.Sprintf("for (size_t i = 0; i < len_%s; ++i) {\n", *arg.Name))
	}

	if isArrayItem {
		g.writeIndent(sb, indents)
		sb.WriteString(fmt.Sprintf("Napi::Value arrayItem = %s[i];\n", *arrName))
	}

	var error_conditional, err_msg string
	if isArrayItem {
		indents = 2
		error_conditional = fmt.Sprintf("arrayItem.%s()", checker)
		err_msg = fmt.Sprintf("(%q + std::to_string(i) + %q)", fmt.Sprintf("`%s` expects args[%d][", name, idx), fmt.Sprintf("] to be %s", msg))
	} else {
		error_conditional = fmt.Sprintf("!info[%d].%s()", idx, checker)
		err_msg = fmt.Sprintf("`%s` expects args[%d] to be %s", name, idx, msg)
	}
	g.WriteErrorHandler(sb, error_conditional, err_msg, indents, is_void)
}

// WIP: rewriting to clean up logic
func (g *PackageGenerator) writeMethod(sb *strings.Builder, m *CPPMethod) {
	fmt.Printf("method: %s\n", *m.Name)
	arg_helpers := g.GetArgData(m.Overloads[0])
	parsedName := "_" + *m.Name
	return_type := "Napi::Value"
	is_void := m.IsVoid()
	if is_void {
		return_type = "void"
	}
	sb.WriteString(fmt.Sprintf("static %s %s(const Napi::CallbackInfo& info) {\n", return_type, parsedName))
	g.writeIndent(sb, 1)
	sb.WriteString("Napi::Env env = info.Env();\n")
	// TODO: handle transforms (and overloads)
	g.WriteArgChecks(sb, arg_helpers, is_void, *m.Name)
	g.WriteArgGetters(sb, arg_helpers, is_void)
	gen_return_name := g.WriteFnCall(sb, fmt.Sprintf("%s::%s", *m.NameSpace, *m.Name), m.Returns, arg_helpers, is_void)
	g.WriteReturnVal(sb, m.Returns.ParseReturnData(g), is_void, gen_return_name)

	sb.WriteString("}\n\n")
}

func (g *PackageGenerator) WriteClassMethod(sb *strings.Builder, f *CPPFieldDecl, class_name string) {
	if f.Name != nil && !g.conf.IsFieldIgnored(class_name, *f.Name) {
		fmt.Printf("class method: %s\n", *f.Name)
		arg_helpers := g.GetArgData(&f.Args)
		parsedName := "_" + *f.Name
		return_type := "Napi::Value"
		is_void := f.IsVoid()
		if is_void {
			return_type = "void"
		}
		sb.WriteString(fmt.Sprintf("static %s %s(const Napi::CallbackInfo& info) {\n", return_type, parsedName))
		g.writeIndent(sb, 1)
		sb.WriteString("Napi::Env env = info.Env();\n")
		// TODO: handle transforms (and overloads)
		g.WriteArgChecks(sb, arg_helpers, is_void, *f.Name)
		g.WriteArgGetters(sb, arg_helpers, is_void)
		gen_return_name := g.WriteFnCall(sb, fmt.Sprintf("%s->%s", *f.Args[0].Name, *f.Name), f.Returns, arg_helpers, is_void, true)
		g.WriteReturnVal(sb, f.Returns.ParseReturnData(g), is_void, gen_return_name)

		sb.WriteString("}\n\n")
	}
}

func (g *PackageGenerator) writeAddonExport(sb *strings.Builder, name string) {
	g.writeIndent(sb, 1)
	parsedName := ("_" + name)
	sb.WriteString(fmt.Sprintf("exports.Set(Napi::String::New(env, %q), Napi::Function::New(env, %s));\n", parsedName, parsedName))
}

func (g *PackageGenerator) WriteBindingsLogic() string {
	sb := &strings.Builder{}
	g.writeFileSourceHeader(sb, *g.Path)

	// TODO: scan force methods to check if this needs to be included (and check gen methods)
	g.WriteVectorArrayBufferDeleter()

	sb.WriteString("// exported functions\n\n")

	for name, c := range g.ParsedData.Classes {
		if c.Decl != nil {
			g.writeClassDeleter(c, name)
			g.writeClassExternalizer(c, name)
			g.writeClassUnExternalizer()

			// TODO: fix once hashed out updated gen flow

			for _, f := range *c.FieldDecl {
				g.WriteClassMethod(sb, f, name)
			}

			if v, ok := g.conf.ClassOpts[name]; ok && len(v.ForcedMethods) > 0 {
				for _, f := range v.ForcedMethods {
					sb.WriteString(strings.Replace(f.FnBody, f.Name, "_"+f.Name, 1))
					sb.WriteString("\n\n")
				}
			}
		}
	}

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
		// fmt.Printf("Name: %s", f.Name)
		sb.WriteString(fmt.Sprintf("%s\n\n", strings.Replace(f.FnBody, f.Name, "_"+f.Name, 1)))
	}

	// writes NAPI `Init` function (init NAPI exports)
	sb.WriteString("// NAPI exports\n\n")
	sb.WriteString("Napi::Object Init(Napi::Env env, Napi::Object exports) {\n")
	fmt.Println("class count: ", len(g.ParsedData.Classes))

	for name, c := range g.ParsedData.Classes {
		// check if header contained class constructor declaration(s)
		if c.Decl != nil {
			// write exports for wrapped class fields (specified in config)
			if c.FieldDecl != nil {
				for _, f := range *c.FieldDecl {
					if f.Name != nil && !g.conf.IsFieldIgnored(name, *f.Name) {
						g.writeAddonExport(sb, *f.Name)
					}
				}
			}
			// write exports for any optionally forced class methods (specified in config)
			if v, ok := g.conf.ClassOpts[name]; ok {
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

	return sb.String()
}

// makes calls to functions that write bindings
func (g *PackageGenerator) WriteBindings(sb *strings.Builder) {
	g.writeRequiredIncludes(sb)
	g.writeBindingsFrontmatter(sb)

	g.writeGlobalVars(sb)

	gen_logic := g.WriteBindingsLogic()

	// write any helpers functions (non-exported; specified in config)
	g.writeHelpers(sb)
	sb.WriteString(gen_logic)
}
