// Package server implements methods for HTTP server configuration, start, shutdown.
package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"test_task/internal/config"
	"test_task/internal/logger"
)

// NewServer creates echo.Echo with configuration.
func NewServer(cfg config.HTTP) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Server.ReadTimeout = cfg.ReadTimeout
	e.Server.WriteTimeout = cfg.WriteTimeout
	e.Server.IdleTimeout = cfg.IdleTimeout
	e.HTTPErrorHandler = NewEchoCustomError().Handler
	e.Use(middleware.Recover())
	e.GET("/ping", func(ctx echo.Context) error {
		return ctx.NoContent(http.StatusOK)
	})

	return e
}

// ShutdownServer stops echo server.
func ShutdownServer(l logger.Logger, e *echo.Echo, cfg config.HTTP) {
	l.Info("stopping http server")
	timeoutCtx, cancelTimeout := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancelTimeout()
	if err := e.Shutdown(timeoutCtx); err != nil &&
		!errors.Is(err, context.Canceled) &&
		!errors.Is(err, http.ErrServerClosed) {
		l.Error(err)
	}
}
