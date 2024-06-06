package dto

import (
	"test_task/internal/entity"
	"test_task/internal/pagination"
)

type (
	UserCreateRequest struct {
		UserInputCore
	}

	UserUpdateRequest struct {
		ID string `param:"id"`
		UserInputCore
	}

	UserListRequest struct {
		UserFilters
		pagination.Pagination
		pagination.Order
	}

	UserInputCore struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Nickname  string `json:"nickname"`
		Password  string `json:"password"`
		Email     string `json:"email"`
		Country   string `json:"country"`
	}

	UserFilters struct {
		ID        string `query:"filter.id" json:"filter.id"`
		FirstName string `query:"filter.first_name" json:"filter.first_name"`
		LastName  string `query:"filter.last_name" json:"filter.last_name"`
		Nickname  string `query:"filter.nickname" json:"filter.nickname"`
		Email     string `query:"filter.email" json:"filter.email"`
		Country   string `query:"filter.country" json:"filter.country"`
	}

	UserCRUResponse struct {
		ID        string `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Nickname  string `json:"nickname"`
		Password  string `json:"password"`
		Email     string `json:"email"`
		Country   string `json:"country"`
	}

	UserViewResponse struct {
		ID        string `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Nickname  string `json:"nickname"`
		Email     string `json:"email"`
		Country   string `json:"country"`
	}
	UserListResponse struct {
		Users []UserViewResponse
	}
)

func MapUserCreateRequestToEntity(request UserCreateRequest) entity.User {
	return MapUserCoreToEntity(request.UserInputCore)
}
func MapUserUpdateRequestToEntity(request UserUpdateRequest) entity.User {
	res := MapUserCoreToEntity(request.UserInputCore)
	res.ID = request.ID
	return res
}

func MapUserCoreToEntity(core UserInputCore) entity.User {
	return entity.User{
		FirstName: core.FirstName,
		LastName:  core.LastName,
		Nickname:  core.Nickname,
		Password:  core.Password,
		Email:     core.Email,
		Country:   core.Country,
	}
}

func MapUserListRequestToEntity(request UserListRequest) entity.UserFilter {
	return entity.UserFilter{
		ID:         request.UserFilters.ID,
		FirstName:  request.UserFilters.FirstName,
		LastName:   request.UserFilters.LastName,
		Nickname:   request.UserFilters.Nickname,
		Email:      request.UserFilters.Email,
		Country:    request.UserFilters.Country,
		Pagination: request.Pagination,
		Order:      request.Order,
	}
}

func MapUserToEntityUserResponse(entity entity.User) UserCRUResponse {
	return UserCRUResponse{
		ID:        entity.ID,
		FirstName: entity.FirstName,
		LastName:  entity.LastName,
		Nickname:  entity.Nickname,
		Password:  entity.Password,
		Email:     entity.Email,
		Country:   entity.Country,
	}
}

func MapUsersToEntityUserListResponse(entities []entity.User) UserListResponse {
	u := UserListResponse{Users: make([]UserViewResponse, 0, len(entities))}
	for _, v := range entities {
		u.Users = append(u.Users, UserViewResponse{
			ID:        v.ID,
			FirstName: v.FirstName,
			LastName:  v.LastName,
			Nickname:  v.Nickname,
			Email:     v.Email,
			Country:   v.Country,
		})
	}
	return u
}
