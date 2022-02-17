package repository

import (
	"context"
	"database/sql"

	"github.com/JamesHsu333/go-twitter/internal/models"
	"github.com/JamesHsu333/go-twitter/internal/tweet"
	"github.com/JamesHsu333/go-twitter/pkg/tracer"
	"github.com/JamesHsu333/go-twitter/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Tweet Repository
type tweetRepo struct {
	db *sqlx.DB
}

func NewTweetRepository(db *sqlx.DB) tweet.Repository {
	return &tweetRepo{db: db}
}

// Create new user
func (r *tweetRepo) Create(ctx context.Context, tweet *models.Tweet) (*models.Tweet, error) {
	ctx, span := tracer.NewSpan(ctx, "tweetRepo.Create", nil)
	defer span.End()

	t := &models.Tweet{}
	if err := r.db.QueryRowxContext(ctx, createTweetQuery, &tweet.UserID, &tweet.Text, &tweet.Image).StructScan(t); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "tweetRepo.Create.StructScan")
	}
	return t, nil
}

func (r *tweetRepo) CreateReply(ctx context.Context, tweetID uint64, replyID uint64) error {
	ctx, span := tracer.NewSpan(ctx, "tweetRepo.CreateReply", nil)
	defer span.End()

	result, err := r.db.ExecContext(ctx, createReplyQuery, tweetID, replyID)
	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "tweetRepo.CreateReply.ExecContext")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.WithMessage(err, "tweetRepo.CreateReply.RowsAffected")
	}
	if rowsAffected == 0 {
		tracer.AddSpanError(span, sql.ErrNoRows)
		return errors.Wrap(sql.ErrNoRows, "tweetRepo.CreateReply.rowsAffected")
	}

	return nil
}

func (r *tweetRepo) CheckTweetExist(ctx context.Context, tweetID uint64) error {
	ctx, span := tracer.NewSpan(ctx, "tweetRepo.CheckTweetExist", nil)
	defer span.End()

	var tweetExists bool
	if err := r.db.QueryRowxContext(ctx, checkTweetExist, tweetID).Scan(&tweetExists); err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "tweetRepo.CheckTweetExist.StructScan")
	}

	if !tweetExists {
		return errors.Wrap(errors.New("Tweet does not exist"), "tweetRepo.CreateTweetExist.StructScan")
	}

	return nil
}

func (r *tweetRepo) GetTweetByID(ctx context.Context, selfID uuid.UUID, tweetID uint64) (*models.TweetWithUser, error) {
	ctx, span := tracer.NewSpan(ctx, "tweetRepo.GetTweetByID", nil)
	defer span.End()

	t := &models.TweetWithUser{}
	if err := r.db.QueryRowxContext(ctx, getTweetQuery, selfID.String(), tweetID).StructScan(t); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "tweetRepo.GetTweetByID.StructScan")
	}
	return t, nil
}

func (r *tweetRepo) GetTweets(ctx context.Context, selfID uuid.UUID, pq *utils.PaginationQuery) (*models.TweetsList, error) {
	ctx, span := tracer.NewSpan(ctx, "tweetRepo.GetTweets", nil)
	defer span.End()

	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotal); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "tweetRepo.GetTweets.GetContext.getTotal")
	}

	if totalCount == 0 {
		return &models.TweetsList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
			Page:       pq.GetPage(),
			Size:       pq.GetSize(),
			HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			Tweets:     make([]*models.TweetWithUser, 0),
		}, nil
	}

	var tweets = make([]*models.TweetWithUser, 0, pq.GetSize())
	if err := r.db.SelectContext(ctx, &tweets, getTweets, selfID.String(), pq.GetOrderBy(), pq.GetOffset(), pq.GetLimit()); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "tweetRepo.GetTweets.SelectContext")
	}

	return &models.TweetsList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		Tweets:     tweets,
	}, nil
}

func (r *tweetRepo) GetTweetsByUserID(ctx context.Context, selfID uuid.UUID, userID uuid.UUID, pq *utils.PaginationQuery) (*models.TweetsList, error) {
	ctx, span := tracer.NewSpan(ctx, "tweetRepo.GetTweetsByUserID", nil)
	defer span.End()

	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotalByUserID, userID); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "tweetRepo.GetTweetsByUserID.GetContext.getTotalByUserID")
	}

	if totalCount == 0 {
		return &models.TweetsList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
			Page:       pq.GetPage(),
			Size:       pq.GetSize(),
			HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			Tweets:     make([]*models.TweetWithUser, 0),
		}, nil
	}

	var tweets = make([]*models.TweetWithUser, 0, pq.GetSize())
	if err := r.db.SelectContext(ctx, &tweets, getTweetsByUserID, selfID.String(), userID.String(), pq.GetOrderBy(), pq.GetOffset(), pq.GetLimit()); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "tweetRepo.GetTweetsByUserID.SelectContext")
	}

	return &models.TweetsList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		Tweets:     tweets,
	}, nil
}

func (r *tweetRepo) GetReplyTweets(ctx context.Context, selfID uuid.UUID, tweetID uint64, pq *utils.PaginationQuery) (*models.TweetsList, error) {
	ctx, span := tracer.NewSpan(ctx, "tweetRepo.GetReplyTweets", nil)
	defer span.End()

	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getReplysTotal, tweetID); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "tweetRepo.GetReplyTweets.SelectContext.getReplysTotal")
	}

	if totalCount == 0 {
		return &models.TweetsList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
			Page:       pq.GetPage(),
			Size:       pq.GetSize(),
			HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			Tweets:     make([]*models.TweetWithUser, 0),
		}, nil
	}

	var tweets = make([]*models.TweetWithUser, 0, pq.GetSize())
	if err := r.db.SelectContext(ctx, &tweets, getReplyTweetsByID, selfID.String(), tweetID, pq.GetOrderBy(), pq.GetOffset(), pq.GetLimit()); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "tweetRepo.GetReplyTweets.SelectContext")
	}

	return &models.TweetsList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		Tweets:     tweets,
	}, nil
}

func (r *tweetRepo) Delete(ctx context.Context, tweetID uint64) error {
	ctx, span := tracer.NewSpan(ctx, "tweetRepo.Delete", nil)
	defer span.End()

	result, err := r.db.ExecContext(ctx, deleteTweetQuery, tweetID)
	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.WithMessage(err, "tweetRepo.Delete.ExecContext")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "tweetRepo.Delete.RowsAffected")
	}
	if rowsAffected == 0 {
		tracer.AddSpanError(span, sql.ErrNoRows)
		return errors.Wrap(sql.ErrNoRows, "tweetRepo.Delete.rowsAffected")
	}

	return nil
}
