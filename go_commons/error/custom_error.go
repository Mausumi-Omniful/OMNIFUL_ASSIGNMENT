package error

import (
	"errors"
	"fmt"
	"strings"

	errorsPkg "github.com/pkg/errors"
)

type CustomError struct {
	errorCode     Code
	errorMsg      string
	data          interface{}
	errors        map[string]string
	error         error
	exists        bool
	retryable     bool
	shouldNotify  bool
	loggingParams map[string]interface{}
}

func NewCustomError(errorCode Code, error string, options ...func(*CustomError)) CustomError {
	c := CustomError{
		errorCode:    errorCode,
		errorMsg:     error,
		data:         nil,
		exists:       true,
		retryable:    false,
		shouldNotify: true,
	}

	for _, option := range options {
		option(&c)
	}

	e := errors.New(fmt.Sprintf("Code: %s | %s", c.errorCode, c.errorMsg))
	c.error = errorsPkg.WithStack(e)
	c.loggingParams = make(map[string]interface{}, 0)
	c.errors = make(map[string]string, 0)
	return c
}

func NewCustomErrorWithPayload(errorCode Code, error string, data interface{}, options ...func(*CustomError)) CustomError {
	c := NewCustomError(errorCode, error, options...)
	c.data = data
	return c
}

func WithRetryable(retryable bool) func(*CustomError) {
	return func(c *CustomError) {
		c.retryable = retryable
	}
}

func WithShouldNotify(shouldNotify bool) func(*CustomError) {
	return func(c *CustomError) {
		c.shouldNotify = shouldNotify
	}
}

func RequestInvalidError(message string, options ...func(*CustomError)) CustomError {
	c := CustomError{
		errorCode:    RequestInvalid,
		errorMsg:     message,
		data:         nil,
		exists:       true,
		retryable:    false,
		shouldNotify: true,
	}
	e := errors.New(fmt.Sprintf("Code: %s | %s", c.errorCode, c.errorMsg))
	c.error = errorsPkg.WithStack(e)
	c.loggingParams = make(map[string]interface{}, 0)
	c.errors = make(map[string]string, 0)

	for _, option := range options {
		option(&c)
	}
	return c
}

func (c CustomError) Exists() bool {
	return c.exists
}

func (c CustomError) Log() {
	fmt.Println(c.ToString())
}

func (c CustomError) LoggingParams() map[string]interface{} {
	return c.loggingParams
}

func (c CustomError) ErrorCode() Code {
	return c.errorCode
}

func (c CustomError) ToError() error {
	return c.error
}

func (c CustomError) Error() string {
	return c.error.Error()
}

func (c CustomError) ErrorMessage() string {
	return c.errorMsg
}

func (c CustomError) ShouldNotify() bool {
	return c.shouldNotify
}

func (c CustomError) Retryable() bool {
	return c.retryable
}

func (c CustomError) ToString() string {
	logMsg := fmt.Sprintf("Code: %s, Msg: %s", c.errorCode, c.errorMsg)

	paramStrings := make([]string, 0)
	for key, val := range c.loggingParams {
		paramStrings = append(paramStrings, fmt.Sprintf("%s: {%+v}", strings.ToUpper(key), val))
	}
	return fmt.Sprintf("%s, Params: [%+v]", logMsg, strings.Join(paramStrings, " | "))
}

// WithParam value param should not be a pointer
func (c CustomError) WithParam(key string, val interface{}) CustomError {
	if c.loggingParams == nil {
		c.loggingParams = make(map[string]interface{}, 0)
	}
	c.loggingParams[key] = val
	return c
}

func (c CustomError) ErrorString() string {
	return c.errorMsg
}

func (c CustomError) UserMessage() string {
	return c.errorMsg
}

func (c CustomError) ErrorData() interface{} {
	return c.data
}

func (c CustomError) ErrorMap() map[string]string {
	return c.errors
}

// WithErrors allows setting the entire errors map
func WithErrors(errors map[string]string) func(*CustomError) {
	return func(c *CustomError) {
		c.errors = errors
	}
}

// WithData allows setting the data field
func WithData(data any) func(*CustomError) {
	return func(c *CustomError) {
		c.data = data
	}
}
