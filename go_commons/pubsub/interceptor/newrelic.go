package interceptor

import (
	"context"
	"github.com/omniful/go_commons/constants"
	"github.com/omniful/go_commons/env"
	error2 "github.com/omniful/go_commons/error"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/newrelic"
	"github.com/omniful/go_commons/pubsub"
	"github.com/pkg/errors"
	"runtime/debug"
)

func NewRelicInterceptor() Interceptor {
	return func(ctx context.Context, msg *pubsub.Message, handler Handler) (err error) {
		transactionName, ok := ctx.Value(constants.KafkaConsumerTransactionName).(string)
		if !ok {
			transactionName = "event." + msg.Topic
		}

		txn := newrelic.StartTransaction(transactionName, nil, nil)
		nrCtx := newrelic.NewContext(ctx, txn)
		defer txn.End()

		newrelic.SetRequestID(nrCtx, env.GetKafkaRequestID(nrCtx, msg.Headers))

		defer func() {
			if r := recover(); r != nil {
				err = error2.NewCustomError(error2.PanicError,
					errors.Wrapf(errors.New(string(debug.Stack())), "[PANIC] %+v", r).Error())

				newrelic.NoticeError(nrCtx, err)

				log.Errorf("[PANIC] error recovered %s", err.Error())
			}
		}()

		err = handler(nrCtx, msg)
		if err != nil {
			newrelic.NoticeError(nrCtx, err)
		}
		return err
	}
}
