package user

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                      uint           `json:"id" gorm:"primaryKey"`
	Email                   string         `json:"email" gorm:"unique;not null"`
	Password                string         `json:"-" gorm:"not null"`
	Name                    string         `json:"name" gorm:"not null"`
	Role                    string         `json:"role" gorm:"default:'user'"`
	IsEmailVerified         bool           `json:"is_email_verified" gorm:"default:false"`
	EmailVerificationToken  *string        `json:"-" gorm:"uniqueIndex;default:null"`
	EmailVerificationExpiry *time.Time     `json:"-"`
	PasswordResetToken      *string        `json:"-" gorm:"uniqueIndex;default:null"`
	PasswordResetExpiry     *time.Time     `json:"-"`
	CreatedAt               time.Time      `json:"created_at"`
	UpdatedAt               time.Time      `json:"updated_at"`
	DeletedAt               gorm.DeletedAt `json:"-" gorm:"index"`
}

// Add these to your existing request/response structs
type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordConfirmRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
