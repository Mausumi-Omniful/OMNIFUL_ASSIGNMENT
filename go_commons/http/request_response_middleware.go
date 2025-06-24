package http

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/constants"
	"github.com/omniful/go_commons/env"
	"github.com/omniful/go_commons/log"
)

type LoggingMiddlewareOptions struct {
	Format      string
	Level       string
	LogRequest  bool
	LogResponse bool
	LogHeader   bool
}

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func RequestLogMiddleware(opts LoggingMiddlewareOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Ignore Health Requests
		if path == "/health" {
			c.Next()

			return
		}

		query := c.Request.URL.RawQuery
		reqID := env.GetRequestID(c)
		start := time.Now()
		requestBodyString := "<Disabled>"
		bodyWriter := &responseWriter{
			body: bytes.NewBufferString("<Disabled>"),
		}

		l := log.DefaultLogger()
		l = l.With(
			log.String(constants.HeaderXOmnifulRequestID, reqID),
		)

		ctx := log.ContextWithLogger(c, l).(*gin.Context)

		// Create a custom ResponseWriter to capture the response body
		if opts.LogResponse {
			bodyWriter.body = bytes.NewBufferString("")
			bodyWriter.ResponseWriter = c.Writer
			ctx.Writer = bodyWriter
		}

		if opts.LogRequest {
			requestBody, _ := ctx.GetRawData()
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
			requestBodyString = string(requestBody)
		}

		defer func() {
			// Capture the request complete timestamp
			end := time.Now()

			logFields := []log.Field{
				log.String("Amzn-Trace-ID", ctx.GetHeader("x-amzn-trace-id")),
				log.String("Method", ctx.Request.Method),
				log.String("Domain", ctx.Request.Host),
				log.String("Path", path),
				log.Int("Status", ctx.Writer.Status()),
				log.String("RequestBody", requestBodyString),
				log.String("Query", query),
				log.String("ResponseBody", bodyWriter.body.String()),
				log.String("IP", ctx.ClientIP()),
				log.String("User-Agent", ctx.Request.UserAgent()),
				log.Duration("Latency", time.Since(start)),
				log.String("RequestReceivedAt", start.Format(time.RFC3339)),
				log.String("RequestCompletedAt", end.Format(time.RFC3339)),
			}

			if correlationID := ctx.GetHeader(constants.HeaderXOmnifulCorrelationID); len(correlationID) > 0 {
				logFields = append(logFields, log.String(constants.HeaderXOmnifulCorrelationID, correlationID))
			}

			if clientService := ctx.GetHeader(constants.HeaderXClientService); len(clientService) > 0 {
				logFields = append(logFields, log.String(constants.HeaderXClientService, clientService))
			}

			if opts.LogHeader {
				for k, v := range ctx.Request.Header {
					if len(v) > 0 {
						logFields = append(logFields, log.String(k, v[0]))
					}
				}
			}

			l.With(logFields...).Info(path)
		}()

		ctx.Next()
	}
}
