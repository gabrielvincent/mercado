package expense

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

var GROCERY_STORES = []string{
	"Assaí",
	"Campeão",
	"Carrefour",
	"Casas Pedro",
	"Guanabara",
	"Hortifruti Genérico",
	"Hortifruti",
	"Inter",
	"Mundial",
	"Merca Dez",
	"Padaria",
	"Prezunic",
	"Princesa",
	"Pão de Açúcar",
	"Rede Economia",
	"SuperPrix",
	"Supermarket",
	"Ultra",
	"Zona Sul",
	"Outro",
}

// var GROCERY_STORES = []string{
// 	"Aldi",
// 	"Auchan",
// 	"Continente",
// 	"El Corte Inglés",
// 	"Froiz",
// 	"Lidl",
// 	"Mercadona",
// 	"Minipreço",
// 	"Padaria",
// 	"Pingo Doce",
// 	"Pomar",
// }

type Expense struct {
	ID           int
	Value        int
	GroceryStore string
	Date         time.Time
}

type GroceryStoresRankingItem struct {
	GroceryStore string
	Sessions     int
}

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

func CalcAvg(expenses []Expense) float64 {
	total := CalcTotal(expenses)
	return float64(total) / float64(len(expenses))
}

func CalcTotal(expenses []Expense) int {
	total := 0
	for _, expense := range expenses {
		total += expense.Value
	}

	return total
}

func GetExpenses(startDate time.Time, endDate time.Time) ([]Expense, error) {
	formattedStartDate := startDate.Format("2006-01-02 15:04:05")
	formattedEndDate := endDate.Format("2006-01-02 15:04:05")

	db, err := openDB()
	defer db.Close()

	if err != nil {
		fmt.Println("--- error opening db: ", err)
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
		fmt.Println("--- error running query: ", err)
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
