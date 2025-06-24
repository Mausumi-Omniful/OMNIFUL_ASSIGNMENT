package auth

import "github.com/omniful/go_commons/httpclient"

func NewBearerAuth(token string) httpclient.Auth {
	return bearer{token}
}

type bearer struct {
	token string
}

func (a bearer) AuthScheme() string { return "Bearer" }
func (a bearer) AuthToken() string  { return a.token }
