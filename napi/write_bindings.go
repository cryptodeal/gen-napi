package napi

import (
	"fmt"
	"path"
	"strings"
)

func (g *PackageGenerator) writeBindings(sb *strings.Builder, classes map[string]*CPPClass, methods map[string]*CPPMethod) {
	// lower_caser := cases.Lower(language.AmericanEnglish)

	sb.WriteString(fmt.Sprintf("#include %q\n", path.Base(g.conf.HeaderOutPath)))
	g.writeBindingsFrontmatter(sb)
	g.writeFileSourceHeader(sb, *g.Path)

}
