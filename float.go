package go_utils

import (
	"fmt"
	"strconv"
)

// 两位小数
func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}
