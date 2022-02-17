package http

import (
	"github.com/JamesHsu333/go-twitter/internal/middleware"
	"github.com/JamesHsu333/go-twitter/internal/tweet"
	"github.com/labstack/echo/v4"
)

// Map tweet routes
func MapTweetRoutes(tweetGroup *echo.Group, h tweet.Handlers, mw *middleware.MiddlewareManager) {
	tweetGroup.Use(mw.AuthSessionMiddleware)
	tweetGroup.GET("/:tweet_id", h.GetTweetByID())
	tweetGroup.GET("", h.GetTweets())
	tweetGroup.GET("/:tweet_id/replys", h.GetReplyTweets())
	tweetGroup.GET("/:tweet_id/liking_users", h.GetLikedUsers())
	tweetGroup.DELETE("/:tweet_id", h.Delete(), mw.CSRF)
	tweetGroup.POST("", h.Create(), mw.CSRF)
	tweetGroup.POST("/:tweet_id/reply", h.CreateReply(), mw.CSRF)
}
