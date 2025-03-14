package sessionstorage

import (
	"errors"
	"github.com/google/uuid"
	"github.com/ubahwin/vdovin-auth/internal/core/model"
	"github.com/ubahwin/vdovin-auth/pkg/strrand"
	"sync"
	"time"
)

type Session struct {
	UserID         uuid.UUID
	Scope          model.SessionScope
	AccessToken    string
	RefreshToken   string
	AccessTokenTTL time.Duration
	UpdatedAt      time.Time
	CreatedAt      time.Time
}

type Mem struct {
	sessions           map[string]Session
	mu                 sync.Mutex
	accessTokenLength  int
	refreshTokenLength int
	accessTokenTTL     time.Duration
}

func NewMem(accessTokenLength, refreshTokenLength int, accessTokenTTL time.Duration) *Mem {
	return &Mem{
		sessions:           make(map[string]Session),
		mu:                 sync.Mutex{},
		accessTokenLength:  accessTokenLength,
		refreshTokenLength: refreshTokenLength,
		accessTokenTTL:     accessTokenTTL,
	}
}

func (storage *Mem) Create(id uuid.UUID, scope model.SessionScope) (*model.Session, error) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	accessToken := strrand.RandSeqStr(storage.accessTokenLength)
	refreshToken := strrand.RandSeqStr(storage.refreshTokenLength)

	now := time.Now().UTC()

	session := Session{
		UserID:         id,
		Scope:          scope,
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
		AccessTokenTTL: storage.accessTokenTTL,
		UpdatedAt:      now,
		CreatedAt:      now,
	}

	storage.sessions[accessToken] = session

	return &model.Session{
		UserID:         id,
		Scope:          session.Scope,
		AccessToken:    session.AccessToken,
		RefreshToken:   session.RefreshToken,
		AccessTokenTTL: session.AccessTokenTTL,
		UpdatedAt:      session.UpdatedAt,
		CreatedAt:      session.CreatedAt,
	}, nil
}

func (storage *Mem) Get(accessToken string) (*model.Session, error) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	session, ok := storage.sessions[accessToken]
	if !ok {
		return nil, ErrNotFound
	}

	if session.UpdatedAt.Add(session.AccessTokenTTL).Before(time.Now().UTC()) {
		delete(storage.sessions, accessToken)
		return nil, ErrNotFound
	}

	return &model.Session{
		UserID:         session.UserID,
		Scope:          session.Scope,
		AccessToken:    session.AccessToken,
		RefreshToken:   session.RefreshToken,
		AccessTokenTTL: session.AccessTokenTTL,
		UpdatedAt:      session.UpdatedAt,
		CreatedAt:      session.CreatedAt,
	}, nil
}

func (storage *Mem) Refresh(refreshToken string) (*model.Session, error) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	for _, session := range storage.sessions {
		if session.RefreshToken == refreshToken {
			delete(storage.sessions, session.AccessToken)

			accessToken := strrand.RandSeqStr(storage.accessTokenLength)
			refreshToken = strrand.RandSeqStr(storage.refreshTokenLength)

			session := Session{
				UserID:         session.UserID,
				Scope:          session.Scope,
				AccessToken:    accessToken,
				RefreshToken:   refreshToken,
				AccessTokenTTL: storage.accessTokenTTL,
				UpdatedAt:      time.Now().UTC(),
				CreatedAt:      session.CreatedAt,
			}

			storage.sessions[accessToken] = session

			return &model.Session{
				UserID:         session.UserID,
				Scope:          session.Scope,
				AccessToken:    session.AccessToken,
				RefreshToken:   session.RefreshToken,
				AccessTokenTTL: session.AccessTokenTTL,
				UpdatedAt:      session.UpdatedAt,
				CreatedAt:      session.CreatedAt,
			}, nil
		}
	}

	return nil, ErrNotFound
}

func (storage *Mem) Delete(accessToken string) error {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	delete(storage.sessions, accessToken)

	return nil
}

var ErrNotFound = errors.New("not found")
