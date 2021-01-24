package main

import (
	"strconv"
	"strings"
)

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	precision := 4
	return strconv.FormatFloat(input_num, 'f', precision, 64)
}

func Pad(prefix string, reqd int) string {
	return strings.Repeat("0", 5-len(prefix)) + prefix
}
