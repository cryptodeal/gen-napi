package napi

import (
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/cpp"
)

func parseIncludes(n *sitter.Node, input []byte) string {
	includes := &strings.Builder{}
	q, err := sitter.NewQuery([]byte("(preproc_include) @includes"), cpp.GetLanguage())
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	qc := sitter.NewQueryCursor()
	qc.Exec(q, n)

	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}
		for _, c := range m.Captures {
			content := c.Node.Content(input)
			if strings.Contains(content, "#include <napi.h>") {
				continue
			}
			includes.WriteString(c.Node.Content(input))
			includes.WriteByte('\n')
		}
	}
	return includes.String()
}

type NameSpaceGroup struct {
	NameSpace   *string
	IsClass     *string
	Methods     []*string
	IsGlobalVar []*string
}

func parseNameSpaces(n *sitter.Node, input []byte) []*string {
	includes := []*string{}
	q, err := sitter.NewQuery([]byte("(namespace_definition) @namespace"), cpp.GetLanguage())
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	qc := sitter.NewQueryCursor()
	qc.Exec(q, n)

	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}
		for _, c := range m.Captures {
			content := c.Node.Content(input)
			includes = append(includes, &content)
		}
	}
	return includes
}
