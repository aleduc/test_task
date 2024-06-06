package notificator

import (
	"context"
)

//go:generate go run github.com/golang/mock/mockgen --source=notificator.go --destination=notificator_mock.go --package=notificator

type OperationType string

const (
	Insert OperationType = "Insert"
	Update OperationType = "Update"
	Delete OperationType = "Delete"
)

type Notification struct {
	Type OperationType
	Data interface{}
}

type Notificator interface {
	Push(ctx context.Context, data Notification) error
}

type RepositoryNotificator interface {
	Push(ctx context.Context, data []byte) error
}
