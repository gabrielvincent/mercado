package main

import (
	authRoutes "mercado/routes/auth"
	homeRoutes "mercado/routes/home"
	"mercado/utils"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

func validateEnv() {
	var varsNotSet []string = []string{}
	stage := os.Getenv("STAGE")
	if stage == "" {
		varsNotSet = append(varsNotSet, "STAGE")
	}

	if stage == "prod" {
		tursoDBConnectionString := os.Getenv("TURSO_DB_CONNECTION_STRING")
		if tursoDBConnectionString == "" {
			varsNotSet = append(varsNotSet, "TURSO_DB_CONNECTION_STRING")
		}
	}

	if len(varsNotSet) == 0 {
		return
	}

	varsNotSetMsg := strings.Join(varsNotSet, ", ")

	panic("Missing environment variables: " + varsNotSetMsg)
}

func main() {
	validateEnv()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e := echo.New()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			path := c.Request().URL.Path
			if strings.HasPrefix(path, "/public") {
				return next(c)
			}

			if c.Request().URL.Path == "/login" {
				return next(c)
			}

			cookie, err := c.Cookie("password")

			if err != nil || cookie.Value == "" {
				return c.Redirect(http.StatusFound, "/login")
			}

			password := strings.ToLower(cookie.Value)

			if !utils.Contains(authRoutes.PASSWORDS, password) {
				return c.Redirect(http.StatusFound, "/login")
			}

			return next(c)
		}
	})

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path
			if strings.HasSuffix(path, ".min.js") {
				c.Response().
					Header().
					Set("cache-control", "max-age=31536000, public")
			}
			return next(c)
		}
	})

	e.GET("/", homeRoutes.Index)
	e.POST("/", homeRoutes.AddExpense)
	e.DELETE("/:id", homeRoutes.DeleteExpense)
	e.POST("/edit/:id", homeRoutes.EditExpense)

	e.GET("/login", authRoutes.Index)
	e.POST("/login", authRoutes.Login)

	e.Static("/public", "public")

	e.Logger.Fatal(e.Start(":" + port))
}
