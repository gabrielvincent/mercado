package stats

import (
	"errors"
	"fmt"
	expense "mercado/app/models/expense"
	stats "mercado/app/models/stats"
	v "mercado/app/views/stats"
	"mercado/utils"
	"net/http"
	"slices"
	"time"

	"github.com/labstack/echo/v4"
)

var NO_EXPENSES_ERR = errors.New("no expenses found")

func getFirstDayOfMonth(date time.Time) time.Time {
	return time.Date(
		date.Year(),
		date.Month(),
		1,
		0,
		0,
		0,
		0,
		time.UTC,
	)
}

func getLastDayOfMonth(date time.Time) time.Time {
	return time.Date(
		date.Year(),
		date.Month()+1,
		0,
		23,
		59,
		59,
		0,
		time.UTC,
	)
}

func getCurrentDayOfMonth() int {
	now := time.Now()
	return now.Day()
}

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

func getGroceryStoresRanking(
	expenses []expense.Expense,
) []expense.GroceryStoresRankingItem {

	groceryStoreToSessions := make(map[string]int)
	for _, expense := range expenses {
		if _, ok := groceryStoreToSessions[expense.GroceryStore]; !ok {
			groceryStoreToSessions[expense.GroceryStore] = 0
		}
		groceryStoreToSessions[expense.GroceryStore] += 1
	}

	rankingItems := make([]expense.GroceryStoresRankingItem, 0)
	for key, value := range groceryStoreToSessions {
		item := expense.GroceryStoresRankingItem{
			GroceryStore: key,
			Sessions:     value,
		}
		rankingItems = append(rankingItems, item)
	}

	slices.SortFunc(
		rankingItems,
		func(a, b expense.GroceryStoresRankingItem) int {
			if a.Sessions == b.Sessions {
				return 0
			}
			if a.Sessions > b.Sessions {
				return -1
			}
			return 1
		},
	)

	return rankingItems
}

func getPeriodComparison(
	targetDate time.Time,
	compareDate time.Time,
	durationInDays int,
) (*stats.PeriodComparison, error) {

	targetDateStart := targetDate.AddDate(0, 0, -durationInDays)
	compareDateStart := compareDate.AddDate(0, 0, -durationInDays)

	compareExpenses, err := expense.GetExpenses(compareDateStart, compareDate)
	if err != nil {
		fmt.Println("--- error getting expenses", err)
		return nil, err
	}
	if len(compareExpenses) == 0 {
		fmt.Println("--- no expenses found for period")
		return nil, NO_EXPENSES_ERR
	}

	targetExpenses, err := expense.GetExpenses(targetDateStart, targetDate)
	if err != nil {
		return nil, err
	}
	if len(targetExpenses) == 0 {
		return nil, NO_EXPENSES_ERR
	}

	targetTotal := expense.CalcTotal(targetExpenses)
	compareTotal := expense.CalcTotal(compareExpenses)
	targetSessions := len(targetExpenses)
	compareSessions := len(compareExpenses)

	compare := stats.PeriodComparison{
		TargetDateStart:  targetDateStart,
		TargetDateEnd:    targetDate,
		CompareDateStart: compareDateStart,
		CompareDateEnd:   compareDate,
		Metrics: map[string]stats.Metric{
			"Spent": stats.ComparisonMetric[int]{
				Name:          "Spent",
				Type:          stats.Amount,
				TargetValue:   targetTotal,
				CompareValue:  compareTotal,
				IncreaseValue: targetTotal - compareTotal,
			},
			"Sessions": stats.ComparisonMetric[int]{
				Name:          "Sessions",
				Type:          stats.Amount,
				TargetValue:   targetSessions,
				CompareValue:  compareSessions,
				IncreaseValue: targetSessions - compareSessions,
			},
		},
	}

	fmt.Println("--- compare:", compare.TargetDateStart, compare.Metrics)

	return &compare, nil
}

func getMoMComparison(targetDate time.Time) (*stats.PeriodComparison, error) {
	durationInDays := 30
	compareDate := targetDate.AddDate(0, 0, -durationInDays)
	return getPeriodComparison(targetDate, compareDate, durationInDays)
}

func getPrevMonthComparison() (*stats.PeriodComparison, error) {
	today := time.Now()
	todayLastMonth := today.AddDate(0, -1, 0)
	firstDayOfMonth := getFirstDayOfCurrentMonth()
	firstDayLastMonth := firstDayOfMonth.AddDate(0, -1, 0)

	if todayLastMonth.Month() != firstDayLastMonth.Month() {
		todayLastMonth = getLastDayOfMonth(firstDayLastMonth)
	}

	return getPeriodComparison(
		today,
		todayLastMonth,
		todayLastMonth.Day(),
	)

}

func Index(c echo.Context) error {
	currentDate := time.Now()
	dateParam := c.QueryParam("date")

	var date time.Time
	var err error

	if dateParam != "" {
		date, err = time.Parse("2006-01-02", dateParam)
	} else {
		date = currentDate
	}

	isCurrentMonth := date.Year() == currentDate.Year() &&
		date.Month() == currentDate.Month()

	if err != nil {
		fmt.Println("--- error parsing date param:", err)
		return err
	}

	startDate := getFirstDayOfMonth(date)
	endDate := getLastDayOfMonth(date)
	expenses, err := expense.GetExpenses(startDate, endDate)

	fmt.Println("--- found", len(expenses), "expenses")

	ranking := getGroceryStoresRanking(expenses)

	fmt.Println("--- found ", len(ranking), "items in ranking")

	ctx := utils.TemplContext{
		"name":           "Gastão",
		"isCurrentMonth": isCurrentMonth,
	}

	prevMonthCompare, err := getPrevMonthComparison()

	if err != nil {
		if err == NO_EXPENSES_ERR {
			return utils.Render(
				c,
				http.StatusOK,
				v.Index(date, expenses, nil, ranking),
				ctx,
			)
		}
		return err
	}

	return utils.Render(
		c,
		http.StatusOK,
		v.Index(
			date,
			expenses,
			prevMonthCompare,
			ranking,
		),
		ctx,
	)
}
