package like

import (
	"context"

	"github.com/JamesHsu333/go-twitter/internal/models"
	"github.com/JamesHsu333/go-twitter/pkg/utils"
	"github.com/google/uuid"
)

type UseCase interface {
	Like(ctx context.Context, userID uuid.UUID, tweetID uint64) error
	GetLikedTweets(ctx context.Context, userID uuid.UUID, pq *utils.PaginationQuery) (*models.TweetsList, error)
	GetLikedUsers(ctx context.Context, selfID uuid.UUID, tweetID uint64, pq *utils.PaginationQuery) (*models.UsersList, error)
	Delete(ctx context.Context, userID uuid.UUID, tweetID uint64) error
}
