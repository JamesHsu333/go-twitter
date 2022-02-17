//go:generate mockgen -source usecase.go -destination mock/usecase_mock.go -package mock
package user

import (
	"context"

	"github.com/JamesHsu333/go-twitter/internal/models"
	"github.com/JamesHsu333/go-twitter/pkg/utils"
	"github.com/google/uuid"
)

// User repository interface
type UseCase interface {
	Register(ctx context.Context, user *models.User) (*models.UserWithToken, error)
	Login(ctx context.Context, user *models.User) (*models.UserWithToken, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	UpdateRole(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, userID uuid.UUID) error
	GetByID(ctx context.Context, selfID uuid.UUID, userID uuid.UUID) (*models.User, error)
	GetCacheByID(ctx context.Context, selfID uuid.UUID, userID uuid.UUID) (*models.User, error)
	GetByUserName(ctx context.Context, selfID uuid.UUID, userName string) (*models.User, error)
	FindByName(ctx context.Context, selfID uuid.UUID, name string, query *utils.PaginationQuery) (*models.UsersList, error)
	GetUsers(ctx context.Context, selfID uuid.UUID, pq *utils.PaginationQuery) (*models.UsersList, error)
}
