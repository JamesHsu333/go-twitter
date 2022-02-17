package usecase

import (
	"context"

	"github.com/JamesHsu333/go-twitter/config"
	"github.com/JamesHsu333/go-twitter/internal/models"
	"github.com/JamesHsu333/go-twitter/internal/tweet"
	"github.com/JamesHsu333/go-twitter/pkg/httpErrors"
	"github.com/JamesHsu333/go-twitter/pkg/logger"
	"github.com/JamesHsu333/go-twitter/pkg/tracer"
	"github.com/JamesHsu333/go-twitter/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Tweet Usecase
type tweetUC struct {
	cfg       *config.Config
	tweetRepo tweet.Repository
	logger    logger.Logger
}

// New Usecase
func NewTweetUseCase(cfg *config.Config, tweetRepo tweet.Repository, logger logger.Logger) tweet.UseCase {
	return &tweetUC{cfg: cfg, tweetRepo: tweetRepo, logger: logger}
}

// Create new tweet
func (u *tweetUC) Create(ctx context.Context, tweet *models.Tweet) (*models.Tweet, error) {
	ctx, span := tracer.NewSpan(ctx, "tweetUC.Create", nil)
	defer span.End()

	self, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, httpErrors.NewUnauthorizedError(errors.WithMessage(err, "tweetUC.Create.GetUserFromCtx"))
	}

	tweet.UserID = self.UserID
	if err = utils.ValidateStruct(ctx, tweet); err != nil {
		tracer.AddSpanError(span, err)
		return nil, httpErrors.NewBadRequestError(errors.WithMessage(err, "tweetUC.Create.ValidateStruct"))
	}

	createdTweet, err := u.tweetRepo.Create(ctx, tweet)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	return createdTweet, nil
}

func (u *tweetUC) CreateReply(ctx context.Context, tweetID uint64, tweet *models.Tweet) (*models.Tweet, error) {
	ctx, span := tracer.NewSpan(ctx, "tweetUC.CreateReply", nil)
	defer span.End()

	if err := u.tweetRepo.CheckTweetExist(ctx, tweetID); err != nil {
		tracer.AddSpanError(span, err)
		return nil, httpErrors.NewNotFoundError(errors.WithMessage(err, "tweetUC.CreateReply.CheckTweetExist"))
	}

	self, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, httpErrors.NewUnauthorizedError(errors.WithMessage(err, "tweetUC.CreateReply.GetUserFromCtx"))
	}

	tweet.UserID = self.UserID
	if err = utils.ValidateStruct(ctx, tweet); err != nil {
		tracer.AddSpanError(span, err)
		return nil, httpErrors.NewBadRequestError(errors.WithMessage(err, "tweetUC.CreateReply.ValidateStruct"))
	}

	createdTweet, err := u.tweetRepo.Create(ctx, tweet)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	if err = u.tweetRepo.CreateReply(ctx, tweetID, createdTweet.ID); err != nil {
		tracer.AddSpanError(span, err)
		if err = u.tweetRepo.Delete(ctx, createdTweet.ID); err != nil {
			tracer.AddSpanError(span, err)
			return nil, err
		}
		return nil, err
	}

	return createdTweet, nil
}

// Get tweet by id
func (u *tweetUC) GetTweetByID(ctx context.Context, tweetID uint64) (*models.TweetWithUser, error) {
	ctx, span := tracer.NewSpan(ctx, "tweetUC.GetTweetByID", nil)
	defer span.End()

	self, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, httpErrors.NewUnauthorizedError(errors.WithMessage(err, "tweetUC.GetTweetByID.GetUserFromCtx"))
	}

	tweet, err := u.tweetRepo.GetTweetByID(ctx, self.UserID, tweetID)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	return tweet, nil
}

// Get tweets
func (u *tweetUC) GetTweets(ctx context.Context, pq *utils.PaginationQuery) (*models.TweetsList, error) {
	ctx, span := tracer.NewSpan(ctx, "tweetUC.GetTweets", nil)
	defer span.End()

	self, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, httpErrors.NewUnauthorizedError(errors.WithMessage(err, "tweetUC.GetTweets.GetUserFromCtx"))
	}

	return u.tweetRepo.GetTweets(ctx, self.UserID, pq)
}

// Get tweets by user id
func (u *tweetUC) GetTweetsByUserID(ctx context.Context, userID uuid.UUID, pq *utils.PaginationQuery) (*models.TweetsList, error) {
	ctx, span := tracer.NewSpan(ctx, "tweetUC.GetTweetsByUserID", nil)
	defer span.End()

	self, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, httpErrors.NewUnauthorizedError(errors.WithMessage(err, "tweetUC.GetTweetByUserID.GetUserFromCtx"))
	}

	return u.tweetRepo.GetTweetsByUserID(ctx, self.UserID, userID, pq)
}

// Get reply tweets by tweet id
func (u *tweetUC) GetReplyTweets(ctx context.Context, tweetID uint64, pq *utils.PaginationQuery) (*models.TweetsList, error) {
	ctx, span := tracer.NewSpan(ctx, "tweetUC.GetReplyTweets", nil)
	defer span.End()

	self, err := utils.GetUserFromCtx(ctx)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, httpErrors.NewUnauthorizedError(errors.WithMessage(err, "tweetUC.GetReplyTweets.GetUserFromCtx"))
	}

	return u.tweetRepo.GetReplyTweets(ctx, self.UserID, tweetID, pq)
}

// Delete tweet
func (u *tweetUC) Delete(ctx context.Context, tweetID uint64) error {
	ctx, span := tracer.NewSpan(ctx, "tweetUC.Delete", nil)
	defer span.End()

	if err := u.tweetRepo.Delete(ctx, tweetID); err != nil {
		tracer.AddSpanError(span, err)
		return err
	}

	return nil
}
