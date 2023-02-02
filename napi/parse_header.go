package napi

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/cpp"
)

type CPPArgDefault struct {
	NameSpace *string
	Val       *string
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
	IsPointer     bool
	Ident         *string
	DefaultValue  *CPPArgDefault
}

type CPPMethod struct {
	Ident            *string
	Overloads        []*[]*CPPArg
	Returns          *string
	ReturnsPrimitive bool
	ReturnsPointer   bool
	ExpectedArgs     int
	OptionalArgs     int
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
	Ident            *string
	Args             *[]*CPPArg
	Returns          *string
	ReturnsPrimitive bool
	ReturnsPointer   bool
}

type Enum struct {
	Ident *string
	Value int
}

type ParsedEnum struct {
	Ident     *string
	NameSpace *string
	Values    []*Enum
}

func getRootNode(path string) (*sitter.Node, []byte) {
	input, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	parser := sitter.NewParser()
	parser.SetLanguage(cpp.GetLanguage())

	tree, err := parser.ParseCtx(context.Background(), nil, input)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	n := tree.RootNode()
	return n, input
}

func (g *PackageGenerator) parseEnums(n *sitter.Node, input []byte, parseIncludes bool) []*ParsedEnum {
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
			fmt.Printf("Parsing enums in: %q\n", usedPath)
			rootNode, byteData := getRootNode(usedPath)
			if rootNode != nil {
				tmp_enums := g.parseEnums(rootNode, byteData, false)
				enums = append(enums, tmp_enums...)
			}
		}
	}
	return enums
}

func parseEnum(n *sitter.Node, input []byte) *ParsedEnum {
	enum_val := &ParsedEnum{}
	// parse enum namespace
	var namespace_node *sitter.Node
	p := n
	for namespace_node == nil {
		if p.Type() == "namespace_definition" {
			namespace_node = p
			break
		}
		p = p.Parent()
	}
	if namespace_node != nil {
		nameNode := namespace_node.ChildByFieldName("name")
		if nameNode != nil {
			name := nameNode.Content(input)
			enum_val.NameSpace = &name
		}
	}
	// parse enum name
	nameNode := n.ChildByFieldName("name")
	if nameNode != nil {
		name := nameNode.Content(input)
		enum_val.Ident = &name
	}

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
			val_name_node := child.ChildByFieldName("name")
			if val_name_node != nil {
				name := val_name_node.Content(input)
				parsedEnum.Ident = &name
			}
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

func parseLocalIncludes(n *sitter.Node, input []byte) []*string {
	includes := []*string{}
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
			if strings.Contains(content, "\"") && strings.Contains(content, ".h") {
				content = strings.ReplaceAll(content, "\"", "")
				content = strings.ReplaceAll(content, "#include", "")
				content = strings.TrimSpace(content)
				includes = append(includes, &content)
			}
		}
	}
	return includes
}

func (g *PackageGenerator) parseMethods(n *sitter.Node, input []byte) map[string]*CPPMethod {
	methods := map[string]*CPPMethod{}
	q, err := sitter.NewQuery([]byte("(declaration [type: (type_identifier) @type type: (primitive_type) @primitive type: (sized_type_specifier) @sized type: (qualified_identifier) @qual_type] [declarator: (function_declarator) @func declarator: (pointer_declarator) @ptr_func])"), cpp.GetLanguage())
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
				Ident:            parsed.Ident,
				Overloads:        []*[]*CPPArg{parsed.Args},
				Returns:          parsed.Returns,
				ReturnsPrimitive: parsed.ReturnsPrimitive,
				ReturnsPointer:   parsed.ReturnsPointer,
			}
			if !g.conf.IsMethodIgnored(*parsed.Ident) {
				methods[*parsed.Ident] = new_method
			}
		}
	}
	return methods
}

func parseNamespace(n *sitter.Node, input []byte) string {
	var out string
	q, err := sitter.NewQuery([]byte("(namespace_definition) @namespace"), cpp.GetLanguage())
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	qc := sitter.NewQueryCursor()
	qc.Exec(q, n)
	count := 0
	for count < 1 {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}
		ns := m.Captures[0].Node.ChildByFieldName("name")
		if ns != nil {
			out = ns.Content(input)
			count++
		}
	}
	return out
}

func parseCPPMethod(r *sitter.Node, b *sitter.Node, content []byte) *ParsedMethod {
	funcDeclNode := b
	returnsPointer := false
	if b.Type() == "pointer_declarator" {
		returnsPointer = true
		funcDeclNode = b.ChildByFieldName("declarator")
	}
	args := parseCPPArg(content, funcDeclNode.ChildByFieldName("parameters"))
	name := funcDeclNode.ChildByFieldName("declarator").Content(content)
	parsed := &ParsedMethod{
		Args:           args,
		ReturnsPointer: returnsPointer,
		Ident:          &name,
	}
	if r != nil {
		nodeType := r.Type()
		if nodeType == "primitive_type" || nodeType == "sized_type_specifier" {
			parsed.ReturnsPrimitive = true
		}
		tempReturns := r.Content(content)
		parsed.Returns = &tempReturns
		// mark `primitive` if it's a string (simplifies stuff a bit)
		if nodeType == "qualified_identifier" && (tempReturns == "std::string" || tempReturns == "string") {
			parsed.ReturnsPrimitive = true
		}
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
		tmp_type := type_node.Type()
		if tmp_type == "primitive_type" || tmp_type == "sized_type_specifier" {
			isPrimitive = true
		} else if tmp_type == "qualified_identifier" {
			// flag `std::string` as primitive
			if argType == "std::string" || argType == "string" {
				isPrimitive = true
			}
		}
		parsed_arg := &CPPArg{
			Type:          &argType,
			TypeQualifier: typeQualifier,
			IsPrimitive:   isPrimitive,
		}
		if node_type == "optional_parameter_declaration" {
			defaultNode := scoped_arg.ChildByFieldName("default_value")
			if defaultNode != nil {
				parsed_arg.DefaultValue = &CPPArgDefault{}
				ns_node := defaultNode.ChildByFieldName("scope")
				if ns_node != nil {
					ns := ns_node.Content(content)
					parsed_arg.DefaultValue.NameSpace = &ns
				}
				if val_node := defaultNode.ChildByFieldName("name"); val_node != nil {
					val := val_node.Content(content)
					parsed_arg.DefaultValue.Val = &val
				} else if val_node := defaultNode.ChildByFieldName("value"); val_node != nil {
					val := val_node.Content(content)
					parsed_arg.DefaultValue.Val = &val
				}
			}
		}
		refNode := scoped_arg.ChildByFieldName("declarator")
		// switch case to handle per node type
		switch refNode.Type() {

		case "pointer_declarator":
			{
				identNode := refNode.ChildByFieldName("declarator")
				if identNode != nil {
					identStr := identNode.Content(content)
					parsed_arg.Ident = &identStr
					parsed_arg.IsPointer = true
				}
			}
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
							// TODO: handle
							// fmt.Println(temp_child.Content(input))
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
							// TODO: handle
							fmt.Println("TODO: handle:", temp_child.Content(input))
						}
						*classes[class_name].FieldDecl = append(*classes[class_name].FieldDecl, parsed)
					}
				}
			}
			// fmt.Println("Matched: ", matched)
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
		params := declarator.ChildByFieldName("parameters")
		if params != nil {
			args := parseCPPArg(input, declarator.ChildByFieldName("parameters"))
			field_decl.Args = args
		}
		child_decl := declarator.ChildByFieldName("declarator")
		if child_decl != nil {
			identStr := child_decl.Content(input)
			field_decl.Ident = &identStr
		} else {
			func_decl := findChildNodeByType(declarator, "function_declarator")
			if func_decl != nil {
				child_decl := func_decl.ChildByFieldName("declarator")
				if child_decl != nil {
					identStr := child_decl.Content(input)
					field_decl.Ident = &identStr
				}

			}
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

type PreprocessBlock struct {
	Node    *sitter.Node
	RawArgs *string
	RawRes  *string
	Expr    bool
}

func (g *PackageGenerator) parseLiterals(n *sitter.Node, input []byte) map[string]*CPPMethod {
	preprocess_funcs := map[string]*PreprocessBlock{}
	linked_expr := map[string]*[]*string{}

	processed_expr := map[string]*CPPMethod{}
	q, err := sitter.NewQuery([]byte("(preproc_function_def) @function_literal"), cpp.GetLanguage())
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
			node := c.Node
			name_node := node.ChildByFieldName("name")
			name := name_node.Content(input)
			preprocess_funcs[name] = &PreprocessBlock{Node: node, Expr: false}
			value_node := node.ChildByFieldName("value")
			if value_node != nil {
				value := value_node.Content(input)
				val_split := strings.Split(value, "\n")
				for _, val := range val_split {
					if strings.Contains(val, "FUNC") && !strings.Contains(val, "OP") && !strings.Contains(val, "operator") {
						end_idx := strings.Index(val, ";")
						used_portion := strings.TrimSpace(val[:end_idx])
						used_split := strings.Split(used_portion, "FUNC")
						res_type := strings.TrimSpace(used_split[0])
						raw_args := strings.TrimSpace(used_split[1][1 : len(used_split[1])-1])
						preprocess_funcs[name].RawRes = &res_type
						preprocess_funcs[name].RawArgs = &raw_args
						break
					}
				}
			}
		}
	}

	q, err = sitter.NewQuery([]byte("(expression_statement) @expr"), cpp.GetLanguage())
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	qc = sitter.NewQueryCursor()
	qc.Exec(q, n)

	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}

		for _, c := range m.Captures {
			node := c.Node
			preprocess_name := findChildNodeByType(node, "call_expression")
			preprocess_alt := findChildNodeByType(node, "binary_expression")
			if preprocess_name != nil {
				scope_node := preprocess_name.ChildByFieldName("function")
				if scope_node != nil {
					scope_name := strings.Replace(scope_node.Content(input), "(", "", -1)
					if _, ok := preprocess_funcs[scope_name]; ok {
						preprocess_funcs[scope_name].Expr = true
						if linked_expr[scope_name] == nil {
							linked_expr[scope_name] = &[]*string{}
						}

						args_node := preprocess_name.ChildByFieldName("arguments")
						unary_expr := findChildNodeByType(args_node, "unary_expression")
						ptr_expr := findChildNodeByType(args_node, "pointer_expression")

						if unary_expr != nil {
							name := strings.TrimSpace(strings.Split(unary_expr.Content(input), ",")[1])
							*linked_expr[scope_name] = append(*linked_expr[scope_name], &name)
						} else if ptr_expr != nil {
							name := strings.TrimSpace(strings.Split(ptr_expr.Content(input), ",")[1])
							*linked_expr[scope_name] = append(*linked_expr[scope_name], &name)
						} else {
							name := strings.Replace(strings.TrimSpace(strings.Split(args_node.Content(input), ",")[1]), ")", "", -1)
							*linked_expr[scope_name] = append(*linked_expr[scope_name], &name)
						}
					}
				}
			} else if preprocess_alt != nil {
				scope_node := preprocess_alt.ChildByFieldName("left")
				if scope_node != nil {
					scope_name := strings.Replace(scope_node.Content(input), "(", "", -1)
					if _, ok := preprocess_funcs[scope_name]; ok {
						preprocess_funcs[scope_name].Expr = true
						if linked_expr[scope_name] == nil {
							linked_expr[scope_name] = &[]*string{}
						}
						name_node := preprocess_alt.ChildByFieldName("right")
						if name_node != nil {
							name := strings.TrimSpace(name_node.Content(input))
							*linked_expr[scope_name] = append(*linked_expr[scope_name], &name)
						}
					}
				}
			}
		}
	}

	// now parse/link data for generating output for preprocess def/expr pairs
	for name, block := range preprocess_funcs {
		if !block.Expr {
			continue
		}
		sb := strings.Builder{}
		for _, fn := range *linked_expr[name] {
			sb.WriteString(fmt.Sprintf("%s %s(%s);\n", *block.RawRes, *fn, *block.RawArgs))
		}
		parser := sitter.NewParser()
		parser.SetLanguage(cpp.GetLanguage())
		temp_input := []byte(sb.String())
		tree, err := parser.ParseCtx(context.Background(), nil, temp_input)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		n := tree.RootNode()
		processed_methods := g.parseMethods(n, temp_input)
		for fn, method := range processed_methods {
			processed_expr[fn] = method
		}
	}
	return processed_expr
}
