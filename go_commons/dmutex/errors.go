package dmutex

import (
	"fmt"
)

type ErrFailedToLock struct {
	errs []error
}

func NewErrFailedToLock(errs ...error) error {
	return &ErrFailedToLock{errs}
}

func (e *ErrFailedToLock) Error() string {
	msg := "failed to lock"
	for _, err := range e.errs {
		msg = fmt.Sprintf("%s: %v", msg, err)
	}
	return msg
}

type ErrFailedToUnlock struct {
	errs []error
}

func NewErrFailedToUnlock(errs ...error) error {
	return &ErrFailedToUnlock{errs}
}

func (e *ErrFailedToUnlock) Error() string {
	msg := "failed to unlock"
	for _, err := range e.errs {
		msg = fmt.Sprintf("%s: %v", msg, err)
	}
	return msg
}
