package usecase

import (
	"context"
	"fmt"

	"test_task/internal/entity"
	"test_task/internal/logger"
	"test_task/internal/notificator"
)

//go:generate go run github.com/golang/mock/mockgen --source=user.go --destination=user_mock.go --package=usecase

type UserRepository interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	Update(ctx context.Context, user entity.User) (entity.User, error)
	Delete(ctx context.Context, id string) error
	GetList(ctx context.Context, query entity.UserFilter) ([]entity.User, int64, error)
}

type Notificator interface {
	Push(ctx context.Context, data notificator.Notification) error
}

type User struct {
	repo        UserRepository
	notificator Notificator
	logger      logger.Logger
}

func NewUser(repo UserRepository, notificator Notificator, l logger.Logger) *User {
	return &User{repo: repo, notificator: notificator, logger: l}
}

func (u *User) Create(ctx context.Context, user entity.User) (entity.User, error) {
	createdUser, err := u.repo.Create(ctx, user)
	if err != nil {
		return entity.User{}, fmt.Errorf("repo create user: %w", err)
	}
	err = u.notificator.Push(ctx, notificator.Notification{
		Type: notificator.Insert,
		Data: createdUser,
	})
	if err != nil {
		u.logger.Error(fmt.Errorf("user create: push notification: %w", err))
	}
	return createdUser, nil
}

func (u *User) Update(ctx context.Context, user entity.User) (entity.User, error) {
	updatedUser, err := u.repo.Update(ctx, user)
	if err != nil {
		return entity.User{}, fmt.Errorf("repo update user: %w", err)
	}
	err = u.notificator.Push(ctx, notificator.Notification{
		Type: notificator.Update,
		Data: updatedUser,
	})
	if err != nil {
		u.logger.Error(fmt.Errorf("user update: push notification: %w", err))
	}
	return updatedUser, nil
}

func (u *User) Delete(ctx context.Context, id string) error {
	err := u.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("repo delete user: %w", err)
	}
	err = u.notificator.Push(ctx, notificator.Notification{
		Type: notificator.Delete,
		Data: entity.User{ID: id},
	})
	if err != nil {
		u.logger.Error(fmt.Errorf("user delete: push notification: %w", err))
	}
	return nil
}

func (u *User) GetList(ctx context.Context, filter entity.UserFilter) ([]entity.User, int64, error) {
	res, total, err := u.repo.GetList(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("repo getList user: %w", err)
	}
	return res, total, nil
}
