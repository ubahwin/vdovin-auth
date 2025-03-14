package userstorage

import (
	"errors"
	"github.com/google/uuid"
	"github.com/ubahwin/vdovin-auth/internal/core/model"
	"sync"
)

type User struct {
	ID           uuid.UUID
	Username     string
	PasswordHash string
	FirstName    string
	LastName     string
	Phone        string
	Email        string
}

type Mem struct {
	users map[uuid.UUID]User
	mu    sync.Mutex
}

func NewMem() *Mem {
	return &Mem{users: make(map[uuid.UUID]User)}
}

func (storage *Mem) Create(user *model.User) (*uuid.UUID, error) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	id := uuid.New()

	storage.users[id] = User{
		ID:           id,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Phone:        user.Phone,
		Email:        user.Email,
	}

	return &id, nil
}

func (storage *Mem) GetByID(id uuid.UUID) (*model.User, error) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	user, ok := storage.users[id]
	if !ok {
		return nil, ErrNotFound
	}

	return &model.User{
		ID:           user.ID,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Phone:        user.Phone,
		Email:        user.Email,
	}, nil
}

func (storage *Mem) GetByUsername(username string) (*model.User, error) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	for _, user := range storage.users {
		if user.Username == username {
			return &model.User{
				ID:           user.ID,
				Username:     user.Username,
				PasswordHash: user.PasswordHash,
				FirstName:    user.FirstName,
				LastName:     user.LastName,
				Phone:        user.Phone,
				Email:        user.Email,
			}, nil
		}
	}

	return nil, ErrNotFound
}

func (storage *Mem) UpdateByID(user model.User) error {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	storage.users[user.ID] = User{
		ID:           user.ID,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Phone:        user.Phone,
		Email:        user.Email,
	}

	return nil
}

func (storage *Mem) DeleteByID(id uuid.UUID) error {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	delete(storage.users, id)

	return nil
}

var ErrNotFound = errors.New("not found")
