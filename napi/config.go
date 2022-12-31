package napi

import (
	"log"
	"path/filepath"
	"strings"
)

const defaultOutBindingsFileName = "bindings.cc"
const defaultOutHeaderFileName = "bindings.h"

type TypeHandler struct {
	OutType string `yaml:"out_type"`
	OutVar  string `yaml:"out_var"`
	Handler string `yaml:"handler"`
}

type TypeMap struct {
	TSType   string `yaml:"ts"`
	NapiType string `yaml:"napi"`
	CastsTo  string `yaml:"casts_to"`
	CastNapi string `yaml:"cast_napi"`
}

type FnOpts struct {
	Name          string   `yaml:"name"`
	JSWrapperName string   `yaml:"js_wrapper_name"`
	JSWrapperAlts []string `yaml:"js_wrapper_alts"`
	FnBody        string   `yaml:"body"`
	IsVoid        bool     `yaml:"is_void"`
}

type ClassOpts struct {
	Fields        []string `yaml:"fields"`
	Methods       []string `yaml:"methods"`
	ForcedMethods []FnOpts `yaml:"forced_methods"`
	Constructor   string   `yaml:"constructor"`
}

type PackageConfig struct {
	// The package path just like you would import it in Go
	Path string `yaml:"path"`

	// Where this output should be written to.
	// If you specify a folder it will be written to a file `index.ts` within that folder. By default it is written into the Golang package folder.
	BindingsOutPath string `yaml:"bindings_out_path"`
	HeaderOutPath   string `yaml:"header_out_path"`

	// Customize the indentation (use \t if you want tabs)
	Indent string `yaml:"indent"`

	// Specify your own custom type translations, useful for custom types, `time.Time` and `null.String`.
	// Be default unrecognized types will be output as `any /* name */`.
	TypeMappings map[string]TypeMap `yaml:"type_mappings"`

	TypeHandlers map[string]TypeHandler `yaml:"type_handlers"`

	ClassOpts map[string]ClassOpts `yaml:"class_opts"`

	// This content will be put at the top of the output Typescript file.
	// You would generally use this to import custom types.
	HeaderFrontmatter   string `yaml:"header_frontmatter"`
	BindingsFrontmatter string `yaml:"bindings_frontmatter"`

	IgnoredMethods []string `yaml:"ignored_methods"`

	GlobalForcedMethods []FnOpts `yaml:"global_methods"`

	GlobalTypeOutTransforms map[string]string `yaml:"global_type_out_transforms"`

	MethodArgCheckTransforms map[string]string `yaml:"method_arg_check_transforms"`

	MethodReturnTransforms map[string]string            `yaml:"method_return_transforms"`
	MethodArgTransforms    map[string]map[string]string `yaml:"method_arg_transforms"`

	GlobalVars  string            `yaml:"global_vars"`
	HelperFuncs map[string]string `yaml:"helper_funcs"`
}

type Config struct {
	Packages []*PackageConfig `yaml:"packages"`
}

func (c Config) PackageNames() []string {
	names := make([]string, len(c.Packages))

	for i, p := range c.Packages {
		names[i] = p.Path
	}
	return names
}

func (c Config) PackageConfig(packagePath string) *PackageConfig {
	for _, pc := range c.Packages {
		if pc.Path == packagePath {
			if pc.Indent == "" {
				pc.Indent = "  "
			}
			return pc
		}
	}
	log.Fatalf("Config not found for package %s", packagePath)
	return nil
}

func (c PackageConfig) IsMethodWrapped(className string, fnName string) bool {
	if v, ok := c.ClassOpts[className]; ok {
		for _, m := range v.Methods {
			if strings.EqualFold(m, fnName) {
				return true
			}
		}
	}
	return false
}

func (c PackageConfig) IsFieldWrapped(className string, fnName string) bool {
	if v, ok := c.ClassOpts[className]; ok {
		for _, f := range v.Fields {
			if strings.EqualFold(f, fnName) {
				return true
			}
		}
	}
	return false
}

func (c PackageConfig) TypeHasHandler(name string) *TypeHandler {
	var handler *TypeHandler
	for hName, h := range c.TypeHandlers {
		if strings.EqualFold(hName, name) {
			handler = &h
			break
		}
	}
	return handler
}

func (c PackageConfig) IsMethodIgnored(name string) bool {
	for _, n := range c.IgnoredMethods {
		if strings.EqualFold(n, name) {
			return true
		}
	}
	return false
}

func (c PackageConfig) ResolvedBindingsOutPath(packageDir string) string {
	if c.BindingsOutPath == "" {
		return filepath.Join(packageDir, defaultOutBindingsFileName)
	} else if !strings.HasSuffix(c.BindingsOutPath, ".cc") {
		return filepath.Join(c.BindingsOutPath, defaultOutBindingsFileName)
	}
	return c.BindingsOutPath
}

func (c PackageConfig) ResolvedHeaderOutPath(packageDir string) string {
	if c.HeaderOutPath == "" {
		return filepath.Join(packageDir, defaultOutHeaderFileName)
	} else if !strings.HasSuffix(c.HeaderOutPath, ".h") {
		return filepath.Join(c.HeaderOutPath, defaultOutHeaderFileName)
	}
	return c.HeaderOutPath
}
