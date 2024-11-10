package model

import "github.com/google/uuid"

type User struct {
	ID           uuid.UUID
	Username     string
	Password     string
	PasswordHash string
	FirstName    string
	LastName     string
	Phone        string
	Email        string

	_ struct{}
}
