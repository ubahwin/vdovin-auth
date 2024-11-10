package model

import "net/url"

type QRSession struct {
	RedirectURI url.URL
	Scope       SessionScope
}
