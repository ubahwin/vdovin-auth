package userapi

import (
	"github.com/ubahwin/vdovin-auth/internal/core/user"
	"log"
)

type Group struct {
	userManager *user.Manager
	log         *log.Logger
}

func New(log *log.Logger, um *user.Manager) *Group {
	return &Group{
		userManager: um,
		log:         log,
	}
}
