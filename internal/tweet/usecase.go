//go:generate mockgen -source usecase.go -destination mock/usecase_mock.go -package mock
package tweet

import (
	"context"

	"github.com/JamesHsu333/go-twitter/internal/models"
	"github.com/JamesHsu333/go-twitter/pkg/utils"
	"github.com/google/uuid"
)

// Tweet usecase interface
type UseCase interface {
	Create(ctx context.Context, tweet *models.Tweet) (*models.Tweet, error)
	CreateReply(ctx context.Context, tweetID uint64, tweet *models.Tweet) (*models.Tweet, error)
	GetTweetByID(ctx context.Context, tweetID uint64) (*models.TweetWithUser, error)
	GetTweets(ctx context.Context, pq *utils.PaginationQuery) (*models.TweetsList, error)
	GetTweetsByUserID(ctx context.Context, userID uuid.UUID, pq *utils.PaginationQuery) (*models.TweetsList, error)
	GetReplyTweets(ctx context.Context, tweetID uint64, pq *utils.PaginationQuery) (*models.TweetsList, error)
	Delete(ctx context.Context, tweetID uint64) error
}
