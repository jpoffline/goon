package goon

import (
	"strconv"
	"strings"
)

// FloatToString converts a float to a string
func FloatToString(fl float64) string {
	// to convert a float number to a string
	precision := 4
	return strconv.FormatFloat(fl, 'f', precision, 64)
}

// Pad a string to be a certain length
func Pad(prefix string, reqd int) string {
	return strings.Repeat("0", 5-len(prefix)) + prefix
}
