package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"test_task/internal/controller/http/dto"
	"test_task/internal/entity"
	"test_task/internal/logger"
	"test_task/internal/pagination"
)

//go:generate go run github.com/golang/mock/mockgen --source=user.go --destination=user_mock.go --package=http

// UserUseCase describes service methods for entity.User.
type UserUseCase interface {
	Create(ctx context.Context, user entity.User) (entity.User, error)
	Update(ctx context.Context, user entity.User) (entity.User, error)
	Delete(ctx context.Context, id string) error
	GetList(ctx context.Context, filter entity.UserFilter) ([]entity.User, int64, error)
}

// User is responsible for handling any user-related requests.
type User struct {
	userService UserUseCase
	logger      logger.Logger
}

// NewUserHandler created new Handler entity.
func NewUserHandler(userService UserUseCase, l logger.Logger) *User {
	return &User{userService: userService, logger: l}
}

func (u *User) Create(ctx echo.Context) error {
	var req dto.UserCreateRequest
	err := ctx.Bind(&req)
	if err != nil {
		u.logger.Error(fmt.Errorf("user create: bind: %w", err))
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": http.StatusText(http.StatusBadRequest)})
	}
	result, err := u.userService.Create(ctx.Request().Context(), dto.MapUserCreateRequestToEntity(req))
	if err != nil {
		u.logger.Error(fmt.Errorf("user create: %w", err))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, BaseResponse{Data: dto.MapUserToEntityUserResponse(result)})

}

func (u *User) Update(ctx echo.Context) error {
	var req dto.UserUpdateRequest
	err := ctx.Bind(&req)
	if err != nil {
		u.logger.Error(fmt.Errorf("user update: bind: %w", err))
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": http.StatusText(http.StatusBadRequest)})
	}
	result, err := u.userService.Update(ctx.Request().Context(), dto.MapUserUpdateRequestToEntity(req))
	if err != nil {
		u.logger.Error(fmt.Errorf("user update: %w", err))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, BaseResponse{Data: dto.MapUserToEntityUserResponse(result)})

}

func (u *User) Delete(ctx echo.Context) error {
	id := ctx.Param("id")
	err := u.userService.Delete(ctx.Request().Context(), id)
	if err != nil {
		u.logger.Error(fmt.Errorf("user delete: %w", err))
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.NoContent(http.StatusOK)
}

func (u *User) List(ctx echo.Context) error {
	var req dto.UserListRequest
	err := ctx.Bind(&req)
	if err != nil {
		u.logger.Error(fmt.Errorf("user list: bind: %w", err))
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": http.StatusText(http.StatusBadRequest)})
	}
	result, total, err := u.userService.GetList(ctx.Request().Context(), dto.MapUserListRequestToEntity(req))
	if err != nil {
		u.logger.Error(fmt.Errorf("user list: %w", err))
		/* Just an example of how it can be, but according to YAGNI I can omit it .
		if errors.Is(err, datastore.ErrNotFound)
			return ErrNotFound
		}*/
		return ctx.NoContent(http.StatusInternalServerError)
	}
	return ctx.JSON(http.StatusOK, PaginatedBaseResponse{
		Pagination: ResponsePagination{
			CurrentPage: req.Pagination.Number,
			LastPage:    pagination.CalculateLastPage(int(total), req.Pagination.Size),
			Total:       total,
		},
		Data: dto.MapUsersToEntityUserListResponse(result),
	})
}
