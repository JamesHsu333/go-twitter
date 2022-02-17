package middleware

import (
	"github.com/JamesHsu333/go-twitter/config"
	"github.com/JamesHsu333/go-twitter/internal/session"
	"github.com/JamesHsu333/go-twitter/internal/user"
	"github.com/JamesHsu333/go-twitter/pkg/logger"
)

// Middleware manager
type MiddlewareManager struct {
	sessUC  session.UCSession
	userUC  user.UseCase
	cfg     *config.Config
	origins []string
	logger  logger.Logger
}

// Middleware manager constructor
func NewMiddlewareManager(sessUC session.UCSession, userUC user.UseCase, cfg *config.Config, origins []string, logger logger.Logger) *MiddlewareManager {
	return &MiddlewareManager{sessUC: sessUC, userUC: userUC, cfg: cfg, origins: origins, logger: logger}
}
