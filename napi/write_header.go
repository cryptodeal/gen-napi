package napi

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func (g *PackageGenerator) writeHeader(sb *strings.Builder, classes map[string]*CPPClass, methods map[string]*CPPMethod, prepocessedMethods map[string]*CPPMethod) {
	lower_caser := cases.Lower(language.AmericanEnglish)

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
			sb.WriteString(fmt.Sprintf("%s::%s* _%s;\n", *cf.NameSpace, class, lower_caser.String(class)[0:1]+class[1:]))
			g.writeIndent(sb, 2)
			sb.WriteString("static Napi::Function GetClass(Napi::Env);\n\n")
			g.writeIndent(sb, 2)
			sb.WriteString("// methods defined in src, wrapped as class methods\n")
			if cf.FieldDecl != nil {
				for _, f := range *cf.FieldDecl {
					if f.Ident != nil && g.conf.IsFieldWrapped(class, *f.Ident) {
						g.writeIndent(sb, 2)
						sb.WriteString(fmt.Sprintf("Napi::Value %s(const Napi::CallbackInfo&);\n", *f.Ident))
					}
				}
			}
			for _, f := range methods {
				if g.conf.IsMethodWrapped(class, *f.Ident) && strings.EqualFold(class, *f.Returns) {
					g.writeIndent(sb, 2)
					sb.WriteString(fmt.Sprintf("Napi::Value %s(const Napi::CallbackInfo&);\n", *f.Ident))
				}
			}
			for _, f := range prepocessedMethods {
				if g.conf.IsMethodWrapped(class, *f.Ident) && strings.EqualFold(class, *f.Returns) {
					g.writeIndent(sb, 2)
					sb.WriteString(fmt.Sprintf("Napi::Value %s(const Napi::CallbackInfo&);\n", *f.Ident))
				}
			}
			if v, ok := g.conf.ClassOpts[class]; ok {
				for _, f := range v.ForcedMethods {
					g.writeIndent(sb, 2)
					sb.WriteString(fmt.Sprintf("Napi::Value %s(const Napi::CallbackInfo&);\n", f.Name))
				}
			}
			sb.WriteByte('\n')
			g.writeIndent(sb, 1)
			sb.WriteString("private:\n")
			sb.WriteString("};")
		}
	}
}
