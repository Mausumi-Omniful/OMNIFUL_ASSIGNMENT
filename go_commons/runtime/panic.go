package runtime

import (
	"github.com/omniful/go_commons/log"
	"runtime/debug"
)

func HandleCrash(additionalHandlers ...func(interface{})) {
	if r := recover(); r != nil {
		log.Errorf("Panic - %v\n%s", r, string(debug.Stack()))
		for _, fn := range additionalHandlers {
			fn(r)
		}
	}
}
