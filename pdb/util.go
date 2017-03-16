package pdb

import (
	"strconv"
	"strings"
)

func parseString(x string) string {
	return strings.TrimSpace(x)
}

func parseInt(x string) int {
	x = parseString(x)
	i, _ := strconv.ParseInt(x, 0, 0)
	return int(i)
}

func parseFloat(x string) float64 {
	x = parseString(x)
	f, _ := strconv.ParseFloat(x, 64)
	return f
}
