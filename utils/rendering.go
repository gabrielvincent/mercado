package utils

import (
	"context"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type TemplContext map[string]interface{}

func Render(
	echoCtx echo.Context,
	statusCode int,
	t templ.Component,
	templCtx TemplContext,
) error {

	echoCtx.Response().Writer.WriteHeader(statusCode)
	echoCtx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	ctx := context.Background()
	if templCtx != nil {
		for k, v := range templCtx {
			ctx = context.WithValue(ctx, k, v)
		}
	}

	return t.Render(ctx, echoCtx.Response().Writer)
}
