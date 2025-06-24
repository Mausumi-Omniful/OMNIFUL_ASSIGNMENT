package public

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

func GetRSAPublicKey(ctx context.Context, publicKey string) (rsaPubKey *rsa.PublicKey, err error) {
	pubPem, _ := pem.Decode([]byte(publicKey))

	if pubPem.Type != "PUBLIC KEY" {
		err = errors.New(fmt.Sprintf("RSA public key is of the wrong type, Pem Type :%s", pubPem.Type))
		return
	}

	parsedKey, parseErr := x509.ParsePKIXPublicKey(pubPem.Bytes)
	if parseErr != nil {
		err = errors.New(fmt.Sprintf("RSA public key is of the wrong type, Pem Type :%s", pubPem.Type))
		return
	}

	rsaPubKey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		err = errors.New("unable to parse pub key")
		return
	}

	return
}
