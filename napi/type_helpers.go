package napi

import (
	"fmt"
	"strings"
)

func stripNameSpace(v string) string {
	if strings.Contains(v, "::") {
		lastIdx := strings.LastIndex(v, "::")
		v = v[lastIdx+2:]
	}
	return v
}

func (g *PackageGenerator) isClass(argType string) bool {
	argType = stripNameSpace(argType)
	_, ok := g.ParsedData.Classes[argType]
	return ok
}

func (g *PackageGenerator) getClass(argType string) *CPPClass {
	argType = stripNameSpace(argType)
	if v, ok := g.ParsedData.Classes[argType]; ok {
		return v
	}
	return nil
}

func PrimitivePtrToTS(t string) (string, string, *string, string) {
	jsTypeEquivalent := ""
	var needsCast *string
	napi_short_type := ""
	arrayType := ""
	switch t {
	case "float":
		arrayType = "float"
		napi_short_type = "float32"
		jsTypeEquivalent = "Float32Array"
	case "double":
		arrayType = "double"
		napi_short_type = "float64"
		jsTypeEquivalent = "Float64Array"
	case "uint8_t":
		arrayType = "uint8_t"
		napi_short_type = "uint8"
		jsTypeEquivalent = "Uint8Array"
	case "bool":
		arrayType = "uint8_t"
		napi_short_type = "uint8"
		jsTypeEquivalent = "Uint8Array"
		needsCast = &t
	case "int8_t":
		arrayType = "int8_t"
		napi_short_type = "int8"
		jsTypeEquivalent = "Int8Array"
	case "uint16_t":
		arrayType = "uint16_t"
		napi_short_type = "uint16"
		jsTypeEquivalent = "Uint16Array"
	case "unsigned short":
		arrayType = "uint16_t"
		napi_short_type = "uint16"
		jsTypeEquivalent = "Uint16Array"
		needsCast = &t
	case "int16_t":
		arrayType = "int16_t"
		napi_short_type = "int16"
		jsTypeEquivalent = "Int16Array"
	case "short", "signed short":
		arrayType = "int16_t"
		napi_short_type = "int16"
		jsTypeEquivalent = "Int16Array"
		needsCast = &t
	case "uint32_t":
		arrayType = "uint32_t"
		napi_short_type = "uint32"
		jsTypeEquivalent = "Uint32Array"
	case "unsigned int":
		arrayType = "uint32_t"
		napi_short_type = "uint32"
		jsTypeEquivalent = "Uint32Array"
		needsCast = &t
	case "int32_t":
		arrayType = "int32_t"
		napi_short_type = "int32"
		jsTypeEquivalent = "Int32Array"
	case "int", "signed int":
		arrayType = "int32_t"
		napi_short_type = "int32"
		jsTypeEquivalent = "Int32Array"
		needsCast = &t
	case "int64_t":
		arrayType = "int64_t"
		napi_short_type = "bigint64"
		jsTypeEquivalent = "BigInt64Array"
	case "long long", "long long int", "signed long long", "signed long long int":
		arrayType = "int64_t"
		napi_short_type = "bigint64"
		jsTypeEquivalent = "BigInt64Array"
		needsCast = &t
	case "uint64_t":
		napi_short_type = "biguint64"
		arrayType = "uint64_t"
		jsTypeEquivalent = "BigUint64Array"
	case "unsigned long long", "unsigned long long int", "size_t":
		arrayType = "uint64_t"
		napi_short_type = "biguint64"
		jsTypeEquivalent = "BigUint64Array"
		needsCast = &t
	}
	return jsTypeEquivalent, arrayType, needsCast, napi_short_type
}

func IsTypeNumber(t string) bool {
	switch strings.TrimSpace(t) {
	case "short", "int", "int8_t", "uint8_t", "int16_t", "uint16_t", "int32_t", "uint32_t", "long", "float", "float_t", "double", "double_t", "long double":
		return true
	default:
		return false
	}
}

func IsArgTemplate(t *CPPArg) bool {
	return t.Type != nil && t.Type.Template != nil
}

func IsTypeBigInt(t string) bool {
	switch strings.TrimSpace(t) {
	case "long long", "size_t", "int64_t", "uint64_t":
		return true
	default:
		return false
	}
}

func IsTypeString(t string) bool {
	switch strings.TrimSpace(t) {
	case "string", "std::string", "char", "wchar_t", "char16_t", "char32_t":
		return true
	default:
		return false
	}
}

func (g *PackageGenerator) CPPTypeToTS(t string, isPointer bool) (string, bool) {
	isEnum, _ := g.IsTypeEnum(t)
	if isEnum {
		return t, false
	}
	if isPointer {
		jsArrayType, _, _, _ := PrimitivePtrToTS(t)
		if jsArrayType != "" {
			return jsArrayType, false
		}
	}
	switch t {
	case "int", "int8_t", "uint8_t", "signed", "unsigned", "short", "long", "long int", "signed long", "signed long int", "unsigned long", "unsigned long int", "long double", "short int", "signed short", "unsigned_short", "signed int", "unsigned int", "unsigned short int", "signed short int", "uint16_t", "uint32_t", "int16_t", "int32_t", "float", "double":
		return "number", false
	case "int64_t", "uint64_t", "long long", "long long int", "signed long long", "signed long long int", "unsigned long long", "unsigned long long int", "size_t":
		return "bigint", false
	case "string", "std::string", "char", "wchar_t", "char16_t", "char32_t":
		return "string", false
	case "bool":
		return "boolean", false
	default:
		return t, true
	}
}

func (g *PackageGenerator) IsArgEnum(t *CPPType) (bool, *string) {
	if t == nil {
		return false, nil
	}
	for _, e := range g.ParsedData.Enums {
		fullName := fmt.Sprintf("%s::%s", *e.NameSpace, e.Name)
		if strings.EqualFold(fullName, t.Name) || strings.EqualFold(e.Name, t.Name) {
			return true, &fullName
		}
	}
	return false, nil
}

func (g *PackageGenerator) IsEnumVal(t *CPPType, v string) bool {
	for _, e := range g.ParsedData.Enums {
		fullName := fmt.Sprintf("%s::%s", *e.NameSpace, e.Name)
		if strings.EqualFold(fullName, t.Name) || strings.EqualFold(e.Name, t.Name) {
			for _, val := range e.Values {
				if strings.EqualFold(*val.Name, v) {
					return true
				}
			}
		}
	}
	return false
}

func (g *PackageGenerator) IsTemplateEnum(t *TemplateType) (bool, *string) {
	if t == nil {
		return false, nil
	}
	for _, e := range g.ParsedData.Enums {
		fullName := fmt.Sprintf("%s::%s", *e.NameSpace, e.Name)
		if strings.EqualFold(fullName, *t.Name) || strings.EqualFold(e.Name, *t.Name) {
			return true, &fullName
		}
	}
	return false, nil
}

func (g *PackageGenerator) IsTypeEnum(t string) (bool, *string) {
	for _, e := range g.ParsedData.Enums {
		fullName := fmt.Sprintf("%s::%s", *e.NameSpace, e.Name)
		if strings.EqualFold(fullName, t) || strings.EqualFold(e.Name, t) {
			return true, &fullName
		}
	}
	return false, nil
}

func TypeIsTypedArray(t string) (*string, bool) {
	numArray := fmt.Sprintf("number[] | %s", t)
	bigIntArray := fmt.Sprintf("Array<number | bigint> | %s", t)

	switch t {
	case "Float32Array", "Float64Array", "Int8Array", "Int16Array", "Int32Array", "Uint8Array", "Uint16Array", "Uint32Array":
		return &numArray, false
	case "BigInt64Array", "BigUint64Array":
		return &bigIntArray, true
	default:
		return nil, false
	}
}
