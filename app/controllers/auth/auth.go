package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"

	v "mercado/app/views/auth"
	"mercado/utils"
)

var PASSWORDS = []string{
	"piá",
	"pia",
	"piazote",
	"piazin",
	"nilda",
	"nildas",
	"valentina",
	"poceta",
	"poce",
	"coelhuda",
	"buchinho",
	"buchinhos",
	"buchos",
	"tchueia",
	"tchueias",
	"tchuba",
	"tchubas",
	"tchubilda",
	"tchubicuda",
	"tchubicudas",
	"mozinho",
	"muzis",
	"chuncho",
	"chunchuncho",
}

func renderLoginPage(c echo.Context, errorMsg string) error {
	return utils.Render(c, http.StatusOK, v.Index(errorMsg), nil)
}

func Index(c echo.Context) error {
	return renderLoginPage(c, "")
}

func Login(c echo.Context) error {
	password := c.Request().FormValue("password")

	if !utils.Contains(PASSWORDS, password) {
		return renderLoginPage(c, "Senha inválida. Assim como você.")

	}

	cookie := new(http.Cookie)
	cookie.Name = "password"
	cookie.Value = password
	cookie.MaxAge = 60 * 60 * 24 * 365
	cookie.Path = "/"
	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "/")
}
