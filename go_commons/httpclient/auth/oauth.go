package auth

import "github.com/omniful/go_commons/httpclient"

func NewOAuth() httpclient.Auth {
	// TODO: implement
	return oauth{}
}

type oauth struct{}

func (a oauth) AuthScheme() string { return "Bearer" }
func (a oauth) AuthToken() string  { return "" }
