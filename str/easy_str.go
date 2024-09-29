package str

import (
	"fmt"
	"strings"
)

func IsBlank(text string) bool {
	return strings.TrimSpace(text) == ""
}

func BytesToHex(data []byte) string {
	var hexStr string
	for _, b := range data {
		// 使用格式化字符串将每个字节转换为十六进制
		hexStr += fmt.Sprintf("%02x", b)
	}
	return hexStr
}
