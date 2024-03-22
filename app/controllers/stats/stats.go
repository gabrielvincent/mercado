package stats

import (
	"mercado/app/models/expense"
	v "mercado/app/views/stats"
	"mercado/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func getFirstDayOfCurrentMonth() time.Time {
	now := time.Now()
	return time.Date(
		now.Year(),
		now.Month(),
		1,
		0,
		0,
		0,
		0,
		now.UTC().Location(),
	)
}

func getLastDayOfCurrentMonth() time.Time {
	now := time.Now()
	return time.Date(
		now.Year(),
		now.Month()+1,
		0,
		23,
		59,
		59,
		0,
		now.UTC().Location(),
	)
}

func Index(context echo.Context) error {

	startDate := getFirstDayOfCurrentMonth()
	endDate := getLastDayOfCurrentMonth()
	expenses, err := expense.GetExpenses(startDate, endDate)

	if err != nil {
		return err
	}

	return utils.Render(context, http.StatusOK, v.Index(expenses))
}
