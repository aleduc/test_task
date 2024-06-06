package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestHealth_View(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(mockHealthChecker1 *MockHealthChecker, mockHealthChecker2 *MockHealthChecker)
		expectedStatus int
		expectedErr    error
	}{
		{
			name: "all health checks pass",
			mockSetup: func(mockHealthChecker1 *MockHealthChecker, mockHealthChecker2 *MockHealthChecker) {
				mockHealthChecker1.EXPECT().Health(gomock.Any()).Return(nil)
				mockHealthChecker2.EXPECT().Health(gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedErr:    nil,
		},
		{
			name: "first health check fails",
			mockSetup: func(mockHealthChecker1 *MockHealthChecker, mockHealthChecker2 *MockHealthChecker) {
				mockHealthChecker1.EXPECT().Health(gomock.Any()).Return(errors.New("health check failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedErr:    nil,
		},
		{
			name: "second health check fails",
			mockSetup: func(mockHealthChecker1 *MockHealthChecker, mockHealthChecker2 *MockHealthChecker) {
				mockHealthChecker1.EXPECT().Health(gomock.Any()).Return(nil)
				mockHealthChecker2.EXPECT().Health(gomock.Any()).Return(errors.New("health check failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ao := assert.New(t)
			ctrl := gomock.NewController(t)
			mockHealthChecker1 := NewMockHealthChecker(ctrl)
			mockHealthChecker2 := NewMockHealthChecker(ctrl)

			tt.mockSetup(mockHealthChecker1, mockHealthChecker2)

			e := echo.New()
			handler := NewHealthController([]HealthChecker{mockHealthChecker1, mockHealthChecker2})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.View(c)

			ao.Equal(tt.expectedStatus, rec.Code)
			ao.Equal(tt.expectedErr, err)
		})
	}
}
