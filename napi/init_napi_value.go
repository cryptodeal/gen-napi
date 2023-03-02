package napi

import (
	"fmt"
	"strings"
)

// utils
func fmt_num_or_string(v interface{}) string {
	switch v.(type) {
	case int:
		return fmt.Sprintf("%d", v)
	case string:
		return fmt.Sprintf("%s", v)
	default:
		panic("invalid type for array length")
	}
}

// init primitive types (number, bool, bigint, string, undefined, null)

func init_napi_undefined(sb *strings.Builder) {
	sb.WriteString("env.Undefined();\n")
}

func init_napi_null(sb *strings.Builder) {
	sb.WriteString("env.Null();\n")
}

func init_napi_boolean(sb *strings.Builder, var_name string) {
	sb.WriteString(fmt.Sprintf("Napi::Boolean::New(env, %s);\n", var_name))
}

func init_napi_string(sb *strings.Builder, var_name string, length ...interface{}) {
	sb.WriteString("Napi::String::New(env, ")
	sb.WriteString(var_name)
	if len(length) > 0 {
		sb.WriteString(fmt.Sprintf(", %s", fmt_num_or_string(length[0])))
	}
	sb.WriteString(");\n")

}

func init_napi_number(sb *strings.Builder, var_name string, needs_cast *string) {
	sb.WriteString("Napi::Number::New(env, ")
	if needs_cast != nil && *needs_cast != "" {
		sb.WriteString(fmt.Sprintf("static_cast<%s>(", *needs_cast))
	}
	sb.WriteString(var_name)
	if needs_cast != nil && *needs_cast != "" {
		sb.WriteByte(')')
	}
	sb.WriteString(");\n")
}

func init_napi_bigint(sb *strings.Builder, var_name string, needs_cast *string) {
	sb.WriteString("Napi::BigInt::New(env, ")
	if needs_cast != nil {
		sb.WriteString(fmt.Sprintf("static_cast<%s>(", *needs_cast))
	}
	sb.WriteString(var_name)
	if needs_cast != nil {
		sb.WriteByte(')')
	}
	sb.WriteString(");\n")
}

func init_napi_date(sb *strings.Builder, var_name string, needs_cast *string) {
	// TODO: automatically coerce value to `double` if necessary
	sb.WriteString("Napi::Date::New(env, ")
	sb.WriteString(var_name)
	sb.WriteString(");\n")
}

// init more complex types

func init_napi_array(sb *strings.Builder, init_len ...interface{}) {
	var len_val interface{}
	if len(init_len) > 0 {
		len_val = init_len[0]
	}
	sb.WriteString("Napi::Array::New(env")
	if len_val != nil {
		sb.WriteString(", ")
	}
	sb.WriteString(fmt_num_or_string(len_val))
	sb.WriteString(");\n")
}

func init_napi_buffer(sb *strings.Builder, ptr_var_name string, byte_len interface{}, finalizer *string, finalizer_hint *string) {
	sb.WriteString(fmt.Sprintf("Napi::Buffer::New(env, %s, ", ptr_var_name))
	sb.WriteString(fmt_num_or_string(byte_len))
	if finalizer != nil {
		sb.WriteString(fmt.Sprintf(", %s", *finalizer))
	}
	if finalizer_hint != nil {
		sb.WriteString(fmt.Sprintf(", %s", *finalizer_hint))
	}
	sb.WriteString(");\n")
}

func init_napi_buffer_new_or_copy(sb *strings.Builder, ptr_var_name string, byte_len interface{}, finalizer *string, finalizer_hint *string) {
	sb.WriteString(fmt.Sprintf("Napi::Buffer::NewOrCopy(env, %s, ", ptr_var_name))
	sb.WriteString(fmt_num_or_string(byte_len))
	if finalizer != nil {
		sb.WriteString(fmt.Sprintf(", %s", *finalizer))
	}
	if finalizer_hint != nil {
		sb.WriteString(fmt.Sprintf(", %s", *finalizer_hint))
	}
	sb.WriteString(");\n")
}

func init_napi_buffer_copy(sb *strings.Builder, ptr_var_name string, byte_len interface{}, finalizer *string, finalizer_hint *string) {
	sb.WriteString(fmt.Sprintf("Napi::Buffer::Copy(env, %s, ", ptr_var_name))
	sb.WriteString(fmt_num_or_string(byte_len))
	if finalizer != nil {
		sb.WriteString(fmt.Sprintf(", %s", *finalizer))
	}
	if finalizer_hint != nil {
		sb.WriteString(fmt.Sprintf(", %s", *finalizer_hint))
	}
	sb.WriteString(");\n")
}

func init_napi_arraybuffer(sb *strings.Builder, ptr_var_name string, byte_len interface{}, finalizer *string, finalizer_hint *string) {
	sb.WriteString(fmt.Sprintf("Napi::ArrayBuffer::New(env, %s, ", ptr_var_name))
	sb.WriteString(fmt_num_or_string(byte_len))
	if finalizer != nil {
		sb.WriteString(fmt.Sprintf(", %s", *finalizer))
	}
	if finalizer_hint != nil {
		sb.WriteString(fmt.Sprintf(", %s", *finalizer_hint))
	}
	sb.WriteString(");\n")
}

func init_napi_typedarray(sb *strings.Builder, arraybuffer_var_name string, byte_offset interface{}, byte_len interface{}, array_type string) {
	sb.WriteString(fmt.Sprintf("Napi::%s::New(env, %s, ", array_type, arraybuffer_var_name))
	sb.WriteString(fmt_num_or_string(byte_offset))
	sb.WriteString(", ")
	sb.WriteString(fmt_num_or_string(byte_len))
	sb.WriteString(");\n")
}
