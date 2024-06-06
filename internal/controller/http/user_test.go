package http

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"test_task/internal/entity"
	"test_task/internal/logger"
)

func TestUser_Create(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(mockUserUseCase *MockUserUseCase, mockLogger *logger.MockLogger)
		expectedStatus int
		expectedBody   string
		expectedErr    error
	}{
		{
			name:        "successful creation",
			requestBody: `{"first_name":"John", "last_name":"Doe", "nickname":"jdoe", "password":"password123", "email":"jdoe@example.com", "country":"USA"}`,
			mockSetup: func(mockUserUseCase *MockUserUseCase, mockLogger *logger.MockLogger) {
				mockUserUseCase.EXPECT().Create(gomock.Any(), gomock.Any()).Return(entity.User{
					ID:        "1",
					FirstName: "John",
					LastName:  "Doe",
					Nickname:  "jdoe",
					Password:  "password123",
					Email:     "jdoe@example.com",
					Country:   "USA",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"data":{"id":"1","first_name":"John","last_name":"Doe","nickname":"jdoe","password":"password123","email":"jdoe@example.com","country":"USA"}}` + "\n",
			expectedErr:    nil,
		},
		{
			name:        "failed creation due to bind error",
			requestBody: `invalid json`,
			mockSetup: func(mockUserUseCase *MockUserUseCase, mockLogger *logger.MockLogger) {
				mockLogger.EXPECT().Error(gomock.Any())
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Bad Request"}` + "\n",
			expectedErr:    nil,
		},
		{
			name:        "failed creation due to service error",
			requestBody: `{"first_name":"John", "last_name":"Doe", "nickname":"jdoe", "password":"password123", "email":"jdoe@example.com", "country":"USA"}`,
			mockSetup: func(mockUserUseCase *MockUserUseCase, mockLogger *logger.MockLogger) {
				mockUserUseCase.EXPECT().Create(gomock.Any(), gomock.Any()).Return(entity.User{}, errors.New("service error"))
				mockLogger.EXPECT().Error(fmt.Errorf("user create: %w", errors.New("service error")))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedErr:    nil,
			expectedBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ao := assert.New(t)
			ctrl := gomock.NewController(t)
			mockUserUseCase := NewMockUserUseCase(ctrl)
			mockLogger := logger.NewMockLogger(ctrl)

			tt.mockSetup(mockUserUseCase, mockLogger)

			e := echo.New()
			handler := NewUserHandler(mockUserUseCase, mockLogger)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte(tt.requestBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.Create(c)
			ao.Equal(tt.expectedErr, err)
			ao.Equal(tt.expectedStatus, rec.Code)
			ao.Equal(tt.expectedBody, rec.Body.String())

		})
	}
}

func TestUser_Update(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(mockUserUseCase *MockUserUseCase, mockLogger *logger.MockLogger)
		expectedStatus int
		expectedBody   string
		expectedErr    error
	}{
		{
			name:        "successful update",
			requestBody: `{"id":"1","first_name":"John","last_name":"Doe","nickname":"jdoe","password":"password123","email":"jdoe@example.com","country":"USA"}`,
			mockSetup: func(mockUserUseCase *MockUserUseCase, mockLogger *logger.MockLogger) {
				mockUserUseCase.EXPECT().Update(gomock.Any(), gomock.Any()).Return(entity.User{
					ID:        "1",
					FirstName: "John",
					LastName:  "Doe",
					Nickname:  "jdoe",
					Password:  "password123",
					Email:     "jdoe@example.com",
					Country:   "USA",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"data":{"id":"1","first_name":"John","last_name":"Doe","nickname":"jdoe","password":"password123","email":"jdoe@example.com","country":"USA"}}` + "\n",
			expectedErr:    nil,
		},
		{
			name:        "failed update due to bind error",
			requestBody: `invalid json`,
			mockSetup: func(mockUserUseCase *MockUserUseCase, mockLogger *logger.MockLogger) {
				mockLogger.EXPECT().Error(gomock.Any())
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Bad Request"}` + "\n",
			expectedErr:    nil,
		},
		{
			name:        "failed update due to service error",
			requestBody: `{"id":"1","first_name":"John","last_name":"Doe","nickname":"jdoe","password":"password123","email":"jdoe@example.com","country":"USA"}`,
			mockSetup: func(mockUserUseCase *MockUserUseCase, mockLogger *logger.MockLogger) {
				mockUserUseCase.EXPECT().Update(gomock.Any(), gomock.Any()).Return(entity.User{}, errors.New("service error"))
				mockLogger.EXPECT().Error(fmt.Errorf("user update: %w", errors.New("service error")))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   ``,
			expectedErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ao := assert.New(t)
			ctrl := gomock.NewController(t)
			mockUserUseCase := NewMockUserUseCase(ctrl)
			mockLogger := logger.NewMockLogger(ctrl)

			tt.mockSetup(mockUserUseCase, mockLogger)

			e := echo.New()
			handler := NewUserHandler(mockUserUseCase, mockLogger)

			req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader([]byte(tt.requestBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.Update(c)

			ao.Equal(tt.expectedStatus, rec.Code)
			ao.Equal(tt.expectedBody, rec.Body.String())
			ao.Equal(tt.expectedErr, err)
		})
	}
}

func TestUser_Delete(t *testing.T) {
	tests := []struct {
		name           string
		paramID        string
		mockSetup      func(mockUserUseCase *MockUserUseCase, mockLogger *logger.MockLogger)
		expectedStatus int
		expectedBody   string
		expectedErr    error
	}{
		{
			name:    "successful deletion",
			paramID: "1",
			mockSetup: func(mockUserUseCase *MockUserUseCase, mockLogger *logger.MockLogger) {
				mockUserUseCase.EXPECT().Delete(gomock.Any(), "1").Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   ``,
			expectedErr:    nil,
		},
		{
			name:    "failed deletion due to service error",
			paramID: "1",
			mockSetup: func(mockUserUseCase *MockUserUseCase, mockLogger *logger.MockLogger) {
				mockUserUseCase.EXPECT().Delete(gomock.Any(), "1").Return(errors.New("service error"))
				mockLogger.EXPECT().Error(gomock.Any()).Times(1)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   ``,
			expectedErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ao := assert.New(t)
			ctrl := gomock.NewController(t)
			mockUserUseCase := NewMockUserUseCase(ctrl)
			mockLogger := logger.NewMockLogger(ctrl)

			tt.mockSetup(mockUserUseCase, mockLogger)

			e := echo.New()
			handler := NewUserHandler(mockUserUseCase, mockLogger)

			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.paramID)

			err := handler.Delete(c)

			ao.Equal(tt.expectedStatus, rec.Code)
			ao.Equal(tt.expectedBody, rec.Body.String())
			ao.Equal(tt.expectedErr, err)
		})
	}
}

func TestUser_List(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(mockUserUseCase *MockUserUseCase, mockLogger *logger.MockLogger)
		expectedStatus int
		expectedBody   string
		expectedErr    error
	}{
		{
			name:        "successful list retrieval",
			requestBody: `{"page":1,"size":10}`,
			mockSetup: func(mockUserUseCase *MockUserUseCase, mockLogger *logger.MockLogger) {
				mockUserUseCase.EXPECT().GetList(gomock.Any(), gomock.Any()).Return([]entity.User{
					{
						ID:        "1",
						FirstName: "John",
						LastName:  "Doe",
						Nickname:  "jdoe",
						Email:     "jdoe@example.com",
						Country:   "USA",
					},
				}, int64(1), nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"pagination":{"current_page":1,"last_page":1,"total":1},"data":{"Users":[{"id":"1","first_name":"John","last_name":"Doe","nickname":"jdoe","email":"jdoe@example.com","country":"USA"}]}}` + "\n",
			expectedErr:    nil,
		},
		{
			name:        "failed list retrieval due to bind error",
			requestBody: `invalid json`,
			mockSetup: func(mockUserUseCase *MockUserUseCase, mockLogger *logger.MockLogger) {
				mockLogger.EXPECT().Error(gomock.Any()).Times(1)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Bad Request"}` + "\n",
			expectedErr:    nil,
		},
		{
			name:        "failed list retrieval due to service error",
			requestBody: `{"pagination":{"number":1,"size":10}}`,
			mockSetup: func(mockUserUseCase *MockUserUseCase, mockLogger *logger.MockLogger) {
				mockUserUseCase.EXPECT().GetList(gomock.Any(), gomock.Any()).Return(nil, int64(0), errors.New("service error"))
				mockLogger.EXPECT().Error(gomock.Any()).Times(1)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   ``,
			expectedErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ao := assert.New(t)
			ctrl := gomock.NewController(t)
			mockUserUseCase := NewMockUserUseCase(ctrl)
			mockLogger := logger.NewMockLogger(ctrl)

			tt.mockSetup(mockUserUseCase, mockLogger)

			e := echo.New()
			handler := NewUserHandler(mockUserUseCase, mockLogger)

			req := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader([]byte(tt.requestBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.List(c)

			ao.Equal(tt.expectedStatus, rec.Code)
			ao.Equal(tt.expectedBody, rec.Body.String())
			ao.Equal(tt.expectedErr, err)
		})
	}
}
