package napi

import (
	"fmt"
	"path/filepath"
	"strings"
)

func (g *PackageGenerator) writeMethod(sb *strings.Builder, m *CPPMethod) {
	sb.WriteString(fmt.Sprintf("static Napi::Value %s(const Napi::CallbackInfo& info) {\n", *f.Ident))
	g.writeIndent(sb, 1)
	sb.WriteString("Napi::Env env = info.Env();\n")
	if len(m.Overloads) == 1 {
		// single overload, parse args
		g.writeIndent(sb, 1)
		sb.WriteString(fmt.Sprintf("if (info.Length() != %d) {\n", len(*m.Overloads[0])))
		g.writeIndent(sb, 2)
		sb.WriteString(fmt.Sprintf("Napi::TypeError::New(info.Env(), %q).ThrowAsJavaScriptException();\n", fmt.Sprintf("`%s` expects exactly 2 args", *m.Ident)))
		g.writeIndent(sb, 2)
		sb.WriteString("return env.Null();\n")
		g.writeIndent(sb, 1)
		sb.WriteString("}\n")
	}
	g.writeIndent(sb, 1)
	sb.WriteString("return env.Null();\n")
	sb.WriteString("}\n\n")
}

func (g *PackageGenerator) writeBindings(sb *strings.Builder, classes map[string]*CPPClass, methods map[string]*CPPMethod) {
	// lower_caser := cases.Lower(language.AmericanEnglish)

	sb.WriteString(fmt.Sprintf("#include %q\n", filepath.Base(g.conf.ResolvedHeaderOutPath(filepath.Dir(*g.Path)))))
	g.writeBindingsFrontmatter(sb)
	sb.WriteString("using namespace Napi;\n")
	g.writeFileSourceHeader(sb, *g.Path)
	g.writeGlobalVars(sb)
	g.writeHelpers(sb)

	sb.WriteString("// exported functions\n")
	for _, f := range methods {
		g.writeMethod(sb, f)
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
