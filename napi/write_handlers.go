package napi

import (
	"fmt"
	"strings"
)

func (g *PackageGenerator) writeHelpers(w *strings.Builder, classes map[string]*CPPClass) {
	if len(g.conf.HelperFuncs) > 0 {
		w.WriteString("// non-exported helpers\n")
		g.writeArrayBufferDeleter(w)
		hasUnexternalizer := false
		for name, c := range classes {
			if c.Decl != nil {
				g.writeClassDeleter(w, c, name)
				g.writeClassExternalizer(w, c, name)
				if !hasUnexternalizer {
					g.writeClassUnExternalizer(w)
					hasUnexternalizer = true
				}

				if c.FieldDecl != nil {
					for _, f := range *c.FieldDecl {
						g.writeClassField(w, f, name, classes)
					}
				}

				if v, ok := g.conf.ClassOpts[name]; ok && len(v.ForcedMethods) > 0 {
					for _, f := range v.ForcedMethods {
						w.WriteString(strings.Replace(f.FnBody, f.Name, "_"+f.Name, 1))
						w.WriteByte('\n')
					}
				}
			}
		}
		for _, helper := range g.conf.HelperFuncs {
			w.WriteString(helper)
			w.WriteByte('\n')
		}
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
