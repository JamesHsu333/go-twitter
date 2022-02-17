package http

import (
	"github.com/JamesHsu333/go-twitter/internal/middleware"
	"github.com/JamesHsu333/go-twitter/internal/user"
	"github.com/labstack/echo/v4"
)

// Map users routes
func MapUserRoutes(userGroup *echo.Group, h user.Handlers, mw *middleware.MiddlewareManager) {
	userGroup.POST("/register", h.Register())
	userGroup.POST("/login", h.Login())
	userGroup.POST("/logout", h.Logout())
	userGroup.Use(mw.AuthSessionMiddleware)
	userGroup.GET("/me", h.GetMe())
	userGroup.GET("/by", h.FindByName())
	userGroup.GET("", h.GetUsers())
	userGroup.GET("/:user_id", h.GetUserByID())
	userGroup.GET("/username/:user_name", h.GetUserByUserName())
	userGroup.GET("/:user_id/followers", h.GetFollowers())
	userGroup.GET("/:user_id/following", h.GetFollowing())
	userGroup.GET("/token", h.GetCSRFToken())
	userGroup.GET("/:user_id/tweets", h.GetTweetsByUserID())
	userGroup.GET("/:user_id/liked_tweets", h.GetLikedTweets())
	userGroup.POST("/:user_id/avatar", h.UploadAvatar(), mw.OwnerMiddleware(), mw.CSRF)
	userGroup.POST("/:user_id/header", h.UploadHeader(), mw.OwnerMiddleware(), mw.CSRF)
	userGroup.POST("/:user_id/following", h.Follow(), mw.OwnerMiddleware(), mw.CSRF)
	userGroup.POST("/:user_id/liked", h.Like(), mw.OwnerMiddleware(), mw.CSRF)
	userGroup.PATCH("/:user_id", h.Update(), mw.OwnerMiddleware(), mw.CSRF)
	userGroup.PATCH("/:user_id/role", h.UpdateRole(), mw.RoleBasedAuthMiddleware([]string{"admin"}), mw.CSRF)
	userGroup.DELETE("/:user_id", h.Delete(), mw.CSRF, mw.RoleBasedAuthMiddleware([]string{"admin"}))
	userGroup.DELETE("/:user_id/following/:following_id", h.DeleteFollowing(), mw.OwnerMiddleware(), mw.CSRF)
	userGroup.DELETE("/:user_id/liked/:tweet_id", h.DeleteLiked(), mw.OwnerMiddleware(), mw.CSRF)
}
