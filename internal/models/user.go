package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User full model
type User struct {
	UserID      uuid.UUID  `json:"user_id" db:"user_id" redis:"user_id" validate:"omitempty"`
	UserName    string     `json:"user_name" db:"user_name" redis:"user_name" validate:"omitempty,required,lte=32"`
	Name        string     `json:"name" db:"name" redis:"name" validate:"omitempty,required,lte=32"`
	Email       string     `json:"email,omitempty" db:"email" redis:"email" validate:"omitempty,lte=64,email"`
	Password    string     `json:"password,omitempty" db:"password" redis:"password" validate:"omitempty,required,gte=6"`
	Role        *string    `json:"role,omitempty" db:"role" redis:"role" validate:"omitempty,lte=10"`
	About       *string    `json:"about,omitempty" db:"about" redis:"about" validate:"omitempty,lte=160"`
	Avatar      *string    `json:"avatar,omitempty" db:"avatar" redis:"avatar" validate:"omitempty,lte=512,url"`
	Header      *string    `json:"header,omitempty" db:"header" redis:"header" validate:"omitempty,lte=512,url"`
	PhoneNumber *string    `json:"phone_number,omitempty" db:"phone_number" redis:"phone_number" validate:"omitempty,lte=20"`
	Country     *string    `json:"country,omitempty" db:"country" redis:"country" validate:"omitempty,lte=20"`
	Gender      *string    `json:"gender,omitempty" db:"gender" redis:"gender" validate:"omitempty,lte=10"`
	Birthday    *time.Time `json:"birthday,omitempty" db:"birthday" redis:"birthday" validate:"omitempty,lte=10"`
	Followers   *int64     `json:"followers" db:"followers" redis:"followers" validate:"omitempty"`
	Following   *int64     `json:"following" db:"following" redis:"following" validate:"omitempty"`
	IsFollowing *bool      `json:"is_following" db:"is_following" redis:"is_following" validate:"omitempty"`
	CreatedAt   time.Time  `json:"created_at,omitempty" db:"created_at" redis:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at,omitempty" db:"updated_at" redis:"updated_at"`
	LoginDate   time.Time  `json:"login_date" db:"login_date" redis:"login_date"`
}

// Hash user password with bcrypt
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Compare user password and payload
func (u *User) ComparePasswords(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}

// Sanitize user password
func (u *User) SanitizePassword() {
	u.Password = ""
}

func (u *User) SanitizeFollow() {
	u.Followers = nil
	u.Following = nil
	u.IsFollowing = nil
}

// Prepare user for register
func (u *User) PrepareCreate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Password = strings.TrimSpace(u.Password)

	if err := u.HashPassword(); err != nil {
		return err
	}

	if u.PhoneNumber != nil {
		*u.PhoneNumber = strings.TrimSpace(*u.PhoneNumber)
	}
	if u.Role != nil {
		*u.Role = strings.ToLower(strings.TrimSpace(*u.Role))
	}
	return nil
}

// Prepare user for register
func (u *User) PrepareUpdate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))

	if u.PhoneNumber != nil {
		*u.PhoneNumber = strings.TrimSpace(*u.PhoneNumber)
	}
	if u.Role != nil {
		*u.Role = strings.ToLower(strings.TrimSpace(*u.Role))
	}
	return nil
}

// All Users response
type UsersList struct {
	TotalCount int     `json:"total_count"`
	TotalPages int     `json:"total_pages"`
	Page       int     `json:"page"`
	Size       int     `json:"size"`
	HasMore    bool    `json:"has_more"`
	Users      []*User `json:"users"`
}

// Find user query
type UserWithToken struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}
