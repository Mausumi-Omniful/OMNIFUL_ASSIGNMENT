package json

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/omniful/go_commons/env"
	commmonError "github.com/omniful/go_commons/error"
	"github.com/omniful/go_commons/newrelic"
)

// Unmarshal is a wrapper for json.Unmarshal that reports errors to New Relic
func Unmarshal(ctx context.Context, data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		errorMsg := fmt.Sprintf("%s | Unmarshal Error: %s", env.GetRequestID(ctx), err.Error())
		newrelic.NoticeCustomError(ctx, commmonError.NewCustomError(commmonError.JsonDeserializationError, errorMsg))
		return err
	}

	return nil
}

// Marshal is a wrapper for json.Marshal that reports errors to New Relic
func Marshal(ctx context.Context, v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		errorMsg := fmt.Sprintf("%s | Marshal Error: %s", env.GetRequestID(ctx), err.Error())
		newrelic.NoticeCustomError(ctx, commmonError.NewCustomError(commmonError.JsonSerializationError, errorMsg))
		return nil, err
	}

	return data, nil
}
