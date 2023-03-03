package napi

import (
	"fmt"
	"strings"
)

const gen_var_prefix = "_gen_tmp"

func GetPrefixedVarName(name ...string) string {
	res_name := gen_var_prefix
	if len(name) > 0 && strings.Contains(name[0], gen_var_prefix) {
		res_name = name[0]
		name = name[1:]
	}
	for _, n := range name {
		res_name = strings.Join([]string{res_name, n}, "_")
	}
	return res_name
}

type NapiType int

func (t NapiType) IsTypedArray() bool {
	switch t {
	case Uint8Array, Int8Array, Uint16Array, Int16Array, Uint32Array, Int32Array, Float32Array, Float64Array, BigInt64Array, BigUint64Array:
		return true
	default:
		return false
	}
}

func (t NapiType) String() *string {
	var out *string = new(string)
	switch t {
	case String:
		*out = "String"
	case Number, NumberEnum:
		*out = "Number"
	case Date:
		*out = "Date"
	case BigInt:
		*out = "BigInt"
	case Boolean:
		*out = "Boolean"
	case External:
		*out = "External"
	case Object:
		*out = "Object"
	// TODO: probably should remove `pair` as a type longer term
	case Pair, Array:
		*out = "Array"
	case Buffer:
		*out = "Buffer"
	case ArrayBuffer:
		*out = "ArrayBuffer"
	case TypedArray:
		*out = "TypedArray"
	case Int8Array, Uint8Array, Int16Array, Uint16Array, Int32Array, Uint32Array, BigInt64Array, BigUint64Array, Float32Array, Float64Array:
		*out = fmt.Sprintf("TypedArrayOf<%s>", t.GetArrayType().String())
	case DataView:
		*out = "DataView"
	case Function:
		*out = "Function"
	case Null:
		*out = "Null"
	case Promise:
		*out = "Promise"
	case Undefined:
		*out = "Undefined"
	default:
		panic("Unknown NapiType")
	}
	return out
}

func (t NapiType) TypedArrayType() string {
	switch t {
	case Uint8Array:
		return "Uint8Array"
	case Int8Array:
		return "Int8Array"
	case Uint16Array:
		return "Uint16Array"
	case Int16Array:
		return "Int16Array"
	case Uint32Array:
		return "Uint32Array"
	case Int32Array:
		return "Int32Array"
	case Float32Array:
		return "Float32Array"
	case Float64Array:
		return "Float64Array"
	case BigInt64Array:
		return "BigInt64Array"
	case BigUint64Array:
		return "BigUint64Array"
	default:
		panic("Unknown NapiType")
	}
}

func (t NapiType) JSTypeString() string {
	switch t {
	case String:
		return "string"
	case Number:
		return "number"
	case Date:
		return "Date"
	case BigInt:
		return "bigint"
	case Boolean:
		return "boolean"
	case External:
		return "never"
	case Buffer:
		return "Buffer"
	case ArrayBuffer:
		return "ArrayBuffer"
	case Uint8Array:
		return "number[] | Uint8Array"
	case Int8Array:
		return "number[] | Int8Array"
	case Uint16Array:
		return "number[] | Uint16Array"
	case Int16Array:
		return "number[] | Int16Array"
	case Uint32Array:
		return "number[] | Uint32Array"
	case Int32Array:
		return "number[] | Int32Array"
	case Float32Array:
		return "number[] | Float32Array"
	case Float64Array:
		return "number[] | Float64Array"
	case BigInt64Array:
		return "Array<number | bigint> | BigInt64Array"
	case BigUint64Array:
		return "Array<number | bigint> | BigUint64Array"
	case DataView:
		return "DataView"
	case Function:
		return "Function"
	case Null:
		return "null"
	case Undefined:
		return "undefined"
	default:
		return "any"
	}
}

func (t NapiType) GetTypeChecker() NapiTypeChecker {
	switch t {
	case String:
		return IsString
	case NumberEnum, Number:
		return IsNumber
	case Date:
		return IsDate
	case BigInt:
		return IsBigInt
	case Boolean:
		return IsBoolean
	case External:
		return IsExternal
	case Object:
		return IsObject
	case PairArray, Pair, Array:
		return IsArray
	case Buffer:
		return IsBuffer
	case ArrayBuffer:
		return IsArrayBuffer
	case TypedArray, Int8Array, Uint8Array, Int16Array, Uint16Array, Int32Array, Uint32Array, BigInt64Array, BigUint64Array, Float32Array, Float64Array:
		return IsTypedArray
	case DataView:
		return IsDataView
	case Function:
		return IsFunction
	case Null:
		return IsNull
	case Promise:
		return IsPromise
	case Symbol:
		return IsSymbol
	case Undefined:
		return IsUndefined
	}
	panic("Unknown NapiType")
}

func (t NapiType) GetArrayType() *NapiArrayType {
	var arrayType = new(NapiArrayType)
	switch t {
	case Int8Array:
		*arrayType = Int8
	case Uint8Array:
		*arrayType = Uint8
	case Int16Array:
		*arrayType = Int16
	case Uint16Array:
		*arrayType = Uint16
	case Int32Array:
		*arrayType = Int32
	case Uint32Array:
		*arrayType = Uint32
	case BigInt64Array:
		*arrayType = BigInt64
	case BigUint64Array:
		*arrayType = BigUint64
	case Float32Array:
		*arrayType = Float32
	case Float64Array:
		*arrayType = Float64
	}
	return arrayType
}

type NapiArrayType int

const (
	Float32 NapiArrayType = iota
	Float64
	Uint8
	Int8
	Uint16
	Int16
	Uint32
	Int32
	BigInt64
	BigUint64
)

func (t NapiArrayType) NapiString() string {
	switch t {
	case Float32:
		return "float32"
	case Float64:
		return "float64"
	case Uint8:
		return "uint8"
	case Int8:
		return "int8"
	case Uint16:
		return "uint16"
	case Int16:
		return "int16"
	case Uint32:
		return "uint32"
	case Int32:
		return "int32"
	case BigInt64:
		return "bigint64"
	case BigUint64:
		return "biguint64"
	}
	panic("Unknown NapiArrayType")
}

func (t NapiArrayType) String() string {
	switch t {
	case Float32:
		return "float"
	case Float64:
		return "double"
	case Uint8:
		return "uint8_t"
	case Int8:
		return "int8_t"
	case Uint16:
		return "uint16_t"
	case Int16:
		return "int16_t"
	case Uint32:
		return "uint32_t"
	case Int32:
		return "int32_t"
	case BigInt64:
		return "int64_t"
	case BigUint64:
		return "uint64_t"
	}
	panic("Unknown NapiArrayType")
}

const (
	String NapiType = iota
	Number
	Date
	BigInt
	Boolean
	External
	Object
	Array
	Buffer
	ArrayBuffer
	TypedArray
	Int8Array
	Uint8Array
	Int16Array
	Uint16Array
	Int32Array
	Uint32Array
	BigInt64Array
	BigUint64Array
	Float32Array
	Float64Array
	DataView
	Function
	Null
	Promise
	Symbol
	Undefined
	// library includes dedicated handlers for the following
	BufferString // handles `char*``
	NumberEnum
	PairArray
	Pair
)

type STLType string

const (
	std_vector = "std::vector"
	std_array  = "std::array"
	std_pair   = "std::pair"
	// string types
	std_string    = "std::string"
	std_u16string = "u16string"
	// std_list   = "std::list"
	// std_set = "std::set"
	// std_map = "std::map"
)

type NapiTypeChecker int

const (
	IsString NapiTypeChecker = iota
	IsNumber
	IsDate
	IsBigInt
	IsBoolean
	IsExternal
	IsObject
	IsArray
	IsBuffer
	IsArrayBuffer
	IsTypedArray
	IsDataView
	IsFunction
	IsNull
	IsPromise
	IsSymbol
	IsUndefined
)

func (t NapiTypeChecker) String() *string {
	var out *string = new(string)
	switch t {
	case IsString:
		*out = "IsString"
	case IsNumber:
		*out = "IsNumber"
	case IsDate:
		*out = "IsDate"
	case IsBigInt:
		*out = "IsBigInt"
	case IsBoolean:
		*out = "IsBoolean"
	case IsExternal:
		*out = "IsExternal"
	case IsObject:
		*out = "IsObject"
	case IsArray:
		*out = "IsArray"
	case IsBuffer:
		*out = "IsBuffer"
	case IsArrayBuffer:
		*out = "IsArrayBuffer"
	case IsTypedArray:
		*out = "IsTypedArray"
	case IsDataView:
		*out = "IsDataView"
	case IsFunction:
		*out = "IsFunction"
	case IsNull:
		*out = "IsNull"
	case IsPromise:
		*out = "IsPromise"
	case IsUndefined:
		*out = "IsUndefined"
	default:
		panic("Unknown NapiTypeChecker")
	}
	return out
}

type NapiTypeGetter int

const (
	FloatValue NapiTypeGetter = iota
	DoubleValue
	Int32Value
	Uint32Value
	Int64Value
	Uint64Value
	Value
	ValueOf
	Utf8Value
	Utf16Value
	Data
	// used to access `DataView`
	ArrayBufferAccessor
)

func (t NapiType) FindFlexibleGetter() NapiTypeGetter {
	switch t {
	case String:
		return Utf8Value
	case Number:
		return Int64Value
	case Date:
		return ValueOf
	case BigInt:
		return Int64Value
	case Boolean:
		return Value
	case External:
		return Data
	case Buffer:
		return Data
	case ArrayBuffer:
		return Data
	case NumberEnum:
		return Int32Value
	case TypedArray, Int8Array, Uint8Array, Int16Array, Uint16Array, Int32Array, Uint32Array, BigInt64Array, BigUint64Array, Float32Array, Float64Array:
		return Data
	case DataView:
		return ArrayBufferAccessor
	}
	panic("Unknown NapiType")
}

func (t NapiTypeGetter) String() *string {
	var out *string = new(string)
	switch t {
	// number getters
	case FloatValue:
		*out = "FloatValue"
	case DoubleValue:
		*out = "DoubleValue"
	case Int32Value:
		*out = "Int32Value"
	case Uint32Value:
		*out = "Uint32Value"
	// shared number/bigint getters
	case Int64Value:
		*out = "Int64Value"
	case Uint64Value:
		*out = "Uint64Value"
	// string getters
	case Utf8Value:
		*out = "Utf8Value"
	case Utf16Value:
		*out = "Utf16Value"
	// pointer getter
	case Data:
		*out = "Data"
	// boolean getter
	case Value:
		*out = "Value"
	// Date getter (ms since epoch)
	case ValueOf:
		*out = "ValueOf"
	case ArrayBufferAccessor:
		*out = "ArrayBuffer"
	default:
		panic("Unknown NapiTypeChecker")
	}
	return out
}

func TypedArrayToNapiType(array_type string) NapiType {
	switch array_type {
	case "int8_t":
		return Int8Array
	case "uint8_t":
		return Uint8Array
	case "int16_t":
		return Int16Array
	case "uint16_t":
		return Uint16Array
	case "int32_t":
		return Int32Array
	case "uint32_t":
		return Uint32Array
	case "int64_t":
		return BigInt64Array
	case "uint64_t":
		return BigUint64Array
	case "float":
		return Float32Array
	case "double":
		return Float64Array
	default:
		panic("Unknown TypedArray type")
	}
}
