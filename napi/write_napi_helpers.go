package napi

import (
	"fmt"
	"strings"
)

func write_mem_adjustment(sb *strings.Builder, get_bytes string, track_mem_usage *string) {
	bytes_var_name := GetPrefixedVarName("bytes", "used")
	sb.WriteString(fmt.Sprintf("auto %s = %s;\n", bytes_var_name, get_bytes))
	if track_mem_usage != nil {
		sb.WriteString(fmt.Sprintf("%s += %s;\n", *track_mem_usage, bytes_var_name))
	}
	sb.WriteString(fmt.Sprintf("Napi::MemoryManagement::AdjustExternalMemory(env, %s);\n", bytes_var_name))
}

func write_external_finalizer(sb *strings.Builder, logic string, hint *string) {
	sb.WriteString("[](Napi::Env env, void* data")
	if hint != nil {
		sb.WriteString(fmt.Sprintf(", %s", *hint))
	}
	sb.WriteString(") {\n")
	sb.WriteString(logic)
	sb.WriteString("}\n")
}
