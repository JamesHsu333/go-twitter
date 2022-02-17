package usecase

import (
	"context"
	"fmt"

	"github.com/JamesHsu333/go-twitter/config"
	"github.com/JamesHsu333/go-twitter/internal/follow"
	"github.com/JamesHsu333/go-twitter/internal/models"
	"github.com/JamesHsu333/go-twitter/pkg/logger"
	"github.com/JamesHsu333/go-twitter/pkg/tracer"
	"github.com/JamesHsu333/go-twitter/pkg/utils"
	"github.com/google/uuid"
)

const (
	basePrefix = "api-twitter:"
)

// Follow Usecase
type followUC struct {
	cfg             *config.Config
	followRepo      follow.Repository
	followRedisRepo follow.RedisRepository
	logger          logger.Logger
}

// New Usecase
func NewFollowUseCase(cfg *config.Config, followRepo follow.Repository, followRedisRepo follow.RedisRepository, logger logger.Logger) follow.UseCase {
	return &followUC{cfg: cfg, followRepo: followRepo, followRedisRepo: followRedisRepo, logger: logger}
}

func (u *followUC) Follow(ctx context.Context, follower uuid.UUID, following uuid.UUID) error {
	ctx, span := tracer.NewSpan(ctx, "followUC.Follow", nil)
	defer span.End()

	if err := u.followRepo.Follow(ctx, follower, following); err != nil {
		tracer.AddSpanError(span, err)
		u.logger.Errorf("followUC.Follow.Follow: %v", err)
		return err
	}

	if err := u.followRedisRepo.DeleteFollowCtx(ctx, u.generateFollowKey("following of", follower.String())); err != nil {
		tracer.AddSpanError(span, err)
		u.logger.Infof("followUC.Follow.DeleteFollowCtx: %v", err)
	}

	if err := u.followRedisRepo.DeleteFollowCtx(ctx, u.generateFollowKey("followers of", following.String())); err != nil {
		tracer.AddSpanError(span, err)
		u.logger.Infof("followUC.Follow.DeleteFollowCtx: %v", err)
	}

	return nil
}

func (u *followUC) GetFollowers(ctx context.Context, selfID uuid.UUID, userID uuid.UUID, pq *utils.PaginationQuery) (*models.UsersList, error) {
	ctx, span := tracer.NewSpan(ctx, "followUC.GetFollowers", nil)
	defer span.End()

	return u.followRepo.GetFollowers(ctx, selfID, userID, pq)
}

func (u *followUC) GetFollowing(ctx context.Context, selfID uuid.UUID, userID uuid.UUID, pq *utils.PaginationQuery) (*models.UsersList, error) {
	ctx, span := tracer.NewSpan(ctx, "followUC.GetFollowing", nil)
	defer span.End()

	return u.followRepo.GetFollowing(ctx, selfID, userID, pq)
}

func (u *followUC) Delete(ctx context.Context, follower uuid.UUID, following uuid.UUID) error {
	ctx, span := tracer.NewSpan(ctx, "followUC.Delete", nil)
	defer span.End()

	if err := u.followRepo.Delete(ctx, follower, following); err != nil {
		tracer.AddSpanError(span, err)
		u.logger.Infof("followUC.Delete.Delete: %v", err)
		return err
	}

	if err := u.followRedisRepo.DeleteFollowCtx(ctx, u.generateFollowKey("following of", follower.String())); err != nil {
		tracer.AddSpanError(span, err)
		u.logger.Infof("followUC.Delete.DeleteFollowCtx: %v", err)
	}

	if err := u.followRedisRepo.DeleteFollowCtx(ctx, u.generateFollowKey("followers of", following.String())); err != nil {
		tracer.AddSpanError(span, err)
		u.logger.Infof("followUC.Delete.DeleteFollowCtx: %v", err)
	}

	return nil
}

func (u *followUC) generateFollowKey(follow string, user string) string {
	return fmt.Sprintf("%s: %s %s", basePrefix, follow, user)
}
