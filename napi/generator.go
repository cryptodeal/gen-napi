package napi

import (
	"context"
	"fmt"
	"go/ast"
	"os"
	"os/exec"
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
	conf      *PackageConfig
	NameSpace *string
	Name      *string
	Path      *string
	RootNode  *sitter.Node
	Input     *[]byte
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
		namespace := parseNamespace(n, input)

		split_path := strings.Split(path, "/")
		name := strings.Replace(split_path[len(split_path)-1], ".h", "", 1)

		napiGen := &PackageGenerator{
			conf:      napiConfig,
			NameSpace: &namespace,
			Name:      &name,
			RootNode:  n,
			Path:      &path,
			Input:     &input,
		}
		g.packageGenerators[*napiGen.Path] = napiGen
		bindings, header, env_wrapper, err := napiGen.Generate()
		if err != nil {
			return err
		}

		cmd_str := []string{"-i"}
		outPath := napiGen.conf.ResolvedBindingsOutPath(filepath.Dir(napiConfig.Path))
		cmd_str = append(cmd_str, outPath)
		err = os.MkdirAll(filepath.Dir(outPath), os.ModePerm)
		if err != nil {
			return nil
		}

		err = os.WriteFile(outPath, []byte(bindings), os.ModePerm)
		if err != nil {
			return nil
		}

		outPath = napiGen.conf.ResolvedHeaderOutPath(filepath.Dir(napiConfig.Path))
		cmd_str = append(cmd_str, outPath)
		err = os.MkdirAll(filepath.Dir(outPath), os.ModePerm)
		if err != nil {
			return nil
		}

		err = os.WriteFile(outPath, []byte(header), os.ModePerm)
		if err != nil {
			return nil
		}

		// programatically exec clang-format
		cmd := exec.Command("clang-format", cmd_str...)
		err = cmd.Run()
		if err != nil {
			return nil
		}

		outPath = napiGen.conf.ResolvedWrapperOutPath(filepath.Dir(napiConfig.Path))
		err = os.MkdirAll(filepath.Dir(outPath), os.ModePerm)
		if err != nil {
			return nil
		}

		err = os.WriteFile(outPath, []byte(env_wrapper), os.ModePerm)
		if err != nil {
			return nil
		}
	}
	return nil
}
