package utils

import (
	"context"
	"errors"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Render(
	echoCtx echo.Context,
	statusCode int,
	component templ.Component,
	ctx context.Context,
) error {
	echoCtx.Response().Writer.WriteHeader(statusCode)
	echoCtx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	if ctx == nil {
		ctx = echoCtx.Request().Context()
	}

	err := component.Render(ctx, echoCtx.Response().Writer)

	return err
}

func RenderAsync(
	echoCtx echo.Context,
	statusCode int,
	component templ.Component,
	ctx context.Context,
	channel chan templ.Component,
) error {

	w := echoCtx.Response().Writer

	flusher, ok := w.(http.Flusher)
	if !ok {
		return errors.New("streaming unsupported")
	}

	w.WriteHeader(statusCode)
	echoCtx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	if ctx == nil {
		ctx = echoCtx.Request().Context()
	}

	err := component.Render(ctx, w)
	flusher.Flush()

	deferComponent := <-channel

	if deferComponent != nil {
		err = deferComponent.Render(ctx, w)
		flusher.Flush()
	}

	return err
}
