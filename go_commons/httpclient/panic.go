package httpclient

import "github.com/omniful/go_commons/log"

type PanicHandler func(err error)

var PanicLogger PanicHandler = func(err error) {
	log.Errorf("Panicked from error: %v", err)
}
