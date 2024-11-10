package user

import (
	"errors"
	"github.com/google/uuid"
	"github.com/ubahwin/vdovin-auth/internal/core/model"
)

type UserStorage interface {
	Create(user *model.User) (*uuid.UUID, error)
	GetByID(id uuid.UUID) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	UpdateByID(user model.User) error
	DeleteByID(id uuid.UUID) error
}

type SessionStorage interface {
	Create(id uuid.UUID, scope model.SessionScope) (*model.Session, error)
	Get(accessToken string) (*model.Session, error)
	Refresh(refreshToken string) (*model.Session, error)
	Delete(accessToken string) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) (bool, error)
}

type Manager struct {
	userStorage    UserStorage
	sessionStorage SessionStorage
	passwordHasher PasswordHasher
}

func NewManager(us UserStorage, ss SessionStorage, ph PasswordHasher) *Manager {
	return &Manager{userStorage: us, sessionStorage: ss, passwordHasher: ph}
}

func (m *Manager) Register(user *model.User) (*uuid.UUID, error) {
	_, getUserError := m.userStorage.GetByUsername(user.Username)
	if getUserError == nil {
		return nil, errors.New("username already exist")
	}

	var hashError error
	user.PasswordHash, hashError = m.passwordHasher.Hash(user.Password)
	if hashError != nil {
		return nil, hashError
	}

	return m.userStorage.Create(user)
}

func (m *Manager) SignIn(username, password string, scope model.SessionScope) (*model.Session, error) {
	user, err := m.userStorage.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	ok, err := m.passwordHasher.Compare(user.PasswordHash, password)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, ErrInvalidPassword
	}

	return m.sessionStorage.Create(user.ID, scope)
}

var ErrInvalidPassword = errors.New("invalid password")
