package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestEchoCustomError_Handler(t *testing.T) {
	e := echo.New()
	type args struct {
		err error
	}
	tests := []struct {
		name            string
		args            args
		expectedCode    int
		expectedMessage string
	}{
		{
			name:            "Internal Server Error",
			args:            args{err: fmt.Errorf("some internal error")},
			expectedCode:    http.StatusInternalServerError,
			expectedMessage: "Internal Server Error",
		},
		{
			name:            "HTTP Error with string message",
			args:            args{err: &echo.HTTPError{Code: http.StatusNotFound, Message: "Not Found"}},
			expectedCode:    http.StatusNotFound,
			expectedMessage: "Not Found",
		},
		{
			name:            "HTTP Error with non-string message",
			args:            args{err: &echo.HTTPError{Code: http.StatusBadRequest, Message: map[string]string{"error": "Bad Request"}}},
			expectedCode:    http.StatusBadRequest,
			expectedMessage: "map[error:Bad Request]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ao := assert.New(t)
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := &EchoCustomError{}
			handler.Handler(tt.args.err, c)
			ao.Equal(tt.expectedCode, rec.Code)
			ao.JSONEq(fmt.Sprintf(`{"code":%d,"message":"%s"}`, tt.expectedCode, tt.expectedMessage), rec.Body.String())
		})
	}
}
