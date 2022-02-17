package tweet

import (
	"context"

	"github.com/JamesHsu333/go-twitter/internal/models"
	"github.com/JamesHsu333/go-twitter/pkg/utils"
	"github.com/google/uuid"
)

// Tweet repository interface
type Repository interface {
	Create(ctx context.Context, tweet *models.Tweet) (*models.Tweet, error)
	CreateReply(ctx context.Context, tweetID uint64, replyID uint64) error
	CheckTweetExist(ctx context.Context, tweetID uint64) error
	GetTweetByID(ctx context.Context, selfID uuid.UUID, tweetID uint64) (*models.TweetWithUser, error)
	GetTweets(ctx context.Context, selfID uuid.UUID, pq *utils.PaginationQuery) (*models.TweetsList, error)
	GetTweetsByUserID(ctx context.Context, self uuid.UUID, userID uuid.UUID, pq *utils.PaginationQuery) (*models.TweetsList, error)
	GetReplyTweets(ctx context.Context, selfID uuid.UUID, tweetID uint64, pq *utils.PaginationQuery) (*models.TweetsList, error)
	Delete(ctx context.Context, tweetID uint64) error
}
