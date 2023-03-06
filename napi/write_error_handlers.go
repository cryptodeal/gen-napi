package napi

import (
	"fmt"
	"strings"
)

func (g *PackageGenerator) WriteErrorHandler(sb *strings.Builder, err_conditional string, err_message string, indent int, isVoid ...bool) {
	is_void := false
	if len(isVoid) > 0 {
		is_void = isVoid[0]
	}

	g.writeIndent(sb, indent)
	sb.WriteString(fmt.Sprintf("if (%s) {\n", err_conditional))
	g.writeIndent(sb, indent+1)
	sb.WriteString(fmt.Sprintf("Napi::TypeError::New(env, %q).ThrowAsJavaScriptException();\n", err_message))
	g.writeIndent(sb, indent+1)

	sb.WriteString("return")
	// return `undefined` if `void`
	if !is_void {
		g.writeIndent(sb, indent+1)
		init_napi_undefined(sb)
	} else {
		sb.WriteString(";\n")
	}
	g.writeIndent(sb, indent)
	sb.WriteString("}\n")
}

func (g *PackageGenerator) WriteArgCountCheck(sb *strings.Builder, name string, expected_arg_count int, optional int, is_void bool) {
	if expected_arg_count == 0 && optional == 0 {
		return
	}

	g.writeIndent(sb, 1)
	sb.WriteString("const auto _arg_count = info.Length();\n")

	var err_conditional, err_msg string
	if optional == 0 {
		err_conditional = fmt.Sprintf("_arg_count != %d", expected_arg_count)
		err_msg = fmt.Sprintf("`%s` expects exactly %d arg", name, expected_arg_count)
		if expected_arg_count > 1 {
			err_msg += "s"
		}
	} else {
		err_conditional = fmt.Sprintf("_arg_count < %d || _arg_count > %d", expected_arg_count, expected_arg_count+optional)
		err_msg = fmt.Sprintf("`%s` expects %d to %d args", name, expected_arg_count, expected_arg_count+optional)
	}

	g.WriteErrorHandler(sb, err_conditional, err_msg, 1, is_void)
}

func (g *PackageGenerator) WriteArgChecks(sb *strings.Builder, args *[]GenArgData, is_void bool, fn_name string) {
	g.WriteArgCountCheck(sb, fn_name, len(*args), 0, is_void)
	for _, a := range *args {
		switch a.NapiType {
		// handle string
		case String:
			g.WriteErrorHandler(sb, fmt.Sprintf("!info[%d].%s()", a.Idx, *a.NapiType.GetTypeChecker().String()), fmt.Sprintf("`%s` expects `%s` (args[%d]) to be typeof `string`", fn_name, a.Name, a.Idx), 1, is_void)
		// handle boolean
		case Boolean:
			g.WriteErrorHandler(sb, fmt.Sprintf("!info[%d].%s()", a.Idx, *a.NapiType.GetTypeChecker().String()), fmt.Sprintf("`%s` expects `%s` (args[%d]) to be typeof `boolean`", fn_name, a.Name, a.Idx), 1, is_void)
		// handle number
		case NumberEnum, Number:
			g.WriteErrorHandler(sb, fmt.Sprintf("!info[%d].%s()", a.Idx, *a.NapiType.GetTypeChecker().String()), fmt.Sprintf("`%s` expects `%s` (args[%d]) to be typeof `number`", fn_name, a.Name, a.Idx), 1, is_void)
		// handle `Array<T>` -> `std::vector<T>` or `std::array<T>` (where `T` is not primitive)
		case PairArray, Pair, Array:
			g.WriteErrorHandler(sb, fmt.Sprintf("!info[%d].%s()", a.Idx, *a.NapiType.GetTypeChecker().String()), fmt.Sprintf("`%s` expects `%s` (args[%d]) to be `Array`", fn_name, a.Name, a.Idx), 1, is_void)
		// handle bigint
		case BigInt:
			g.WriteErrorHandler(sb, fmt.Sprintf("!info[%d].%s()", a.Idx, *a.NapiType.GetTypeChecker().String()), fmt.Sprintf("`%s` expects `%s` (args[%d]) to be typeof `bigint`", fn_name, a.Name, a.Idx), 1, is_void)
		// handle date
		case Date:
			g.WriteErrorHandler(sb, fmt.Sprintf("!info[%d].%s()", a.Idx, *a.NapiType.GetTypeChecker().String()), fmt.Sprintf("`%s` expects `%s` (args[%d]) to be instanceof `Date`", fn_name, a.Name, a.Idx), 1, is_void)
		// handle typed arrays
		case TypedArray, Int8Array, Uint8Array, Int16Array, Uint16Array, Int32Array, Uint32Array, Float32Array, Float64Array, BigInt64Array, BigUint64Array:
			g.WriteErrorHandler(sb, fmt.Sprintf("!info[%d].%s()", a.Idx, *a.NapiType.GetTypeChecker().String()), fmt.Sprintf("`%s` expects `%s` (args[%d]) to be instanceof `%s`", fn_name, a.Name, a.Idx, a.NapiType.JSTypeString()), 1, is_void)
		// handle `Array<[T1, T2]>` -> `std::vector<std::pair<T1, T2>>`
		// handle returning pointer to JS
		case External:
			g.WriteErrorHandler(sb, fmt.Sprintf("!info[%d].%s()", a.Idx, *a.NapiType.GetTypeChecker().String()), fmt.Sprintf("`%s` expects `%s` (args[%d]) to be native `%s` (typeof `Napi::External<%s>)`", fn_name, a.Name, a.Idx, a.RawType.Name, a.RawType.GetFullType()), 1, is_void)
		}
	}
}
