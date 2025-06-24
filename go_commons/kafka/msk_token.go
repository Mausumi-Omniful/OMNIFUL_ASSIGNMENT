package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/aws/aws-msk-iam-sasl-signer-go/signer"
)

type MSKAccessTokenProvider struct {
	region string
}

func (m *MSKAccessTokenProvider) Token() (*sarama.AccessToken, error) {
	token, _, err := signer.GenerateAuthToken(context.TODO(), m.region)
	return &sarama.AccessToken{Token: token}, err
}
