package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/purwandi/grafana-stack/src/pkg/logger"
	"github.com/purwandi/grafana-stack/src/service/location"
	"github.com/purwandi/grafana-stack/src/service/user"
	"go.uber.org/zap"
)

// Authorize ....
func Authorize(c echo.Context) error {
	userService := user.Service{}
	locationService := location.Service{}
	span, ctx := opentracing.StartSpanFromContext(c.Request().Context(), "authorize")
	span.SetTag("net.service", "user")
	span.SetTag("net.endpoint", "/")
	logger.InfoWithContext(ctx, "Info", zap.String("hello", "hello"))

	// find user
	if err := userService.FindAll(ctx); err != nil {
		logger.ErrorWithContext(ctx, "")
	}

	// find user
	if err := locationService.FindAll(ctx); err != nil {
		logger.ErrorWithContext(ctx, "")
	}

	span.Finish()

	return c.JSON(http.StatusOK, map[string]string{
		"message": "ok",
	})
}
