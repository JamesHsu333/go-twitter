package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/JamesHsu333/go-twitter/config"
	"github.com/JamesHsu333/go-twitter/internal/models"
	"github.com/JamesHsu333/go-twitter/internal/session"
	"github.com/JamesHsu333/go-twitter/pkg/tracer"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	basePrefix = "api-session:"
)

// Session repository
type SessionRepo struct {
	redisClient *redis.Client
	basePrefix  string
	cfg         *config.Config
}

// Session repository init
func NewSessionRepository(redisClient *redis.Client, cfg *config.Config) session.SessRepository {
	return &SessionRepo{redisClient: redisClient, basePrefix: basePrefix, cfg: cfg}
}

// Create session in redis
func (s *SessionRepo) CreateSession(ctx context.Context, sess *models.Session, expire int) (string, error) {
	ctx, span := tracer.NewSpan(ctx, "SessionRepo.CreateSession", nil)
	defer span.End()

	sess.SessionID = uuid.New().String()
	sessionKey := s.createKey(sess.SessionID)

	sessBytes, err := json.Marshal(&sess)
	if err != nil {
		tracer.AddSpanError(span, err)
		return "", errors.WithMessage(err, "SessionRepo.CreateSession.json.Marshal")
	}
	if err = s.redisClient.Set(ctx, sessionKey, sessBytes, time.Second*time.Duration(expire)).Err(); err != nil {
		tracer.AddSpanError(span, err)
		return "", errors.Wrap(err, "SessionRepo.CreateSession.redisClient.Set")
	}
	return sessionKey, nil
}

// Get session by id
func (s *SessionRepo) GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error) {
	ctx, span := tracer.NewSpan(ctx, "SessionRepo.GetSessionByID", nil)
	defer span.End()

	sessBytes, err := s.redisClient.Get(ctx, sessionID).Bytes()
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "SessionRep.GetSessionByID.redisClient.Get")
	}

	sess := &models.Session{}
	if err = json.Unmarshal(sessBytes, &sess); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "SessionRepo.GetSessionByID.json.Unmarshal")
	}
	return sess, nil
}

// Delete session by id
func (s *SessionRepo) DeleteByID(ctx context.Context, sessionID string) error {
	ctx, span := tracer.NewSpan(ctx, "SessionRepo.DeleteByID", nil)
	defer span.End()

	if err := s.redisClient.Del(ctx, sessionID).Err(); err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "sessionRepo.DeleteByID")
	}
	return nil
}

func (s *SessionRepo) createKey(sessionID string) string {
	return fmt.Sprintf("%s: %s", s.basePrefix, sessionID)
}
