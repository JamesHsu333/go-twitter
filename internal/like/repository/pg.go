package repository

import (
	"context"
	"database/sql"

	"github.com/JamesHsu333/go-twitter/internal/like"
	"github.com/JamesHsu333/go-twitter/internal/models"
	"github.com/JamesHsu333/go-twitter/pkg/tracer"
	"github.com/JamesHsu333/go-twitter/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type likeRepo struct {
	db *sqlx.DB
}

func NewLikeRepository(db *sqlx.DB) like.Repository {
	return &likeRepo{db: db}
}

func (r *likeRepo) Like(ctx context.Context, userID uuid.UUID, tweetID uint64) error {
	ctx, span := tracer.NewSpan(ctx, "likeRepo.Like", nil)
	defer span.End()

	result, err := r.db.ExecContext(ctx, likeQuery, userID.String(), tweetID)
	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.WithMessage(err, "likeRepo.Like.ExecContext")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.WithMessage(err, "likeRepo.Like.RowsAffected")
	}
	if rowsAffected == 0 {
		tracer.AddSpanError(span, sql.ErrNoRows)
		return errors.Wrap(sql.ErrNoRows, "likeRepo.Like.rowsAffected")
	}

	return nil
}

func (r *likeRepo) GetLikedTweets(ctx context.Context, userID uuid.UUID, pq *utils.PaginationQuery) (*models.TweetsList, error) {
	ctx, span := tracer.NewSpan(ctx, "likeRepo.GetLikedTweets", nil)
	defer span.End()

	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotalLikedTweets, userID.String()); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "likeRepo.GetLikedTweets.GetContext.getTotal")
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
	if err := r.db.SelectContext(ctx, &tweets, getLikedTweets, userID.String(), pq.GetOrderBy(), pq.GetOffset(), pq.GetLimit()); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "likeRepo.GetLikedTweets.SelectContext")
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

func (r *likeRepo) GetLikedUsers(ctx context.Context, selfID uuid.UUID, tweetID uint64, pq *utils.PaginationQuery) (*models.UsersList, error) {
	ctx, span := tracer.NewSpan(ctx, "likeRepo.GetLikedUsers", nil)
	defer span.End()

	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotalLikedUsers, tweetID); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "likeRepo.GetLikedUsers.GetContext.getTotal")
	}

	if totalCount == 0 {
		return &models.UsersList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
			Page:       pq.GetPage(),
			Size:       pq.GetSize(),
			HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
			Users:      make([]*models.User, 0),
		}, nil
	}

	var users = make([]*models.User, 0, pq.GetSize())
	if err := r.db.SelectContext(ctx, &users, getLikedUsers, selfID.String(), tweetID, pq.GetOrderBy(), pq.GetOffset(), pq.GetLimit()); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "likeRepo.GetLikedTweets.SelectContext")
	}

	return &models.UsersList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pq.GetSize()),
		Page:       pq.GetPage(),
		Size:       pq.GetSize(),
		HasMore:    utils.GetHasMore(pq.GetPage(), totalCount, pq.GetSize()),
		Users:      users,
	}, nil
}

func (r *likeRepo) Delete(ctx context.Context, userID uuid.UUID, tweetID uint64) error {
	ctx, span := tracer.NewSpan(ctx, "likeRepo.Delete", nil)
	defer span.End()

	result, err := r.db.ExecContext(ctx, deleteQuery, userID.String(), tweetID)
	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.WithMessage(err, "likeRepo.Delete.ExecContext")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "likeRepo.Delete.RowsAffected")
	}
	if rowsAffected == 0 {
		tracer.AddSpanError(span, sql.ErrNoRows)
		return errors.Wrap(sql.ErrNoRows, "likeRepo.Delete.rowsAffected")
	}
	return nil
}
