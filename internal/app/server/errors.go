package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// EchoCustomError is responsible for echo unified error handling.
type EchoCustomError struct{}

// NewEchoCustomError creates new instance of EchoCustomError.
func NewEchoCustomError() *EchoCustomError {
	return &EchoCustomError{}
}

// Handler handles echo.Handlers errors.
func (EchoCustomError) Handler(err error, ctx echo.Context) {
	code := http.StatusInternalServerError
	message := "Internal Server Error"
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		if m, ok := he.Message.(string); ok {
			message = m
		} else {
			message = fmt.Sprintf("%v", he.Message)
		}
	}
	if !ctx.Response().Committed {
		_ = ctx.JSON(code, map[string]interface{}{
			"code":    code,
			"message": message,
		})
	}
}
