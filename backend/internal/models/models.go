package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	LinkedAccounts []LinkedAccount `gorm:"foreignKey:UserID" json:"linked_accounts,omitempty"`
}

// LinkedAccount represents a connected social media account
type LinkedAccount struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	Provider    string         `gorm:"not null;default:'linkedin'" json:"provider"`
	AccountID   string         `gorm:"not null" json:"account_id"`
	AccountName string         `json:"account_name,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// Request/Response DTOs

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  struct {
		ID    uint   `json:"id"`
		Email string `json:"email"`
	} `json:"user"`
}

type LinkedInCookieRequest struct {
	Cookie string `json:"cookie" binding:"required"`
}

type LinkedInCredentialsRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LinkedInConnectResponse struct {
	Message   string        `json:"message"`
	AccountID string        `json:"account_id"`
	Account   LinkedAccount `json:"account"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
