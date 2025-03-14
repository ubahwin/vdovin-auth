package model

import (
	"net/url"
)

type QRSession struct {
	AuthID      string
	RedirectURI url.URL
	Scope       SessionScope
}
