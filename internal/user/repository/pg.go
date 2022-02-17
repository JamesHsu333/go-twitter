package repository

import (
	"context"
	"database/sql"

	"github.com/JamesHsu333/go-twitter/internal/models"
	"github.com/JamesHsu333/go-twitter/internal/user"
	"github.com/JamesHsu333/go-twitter/pkg/tracer"
	"github.com/JamesHsu333/go-twitter/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// User Repository
type userRepo struct {
	db *sqlx.DB
}

// User Repository constructor
func NewUserRepository(db *sqlx.DB) user.Repository {
	return &userRepo{db: db}
}

// Create new user
func (r *userRepo) Register(ctx context.Context, user *models.User) (*models.User, error) {
	ctx, span := tracer.NewSpan(ctx, "userRepo.Register", nil)
	defer span.End()

	u := &models.User{}
	if err := r.db.QueryRowxContext(ctx, createUserQuery, &user.UserName, &user.Name, &user.Email, &user.Password,
		&user.About, &user.Avatar, &user.Header, &user.PhoneNumber, &user.Country, &user.Gender, utils.ParseTimeFormat(user.Birthday),
	).StructScan(u); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "userRepo.Register.StructScan")
	}

	return u, nil
}

// Update existing user
func (r *userRepo) Update(ctx context.Context, user *models.User) (*models.User, error) {
	ctx, span := tracer.NewSpan(ctx, "userRepo.Update", nil)
	defer span.End()

	u := &models.User{}
	if err := r.db.GetContext(ctx, u, updateUserQuery, &user.UserName, &user.Name, &user.Email,
		&user.About, &user.Avatar, &user.Header, &user.PhoneNumber, &user.Country, &user.Gender,
		utils.ParseTimeFormat(user.Birthday), &user.UserID,
	); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "userRepo.Update.GetContext")
	}

	return u, nil
}

// Delete existing user
func (r *userRepo) Delete(ctx context.Context, userID uuid.UUID) error {
	ctx, span := tracer.NewSpan(ctx, "userRepo.Delete", nil)
	defer span.End()

	result, err := r.db.ExecContext(ctx, deleteUserQuery, userID)
	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.WithMessage(err, "userRepo Delete ExecContext")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "userRepo.Delete.RowsAffected")
	}
	if rowsAffected == 0 {
		return errors.Wrap(sql.ErrNoRows, "userRepo.Delete.rowsAffected")
	}

	return nil
}

// Get user by id
func (r *userRepo) GetByID(ctx context.Context, selfID uuid.UUID, userID uuid.UUID) (*models.User, error) {
	ctx, span := tracer.NewSpan(ctx, "userRepo.GetByID", nil)
	defer span.End()

	user := &models.User{}
	if err := r.db.QueryRowxContext(ctx, getUserQuery, selfID, userID).StructScan(user); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "userRepo.GetByID.QueryRowxContext")
	}
	return user, nil
}

// Get user by user name
func (r *userRepo) GetByUserName(ctx context.Context, selfID uuid.UUID, name string) (*models.User, error) {
	ctx, span := tracer.NewSpan(ctx, "userRepo.GetByUserName", nil)
	defer span.End()

	user := &models.User{}
	if err := r.db.QueryRowxContext(ctx, getUserByUserNameQuery, selfID, name).StructScan(user); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "userRepo.GetByUserName.QueryRowxContext")
	}
	return user, nil
}

// Find users by name
func (r *userRepo) FindByName(ctx context.Context, selfID uuid.UUID, name string, query *utils.PaginationQuery) (*models.UsersList, error) {
	ctx, span := tracer.NewSpan(ctx, "userRepo.FindByName", nil)
	defer span.End()

	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotalCount, name); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "userRepo.FindByName.GetContext.totalCount")
	}

	if totalCount == 0 {
		return &models.UsersList{
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, query.GetSize()),
			Page:       query.GetPage(),
			Size:       query.GetSize(),
			HasMore:    utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
			Users:      make([]*models.User, 0),
		}, nil
	}

	rows, err := r.db.QueryxContext(ctx, findUsers, selfID, name, query.GetOffset(), query.GetLimit())
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "userRepo.FindByName.QueryxContext")
	}
	defer rows.Close()

	var users = make([]*models.User, 0, query.GetSize())
	for rows.Next() {
		var user models.User
		if err = rows.StructScan(&user); err != nil {
			tracer.AddSpanError(span, err)
			return nil, errors.Wrap(err, "userRepo.FindByName.StructScan")
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "userRepo.FindByName.rows.Err")
	}

	return &models.UsersList{
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, query.GetSize()),
		Page:       query.GetPage(),
		Size:       query.GetSize(),
		HasMore:    utils.GetHasMore(query.GetPage(), totalCount, query.GetSize()),
		Users:      users,
	}, nil
}

// Get users with pagination
func (r *userRepo) GetUsers(ctx context.Context, selfID uuid.UUID, pq *utils.PaginationQuery) (*models.UsersList, error) {
	ctx, span := tracer.NewSpan(ctx, "userRepo.GetUsers", nil)
	defer span.End()

	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, getTotal); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "userRepo.GetUsers.GetContext.totalCount")
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
	if err := r.db.SelectContext(
		ctx,
		&users,
		getUsers,
		selfID,
		pq.GetOrderBy(),
		pq.GetOffset(),
		pq.GetLimit(),
	); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "userRepo.GetUsers.SelectContext")
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

// Find user by email
func (r *userRepo) FindByEmail(ctx context.Context, user *models.User) (*models.User, error) {
	ctx, span := tracer.NewSpan(ctx, "userRepo.FindByEmail", nil)
	defer span.End()

	foundUser := &models.User{}
	if err := r.db.QueryRowxContext(ctx, findUserByEmail, user.Email).StructScan(foundUser); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "userRepo.FindByEmail.QueryRowxContext")
	}
	return foundUser, nil
}

// Update existing user role
func (r *userRepo) UpdateRole(ctx context.Context, user *models.User) (*models.User, error) {
	ctx, span := tracer.NewSpan(ctx, "userRepo.UpdateRole", nil)
	defer span.End()

	u := &models.User{}
	if err := r.db.GetContext(ctx, u, updateUserRoleQuery, &user.Role, &user.UserID); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "userRepo.UpdateRole.GetContext")
	}
	return u, nil
}
