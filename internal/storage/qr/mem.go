package qrstorage

import (
	"errors"
	"github.com/google/uuid"
	"github.com/ubahwin/vdovin-auth/internal/core/model"
	"github.com/ubahwin/vdovin-auth/pkg/strrand"
	"net/url"
	"sync"
	"time"
)

type QRSession struct {
	ID          uuid.UUID
	RedirectURI url.URL
	Scope       model.SessionScope
	CreatedAt   time.Time
}

type Mem struct {
	sessions     map[string]QRSession
	mu           sync.Mutex
	qrCodeLength int
	qrTTL        time.Duration
}

func NewMem(qrCodeLength int, qrTTL time.Duration) *Mem {
	return &Mem{
		sessions:     make(map[string]QRSession),
		mu:           sync.Mutex{},
		qrCodeLength: qrCodeLength,
		qrTTL:        qrTTL,
	}
}

func (storage *Mem) Create(scope model.SessionScope, redirectURI url.URL) string {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	qrCode := strrand.RandSeqStr(storage.qrCodeLength)

	storage.sessions[qrCode] = QRSession{
		ID:          uuid.New(),
		RedirectURI: redirectURI,
		Scope:       scope,
		CreatedAt:   time.Now().UTC(),
	}

	return qrCode
}

func (storage *Mem) Get(qrCode string) (*model.QRSession, error) {
	storage.mu.Lock()
	defer storage.mu.Unlock()

	session, ok := storage.sessions[qrCode]
	if !ok {
		return nil, ErrNotFound
	}

	if session.CreatedAt.Add(storage.qrTTL).Before(time.Now().UTC()) {
		delete(storage.sessions, qrCode)
		return nil, ErrNotFound
	}

	return &model.QRSession{
		RedirectURI: session.RedirectURI,
		Scope:       session.Scope,
	}, nil
}

var ErrNotFound = errors.New("not found")
