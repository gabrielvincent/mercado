package routes

import (
	"database/sql"
	"fmt"
	"log"
	"mercado/utils"
	htmx "mercado/utils/htmx"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type Expense struct {
	ID           int
	Value        int
	GroceryStore string
	Date         time.Time
}

type LayoutData struct {
	Expenses      []Expense
	GroceryStores []string
}

var DB_URL = "libsql://mercado-gabrielvincent.turso.io?authToken=eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3MDk2NTAzNjQsImlkIjoiYWMzNGNmNGUtZjgwYi00YTMyLTgyMTktMjcyNTJkMzcwZDMzIn0.h-ACyvyjzis-XbA4z_mnEySGfT0S0-Ark9QRe47-r6ovUdn3QkEA3l1AYcJCEK-6zUTQ66nAtL-zhR7jJ5GBCQ"

var GROCERY_STORES = []string{
	"Lidl",
	"Pingo Doce",
	"Continente",
	"Mercadona",
	"Froiz",
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

func formatDate(date time.Time) string {

	formatted := strconv.Itoa(
		date.Day(),
	) + " de " + PT_MONTHS[date.Month()-1] + ", " + date.Format(
		"15h:04m",
	)

	return formatted
}

func getGroceryStores() []string {
	return GROCERY_STORES
}

func validateGroceryStore(groceryStore string) bool {
	return utils.Contains(GROCERY_STORES, groceryStore)
}

func formatCurrency(value int) string {

	whole := value / 100
	decimal := value % 100
	decimalStr := strconv.Itoa(decimal)
	paddedDecimal := fmt.Sprintf("%02s", decimalStr)

	return "€" + strconv.Itoa(whole) + "," + paddedDecimal
}

func Index(c echo.Context) error {

	db, err := sql.Open("libsql", DB_URL)
	defer db.Close()

	if err != nil {
		return err
	}

	rows, err := db.Query(`
        SELECT *
        FROM expenses
        WHERE strftime('%Y-%m', date) = strftime('%Y-%m', 'now')
        ORDER BY date DESC;
    `)
	defer rows.Close()

	if err != nil {
		return err
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

		date, error := time.Parse("2006-01-02 15:04:05", dateStr)

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

	return utils.RenderInLayout(
		c,
		"home",
		LayoutData{
			Expenses:      expenses,
			GroceryStores: GROCERY_STORES,
		},
		template.FuncMap{
			"getGroceryStores": getGroceryStores,
			"formatCurrency":   formatCurrency,
			"formatDate":       formatDate,
		},
	)
}

func AddExpense(c echo.Context) error {

	setErrorHeaders := func() {
		header := c.Response().Header()
		htmx.Retarget(header, "previous span.error-message")
		htmx.Reswap(header, "innerHTML")
	}

	valueStr := c.Request().FormValue("value")
	groceryStore := c.Request().FormValue("grocery-store")

	if valueStr == "" {
		setErrorHeaders()
		return c.String(http.StatusOK, "O valor deve ser informado")
	}

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		setErrorHeaders()
		return c.String(http.StatusOK, "O valor deve ser um número")
	}

	value = value * 100

	if value <= 0 {
		setErrorHeaders()
		return c.String(http.StatusOK, "O valor deve ser positivo e não-zero")
	}

	if !validateGroceryStore(groceryStore) {
		return c.String(
			http.StatusOK,
			"Nome de mercado inválido: "+groceryStore,
		)
	}

	db, err := sql.Open("libsql", DB_URL)
	defer db.Close()

	date := time.Now().In(time.UTC)
	formattedDate := date.Format("2006-01-02 15:04:05")

	result, err := db.Exec(
		"INSERT INTO expenses (value, grocery_store, date) VALUES (?, ?, ?)",
		int(value),
		groceryStore,
		formattedDate,
	)

	rowsAffected, err := result.RowsAffected()

	if rowsAffected == 0 {
		return c.String(http.StatusOK, "Falha ao adicionar")
	}

	if err != nil {
		log.Printf("Error: %v", err)
		return c.String(http.StatusOK, "Falha ao adicionar")
	}

	id, _ := result.LastInsertId()

	return utils.RenderPartial(
		c,
		"expense-list-item",
		Expense{
			ID:           int(id),
			Value:        int(value),
			GroceryStore: groceryStore,
			Date:         date,
		},
	)

}

func DeleteExpense(c echo.Context) error {

	expenseIDStr := c.Param("id")
	expenseID, err := strconv.Atoi(expenseIDStr)
	if err != nil {
		return c.String(http.StatusOK, "Error!")
	}

	db, err := sql.Open("libsql", DB_URL)
	defer db.Close()

	result, err := db.Exec("DELETE FROM expenses WHERE id = ?", expenseID)
	rowsAffected, err := result.RowsAffected()

	if err != nil || rowsAffected == 0 {
		return c.String(http.StatusOK, "Falha ao apagar!")
	}

	return c.NoContent(http.StatusOK)
}

func EditExpense(c echo.Context) error {
	id := c.Param("id")
	value := c.Request().FormValue("value")
	groceryStore := c.Request().FormValue("grocery-store")
	idInt, err := strconv.Atoi(id)
	valueFloat, err := strconv.ParseFloat(value, 64)

	if !validateGroceryStore(groceryStore) {
		return c.String(
			http.StatusOK,
			"Nome de mercado inválido: "+groceryStore,
		)
	}

	if id == "" || err != nil {
		log.Printf(`--- error: %v`, err)
		return c.String(http.StatusOK, "Error!")
	}

	valueFloat = valueFloat * 100

	db, err := sql.Open("libsql", DB_URL)
	defer db.Close()

	result, err := db.Exec(
		"UPDATE expenses SET value = ?, grocery_store = ? WHERE id = ?",
		int(valueFloat),
		groceryStore,
		idInt,
	)
	rowsAffected, err := result.RowsAffected()

	if err != nil || rowsAffected == 0 {
		return c.String(http.StatusOK, "Falha ao editar")
	}

	rows, err := db.Query("SELECT date from expenses WHERE id = ?", idInt)
	rows.Next()

	var date string
	rows.Scan(&date)

	dateTime, _ := time.Parse("2006-01-02 15:04:05", date)

	return utils.RenderPartial(
		c,
		"expense-list-item",
		Expense{
			ID:           idInt,
			Value:        int(valueFloat),
			GroceryStore: groceryStore,
			Date:         dateTime,
		},
	)

}
