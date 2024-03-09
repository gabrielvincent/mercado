package routes

import (
	"database/sql"
	"fmt"
	"log"
	homeComponents "mercado/components/home"
	model "mercado/models/home"
	"mercado/utils"
	htmx "mercado/utils/htmx"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type LayoutData struct {
	Expenses      []model.Expense
	GroceryStores []string
}

var GROCERY_STORES = []string{
	"Lidl",
	"Pingo Doce",
	"Minipreço",
	"El Corte Inglés",
	"Continente",
	"Mercadona",
	"Froiz",
	"Aldi",
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

	db, err := sql.Open(dbEngine, dbURL)
	if err != nil {
		log.Fatal(err)
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

	db, err := openDB()
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

	var expenses []model.Expense
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

		expenses = append(expenses, model.Expense{
			ID:           id,
			Value:        value,
			GroceryStore: groceryStore,
			Date:         date,
		})
	}

	return utils.Render(c, http.StatusOK, homeComponents.Index(expenses))
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
		homeComponents.ExpensesListItem(model.Expense{
			ID:           int(id),
			Value:        int(value),
			GroceryStore: groceryStore,
			Date:         date,
		}),
	)

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

	expense := model.Expense{
		ID:           idInt,
		Value:        int(valueFloat),
		GroceryStore: groceryStore,
		Date:         dateTime,
	}

	return utils.Render(
		c,
		http.StatusOK,
		homeComponents.ExpensesListItem(expense),
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
