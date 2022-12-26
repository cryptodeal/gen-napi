package napi

import (
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/cpp"
)

type CPPArgDefault struct {
	Val *string
}

type TemplateArg struct {
	Identifier *string
}

type TemplateDeclArg struct {
	Identifier *string
	MetaType   *string
}

type TemplateDecl struct {
	Args *[]*TemplateDeclArg
}

type QualifiedIdentifier struct {
	Scope        *string
	Name         *string
	TemplateArgs *[]*TemplateArg
}

type CPPFriendFunc struct {
	QualifiedIdent *QualifiedIdentifier
	Args           *[]*CPPArg
}

type CPPFieldDecl struct {
	Ident         *string
	Args          *[]*CPPArg
	Returns       *CPPType
	TypeQualifier *string
}

type CPPFriend struct {
	Ident          *string
	QualifiedIdent *QualifiedIdentifier
	IsClass        bool
	Type           *CPPType
	FuncDecl       *CPPFriendFunc
}

type CPPType struct {
	FullType     *string
	Scope        *string
	Name         *string
	NameSpace    *string
	TemplateType *[]*TemplateArg
}

type CPPClass struct {
	NameSpace    *string
	FieldDecl    *[]*CPPFieldDecl
	FriendDecl   *[]*CPPFriend
	Decl         *[]*ParsedClassDecl
	TemplateDecl *[]*TemplateMethod
}

type CPPArg struct {
	TypeQualifier *string
	IsPrimitive   bool
	Type          *string
	RefDecl       *string
	Ident         *string
	DefaultValue  *CPPArgDefault
}

type CPPMethod struct {
	Ident     *string
	Overloads []*[]*CPPArg
	Returns   *string
}

type TemplateMethod struct {
	TemplateDecl          *TemplateDecl
	Returns               *string
	PointerMethod         bool
	StorageClassSpecifier *string
	RefDecl               *string
	Ident                 *string
	Args                  *[]*CPPArg
	TypeQualifier         *string
}

type ParsedClassDecl struct {
	Ident        *string
	Args         *[]*CPPArg
	Returns      *string
	Explicit     bool
	Virtual      bool
	IsDestructor bool
}

type ParsedMethod struct {
	Ident   *string
	Args    *[]*CPPArg
	Returns *string
}

func parseMethods(n *sitter.Node, input []byte) map[string]*CPPMethod {
	methods := map[string]*CPPMethod{}
	q, err := sitter.NewQuery([]byte("(declaration type: (type_identifier) @type declarator: (function_declarator) @func)"), cpp.GetLanguage())
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
		res, body := splitMatches(m.Captures)
		parsed := parseCPPMethod(res.Node, body.Node, input)
		if v, ok := methods[*parsed.Ident]; ok {
			// encountered method previously (fn overloading)
			v.Overloads = append(v.Overloads, parsed.Args)
		} else {
			// first time having encountered this method, so create a new entry
			new_method := &CPPMethod{
				Ident:     parsed.Ident,
				Overloads: []*[]*CPPArg{parsed.Args},
				Returns:   parsed.Returns,
			}
			methods[*parsed.Ident] = new_method
		}
	}
	return methods
}

func parseCPPMethod(r *sitter.Node, b *sitter.Node, content []byte) *ParsedMethod {
	args := parseCPPArg(content, b.ChildByFieldName("parameters"))
	name := b.ChildByFieldName("declarator").Content(content)
	parsed := &ParsedMethod{
		Args:  args,
		Ident: &name,
	}
	if r != nil {
		tempReturns := r.Content(content)
		parsed.Returns = &tempReturns
	}
	return parsed
}

func parseCPPArg(content []byte, arg_list *sitter.Node) *[]*CPPArg {
	args := []*CPPArg{}
	if arg_list == nil {
		return &args
	}
	child_count := int(arg_list.ChildCount())
	for i := 0; i < child_count; i++ {
		scoped_arg := arg_list.Child(i)
		node_type := scoped_arg.Type()
		if node_type != "parameter_declaration" && node_type != "optional_parameter_declaration" {
			continue
		}
		type_node := scoped_arg.ChildByFieldName("type")
		argType := type_node.Content(content)
		typeQualifier := getTypeQualifier(scoped_arg, content)
		isPrimitive := false
		if type_node.Type() == "primitive_type" || type_node.Type() == "sized_type_specifier" {
			isPrimitive = true
		}
		parsed_arg := &CPPArg{
			Type:          &argType,
			TypeQualifier: typeQualifier,
			IsPrimitive:   isPrimitive,
		}
		refNode := scoped_arg.ChildByFieldName("declarator")
		// switch case to handle per node type
		switch refNode.Type() {
		case "reference_declarator":
			{
				identNode := findChildNodeByType(refNode, "identifier")
				if identNode != nil {
					identStr := identNode.Content(content)
					parsed_arg.Ident = &identStr
					refDeclStr := strings.ReplaceAll(refNode.Content(content), identStr, "")
					parsed_arg.RefDecl = &refDeclStr
				}
			}
		case "identifier":
			{
				identStr := refNode.Content(content)
				parsed_arg.Ident = &identStr
			}
		case "optional_parameter_declaration":
			{
				tempDefault := scoped_arg.ChildByFieldName("default_value").Content(content)
				parsed_arg.DefaultValue = &CPPArgDefault{Val: &tempDefault}
			}
		}
		args = append(args, parsed_arg)
	}
	return &args
}

func parseClasses(n *sitter.Node, input []byte) map[string]*CPPClass {
	classes := map[string]*CPPClass{}

	q, err := sitter.NewQuery([]byte("(class_specifier) @class_def"), cpp.GetLanguage())
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
		// fmt.Println(len(m.Captures))
		for _, c := range m.Captures {
			namespace := getNameSpace(c.Node, input)
			class_name := c.Node.ChildByFieldName("name").Content(input)
			classes[class_name] = &CPPClass{
				NameSpace: namespace,
			}
			class_body := c.Node.ChildByFieldName("body")
			class_friends := &[]*CPPFriend{}
			if class_body == nil {
				// TODO: should probably parse class def w/o body as well
				continue
			}

			child_count := int(class_body.ChildCount())
			matched := 0
			for i := 0; i < child_count; i++ {
				temp_child := class_body.Child(i)
				// switch case to handle per node type
				switch temp_child.Type() {
				case "friend_declaration": // WORKING?? (needs unit tests)
					{
						new_friend := parseClassFriend(temp_child, input)
						*class_friends = append(*class_friends, new_friend)
					}
				case "template_declaration": // WORKING?? (needs unit tests)
					{
						if classes[class_name].TemplateDecl == nil {
							classes[class_name].TemplateDecl = &[]*TemplateMethod{}
						}
						temp_decl := parseClassTemplateMethod(temp_child, input)
						if temp_decl.Ident == nil {
							fmt.Println(temp_child.Content(input))
						}
						*classes[class_name].TemplateDecl = append(*classes[class_name].TemplateDecl, temp_decl)
					}
				case "declaration": // WORKING?? (needs unit tests)
					{
						if classes[class_name].Decl == nil {
							classes[class_name].Decl = &[]*ParsedClassDecl{}
						}
						parsed := parseClassDecl(temp_child, input)
						*classes[class_name].Decl = append(*classes[class_name].Decl, parsed)
					}
				case "field_declaration":
					{
						// TODO: parse top level nodes w type `field_declaration`
						if classes[class_name].FieldDecl == nil {
							classes[class_name].FieldDecl = &[]*CPPFieldDecl{}
						}
						matched++
						parsed := parseFieldDecl(temp_child, input)
						if parsed.Ident == nil {
							fmt.Println(temp_child.Content(input))
							fmt.Println(temp_child.Content(input))
						}
						*classes[class_name].FieldDecl = append(*classes[class_name].FieldDecl, parsed)
					}
				}
			}
			fmt.Println("Matched: ", matched)
			classes[class_name].FriendDecl = class_friends
		}
	}
	return classes
}

func parseClassDecl(n *sitter.Node, input []byte) *ParsedClassDecl {
	explicitNode := n.ChildByFieldName("explicit_function_specifier")
	parsed := &ParsedClassDecl{}
	if explicitNode != nil {
		parsed.Explicit = true
	}
	decl := n.ChildByFieldName("declarator")
	if decl != nil {
		parsed.Args = parseCPPArg(input, decl.ChildByFieldName("parameters"))
		nameNode := decl.ChildByFieldName("declarator")
		if nameNode != nil {
			if nameNode.Type() == "destructor_name" {
				parsed.IsDestructor = true
				identNode := findChildNodeByType(nameNode, "identifier")
				identStr := identNode.Content(input)
				parsed.Ident = &identStr
			} else {
				nameStr := nameNode.Content(input)
				parsed.Ident = &nameStr
			}
		}
	}
	return parsed
}

func parseFieldDecl(n *sitter.Node, input []byte) *CPPFieldDecl {
	field_decl := CPPFieldDecl{}

	type_node := n.ChildByFieldName("type")
	if type_node != nil {
		parsed_type := CPPType{}
		full_type := type_node.Content(input)
		parsed_type.FullType = &full_type
		field_decl.Returns = &parsed_type
	}

	declarator := n.ChildByFieldName("declarator")
	if declarator != nil {
		child_decl := declarator.ChildByFieldName("declarator")
		if child_decl != nil {
			identStr := child_decl.Content(input)
			field_decl.Ident = &identStr
		}
	}
	return &field_decl
}

func parseClassTemplateMethod(n *sitter.Node, input []byte) *TemplateMethod {
	template_method := &TemplateMethod{
		TemplateDecl: &TemplateDecl{},
	}
	params := n.ChildByFieldName("parameters")
	template_method.TemplateDecl.Args = &[]*TemplateDeclArg{}
	param_count := int(params.ChildCount())
	for i := 0; i < param_count; i++ {
		param := params.Child(i)
		if param.Type() != "type_parameter_declaration" {
			continue
		}
		param_split := strings.Split(param.Content(input), " ")
		DeclArg := &TemplateDeclArg{
			Identifier: &param_split[1],
			MetaType:   &param_split[0],
		}
		*template_method.TemplateDecl.Args = append(*template_method.TemplateDecl.Args, DeclArg)
	}

	childCount := int(n.ChildCount())
	for i := 0; i < childCount; i++ {
		tempChild := n.Child(i)
		childType := tempChild.Type()
		// switch case to handle per node type
		switch childType {
		case "declaration":
			{
				template_method.Returns = getTypeVal(tempChild, input)
				declarator := tempChild.ChildByFieldName("declarator")
				// switch case to handle per node type
				switch declarator.Type() {
				case "pointer_declarator":
					{
						template_method.PointerMethod = true
						decl := declarator.ChildByFieldName("declarator")
						// switch case to handle per node type
						switch decl.Type() {
						case "function_declarator":
							{
								template_method.Args = parseCPPArg(input, decl.ChildByFieldName("parameters"))
								template_method.TypeQualifier = getTypeQualifier(decl, input)
								nameNode := decl.ChildByFieldName("name")
								if nameNode != nil {
									name := nameNode.Content(input)
									template_method.Ident = &name
								} else {
									decl := decl.ChildByFieldName("declarator")
									if decl != nil {
										parseTemplateFuncIdent(decl, input, template_method)
									}
								}
							}
						case "function_definition":
							{
								parseTemplateFuncDefNode(tempChild, input, template_method)
							}
						}
					}
				case "function_declarator":
					{
						decl := declarator.ChildByFieldName("declarator")
						template_method.Args = parseCPPArg(input, decl.ChildByFieldName("parameters"))
						parseTemplateFuncIdent(decl, input, template_method)
						template_method.TypeQualifier = getTypeQualifier(decl, input)
					}
				}
			}
		case "function_definition":
			{
				parseTemplateFuncDefNode(tempChild, input, template_method)
			}
		}
	}
	return template_method
}

func parseTemplateFuncIdent(n *sitter.Node, input []byte, method *TemplateMethod) {
	nameNode := n.ChildByFieldName("name")
	if nameNode != nil {
		name := nameNode.Content(input)
		method.Ident = &name
	} else {
		name := n.Content(input)
		method.Ident = &name
	}
}

func getStorageClassSpecifier(n *sitter.Node, input []byte) *string {
	storageNode := findChildNodeByType(n, "storage_class_specifier")
	if storageNode != nil {
		storage := storageNode.Content(input)
		return &storage
	}
	return nil
}

func getTypeVal(n *sitter.Node, input []byte) *string {
	var res *string
	tempType := n.ChildByFieldName("type")
	if tempType != nil {
		content := tempType.Content(input)
		res = &content
	}
	return res
}

func parseTemplateFuncDefNode(n *sitter.Node, input []byte, method *TemplateMethod) {
	method.StorageClassSpecifier = getStorageClassSpecifier(n, input)
	method.Returns = getTypeVal(n, input)
	declarator := n.ChildByFieldName("declarator")
	if declarator != nil {
		// switch case to handle per node type
		switch declarator.Type() {
		case "function_declarator":
			{
				nameNode := declarator.ChildByFieldName("declarator")
				if nameNode != nil {
					ident := nameNode.Content(input)
					method.Ident = &ident
				}
				method.Args = parseCPPArg(input, declarator.ChildByFieldName("parameters"))
			}
		case "reference_declarator":
			{
				funcDecl := findChildNodeByType(declarator, "function_declarator")
				method.Args = parseCPPArg(input, funcDecl.ChildByFieldName("parameters"))
				decl := funcDecl.ChildByFieldName("declarator")
				name := decl.Content(input)
				refDecl := strings.ReplaceAll(funcDecl.Content(input), name, "")
				method.RefDecl = &refDecl
				method.Ident = &name
			}
		}
	}

}

func parseClassFriend(n *sitter.Node, input []byte) *CPPFriend {
	new_friend := &CPPFriend{}
	child_count := int(n.ChildCount())
	for j := 0; j < child_count; j++ {
		grandchild := n.Child(j)
		tempType := grandchild.Type()
		switch tempType {
		case "type_identifier":
			tempName := grandchild.Content(input)
			new_friend.Ident = &tempName
			new_friend.IsClass = true
		case "declaration":
			great_grandchild_count := int(grandchild.ChildCount())
			for k := 0; k < great_grandchild_count; k++ {
				great_grandchild := grandchild.Child(k)
				temp_great_type := great_grandchild.Type()
				/* nested switch, is a bit ugly, but good perf */
				switch temp_great_type {
				case "qualified_identifier":
					qualID := &QualifiedIdentifier{}
					scope := great_grandchild.ChildByFieldName("scope")
					if scope != nil {
						tempScope := scope.Content(input)
						qualID.Scope = &tempScope
					}
					_name := great_grandchild.ChildByFieldName("name")
					if _name != nil {
						name := _name.ChildByFieldName("name")
						if name != nil {
							tempName := name.Content(input)
							qualID.Name = &tempName
						}
						arguments := _name.ChildByFieldName("arguments")
						arg_childs := int(arguments.ChildCount())
						template_args := &[]*TemplateArg{}
						for l := 0; l < arg_childs; l++ {
							arg := arguments.Child(l)
							argType := arg.Type()
							if argType == "type_descriptor" {
								temp_arg_type := arg.ChildByFieldName("type")
								if temp_arg_type != nil {
									parsed_temp_arg := temp_arg_type.Content(input)
									*template_args = append(*template_args, &TemplateArg{&parsed_temp_arg})
								}
							}
						}
						qualID.TemplateArgs = template_args
					}
					new_friend.QualifiedIdent = qualID
				case "function_declarator":
					decl := great_grandchild.ChildByFieldName("declarator")
					friend_func := &CPPFriendFunc{QualifiedIdent: &QualifiedIdentifier{}}
					if decl != nil {
						scope := decl.ChildByFieldName("scope")
						if scope != nil {
							tempScope := scope.Content(input)
							friend_func.QualifiedIdent.Scope = &tempScope
						}
						name := decl.ChildByFieldName("name")
						if name != nil {
							tempName := name.Content(input)
							friend_func.QualifiedIdent.Name = &tempName
						}
					}
					params := great_grandchild.ChildByFieldName("parameters")
					friend_func.Args = parseCPPArg(input, params)
					new_friend.FuncDecl = friend_func
				} /* end nested switch/case */
			} /* end outer switch/case */
		}
	}
	return new_friend
}

// helper functions

func getNameSpace(n *sitter.Node, input []byte) *string {
	var nameSpace *string
	test_node := n
	for test_node.Type() != "namespace_definition" {
		test_node = test_node.Parent()
	}
	name_node := test_node.ChildByFieldName("name")
	if name_node != nil {
		name := name_node.Content(input)
		nameSpace = &name
	}
	return nameSpace
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

func splitMatches(matched []sitter.QueryCapture) (sitter.QueryCapture, sitter.QueryCapture) {
	return matched[0], matched[1]
}
