package napi

func isClass(argType string, classes map[string]*CPPClass) bool {
	_, ok := classes[argType]
	return ok
}

func CPPTypeToTS(t string) (string, bool) {
	switch t {
	case "int", "int8_t", "uint8_t", "signed", "unsigned", "short", "long", "long int", "size_t", "signed long", "signed long int", "unsigned long", "unsigned long int", "long long", "long long int", "signed long long", "signed long long int", "unsigned long long", "unsigned long long int", "long double", "signed char", "unsigned char", "short int", "signed short", "unsigned_short", "signed int", "unsigned int", "unsigned short int", "signed short int", "uint16_t", "uint32_t", "uint64_t", "int16_t", "int32_t", "int64_t", "float", "double":
		return "number", false
	case "string", "std::string", "char":
		return "string", false
	case "bool":
		return "boolean", false
	default:
		return t, true
	}
}
