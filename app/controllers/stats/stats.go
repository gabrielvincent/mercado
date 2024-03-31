package stats

import (
	"errors"
	"fmt"
	expense "mercado/app/models/expense"
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

func getPrevMonthCompareInfo(
	expenses []expense.Expense,
	date time.Time,
) (*expense.PreviousMonthCompareInfo, error) {
	previousMonth := date.AddDate(0, -1, 0)
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
	previousMonthLastDay := time.Date(
		previousMonth.Year(),
		previousMonth.Month()+1,
		0,
		23,
		59,
		59,
		0,
		time.UTC,
	)
	fmt.Println(
		"--- preivous month first day: ",
		previousMonthFirstDay,
		" previous month last day: ", previousMonthLastDay,
	)
	previousMonthExpenses, err := expense.GetExpenses(
		previousMonthFirstDay,
		previousMonthLastDay,
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

	prevMonthCompareInfo, err := getPrevMonthCompareInfo(expenses, date)

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
		v.Index(date, expenses, prevMonthCompareInfo, ranking),
		ctx,
	)
}
