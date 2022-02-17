package user

import (
	"context"

	"github.com/JamesHsu333/go-twitter/internal/models"
	"github.com/JamesHsu333/go-twitter/pkg/utils"
	"github.com/google/uuid"
)

// User repository interface
type Repository interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, userID uuid.UUID) error
	GetByID(ctx context.Context, selfID uuid.UUID, userID uuid.UUID) (*models.User, error)
	GetByUserName(ctx context.Context, selfID uuid.UUID, userName string) (*models.User, error)
	FindByName(ctx context.Context, selfID uuid.UUID, name string, query *utils.PaginationQuery) (*models.UsersList, error)
	FindByEmail(ctx context.Context, user *models.User) (*models.User, error)
	GetUsers(ctx context.Context, selfID uuid.UUID, pq *utils.PaginationQuery) (*models.UsersList, error)
	UpdateRole(ctx context.Context, user *models.User) (*models.User, error)
}
