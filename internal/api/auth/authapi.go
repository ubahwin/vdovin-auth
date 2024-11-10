package authapi

import (
	"github.com/ubahwin/vdovin-auth/internal/core/auth"
	"log"
)

type Group struct {
	authorizer *auth.Authorizer
	log        *log.Logger
}

func New(log *log.Logger, authorizer *auth.Authorizer) *Group {
	return &Group{
		authorizer: authorizer,
		log:        log,
	}
}
