package middleware

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
)

// OpenTracing ...
func OpenTracing(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		serverSpan := opentracing.StartSpan(fmt.Sprintf("%s %s", c.Path(), c.Request().Method))
		defer serverSpan.Finish()

		ctx := opentracing.ContextWithSpan(c.Request().Context(), serverSpan)
		c.SetRequest(c.Request().WithContext(ctx))

		serverSpan.SetTag("http.method", c.Request().Method)
		serverSpan.SetTag("http.url", c.Request().URL)
		serverSpan.SetTag("http.status_code", c.Response().Status)

		err := next(c)
		if err != nil {
			serverSpan.SetTag("error", true)
		}

		return err
	}
}
