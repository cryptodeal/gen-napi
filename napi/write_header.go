package napi

import (
	"fmt"
	"strings"
)

func (g *PackageGenerator) writeHeader(sb *strings.Builder, classes map[string]*CPPClass, methods map[string]*CPPMethod, prepocessedMethods map[string]*CPPMethod) {
	sb.WriteString("#pragma once\n")
	sb.WriteString("#include <napi.h>\n")
	g.writeHeaderFrontmatter(sb)
	g.writeFileSourceHeader(sb, *g.Path)

	for class, cf := range classes {
		if cf.Decl != nil {
			sb.WriteString(fmt.Sprintf("class %s : public Napi::ObjectWrap<%s> {\n", class, class))
			g.writeIndent(sb, 1)
			sb.WriteString("public:\n")
			g.writeIndent(sb, 2)
			sb.WriteString(fmt.Sprintf("%s(const Napi::CallbackInfo&);\n", class))
			g.writeIndent(sb, 2)
			sb.WriteString("static Napi::FunctionReference* constructor;\n")
			g.writeIndent(sb, 2)
			sb.WriteString(fmt.Sprintf("%s::%s* _%s;\n", *cf.NameSpace, class, g.casers.lower.String(class)[0:1]+class[1:]))
			g.writeIndent(sb, 2)
			sb.WriteString("static Napi::Function GetClass(Napi::Env);\n")
			/* TODO: fix finalize method
			g.writeIndent(sb, 2)
			sb.WriteString("virtual void Finalize(Napi::Env env);\n\n")
			*/
			g.writeIndent(sb, 2)
			sb.WriteString("// methods defined in src, wrapped as class methods\n")
			if cf.FieldDecl != nil {
				for _, f := range *cf.FieldDecl {
					if f.Name != nil && g.conf.IsFieldWrapped(class, *f.Name) {
						g.writeIndent(sb, 2)
						sb.WriteString(fmt.Sprintf("Napi::Value %s(const Napi::CallbackInfo&);\n", *f.Name))
					}
				}
			}
			if v, ok := g.conf.ClassOpts[class]; ok {
				for _, f := range v.ForcedMethods {
					g.writeIndent(sb, 2)
					if f.IsVoid {
						sb.WriteString(fmt.Sprintf("void %s(const Napi::CallbackInfo&);\n", f.Name))
					} else {
						sb.WriteString(fmt.Sprintf("Napi::Value %s(const Napi::CallbackInfo&);\n", f.Name))
					}
				}
			}
			sb.WriteByte('\n')
			g.writeIndent(sb, 1)
			sb.WriteString("private:\n")
			sb.WriteString("};")
		}
	}
}
