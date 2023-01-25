package napi

import (
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/cpp"
)

func isAlreadyIncluded(val string) bool {
	if strings.Contains(val, "napi.h") || strings.Contains(val, "atomic") || strings.Contains(val, "string") {
		return true
	}
	return false
}

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
			if isAlreadyIncluded(string(content)) {
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

func parseScopedForcedMethods(n *sitter.Node, input []byte) []FnOpts {
	methods := []FnOpts{}
	child_count := int(n.ChildCount())
	i := 0
	var unusedComment *sitter.Node
	for i < child_count {
		child := n.Child(i)
		if child == nil {
			break
		}
		fnData := FnOpts{}
		switch child.Type() {
		case "comment":
			{
				unusedComment = child
			}
		case "function_definition":
			{
				decl := child.ChildByFieldName("declarator")
				if decl != nil {
					fnData.Name = getFuncName(decl, input)
					fnData.FnBody = child.Content(input)
					if unusedComment != nil {
						comment := unusedComment.Content(input)
						comment_split := strings.Split(comment, "\n")
						length := len(comment_split)
						for j, line := range comment_split {
							if j == 0 {
								continue
							}
							if strings.Contains(line, "`ts_return_type`") {
								fnData.TSReturnType = strings.TrimSpace(strings.Replace(line, "@gen-napi-`ts_return_type`:", "", 1))
							} else if strings.Contains(line, "`ts_args`") {
								argString := strings.TrimSpace(strings.Replace(line, "@gen-napi-`ts_args`:", "", 1))
								argString = argString[1 : len(argString)-1]
								argSplit := strings.Split(argString, ",")
								for _, arg := range argSplit {
									fnArg := FnArg{}
									hasType := strings.Contains(arg, ":")
									hasDefault := strings.Contains(arg, "=")
									if !hasType && !hasDefault {
										fnArg.Name = strings.TrimSpace(arg)
										fnArg.TSType = "any"
									} else if hasType && !hasDefault {
										argSplit := strings.Split(arg, ":")
										fnArg.Name = strings.TrimSpace(argSplit[0])
										fnArg.TSType = strings.TrimSpace(argSplit[1])
									} else if !hasType && hasDefault {
										argSplit := strings.Split(arg, "=")
										fnArg.Name = strings.TrimSpace(argSplit[0])
										fnArg.TSType = "any"
										fnArg.Default = strings.TrimSpace(argSplit[1])
									} else {
										argSplit := strings.Split(arg, ":")
										fnArg.Name = strings.TrimSpace(argSplit[0])
										defaultSplit := strings.Split(argSplit[1], "=")
										fnArg.TSType = strings.TrimSpace(defaultSplit[0])
										fnArg.Default = strings.TrimSpace(defaultSplit[1])
									}
									fnData.Args = append(fnData.Args, fnArg)
								}
							}
							if j == length-2 {
								break
							}
						}
						// comment is used, reset `unusedComment`
						unusedComment = nil
					}
					methods = append(methods, fnData)
				}
			}
		}
		i++
	}
	return methods
}

func parseScopedFnBlock(n *sitter.Node, input []byte, scopeName string) *[]FnOpts {
	var methods []FnOpts
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
				if !strings.EqualFold(nameContent, scopeName) {
					continue
				}
				bodyNode := c.Node.ChildByFieldName("body")
				if bodyNode != nil {
					methods = parseScopedForcedMethods(bodyNode, input)
				}
			}
		}
	}
	return &methods
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
