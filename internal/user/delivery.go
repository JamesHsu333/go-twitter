package user

import "github.com/labstack/echo/v4"

// User HTTP Handlers interface
type Handlers interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc
	Logout() echo.HandlerFunc
	Update() echo.HandlerFunc
	UpdateRole() echo.HandlerFunc
	Delete() echo.HandlerFunc
	GetUserByID() echo.HandlerFunc
	GetUserByUserName() echo.HandlerFunc
	Follow() echo.HandlerFunc
	GetFollowers() echo.HandlerFunc
	GetFollowing() echo.HandlerFunc
	DeleteFollowing() echo.HandlerFunc
	FindByName() echo.HandlerFunc
	GetUsers() echo.HandlerFunc
	GetMe() echo.HandlerFunc
	UploadAvatar() echo.HandlerFunc
	UploadHeader() echo.HandlerFunc
	GetCSRFToken() echo.HandlerFunc
	Like() echo.HandlerFunc
	GetLikedTweets() echo.HandlerFunc
	DeleteLiked() echo.HandlerFunc
	GetTweetsByUserID() echo.HandlerFunc
}
