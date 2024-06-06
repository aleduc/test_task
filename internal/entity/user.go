package entity

import (
	"test_task/internal/pagination"
)

type User struct {
	ID        string
	FirstName string
	LastName  string
	Nickname  string
	Password  string
	Email     string
	Country   string
}

// UserFilter business layer filter struct.
// I used separate structure intentionally. might be the case when UserFilter != User(email domain for example).
type UserFilter struct {
	ID        string
	FirstName string
	LastName  string
	Nickname  string
	Email     string
	Country   string
	pagination.Pagination
	pagination.Order
}
