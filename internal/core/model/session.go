package model

import (
	"github.com/google/uuid"
	"time"
)

type Session struct {
	UserID         uuid.UUID
	Scope          SessionScope
	AccessToken    string
	RefreshToken   string
	AccessTokenTTL time.Duration
	UpdatedAt      time.Time
	CreatedAt      time.Time
}
