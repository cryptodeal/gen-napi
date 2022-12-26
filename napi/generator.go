package napi

import (
	"context"
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/cpp"
)

// Generator for one or more input packages, responsible for linking
// them together if necessary.
type TSGo struct {
	conf              *Config
	packageGenerators map[string]*PackageGenerator
}

type ArgHelpers struct {
	FFIType     string
	CGoWrapType string
	OGGoType    string
	Name        string
	ASTField    *ast.Field
}

type ResHelpers struct {
	FFIType     string
	CGoWrapType string
	OGGoType    string
	ASTType     *ast.Expr
}

type FFIFunc struct {
	args           []*ArgHelpers
	returns        []*ResHelpers
	isHandleFn     bool
	isStarExpr     bool
	name           *string
	fieldAccessors []*StructAccessor
	disposeHandle  *DisposeStructFunc
}

type DisposeStructFunc struct {
	args   []*ArgHelpers
	fnName string
	name   string
}

type StructAccessor struct {
	args           []*ArgHelpers
	returns        []*ResHelpers
	isHandleFn     *string
	isStarExpr     bool
	isOptional     bool
	name           *string
	fnName         *string
	arrayType      *string
	structType     *string
	fieldAccessors []*StructAccessor
	disposeHandle  *DisposeStructFunc
}

type ClassWrapper struct {
	args    []*ArgHelpers
	returns []*ResHelpers
	// TODO: might be useful in future?
	// isHandleFn     *string
	// isStarExpr     bool
	// isOptional     bool
	structType     *string
	name           *string
	fieldAccessors []*StructAccessor
	disposeHandle  *DisposeStructFunc
}

type FFIState struct {
	GoImports        map[string]bool
	CImports         map[string]bool
	FFIHelpers       map[string]bool
	CHelpers         map[string]bool
	FFIFuncs         map[string]*FFIFunc
	StructHelpers    map[string][]*StructAccessor
	ParsedStructs    map[string]bool
	TypeHelpers      map[string]string
	GoWrappedStructs map[string]bool
}

// Responsible for generating the code for an input package
type PackageGenerator struct {
	conf     *PackageConfig
	Name     *string
	Path     *string
	RootNode *sitter.Node
}

type EnumField struct {
	Name string
	Val  string
}

type TSHelpers struct {
	EnumStructs map[string][]*EnumField
}

func New(config *Config) *TSGo {
	return &TSGo{
		conf:              config,
		packageGenerators: make(map[string]*PackageGenerator),
	}
}

func (g *TSGo) SetTypeMapping(goType string, tsType string) {
	for _, p := range g.conf.Packages {
		p.TypeMappings[goType] = tsType
	}
}

func (g *TSGo) Generate() error {
	hdr_paths := g.conf.PackageNames()

	for _, path := range hdr_paths {
		input, err := os.ReadFile(path)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		parser := sitter.NewParser()
		parser.SetLanguage(cpp.GetLanguage())

		tree, err := parser.ParseCtx(context.Background(), nil, input)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		n := tree.RootNode()

		napiConfig := g.conf.PackageConfig(path)

		split_path := strings.Split(path, "/")
		name := strings.Replace(split_path[len(split_path)-1], ".h", "", 1)

		napiGen := &PackageGenerator{
			conf:     napiConfig,
			Name:     &name,
			RootNode: n,
			Path:     &path,
		}
		g.packageGenerators[*napiGen.Path] = napiGen
		bindings, header, err := napiGen.Generate()
		if err != nil {
			return err
		}

		outPath := napiGen.conf.ResolvedBindingsOutPath(filepath.Dir(napiConfig.Path))
		err = os.MkdirAll(filepath.Dir(outPath), os.ModePerm)
		if err != nil {
			return nil
		}

		err = os.WriteFile(outPath, []byte(bindings), os.ModePerm)
		if err != nil {
			return nil
		}

		outPath = napiGen.conf.ResolvedHeaderOutPath(filepath.Dir(napiConfig.Path))
		err = os.MkdirAll(filepath.Dir(outPath), os.ModePerm)
		if err != nil {
			return nil
		}

		err = os.WriteFile(outPath, []byte(header), os.ModePerm)
		if err != nil {
			return nil
		}
	}
	return nil
}
