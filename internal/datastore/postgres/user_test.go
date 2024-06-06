package postgres

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"test_task/internal/entity"
	"test_task/internal/pagination"
)

func TestUserRepository_Create(t *testing.T) {
	type testCase struct {
		name        string
		input       entity.User
		mockSetup   func(sqlmock.Sqlmock)
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "successful creation",
			input: entity.User{
				ID:        "3d6f0eb1-2b1e-4d0f-b1a0-52f2b9249e85",
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "jdoe",
				Password:  "password123",
				Email:     "jdoe@example.com",
				Country:   "USA",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {

				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO \"users\" (\"first_name\",\"last_name\",\"nickname\",\"password\",\"email\",\"country\",\"created_at\",\"updated_at\",\"id\") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING \"id\"")).
					WithArgs("John", "Doe", "jdoe", "password123", "jdoe@example.com", "USA", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("3d6f0eb1-2b1e-4d0f-b1a0-52f2b9249e85"))
				mock.ExpectCommit()
			},
			expectedErr: nil,
		},
		{
			name: "failed creation due to mapping error",
			input: entity.User{
				ID:        "invalid-uuid",
				FirstName: "Invalid",
			},
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expectedErr: errors.New("MapEntityUserToModelUser: user ID is not uuid compatible: invalid UUID length: 12"),
		},
		{
			name: "failed creation due to database error",
			input: entity.User{
				ID:        "3d6f0eb1-2b1e-4d0f-b1a0-52f2b9249e85",
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "jdoe",
				Password:  "password123",
				Email:     "jdoe@example.com",
				Country:   "USA",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO \"users\" (\"first_name\",\"last_name\",\"nickname\",\"password\",\"email\",\"country\",\"created_at\",\"updated_at\",\"id\") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING \"id\"")).
					WithArgs("John", "Doe", "jdoe", "password123", "jdoe@example.com", "USA", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ao := assert.New(t)
			db, mock, err := sqlmock.New()
			ao.NoError(err)
			defer db.Close()

			tc.mockSetup(mock)

			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: db,
			}), &gorm.Config{})
			ao.NoError(err)

			repo := NewUserRepository(gormDB)
			_, err = repo.Create(context.Background(), tc.input)
			if tc.expectedErr != nil {
				ao.EqualError(err, tc.expectedErr.Error())
			} else {
				ao.NoError(err)
			}

			ao.NoError(mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	type testCase struct {
		name        string
		input       entity.User
		mockSetup   func(sqlmock.Sqlmock)
		expectedErr error
	}

	testCases := []testCase{
		{
			name: "successful update",
			input: entity.User{
				ID:        "3d6f0eb1-2b1e-4d0f-b1a0-52f2b9249e85",
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "jdoe",
				Password:  "password123",
				Email:     "jdoe@example.com",
				Country:   "USA",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE \"users\" SET \"first_name\"=$1,\"last_name\"=$2,\"nickname\"=$3,\"password\"=$4,\"email\"=$5,\"country\"=$6,\"updated_at\"=$7 WHERE id = $8")).
					WithArgs("John", "Doe", "jdoe", "password123", "jdoe@example.com", "USA", sqlmock.AnyArg(), "3d6f0eb1-2b1e-4d0f-b1a0-52f2b9249e85").
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			expectedErr: nil,
		},
		{
			name: "failed update due to mapping error",
			input: entity.User{
				ID:        "invalid-uuid",
				FirstName: "Invalid",
			},
			mockSetup:   func(mock sqlmock.Sqlmock) {},
			expectedErr: errors.New("MapEntityUserToModelUser: user ID is not uuid compatible: invalid UUID length: 12"),
		},
		{
			name: "failed update due to database error",
			input: entity.User{
				ID:        "3d6f0eb1-2b1e-4d0f-b1a0-52f2b9249e85",
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "jdoe",
				Password:  "password123",
				Email:     "jdoe@example.com",
				Country:   "USA",
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("UPDATE \"users\" SET \"first_name\"=$1,\"last_name\"=$2,\"nickname\"=$3,\"password\"=$4,\"email\"=$5,\"country\"=$6,\"updated_at\"=$7 WHERE id = $8")).
					WithArgs("John", "Doe", "jdoe", "password123", "jdoe@example.com", "USA", sqlmock.AnyArg(), "3d6f0eb1-2b1e-4d0f-b1a0-52f2b9249e85").
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ao := assert.New(t)
			db, mock, err := sqlmock.New()
			ao.NoError(err)
			defer db.Close()

			tc.mockSetup(mock)

			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: db,
			}), &gorm.Config{})
			ao.NoError(err)

			repo := NewUserRepository(gormDB)
			_, err = repo.Update(context.Background(), tc.input)
			if tc.expectedErr != nil {
				ao.EqualError(err, tc.expectedErr.Error())
			} else {
				ao.NoError(err)
			}

			ao.NoError(mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_Delete(t *testing.T) {
	type testCase struct {
		name        string
		input       string
		mockSetup   func(sqlmock.Sqlmock)
		expectedErr error
	}

	testCases := []testCase{
		{
			name:  "successful deletion",
			input: "3d6f0eb1-2b1e-4d0f-b1a0-52f2b9249e85",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "users"`).
					WithArgs("3d6f0eb1-2b1e-4d0f-b1a0-52f2b9249e85").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedErr: nil,
		},
		{
			name:  "failed deletion due to database error",
			input: "3d6f0eb1-2b1e-4d0f-b1a0-52f2b9249e85",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "users"`).
					WithArgs("3d6f0eb1-2b1e-4d0f-b1a0-52f2b9249e85").
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			expectedErr: errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ao := assert.New(t)
			db, mock, err := sqlmock.New()
			ao.NoError(err)
			defer db.Close()

			tc.mockSetup(mock)

			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: db,
			}), &gorm.Config{})
			ao.NoError(err)

			repo := NewUserRepository(gormDB)
			err = repo.Delete(context.Background(), tc.input)
			if tc.expectedErr != nil {
				ao.EqualError(err, tc.expectedErr.Error())
			} else {
				ao.NoError(err)
			}

			ao.NoError(mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_GetList(t *testing.T) {
	type testCase struct {
		name          string
		input         entity.UserFilter
		mockSetup     func(sqlmock.Sqlmock)
		expectedRes   []entity.User
		expectedTotal int64
		expectedErr   error
	}

	testCases := []testCase{
		{
			name: "successful retrieval",
			input: entity.UserFilter{
				Pagination: pagination.Pagination{
					Number: 1,
					Size:   10,
				},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "nickname", "password", "email", "country"}).
					AddRow("3d6f0eb1-2b1e-4d0f-b1a0-52f2b9249e85", "John", "Doe", "jdoe", "password123", "jdoe@example.com", "USA")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM \"users\"")).
					WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow("33"))
				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WillReturnRows(rows)
			},
			expectedRes: []entity.User{
				{
					ID:        "3d6f0eb1-2b1e-4d0f-b1a0-52f2b9249e85",
					FirstName: "John",
					LastName:  "Doe",
					Nickname:  "jdoe",
					Password:  "password123",
					Email:     "jdoe@example.com",
					Country:   "USA",
				},
			},
			expectedTotal: 33,
			expectedErr:   nil,
		},
		{
			name: "failed retrieval due to database error",
			input: entity.UserFilter{
				Pagination: pagination.Pagination{
					Number: 1,
					Size:   10,
				},
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM \"users\"")).
					WillReturnRows(sqlmock.NewRows([]string{"count(*)"}).AddRow("33"))
				mock.ExpectQuery(`SELECT \* FROM "users"`).
					WillReturnError(errors.New("db error"))
			},
			expectedRes: nil,
			expectedErr: errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ao := assert.New(t)
			db, mock, err := sqlmock.New()
			ao.NoError(err)
			defer db.Close()

			tc.mockSetup(mock)

			gormDB, err := gorm.Open(postgres.New(postgres.Config{
				Conn: db,
			}), &gorm.Config{})
			ao.NoError(err)

			repo := NewUserRepository(gormDB)
			result, total, err := repo.GetList(context.Background(), tc.input)
			if tc.expectedErr != nil {
				ao.EqualError(err, tc.expectedErr.Error())
			} else {
				ao.NoError(err)
				ao.Equal(tc.expectedRes, result)
				ao.Equal(tc.expectedTotal, total)
			}

			ao.NoError(mock.ExpectationsWereMet())
		})
	}
}
