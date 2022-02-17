package usecase

import (
	"context"
	"fmt"
	"net/http"

	"github.com/JamesHsu333/go-twitter/config"
	"github.com/JamesHsu333/go-twitter/internal/follow"
	"github.com/JamesHsu333/go-twitter/internal/models"
	"github.com/JamesHsu333/go-twitter/internal/user"
	"github.com/JamesHsu333/go-twitter/pkg/httpErrors"
	"github.com/JamesHsu333/go-twitter/pkg/logger"
	"github.com/JamesHsu333/go-twitter/pkg/tracer"
	"github.com/JamesHsu333/go-twitter/pkg/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	basePrefix    = "api-twitter:"
	cacheDuration = 3600
)

// User UseCase
type userUC struct {
	cfg             *config.Config
	userRepo        user.Repository
	redisRepo       user.RedisRepository
	followRedisRepo follow.RedisRepository
	logger          logger.Logger
}

// User UseCase constructor
func NewUserUseCase(cfg *config.Config, userRepo user.Repository, redisRepo user.RedisRepository, followRedisRepo follow.RedisRepository, log logger.Logger) user.UseCase {
	return &userUC{cfg: cfg, userRepo: userRepo, redisRepo: redisRepo, followRedisRepo: followRedisRepo, logger: log}
}

// Create new user
func (u *userUC) Register(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	ctx, span := tracer.NewSpan(ctx, "userUC.Register", nil)
	defer span.End()

	existsUser, err := u.userRepo.FindByEmail(ctx, user)
	if existsUser != nil || err == nil {
		return nil, httpErrors.NewRestErrorWithMessage(http.StatusBadRequest, httpErrors.ErrEmailAlreadyExists, nil)
	}

	if err = user.PrepareCreate(); err != nil {
		tracer.AddSpanError(span, err)
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "userUC.Register.PrepareCreate"))
	}

	createdUser, err := u.userRepo.Register(ctx, user)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}
	createdUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(createdUser, u.cfg)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "userUC.Register.GenerateJWTToken"))
	}

	return &models.UserWithToken{
		User:  createdUser,
		Token: token,
	}, nil
}

// Update existing user
func (u *userUC) Update(ctx context.Context, user *models.User) (*models.User, error) {
	ctx, span := tracer.NewSpan(ctx, "userUC.Update", nil)
	defer span.End()

	if err := user.PrepareUpdate(); err != nil {
		tracer.AddSpanError(span, err)
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "userUC.Update.PrepareUpdate"))
	}

	updatedUser, err := u.userRepo.Update(ctx, user)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	updatedUser.SanitizePassword()

	if err = u.redisRepo.DeleteUserCtx(ctx, u.GenerateUserKey(user.UserID.String())); err != nil {
		tracer.AddSpanError(span, err)
		u.logger.Errorf("userUC.Update.DeleteUserCtx: %s", err)
	}

	updatedUser.SanitizePassword()

	return updatedUser, nil
}

// Delete new user
func (u *userUC) Delete(ctx context.Context, userID uuid.UUID) error {
	ctx, span := tracer.NewSpan(ctx, "userUC.Delete", nil)
	defer span.End()

	if err := u.userRepo.Delete(ctx, userID); err != nil {
		tracer.AddSpanError(span, err)
		return err
	}

	if err := u.redisRepo.DeleteUserCtx(ctx, u.GenerateUserKey(userID.String())); err != nil {
		tracer.AddSpanError(span, err)
		u.logger.Errorf("userUC.Delete.DeleteUserCtx: %s", err)
	}

	return nil
}

// Get user by id
func (u *userUC) GetByID(ctx context.Context, selfID uuid.UUID, userID uuid.UUID) (*models.User, error) {
	ctx, span := tracer.NewSpan(ctx, "userUC.GetByID", nil)
	defer span.End()

	user, err := u.userRepo.GetByID(ctx, selfID, userID)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	if err = u.followRedisRepo.SetFollowCtx(ctx, u.GenerateFollowKey(selfID.String(), userID.String()), cacheDuration, *user.IsFollowing); err != nil {
		tracer.AddSpanError(span, err)
		u.logger.Errorf("userUC.GetByID.SetFollowCtx: %v", err)
	}

	if err = u.followRedisRepo.SetFollowCountCtx(ctx, u.GenerateFollowCountKey("followers of", userID.String()), cacheDuration, *user.Followers); err != nil {
		tracer.AddSpanError(span, err)
		u.logger.Errorf("userUC.GetByID.SetFollowCtx: %v", err)
	}

	if err = u.followRedisRepo.SetFollowCountCtx(ctx, u.GenerateFollowCountKey("following of", userID.String()), cacheDuration, *user.Following); err != nil {
		tracer.AddSpanError(span, err)
		u.logger.Errorf("userUC.GetByID.SetFollowCtx: %v", err)
	}

	if err = u.redisRepo.SetUserCtx(ctx, u.GenerateUserKey(userID.String()), cacheDuration, *user); err != nil {
		tracer.AddSpanError(span, err)
		u.logger.Errorf("userUC.GetByID.SetUserCtx: %v", err)
	}

	user.SanitizePassword()

	return user, nil
}

func (u *userUC) GetCacheByID(ctx context.Context, selfID uuid.UUID, userID uuid.UUID) (*models.User, error) {
	ctx, span := tracer.NewSpan(ctx, "userUC.GetCacheByID", nil)
	defer span.End()

	cachedFollowerCount, err := u.followRedisRepo.GetFollowCountCtx(ctx, u.GenerateFollowCountKey("followers of", userID.String()))
	if err != nil {
		tracer.AddSpanError(span, err)
		u.logger.Infof("userUC.GetCacheByID.GetFollowCountCtx: %v", err)
		return nil, err
	}

	cachedFollowingCount, err := u.followRedisRepo.GetFollowCountCtx(ctx, u.GenerateFollowCountKey("following of", userID.String()))
	if err != nil {
		tracer.AddSpanError(span, err)
		u.logger.Infof("userUC.GetCacheByID.GetFollowCountCtx: %v", err)
		return nil, err
	}

	cachedFollow, err := u.followRedisRepo.GetFollowCtx(ctx, u.GenerateFollowKey(selfID.String(), userID.String()))
	if err != nil {
		tracer.AddSpanError(span, err)
		u.logger.Infof("userUC.GetCacheByID.GetFollowCtx: %v", err)
		return nil, err
	}

	cachedUser, err := u.redisRepo.GetByIDCtx(ctx, u.GenerateUserKey(userID.String()))
	if err != nil {
		tracer.AddSpanError(span, err)
		u.logger.Infof("userUC.GetCacheByID.GetByIDCtx: %v", err)
		return nil, err
	}

	cachedUser.IsFollowing = cachedFollow
	cachedUser.Followers = cachedFollowerCount
	cachedUser.Following = cachedFollowingCount
	return cachedUser, nil
}

// Get user by user name
func (u *userUC) GetByUserName(ctx context.Context, selfID uuid.UUID, name string) (*models.User, error) {
	ctx, span := tracer.NewSpan(ctx, "userUC.GetByUserName", nil)
	defer span.End()

	return u.userRepo.GetByUserName(ctx, selfID, name)
}

// Find users by name
func (u *userUC) FindByName(ctx context.Context, selfID uuid.UUID, name string, query *utils.PaginationQuery) (*models.UsersList, error) {
	ctx, span := tracer.NewSpan(ctx, "userUC.FindByName", nil)
	defer span.End()

	return u.userRepo.FindByName(ctx, selfID, name, query)
}

// Get users with pagination
func (u *userUC) GetUsers(ctx context.Context, selfID uuid.UUID, pq *utils.PaginationQuery) (*models.UsersList, error) {
	ctx, span := tracer.NewSpan(ctx, "userUC.GetUsers", nil)
	defer span.End()

	return u.userRepo.GetUsers(ctx, selfID, pq)
}

// Login user, returns user model with jwt token
func (u *userUC) Login(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	ctx, span := tracer.NewSpan(ctx, "userUC.Login", nil)
	defer span.End()

	foundUser, err := u.userRepo.FindByEmail(ctx, user)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	if err = foundUser.ComparePasswords(user.Password); err != nil {
		tracer.AddSpanError(span, err)
		return nil, httpErrors.NewUnauthorizedError(errors.Wrap(err, "userUC.Login.ComparePasswords"))
	}

	foundUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(foundUser, u.cfg)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, httpErrors.NewInternalServerError(errors.Wrap(err, "userUC.Login.GenerateJWTToken"))
	}

	return &models.UserWithToken{
		User:  foundUser,
		Token: token,
	}, nil
}

// Update user role
func (u *userUC) UpdateRole(ctx context.Context, user *models.User) (*models.User, error) {
	ctx, span := tracer.NewSpan(ctx, "userUC.UpdateRole", nil)
	defer span.End()

	if err := user.PrepareUpdate(); err != nil {
		tracer.AddSpanError(span, err)
		return nil, httpErrors.NewBadRequestError(errors.Wrap(err, "userUC.UpdateRole.PrepareUpdate"))
	}

	updatedUser, err := u.userRepo.UpdateRole(ctx, user)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, err
	}

	updatedUser.SanitizePassword()

	if err = u.redisRepo.DeleteUserCtx(ctx, u.GenerateUserKey(user.UserID.String())); err != nil {
		tracer.AddSpanError(span, err)
		u.logger.Errorf("userUC.UpdateRole.DeleteUserCtx: %s", err)
	}

	updatedUser.SanitizePassword()

	return updatedUser, nil
}

func (u *userUC) GenerateUserKey(userID string) string {
	return fmt.Sprintf("%s: %s", basePrefix, userID)
}

func (u *userUC) GenerateFollowKey(follower string, following string) string {
	return fmt.Sprintf("%s: %s+%s", basePrefix, follower, following)
}

func (u *userUC) GenerateFollowCountKey(follow string, user string) string {
	return fmt.Sprintf("%s: %s %s", basePrefix, follow, user)
}
