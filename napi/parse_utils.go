package napi

import (
	"fmt"
	"path/filepath"
	"strconv"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/cpp"
)

// utils
func splitMatches(matched []sitter.QueryCapture) (sitter.QueryCapture, sitter.QueryCapture) {
	return matched[0], matched[1]
}

func getTypeQualifier(n *sitter.Node, input []byte) *string {
	qualNode := findChildNodeByType(n, "type_qualifier")
	if qualNode != nil {
		type_qualifier := qualNode.Content(input)
		return &type_qualifier
	}
	return nil
}

func findChildNodeByType(n *sitter.Node, node_type string) *sitter.Node {
	child_count := int(n.ChildCount())
	for i := 0; i < child_count; i++ {
		tmp := n.Child(i)
		if tmp.Type() == node_type {
			return tmp
		}
	}
	return nil
}

func ParseMappedType(mapped_type string, all_type_mappings map[string]TypeMap, parser *sitter.Parser, og_type string) *CPPType {
	var parsed_type *CPPType = new(CPPType)
	input := []byte(fmt.Sprintf("using FakeType = %s;", mapped_type))
	root, _ := ParseByteToTree(input, parser)

	alias_node := QueryAliasDecls(root, input)[0]
	type_descriptor_node := alias_node.ChildByFieldName("type")
	if type_descriptor_node != nil {
		type_node := type_descriptor_node.ChildByFieldName("type")
		if type_node != nil {
			parsed_type = ParseType(type_node, input, all_type_mappings, parser)
			parsed_type.OGType = &og_type
		}
	}
	return parsed_type
}

// parse full namespace for declaration node
func parseNameSpace(n *sitter.Node, input []byte) *string {
	ns := []string{}
	p := n
	for p != nil && !p.IsMissing() {
		if p.Type() == "namespace_definition" {
			nameNode := p.ChildByFieldName("name")
			if nameNode != nil {
				ns = append(ns, nameNode.Content(input))
			}
		}
		p = p.Parent()
	}

	var nameSpace *string = new(string)
	for i := len(ns) - 1; i >= 0; i-- {
		*nameSpace += ns[i]
		if i > 0 {
			*nameSpace += "::"
		}
	}
	return nameSpace
}

func ParseType(r *sitter.Node, content []byte, mapped_types map[string]TypeMap, parser *sitter.Parser) *CPPType {
	out := &CPPType{}
	nodeType := r.Type()
	var type_name string
	var type_namespace *string

	if nodeType == "qualified_identifer" {
		scope := r.ChildByFieldName("scope")
		if scope != nil {
			ns := scope.Content(content)
			type_namespace = &ns
		}
		type_name = r.ChildByFieldName("name").Content(content)
	} else {
		type_name = stripNameSpace(r.Content(content))
		if nodeType == "type_identifier" {
			type_namespace = parseNameSpace(r, content)
		}
	}
	out.Name = type_name
	out.NameSpace = type_namespace
	out.UseMappedType(mapped_types, parser)

	isString := type_name == "std::string" || type_name == "string"
	// mark `primitive` if it's a string (simplifies stuff a bit)
	out.IsPrimitive = nodeType == "primitive_type" || nodeType == "sized_type_specifier" || (nodeType == "qualified_identifier" && isString)
	template_type := ParseTemplateArg(r, content)
	out.Template = template_type
	return out
}

func (g *PackageGenerator) parseReturnType(r *sitter.Node, content []byte) *CPPType {
	return ParseType(r, content, g.conf.TypeMappings, g.Parser)
}

// parse template types
type TemplateType struct {
	Name        *string
	NameSpace   *string
	Args        []*TemplateType
	IsPointer   bool
	IsPrimitive bool
}

func (t *TemplateType) GetFullType() string {
	var full_type string
	if t.NameSpace != nil {
		full_type = fmt.Sprintf("%s::", *t.NameSpace)
	}
	return full_type + *t.Name
}

func ParseTemplateType(n *sitter.Node, input []byte, ptr_declarator ...*sitter.Node) *TemplateType {
	var is_ptr_node *sitter.Node
	if len(ptr_declarator) > 0 {
		is_ptr_node = ptr_declarator[0]
	}
	template_type := &TemplateType{}
	switch n.Type() {
	case "qualified_identifier":
		{
			scope_node := n.ChildByFieldName("scope")
			if scope_node != nil && scope_node.Type() == "namespace_identifier" {
				name := scope_node.Content(input)
				template_type.NameSpace = &name
			}
			template_node := n.ChildByFieldName("name")
			name_node := template_node.ChildByFieldName("name")
			if name_node != nil {
				name := name_node.Content(input)
				template_type.Name = &name
			} else {
				name := template_node.Content(input)
				template_type.Name = &name
			}
			template_arg_node := template_node.ChildByFieldName("arguments")
			if template_arg_node != nil {
				arg_count := int(template_arg_node.ChildCount())
				for i := 0; i < arg_count; i++ {
					arg_node := template_arg_node.Child(i)
					tempArg := ParseTemplateArg(arg_node, input)
					if tempArg != nil {
						template_type.Args = append(template_type.Args, tempArg)
					}
				}
			}
		}
	default:
		{
			name := stripNameSpace(n.Content(input))
			if is_ptr_node != nil && is_ptr_node.Type() == "abstract_pointer_declarator" {
				template_type.IsPointer = true
			}
			if n.Type() == "primitive_type" {
				template_type.IsPrimitive = true
			}
			template_type.NameSpace = parseNameSpace(n, input)
			template_type.Name = &name
		}
	}
	if *template_type.Name == "string" {
		template_type.IsPrimitive = true
	}

	return template_type
}

func ParseTemplateArg(n *sitter.Node, input []byte) *TemplateType {
	switch n.Type() {
	case "parameter_declaration", "optional_parameter_declaration":
		{
			type_node := n.ChildByFieldName("type")
			if type_node != nil && type_node.Type() == "qualified_identifier" {
				return ParseTemplateArg(type_node, input)
			}
		}
	case "qualified_identifier":
		{
			template_node := n.ChildByFieldName("name")
			if template_node != nil && template_node.Type() == "template_type" {
				return ParseTemplateType(n, input)
			}
		}
	case "type_descriptor":
		{
			type_node := n.ChildByFieldName("type")
			declarator := n.ChildByFieldName("declarator")
			if type_node != nil {
				return ParseTemplateType(type_node, input, declarator)
			}
		}
	case "number_literal":
		{
			return ParseTemplateType(n, input)
		}
	}
	return nil
}

// parse enums
type Enum struct {
	Name  *string
	Value int
}

type ParsedEnum struct {
	Name      string
	NameSpace *string
	Values    []*Enum
}

func (g *PackageGenerator) parseEnumDecls(n *sitter.Node, input []byte, parseIncludes bool) []*ParsedEnum {
	enums := []*ParsedEnum{}
	q, err := sitter.NewQuery([]byte("(enum_specifier) @enums"), cpp.GetLanguage())
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
			enums = append(enums, parseEnum(c.Node, input))
		}
	}
	if parseIncludes {
		for _, local := range g.LocalIncludes {
			usedPath := filepath.Join(g.conf.LibRootDir, *local)
			rootNode, byteData := getRootNode(usedPath, g.Parser)
			if rootNode != nil {
				tmp_enums := g.parseEnumDecls(rootNode, byteData, false)
				enums = append(enums, tmp_enums...)
			}
		}
	}
	return enums
}

func parseEnum(n *sitter.Node, input []byte) *ParsedEnum {
	enum_val := &ParsedEnum{}
	// parse enum namespace
	enum_val.NameSpace = parseNameSpace(n, input)
	// parse enum name
	enum_val.Name = n.ChildByFieldName("name").Content(input)

	// parse enum values
	bodyNode := n.ChildByFieldName("body")
	if bodyNode != nil && bodyNode.Type() == "enumerator_list" {
		child_count := int(bodyNode.ChildCount())
		enum_children := []*sitter.Node{}
		for i := 0; i < child_count; i++ {
			tmp_child := bodyNode.Child(i)
			if tmp_child.Type() != "enumerator" {
				continue
			}
			enum_children = append(enum_children, tmp_child)
		}
		for i, child := range enum_children {
			parsedEnum := &Enum{}
			name := child.ChildByFieldName("name").Content(input)
			parsedEnum.Name = &name
			val_node := child.ChildByFieldName("value")
			if val_node != nil {
				v, err := strconv.Atoi(val_node.Content(input))
				if err == nil {
					parsedEnum.Value = v
				}
			} else {
				parsedEnum.Value = i
			}
			enum_val.Values = append(enum_val.Values, parsedEnum)
		}
	}
	return enum_val
}
