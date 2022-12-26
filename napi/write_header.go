package napi

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func (g *PackageGenerator) writeHeader(sb *strings.Builder, classes map[string]*CPPClass, methods map[string]*CPPMethod) {
	lower_caser := cases.Lower(language.AmericanEnglish)

	sb.WriteString("#pragma once\n")
	sb.WriteString("#include <napi.h>\n")
	for class, cf := range classes {
		if cf.FieldDecl != nil {
			sb.WriteString(fmt.Sprintf("class %s : public Napi::ObjectWrap<%s> {\n", class, class))
			g.writeIndent(sb, 1)
			sb.WriteString("public:\n")
			g.writeIndent(sb, 2)
			sb.WriteString(fmt.Sprintf("%s(const Napi::CallbackInfo&);\n", class))
			g.writeIndent(sb, 2)
			sb.WriteString("static Napi::FunctionReference* constructor;\n")
			g.writeIndent(sb, 2)
			sb.WriteString(fmt.Sprintf("%s::%s* _%s;\n\n", *cf.NameSpace, class, lower_caser.String(class)[0:1]+class[1:]))

			for _, f := range methods {
				if g.conf.IsMethodWrapped(class, *f.Ident) {
					g.writeIndent(sb, 2)
					sb.WriteString(fmt.Sprintf("\t\tNapi::Value %s(const Napi::CallbackInfo&);\n", *f.Ident))
				}
			}
			g.writeIndent(sb, 1)
			sb.WriteString("private:\n")
			sb.WriteString("};")
		}
	}
}
