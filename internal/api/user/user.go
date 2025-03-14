package userapi

import (
	"github.com/ubahwin/vdovin-auth/internal/api"
	"net/http"
)

type UserInfoReq struct {
	AccessToken string `json:"access_token"`
}

type UserInfoResp struct {
	Success bool                   `json:"success"`
	Comment string                 `json:"comment,omitempty"`
	User    map[string]interface{} `json:"user,omitempty"`
}

func (req UserInfoReq) Validate(_ *api.Context) error {
	return nil
}

func (g *Group) UserInfo(_ *api.Context, req *UserInfoReq) (*UserInfoResp, int) {
	user, err := g.userManager.UserInfo(req.AccessToken)
	if err != nil {
		return &UserInfoResp{
			Success: false,
			Comment: err.Error(),
		}, http.StatusOK
	}

	return &UserInfoResp{
		Success: true,
		User:    user,
	}, http.StatusOK
}
