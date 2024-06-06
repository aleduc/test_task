package http

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestInitRoutes(t *testing.T) {
	e := echo.New()
	InitRoutes(e, Controllers{User: nil})
	tests := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodPut,
			path:   fmt.Sprintf(APIv1 + "users/:id"),
		},
		{
			method: http.MethodDelete,
			path:   fmt.Sprintf(APIv1 + "users/:id"),
		},
		{
			method: http.MethodGet,
			path:   fmt.Sprintf(APIv1 + "users"),
		},
		{
			method: http.MethodPost,
			path:   fmt.Sprintf(APIv1 + "users"),
		},
		{
			method: http.MethodGet,
			path:   fmt.Sprintf(APIv1 + "health"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			found := false
			for _, r := range e.Routes() {
				if r.Method == tt.method && r.Path == tt.path {
					found = true
					break
				}
			}
			assert.False(t, found, "Route %s %s should be registered", tt.method, tt.path)
		})
	}
}
