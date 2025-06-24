package response

import (
	"encoding/json"
	"github.com/omniful/go_commons/http"
)

// ExecResponse is a success response returned by go commons lambda client
type ExecResponse struct {
	ExecutedVersion string
	StatusCode      http.StatusCode
	Data            interface{}
	MetaData        map[string]interface{}
}

// LambdaRes is response struct returned aws lambda handler
type LambdaRes struct {
	StatusCode int                    `json:"status_code"`
	Data       json.RawMessage        `json:"data"`
	MetaData   map[string]interface{} `json:"meta_data"`
}

type ErrorLambdaData struct {
	ErrorCode    string `json:"error_code"` // Custom error code
	ErrorMessage string `json:"error_message"`
}

// IsSuccess method returns true if HTTP status is 2xx else false
func (r *ExecResponse) IsSuccess() bool {
	return r.StatusCode.Is2xx()
}
