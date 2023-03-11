package napi

import (
	"fmt"
	"strings"
)

func WriteBooleanReturn(sb *strings.Builder, res_name string) {
	sb.WriteString("return ")
	init_napi_boolean(sb, res_name)
}

func WriteNumberReturn(sb *strings.Builder, res_name string, needs_cast *string) {
	sb.WriteString("return ")
	init_napi_number(sb, res_name, needs_cast)
}

func WriteBigIntReturn(sb *strings.Builder, res_name string, needs_cast *string) {
	sb.WriteString("return ")
	init_napi_bigint(sb, res_name, needs_cast)
}

func WriteStringReturn(sb *strings.Builder, res_name string) {
	sb.WriteString("return ")
	init_napi_string(sb, res_name)
}

func (g *PackageGenerator) WriteExternalReturn(sb *strings.Builder, res_name string, scoped_name string) {
	new_instance_name := GetPrefixedVarName(res_name, "new", stripNameSpace(scoped_name))
	sb.WriteString(fmt.Sprintf("auto* %s = new %s(%s);\n", new_instance_name, scoped_name, res_name))
	bytes_accessor := g.conf.GetBytesAccessor(stripNameSpace(scoped_name))
	if bytes_accessor != nil {
		write_mem_adjustment(sb, fmt.Sprintf("%s->%s()", new_instance_name, *bytes_accessor), g.conf.TrackExternalMemory)
	}
	sb.WriteString(fmt.Sprintf("return Externalize%s(env, %s);\n", stripNameSpace(scoped_name), new_instance_name))
}

func WritePairReturn(sb *strings.Builder, res_name string, template_args []TemplateNapiHandlers) {
	out_var_name := GetPrefixedVarName(res_name, "js", "tform")
	sb.WriteString(fmt.Sprintf("Napi::Array %s = ", out_var_name))
	init_napi_array(sb, 2)

	for i, t := range template_args {
		sb.WriteString(fmt.Sprintf("%s[%d] = ", out_var_name, i))
		pair_accessor := "first"
		if i == 1 {
			pair_accessor = "second"
		}
		switch t.NapiType {
		// TODO: optimize and pass as TypedArray
		case Number, NumberEnum:
			init_napi_number(sb, fmt.Sprintf("%s.%s", res_name, pair_accessor), t.NeedsCast)
		}
	}

	sb.WriteString(fmt.Sprintf("return %s;\n", out_var_name))
}

func WriteArrayReturn(sb *strings.Builder, res_name string, template_args []TemplateNapiHandlers) {
	out_var_name := GetPrefixedVarName(res_name, "js", "tform")
	len_var_name := fmt.Sprintf("%s_len", res_name)
	sb.WriteString(fmt.Sprintf("auto %s = %s.size();\n", len_var_name, res_name))
	sb.WriteString(fmt.Sprintf("Napi::Array %s = ", out_var_name))
	init_napi_array(sb, len_var_name)
	sb.WriteString(fmt.Sprintf("for (size_t i = 0; i < %s; ++i) {\n", len_var_name))
	for _, t := range template_args {
		sb.WriteString(fmt.Sprintf("%s[i] = ", out_var_name))
		switch t.NapiType {
		// TODO: optimize and pass as TypedArray
		case Number, NumberEnum:
			init_napi_number(sb, fmt.Sprintf("%s[i]", res_name), t.NeedsCast)
		case String:
			init_napi_string(sb, fmt.Sprintf("%s[i]", res_name))
		default:
			panic(fmt.Sprintf("unsupported array return type: `%s`; please file issue w info to add support", *t.NapiType.String()))
		}
	}
	sb.WriteString("}\n")

	sb.WriteString(fmt.Sprintf("return %s;\n", out_var_name))
}

func (g *PackageGenerator) WriteTypedArrayReturn(sb *strings.Builder, res_name string, r GenReturnData, len ...interface{}) {
	g.WriteVectorArrayBufferDeleter()
	used_var := res_name
	if r.STLType == "" && r.TypedArrayInfo.NeedsCast {
		used_var = fmt.Sprintf("%s_cast", res_name)
		sb.WriteString(fmt.Sprintf("auto* %s = static_cast<%s*>(%s);\n", used_var, *r.NeedsCast, res_name))
		// TODO: need better means of getting len
		vector_wrapper := fmt.Sprintf("%s_vec", res_name)
		sb.WriteString(fmt.Sprintf("std::unique_ptr<std::vector<%s>> %s = std::make_unique<std::vector<%s>>(%s, %s + %s);\n", *r.NeedsCast, vector_wrapper, *r.NeedsCast, used_var, used_var, fmt_num_or_string(len[0])))
		used_var = vector_wrapper
	} else if r.TypedArrayInfo.NeedsCast {
		used_var = fmt.Sprintf("%s_cast", res_name)
		sb.WriteString(fmt.Sprintf("std::unique_ptr<std::vector<%s>> %s;\n", *r.NeedsCast, used_var))
	} else {
		used_var = fmt.Sprintf("%s_vec", res_name)
		sb.WriteString(fmt.Sprintf("std::unique_ptr<std::vector<%s>> %s;\n", *r.NeedsCast, used_var))
		sb.WriteString(fmt.Sprintf("%s.reset(&%s);\n", used_var, res_name))
	}
	// element length of `std::vector`
	elem_len_name := fmt.Sprintf("%s_elem_len", used_var)
	sb.WriteString(fmt.Sprintf("auto %s = %s->size();\n", elem_len_name, used_var))

	if r.STLType != "" && r.TypedArrayInfo.NeedsCast {
		sb.WriteString(fmt.Sprintf("%s->reserve(%s);\n", used_var, elem_len_name))
		sb.WriteString(fmt.Sprintf("for (size_t i = 0; i < %s; ++i) {\n", elem_len_name))
		sb.WriteString(fmt.Sprintf("(*%s)[i] = static_cast<%s>(%s[i]);\n", used_var, *r.NeedsCast, res_name))
		sb.WriteString("}\n")
	}

	// byte length of `std::vector`
	byte_len_name := fmt.Sprintf("%s_byte_len", used_var)
	sb.WriteString(fmt.Sprintf("auto %s = %s * sizeof(%s);\n", byte_len_name, elem_len_name, *r.NeedsCast))

	// helpers to write external array w finalizer
	finalizer := fmt.Sprintf("DeleteArrayBufferFromVector<%s>", *r.NeedsCast)
	hint := fmt.Sprintf("%s.get()", used_var)
	out_var_name := GetPrefixedVarName(res_name, "js", "tform")
	sb.WriteString(fmt.Sprintf("Napi::ArrayBuffer %s = ", out_var_name))
	init_napi_arraybuffer(sb, fmt.Sprintf("%s->data()", used_var), byte_len_name, &finalizer, &hint)
	sb.WriteString(fmt.Sprintf("%s.release();\n", used_var))

	write_mem_adjustment(sb, byte_len_name, g.conf.TrackExternalMemory)
	sb.WriteString(fmt.Sprintf("return Napi::%s::New(env, %s, %s, 0);\n", *r.NapiType.String(), elem_len_name, out_var_name))
}

func (g *PackageGenerator) WriteStdArrayToArrayBuffer(sb *strings.Builder, res_name string, r GenReturnData) {
	g.WriteVectorArrayBufferDeleter()
	used_var := res_name
	if r.TypedArrayInfo.NeedsCast {
		used_var = fmt.Sprintf("%s_cast", res_name)
		sb.WriteString(fmt.Sprintf("std::unique_ptr<std::array<%s, %s>> %s = std::make_unique<std::array<%s, %s>>(%s)", *r.NeedsCast, *r.RawType.Template.Args[1].Name, used_var, *r.NeedsCast, *r.RawType.Template.Args[1].Name, res_name))
	} else {
		used_var = fmt.Sprintf("%s_std_array", res_name)
		sb.WriteString(fmt.Sprintf("std::unique_ptr<std::array<%s, %s>> %s;\n", *r.NeedsCast, *r.RawType.Template.Args[1].Name, used_var))
		sb.WriteString(fmt.Sprintf("%s.reset(&%s);\n", used_var, res_name))
	}

	// byte length of `std::vector`
	byte_len_name := fmt.Sprintf("%s_byte_len", used_var)
	sb.WriteString(fmt.Sprintf("auto %s = %s * sizeof(%s);\n", byte_len_name, *r.RawType.Template.Args[1].Name, *r.NeedsCast))

	// helpers to write external array w finalizer
	finalizer := &strings.Builder{}
	finalizer.WriteString(fmt.Sprintf("[](Napi::Env env, void* /*data*/, std::array<%s, %s>* hint) {\n", *r.NeedsCast, *r.RawType.Template.Args[1].Name))
	finalizer.WriteString(fmt.Sprintf("std::unique_ptr<std::array<%s, %s>> arrayPtrToDelete(hint);\n", *r.NeedsCast, *r.RawType.Template.Args[1].Name))
	write_mem_adjustment(finalizer, fmt.Sprintf("-(%s * sizeof(%s))", *r.RawType.Template.Args[1].Name, *r.NeedsCast), g.conf.TrackExternalMemory)
	finalizer.WriteString("}\n")
	finalizer_string := finalizer.String()
	hint := fmt.Sprintf("%s.get()", used_var)
	out_var_name := GetPrefixedVarName(res_name, "js", "tform")
	sb.WriteString(fmt.Sprintf("Napi::ArrayBuffer %s = ", out_var_name))
	init_napi_arraybuffer(sb, *r.RawType.Template.Args[1].Name, byte_len_name, &finalizer_string, &hint)
	sb.WriteString(fmt.Sprintf("%s.release();\n", used_var))

	write_mem_adjustment(sb, byte_len_name, g.conf.TrackExternalMemory)
	sb.WriteString(fmt.Sprintf("return Napi::%s::New(env, %s, %s, 0);\n", *r.NapiType.String(), *r.RawType.Template.Args[1].Name, out_var_name))
}

func (g *PackageGenerator) WriteReturnVal(sb *strings.Builder, r GenReturnData, is_void bool, gen_result_name *string) {
	if gen_result_name == nil || is_void {
		return
	}
	// fmt.Printf("NapiType: %q; STLType: %q\n", *r.NapiType.String(), r.STLType)
	switch r.NapiType {
	// handle boolean
	case Boolean:
		WriteBooleanReturn(sb, *gen_result_name)
	// handle number
	case NumberEnum, Number:
		WriteNumberReturn(sb, *gen_result_name, r.NeedsCast)
	// handle bigint
	case BigInt:
		WriteBigIntReturn(sb, *gen_result_name, r.NeedsCast)
	// handle string
	case String:
		WriteStringReturn(sb, *gen_result_name)
	// handle returning wrapped pointer to struct/class
	case External:
		g.WriteExternalReturn(sb, *gen_result_name, r.NativeType)
	// handle `std::pair<T1, T2>` -> `[T1, T2]`
	// TODO: optimize, leveraging `TypedArray` where possible
	case Pair:
		template_args := []TemplateNapiHandlers{}
		for _, t := range r.RawType.Template.Args {
			template_args = append(template_args, t.GetNapiHandlers(g))
		}
		WritePairReturn(sb, *gen_result_name, template_args)
	// handle `std::vector<T>` & `std::array<T>` -> `Array<T>`
	// TODO: optimize, leveraging `TypedArray` where possible
	case Array:
		template_args := []TemplateNapiHandlers{}
		for _, t := range r.RawType.Template.Args {
			template_args = append(template_args, t.GetNapiHandlers(g))
		}
		WriteArrayReturn(sb, *gen_result_name, template_args)
	// handle typed arrays
	case TypedArray, Int8Array, Uint8Array, Int16Array, Uint16Array, Int32Array, Uint32Array, Float32Array, Float64Array, BigInt64Array, BigUint64Array:
		if r.STLType == "" || r.STLType == std_vector {
			g.WriteTypedArrayReturn(sb, *gen_result_name, r)
		} else if r.STLType == std_array {
			g.WriteStdArrayToArrayBuffer(sb, *gen_result_name, r)
		} else {
			panic(fmt.Sprintf("unhandled STLType: %q", r.STLType))
		}
		// TODO: need to handle the following types:
		/*
			// handle date
			case Date:
				WriteDateGetter(sb, a.Name, a.Idx)
			// handle typed arrays
			// handle `Array<[T1, T2]>` -> `std::vector<std::pair<T1, T2>>`
			case PairArray:
				g.WritePairHandler(sb, a.RawType, a.Name, a.Idx, is_void)
		*/
	}
}
