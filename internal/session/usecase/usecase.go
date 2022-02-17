package usecase

import (
	"context"

	"github.com/JamesHsu333/go-twitter/config"
	"github.com/JamesHsu333/go-twitter/internal/models"
	"github.com/JamesHsu333/go-twitter/internal/session"
	"github.com/JamesHsu333/go-twitter/pkg/tracer"
)

// Session use case
type SessionUC struct {
	sessionRepo session.SessRepository
	cfg         *config.Config
}

// New session use case constructor
func NewSessionUseCase(sessionRepo session.SessRepository, cfg *config.Config) session.UCSession {
	return &SessionUC{sessionRepo: sessionRepo, cfg: cfg}
}

// Create new session
func (u *SessionUC) CreateSession(ctx context.Context, session *models.Session, expire int) (string, error) {
	ctx, span := tracer.NewSpan(ctx, "sessionUC.CreateSession", nil)
	defer span.End()

	return u.sessionRepo.CreateSession(ctx, session, expire)
}

// Delete session by id
func (u *SessionUC) DeleteByID(ctx context.Context, sessionID string) error {
	ctx, span := tracer.NewSpan(ctx, "sessionUC.DeleteByID", nil)
	defer span.End()

	return u.sessionRepo.DeleteByID(ctx, sessionID)
}

// get session by id
func (u *SessionUC) GetSessionByID(ctx context.Context, sessionID string) (*models.Session, error) {
	ctx, span := tracer.NewSpan(ctx, "sessionUC.GetSessionByID", nil)
	defer span.End()

	return u.sessionRepo.GetSessionByID(ctx, sessionID)
}
