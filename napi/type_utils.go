package napi

import (
	"fmt"
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

func (g *PackageGenerator) PairToJsType(t *TemplateType) string {
	p1, _ := g.CPPTypeToTS(*t.Args[0].Name, false)
	p2, _ := g.CPPTypeToTS(*t.Args[1].Name, false)
	return fmt.Sprintf("[%s, %s]", p1, p2)
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
