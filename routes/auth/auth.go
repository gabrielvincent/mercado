package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"mercado/utils"
)

var PASSWORDS = []string{
	"piá",
	"pia",
	"piazote",
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
	return utils.RenderInLayout(
		c,
		"login",
		map[string]interface{}{
			"Error": errorMsg,
		},
		nil,
	)
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
	cookie.Path = "/"
	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "/")
}
