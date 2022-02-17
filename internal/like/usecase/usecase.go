package usecase

import (
	"context"

	"github.com/JamesHsu333/go-twitter/config"
	"github.com/JamesHsu333/go-twitter/internal/like"
	"github.com/JamesHsu333/go-twitter/internal/models"
	"github.com/JamesHsu333/go-twitter/pkg/logger"
	"github.com/JamesHsu333/go-twitter/pkg/tracer"
	"github.com/JamesHsu333/go-twitter/pkg/utils"
	"github.com/google/uuid"
)

type likeUC struct {
	cfg      *config.Config
	likeRepo like.Repository
	logger   logger.Logger
}

func NewLikeUseCase(cfg *config.Config, likeRepo like.Repository, logger logger.Logger) like.UseCase {
	return &likeUC{cfg: cfg, likeRepo: likeRepo, logger: logger}
}

func (u *likeUC) Like(ctx context.Context, userID uuid.UUID, tweetID uint64) error {
	ctx, span := tracer.NewSpan(ctx, "likeUC.Like", nil)
	defer span.End()

	return u.likeRepo.Like(ctx, userID, tweetID)
}

func (u *likeUC) GetLikedTweets(ctx context.Context, userID uuid.UUID, pq *utils.PaginationQuery) (*models.TweetsList, error) {
	ctx, span := tracer.NewSpan(ctx, "likeUC.GetLikedTweets", nil)
	defer span.End()

	return u.likeRepo.GetLikedTweets(ctx, userID, pq)
}

func (u *likeUC) GetLikedUsers(ctx context.Context, selfID uuid.UUID, tweetID uint64, pq *utils.PaginationQuery) (*models.UsersList, error) {
	ctx, span := tracer.NewSpan(ctx, "likeUC.GetLikedUsers", nil)
	defer span.End()

	return u.likeRepo.GetLikedUsers(ctx, selfID, tweetID, pq)
}

func (u *likeUC) Delete(ctx context.Context, userID uuid.UUID, tweetID uint64) error {
	ctx, span := tracer.NewSpan(ctx, "likeUC.Delete", nil)
	defer span.End()

	return u.likeRepo.Delete(ctx, userID, tweetID)
}
