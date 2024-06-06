package http

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:generate go run github.com/golang/mock/mockgen --source=health.go --destination=health_mock.go --package=http

type HealthChecker interface {
	Health(ctx context.Context) error
}

type Health struct {
	healthCheckers []HealthChecker
}

func NewHealthController(healthCheckers []HealthChecker) *Health {
	return &Health{healthCheckers: healthCheckers}
}

func (h *Health) View(ctx echo.Context) error {
	for _, v := range h.healthCheckers {
		err := v.Health(ctx.Request().Context())
		if err != nil {
			return ctx.NoContent(http.StatusInternalServerError)
		}
	}
	return ctx.NoContent(http.StatusOK)
}
