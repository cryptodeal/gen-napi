package napi

import (
	"fmt"
	"strings"
)

func write_mem_adjustment(sb *strings.Builder, get_bytes string, track_mem_usage *string) {
	if track_mem_usage != nil {
		sb.WriteString(fmt.Sprintf("%s = ", *track_mem_usage))
	}
	sb.WriteString(fmt.Sprintf("Napi::MemoryManagement::AdjustExternalMemory(env, %s);\n", get_bytes))
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
