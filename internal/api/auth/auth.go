package authapi

import (
	"github.com/ubahwin/vdovin-auth/internal/api"
	"github.com/ubahwin/vdovin-auth/internal/core/model"
	"net/http"
	"net/url"
)

type AuthReq struct {
	RedirectURI string `json:"redirect_uri"`
	Scope       string `json:"scope"`
}

type AuthResp struct {
	QRCode string `json:"qr_code"`
}

func (req AuthReq) Validate(_ *api.Context) error {
	return nil
}

func (g *Group) Auth(_ *api.Context, req *AuthReq) (*AuthResp, int) {
	scope, err := model.ParseSessionScope(req.Scope)
	if err != nil {
		return nil, http.StatusBadRequest
	}

	redirectURI, err := url.Parse(req.RedirectURI)
	if err != nil {
		return nil, http.StatusBadRequest
	}

	qrCode, err := g.authorizer.Authorize(scope, *redirectURI)
	if err != nil {
		return nil, http.StatusBadRequest
	}

	return &AuthResp{
		QRCode: qrCode,
	}, http.StatusOK
}
