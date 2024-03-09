package home

import (
	"strconv"
	"time"
)

type Expense struct {
	ID           int
	Value        int
	GroceryStore string
	Date         time.Time
}

var PT_MONTHS = []string{
	"Janeiro",
	"Fevereiro",
	"Março",
	"Abril",
	"Maio",
	"Junho",
	"Julho",
	"Agosto",
	"Setembro",
	"Outubro",
	"Novembro",
	"Dezembro",
}

func FormatDate(date time.Time) string {

	formatted := strconv.Itoa(
		date.Day(),
	) + " de " + PT_MONTHS[date.Month()-1] + ", " + date.Format(
		"15h:04m",
	)

	return formatted
}
