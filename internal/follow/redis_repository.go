package follow

import (
	"context"
)

// Follow Redis repository interface
type RedisRepository interface {
	GetFollowCtx(ctx context.Context, key string) (*bool, error)
	SetFollowCtx(ctx context.Context, key string, seconds int, isFollowing bool) error
	GetFollowCountCtx(ctx context.Context, key string) (*int64, error)
	SetFollowCountCtx(ctx context.Context, key string, seconds int, value int64) error
	DeleteFollowCtx(ctx context.Context, key string) error
}
