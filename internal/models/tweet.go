package models

import (
	"time"

	"github.com/google/uuid"
)

// Tweet Model
type Tweet struct {
	ID        uint64    `json:"id" form:"id" db:"id" redis:"id" validate:"omitempty"`
	UserID    uuid.UUID `json:"user_id" form:"user_id" db:"user_id" redis:"user_id" validate:"omitempty,required"`
	Text      string    `json:"text" form:"text" db:"text" redis:"text" validate:"omitempty,required,lte=260"`
	Image     *string   `json:"image,omitempty" form:"image" db:"image" redis:"image" validate:"omitempty,lte=512"`
	Likes     int64     `json:"likes" form:"likes" db:"likes" redis:"likes" validate:"omitempty"`
	Replys    int64     `json:"replys" form:"replys" db:"replys" redis:"replys" validate:"omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" form:"created_at" db:"created_at" redis:"created_at"`
}

type TweetWithUser struct {
	ID           uint64    `json:"id" db:"id" redis:"id" validate:"omitempty"`
	Text         string    `json:"text" db:"text" redis:"text" validate:"omitempty,required,lte=260"`
	Image        *string   `json:"image,omitempty" db:"image" redis:"image" validate:"omitempty,lte=512,url"`
	CreatedAt    time.Time `json:"created_at,omitempty" db:"created_at" redis:"created_at"`
	UserID       uuid.UUID `json:"user_id" db:"user_id" redis:"user_id" validate:"omitempty"`
	UserName     string    `json:"user_name" db:"user_name" redis:"user_name" validate:"omitempty,required,lte=32"`
	Name         string    `json:"name" db:"name" redis:"name" validate:"omitempty,required,lte=32"`
	About        *string   `json:"about,omitempty" db:"about" redis:"about" validate:"omitempty,lte=160"`
	Avatar       *string   `json:"avatar,omitempty" db:"avatar" redis:"avatar" validate:"omitempty,lte=512,url"`
	Likes        int64     `json:"likes" db:"likes" redis:"likes" validate:"omitempty"`
	Replys       int64     `json:"replys" db:"replys" redis:"replys" validate:"omitempty"`
	AlreadyLiked bool      `json:"already_liked" db:"already_liked" redis:"already_liked"`
}

// All Tweets response
type TweetsList struct {
	TotalCount int              `json:"total_count"`
	TotalPages int              `json:"total_pages"`
	Page       int              `json:"page"`
	Size       int              `json:"size"`
	HasMore    bool             `json:"has_more"`
	Tweets     []*TweetWithUser `json:"tweets"`
}
