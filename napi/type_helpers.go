package napi

func isClass(argType string, classes map[string]*CPPClass) bool {
	_, ok := classes[argType]
	return ok
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
	case "unsigned char":
		arrayType = "uint8_t"
		napi_short_type = "uint8"
		jsTypeEquivalent = "Uint8Array"
		needsCast = &t
	case "int8_t":
		arrayType = "int8_t"
		napi_short_type = "int8"
		jsTypeEquivalent = "Int8Array"
	case "bool", "char", "signed char":
		arrayType = "int8_t"
		napi_short_type = "int8"
		jsTypeEquivalent = "Int8Array"
		needsCast = &t
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

func CPPTypeToTS(t string, isPointer bool) (string, bool) {
	if isPointer {
		jsArrayType, _, _, _ := PrimitivePtrToTS(t)
		if jsArrayType != "" {
			return jsArrayType, false
		}
	}
	switch t {
	case "int", "int8_t", "char", "uint8_t", "signed", "unsigned", "short", "long", "long int", "signed long", "signed long int", "unsigned long", "unsigned long int", "long double", "signed char", "unsigned char", "short int", "signed short", "unsigned_short", "signed int", "unsigned int", "unsigned short int", "signed short int", "uint16_t", "uint32_t", "int16_t", "int32_t", "float", "double":
		return "number", false
	case "int64_t", "uint64_t", "long long", "long long int", "signed long long", "signed long long int", "unsigned long long", "unsigned long long int", "size_t":
		return "bigint", false
	case "string", "std::string":
		return "string", false
	case "bool":
		return "boolean", false
	default:
		return t, true
	}
}
