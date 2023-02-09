package napi

import (
	"fmt"
	"strings"
)

func WriteTSMap(sb *strings.Builder, keyType string, valType string) {
	sb.WriteString(fmt.Sprintf("Map<%s, %s[]>", keyType, valType))
}

func WriteTSMultiMap(sb *strings.Builder, keyType string, valType string) {
	WriteTSMap(sb, keyType, valType+"[]")
}
