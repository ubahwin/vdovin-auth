package authapi

import (
	"github.com/ubahwin/vdovin-auth/internal/api"
	"net/http"
)

type CodeReq struct {
	/// Only trusted applications can have this token
	Token string `json:"token"`

	/// QR code transform to this code in trusted applications
	Code string `json:"code"`
}

type CodeResp struct {
	Success bool   `json:"success"`
	Comment string `json:"comment,omitempty"`
}

func (req CodeReq) Validate(_ *api.Context) error {
	return nil
}

func (g *Group) Code(_ *api.Context, req *CodeReq) (*CodeResp, int) {
	err := g.authorizer.FindActiveSession(req.Token, req.Code)
	if err != nil {
		return &CodeResp{
			Success: false,
			Comment: err.Error(),
		}, http.StatusBadRequest
	}

	return &CodeResp{
		Success: true,
	}, http.StatusOK
}
