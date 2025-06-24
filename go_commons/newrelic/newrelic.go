package newrelic

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/newrelic/go-agent/v3/integrations/nrpkgerrors"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/omniful/go_commons/constants"
	customError "github.com/omniful/go_commons/error"
	"github.com/omniful/go_commons/shutdown"
)

type Event string

var Error Event = "Error"

var (
	application *newrelic.Application = nil
)

type Options struct {
	Name                   string
	License                string
	Enabled                bool
	CrossApplicationTracer bool
	DistributedTracer      bool
}

type Attribute struct {
	Name  string
	Value interface{}
}

// Initialize newrelic
func Initialize(op *Options) *newrelic.Application {
	if !op.Enabled {
		return nil
	}

	var err error
	application, err = newrelic.NewApplication(
		newrelic.ConfigAppName(op.Name),
		newrelic.ConfigLicense(op.License),
		newrelic.ConfigDistributedTracerEnabled(op.DistributedTracer),
		newrelic.ConfigEnabled(op.Enabled),
		func(c *newrelic.Config) {
			c.CrossApplicationTracer.Enabled = op.CrossApplicationTracer
		},
		func(c *newrelic.Config) {
			c.ErrorCollector.RecordPanics = true
		},
	)

	if err != nil {
		panic("Could not initialize newrelic: " + err.Error())
	}

	// Registering for shutdown
	shutdown.RegisterShutdownCallback(constants.Newrelic, &Newrelic{application})

	return application
}

// GetApplication Returns newrelic application
func GetApplication() *newrelic.Application {
	return application
}

func StartSegmentWithContext(ctx context.Context, name string) *newrelic.Segment {
	txn := newrelic.FromContext(ctx)
	if txn == nil {
		return nil
	}
	return txn.StartSegment(name)
}

func StartSegment(name string, txn *newrelic.Transaction) *newrelic.Segment {
	return txn.StartSegment(name)
}

func StartTransaction(name string, w http.ResponseWriter, r *http.Request) *newrelic.Transaction {
	txn := application.StartTransaction(name)
	if r != nil {
		txn.SetWebRequestHTTP(r)
	}

	if w != nil {
		txn.SetWebResponse(w)
	}

	return txn
}

func StartNonWebTransaction(name string) *newrelic.Transaction {
	return application.StartTransaction(name)
}

func AddAttributeWithContext(ctx context.Context, attrs ...Attribute) {
	txn := newrelic.FromContext(ctx)
	if application == nil || txn == nil {
		return
	}

	for _, attr := range attrs {
		txn.AddAttribute(attr.Name, attr.Value)
	}
}

func getURL(method, target string) *url.URL {
	var host string
	// target can be anything from
	// https://github.com/grpc/grpc/blob/master/doc/naming.md
	// see https://godoc.org/google.golang.org/grpc#DialContext
	if strings.HasPrefix(target, "unix:") {
		host = "localhost"
	} else {
		host = strings.TrimPrefix(target, "dns:///")
	}
	return &url.URL{
		Scheme: "grpc",
		Host:   host,
		Path:   method,
	}
}

func NoticeCustomError(ctx context.Context, err customError.CustomError) {
	txn := newrelic.FromContext(ctx)
	if txn == nil {
		return
	}

	txn.NoticeError(newrelic.Error{
		Message: err.ErrorMessage(),
		Class:   string(err.ErrorCode()),
	})
}

func NoticeError(ctx context.Context, err error) {
	txn := newrelic.FromContext(ctx)
	if txn == nil {
		return
	}
	// If it's not a custom error, use the default error wrapping
	txn.NoticeError(nrpkgerrors.Wrap(err))
}

func NoticeExpectedCustomError(ctx context.Context, err customError.CustomError) {
	txn := newrelic.FromContext(ctx)
	if txn == nil {
		return
	}

	txn.NoticeExpectedError(newrelic.Error{
		Message: err.ErrorMessage(),
		Class:   string(err.ErrorCode()),
	})
}

func NoticeExpectedError(ctx context.Context, err error) {
	txn := newrelic.FromContext(ctx)
	if txn == nil {
		return
	}

	// If it's not a custom error, use the default error wrapping
	txn.NoticeExpectedError(nrpkgerrors.Wrap(err))
}

// RecordEvent Each value in the properties map must be a number, string, or boolean.
func RecordEvent(event Event, properties map[string]interface{}) {
	application.RecordCustomEvent(string(event), properties)
}

func StartCustomDataSegment(ctx context.Context, product string, operation string) *newrelic.DatastoreSegment {
	productName := newrelic.DatastoreProduct(product)
	txn := newrelic.FromContext(ctx)
	return &newrelic.DatastoreSegment{
		StartTime: txn.StartSegmentNow(),
		Product:   productName,
		Operation: operation,
	}
}

func NewContext(ctx context.Context, txn *newrelic.Transaction) context.Context {
	return newrelic.NewContext(ctx, txn)
}

func FromContext(ctx context.Context) *newrelic.Transaction {
	return newrelic.FromContext(ctx)
}

func GetContext(txn *newrelic.Transaction) context.Context {
	return newrelic.NewContext(context.Background(), txn)
}

func SetRequestID(ctx context.Context, requestID string) {
	attribute := Attribute{
		Name:  constants.HeaderXOmnifulRequestID,
		Value: requestID,
	}

	AddAttributeWithContext(ctx, attribute)
}
