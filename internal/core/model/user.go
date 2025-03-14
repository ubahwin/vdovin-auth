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

func (u *User) GetAvailableValues(scopeAccess SessionScope) map[string]interface{} {
	values := map[string]interface{}{}

	switch {
	case scopeAccess.IsAllowed(SessionScope(SessionScopeBasicInfoEntry)):
		values["first_name"] = u.FirstName
		values["last_name"] = u.LastName
	case scopeAccess.IsAllowed(SessionScope(SessionScopePhoneEntry)):
		values["phone"] = u.Phone
	case scopeAccess.IsAllowed(SessionScope(SessionScopeEmailEntry)):
		values["email"] = u.Email
	}

	return values
}
