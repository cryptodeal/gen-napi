package napi

import (
	"fmt"
	"path/filepath"
	"strings"
)

func (g *PackageGenerator) writeBindings(sb *strings.Builder, classes map[string]*CPPClass, methods map[string]*CPPMethod) {
	// lower_caser := cases.Lower(language.AmericanEnglish)

	sb.WriteString(fmt.Sprintf("#include %q\n", filepath.Base(g.conf.Path)))
	g.writeBindingsFrontmatter(sb)
	g.writeFileSourceHeader(sb, *g.Path)

}
