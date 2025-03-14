package userapi

import (
	"github.com/ubahwin/vdovin-auth/internal/api"
	"net/http"
	"time"
)

type SignInReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInResp struct {
	Success      bool      `json:"success"`
	Comment      string    `json:"comment,omitempty"`
	AccessToken  string    `json:"access_token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	ExpiresAt    time.Time `json:"expires_at,omitempty"`
	User         UserResp  `json:"user,omitempty"`
}

type UserResp struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

func (req SignInReq) Validate(_ *api.Context) error {
	return nil
}

func (g *Group) SignIn(_ *api.Context, req *SignInReq) (*SignInResp, int) {
	session, user, err := g.userManager.SignIn(req.Username, req.Password)
	if err != nil {
		return &SignInResp{
			Success: false,
			Comment: err.Error(),
		}, http.StatusOK
	}

	return &SignInResp{
		Success:      true,
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
		ExpiresAt:    session.UpdatedAt.Add(session.AccessTokenTTL),
		User: UserResp{
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			Email:     user.Email,
		},
	}, http.StatusOK
}
