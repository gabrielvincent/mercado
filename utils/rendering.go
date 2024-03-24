package utils

import (
	"context"
	"io"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func concatTemplates(
	templatesA *template.Template,
	templatesB *template.Template,
) (*template.Template, error) {

	templates := templatesA

	for _, t := range templatesB.Templates() {
		var err error = nil
		templates, err = templates.AddParseTree(t.Name(), t.Tree)
		if err != nil {
			return nil, err
		}
	}

	return templates, nil
}

func (t *Template) Render(
	w io.Writer,
	name string,
	data interface{},
	c echo.Context,
) error {

	err := t.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		// Log the error
		c.Logger().Error(err)
	}
	return err
}

func RenderInLayout(
	c echo.Context,
	name string,
	data interface{},
	customFuncs template.FuncMap,
) error {
	partialsGlob := "public/views/partials/" + name + "/*.html"
	name = "public/views/" + name + ".html"
	partials, err := filepath.Glob(partialsGlob)

	allPaths := append(partials, name, "public/views/layout.html")

	tmpl, err := template.New(name).Funcs(customFuncs).ParseFiles(
		allPaths...,
	)

	if err != nil {
		c.Logger().Error(err)
		return err
	}

	c.Echo().Renderer = &Template{
		templates: tmpl,
	}

	return c.Render(http.StatusOK, "layout", data)
}

func RenderPartial(
	c echo.Context,
	partial string,
	data interface{},
) error {
	return c.Render(http.StatusOK, partial, data)
}

func Render(
	echoCtx echo.Context,
	statusCode int,
	t templ.Component,
	ctx context.Context,
) error {
	echoCtx.Response().Writer.WriteHeader(statusCode)
	echoCtx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	if ctx == nil {
		ctx = echoCtx.Request().Context()
	}

	return t.Render(ctx, echoCtx.Response().Writer)
}
