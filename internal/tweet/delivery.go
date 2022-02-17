package tweet

import "github.com/labstack/echo/v4"

// Tweet HTTP Handlers interface
type Handlers interface {
	Create() echo.HandlerFunc
	CreateReply() echo.HandlerFunc
	GetTweetByID() echo.HandlerFunc
	GetTweets() echo.HandlerFunc
	GetReplyTweets() echo.HandlerFunc
	GetLikedUsers() echo.HandlerFunc
	Delete() echo.HandlerFunc
}
