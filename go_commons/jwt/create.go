package jwt

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
)

func CreateJWTTokenWithSignature(method jwt.SigningMethod, claims jwt.Claims, privateKey *rsa.PrivateKey) (t string, err error) {
	token := jwt.NewWithClaims(method, claims)

	t, err = token.SignedString(privateKey)
	if err != nil {
		return
	}

	return
}

func CreateJWTTokenWithoutSignature(claims jwt.Claims) (t string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)

	t, err = token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		return
	}

	return
}
