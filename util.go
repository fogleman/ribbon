package ribbon

import "strconv"

func parseInt(x string) int {
	i, _ := strconv.ParseInt(x, 0, 0)
	return int(i)
}

func parseFloat(x string) float64 {
	f, _ := strconv.ParseFloat(x, 64)
	return f
}
