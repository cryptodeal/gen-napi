package napi

import (
	"fmt"
	"strings"
)

type GetterArg struct {
	NameSuffix string
	Type       string
	Value      string
	IsPointer  bool
	ErrorCheck *string
}

type ArgChecks struct {
	IsMapped     bool
	GetterArgs   *[]GetterArg
	NapiChecker  NapiTypeChecker
	NapiGetter   NapiTypeGetter
	CastTo       *string
	ErrorDetails string
	IsArray      *string
	NapiType     NapiType
	JSType       string
	NameSpace    *string
}

// js types
const js_number_type = "number"
const js_bigint_type = "bigint"
const js_boolean_type = "boolean"
const js_string_type = "string"
const js_date_type = "Date"

func GetArrayName(a *CPPArg) *string {
	arrName := GetPrefixedVarName("parsed", *a.Name)
	return &arrName
}

func (g *PackageGenerator) PairToJsType(t *TemplateType) string {
	p1, _ := g.CPPTypeToTS(*t.Args[0].Name, false)
	p2, _ := g.CPPTypeToTS(*t.Args[1].Name, false)
	return fmt.Sprintf("[%s, %s]", p1, p2)
}

func (g *PackageGenerator) WriteArgCheck(sb *strings.Builder, checks *ArgChecks, name string, i int, a *CPPArg, is_void bool) {
	if checks != nil {
		g.WriteArgTypeCheck(sb, name, *checks.NapiChecker.String(), i, checks.ErrorDetails, nil, a, is_void)
		if checks.IsArray != nil {
			tsType, isObject := g.CPPTypeToTS(*a.Type.Template.Args[0].Name, false)
			if isObject {
				g.WriteArgTypeCheck(sb, name, "IsExternal", i, fmt.Sprintf("native `%s` (typeof `Napi::External<%s::%s>`)", tsType, *g.NameSpace, tsType), checks.IsArray, nil, is_void)
			} else {
				g.WriteArgTypeCheck(sb, name, fmt.Sprintf("Is%s", g.casers.upper.String(tsType[0:1])+tsType[1:]), i, fmt.Sprintf("typeof `%s`", tsType), checks.IsArray, nil, is_void)
			}
		}
	}
}

func (g *PackageGenerator) WritePrimitiveGetter(sb *strings.Builder, t *CPPType, name string, checks *ArgChecks, i int) {
	has_args := checks.GetterArgs != nil
	if has_args {
		for _, arg := range *checks.GetterArgs {
			g.writeIndent(sb, 1)
			sb.WriteString(fmt.Sprintf("%s _%s_%s = %s;\n", arg.Type, name, arg.NameSuffix, arg.Value))
		}
	}
	g.writeIndent(sb, 1)
	sb.WriteString(fmt.Sprintf("%s %s = ", t.Name, name))
	// only cast when necessary
	if checks.CastTo != nil {
		sb.WriteString(fmt.Sprintf("static_cast<%s>(", *checks.CastTo))
	}
	sb.WriteString(fmt.Sprintf("info[%d].As<Napi::%s>().%s(", i, *checks.NapiType.String(), *checks.NapiGetter.String()))
	if has_args {
		arg_count := len(*checks.GetterArgs)
		for i, arg := range *checks.GetterArgs {
			if i > 0 && i < arg_count {
				sb.WriteString(", ")
			}
			if arg.IsPointer {
				sb.WriteByte('&')
			}
			sb.WriteString(fmt.Sprintf("_%s_%s", name, arg.NameSuffix))
		}
	}
	sb.WriteByte(')')
	if checks.CastTo != nil {
		sb.WriteByte(')')
	}
	sb.WriteString(";\n")
	if has_args {
		for _, arg := range *checks.GetterArgs {
			if arg.ErrorCheck == nil {
				continue
			}
			g.WriteErrorHandler(sb, fmt.Sprintf("!_%s_%s", name, arg.NameSuffix), fmt.Sprintf("failed to parse `%s`; %s", name, *arg.ErrorCheck), 1)
		}
	}
}

func (t *CPPType) GetScopedType() string {
	var full_type string
	if t.NameSpace != nil {
		full_type = fmt.Sprintf("%s::", *t.NameSpace)
	}
	return full_type + t.Name
}

func (g *PackageGenerator) WriteArgGetter(sb *strings.Builder, checks *ArgChecks, name string, a *CPPArg, i int) {
	isArgTransform, argTransformVal := g.conf.IsArgTransform(name, *a.Name)
	if isArgTransform && !strings.Contains(*argTransformVal, "/arg_") {
		g.writeIndent(sb, 1)
		sb.WriteString(strings.ReplaceAll(*argTransformVal, "/arg/", fmt.Sprintf("info[%d]", i)))
		return
	}

	if checks != nil {
		isEnum, enumName := g.IsArgEnum(a.Type)
		// hasDefault := a != nil && a.DefaultValue != nil
		if isEnum {
			g.writeIndent(sb, 1)
			sb.WriteString(fmt.Sprintf("%s %s = static_cast<%s>(info[%d].As<Napi::Number>().Int32Value());\n\n", *enumName, *a.Name, *enumName, i))
			return
		}

		if checks.IsMapped && !isArgTransform {
			g.WritePrimitiveGetter(sb, a.Type.MappedType, *a.Name, checks, i)
			return
		}

		if a.Type.IsPrimitive && !isArgTransform {
			g.WritePrimitiveGetter(sb, a.Type, *a.Name, checks, i)
			return
		}

		if g.isClass(a.Type.Name) && !isArgTransform {
			classData := g.getClass(a.Type.Name)
			g.writeIndent(sb, 1)
			sb.WriteString(fmt.Sprintf("%s::%s* %s = UnExternalize<%s::%s>(info[%d]);\n", *classData.NameSpace, stripNameSpace(a.Type.Name), *a.Name, *classData.NameSpace, stripNameSpace(a.Type.Name), i))
			return
		}
	}
}

func (g *PackageGenerator) GetTypeHelpers(t string) *ArgChecks {
	checks := &ArgChecks{}

	isEnum, enumName := g.IsTypeEnum(t)
	if isEnum {
		checks.NapiType = Number
		checks.NapiChecker = checks.NapiType.GetTypeChecker()
		checks.JSType = js_number_type
		checks.NapiGetter = Int32Value
		checks.CastTo = enumName
		return checks
	}

	if g.isClass(t) {
		checks.NapiType = External
		checks.NapiChecker = checks.NapiType.GetTypeChecker()
		checks.JSType = t
		checks.NapiGetter = Data
		return checks
	}

	/* TODO: handle
	if t.IsPrimitive && t.IsPointer {
		return checks
	}
	*/

	checks.NapiType = Number
	checks.JSType = js_number_type
	checks.NapiChecker = checks.NapiType.GetTypeChecker()
	checks.ErrorDetails = "typeof `number`"
	switch t {
	case "float":
		{
			checks.NapiChecker = checks.NapiType.GetTypeChecker()
			checks.NapiGetter = FloatValue
		}
	case "double":
		{
			checks.NapiGetter = DoubleValue
		}
	case "int32_t":
		{
			checks.NapiGetter = Int32Value
		}
	case "signed", "int8_t", "int16_t", "short", "int":
		{
			checks.NapiGetter = Int32Value
			checks.CastTo = &t
		}
	case "unsigned", "uint8_t", "uint16_t", "unsigned short", "unsigned int":
		{
			checks.NapiGetter = Uint32Value
			checks.CastTo = &t
		}
	case "uint32_t":
		{
			checks.NapiGetter = Uint32Value
		}
	case "int64_t":
		{
			checks.NapiType = BigInt
			checks.JSType = js_bigint_type
			checks.ErrorDetails = "typeof `bigint`"
			error_msg := "failed to convert `bigint` to `int64_t` losslessly"
			checks.GetterArgs = &[]GetterArg{
				{
					NameSuffix: "lossless",
					Value:      "true",
					Type:       "bool",
					IsPointer:  true,
					ErrorCheck: &error_msg,
				},
			}
			checks.NapiGetter = Int64Value
		}
	case "long long":
		{
			checks.NapiType = BigInt
			checks.NapiChecker = checks.NapiType.GetTypeChecker()
			checks.JSType = js_bigint_type
			checks.ErrorDetails = "typeof `bigint`"
			checks.CastTo = &t
			error_msg := "failed to convert `bigint` to `int64_t` losslessly"
			checks.GetterArgs = &[]GetterArg{
				{
					NameSuffix: "lossless",
					Value:      "true",
					Type:       "bool",
					IsPointer:  true,
					ErrorCheck: &error_msg,
				},
			}
			checks.NapiGetter = Int64Value
		}
	case "uint64_t":
		{
			checks.NapiType = BigInt
			checks.NapiChecker = checks.NapiType.GetTypeChecker()
			checks.JSType = js_bigint_type
			checks.ErrorDetails = "typeof `bigint`"
			checks.NapiGetter = Uint64Value
			error_msg := "failed to convert `bigint` to `uint64_t` losslessly"
			checks.GetterArgs = &[]GetterArg{
				{
					NameSuffix: "lossless",
					Value:      "true",
					Type:       "bool",
					IsPointer:  true,
					ErrorCheck: &error_msg,
				},
			}
		}
	case "unsigned long long", "size_t", "uintptr_t":
		{
			checks.NapiType = BigInt
			checks.NapiChecker = checks.NapiType.GetTypeChecker()
			checks.JSType = js_bigint_type
			checks.ErrorDetails = "typeof `bigint`"
			checks.NapiGetter = Uint64Value
			error_msg := "failed to convert `bigint` to `uint64_t` losslessly"
			checks.GetterArgs = &[]GetterArg{
				{
					NameSuffix: "lossless",
					Value:      "true",
					Type:       "bool",
					IsPointer:  true,
					ErrorCheck: &error_msg,
				},
			}
			checks.CastTo = &t
		}
	case "bool":
		{
			checks.NapiType = Boolean
			checks.NapiChecker = checks.NapiType.GetTypeChecker()
			checks.ErrorDetails = "typeof `boolean`"
			checks.JSType = js_boolean_type
			checks.NapiGetter = Value
		}
	case "std::string", "string":
		{
			checks.NapiType = String
			checks.NapiChecker = checks.NapiType.GetTypeChecker()
			checks.ErrorDetails = "typeof `string`"
			checks.JSType = js_string_type
			checks.NapiGetter = Utf8Value
		}
	case "std::u16string":
		{
			checks.NapiType = String
			checks.NapiChecker = checks.NapiType.GetTypeChecker()
			checks.ErrorDetails = "typeof `string`"
			checks.JSType = js_string_type
			checks.NapiGetter = Utf16Value
		}
	}

	return checks
}

func (t *CPPType) GetTypeHandlers(g *PackageGenerator, isMapped ...bool) *ArgChecks {
	is_mapped_type := false
	if len(isMapped) > 0 {
		is_mapped_type = isMapped[0]
	}
	if t.MappedType != nil {
		return t.MappedType.GetTypeHandlers(g, true)
	}
	checks := &ArgChecks{
		IsMapped: is_mapped_type,
	}

	isEnum, enumName := g.IsArgEnum(t)
	// enums are passed as `number`
	if isEnum {
		checks.NapiType = Number
		checks.NapiChecker = checks.NapiType.GetTypeChecker()
		checks.ErrorDetails = "typeof `number`"
		checks.JSType = js_number_type
		checks.NapiGetter = Int32Value
		checks.CastTo = enumName
		return checks
	}

	if g.isClass(t.Name) {
		checks.NapiType = External
		checks.NapiChecker = checks.NapiType.GetTypeChecker()
		checks.ErrorDetails = fmt.Sprintf("native `%s` (typeof `Napi::External<%s::%s>`)", t.Name, *g.NameSpace, t.Name)
		checks.JSType = t.Name
		checks.NapiGetter = Data
		return checks
	}

	// handles `TypedArray` check
	if t.IsPrimitive && t.IsPointer {
		jsType, arrayType, needsCast, _ := PrimitivePtrToTS(t.Name)
		checks.JSType = jsType
		checks.NapiGetter = Data
		checks.NapiType = TypedArrayToNapiType(arrayType)
		checks.NapiChecker = checks.NapiType.GetTypeChecker()
		checks.ErrorDetails = fmt.Sprintf("instanceof `%s`", jsType)
		checks.CastTo = needsCast
		return checks
	}

	// handles Temlate `std::*` types
	if t.Template != nil {
		switch *t.Template.Name {
		case "vector":
			{
				checks.NapiType = Array
				checks.NapiChecker = checks.NapiType.GetTypeChecker()
				if len(t.Template.Args[0].Args) > 0 {
					if *t.Template.Name == "pair" {
						checks.JSType = g.PairToJsType(t.Template.Args[0])
					}
				} else {
					helpers := g.GetTypeHelpers(*t.Template.Args[0].Name)
					checks.JSType = fmt.Sprintf("Array<%s>", helpers.JSType)
				}
				checks.ErrorDetails = fmt.Sprintf("`%s`", checks.JSType)
				return checks
			}
		case "pair":
			{
				checks.NapiType = Array
				checks.NapiChecker = checks.NapiType.GetTypeChecker()
				checks.ErrorDetails = fmt.Sprintf("`%s`", checks.JSType)
				checks.JSType = g.PairToJsType(t.Template)
				return checks
			}
			// TODO: need to handle addtl types (e.g. `std::map`)
		}
	}

	// TODO: fix handling of `char` and other string types
	// handles primitive types (and `std::string`/`string`)
	if t.IsPrimitive {
		checks.NapiType = Number
		checks.NapiChecker = checks.NapiType.GetTypeChecker()
		checks.JSType = js_number_type
		checks.ErrorDetails = "typeof `number`"
		switch t.Name {
		case "float":
			checks.NapiGetter = FloatValue
		case "double":
			checks.NapiGetter = DoubleValue
		case "int32_t":
			checks.NapiGetter = Int32Value
		case "signed", "int8_t", "int16_t", "short", "int":
			{
				checks.NapiGetter = Int32Value
				checks.CastTo = &t.Name
			}
		case "unsigned", "uint8_t", "uint16_t", "unsigned short", "unsigned int":
			{
				checks.NapiGetter = Uint32Value
				checks.CastTo = &t.Name
			}
		case "uint32_t":
			checks.NapiGetter = Uint32Value
		case "int64_t":
			{
				checks.NapiType = BigInt
				checks.NapiChecker = checks.NapiType.GetTypeChecker()
				checks.JSType = js_bigint_type
				checks.ErrorDetails = "typeof `bigint`"
				checks.NapiGetter = Int64Value
				error_msg := "failed to convert `bigint` to `int64_t` losslessly"
				checks.GetterArgs = &[]GetterArg{
					{
						NameSuffix: "lossless",
						Value:      "true",
						Type:       "bool",
						IsPointer:  true,
						ErrorCheck: &error_msg,
					},
				}
			}
		case "long long":
			{
				checks.NapiType = BigInt
				checks.NapiChecker = checks.NapiType.GetTypeChecker()
				checks.JSType = js_bigint_type
				checks.ErrorDetails = "typeof `bigint`"
				checks.NapiGetter = Int64Value
				error_msg := "failed to convert `bigint` to `int64_t` losslessly"
				checks.GetterArgs = &[]GetterArg{
					{
						NameSuffix: "lossless",
						Value:      "true",
						Type:       "bool",
						IsPointer:  true,
						ErrorCheck: &error_msg,
					},
				}
				checks.CastTo = &t.Name
			}
		case "uint64_t":
			{
				checks.NapiType = BigInt
				checks.NapiChecker = checks.NapiType.GetTypeChecker()
				checks.JSType = js_bigint_type
				checks.ErrorDetails = "typeof `bigint`"
				checks.NapiGetter = Uint64Value
				error_msg := "failed to convert `bigint` to `uint64_t` losslessly"
				checks.GetterArgs = &[]GetterArg{
					{
						NameSuffix: "lossless",
						Value:      "true",
						Type:       "bool",
						IsPointer:  true,
						ErrorCheck: &error_msg,
					},
				}
			}
		case "unsigned long long", "size_t", "uintptr_t":
			{
				checks.NapiType = BigInt
				checks.NapiChecker = checks.NapiType.GetTypeChecker()
				checks.JSType = js_bigint_type
				checks.ErrorDetails = "typeof `bigint`"
				checks.NapiGetter = Uint64Value
				error_msg := "failed to convert `bigint` to `uint64_t` losslessly"
				checks.GetterArgs = &[]GetterArg{
					{
						NameSuffix: "lossless",
						Value:      "true",
						Type:       "bool",
						IsPointer:  true,
						ErrorCheck: &error_msg,
					},
				}
				checks.CastTo = &t.Name
			}
		case "bool":
			{
				checks.NapiType = Boolean
				checks.NapiChecker = checks.NapiType.GetTypeChecker()
				checks.ErrorDetails = "typeof `boolean`"
				checks.JSType = js_boolean_type
				checks.NapiGetter = Value
			}
		case "std::string", "string":
			{
				checks.NapiType = String
				checks.NapiChecker = checks.NapiType.GetTypeChecker()
				checks.ErrorDetails = "typeof `string`"
				checks.JSType = js_string_type
				checks.NapiGetter = Utf8Value
			}
		case "std::u16string":
			{
				checks.NapiType = String
				checks.NapiChecker = checks.NapiType.GetTypeChecker()
				checks.ErrorDetails = "typeof `string`"
				checks.JSType = js_string_type
				checks.NapiGetter = Utf16Value
			}
		}
		return checks
	}

	// return nil; we don't know how to handle this type
	return nil
}

func (a *CPPArg) GetArgChecks(g *PackageGenerator) *ArgChecks {
	checks := a.Type.GetTypeHandlers(g)
	// handles Temlate `std::*` types
	if a.Type.Template != nil {
		switch *a.Type.Template.Name {
		case "vector":
			{
				checks.IsArray = GetArrayName(a)
			}
		case "pair":
			{
				checks.IsArray = GetArrayName(a)
			}
			// TODO: need to handle addtl types (e.g. `std::map`)
		}
	}
	return checks
}
