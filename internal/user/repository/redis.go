package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/JamesHsu333/go-twitter/internal/models"
	"github.com/JamesHsu333/go-twitter/internal/user"
	"github.com/JamesHsu333/go-twitter/pkg/tracer"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

// User redis repository
type userRedisRepo struct {
	redisClient *redis.Client
}

// User redis repository constructor
func NewUserRedisRepo(redisClient *redis.Client) user.RedisRepository {
	return &userRedisRepo{redisClient: redisClient}
}

// Get user by id
func (a *userRedisRepo) GetByIDCtx(ctx context.Context, key string) (*models.User, error) {
	ctx, span := tracer.NewSpan(ctx, "userRedisRepo.GetByIDCtx", nil)
	defer span.End()

	userBytes, err := a.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "userRedisRepo.GetByIDCtx.redisClient.Get")
	}
	user := &models.User{}
	if err = json.Unmarshal(userBytes, user); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "userRedisRepo.GetByIDCtx.json.Unmarshal")
	}
	return user, nil
}

// Cache user with duration in seconds
func (a *userRedisRepo) SetUserCtx(ctx context.Context, key string, seconds int, user models.User) error {
	ctx, span := tracer.NewSpan(ctx, "userRedisRepo.SetUserCtx", nil)
	defer span.End()

	user.SanitizeFollow()
	userBytes, err := json.Marshal(user)
	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "userRedisRepo.SetUserCtx.json.Unmarshal")
	}
	if err = a.redisClient.Set(ctx, key, userBytes, time.Second*time.Duration(seconds)).Err(); err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "userRedisRepo.SetUserCtx.redisClient.Set")
	}
	return nil
}

// Delete user by key
func (a *userRedisRepo) DeleteUserCtx(ctx context.Context, key string) error {
	ctx, span := tracer.NewSpan(ctx, "userRedisRepo.DeleteUserCtx", nil)
	defer span.End()

	if err := a.redisClient.Del(ctx, key).Err(); err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "userRedisRepo.DeleteUserCtx.redisClient.Del")
	}
	return nil
}
