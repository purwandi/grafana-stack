package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/opentracing-contrib/go-zap/log"
	"github.com/opentracing/opentracing-go"
	"github.com/purwandi/grafana-stack/src/cmd/handler"
	"github.com/purwandi/grafana-stack/src/pkg/logger"
	appmiddleware "github.com/purwandi/grafana-stack/src/pkg/middleware"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

var (
	shutdowns []func() error
)

func init() {
	godotenv.Load()
	cfg, err := config.FromEnv()
	logger.Init("./storage/log/user/user.log")

	if err != nil {
		// parsing errors might happen here, such as when we get a string where we expect a number
		log.Error(fmt.Sprintf("Could not parse Jaeger env vars: %s", err.Error()))
		return
	}

	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		log.Error(fmt.Sprintf("Could not initialize jaeger tracer: %s", err.Error()))
		return
	}

	shutdowns = append(shutdowns, closer.Close)

	opentracing.SetGlobalTracer(tracer)
}

func main() {
	shutdown := make(chan struct{})
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(appmiddleware.OpenTracing)
	e.Use(appmiddleware.Logger)

	e.GET("/", handler.Authorize)

	go func() {
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}

		// close another module
		for i := range shutdowns {
			shutdowns[i]()
		}

		close(shutdown)
	}()

	log.Info("Server has been started successfully", zap.String("info", "Ok"))

	go func() {
		if err := e.Start(":8080"); err != nil {
			log.Info("shutting down server gracefully")
		}
	}()
	<-shutdown
}
