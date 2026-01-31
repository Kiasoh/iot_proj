package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID           int        `json:"id"`
	Email        string     `json:"email"`
	Password     string     `json:"-"`
	KeyCard      *string `json:"key_card"`
	AccessLevel  int        `json:"access_level"`
	LastAccessed *time.Time `json:"last_accessed"`
	CreatedAt    time.Time  `json:"created_at"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Claims struct {
	ID int `json:"id"`
	jwt.RegisteredClaims
}

type EntryLog struct {
	ID        int        `json:"id"`
	KeyCard   *string `json:"key_card"`
	Status    string     `json:"status"`
	Message   *string    `json:"message"`
	CreatedAt time.Time  `json:"created_at"`
}

func StringPtr(s string) *string {
	return &s
}
