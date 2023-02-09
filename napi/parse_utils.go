package napi

import (
	sitter "github.com/smacker/go-tree-sitter"
)

type TemplateType struct {
	Name      *string
	NameSpace *string
	Args      []*TemplateType
}

func ParseTemplateType(n *sitter.Node, input []byte) *TemplateType {
	template_type := &TemplateType{}
	switch n.Type() {
	case "qualified_identifier":
		{
			scope_node := n.ChildByFieldName("scope")
			if scope_node != nil && scope_node.Type() == "namespace_identifier" {
				name := scope_node.Content(input)
				template_type.NameSpace = &name
			}
			used_node := n.ChildByFieldName("name")
			name_node := used_node.ChildByFieldName("name")
			if name_node != nil {
				name := name_node.Content(input)
				template_type.Name = &name
			}
			template_arg_node := used_node.ChildByFieldName("arguments")
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
			name := n.Content(input)
			template_type.Name = &name
		}
	}

	return template_type
}

func ParseTemplateArg(n *sitter.Node, input []byte) *TemplateType {
	var template_type *TemplateType = nil
	switch n.Type() {
	case "parameter_declaration", "optional_parameter_declaration":
		{
			type_node := n.ChildByFieldName("type")
			if type_node != nil && type_node.Type() == "qualified_identifier" {
				template_node := type_node.ChildByFieldName("name")
				if template_node != nil && template_node.Type() == "template_type" {
					template_type = ParseTemplateType(type_node, input)
				}
			}
		}
	case "type_descriptor":
		{
			type_node := n.ChildByFieldName("type")
			if type_node != nil {
				template_type = ParseTemplateType(type_node, input)
			}
		}
	}
	return template_type
}
