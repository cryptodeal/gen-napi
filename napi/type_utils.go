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
const napi_object_type = "Object"
const napi_bigint_type = "BigInt"
const napi_typedarray_type = "TypedArray"
const napi_array_type = "Array"
const napi_boolean_type = "Boolean"
const napi_string_type = "String"
const napi_external_type = "External"

// js types
const js_number_type = "number"
const js_bigint_type = "bigint"
const js_boolean_type = "boolean"
const js_string_type = "string"
const js_array_type = "Array"

// `node-addon-api` type checkers
var numCheck = "IsNumber"
var objCheck = "IsObject"
var bigIntCheck = "IsBigInt"
var typedArrayCheck = "IsTypedArray"
var arrayCheck = "IsArray"
var boolCheck = "IsBoolean"
var stringCheck = "IsString"
var externalCheck = "IsExternal"

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
			tsType, isObject := g.CPPTypeToTS(*a.Template.Args[0].Name, false)
			if isObject {
				g.writeArgTypeChecker(sb, name, "IsExternal", i, fmt.Sprintf("native `%s` (typeof `Napi::External<%s::%s>`)", tsType, *g.NameSpace, tsType), 2, checks.IsArray, nil, is_void)
			} else {
				g.writeArgTypeChecker(sb, name, fmt.Sprintf("Is%s", g.casers.upper.String(tsType[0:1])+tsType[1:]), i, fmt.Sprintf("typeof `%s`", tsType), 2, checks.IsArray, nil, is_void)
			}
		}
	}
}

func (g *PackageGenerator) WritePrimitiveGetter(sb *strings.Builder, arg *CPPArg, checks *ArgChecks, i int) {
	g.writeIndent(sb, 1)
	sb.WriteString(fmt.Sprintf("%s %s = ", *arg.Type, *arg.Name))
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
		isEnum, enumName := g.IsArgEnum(a)
		// hasDefault := a != nil && a.DefaultValue != nil
		if isEnum {
			g.writeIndent(sb, 1)
			sb.WriteString(fmt.Sprintf("%s %s = static_cast<%s>(info[%d].As<Napi::Number>().Int32Value());\n\n", *enumName, *a.Name, *enumName, i))
			return
		}

		if a.IsPrimitive && !isArgTransform {
			g.WritePrimitiveGetter(sb, a, checks, i)
			return
		}

		if g.isClass(*a.Type) && !isArgTransform {
			classData := g.getClass(*a.Type)
			g.writeIndent(sb, 1)
			sb.WriteString(fmt.Sprintf("%s::%s* %s = UnExternalize<%s::%s>(info[%d]);\n", *classData.NameSpace, stripNameSpace(*a.Type), *a.Name, *classData.NameSpace, stripNameSpace(*a.Type), i))
			return
		}
	}
}

func (a *CPPArg) GetArgChecks(g *PackageGenerator) *ArgChecks {
	checks := &ArgChecks{}

	isEnum, _ := g.IsArgEnum(a)
	// enums are passed as `number`
	if isEnum {
		checks.NapiChecker = &numCheck
		checks.NapiType = napi_number_type
		checks.ErrorDetails = "typeof `number`"
		checks.JSType = js_number_type
		checks.NapiGetter = &i32Getter
		return checks
	}

	if g.isClass(*a.Type) {
		checks.NapiChecker = &externalCheck
		checks.NapiType = napi_external_type
		checks.ErrorDetails = fmt.Sprintf("native `%s` (typeof `Napi::External<%s::%s>`)", *a.Type, *g.NameSpace, *a.Type)
		checks.JSType = *a.Type
		checks.NapiGetter = &ptrGetter
		return checks
	}

	// handles `TypedArray` check
	if a.IsPrimitive && a.IsPointer {
		jsType, arrayType, needsCast, _ := PrimitivePtrToTS(*a.Type)
		checks.NapiChecker = &typedArrayCheck
		checks.JSType = jsType
		checks.NapiGetter = &ptrGetter
		checks.NapiType = fmt.Sprintf("TypedArrayOf<%s>", arrayType)
		checks.ErrorDetails = fmt.Sprintf("instanceof `%s`", jsType)
		checks.CastTo = needsCast
		return checks
	}

	// handles Temlate `std::*` types
	if a.Template != nil {
		switch *a.Template.Name {
		case "vector":
			{
				checks.IsArray = GetArrayName(a)
				checks.NapiChecker = &arrayCheck
				if len(a.Template.Args[0].Args) > 0 {
					if *a.Template.Name == "pair" {
						checks.JSType = g.PairToJsType(a.Template.Args[0])
					}
				} else {
					checks.JSType = fmt.Sprintf("Array<%s>", *a.Template.Args[0].Name)
				}
				checks.ErrorDetails = fmt.Sprintf("`%s`", checks.JSType)
				return checks
			}
		case "pair":
			{
				checks.IsArray = GetArrayName(a)
				checks.NapiChecker = &arrayCheck
				checks.ErrorDetails = fmt.Sprintf("`%s`", checks.JSType)
				checks.JSType = g.PairToJsType(a.Template)
				return checks
			}
			// TODO: need to handle addtl types (e.g. `std::map`)
		}
	}

	// TODO: fix handling of `char` and other string types
	// handles primitive types (and `std::string`/`string`)
	if a.IsPrimitive {
		checks.NapiType = napi_number_type
		checks.JSType = js_number_type
		checks.ErrorDetails = "typeof `number`"
		switch *a.Type {
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
				checks.CastTo = a.Type
			}
		case "unsigned long long", "unsigned", "uint8_t", "uint16_t", "unsigned short", "unsigned int":
			{
				checks.NapiChecker = &numCheck
				checks.NapiGetter = &u32Getter
				checks.CastTo = a.Type
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
				checks.CastTo = a.Type
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

	// return nil; no simple validation for this type
	return nil
}
