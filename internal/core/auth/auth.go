package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/ubahwin/vdovin-auth/internal/core/model"
	"net/http"
	"net/url"
)

//type UserStorage interface {
//	Create(user *model.User) (int, error)
//	GetByID(id int) (*model.User, error)
//	GetByUsername(username string) (*model.User, error)
//	UpdateByID(user model.User) error
//	DeleteByID(id int) error
//}

type SessionStorage interface {
	Create(id uuid.UUID, scope model.SessionScope) (*model.Session, error)
	Get(accessToken string) (*model.Session, error)
	Refresh(refreshToken string) (*model.Session, error)
	Delete(accessToken string) error
}

//type PasswordHasher interface {
//	Hash(password string) (string, error)
//	Compare(hash, password string) (bool, error)
//}

type QRSessionStorage interface {
	// Return QR code
	Create(scope model.SessionScope, redirectURI url.URL) string

	// Return redirect_uri if QR is active
	Get(qrCode string) (*model.QRSession, error)
}

type Authorizer struct {
	//userStorage    UserStorage
	sessionStorage SessionStorage
	//passwordHasher PasswordHasher
	qrStorage QRSessionStorage
}

func NewAuthorizer(
	//us UserStorage,
	ss SessionStorage,
	//ph PasswordHasher,
	qrs QRSessionStorage,
) *Authorizer {
	return &Authorizer{
		//userStorage: us,
		sessionStorage: ss,
		//passwordHasher: ph,
		qrStorage: qrs,
	}
}

func (a *Authorizer) Authorize(scope model.SessionScope, redirectURI url.URL) (string, error) {
	return a.qrStorage.Create(scope, redirectURI), nil
}

func (a *Authorizer) FindActiveSession(token, code string) error {
	// Проверка токена доверенного приложения
	userSession, err := a.sessionStorage.Get(token)
	if err != nil {
		return err
	}

	// Проверяем, а создавался ли такой QR-код
	qrSession, err := a.qrStorage.Get(code)
	if err != nil {
		return err
	}

	// Доступ к данным, который запрашивается от стороннего сервиса, должен быть у доверенного приложения
	if userSession.Scope.IsAllowed(qrSession.Scope) {
		return errors.New("scope error")
	}

	// Создаём новую сессию для сервиса, это уже другая сессия, не для доверенного приложения
	session, err := a.sessionStorage.Create(uuid.New(), qrSession.Scope)
	if err != nil {
		return err
	}

	// Посылаем на `redirect_uri`, `access_token`, `auth_id` (qr_code)
	return a.SendAccessToken(qrSession.AuthID, session.AccessToken, qrSession.RedirectURI) // TODO: create independent service for send request
}

func (a *Authorizer) SendAccessToken(authID, accessToken string, redirectURI url.URL) error {
	data := map[string]interface{}{
		"auth_id":      authID,
		"access_token": accessToken,
		"scope":        "basic,phone",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", redirectURI.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	return nil
}
