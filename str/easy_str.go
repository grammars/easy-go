package str

import (
	"strings"
)

func IsBlank(text string) bool {
	return strings.TrimSpace(text) == ""
}
