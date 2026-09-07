package helpers

import (
	"strings"
)

func Required(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxLength(value string, max int) bool {
	return len([]rune(strings.TrimSpace(value))) <= max
}
