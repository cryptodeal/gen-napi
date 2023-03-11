package napi

import (
	"fmt"
	"strings"
)

type TypedArrayHelpers struct {
	NapiType   NapiType
	NativeType string
	NeedsCast  bool
}

func GetTypedArrayHandlers(var_type string) *TypedArrayHelpers {
	// pairs must be passed as `Array<[T1, T2]>`
	// cannot optimize by passing as TypedArray
	if var_type == "pair" || var_type == "std::pair" {
		return nil
	}
	optimized := &TypedArrayHelpers{
		NativeType: var_type,
	}
	switch var_type {
	case "float":
		optimized.NapiType = Float32Array
		return optimized
	case "double":
		optimized.NapiType = Float64Array
		return optimized
	case "uint8_t":
		optimized.NapiType = Uint8Array
		return optimized
	case "bool":
		optimized.NapiType = Uint8Array
		optimized.NeedsCast = true
		return optimized
	case "int8_t":
		optimized.NapiType = Int8Array
		return optimized
	case "uint16_t":
		optimized.NapiType = Uint16Array
		return optimized
	case "unsigned short":
		optimized.NapiType = Uint16Array
		optimized.NeedsCast = true
		return optimized
	case "int16_t":
		optimized.NapiType = Int16Array
		return optimized
	case "short", "signed short":
		optimized.NapiType = Int16Array
		optimized.NeedsCast = true
		return optimized
	case "uint32_t":
		optimized.NapiType = Uint32Array
		return optimized
	case "unsigned int", "unsigned long":
		optimized.NapiType = Uint32Array
		optimized.NeedsCast = true
		return optimized
	case "int32_t":
		optimized.NapiType = Int32Array
		return optimized
	case "long", "signed long", "int", "signed int":
		optimized.NapiType = Int32Array
		optimized.NeedsCast = true
		return optimized
	case "uint64_t":
		optimized.NapiType = BigUint64Array
		return optimized
	case "unsigned long long", "unsigned long long int", "size_t":
		optimized.NapiType = BigUint64Array
		optimized.NeedsCast = true
		return optimized
	case "int64_t":
		optimized.NapiType = BigInt64Array
		return optimized
	case "long long", "long long int", "signed long long", "signed long long int":
		optimized.NapiType = BigInt64Array
		optimized.NeedsCast = true
		return optimized
	}
	return nil
}

func WriteMaybeCast(sb *strings.Builder, name string, from string, cast_to *string) {
	sb.WriteString(fmt.Sprintf("auto %s = ", name))
	if cast_to != nil && *cast_to != "" {
		sb.WriteString(fmt.Sprintf("static_cast<%s>(", *cast_to))
	}
	sb.WriteString(from)
	if cast_to != nil && *cast_to != "" {
		sb.WriteByte(')')
	}
	sb.WriteString(";\n")
}

func WriteMaybePointerCast(sb *strings.Builder, name string, from string, cast_to *string) {
	sb.WriteString(fmt.Sprintf("auto* %s = ", name))
	if cast_to != nil && *cast_to != "" {
		sb.WriteString(fmt.Sprintf("static_cast<%s*>(", *cast_to))
	}
	sb.WriteString(from)
	if cast_to != nil && *cast_to != "" {
		sb.WriteByte(')')
	}
	sb.WriteString(";\n")
}

type GenReturnData struct {
	NapiType       NapiType
	NativeType     string
	NeedsCast      *string
	TypedArrayInfo *TypedArrayHelpers
	STLType        STLType
	JSType         string
	IsPointer      bool
	IsConst        bool
	RefDecl        *string
	RawType        *CPPType
}

type GenArgData struct {
	Name             string
	NapiType         NapiType
	NapiGetter       NapiTypeGetter
	NativeType       string
	NeedsCast        bool
	NeedsConstructor *string
	TypedArrayInfo   *TypedArrayHelpers
	STLType          STLType
	JSType           string
	DefaultValue     *string
	IsPointer        bool
	IsConst          bool
	RefDecl          *string
	Idx              int
	RawType          *CPPType
}

type TemplateNapiHandlers struct {
	NameSpace      *string
	NativeType     string
	IsPointer      bool
	TypedArrayInfo *TypedArrayHelpers
	NapiType       NapiType
	JSType         string
	STLType        STLType
	NeedsCast      *string
}

func (t *TemplateType) GetNapiHandlers(g *PackageGenerator) TemplateNapiHandlers {
	// fmt.Printf("Template Type: %s, %v\n", *t.Name, t.IsPrimitive)
	type_data := TemplateNapiHandlers{
		NativeType: *t.Name,
		NameSpace:  t.NameSpace,
		IsPointer:  t.IsPointer,
		NeedsCast:  new(string),
	}

	var needs_cast *string = new(string)

	isEnum, enumName := g.IsTemplateEnum(t)
	if isEnum {
		type_data.NapiType = NumberEnum
		type_data.JSType = *enumName
		*needs_cast = "double"
		type_data.NeedsCast = needs_cast
		return type_data
	}

	if t.IsPrimitive && t.IsPointer {
		if *t.Name == "char" {
			type_data.NapiType = String
			type_data.JSType = js_string_type
			type_data.STLType = std_string
			return type_data
		} else if *t.Name == "char16_t" {
			type_data.NapiType = String
			type_data.JSType = js_string_type
			type_data.STLType = std_u16string
			return type_data
		}
		helpers := GetTypedArrayHandlers(*t.Name)
		type_data.NapiType = helpers.NapiType
		type_data.TypedArrayInfo = helpers
		type_data.NativeType = helpers.NativeType
		type_data.JSType = helpers.NapiType.JSTypeString()
		if helpers.NeedsCast {
			*type_data.NeedsCast = helpers.NapiType.GetArrayType().String()
		}
		return type_data
	}

	if t.IsPrimitive {
		type_data.NapiType = Number
		type_data.JSType = js_number_type
		switch *t.Name {
		case "float", "int32_t", "signed", "int8_t", "int16_t", "short", "int", "unsigned", "uint8_t", "uint16_t", "unsigned short", "unsigned int", "uint32_t":
			*type_data.NeedsCast = "double"
			type_data.NativeType = *t.Name
		case "int64_t":
			{
				type_data.JSType = js_bigint_type
				type_data.NapiType = BigInt
				type_data.NativeType = *t.Name
			}
		case "long long", "signed long long":
			{
				type_data.NapiType = BigInt
				type_data.JSType = js_bigint_type
				type_data.NeedsCast = t.Name
				type_data.NativeType = *t.Name
			}
		case "uint64_t":
			{
				type_data.NapiType = BigInt
				type_data.JSType = js_bigint_type
				type_data.NativeType = *t.Name
			}
		case "unsigned long long", "size_t", "uintptr_t":
			{
				type_data.NapiType = BigInt
				type_data.JSType = js_bigint_type
				type_data.NeedsCast = t.Name
				type_data.NativeType = *t.Name
			}
		case "bool":
			{
				type_data.NapiType = Boolean
				type_data.JSType = js_boolean_type
			}
		case "std::string", "string":
			{
				type_data.NapiType = String
				type_data.JSType = js_string_type
			}
		case "std::u16string":
			{
				type_data.NapiType = String
				type_data.JSType = js_string_type
			}
		}
		return type_data
	}

	// handles struct/class types
	class_info := g.getClass(*t.Name)
	if class_info != nil {
		type_data.NativeType = fmt.Sprintf("%s::%s", *class_info.NameSpace, *t.Name)
	} else {
		type_data.NativeType = *t.Name
	}
	type_data.NapiType = External
	type_data.JSType = *t.Name
	return type_data

}

func (t *CPPType) ParseReturnData(g *PackageGenerator) GenReturnData {
	type_data := GenReturnData{
		IsPointer: t.IsPointer,
		IsConst:   t.IsConst,
		RefDecl:   t.RefDecl,
		RawType:   t,
		NeedsCast: new(string),
	}

	usedType := t

	if t.MappedType != nil {
		usedType = t.MappedType
	}

	isEnum, enumName := g.IsArgEnum(usedType)
	if isEnum {
		type_data.NapiType = NumberEnum
		type_data.JSType = *enumName
		needs_cast := "double"
		type_data.NeedsCast = &needs_cast
		return type_data
	}

	if usedType.IsPrimitive && usedType.IsPointer {
		if usedType.Name == "char" {
			type_data.NapiType = String
			type_data.JSType = js_string_type
			type_data.STLType = std_string
			return type_data
		} else if usedType.Name == "char16_t" {
			type_data.NapiType = String
			type_data.JSType = js_string_type
			type_data.STLType = std_u16string
			return type_data
		}
		helpers := GetTypedArrayHandlers(usedType.Name)
		type_data.NapiType = helpers.NapiType
		type_data.TypedArrayInfo = helpers
		type_data.NativeType = helpers.NativeType
		type_data.JSType = helpers.NapiType.TypedArrayType()
		*type_data.NeedsCast = helpers.NapiType.GetArrayType().String()
		return type_data
	}

	if usedType.Template != nil {
		// fmt.Printf("%s is template type!\n", *usedType.Template.Name)
		switch *usedType.Template.Name {
		case "vector":
			{
				type_data.STLType = std_vector
				helpers := GetTypedArrayHandlers(*usedType.Template.Args[0].Name)
				if helpers != nil {
					type_data.NapiType = helpers.NapiType
					type_data.NativeType = helpers.NativeType
					type_data.TypedArrayInfo = helpers
					type_data.JSType = helpers.NapiType.TypedArrayType()
					*type_data.NeedsCast = helpers.NapiType.GetArrayType().String()
				} else {
					type_data.NapiType = Array
					if *usedType.Template.Args[0].Name == "pair" || *usedType.Template.Args[0].Name == "std::pair" {
						type_data.NapiType = PairArray
						type_data.JSType = fmt.Sprintf("Array<%s>", g.PairToJsType(usedType.Template.Args[0]))
					} else {
						helpers := g.GetTypeHelpers(*usedType.Template.Args[0].Name)
						type_data.JSType = fmt.Sprintf("Array<%s>", helpers.JSType)
					}
				}
				return type_data
			}
		case "array":
			{
				// fmt.Println(*usedType.Template.Args[1].Name)
				type_data.STLType = std_array
				helpers := GetTypedArrayHandlers(*usedType.Template.Args[0].Name)
				if helpers != nil {
					type_data.NapiType = helpers.NapiType
					type_data.NativeType = helpers.NativeType
					type_data.TypedArrayInfo = helpers
					type_data.JSType = helpers.NapiType.TypedArrayType()
					*type_data.NeedsCast = helpers.NapiType.GetArrayType().String()
				} else {
					type_data.NapiType = Array
					if *usedType.Template.Args[0].Name == "pair" || *usedType.Template.Args[0].Name == "std::pair" {
						type_data.NapiType = PairArray
						type_data.JSType = fmt.Sprintf("Array<%s>", g.PairToJsType(usedType.Template.Args[0]))
					} else {
						helpers := g.GetTypeHelpers(*usedType.Template.Args[0].Name)
						type_data.JSType = fmt.Sprintf("Array<%s>", helpers.JSType)
					}
				}
				return type_data
			}
		case "pair":
			{
				type_data.STLType = std_pair
				type_data.NapiType = Pair
				type_data.JSType = g.PairToJsType(usedType.Template)
				return type_data
			}
			// TODO: need to handle addtl types (e.g. `std::map`)
		}
	}

	if usedType.IsPrimitive {
		type_data.NapiType = Number
		type_data.JSType = js_number_type
		switch usedType.Name {
		case "float", "int32_t", "signed", "int8_t", "int16_t", "short", "int", "unsigned", "uint8_t", "uint16_t", "unsigned short", "unsigned int", "uint32_t":
			*type_data.NeedsCast = "double"
			type_data.NativeType = "double"
		case "int64_t":
			{
				type_data.JSType = js_bigint_type
				type_data.NapiType = BigInt
			}
		case "long long", "signed long long":
			{
				type_data.NapiType = BigInt
				type_data.JSType = js_bigint_type
				*type_data.NeedsCast = "int64_t"
				type_data.NativeType = usedType.Name
			}
		case "uint64_t":
			{
				type_data.NapiType = BigInt
				type_data.JSType = js_bigint_type
			}
		case "unsigned long long", "size_t", "uintptr_t":
			{
				type_data.NapiType = BigInt
				type_data.JSType = js_bigint_type
				type_data.NativeType = usedType.Name
				*type_data.NeedsCast = "uint64_t"
			}
		case "bool":
			{
				type_data.NapiType = Boolean
				type_data.JSType = js_boolean_type
			}
		case "std::string", "string":
			{
				type_data.NapiType = String
				type_data.JSType = js_string_type
			}
		case "std::u16string":
			{
				type_data.NapiType = String
				type_data.JSType = js_string_type
			}
		}
		return type_data
	}

	// handles struct/class types
	class_info := g.getClass(usedType.Name)
	if class_info != nil {
		type_data.NativeType = fmt.Sprintf("%s::%s", *class_info.NameSpace, usedType.Name)
	} else {
		type_data.NativeType = usedType.Name
	}
	type_data.NapiType = External
	type_data.JSType = usedType.Name
	return type_data
}

func (a *CPPArg) ParseDefaultValue() string {
	val := *a.DefaultValue.Val
	val = strings.ReplaceAll(val, "{", "[")
	val = strings.ReplaceAll(val, "}", "]")
	return val
}

func (a *CPPArg) ParseArgData(g *PackageGenerator, name string, idx int) GenArgData {
	arg_data := GenArgData{
		Name:      name,
		IsPointer: a.Type.IsPointer,
		IsConst:   a.Type.IsConst,
		RefDecl:   a.Type.RefDecl,
		Idx:       idx,
		RawType:   a.Type,
	}

	if a.DefaultValue != nil {
		default_value := a.ParseDefaultValue()
		arg_data.DefaultValue = &default_value
	}

	usedType := a.Type

	if a.Type.MappedType != nil {
		// fmt.Println("mapped typed: ", *a.Name)
		if !a.Type.MappedType.IsPrimitive {
			arg_data.NeedsConstructor = a.Type.MappedType.OGType
		}
		usedType = a.Type.MappedType
	}

	isEnum, enumName := g.IsArgEnum(usedType)
	if isEnum {
		arg_data.NapiType = NumberEnum
		arg_data.NapiGetter = Int32Value
		arg_data.NeedsCast = true
		arg_data.NativeType = *enumName
		arg_data.JSType = usedType.Name
		if a.DefaultValue != nil {
			if g.IsEnumVal(usedType, *a.DefaultValue.Val) {
				default_value := fmt.Sprintf("%s.%s", usedType.Name, *a.DefaultValue.Val)
				arg_data.DefaultValue = &default_value
			} else {
				arg_data.DefaultValue = nil
			}
		}
		return arg_data
	}

	// handles `TypedArray` check
	if usedType.IsPrimitive && usedType.IsPointer {
		if usedType.Name == "char" {
			arg_data.NapiType = String
			arg_data.NapiGetter = Utf8Value
			arg_data.JSType = js_string_type
			arg_data.STLType = std_string
			return arg_data
		} else if usedType.Name == "char16_t" {
			arg_data.NapiType = String
			arg_data.NapiGetter = Utf16Value
			arg_data.JSType = js_string_type
			arg_data.STLType = std_string
			return arg_data
		}
		helpers := GetTypedArrayHandlers(usedType.Name)
		arg_data.NapiType = helpers.NapiType
		arg_data.NapiGetter = Data
		arg_data.TypedArrayInfo = helpers
		arg_data.NativeType = helpers.NativeType
		arg_data.JSType = helpers.NapiType.JSTypeString()
		arg_data.NeedsCast = helpers.NeedsCast
		return arg_data
	}

	// handles Temlate `std::*` types
	if usedType.Template != nil {
		switch *usedType.Template.Name {
		case "vector":
			{
				arg_data.STLType = std_vector
				helpers := GetTypedArrayHandlers(*usedType.Template.Args[0].Name)
				if helpers != nil {
					arg_data.NapiType = helpers.NapiType
					arg_data.NapiGetter = Data
					arg_data.NativeType = helpers.NativeType
					arg_data.TypedArrayInfo = helpers
					arg_data.JSType = helpers.NapiType.JSTypeString()
					arg_data.NeedsCast = helpers.NeedsCast
				} else {
					arg_data.NapiType = Array
					if *usedType.Template.Args[0].Name == "pair" || *usedType.Template.Args[0].Name == "std::pair" {
						arg_data.NapiType = PairArray
						arg_data.JSType = fmt.Sprintf("Array<%s>", g.PairToJsType(usedType.Template.Args[0]))
					} else {
						helpers := g.GetTypeHelpers(*usedType.Template.Args[0].Name)
						arg_data.NapiGetter = Data
						arg_data.JSType = fmt.Sprintf("Array<%s>", helpers.JSType)
					}
				}
				return arg_data
			}
		case "array":
			{
				arg_data.STLType = std_array
				helpers := GetTypedArrayHandlers(*usedType.Template.Args[0].Name)
				if helpers != nil {
					arg_data.NapiType = helpers.NapiType
					arg_data.NapiGetter = Data
					arg_data.NativeType = helpers.NativeType
					arg_data.STLType = std_array
					arg_data.TypedArrayInfo = helpers
					arg_data.JSType = helpers.NapiType.JSTypeString()
					arg_data.NeedsCast = helpers.NeedsCast
				} else {
					arg_data.NapiType = Array
					if *usedType.Template.Args[0].Name == "pair" || *usedType.Template.Args[0].Name == "std::pair" {
						arg_data.NapiType = PairArray
						arg_data.JSType = fmt.Sprintf("Array<%s>", g.PairToJsType(usedType.Template.Args[0]))
					} else {
						helpers := g.GetTypeHelpers(*usedType.Template.Args[0].Name)
						arg_data.NapiGetter = Data
						arg_data.JSType = fmt.Sprintf("Array<%s>", helpers.JSType)
					}
				}
				return arg_data
			}
		case "pair":
			{
				arg_data.STLType = std_pair
				arg_data.NapiType = Pair
				arg_data.JSType = g.PairToJsType(usedType.Template)
				return arg_data
			}
			// TODO: need to handle addtl types (e.g. `std::map`)
		}
	}

	// TODO: fix handling of `char` and other string types
	// handles primitive types (and `std::string`/`string`)
	if usedType.IsPrimitive {
		arg_data.NapiType = Number
		arg_data.JSType = js_number_type
		switch usedType.Name {
		case "float":
			arg_data.NapiGetter = FloatValue
		case "double":
			arg_data.NapiGetter = DoubleValue
		case "int32_t":
			arg_data.NapiGetter = Int32Value
		case "signed", "int8_t", "int16_t", "short", "int":
			{
				arg_data.NapiGetter = Int32Value
				arg_data.NeedsCast = true
				arg_data.NativeType = usedType.Name
			}
		case "unsigned", "uint8_t", "uint16_t", "unsigned short", "unsigned int":
			{
				arg_data.NapiGetter = Uint32Value
				arg_data.NeedsCast = true
				arg_data.NativeType = usedType.Name
			}
		case "uint32_t":
			arg_data.NapiGetter = Uint32Value
		case "int64_t":
			{
				arg_data.NapiType = BigInt
				arg_data.JSType = js_bigint_type
				arg_data.NapiGetter = Int64Value
			}
		case "long long":
			{
				arg_data.NapiType = BigInt
				arg_data.JSType = js_bigint_type
				arg_data.NapiGetter = Int64Value
				arg_data.NeedsCast = true
				arg_data.NativeType = usedType.Name
			}
		case "uint64_t":
			{
				arg_data.NapiType = BigInt
				arg_data.JSType = js_bigint_type
				arg_data.NapiGetter = Uint64Value
			}
		case "unsigned long long", "size_t", "uintptr_t":
			{
				arg_data.NapiType = BigInt
				arg_data.JSType = js_bigint_type
				arg_data.NapiGetter = Uint64Value
				arg_data.NeedsCast = true
				arg_data.NativeType = usedType.Name
			}
		case "bool":
			{
				arg_data.NapiType = Boolean
				arg_data.JSType = js_boolean_type
				arg_data.NapiGetter = Value
			}
		case "std::time_t", "time_t":
			{
				arg_data.NapiType = Date
				arg_data.JSType = js_date_type
				arg_data.NapiGetter = ValueOf
			}
		case "std::string", "string":
			{
				arg_data.NapiType = String
				arg_data.JSType = js_string_type
				arg_data.NapiGetter = Utf8Value
			}
		case "std::u16string":
			{
				arg_data.NapiType = String
				arg_data.JSType = js_string_type
				arg_data.NapiGetter = Utf16Value
			}
		}
		return arg_data
	}

	// handles struct/class types
	class_info := g.getClass(usedType.Name)
	if class_info != nil {
		arg_data.NativeType = fmt.Sprintf("%s::%s", *class_info.NameSpace, usedType.Name)
	} else {
		arg_data.NativeType = usedType.Name
	}
	arg_data.NapiType = External
	arg_data.JSType = usedType.Name
	arg_data.NapiGetter = Data
	return arg_data
}

func (g *PackageGenerator) WriteFnCall(sb *strings.Builder, fn_prefix string, method_name string, return_type *CPPType, args *[]GenArgData, is_void bool, is_class_method ...bool) *string {
	is_class_field := false
	if len(is_class_method) > 0 {
		is_class_field = is_class_method[0]
	}
	var gen_result_name *string = new(string)
	if !is_void {
		*gen_result_name = GetPrefixedVarName("res", "value")
	}
	isTransform, isGrouped, transform := g.conf.IsReturnTransform(method_name)

	if (isTransform && !isGrouped && !strings.Contains(*transform, method_name)) || !isTransform || (isTransform && isGrouped) {
		g.writeIndent(sb, 1)
		if !is_void {
			sb.WriteString("auto")
			if return_type.IsPointer {
				sb.WriteByte('*')
			}
			sb.WriteString(fmt.Sprintf(" %s = ", *gen_result_name))
		}

		sb.WriteString(fmt.Sprintf("%s%s(", fn_prefix, method_name))
		arg_count := len(*args)
		for i, a := range *args {
			if is_class_field && i == 0 {
				continue
			}
			if ((i > 0 && !is_class_field) || i > 1) && i < arg_count {
				sb.WriteString(", ")
			}
			switch a.NapiType {
			case String:
				sb.WriteString(a.Name)
			case Boolean, Number, NumberEnum, BigInt, Date, Array, PairArray, Pair:
				if a.IsPointer {
					sb.WriteByte('&')
				}
				if a.NeedsConstructor != nil {
					sb.WriteString(fmt.Sprintf("%s(", *a.NeedsConstructor))
				}
				sb.WriteString(a.Name)
				if a.NeedsConstructor != nil {
					sb.WriteByte(')')
				}
			case TypedArray, Int8Array, Uint8Array, Int16Array, Uint16Array, Int32Array, Uint32Array, Float32Array, Float64Array, BigInt64Array, BigUint64Array:
				if a.NeedsConstructor != nil {
					sb.WriteString(fmt.Sprintf("%s(", *a.NeedsConstructor))
				}
				sb.WriteString(a.Name)
				if a.NeedsConstructor != nil {
					sb.WriteByte(')')
				}
			case External:
				if !a.IsPointer {
					sb.WriteByte('*')
				}
				sb.WriteString(a.Name)
			}
		}
		sb.WriteString(");\n")
	}

	if isTransform {
		*transform = strings.ReplaceAll(*transform, "/return/", *gen_result_name)
		for i, arg := range *args {
			*transform = strings.ReplaceAll(*transform, fmt.Sprintf("/arg_%d/", i), arg.Name)
		}
		sb.WriteString(*transform)
	}

	return gen_result_name
}
