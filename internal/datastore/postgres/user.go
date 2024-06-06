package postgres

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"test_task/internal/datastore/postgres/model"
	"test_task/internal/entity"
	"test_task/internal/pagination"
)

type UserRepository struct {
	pgClient *gorm.DB
}

func NewUserRepository(pgClient *gorm.DB) *UserRepository {
	return &UserRepository{pgClient: pgClient}
}

func (u *UserRepository) Create(ctx context.Context, user entity.User) (entity.User, error) {
	modelUser, err := model.MapEntityUserToModelUser(user)
	if err != nil {
		return entity.User{}, fmt.Errorf("MapEntityUserToModelUser: %w", err)
	}

	err = u.pgClient.WithContext(ctx).Create(&modelUser).Error
	return model.MapModelUserToEntityUser(modelUser), err
}

func (u *UserRepository) Update(ctx context.Context, user entity.User) (entity.User, error) {
	modelUser, err := model.MapEntityUserToModelUser(user)
	if err != nil {
		return entity.User{}, fmt.Errorf("MapEntityUserToModelUser: %w", err)
	}
	err = u.pgClient.WithContext(ctx).Model(&model.User{}).
		Where("id = ?", modelUser.ID).
		Select("FirstName", "LastName", "Nickname", "Password", "Email", "Country").
		Updates(&modelUser).Error
	return model.MapModelUserToEntityUser(modelUser), err
}

func (u *UserRepository) Delete(ctx context.Context, id string) error {
	return u.pgClient.WithContext(ctx).Unscoped().Delete(&model.User{}, "id = ?", id).Error
}

func (u *UserRepository) GetList(ctx context.Context, query entity.UserFilter) ([]entity.User, int64, error) {
	var (
		res   = make([]model.User, 0)
		total int64
	)
	err := u.pgClient.WithContext(ctx).Model(&model.User{}).Scopes(
		StringEqFilterScope("id", query.ID),
		StringEqFilterScope("first_name", query.FirstName),
		StringEqFilterScope("last_name", query.LastName),
		StringEqFilterScope("nickname", query.Nickname),
		StringEqFilterScope("email", query.Email),
		StringEqFilterScope("country", query.Country),
	).Count(&total).
		Limit(query.Pagination.Size).
		Offset(pagination.CalculateOffset(query.Pagination.Number, query.Pagination.Size)).
		Find(&res).Error

	return model.MapModelUsersToEntityUsers(res), total, err
}
