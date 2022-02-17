package server

import (
	"net/http"
	"strings"

	fileRepository "github.com/JamesHsu333/go-twitter/internal/file/repository"
	fileUseCase "github.com/JamesHsu333/go-twitter/internal/file/usecase"
	followRepository "github.com/JamesHsu333/go-twitter/internal/follow/repository"
	followUseCase "github.com/JamesHsu333/go-twitter/internal/follow/usecase"
	likeRepository "github.com/JamesHsu333/go-twitter/internal/like/repository"
	likeUseCase "github.com/JamesHsu333/go-twitter/internal/like/usecase"
	apiMiddlewares "github.com/JamesHsu333/go-twitter/internal/middleware"
	sessionRepository "github.com/JamesHsu333/go-twitter/internal/session/repository"
	"github.com/JamesHsu333/go-twitter/internal/session/usecase"
	tweetHttp "github.com/JamesHsu333/go-twitter/internal/tweet/delivery/http"
	tweetRepository "github.com/JamesHsu333/go-twitter/internal/tweet/repository"
	tweetUseCase "github.com/JamesHsu333/go-twitter/internal/tweet/usecase"
	userHttp "github.com/JamesHsu333/go-twitter/internal/user/delivery/http"
	userRepository "github.com/JamesHsu333/go-twitter/internal/user/repository"
	userUseCase "github.com/JamesHsu333/go-twitter/internal/user/usecase"
	"github.com/JamesHsu333/go-twitter/pkg/csrf"
	"github.com/JamesHsu333/go-twitter/pkg/metric"
	"github.com/JamesHsu333/go-twitter/pkg/utils"
	"golang.org/x/time/rate"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/swaggo/swag/example/basic/docs"
)

// Map Server Handlers
func (s *Server) MapHandlers(e *echo.Echo) error {
	metrics, err := metric.CreateMetrics(s.cfg.Metrics.URL, s.cfg.Metrics.ServiceName)
	if err != nil {
		s.logger.Errorf("CreateMetrics Error: %s", err)
	}
	s.logger.Infof(
		"Metrics available URL: %s, ServiceName: %s",
		s.cfg.Metrics.URL,
		s.cfg.Metrics.ServiceName,
	)

	// Init repositories
	aRepo := userRepository.NewUserRepository(s.db)
	tRepo := tweetRepository.NewTweetRepository(s.db)
	sRepo := sessionRepository.NewSessionRepository(s.redisClient, s.cfg)
	userRedisRepo := userRepository.NewUserRedisRepo(s.redisClient)
	fileRepo := fileRepository.NewFileRepository(s.cfg)
	followRepo := followRepository.NewFollowRepository(s.db)
	followRedisRepo := followRepository.NewFollowRedisRepo(s.redisClient)
	likeRepo := likeRepository.NewLikeRepository(s.db)

	// Init useCases
	userUC := userUseCase.NewUserUseCase(s.cfg, aRepo, userRedisRepo, followRedisRepo, s.logger)
	sessUC := usecase.NewSessionUseCase(sRepo, s.cfg)
	tweetUC := tweetUseCase.NewTweetUseCase(s.cfg, tRepo, s.logger)
	fileUC := fileUseCase.NewFileUseCase(s.cfg, fileRepo, s.logger)
	followUC := followUseCase.NewFollowUseCase(s.cfg, followRepo, followRedisRepo, s.logger)
	likeUC := likeUseCase.NewLikeUseCase(s.cfg, likeRepo, s.logger)

	// Init handlers
	userHandlers := userHttp.NewUserHandlers(s.cfg, userUC, sessUC, fileUC, followUC, likeUC, tweetUC, s.logger)
	tweetHandlers := tweetHttp.NewTweetHandlers(s.cfg, tweetUC, fileUC, likeUC, s.logger)

	mw := apiMiddlewares.NewMiddlewareManager(sessUC, userUC, s.cfg, []string{"*"}, s.logger)

	e.Use(mw.RequestLoggerMiddleware)

	docs.SwaggerInfo.Title = "Go example REST API"
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Static file
	e.Static(s.cfg.File.FilePath, s.cfg.File.FilePath)

	if s.cfg.Server.SSL {
		e.Pre(middleware.HTTPSRedirect())
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID, csrf.CSRFHeader},
	}))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	e.Use(middleware.RequestID())
	e.Use(mw.MetricsMiddleware(metrics))

	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("3M"))
	if s.cfg.Server.RateLimit {
		e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(s.cfg.Server.RateLimitPerSec))))
	}
	if s.cfg.Server.Debug {
		e.Use(mw.DebugMiddleware)
	}

	v1 := e.Group("/api/v1")

	health := v1.Group("/health")
	userGroup := v1.Group("/users")
	tweetGroup := v1.Group("/tweets")

	userHttp.MapUserRoutes(userGroup, userHandlers, mw)
	tweetHttp.MapTweetRoutes(tweetGroup, tweetHandlers, mw)

	health.GET("", func(c echo.Context) error {
		s.logger.Infof("Health check RequestID: %s", utils.GetRequestID(c))
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	return nil
}
