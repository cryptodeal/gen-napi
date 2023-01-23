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
		}
	}
	return includes.String()
}

func parseGlobalVars(n *sitter.Node, input []byte) string {
	global_vars := &strings.Builder{}
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
			name := c.Node.ChildByFieldName("name")
			if name != nil {
				nameContent := name.Content(input)
				if !strings.EqualFold(nameContent, "global_vars") {
					continue
				}
				bodyNode := c.Node.ChildByFieldName("body")
				if bodyNode != nil {
					splitBody := strings.Split(bodyNode.Content(input), "\n")
					length := len(splitBody)
					for i, line := range splitBody {
						if i == 0 || i == length-1 {
							continue
						}
						global_vars.WriteString(fmt.Sprintf("%s\n", line))
						if i == length-2 {
							global_vars.WriteByte('\n')
						}
					}
				}
			}
		}
	}
	return global_vars.String()
}

func getFuncName(n *sitter.Node, input []byte) string {
	var name string
	nameNode := n.ChildByFieldName("declarator")
	if nameNode != nil {
		name = nameNode.Content(input)
	}
	return name
}

func parsePrivateHelpers(n *sitter.Node, input []byte) map[string]string {
	helpers := map[string]string{}
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
			name := c.Node.ChildByFieldName("name")
			if name != nil {
				nameContent := name.Content(input)
				if !strings.EqualFold(nameContent, "private_helpers") {
					continue
				}
				bodyNode := c.Node.ChildByFieldName("body")
				if bodyNode != nil {
					child_count := int(bodyNode.ChildCount())
					i := 0
					for i < child_count {
						child := bodyNode.Child(i)
						if child != nil {
							childType := child.Type()
							if childType == "function_definition" {
								decl := child.ChildByFieldName("declarator")
								if decl != nil {
									name := getFuncName(decl, input)
									helpers[name] = child.Content(input)
								}
							} else if childType == "template_declaration" {
								funcNode := findChildNodeByType(child, "function_definition")
								if funcNode != nil {
									decl := funcNode.ChildByFieldName("declarator")
									if decl != nil {
										name := getFuncName(decl, input)
										helpers[name] = child.Content(input)
									}
								}
							}
						}
						i++
					}
				}
			}
		}
	}
	return helpers
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
