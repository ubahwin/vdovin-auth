package userapi

import (
	"github.com/ubahwin/vdovin-auth/internal/api"
	"github.com/ubahwin/vdovin-auth/internal/core/model"
	"net/http"
	"time"
)

type SignInReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Scope    string `json:"scope"`
}

type SignInResp struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (req SignInReq) Validate(_ *api.Context) error {
	return nil
}

func (g *Group) SignIn(_ *api.Context, req *SignInReq) (*SignInResp, int) {
	scope, err := model.ParseSessionScope(req.Scope)
	if err != nil {
		return nil, http.StatusBadRequest
	}

	session, err := g.userManager.SignIn(req.Username, req.Password, scope)
	if err != nil {
		return nil, http.StatusBadRequest
	}

	return &SignInResp{
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
		ExpiresAt:    session.UpdatedAt.Add(session.AccessTokenTTL),
	}, http.StatusOK
}
