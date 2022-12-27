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

// Responsible for generating the code for an input package
type PackageGenerator struct {
	conf     *PackageConfig
	Name     *string
	Path     *string
	RootNode *sitter.Node
	Input    *[]byte
}

type EnumField struct {
	Name string
	Val  string
}

func New(config *Config) *TSGo {
	return &TSGo{
		conf:              config,
		packageGenerators: make(map[string]*PackageGenerator),
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
			Input:    &input,
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
