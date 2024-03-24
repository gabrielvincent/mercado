package stats

import (
	"context"
	"errors"
	expense "mercado/app/models/expense"
	v "mercado/app/views/stats"
	"mercado/utils"
	"net/http"
	"slices"
	"time"

	"github.com/labstack/echo/v4"
)

var NO_EXPENSES_ERR = errors.New("no expenses found")

func getCurrentDayOfMonth() int {
	now := time.Now()
	return now.Day()
}

func getPrevMonthCompareInfo(
	expenses []expense.Expense,
) (*expense.PreviousMonthCompareInfo, error) {
	currentDay := getCurrentDayOfMonth()
	previousMonth := time.Now().AddDate(0, -1, 0)
	previousMonthFirstDay := time.Date(
		previousMonth.Year(),
		previousMonth.Month(),
		1,
		0,
		0,
		0,
		0,
		time.UTC,
	)
	thisDayLastMonth := time.Date(
		previousMonth.Year(),
		previousMonth.Month(),
		currentDay,
		23,
		59,
		59,
		0,
		time.UTC,
	)
	previousMonthExpenses, err := expense.GetExpenses(
		previousMonthFirstDay,
		thisDayLastMonth,
	)

	if err != nil {
		return nil, err
	}

	if len(previousMonthExpenses) == 0 {
		return nil, NO_EXPENSES_ERR
	}

	currTotal := expense.CalcTotal(expenses)
	prevTotal := expense.CalcTotal(previousMonthExpenses)

	return &expense.PreviousMonthCompareInfo{
		Total:    prevTotal,
		Sessions: len(previousMonthExpenses),
		CompareTotalPercent: float64(
			prevTotal,
		) / float64(
			currTotal,
		),
		CompareGroceryStoreSessionsPercent: float64(
			len(previousMonthExpenses),
		) / float64(
			len(expenses),
		),
	}, nil
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

func Index(c echo.Context) error {

	startDate := getFirstDayOfCurrentMonth()
	endDate := getLastDayOfCurrentMonth()
	expenses, err := expense.GetExpenses(startDate, endDate)
	if err != nil {
		return err
	}

	groceryStoresRanking := getGroceryStoresRanking(expenses)
	prevMonthCompareInfo, err := getPrevMonthCompareInfo(expenses)
	if err != nil {
		if err == NO_EXPENSES_ERR {
			return utils.Render(
				c,
				http.StatusOK,
				v.Index(expenses, nil, groceryStoresRanking),
				nil,
			)
		}
		return err
	}

	ctx := context.WithValue(c.Request().Context(), "name", "Gastão")

	return utils.Render(
		c,
		http.StatusOK,
		v.Index(expenses, prevMonthCompareInfo, groceryStoresRanking),
		ctx,
	)
}
