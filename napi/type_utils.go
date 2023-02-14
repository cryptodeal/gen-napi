package napi

import (
	"fmt"
	"strings"
)

type ArgChecks struct {
	NapiChecker  *string
	NapiGetter   *string
	CastTo       *string
	ErrorDetails string
	IsArray      *string
	NapiType     string
	JSType       string
}

// `node-addon-api` types
const napi_number_type = "Number"
const napi_bigint_type = "BigInt"
const napi_boolean_type = "Boolean"
const napi_string_type = "String"
const napi_external_type = "External"

// const napi_object_type = "Object"
// const napi_typedarray_type = "TypedArray"
// const napi_array_type = "Array"

// js types
const js_number_type = "number"
const js_bigint_type = "bigint"
const js_boolean_type = "boolean"
const js_string_type = "string"

// const js_array_type = "Array"

// `node-addon-api` type checkers
var numCheck = "IsNumber"
var bigIntCheck = "IsBigInt"
var typedArrayCheck = "IsTypedArray"
var arrayCheck = "IsArray"
var boolCheck = "IsBoolean"
var stringCheck = "IsString"
var externalCheck = "IsExternal"

// var objCheck = "IsObject"

// `node-addon-api` value getters
var floatGetter = "FloatValue"
var doubleGetter = "DoubleValue"
var i32Getter = "Int32Value"
var u32Getter = "Uint32Value"
var i64Getter = "Int64Value"
var defaultGetter = "Value"
var utf8Getter = "Utf8Value"
var utf16Getter = "Utf16Value"
var ptrGetter = "Data"

func GetArrayName(a *CPPArg) *string {
	arrName := fmt.Sprintf("_tmp_parsed_%s", *a.Name)
	return &arrName
}

func (g *PackageGenerator) PairToJsType(t *TemplateType) string {
	p1, _ := g.CPPTypeToTS(*t.Args[0].Name, false)
	p2, _ := g.CPPTypeToTS(*t.Args[1].Name, false)
	return fmt.Sprintf("[%s, %s]", p1, p2)
}

func (g *PackageGenerator) WriteArgCheck(sb *strings.Builder, checks *ArgChecks, name string, i int, a *CPPArg, is_void bool) {
	if checks != nil {
		g.writeArgTypeChecker(sb, name, *checks.NapiChecker, i, checks.ErrorDetails, 1, nil, a, is_void)
		if checks.IsArray != nil {
			tsType, isObject := g.CPPTypeToTS(*a.Type.Template.Args[0].Name, false)
			if isObject {
				g.writeArgTypeChecker(sb, name, "IsExternal", i, fmt.Sprintf("native `%s` (typeof `Napi::External<%s::%s>`)", tsType, *g.NameSpace, tsType), 2, checks.IsArray, nil, is_void)
			} else {
				g.writeArgTypeChecker(sb, name, fmt.Sprintf("Is%s", g.casers.upper.String(tsType[0:1])+tsType[1:]), i, fmt.Sprintf("typeof `%s`", tsType), 2, checks.IsArray, nil, is_void)
			}
		}
	}
}

func (g *PackageGenerator) WritePrimitiveGetter(sb *strings.Builder, t *CPPType, name string, checks *ArgChecks, i int) {
	g.writeIndent(sb, 1)
	sb.WriteString(fmt.Sprintf("%s %s = ", t.Name, name))
	// only cast when necessary
	if checks.CastTo != nil {
		sb.WriteString(fmt.Sprintf("static_cast<%s>(", *checks.CastTo))
	}
	sb.WriteString(fmt.Sprintf("info[%d].As<Napi::%s>().%s()", i, checks.NapiType, *checks.NapiGetter))
	if checks.CastTo != nil {
		sb.WriteByte(')')
	}
	sb.WriteString(";\n")
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
		checks.NapiChecker = &numCheck
		checks.NapiType = napi_number_type
		checks.JSType = js_number_type
		checks.NapiGetter = &i32Getter
		checks.CastTo = enumName
		return checks
	}

	if g.isClass(t) {
		checks.NapiChecker = &externalCheck
		checks.NapiType = napi_external_type
		checks.JSType = t
		checks.NapiGetter = &ptrGetter
		return checks
	}

	/* TODO: handle
	if t.IsPrimitive && t.IsPointer {
		return checks
	}
	*/

	checks.NapiType = napi_number_type
	checks.JSType = js_number_type
	checks.ErrorDetails = "typeof `number`"
	switch t {
	case "float":
		{
			checks.NapiChecker = &numCheck
			checks.NapiGetter = &floatGetter
		}
	case "double":
		{
			checks.NapiChecker = &numCheck
			checks.NapiGetter = &doubleGetter
		}
	case "int32_t":
		{
			checks.NapiChecker = &numCheck
			checks.NapiGetter = &i32Getter
		}
	case "long long", "signed", "int8_t", "int16_t", "short", "int":
		{
			checks.NapiChecker = &numCheck
			checks.NapiGetter = &i32Getter
			checks.CastTo = &t
		}
	case "unsigned long long", "unsigned", "uint8_t", "uint16_t", "unsigned short", "unsigned int":
		{
			checks.NapiChecker = &numCheck
			checks.NapiGetter = &u32Getter
			checks.CastTo = &t
		}
	case "uint32_t":
		{
			checks.NapiChecker = &numCheck
			checks.NapiGetter = &u32Getter
		}
	case "int64_t":
		{
			checks.NapiChecker = &bigIntCheck
			checks.NapiType = napi_bigint_type
			checks.JSType = js_bigint_type
			checks.ErrorDetails = "typeof `bigint`"
			checks.NapiGetter = &i64Getter
		}
	case "uint64_t", "size_t", "uintptr_t":
		{
			checks.NapiChecker = &bigIntCheck
			checks.NapiType = napi_bigint_type
			checks.JSType = js_bigint_type
			checks.ErrorDetails = "typeof `bigint`"
			checks.NapiGetter = &i64Getter
			checks.CastTo = &t
		}
	case "bool":
		{
			checks.NapiChecker = &boolCheck
			checks.NapiType = napi_boolean_type
			checks.ErrorDetails = "typeof `boolean`"
			checks.JSType = js_boolean_type
			checks.NapiGetter = &defaultGetter
		}
	case "std::string", "string":
		{
			checks.NapiChecker = &stringCheck
			checks.ErrorDetails = "typeof `string`"
			checks.NapiType = napi_string_type
			checks.JSType = js_string_type
			checks.NapiGetter = &utf8Getter
		}
	case "std::u16string":
		{
			checks.NapiChecker = &stringCheck
			checks.NapiType = napi_string_type
			checks.ErrorDetails = "typeof `string`"
			checks.JSType = js_string_type
			checks.NapiGetter = &utf16Getter
		}
	}

	return checks
}

func (t *CPPType) GetTypeHandlers(g *PackageGenerator) *ArgChecks {
	checks := &ArgChecks{}

	isEnum, enumName := g.IsArgEnum(t)
	// enums are passed as `number`
	if isEnum {
		checks.NapiChecker = &numCheck
		checks.NapiType = napi_number_type
		checks.ErrorDetails = "typeof `number`"
		checks.JSType = js_number_type
		checks.NapiGetter = &i32Getter
		checks.CastTo = enumName
		return checks
	}

	if g.isClass(t.Name) {
		checks.NapiChecker = &externalCheck
		checks.NapiType = napi_external_type
		checks.ErrorDetails = fmt.Sprintf("native `%s` (typeof `Napi::External<%s::%s>`)", t.Name, *g.NameSpace, t.Name)
		checks.JSType = t.Name
		checks.NapiGetter = &ptrGetter
		return checks
	}

	// handles `TypedArray` check
	if t.IsPrimitive && t.IsPointer {
		jsType, arrayType, needsCast, _ := PrimitivePtrToTS(t.Name)
		checks.NapiChecker = &typedArrayCheck
		checks.JSType = jsType
		checks.NapiGetter = &ptrGetter
		checks.NapiType = fmt.Sprintf("TypedArrayOf<%s>", arrayType)
		checks.ErrorDetails = fmt.Sprintf("instanceof `%s`", jsType)
		checks.CastTo = needsCast
		return checks
	}

	// handles Temlate `std::*` types
	if t.Template != nil {
		switch *t.Template.Name {
		case "vector":
			{
				checks.NapiChecker = &arrayCheck
				if len(t.Template.Args[0].Args) > 0 {
					if *t.Template.Name == "pair" {
						checks.JSType = g.PairToJsType(t.Template.Args[0])
					}
				} else {
					checks.JSType = fmt.Sprintf("Array<%s>", *t.Template.Args[0].Name)
				}
				checks.ErrorDetails = fmt.Sprintf("`%s`", checks.JSType)
				return checks
			}
		case "pair":
			{
				checks.NapiChecker = &arrayCheck
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
		checks.NapiType = napi_number_type
		checks.JSType = js_number_type
		checks.ErrorDetails = "typeof `number`"
		switch t.Name {
		case "float":
			{
				checks.NapiChecker = &numCheck
				checks.NapiGetter = &floatGetter
			}
		case "double":
			{
				checks.NapiChecker = &numCheck
				checks.NapiGetter = &doubleGetter
			}
		case "int32_t":
			{
				checks.NapiChecker = &numCheck
				checks.NapiGetter = &i32Getter
			}
		case "long long", "signed", "int8_t", "int16_t", "short", "int":
			{
				checks.NapiChecker = &numCheck
				checks.NapiGetter = &i32Getter
				checks.CastTo = &t.Name
			}
		case "unsigned long long", "unsigned", "uint8_t", "uint16_t", "unsigned short", "unsigned int":
			{
				checks.NapiChecker = &numCheck
				checks.NapiGetter = &u32Getter
				checks.CastTo = &t.Name
			}
		case "uint32_t":
			{
				checks.NapiChecker = &numCheck
				checks.NapiGetter = &u32Getter
			}
		case "int64_t":
			{
				checks.NapiChecker = &bigIntCheck
				checks.NapiType = napi_bigint_type
				checks.JSType = js_bigint_type
				checks.ErrorDetails = "typeof `bigint`"
				checks.NapiGetter = &i64Getter
			}
		case "uint64_t", "size_t", "uintptr_t":
			{
				checks.NapiChecker = &bigIntCheck
				checks.NapiType = napi_bigint_type
				checks.JSType = js_bigint_type
				checks.ErrorDetails = "typeof `bigint`"
				checks.NapiGetter = &i64Getter
				checks.CastTo = &t.Name
			}
		case "bool":
			{
				checks.NapiChecker = &boolCheck
				checks.NapiType = napi_boolean_type
				checks.ErrorDetails = "typeof `boolean`"
				checks.JSType = js_boolean_type
				checks.NapiGetter = &defaultGetter
			}
		case "std::string", "string":
			{
				checks.NapiChecker = &stringCheck
				checks.ErrorDetails = "typeof `string`"
				checks.NapiType = napi_string_type
				checks.JSType = js_string_type
				checks.NapiGetter = &utf8Getter
			}
		case "std::u16string":
			{
				checks.NapiChecker = &stringCheck
				checks.NapiType = napi_string_type
				checks.ErrorDetails = "typeof `string`"
				checks.JSType = js_string_type
				checks.NapiGetter = &utf16Getter
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
