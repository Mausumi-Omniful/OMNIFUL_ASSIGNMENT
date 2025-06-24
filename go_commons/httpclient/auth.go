package httpclient

import (
	"fmt"
	"github.com/omniful/go_commons/constants"
)

// Auth defines authorization header
type Auth interface {
	AuthScheme() string
	AuthToken() string
}

type AuthRefreshScope string

func ToAuthHeaderValue(a Auth) string {
	return fmt.Sprintf("%s %s", a.AuthScheme(), a.AuthToken())
}

func ToAuthHeader(a Auth) map[string]string {
	return map[string]string{
		constants.HeaderAuthorization: ToAuthHeaderValue(a),
	}
}
