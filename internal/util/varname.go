package util

import "strings"

func GetVariableName(val []byte) string {
	return strings.ReplaceAll(string(val), "$", "")
}
