package napi

import (
	"fmt"
	"strings"
)

func (g *PackageGenerator) GetArgData(args *[]*CPPArg) *[]GenArgData {
	arg_data := &[]GenArgData{}
	for i, a := range *args {
		*arg_data = append(*arg_data, a.ParseArgData(g, *a.Name, i))
	}
	return arg_data
}

func (g *PackageGenerator) WriteTypedArrayGetter(sb *strings.Builder, arg_name string, idx int, array_info *TypedArrayHelpers, to_vector ...bool) {
	target_vector := false
	if len(to_vector) > 0 {
		target_vector = to_vector[0]
	}
	arr_type := array_info.NapiType.GetArrayType().String()
	tmp_name := GetPrefixedVarName(arg_name)
	sb.WriteString(fmt.Sprintf("Napi::TypedArrayOf<%s> %s = info[%d].As<Napi::TypedArrayOf<%s>>();\n", arr_type, tmp_name, idx, arr_type))
	out_name := arg_name
	if target_vector {
		out_name = GetPrefixedVarName(arg_name, "ptr")
	}
	needs_cast := new(string)
	if array_info.NeedsCast {
		*needs_cast = array_info.NativeType
	}
	WriteMaybePointerCast(sb, out_name, fmt.Sprintf("%s.Data()", tmp_name), needs_cast)
}

func (g *PackageGenerator) WriteOptimizedTypedArray(sb *strings.Builder, arg_name string, idx int, vector_info *TypedArrayHelpers, stl_type STLType) {
	ptr_name := GetPrefixedVarName(arg_name, "ptr")
	len_name := GetPrefixedVarName(arg_name, "len")
	g.WriteTypedArrayGetter(sb, arg_name, idx, vector_info, true)
	sb.WriteString(fmt.Sprintf("size_t %s = %s.ElementLength();\n", len_name, GetPrefixedVarName(arg_name)))
	// fmt.Println("stl_type", stl_type)
	sb.WriteString(fmt.Sprintf("%s<%s", stl_type, vector_info.NativeType))
	if stl_type == "std::vector" {
		sb.WriteString(fmt.Sprintf(">%s(%s, %s + %s);\n", arg_name, ptr_name, ptr_name, len_name))
	} else {
		sb.WriteString(fmt.Sprintf(", %s> %s;\n", len_name, arg_name))
		sb.WriteString(fmt.Sprintf("std::copy_n(std::begin(%s), %s, std::begin(%s));\n", ptr_name, len_name, arg_name))
	}
}

func (g *PackageGenerator) WriteBigIntGetter(sb *strings.Builder, arg_name string, arg_getter NapiTypeGetter, idx int, needs_cast *string, is_void bool, opt_array_name ...string) {
	arr_name := "info"
	if len(opt_array_name) > 0 {
		arr_name = opt_array_name[0]
	}
	loss_var := GetPrefixedVarName("is", arg_name, "lossless")
	sb.WriteString(fmt.Sprintf("bool %s = true;\n", loss_var))
	native_type := "int64_t"
	if needs_cast != nil {
		native_type = *needs_cast
	} else if arg_getter == Uint64Value {
		native_type = "uint64_t"
	}
	WriteMaybeCast(sb, arg_name, fmt.Sprintf("%s[%d].As<Napi::BigInt>().%s(&%s)", arr_name, idx, *arg_getter.String(), loss_var), needs_cast)
	g.WriteErrorHandler(sb, fmt.Sprintf("!%s", loss_var), fmt.Sprintf("failed to losslessly convert `%s` from typeof `bigint` to `%s`", arg_name, native_type), 1, false)
}

func WriteNumberGetter(sb *strings.Builder, arg_name string, arg_getter NapiTypeGetter, idx int, needs_cast *string, opt_array_name ...string) {
	arr_name := "info"
	if len(opt_array_name) > 0 {
		arr_name = opt_array_name[0]
	}
	WriteMaybeCast(sb, arg_name, fmt.Sprintf("%s[%d].As<Napi::Number>().%s()", arr_name, idx, *arg_getter.String()), needs_cast)
}

func WriteDateGetter(sb *strings.Builder, arg_name string, idx int, opt_array_name ...string) {
	arr_name := "info"
	if len(opt_array_name) > 0 {
		arr_name = opt_array_name[0]
	}
	sb.WriteString(fmt.Sprintf("time_t %s = static_cast<time_t>(%s[%d].As<Napi::Date>().ValueOf() / 1000);\n", arg_name, arr_name, idx))
}

func WriteStringGetter(sb *strings.Builder, arg_name string, idx int, arg_getter NapiTypeGetter, stl_type STLType, opt_array_name ...string) {
	is_char := stl_type == std_string
	is_u16char := stl_type == std_u16string
	needs_char_array := is_char || is_u16char
	arr_name := "info"
	if len(opt_array_name) > 0 {
		arr_name = opt_array_name[0]
	}
	used_name := arg_name
	if needs_char_array {
		used_name = GetPrefixedVarName("std_str", arg_name)
	}
	sb.WriteString(fmt.Sprintf("auto %s = %s[%d].As<Napi::String>().%s();\n", used_name, arr_name, idx, *arg_getter.String()))
	if needs_char_array {
		sb.WriteString(fmt.Sprintf("auto* %s = %s.c_str();\n", arg_name, used_name))
	}
}

func WriteEnumGetter(sb *strings.Builder, arg_name string, enum_name string, idx int) {
	WriteNumberGetter(sb, arg_name, Int32Value, idx, &enum_name)
}

func WriteExternalGetter(sb *strings.Builder, arg_name string, type_name string, idx int, opt_array_name ...string) {
	used_array_name := "info"
	if len(opt_array_name) > 0 {
		used_array_name = opt_array_name[0]
	}
	sb.WriteString(fmt.Sprintf("auto* %s = UnExternalize<%s>(%s[%d]);\n", arg_name, type_name, used_array_name, idx))
}

func WriteBooleanGetter(sb *strings.Builder, arg_name string, idx int, opt_array_name ...string) {
	used_array_name := "info"
	if len(opt_array_name) > 0 {
		used_array_name = opt_array_name[0]
	}
	sb.WriteString(fmt.Sprintf("bool %s = %s[%d].As<Napi::Boolean>().Value();\n", arg_name, used_array_name, idx))
}

func (g *PackageGenerator) WriteArrayGetter(sb *strings.Builder, arg_name string, idx int, stl_type STLType, raw_type *CPPType, opt_array_name ...string) {
	used_array_name := "info"
	if len(opt_array_name) > 0 {
		used_array_name = opt_array_name[0]
	}
	tmp_arr_name := GetPrefixedVarName(arg_name, "js", "arr")
	sb.WriteString(fmt.Sprintf("auto %s = %s[%d].As<Napi::Array>();\n", tmp_arr_name, used_array_name, idx))
	len_name := GetPrefixedVarName(arg_name, "js", "arr", "len")
	sb.WriteString(fmt.Sprintf("size_t %s = %s.Length();\n", len_name, tmp_arr_name))
	template_type := raw_type.Template.Args[0].GetFullType()
	sb.WriteString(fmt.Sprintf("%s<%s", stl_type, template_type))
	if stl_type == std_array {
		sb.WriteString(fmt.Sprintf(", %s", len_name))
	}
	sb.WriteString(fmt.Sprintf("> %s;\n", arg_name))
	if stl_type == std_vector {
		sb.WriteString(fmt.Sprintf("%s.reserve(%s);\n", arg_name, len_name))
	}
	sb.WriteString(fmt.Sprintf("for (size_t i = 0; i < %s; ++i) {\n", len_name))
	tmp_item_name := GetPrefixedVarName("arr_item")
	sb.WriteString(fmt.Sprintf("Napi::Value %s = %s[i];\n", tmp_item_name, tmp_arr_name))
	helpers := g.GetTypeHelpers(*raw_type.Template.Args[0].Name)

	napi_type := helpers.NapiType.String()
	pointer_str := ""
	if helpers.NapiType == External {
		pointer_str = "*"
		*napi_type += fmt.Sprintf("<%s>", template_type)
	}
	if stl_type == std_vector {
		sb.WriteString(fmt.Sprintf("%s.emplace_back(%s(%s.As<Napi::%s>().%s()));\n", arg_name, pointer_str, tmp_item_name, *napi_type, *helpers.NapiGetter.String()))
	} else {
		sb.WriteString(fmt.Sprintf("%s[i] = %s(%s.As<Napi::%s>().%s());\n", arg_name, pointer_str, tmp_item_name, *napi_type, *helpers.NapiGetter.String()))
	}

	sb.WriteString("}\n")
}

// TODO: improve default handling more complex types
func (g *PackageGenerator) WriteArgGetters(sb *strings.Builder, args *[]GenArgData, method_name string, is_void bool) {
	for i, a := range *args {
		// TODO: handle argument transforms
		isTransformed, transform := g.conf.IsArgTransform(method_name, a.Name)
		if !isTransformed {
			needs_cast := new(string)
			if a.NeedsCast {
				*needs_cast = a.NativeType
			}
			// TODO: handle string types
			switch a.NapiType {
			// handle boolean
			case Boolean:
				WriteBooleanGetter(sb, a.Name, a.Idx)
			// handle string
			case String:
				WriteStringGetter(sb, a.Name, a.Idx, a.NapiGetter, a.STLType)
			// handle number
			case Number:
				WriteNumberGetter(sb, a.Name, a.NapiGetter, a.Idx, needs_cast)
			// handle `Array<T>` -> `std::vector<T>` or `std::array<T>` (where `T` is not primitive)
			case Array:
				g.WriteArrayGetter(sb, a.Name, a.Idx, a.STLType, a.RawType)
			// handle bigint
			case BigInt:
				g.WriteBigIntGetter(sb, a.Name, a.NapiGetter, a.Idx, needs_cast, is_void)
			// handle native enum (passed as number)
			case NumberEnum:
				WriteEnumGetter(sb, a.Name, *needs_cast, a.Idx)
			// handle date
			case Date:
				WriteDateGetter(sb, a.Name, a.Idx)
			// handle typed arrays
			case TypedArray, Int8Array, Uint8Array, Int16Array, Uint16Array, Int32Array, Uint32Array, Float32Array, Float64Array, BigInt64Array, BigUint64Array:
				if a.STLType == "" {
					g.WriteTypedArrayGetter(sb, a.Name, a.Idx, a.TypedArrayInfo)
				} else {
					g.WriteOptimizedTypedArray(sb, a.Name, a.Idx, a.TypedArrayInfo, a.STLType)
				}
			// handle `Array<[T1, T2]>` -> `std::vector<std::pair<T1, T2>>`
			case PairArray:
				g.WritePairArrayHandler(sb, a.RawType, a.Name, a.Idx)
				// handle `[T1, T2]` -> `std::pair<T1, T2>`
			case Pair:
				g.WritePairHandler(sb, a.RawType, a.Name, a.Idx, is_void)
			// handle returning pointer to JS
			case External:
				WriteExternalGetter(sb, a.Name, a.NativeType, a.Idx)
			}
		} else {
			for j, val := range *args {
				*transform = strings.ReplaceAll(*transform, fmt.Sprintf("/arg_%d/", j), val.Name)
			}
			sb.WriteString(strings.ReplaceAll(*transform, "/arg/", fmt.Sprintf("info[%d]", i)))
		}
	}
}
