package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Username  string         `json:"username" gorm:"unique;not null"`
    Email     string         `json:"email" gorm:"unique;not null"`
    Password  string         `json:"-" gorm:"not null"`
    Role      string         `json:"role" gorm:"default:'user'"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`
}

type UserResponse struct {
    ID        uint      `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Role      string    `json:"role"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

type AuthResponse struct {
    Token string       `json:"token"`
    User  UserResponse `json:"user"`
}

func (u *User) ToUserResponse() UserResponse {
    var deletedAt *time.Time
    if u.DeletedAt.Valid {
        deletedAt = &u.DeletedAt.Time
    }
    
    return UserResponse{
        ID:        u.ID,
        Username:  u.Username,
        Email:     u.Email,
        Role:      u.Role,
        CreatedAt: u.CreatedAt,
        UpdatedAt: u.UpdatedAt,
        DeletedAt: deletedAt,
    }
}