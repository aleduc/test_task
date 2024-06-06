package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"test_task/internal/entity"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	FirstName string
	LastName  string
	Nickname  string
	Password  string
	Email     string
	Country   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func MapEntityUserToModelUser(user entity.User) (u User, err error) {
	id := uuid.Nil
	if user.ID != "" {
		id, err = uuid.Parse(user.ID)
		if err != nil {
			return u, fmt.Errorf("user ID is not uuid compatible: %w", err)
		}
	}

	return User{
		ID:        id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Nickname:  user.Nickname,
		Password:  user.Password,
		Email:     user.Email,
		Country:   user.Country,
	}, nil
}

func MapModelUserToEntityUser(user User) entity.User {
	return entity.User{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Nickname:  user.Nickname,
		Password:  user.Password,
		Email:     user.Email,
		Country:   user.Country,
	}
}

func MapModelUsersToEntityUsers(users []User) []entity.User {
	res := make([]entity.User, 0, len(users))
	for _, v := range users {
		res = append(res, MapModelUserToEntityUser(v))
	}
	return res
}
