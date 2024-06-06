package usecase

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"test_task/internal/entity"
	"test_task/internal/logger"
	"test_task/internal/notificator"
	"test_task/internal/pagination"
)

func TestUser_Create(t *testing.T) {
	type testCase struct {
		name          string
		input         entity.User
		repoResult    entity.User
		repoError     error
		notifyError   error
		expectedError error
	}

	testCases := []testCase{
		{
			name:          "success",
			input:         entity.User{ID: "1", FirstName: "John"},
			repoResult:    entity.User{ID: "1", FirstName: "John"},
			repoError:     nil,
			notifyError:   nil,
			expectedError: nil,
		},
		{
			name:          "repo error",
			input:         entity.User{ID: "1", FirstName: "John"},
			repoResult:    entity.User{},
			repoError:     errors.New("repo error"),
			notifyError:   nil,
			expectedError: fmt.Errorf("repo create user: %w", errors.New("repo error")),
		},
		{
			name:          "notificator error",
			input:         entity.User{ID: "1", FirstName: "John"},
			repoResult:    entity.User{ID: "1", FirstName: "John"},
			repoError:     nil,
			notifyError:   errors.New("notification error"),
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ao := assert.New(t)

			mockRepo := NewMockUserRepository(ctrl)
			mockNotificator := notificator.NewMockNotificator(ctrl)
			mockLogger := logger.NewMockLogger(ctrl)
			u := NewUser(mockRepo, mockNotificator, mockLogger)

			mockRepo.EXPECT().Create(gomock.Any(), tc.input).Return(tc.repoResult, tc.repoError)
			if tc.repoError == nil {
				mockNotificator.EXPECT().Push(gomock.Any(), notificator.Notification{
					Type: notificator.Insert,
					Data: tc.input,
				}).Return(tc.notifyError)
				if tc.notifyError != nil {
					mockLogger.EXPECT().Error(gomock.Any()).Times(1)
				}
			}

			result, err := u.Create(context.Background(), tc.input)
			if tc.expectedError != nil {
				ao.Error(err)
				ao.Equal(tc.expectedError.Error(), err.Error())
			} else {
				ao.NoError(err)
				ao.Equal(tc.repoResult, result)
			}
		})
	}
}

func TestUser_Update(t *testing.T) {
	type testCase struct {
		name          string
		input         entity.User
		repoResult    entity.User
		repoError     error
		notifyError   error
		expectedError error
	}

	testCases := []testCase{
		{
			name:          "success",
			input:         entity.User{ID: "1", FirstName: "John"},
			repoResult:    entity.User{ID: "1", FirstName: "John"},
			repoError:     nil,
			notifyError:   nil,
			expectedError: nil,
		},
		{
			name:          "repo error",
			input:         entity.User{ID: "1", FirstName: "John"},
			repoResult:    entity.User{},
			repoError:     errors.New("repo error"),
			notifyError:   nil,
			expectedError: fmt.Errorf("repo update user: %w", errors.New("repo error")),
		},
		{
			name:          "notificator error",
			input:         entity.User{ID: "1", FirstName: "John"},
			repoResult:    entity.User{ID: "1", FirstName: "John"},
			repoError:     nil,
			notifyError:   errors.New("notification error"),
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ao := assert.New(t)
			mockRepo := NewMockUserRepository(ctrl)
			mockNotificator := notificator.NewMockNotificator(ctrl)
			mockLogger := logger.NewMockLogger(ctrl)
			u := NewUser(mockRepo, mockNotificator, mockLogger)

			mockRepo.EXPECT().Update(gomock.Any(), tc.input).Return(tc.repoResult, tc.repoError)
			if tc.repoError == nil {
				mockNotificator.EXPECT().Push(gomock.Any(), notificator.Notification{
					Type: notificator.Update,
					Data: tc.input,
				}).Return(tc.notifyError)
				if tc.notifyError != nil {
					mockLogger.EXPECT().Error(gomock.Any()).Times(1)
				}
			}

			result, err := u.Update(context.Background(), tc.input)
			if tc.expectedError != nil {
				ao.Error(err)
				ao.Equal(tc.expectedError.Error(), err.Error())
			} else {
				ao.NoError(err)
				ao.Equal(tc.repoResult, result)
			}
		})
	}
}

func TestUser_Delete(t *testing.T) {
	type testCase struct {
		name          string
		input         string
		repoError     error
		notifyError   error
		expectedError error
	}

	testCases := []testCase{
		{
			name:          "success",
			input:         "1",
			repoError:     nil,
			notifyError:   nil,
			expectedError: nil,
		},
		{
			name:          "repo error",
			input:         "1",
			repoError:     errors.New("repo error"),
			notifyError:   nil,
			expectedError: fmt.Errorf("repo delete user: %w", errors.New("repo error")),
		},
		{
			name:          "notificator error",
			input:         "1",
			repoError:     nil,
			notifyError:   errors.New("notification error"),
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ao := assert.New(t)

			mockRepo := NewMockUserRepository(ctrl)
			mockNotificator := notificator.NewMockNotificator(ctrl)
			mockLogger := logger.NewMockLogger(ctrl)
			u := NewUser(mockRepo, mockNotificator, mockLogger)

			mockRepo.EXPECT().Delete(gomock.Any(), tc.input).Return(tc.repoError)
			if tc.repoError == nil {
				mockNotificator.EXPECT().Push(gomock.Any(), notificator.Notification{
					Type: notificator.Delete,
					Data: entity.User{ID: tc.input},
				}).Return(tc.notifyError)
				if tc.notifyError != nil {
					mockLogger.EXPECT().Error(gomock.Any()).Times(1)
				}
			}

			err := u.Delete(context.Background(), tc.input)
			if tc.expectedError != nil {
				ao.Error(err)
				ao.Equal(tc.expectedError.Error(), err.Error())
			} else {
				ao.NoError(err)
			}
		})
	}
}

func TestUser_GetList(t *testing.T) {
	type testCase struct {
		name          string
		input         entity.UserFilter
		repoResult    []entity.User
		repoTotal     int64
		repoError     error
		expectedError error
	}

	testCases := []testCase{
		{
			name: "success",
			input: entity.UserFilter{Pagination: pagination.Pagination{
				Number: 1,
				Size:   10,
			}},
			repoResult:    []entity.User{{ID: "1", FirstName: "John"}},
			repoTotal:     1,
			repoError:     nil,
			expectedError: nil,
		},
		{
			name: "repo error",
			input: entity.UserFilter{Pagination: pagination.Pagination{
				Number: 1,
				Size:   10,
			}},
			repoResult:    nil,
			repoTotal:     0,
			repoError:     errors.New("repo error"),
			expectedError: fmt.Errorf("repo getList user: %w", errors.New("repo error")),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ao := assert.New(t)

			mockRepo := NewMockUserRepository(ctrl)
			mockNotificator := notificator.NewMockNotificator(ctrl)
			mockLogger := logger.NewMockLogger(ctrl)
			u := NewUser(mockRepo, mockNotificator, mockLogger)

			mockRepo.EXPECT().GetList(gomock.Any(), tc.input).Return(tc.repoResult, tc.repoTotal, tc.repoError)

			result, total, err := u.GetList(context.Background(), tc.input)
			if tc.expectedError != nil {
				ao.Error(err)
				ao.Equal(tc.expectedError.Error(), err.Error())
			} else {
				ao.NoError(err)
				ao.Equal(tc.repoResult, result)
				ao.Equal(tc.repoTotal, total)
			}
		})
	}
}
