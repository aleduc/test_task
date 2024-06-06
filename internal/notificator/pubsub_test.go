package notificator

import (
	"context"
	"errors"
	"testing"
	"time"

	"test_task/internal/logger"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPubSub_Push(t *testing.T) {
	type testCase struct {
		name      string
		data      Notification
		expectErr bool
	}

	testCases := []testCase{
		{
			name:      "successful push",
			data:      Notification{Type: "test", Data: "test data"},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ao := assert.New(t)
			ctrl := gomock.NewController(t)
			mockNotificator := NewMockRepositoryNotificator(ctrl)
			mockLogger := logger.NewMockLogger(ctrl)
			ps := NewPubSub(mockNotificator, 10, mockLogger)

			err := ps.Push(context.Background(), tc.data)
			if tc.expectErr {
				ao.Error(err)
			} else {
				ao.NoError(err)
			}
		})
	}
}

func TestPubSub_Start(t *testing.T) {
	type testCase struct {
		name          string
		notifications []Notification
		setupMocks    func(mockNotificator *MockRepositoryNotificator, mockLogger *logger.MockLogger)
		expectErr     bool
	}

	testCases := []testCase{
		{
			name:          "successful processing",
			notifications: []Notification{{Type: "test", Data: "test data"}},
			setupMocks: func(mockNotificator *MockRepositoryNotificator, mockLogger *logger.MockLogger) {
				mockNotificator.EXPECT().Push(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectErr: false,
		},
		{
			name:          "marshal error",
			notifications: []Notification{{Type: "test", Data: make(chan int)}},
			setupMocks: func(mockNotificator *MockRepositoryNotificator, mockLogger *logger.MockLogger) {
				mockLogger.EXPECT().Error(gomock.Any()).AnyTimes()
			},
			expectErr: true,
		},
		{
			name:          "push error",
			notifications: []Notification{{Type: "test", Data: "test data"}},
			setupMocks: func(mockNotificator *MockRepositoryNotificator, mockLogger *logger.MockLogger) {
				mockNotificator.EXPECT().Push(gomock.Any(), gomock.Any()).Return(errors.New("push error"))
				mockLogger.EXPECT().Error(gomock.Any()).AnyTimes()
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockNotificator := NewMockRepositoryNotificator(ctrl)
			mockLogger := logger.NewMockLogger(ctrl)
			ps := NewPubSub(mockNotificator, 10, mockLogger)

			tc.setupMocks(mockNotificator, mockLogger)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go ps.Start(ctx)
			for _, n := range tc.notifications {
				_ = ps.Push(ctx, n)
			}
			time.Sleep(100 * time.Millisecond)
		})
	}
}

func TestPubSub_Stop(t *testing.T) {
	type testCase struct {
		name          string
		notifications []Notification
		setupMocks    func(mockLogger *logger.MockLogger)
	}

	testCases := []testCase{
		{
			name:          "empty buffer",
			notifications: nil,
			setupMocks: func(mockLogger *logger.MockLogger) {
				mockLogger.EXPECT().Info("notification consumer stop started")
				mockLogger.EXPECT().Info("notification consumer stop finished")
			},
		},
		{
			name:          "non-empty buffer with timeout",
			notifications: []Notification{{Type: "test", Data: "test data"}},
			setupMocks: func(mockLogger *logger.MockLogger) {
				logs := []string{
					"notification consumer stop started",
					"1 notifications should be processed",
					"1 notifications should have been processed, but they weren't",
					"notification consumer stop finished",
				}
				for _, log := range logs {
					mockLogger.EXPECT().Info(log).Times(1)
				}

			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockNotificator := NewMockRepositoryNotificator(ctrl)
			mockLogger := logger.NewMockLogger(ctrl)
			ps := NewPubSub(mockNotificator, 10, mockLogger)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			tc.setupMocks(mockLogger)
			for _, n := range tc.notifications {
				_ = ps.Push(ctx, n)
			}
			ps.Stop(cancel, 150*time.Millisecond, 200*time.Millisecond)
		})
	}
}
