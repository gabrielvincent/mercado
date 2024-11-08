package utils

import (
	"fmt"
	"strconv"
)

func FormatCurrency(value int) string {
	whole := value / 100
	decimal := value % 100
	decimalStr := strconv.Itoa(decimal)
	paddedDecimal := fmt.Sprintf("%02s", decimalStr)

	return "R$" + strconv.Itoa(whole) + "," + paddedDecimal
}
