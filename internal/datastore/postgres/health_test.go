package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestNewHealth(t *testing.T) {
	ao := assert.New(t)
	db, _, err := sqlmock.New()
	ao.NoError(err)
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	ao.NoError(err)

	h := NewHealth(gormDB)
	ao.NotNil(h)
	ao.Equal(gormDB, h.pgClient)
}

func TestHealth_Health(t *testing.T) {
	type testCase struct {
		name        string
		setupMock   func(mock sqlmock.Sqlmock)
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "success ping",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPing().WillReturnError(nil)
			},
			expectedErr: nil,
		},
		{
			name: "ping failure",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPing().WillReturnError(errors.New("ping error"))
			},
			expectedErr: errors.New("ping error"),
		},
		{
			name: "db failure",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPing().WillReturnError(sql.ErrConnDone)
			},
			expectedErr: sql.ErrConnDone,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ao := assert.New(t)
			db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			ao.NoError(err)
			defer db.Close()
			mock.ExpectPing().WillReturnError(nil) // for close
			tc.setupMock(mock)

			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: db,
			}), &gorm.Config{})
			ao.NoError(err)

			h := NewHealth(gormDB)
			err = h.Health(context.Background())

			ao.Equal(tc.expectedErr, err)
			ao.NoError(mock.ExpectationsWereMet())
		})
	}
}
