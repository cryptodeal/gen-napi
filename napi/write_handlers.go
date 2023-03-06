package napi

import (
	"fmt"
	"strings"
)

func (g *PackageGenerator) WritePairHandler(sb *strings.Builder, t *CPPType, name string, idx int, is_void bool, opt_arr_name ...string) {
	arr_name := "info"
	if len(opt_arr_name) > 0 {
		arr_name = opt_arr_name[0]
	}
	tmp_arr_name := GetPrefixedVarName(name, "arr")
	g.writeIndent(sb, 1)
	sb.WriteString(fmt.Sprintf("Napi::Array %s = %s[%d].As<Napi::Array>();\n", tmp_arr_name, arr_name, idx))
	g.WriteErrorHandler(sb, fmt.Sprintf("%s.Length() != 2", tmp_arr_name), fmt.Sprintf("`%s` expects a pair, but received an array of length != 2", name), 1, is_void)
	item1_name := GetPrefixedVarName(name, "item1")
	item2_name := GetPrefixedVarName(name, "item2")
	g.writeIndent(sb, 1)
	sb.WriteString(fmt.Sprintf("Napi::Value %s = %s[static_cast<size_t>(0)];\n", item1_name, tmp_arr_name))
	g.writeIndent(sb, 1)
	sb.WriteString(fmt.Sprintf("Napi::Value %s = %s[static_cast<size_t>(1)];\n", item2_name, tmp_arr_name))
	helpers_1 := g.GetTypeHelpers(*t.Template.Args[0].Name)
	helpers_2 := g.GetTypeHelpers(*t.Template.Args[1].Name)
	h1_needsCast := helpers_1.CastTo != nil
	h2_needsCast := helpers_2.CastTo != nil
	var h1, h2 string
	if h1_needsCast {
		h1 += fmt.Sprintf("static_cast<%s>(", *helpers_1.CastTo)
	}
	h1 += fmt.Sprintf("%s.As<Napi::%s>().%s()", item1_name, *helpers_1.NapiType.String(), *helpers_1.NapiGetter.String())
	if h1_needsCast {
		h1 += ")"
	}
	if h2_needsCast {
		h2 += fmt.Sprintf("static_cast<%s>(", *helpers_2.CastTo)
	}
	h2 += fmt.Sprintf("%s.As<Napi::%s>().%s()", item2_name, *helpers_2.NapiType.String(), *helpers_2.NapiGetter.String())
	if h2_needsCast {
		h2 += ")"
	}
	g.writeIndent(sb, 1)
	sb.WriteString(fmt.Sprintf("std::pair<%s, %s> %s(%s, %s);\n", *t.Template.Args[0].Name, *t.Template.Args[1].Name, name, h1, h2))
}

func (g *PackageGenerator) WritePairArrayHandler(sb *strings.Builder, t *CPPType, name string, idx int) {
	g.writeIndent(sb, 1)
	sb.WriteString(fmt.Sprintf("Napi::Array _tmp_array_%s = info[%d].As<Napi::Array>();\n", name, idx))
	g.writeIndent(sb, 1)
	sb.WriteString(fmt.Sprintf("auto _tmp_len_%s = _tmp_array_%s.Length();\n", name, name))
	g.writeIndent(sb, 1)
	sb.WriteString(fmt.Sprintf("std::vector<std::pair<%s, %s>> %s;\n", *t.Template.Args[0].Args[0].Name, *t.Template.Args[0].Args[1].Name, name))
	g.writeIndent(sb, 1)
	sb.WriteString(fmt.Sprintf("%s.reserve(_tmp_len_%s);\n", name, name))
	g.writeIndent(sb, 1)
	sb.WriteString(fmt.Sprintf("for (size_t i = 0; i < _tmp_len_%s; i++) {\n", name))
	g.writeIndent(sb, 2)
	sb.WriteString(fmt.Sprintf("Napi::Value _tmp_array_item = _tmp_array_%s[i];\n", name))
	g.writeIndent(sb, 2)
	sb.WriteString("Napi::Array _tmp_pair = _tmp_array_item.As<Napi::Array>();\n")
	g.writeIndent(sb, 2)
	sb.WriteString("size_t idx1 = 0, idx2 = 1;\n")
	g.writeIndent(sb, 2)
	sb.WriteString("Napi::Value _tmp_pair_item1 = _tmp_pair[idx1];\n")
	g.writeIndent(sb, 2)
	sb.WriteString("Napi::Value _tmp_pair_item2 = _tmp_pair[idx2];\n")

	// write handlers for each variable
	helpers_1 := g.GetTypeHelpers(*t.Template.Args[0].Args[0].Name)
	helpers_2 := g.GetTypeHelpers(*t.Template.Args[0].Args[1].Name)
	h1_needsCast := helpers_1.CastTo != nil
	h2_needsCast := helpers_2.CastTo != nil
	var h1, h2 string
	if h1_needsCast {
		h1 += fmt.Sprintf("static_cast<%s>(", *helpers_1.CastTo)
	}
	h1 += fmt.Sprintf("_tmp_pair_item1.As<Napi::%s>().%s()", *helpers_1.NapiType.String(), *helpers_1.NapiGetter.String())
	if h1_needsCast {
		h1 += ")"
	}
	if h2_needsCast {
		h2 += fmt.Sprintf("static_cast<%s>(", *helpers_2.CastTo)
	}
	h2 += fmt.Sprintf("_tmp_pair_item2.As<Napi::%s>().%s()", *helpers_2.NapiType.String(), *helpers_2.NapiGetter.String())
	if h2_needsCast {
		h2 += ")"
	}
	g.writeIndent(sb, 2)
	sb.WriteString(fmt.Sprintf("%s.emplace_back(%s, %s);\n", name, h1, h2))
	g.writeIndent(sb, 1)
	sb.WriteString("}\n")
}

func (g *PackageGenerator) WriteVectorHandler() {
	fn_name := "jsArrayToVector"
	if !g.GenHelperExists(fn_name) {
		sb := &strings.Builder{}
		sb.WriteString("template <typename T>\n")
		sb.WriteString("static inline std::vector<T> jsArrayToVector(Napi::Array arr, bool reverse, int invert) {\n")
		g.writeIndent(sb, 1)
		sb.WriteString("std::vector<T> out;\n")
		g.writeIndent(sb, 1)
		sb.WriteString("const size_t len = arr.Length();\n")
		g.writeIndent(sb, 1)
		sb.WriteString("out.reserve(len);\n")
		g.writeIndent(sb, 1)
		sb.WriteString("for(size_t i = 0; i < len; ++i) {\n")
		g.writeIndent(sb, 2)
		sb.WriteString("const auto idx = reverse ? len - i - 1 : i;\n")
		g.writeIndent(sb, 2)
		sb.WriteString("Napi::Value val = arr[idx];\n")
		g.writeIndent(sb, 2)
		sb.WriteString("auto v = static_cast<T>(val.As<Napi::Number>().Int64Value());\n")
		g.writeIndent(sb, 2)
		sb.WriteString("if (invert && v < 0) {\n")
		g.writeIndent(sb, 3)
		sb.WriteString("v = -v - 1;\n")
		g.writeIndent(sb, 2)
		sb.WriteString("} else if (invert) {\n")
		g.writeIndent(sb, 3)
		sb.WriteString("v = invert - v - 1;\n")
		g.writeIndent(sb, 2)
		sb.WriteString("}\n")
		g.writeIndent(sb, 2)
		sb.WriteString("out.emplace_back(v);\n")
		g.writeIndent(sb, 1)
		sb.WriteString("}\n")
		g.writeIndent(sb, 1)
		sb.WriteString("return out;\n")
		sb.WriteString("}\n\n")
		g.AddGenHelper(fn_name, sb.String())
	}
}

func (g *PackageGenerator) writeHelpers(w *strings.Builder) {
	w.WriteString("// non-exported helpers\n")
	for _, key := range g.GenHelperKeys {
		w.WriteString(g.GenHelpers[key])
		w.WriteByte('\n')
	}

	for _, helper := range g.conf.HelperFuncs {
		w.WriteString(helper)
		w.WriteByte('\n')
	}
}

func (g *PackageGenerator) writeClassDeleter(class *CPPClass, name string) {
	fn_name := fmt.Sprintf("Delete%s", name)
	if !g.GenHelperExists(fn_name) {
		sb := &strings.Builder{}
		sb.WriteString(fmt.Sprintf("static inline void Delete%s(Napi::Env env, void* ptr) {\n", name))
		g.writeIndent(sb, 1)
		sb.WriteString(fmt.Sprintf("auto* val = static_cast<%s::%s*>(ptr);\n", *class.NameSpace, name))
		if v, ok := g.conf.ClassOpts[name]; ok && v.ExternalFinalizer != "" {
			sb.WriteString(strings.ReplaceAll(v.ExternalFinalizer, "/this/", "val"))
		}
		g.writeIndent(sb, 1)
		sb.WriteString("delete val;\n")
		sb.WriteString("}\n\n")
		g.AddGenHelper(fn_name, sb.String())
	}
}

func (g *PackageGenerator) WriteVectorArrayBufferDeleter() {
	name := "DeleteArrayBufferFromVector"
	if !g.GenHelperExists(name) {
		sb := &strings.Builder{}
		sb.WriteString("template <typename T>\n")
		sb.WriteString(fmt.Sprintf("static inline void %s(Napi::Env env, void* /*data*/, std::vector<T>* hint) {\n", name))
		g.writeIndent(sb, 1)
		sb.WriteString("size_t bytes = hint->size() * sizeof(T);\n")
		g.writeIndent(sb, 1)
		sb.WriteString("std::unique_ptr<std::vector<T>> vectorPtrToDelete(hint);\n")
		write_mem_adjustment(sb, "-bytes", g.conf.TrackExternalMemory)
		sb.WriteString("}\n\n")
		g.AddGenHelper(name, sb.String())
	}
}

func (g *PackageGenerator) WriteArrayDeleter(stl_type STLType) {
	name := "DeleteArrayBuffer"
	if !g.GenHelperExists(name) {
		sb := &strings.Builder{}
		sb.WriteString("template <typename T>\n")
		sb.WriteString(fmt.Sprintf("static inline void %s(Napi::Env env, void* /*data*/, std::vector<T>* hint) {\n", name))
		g.writeIndent(sb, 1)
		sb.WriteString("size_t bytes = hint->size() * sizeof(T);\n")
		g.writeIndent(sb, 1)
		sb.WriteString("std::unique_ptr<std::vector<T>> vectorPtrToDelete(hint);\n")
		write_mem_adjustment(sb, "-bytes", g.conf.TrackExternalMemory)
		sb.WriteString("}\n\n")
		g.AddGenHelper(name, sb.String())
	}
}

func (g *PackageGenerator) writeClassExternalizer(class *CPPClass, name string) {
	fn_name := fmt.Sprintf("Externalize%s", name)
	if !g.GenHelperExists(fn_name) {
		sb := &strings.Builder{}
		sb.WriteString(fmt.Sprintf("static inline Napi::External<%s::%s> %s(Napi::Env env, %s::%s* ptr) {\n", *class.NameSpace, name, fn_name, *class.NameSpace, name))
		g.writeIndent(sb, 1)
		sb.WriteString(fmt.Sprintf("return Napi::External<%s::%s>::New(env, ptr, Delete%s);\n", *class.NameSpace, name, name))
		sb.WriteString("}\n\n")
		g.AddGenHelper(fn_name, sb.String())
	}
}

func (g *PackageGenerator) writeClassUnExternalizer() {
	fn_name := "UnExternalize"
	if !g.GenHelperExists(fn_name) {
		sb := &strings.Builder{}
		sb.WriteString("template <typename T>\n")
		sb.WriteString("static inline T* UnExternalize(Napi::Value val) {\n")
		g.writeIndent(sb, 1)
		sb.WriteString("return val.As<Napi::External<T>>().Data();\n")
		sb.WriteString("}\n\n")
		g.AddGenHelper(fn_name, sb.String())
	}
}
