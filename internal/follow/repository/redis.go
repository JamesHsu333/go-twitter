package repository

import (
	"context"
	"time"

	"github.com/JamesHsu333/go-twitter/internal/follow"
	"github.com/JamesHsu333/go-twitter/pkg/tracer"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

// Follow redis repository
type followRedisRepo struct {
	redisClient *redis.Client
}

// Follow redis repository constructor
func NewFollowRedisRepo(redisClient *redis.Client) follow.RedisRepository {
	return &followRedisRepo{redisClient: redisClient}
}

func (a *followRedisRepo) GetFollowCtx(ctx context.Context, key string) (*bool, error) {
	ctx, span := tracer.NewSpan(ctx, "followRedisRepo.GetFollowCtx", nil)
	defer span.End()

	isFollowing, err := a.redisClient.Get(ctx, key).Bool()
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "followRedisRepo.GetFollowCtx.redisClient.Get")
	}
	return &isFollowing, nil
}

func (a *followRedisRepo) SetFollowCtx(ctx context.Context, key string, seconds int, isFollowing bool) error {
	ctx, span := tracer.NewSpan(ctx, "followRedisRepo.SetFollowCtx", nil)
	defer span.End()

	if err := a.redisClient.Set(ctx, key, isFollowing, time.Second*time.Duration(seconds)).Err(); err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "followRedisRepo.SetFollowCtx.redisClient.Set")
	}
	return nil
}

func (a *followRedisRepo) GetFollowCountCtx(ctx context.Context, key string) (*int64, error) {
	ctx, span := tracer.NewSpan(ctx, "followRedisRepo.GetFollowCountCtx", nil)
	defer span.End()

	followers, err := a.redisClient.Get(ctx, key).Int64()
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "followRedisRepo.GetFollowersCtx.redisClient.Get")
	}
	return &followers, nil
}

func (a *followRedisRepo) SetFollowCountCtx(ctx context.Context, key string, seconds int, value int64) error {
	ctx, span := tracer.NewSpan(ctx, "followRedisRepo.SetFollowCountCtx", nil)
	defer span.End()

	if err := a.redisClient.Set(ctx, key, value, time.Second*time.Duration(seconds)).Err(); err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "followRedisRepo.SetFollowCountCtx.redisClient.Set")
	}
	return nil
}

func (a *followRedisRepo) DeleteFollowCtx(ctx context.Context, key string) error {
	ctx, span := tracer.NewSpan(ctx, "followRedisRepo.DeleteFollowCtx", nil)
	defer span.End()

	if err := a.redisClient.Del(ctx, key).Err(); err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "followRedisRepo.DeleteFollowCtx.redisClient.Del")
	}
	return nil
}
