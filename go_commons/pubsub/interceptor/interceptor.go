package interceptor

import (
	"context"

	"github.com/omniful/go_commons/pubsub"
)

// Handler is an alias for handler function type
type Handler func(context.Context, *pubsub.Message) error

// Interceptor allows hooking additional functionality before or after
// producing / consuming an event
type Interceptor func(ctx context.Context, info *pubsub.Message, handler Handler) error

// NoOpInterceptor is the default interceptor that does nothing
var NoOpInterceptor Interceptor = func(ctx context.Context, msg *pubsub.Message, handler Handler) error {
	return handler(ctx, msg)
}

// ChainInterceptors (f1, f2, f3) returns a func g such that
//
//	g = f(args, handler) {
//	    f1(args, f2) {
//		       f2(args, f3) {
//	             f3(args, handler)
//	        }
//	    }
//	}
func ChainInterceptors(interceptors ...Interceptor) Interceptor {
	n := len(interceptors)

	return func(ctx context.Context, info *pubsub.Message, handler Handler) error {
		chainer := func(currentInter Interceptor, currentHandler Handler) Handler {
			return func(currentContext context.Context, currentInfo *pubsub.Message) error {
				return currentInter(currentContext, currentInfo, currentHandler)
			}
		}

		chainedHandler := handler
		for i := n - 1; i >= 0; i-- {
			chainedHandler = chainer(interceptors[i], chainedHandler)
		}

		return chainedHandler(ctx, info)
	}
}
