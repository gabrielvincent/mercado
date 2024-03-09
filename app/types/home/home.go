package home

import "time"

type Expense struct {
	ID           int
	Value        int
	GroceryStore string
	Date         time.Time
}
