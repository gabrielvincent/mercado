package expense

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

func openDB() (*sql.DB, error) {

	stage := os.Getenv("STAGE")
	var dbURL string
	var dbEngine string

	if stage == "prod" {
		dbURL = os.Getenv("TURSO_DB_CONNECTION_STRING")
		dbEngine = "libsql"
	}
	if stage == "dev" {
		dbURL = "./database.sqlite"
		dbEngine = "sqlite3"
	}

	db, err := sql.Open(dbEngine, dbURL)
	if err != nil {
		log.Fatal(err)
	}
	return db, err
}

func parseDate(dateStr string) (time.Time, error) {
	t, error := time.Parse("2006-01-02 15:04:05", dateStr)

	if error != nil {
		t, error = time.Parse("2006-01-02T15:04:05Z", dateStr)
	}

	return t, error
}

type Expense struct {
	ID           int
	Value        int
	GroceryStore string
	Date         time.Time
}

func GetExpenses(startDate time.Time, endDate time.Time) ([]Expense, error) {

	formattedStartDate := startDate.Format("2006-01-02 15:04:05")
	formattedEndDate := endDate.Format("2006-01-02 15:04:05")

	db, err := openDB()
	defer db.Close()

	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
        SELECT *
        FROM expenses
        WHERE date BETWEEN '%s' AND '%s'
        ORDER BY date DESC;
    `, formattedStartDate, formattedEndDate)

	rows, err := db.Query(query)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var expenses []Expense
	for rows.Next() {
		var id int
		var value int
		var groceryStore string
		var dateStr string

		err := rows.Scan(&id, &value, &groceryStore, &dateStr)
		if err != nil {
			log.Fatal(err)
		}

		date, error := parseDate(dateStr)

		if error != nil {
			log.Fatal(error)
			panic(error)
		}

		expenses = append(expenses, Expense{
			ID:           id,
			Value:        value,
			GroceryStore: groceryStore,
			Date:         date,
		})
	}

	return expenses, nil
}
