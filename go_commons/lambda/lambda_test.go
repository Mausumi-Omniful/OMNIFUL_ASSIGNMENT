package lambda

import (
	"context"
	"github.com/omniful/go_commons/lambda/request"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAWSLambdaClient(t *testing.T) {

	t.Run("Default Configuration", func(t *testing.T) {
		t.Parallel()

		client, err := NewAWSLambdaClient(context.TODO())
		assert.NoError(t, err, "should not error when initializing the lambda client")
		assert.NotNil(t, client, "client should be initialized")
		assert.Equal(t, "local", client.GetPrefix(), "default environment should be local")
	})

	t.Run("WithClient Option", func(t *testing.T) {
		t.Parallel()

		client, err := NewAWSLambdaClient(context.TODO(), WithPrefix("staging"))
		assert.NoError(t, err, "should not error when setting client environment")
		assert.NotNil(t, client, "client should be initialized")
		assert.Equal(t, "staging", client.GetPrefix(), "environment should be set to staging")
	})

	t.Run("Nil Client", func(t *testing.T) {
		t.Parallel()

		var client *Client

		var res interface{}
		resp, errResp := client.Execute(context.TODO(), &request.ExecRequest{}, res)

		assert.Nil(t, resp, "response should be nil if client is nil")
		assert.NotNil(t, errResp, "error response should not be nil if client is nil")
		assert.Equal(t, 400, int(errResp.StatusCode), "status code should be 400 for bad request error")
	})

}
