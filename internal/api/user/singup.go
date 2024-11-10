package userapi

import (
	"github.com/google/uuid"
	"github.com/ubahwin/vdovin-auth/internal/api"
	"github.com/ubahwin/vdovin-auth/internal/core/model"
	"net/http"
)

type SignUpReq struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
}

type SignUpResp struct {
	Success bool      `json:"success"`
	UserID  uuid.UUID `json:"user_id,omitempty"`
	Comment string    `json:"comment,omitempty"`
}

func (req SignUpReq) Validate(_ *api.Context) error {
	return nil
}

func (g *Group) SignUp(_ *api.Context, req *SignUpReq) (*SignUpResp, int) {
	userID, err := g.userManager.Register(&model.User{
		Username:  req.Username,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Email:     req.Email,
	})
	if err != nil {
		return &SignUpResp{
			Success: false,
			Comment: err.Error(),
		}, http.StatusInternalServerError
	}

	return &SignUpResp{Success: true, UserID: *userID}, http.StatusOK
}
