package home

import (
	"database/sql"
	"fmt"
	"log"
	expense "mercado/app/models/expense"
	v "mercado/app/views/home"
	"mercado/utils"
	htmx "mercado/utils/htmx"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type LayoutData struct {
	Expenses      []expense.Expense
	GroceryStores []string
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

	fmt.Println(
		"--- will open db with url: ",
		dbURL,
		" and engine: ",
		dbEngine,
		" in stage: ",
		stage,
	)

	db, err := sql.Open(dbEngine, dbURL)
	if err != nil {
		log.Printf(`--- error opening db: %v`, err)
	}
	return db, err
}

func formatDate(date time.Time) string {

	formatted := strconv.Itoa(
		date.Day(),
	) + " de " + PT_MONTHS[date.Month()-1] + ", " + date.Format(
		"15h:04m",
	)

	return formatted
}

func parseDate(dateStr string) (time.Time, error) {
	t, error := time.Parse("2006-01-02 15:04:05", dateStr)

	if error != nil {
		t, error = time.Parse("2006-01-02T15:04:05Z", dateStr)
	}

	return t, error
}

func validateGroceryStore(groceryStore string) bool {
	return utils.Contains(expense.GROCERY_STORES, groceryStore)
}

func formatCurrency(value int) string {

	whole := value / 100
	decimal := value % 100
	decimalStr := strconv.Itoa(decimal)
	paddedDecimal := fmt.Sprintf("%02s", decimalStr)

	return "€" + strconv.Itoa(whole) + "," + paddedDecimal
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

func Index(c echo.Context) error {

	date := time.Now()
	thirtyDaysAgo := date.AddDate(0, 0, -50)
	expenses, err := expense.GetExpenses(thirtyDaysAgo, date)

	if err != nil {
		fmt.Println("--- error getting expenses: ", err)
		return err
	}

	return utils.Render(c, http.StatusOK, v.Index(expenses), nil)
}

func AddExpense(c echo.Context) error {

	setErrorHeaders := func() {
		header := c.Response().Header()
		htmx.Retarget(header, "find span.error-message")
		htmx.Reswap(header, "innerHTML")
	}

	valueStr := c.Request().FormValue("value")
	groceryStore := c.Request().FormValue("grocery-store")

	if valueStr == "" {
		setErrorHeaders()
		return c.String(http.StatusOK, "O valor deve ser informado")
	}

	valueStr = strings.TrimSpace(valueStr)
	valueStr = strings.ReplaceAll(valueStr, ".", "")
	valueStr = strings.ReplaceAll(valueStr, ",", ".")

	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		setErrorHeaders()
		return c.String(http.StatusOK, "O valor deve ser um número")
	}

	if value <= 0 {
		setErrorHeaders()
		return c.String(http.StatusOK, "O valor deve ser positivo e não-zero")
	}

	if !validateGroceryStore(groceryStore) {
		var msg string

		if groceryStore == "" {
			setErrorHeaders()
			msg = "Véi, seleciona um mercado 😤"
		} else {
			msg = "Nome de mercado inválido: " + groceryStore
		}

		return c.String(http.StatusOK, msg)
	}

	db, err := openDB()
	defer db.Close()

	date := time.Now().In(time.UTC)
	formattedDate := date.Format("2006-01-02 15:04:05")

	result, err := db.Exec(
		"INSERT INTO expenses (value, grocery_store, date) VALUES (?, ?, ?)",
		int(value),
		groceryStore,
		formattedDate,
	)

	if err != nil {
		log.Printf("--- error adding expense: %v", err)
		return c.String(http.StatusOK, "Falha ao adicionar")
	}

	rowsAffected, err := result.RowsAffected()

	if rowsAffected == 0 {
		return c.String(http.StatusOK, "Falha ao adicionar")
	}

	if err != nil {
		log.Printf("Error: %v", err)
		return c.String(http.StatusOK, "Falha ao adicionar")
	}

	id, _ := result.LastInsertId()

	return utils.Render(
		c,
		http.StatusOK,
		v.ExpensesListItem(expense.Expense{
			ID:           int(id),
			Value:        int(value),
			GroceryStore: groceryStore,
			Date:         date,
		}),
		nil,
	)

}

func EditExpense(c echo.Context) error {
	id := c.Param("id")
	value := c.Request().FormValue("value")
	groceryStore := c.Request().FormValue("grocery-store")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Printf(`--- error converting id to int: %v`, err)
		return c.String(http.StatusOK, "Error!")
	}

	valueFloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Printf(`--- error value to float: %v`, err)
		return c.String(http.StatusOK, "Error!")
	}

	if !validateGroceryStore(groceryStore) {
		return c.String(
			http.StatusOK,
			"Nome de mercado inválido: "+groceryStore,
		)
	}

	if id == "" || err != nil {
		log.Printf(`--- invalid id: %v`, id)
		return c.String(http.StatusOK, "Error!")
	}

	valueFloat = valueFloat * 100

	db, err := openDB()
	defer db.Close()

	result, err := db.Exec(
		"UPDATE expenses SET value = ?, grocery_store = ? WHERE id = ?",
		int(valueFloat),
		groceryStore,
		idInt,
	)

	if err != nil {
		return c.String(http.StatusOK, "Error!")
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil || rowsAffected == 0 {
		return c.String(http.StatusOK, "Falha ao editar")
	}

	rows, err := db.Query("SELECT date from expenses WHERE id = ?", idInt)
	rows.Next()
	defer rows.Close()

	var dateStr string
	rows.Scan(&dateStr)

	dateTime, _ := parseDate(dateStr)

	expense := expense.Expense{
		ID:           idInt,
		Value:        int(valueFloat),
		GroceryStore: groceryStore,
		Date:         dateTime,
	}

	return utils.Render(
		c,
		http.StatusOK,
		v.ExpensesListItem(expense),
		nil,
	)

}

func DeleteExpense(c echo.Context) error {

	expenseIDStr := c.Param("id")
	expenseID, err := strconv.Atoi(expenseIDStr)
	if err != nil {
		return c.String(http.StatusOK, "Error!")
	}

	db, err := openDB()
	defer db.Close()

	result, err := db.Exec("DELETE FROM expenses WHERE id = ?", expenseID)
	rowsAffected, err := result.RowsAffected()

	if err != nil || rowsAffected == 0 {
		return c.String(http.StatusOK, "Falha ao apagar!")
	}

	return c.NoContent(http.StatusOK)
}
