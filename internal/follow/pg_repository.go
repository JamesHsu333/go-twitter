package follow

import (
	"context"

	"github.com/JamesHsu333/go-twitter/internal/models"
	"github.com/JamesHsu333/go-twitter/pkg/utils"
	"github.com/google/uuid"
)

type Repository interface {
	Follow(ctx context.Context, follower uuid.UUID, following uuid.UUID) error
	GetFollowers(ctx context.Context, selfID uuid.UUID, userID uuid.UUID, pq *utils.PaginationQuery) (*models.UsersList, error)
	GetFollowing(ctx context.Context, selfID uuid.UUID, userID uuid.UUID, pq *utils.PaginationQuery) (*models.UsersList, error)
	Delete(ctx context.Context, follower uuid.UUID, following uuid.UUID) error
}
