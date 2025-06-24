package auth

import (
	"encoding/base64"
	"fmt"
	"github.com/omniful/go_commons/httpclient"
)

func NewBasicAuth(username, password string) httpclient.Auth {
	return basic{username, password}
}

type basic struct {
	username string
	password string
}

func (a basic) AuthScheme() string { return "Basic" }
func (a basic) AuthToken() string {
	s := fmt.Sprintf("%s:%s", a.username, a.password)
	return base64.StdEncoding.EncodeToString([]byte(s))
}
