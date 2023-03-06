package napi

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/cpp"
)

const defaultOutBindingsFileName = "bindings.cc"
const defaultOutHeaderFileName = "bindings.h"
const defaultOutJsWrapperFileName = "gen_method_shims.js"

type TypeHandler struct {
	OutType string `yaml:"out_type"`
	OutVar  string `yaml:"out_var"`
	Handler string `yaml:"handler"`
}

type MethodTransforms struct {
	ArgCount           *int              `yaml:"arg_count"`
	ArgCheckTransforms string            `yaml:"arg_check_transforms"`
	ReturnTransforms   string            `yaml:"return_transforms"`
	ArgTransforms      map[string]string `yaml:"arg_transforms"`
}

type TypeMap struct {
	TSType       string `yaml:"ts"`
	NeedsWrapper bool   `yaml:"needs_wrapper"`
	NapiType     string `yaml:"napi"`
	NativeType   string `yaml:"native_type"`
	CastsTo      string `yaml:"casts_to"`
	CastNapi     string `yaml:"cast_napi"`
}

type FnArg struct {
	Name    string `yaml:"name"`
	TSType  string `yaml:"ts_type"`
	Default string `yaml:"default"`
}

type FnOpts struct {
	Name          string   `yaml:"name"`
	Args          []FnArg  `yaml:"args"`
	JSWrapperName string   `yaml:"js_wrapper_name"`
	JSWrapperAlts []string `yaml:"js_wrapper_alts"`
	FnBody        string   `yaml:"body"`
	IsVoid        bool     `yaml:"is_void"`
	TSReturnType  string   `yaml:"ts_return_type"`
}

type ClassOpts struct {
	IgnoredFields     []string `yaml:"ignored_fields"`
	ExternalFinalizer string   `yaml:"ext_finalizer_transform"`
	BytesUsed         *string  `yaml:"bytes_used"`
	ForcedMethods     []FnOpts `yaml:"forced_methods"`
	PathToImpl        string   `yaml:"path_to_impl"`
}

type JSWrapperOpts struct {
	AddonPath       string `yaml:"addon_path"`
	FrontMatter     string `yaml:"front_matter"`
	WrapperOutPath  string `yaml:"wrapper_out_path"`
	ShimFrontMatter string `yaml:"shim_front_matter"`
	// specifies whether gen JS/TS wrapper code
	EnvType string `yaml:"env_type"`
}

type GroupedMethodTransforms struct {
	AppliesTo        []string          `yaml:"applies_to"`
	ReturnTransforms string            `yaml:"return_transforms"`
	ArgTransforms    map[string]string `yaml:"arg_transforms"`
}

type PackageConfig struct {
	// The package path just like you would import it in Go
	Path       string `yaml:"path"`
	LibRootDir string `yaml:"lib_root_dir"`

	TrackExternalMemory *string `yaml:"track_external_memory"`

	// Where this output should be written to.
	// If you specify a folder it will be written to a file `index.ts` within that folder. By default it is written into the Golang package folder.
	BindingsOutPath   string        `yaml:"bindings_out_path"`
	HeaderOutPath     string        `yaml:"header_out_path"`
	JSWrapperOpts     JSWrapperOpts `yaml:"js_wrapper_opts"`
	PathToForcedLogic string        `yaml:"path_to_forced_logic"`

	// Customize the indentation (use \t if you want tabs)
	Indent string `yaml:"indent"`

	GroupedMethodTransforms []GroupedMethodTransforms `yaml:"grouped_method_transforms"`

	// Specify your own custom type translations, useful for custom types, `time.Time` and `null.String`.
	// Be default unrecognized types will be output as `any /* name */`.
	TypeMappings map[string]TypeMap `yaml:"type_mappings"`

	TypeHandlers map[string]TypeHandler `yaml:"type_handlers"`

	ClassOpts map[string]ClassOpts `yaml:"class_opts"`

	HeaderFrontmatter   string `yaml:"header_frontmatter"`
	BindingsFrontmatter string `yaml:"bindings_frontmatter"`

	IgnoredMethods []string `yaml:"ignored_methods"`

	GlobalForcedMethods []FnOpts `yaml:"global_methods"`

	MethodTransforms map[string]MethodTransforms `yaml:"method_transforms"`

	GlobalVars  string            `yaml:"global_vars"`
	HelperFuncs map[string]string `yaml:"helper_funcs"`
}

type Config struct {
	Packages []*PackageConfig `yaml:"packages"`
}

func (c Config) LoadForcedLogic() {
	for _, p := range c.Packages {
		if p.PathToForcedLogic != "" {
			input, err := os.ReadFile(p.PathToForcedLogic)
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
			for name, class := range p.ClassOpts {
				methods := *parseScopedFnBlock(n, input, fmt.Sprintf("%s_forced_methods", name))
				class.ForcedMethods = methods
				p.ClassOpts[name] = class
			}
			p.GlobalForcedMethods = *parseScopedFnBlock(n, input, "exported_global_methods")
			p.BindingsFrontmatter = parseIncludes(n, input)
			p.GlobalVars = parseGlobalVars(n, input)
			p.HelperFuncs = parsePrivateHelpers(n, input)
		}
	}
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

func (c PackageConfig) GetBytesAccessor(name string) *string {
	var bytes_accessor *string = new(string)
	if v, ok := c.ClassOpts[name]; ok {
		bytes_accessor = v.BytesUsed
	}
	return bytes_accessor
}

func (c PackageConfig) IsReturnTransform(method *CPPMethod) (bool, bool, *string) {
	for _, t := range c.GroupedMethodTransforms {
		for _, a := range t.AppliesTo {
			if strings.EqualFold(a, *method.Name) {
				return true, true, &t.ReturnTransforms
			}
		}
	}
	if v, ok := c.MethodTransforms[*method.Name]; ok && v.ReturnTransforms != "" {
		return true, false, &v.ReturnTransforms
	}
	return false, false, nil
}

func (c PackageConfig) IsArgTransform(fnName string, argName string) (bool, *string) {
	if v, ok := c.MethodTransforms[fnName]; ok {
		if mv, ok := v.ArgTransforms[argName]; ok {
			return true, &mv
		}
	}
	for _, t := range c.GroupedMethodTransforms {
		for _, a := range t.AppliesTo {
			if strings.EqualFold(a, fnName) {
				if v, ok := t.ArgTransforms[argName]; ok {
					return true, &v
				}
			}
		}
	}
	return false, nil
}

func (c PackageConfig) IsEnvTS() bool {
	return c.JSWrapperOpts.EnvType == "ts" || strings.HasSuffix(c.JSWrapperOpts.WrapperOutPath, ".ts")
}

func (c PackageConfig) IsFieldIgnored(className string, fnName string) bool {
	if v, ok := c.ClassOpts[className]; ok {
		for _, f := range v.IgnoredFields {
			if strings.EqualFold(f, fnName) {
				return true
			}
		}
	}
	return false
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

func (c PackageConfig) ResolvedForcedLogicPath(packageDir string) string {
	if c.HeaderOutPath == "" {
		return filepath.Join(packageDir, defaultOutHeaderFileName)
	} else if !strings.HasSuffix(c.HeaderOutPath, ".h") {
		return filepath.Join(c.HeaderOutPath, defaultOutHeaderFileName)
	}
	return c.HeaderOutPath
}

func (c PackageConfig) ResolvedWrapperOutPath(packageDir string) string {
	conf_path := c.JSWrapperOpts.WrapperOutPath
	if conf_path == "" {
		return filepath.Join(packageDir, defaultOutJsWrapperFileName)
	} else if !strings.HasSuffix(conf_path, ".js") && !strings.HasSuffix(conf_path, ".mjs") && !strings.HasSuffix(conf_path, ".cjs") && !strings.HasSuffix(conf_path, ".ts") {
		return filepath.Join(conf_path, defaultOutJsWrapperFileName)
	}
	return conf_path
}

func (c PackageConfig) ResolvedWrappedEnumOutPath(packageDir string) string {
	conf_path := c.JSWrapperOpts.WrapperOutPath
	if conf_path == "" {
		return filepath.Join(packageDir, defaultOutJsWrapperFileName)
	} else if !strings.HasSuffix(conf_path, ".js") && !strings.HasSuffix(conf_path, ".mjs") && !strings.HasSuffix(conf_path, ".cjs") && !strings.HasSuffix(conf_path, ".ts") {
		return filepath.Join(conf_path, "gen_enums.js")
	}
	path_split := strings.Split(conf_path, "/")
	return strings.Join(path_split[:len(path_split)-1], "/") + fmt.Sprintf("/gen_enums.%s", strings.Split(path_split[len(path_split)-1], ".")[1])
}

func (c PackageConfig) ResolvedShimPath(packageDir string, val_name string) string {
	conf_path := c.JSWrapperOpts.WrapperOutPath
	if conf_path == "" {
		return filepath.Join(packageDir, defaultOutJsWrapperFileName)
	} else if !strings.HasSuffix(conf_path, ".js") && !strings.HasSuffix(conf_path, ".mjs") && !strings.HasSuffix(conf_path, ".cjs") && !strings.HasSuffix(conf_path, ".ts") {
		return filepath.Join(conf_path, fmt.Sprintf("gen_%s_methods_shim.js", val_name))
	}
	path_split := strings.Split(conf_path, "/")
	return strings.Join(path_split[:len(path_split)-1], "/") + fmt.Sprintf("/gen_%s_methods_shim.%s", val_name, strings.Split(path_split[len(path_split)-1], ".")[1])
}

func (c PackageConfig) ResolvedBindingsImportPath(packageDir string) string {
	packageDirDepth := len(strings.Split(packageDir, "/"))
	wrapperDirDepth := len(strings.Split(c.ResolvedWrapperOutPath(packageDir), "/"))
	depth_dif := (wrapperDirDepth - packageDirDepth) + 1
	sb := new(strings.Builder)
	for i := 0; i < depth_dif; i++ {
		sb.WriteString("../")
	}
	sb.WriteString(c.JSWrapperOpts.AddonPath)
	return sb.String()
}
