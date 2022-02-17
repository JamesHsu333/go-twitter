package repository

import (
	"context"
	"database/sql"

	"github.com/JamesHsu333/go-twitter/internal/follow"
	"github.com/JamesHsu333/go-twitter/internal/models"
	"github.com/JamesHsu333/go-twitter/pkg/tracer"
	"github.com/JamesHsu333/go-twitter/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Follow repository
type followRepo struct {
	db *sqlx.DB
}

func NewFollowRepository(db *sqlx.DB) follow.Repository {
	return &followRepo{db: db}
}

// Follow an user
func (r *followRepo) Follow(ctx context.Context, follower uuid.UUID, following uuid.UUID) error {
	ctx, span := tracer.NewSpan(ctx, "followRepo.Follow", nil)
	defer span.End()

	result, err := r.db.ExecContext(ctx, followQuery, follower, following)
	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.WithMessage(err, "followRepo.Follow.ExecContext")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "followRepo.Follow.RowsAffected")
	}
	if rowsAffected == 0 {
		tracer.AddSpanError(span, sql.ErrNoRows)
		return errors.Wrap(sql.ErrNoRows, "followRepo.Follow.rowsAffected")
	}

	return nil
}

func (r *followRepo) GetFollowers(ctx context.Context, selfID uuid.UUID, userID uuid.UUID, pq *utils.PaginationQuery) (*models.UsersList, error) {
	ctx, span := tracer.NewSpan(ctx, "followRepo.GetFollowers", nil)
	defer span.End()

	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getFollowersTotal, userID.String()); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "followRepo.GetFollowers.GetContext.getTotal")
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
	if err := r.db.SelectContext(ctx, &users, getFollowers, selfID.String(), userID.String(), pq.GetOrderBy(), pq.GetOffset(), pq.GetLimit()); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "followRepo.GetFollowers.SelectContext")
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

func (r *followRepo) GetFollowing(ctx context.Context, selfID uuid.UUID, userID uuid.UUID, pq *utils.PaginationQuery) (*models.UsersList, error) {
	ctx, span := tracer.NewSpan(ctx, "followRepo.GetFollowing", nil)
	defer span.End()

	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getFollowingTotal, userID.String()); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "followRepo.GetFollowing.GetContext.getTotal")
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
	if err := r.db.SelectContext(ctx, &users, getFollowing, selfID.String(), userID.String(), pq.GetOrderBy(), pq.GetOffset(), pq.GetLimit()); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "followRepo.GetFollowing.SelectContext")
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

func (r *followRepo) Delete(ctx context.Context, follower uuid.UUID, following uuid.UUID) error {
	ctx, span := tracer.NewSpan(ctx, "followRepo.Delete", nil)
	defer span.End()

	result, err := r.db.ExecContext(ctx, deleteQuery, follower, following)
	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.WithMessage(err, "followRepo.Delete.ExecContext")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "followRepo.Delete.RowsAffected")
	}
	if rowsAffected == 0 {
		tracer.AddSpanError(span, sql.ErrNoRows)
		return errors.Wrap(sql.ErrNoRows, "followRepo.Delete.rowsAffected")
	}

	return nil
}
