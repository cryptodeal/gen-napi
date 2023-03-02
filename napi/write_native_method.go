package napi

import "strings"

func WriteReturnSignature(sb *strings.Builder, is_void bool) {
	if !is_void {
		sb.WriteString("Napi::Value")
	} else {
		sb.WriteString("void")
	}
}

func (g *PackageGenerator) WriteArgAccessor(sb *strings.Builder, dest_type string, dest_name string, arg_idx int, napi_type string) {
}
