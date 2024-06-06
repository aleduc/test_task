package http

import (
	"github.com/labstack/echo/v4"
)

const (
	// APIv1 initial API group.
	APIv1 = "api/v1/"

	usersGroupName = "users"
)

// Controllers combines all handlers in one struct for a following routing.
type Controllers struct {
	User             *User
	HealthController *Health
}

// InitRoutes initializes all service routes.
func InitRoutes(e *echo.Echo, handlers Controllers) {
	apiV1Group := e.Group(APIv1)

	// init service routes
	apiV1Group.GET("health", handlers.HealthController.View)

	// init API
	NewUserRoutes(apiV1Group, handlers.User)
}

// NewUserRoutes registers routes for user entity.
func NewUserRoutes(e *echo.Group, h *User) {
	userGroup := e.Group(usersGroupName)
	userGroup.POST("", h.Create)
	userGroup.GET("", h.List)
	userGroup.PUT("/:id", h.Update)
	userGroup.DELETE("/:id", h.Delete)
}
