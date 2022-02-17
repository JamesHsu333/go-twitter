package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/JamesHsu333/go-twitter/config"
	"github.com/JamesHsu333/go-twitter/internal/models"
	"github.com/JamesHsu333/go-twitter/internal/user"
	"github.com/JamesHsu333/go-twitter/pkg/httpErrors"
	"github.com/JamesHsu333/go-twitter/pkg/tracer"
	"github.com/JamesHsu333/go-twitter/pkg/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Auth sessions middleware using redis
func (mw *MiddlewareManager) AuthSessionMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, span := tracer.NewSpan(c.Request().Context(), "MiddlewareManager.AuthSessionMiddleware", nil)
		defer span.End()

		cookie, err := c.Cookie(mw.cfg.Session.Name)
		if err != nil {
			mw.logger.Errorf("AuthSessionMiddleware RequestID: %s, Error: %s",
				utils.GetRequestID(c),
				err.Error(),
			)
			if err == http.ErrNoCookie {
				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(err))
			}
			return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
		}

		sid := cookie.Value

		sess, err := mw.sessUC.GetSessionByID(ctx, cookie.Value)
		if err != nil {
			mw.logger.Errorf("GetSessionByID RequestID: %s, CookieValue: %s, Error: %s",
				utils.GetRequestID(c),
				cookie.Value,
				err.Error(),
			)
			return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
		}

		user, err := mw.userUC.GetCacheByID(ctx, sess.UserID, sess.UserID)
		if err != nil {
			user, err = mw.userUC.GetByID(ctx, sess.UserID, sess.UserID)
			if err != nil {
				mw.logger.Errorf("GetByID RequestID: %s, Error: %s",
					utils.GetRequestID(c),
					err.Error(),
				)
				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
			}
		}

		c.Set("sid", sid)
		c.Set("uid", sess.SessionID)
		c.Set("user", user)

		c.SetRequest(c.Request().WithContext(context.WithValue(ctx, utils.UserCtxKey{}, user)))

		mw.logger.Infof(
			"SessionMiddleware, RequestID: %s,  IP: %s, UserID: %s, CookieSessionID: %s",
			utils.GetRequestID(c),
			utils.GetIPAddress(c),
			user.UserID.String(),
			cookie.Value,
		)

		return next(c)
	}
}

// JWT way of auth using cookie or Authorization header
func (mw *MiddlewareManager) AuthJWTMiddleware(userUC user.UseCase, cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			span := tracer.SpanFromContext(c.Request().Context())
			defer span.End()

			bearerHeader := c.Request().Header.Get("Authorization")

			mw.logger.Infof("auth middleware bearerHeader %s", bearerHeader)

			if bearerHeader != "" {
				headerParts := strings.Split(bearerHeader, " ")
				if len(headerParts) != 2 {
					mw.logger.Error("auth middleware", zap.String("headerParts", "len(headerParts) != 2"))
					return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
				}

				tokenString := headerParts[1]

				if err := mw.validateJWTToken(tokenString, userUC, c, cfg); err != nil {
					mw.logger.Error("middleware validateJWTToken", zap.String("headerJWT", err.Error()))
					return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
				}

				return next(c)
			}

			cookie, err := c.Cookie("jwt-token")
			if err != nil {
				mw.logger.Errorf("c.Cookie", err.Error())
				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
			}

			if err = mw.validateJWTToken(cookie.Value, userUC, c, cfg); err != nil {
				mw.logger.Errorf("validateJWTToken", err.Error())
				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
			}
			return next(c)
		}
	}
}

// Admin role
func (mw *MiddlewareManager) AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		span := tracer.SpanFromContext(c.Request().Context())
		defer span.End()

		user, ok := c.Get("user").(*models.User)
		if !ok || *user.Role != "admin" {
			return c.JSON(http.StatusForbidden, httpErrors.NewUnauthorizedError(httpErrors.PermissionDenied))
		}
		return next(c)
	}
}

// Check owner using ctx user
func (mw *MiddlewareManager) OwnerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			span := tracer.SpanFromContext(c.Request().Context())
			defer span.End()

			user, ok := c.Get("user").(*models.User)
			if !ok {
				mw.logger.Errorf("Error c.Get(user) RequestID: %s, ERROR: %s,", utils.GetRequestID(c), "invalid user ctx")
				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
			}

			if user.UserID.String() != c.Param("user_id") {
				mw.logger.Errorf("Error c.Get(user) RequestID: %s, UserID: %s, ERROR: %s,",
					utils.GetRequestID(c),
					user.UserID.String(),
					"invalid user ctx",
				)
				return c.JSON(http.StatusForbidden, httpErrors.NewForbiddenError(httpErrors.Forbidden))
			}

			return next(c)
		}
	}
}

// Role based auth middleware, using ctx user
func (mw *MiddlewareManager) RoleBasedAuthMiddleware(roles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			span := tracer.SpanFromContext(c.Request().Context())
			defer span.End()

			user, ok := c.Get("user").(*models.User)
			if !ok {
				mw.logger.Errorf("Error c.Get(user) RequestID: %s, UserID: %s, ERROR: %s,",
					utils.GetRequestID(c),
					user.UserID.String(),
					"invalid user ctx",
				)
				return c.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(httpErrors.Unauthorized))
			}

			for _, role := range roles {
				if role == *user.Role {
					return next(c)
				}
			}

			mw.logger.Errorf("Error c.Get(user) RequestID: %s, UserID: %s, ERROR: %s,",
				utils.GetRequestID(c),
				user.UserID.String(),
				"invalid user ctx",
			)

			return c.JSON(http.StatusForbidden, httpErrors.NewForbiddenError(httpErrors.PermissionDenied))
		}
	}
}

func (mw *MiddlewareManager) validateJWTToken(tokenString string, userUC user.UseCase, c echo.Context, cfg *config.Config) error {
	ctx, span := tracer.NewSpan(c.Request().Context(), "MiddlewareManager.validateJWTToken", nil)
	defer span.End()

	if tokenString == "" {
		return httpErrors.InvalidJWTToken
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signin method %v", token.Header["alg"])
		}
		secret := []byte(cfg.Server.JwtSecretKey)
		return secret, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return httpErrors.InvalidJWTToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["id"].(string)
		if !ok {
			return httpErrors.InvalidJWTClaims
		}

		userUUID, err := uuid.Parse(userID)
		if err != nil {
			return err
		}

		u, err := userUC.GetByID(c.Request().Context(), userUUID, userUUID)
		if err != nil {
			return err
		}

		c.Set("user", u)

		ctx := context.WithValue(ctx, utils.UserCtxKey{}, u)
		// req := c.Request().WithContext(ctx)
		c.SetRequest(c.Request().WithContext(ctx))
	}
	return nil
}

// Check auth middleware
func (mw *MiddlewareManager) CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		c, span := tracer.NewSpan(ctx.Request().Context(), "MiddlewareManager.validateJWTToken", nil)
		defer span.End()

		cookie, err := ctx.Cookie("session_id")
		if err != nil {
			mw.logger.Errorf("CheckAuth.ctx.Cookie: %s, Cookie: %#v, Error: %s",
				utils.GetRequestID(ctx),
				cookie,
				err,
			)
			return ctx.JSON(http.StatusUnauthorized, httpErrors.NewUnauthorizedError(err))
		}
		sid := cookie.Value

		session, err := mw.sessUC.GetSessionByID(c, sid)
		if err != nil {
			// Cookie is invalid, delete it from browser
			newCookie := http.Cookie{Name: "session_id", Value: sid, Expires: time.Now().AddDate(-1, 0, 0)}
			ctx.SetCookie(&newCookie)

			mw.logger.Errorf("CheckAuth.sessUC.GetSessionByID: %s, Cookie: %#v, Error: %s",
				utils.GetRequestID(ctx),
				cookie,
				err,
			)
			return ctx.JSON(http.StatusUnauthorized, httpErrors.NoCookie)
		}

		ctx.Set("uid", session.SessionID)
		ctx.Set("sid", sid)
		return next(ctx)
	}
}
