package httpclient

import (
	"github.com/omniful/go_commons/constants"
	"github.com/omniful/go_commons/httpclient/request"
)

// AuthProvider provides authorization for a request
type AuthProvider interface {
	Apply(*Context, request.Request) (request.Request, error)
}

func AddAuthorizationHeader(req request.Request, a Auth) request.Request {
	req.GetHeaders().Add(constants.HeaderAuthorization, ToAuthHeaderValue(a))
	return req
}

/* START SIMPLE AUTH PROVIDER */

func NewSimpleAuthProvider(a Auth) AuthProvider {
	return &simpleAuthProvider{a: a}
}

type simpleAuthProvider struct {
	a Auth
}

func (p simpleAuthProvider) Apply(_ *Context, req request.Request) (request.Request, error) {
	return AddAuthorizationHeader(req, p.a), nil
}

/* END SIMPLE AUTH PROVIDER */
