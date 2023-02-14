package napi

import (
	"fmt"
	"strings"
)

func (g *PackageGenerator) WritePairHandler(sb *strings.Builder, t *CPPType, name string, idx int) {
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
	h1 += fmt.Sprintf("_tmp_pair_item1.As<Napi::%s>().%s()", helpers_1.NapiType, *helpers_1.NapiGetter)
	if h1_needsCast {
		h1 += ")"
	}
	if h2_needsCast {
		h2 += fmt.Sprintf("static_cast<%s>(", *helpers_2.CastTo)
	}
	h2 += fmt.Sprintf("_tmp_pair_item2.As<Napi::%s>().%s()", helpers_2.NapiType, *helpers_2.NapiGetter)
	if h2_needsCast {
		h2 += ")"
	}
	g.writeIndent(sb, 2)
	sb.WriteString(fmt.Sprintf("%s.emplace_back(%s, %s);\n", name, h1, h2))
	g.writeIndent(sb, 1)
	sb.WriteString("}\n")
}

func (g *PackageGenerator) writeJsArrayToVectorFn(sb *strings.Builder) {
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
	sb.WriteString("auto v = static_cast<const T>(val.As<Napi::Number>().Int64Value());\n")
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
}

func (g *PackageGenerator) writeHelpers(w *strings.Builder) {
	w.WriteString("// non-exported helpers\n")
	g.writeJsArrayToVectorFn(w)
	g.writeArrayBufferDeleter(w)
	hasUnexternalizer := false
	for name, c := range g.ParsedData.Classes {
		if c.Decl != nil {
			g.writeClassDeleter(w, c, name)
			g.writeClassExternalizer(w, c, name)
			if !hasUnexternalizer {
				g.writeClassUnExternalizer(w)
				hasUnexternalizer = true
			}

			if c.FieldDecl != nil {
				for _, f := range *c.FieldDecl {
					g.writeClassField(w, f, name)
				}
			}

			if v, ok := g.conf.ClassOpts[name]; ok && len(v.ForcedMethods) > 0 {
				for _, f := range v.ForcedMethods {
					w.WriteString(strings.Replace(f.FnBody, f.Name, "_"+f.Name, 1))
					w.WriteString("\n\n")
				}
			}
		}
	}
	for _, helper := range g.conf.HelperFuncs {
		w.WriteString(helper)
		w.WriteByte('\n')
	}
}

func (g *PackageGenerator) writeClassDeleter(sb *strings.Builder, class *CPPClass, name string) {
	sb.WriteString(fmt.Sprintf("static inline void Delete%s(Napi::Env env, void* ptr) {\n", name))
	g.writeIndent(sb, 1)
	sb.WriteString(fmt.Sprintf("auto* val = static_cast<%s::%s*>(ptr);\n", *g.NameSpace, name))
	if v, ok := g.conf.ClassOpts[name]; ok && v.ExternalFinalizer != "" {
		sb.WriteString(strings.ReplaceAll(v.ExternalFinalizer, "/this/", "val"))
	}
	g.writeIndent(sb, 1)
	sb.WriteString("delete val;\n")
	sb.WriteString("}\n\n")
}

func (g *PackageGenerator) writeArrayBufferDeleter(sb *strings.Builder) {
	sb.WriteString("template <typename T>\n")
	sb.WriteString("static inline void DeleteArrayBuffer(Napi::Env env, void* /*data*/, std::vector<T>* hint) {\n")
	g.writeIndent(sb, 1)
	sb.WriteString("size_t bytes = hint->size() * sizeof(T);\n")
	g.writeIndent(sb, 1)
	sb.WriteString("std::unique_ptr<std::vector<T>> vectorPtrToDelete(hint);\n")
	g.writeIndent(sb, 1)
	sb.WriteString("Napi::MemoryManagement::AdjustExternalMemory(env, -bytes);\n")
	sb.WriteString("}\n\n")
}

func (g *PackageGenerator) writeClassExternalizer(sb *strings.Builder, class *CPPClass, name string) {
	sb.WriteString(fmt.Sprintf("static inline Napi::External<%s::%s> Externalize%s(Napi::Env env, %s::%s* ptr) {\n", *g.NameSpace, name, name, *g.NameSpace, name))
	g.writeIndent(sb, 1)
	sb.WriteString(fmt.Sprintf("return Napi::External<%s::%s>::New(env, ptr, Delete%s);\n", *g.NameSpace, name, name))
	sb.WriteString("}\n\n")
}

func (g *PackageGenerator) writeClassUnExternalizer(sb *strings.Builder) {
	sb.WriteString("template <typename T>\n")
	sb.WriteString("static inline T* UnExternalize(Napi::Value val) {\n")
	g.writeIndent(sb, 1)
	sb.WriteString("return val.As<Napi::External<T>>().Data();\n")
	sb.WriteString("}\n\n")
}
