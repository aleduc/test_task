package model

import (
	"testing"

	"test_task/internal/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMapEntityUserToModelUser(t *testing.T) {
	type testCase struct {
		name     string
		user     entity.User
		expected User
		hasError bool
	}

	testCases := []testCase{
		{
			name: "valid user with UUID",
			user: entity.User{
				ID:        "3d6f0eb1-2b1e-4d0f-b1a0-52f2b9249e85",
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "jdoe",
				Password:  "password123",
				Email:     "jdoe@example.com",
				Country:   "USA",
			},
			expected: User{
				ID:        uuid.MustParse("3d6f0eb1-2b1e-4d0f-b1a0-52f2b9249e85"),
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "jdoe",
				Password:  "password123",
				Email:     "jdoe@example.com",
				Country:   "USA",
			},
			hasError: false,
		},
		{
			name: "empty ID",
			user: entity.User{
				FirstName: "Jane",
				LastName:  "Smith",
				Nickname:  "jsmith",
				Password:  "password456",
				Email:     "jsmith@example.com",
				Country:   "Canada",
			},
			expected: User{
				ID:        uuid.Nil,
				FirstName: "Jane",
				LastName:  "Smith",
				Nickname:  "jsmith",
				Password:  "password456",
				Email:     "jsmith@example.com",
				Country:   "Canada",
			},
			hasError: false,
		},
		{
			name: "invalid UUID",
			user: entity.User{
				ID:        "invalid-uuid",
				FirstName: "Invalid",
			},
			expected: User{},
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := MapEntityUserToModelUser(tc.user)
			ao := assert.New(t)
			if tc.hasError {
				ao.Error(err)
			} else {
				ao.NoError(err)
				ao.Equal(tc.expected, result)
			}
		})
	}
}

func TestMapModelUserToEntityUser(t *testing.T) {
	type testCase struct {
		name     string
		user     User
		expected entity.User
	}

	testCases := []testCase{
		{
			name: "valid user",
			user: User{
				ID:        uuid.MustParse("3d6f0eb1-2b1e-4d0f-b1a0-52f2b9249e85"),
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "jdoe",
				Password:  "password123",
				Email:     "jdoe@example.com",
				Country:   "USA",
			},
			expected: entity.User{
				ID:        "3d6f0eb1-2b1e-4d0f-b1a0-52f2b9249e85",
				FirstName: "John",
				LastName:  "Doe",
				Nickname:  "jdoe",
				Password:  "password123",
				Email:     "jdoe@example.com",
				Country:   "USA",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := MapModelUserToEntityUser(tc.user)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMapModelUsersToEntityUsers(t *testing.T) {
	type testCase struct {
		name     string
		users    []User
		expected []entity.User
	}

	testCases := []testCase{
		{
			name: "multiple users",
			users: []User{
				{
					ID:        uuid.MustParse("3d6f0eb1-2b1e-4d0f-b1a0-52f2b9249e85"),
					FirstName: "John",
					LastName:  "Doe",
					Nickname:  "jdoe",
					Password:  "password123",
					Email:     "jdoe@example.com",
					Country:   "USA",
				},
				{
					ID:        uuid.MustParse("1e2d3c4b-5a6f-7e8d-9b0c-1d2e3f4a5b6c"),
					FirstName: "Jane",
					LastName:  "Smith",
					Nickname:  "jsmith",
					Password:  "password456",
					Email:     "jsmith@example.com",
					Country:   "Canada",
				},
			},
			expected: []entity.User{
				{
					ID:        "3d6f0eb1-2b1e-4d0f-b1a0-52f2b9249e85",
					FirstName: "John",
					LastName:  "Doe",
					Nickname:  "jdoe",
					Password:  "password123",
					Email:     "jdoe@example.com",
					Country:   "USA",
				},
				{
					ID:        "1e2d3c4b-5a6f-7e8d-9b0c-1d2e3f4a5b6c",
					FirstName: "Jane",
					LastName:  "Smith",
					Nickname:  "jsmith",
					Password:  "password456",
					Email:     "jsmith@example.com",
					Country:   "Canada",
				},
			},
		},
		{
			name:     "no users",
			users:    []User{},
			expected: []entity.User{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := MapModelUsersToEntityUsers(tc.users)
			assert.Equal(t, tc.expected, result)
		})
	}
}
