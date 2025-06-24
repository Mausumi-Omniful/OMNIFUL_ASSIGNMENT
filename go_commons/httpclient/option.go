package httpclient

import (
	"github.com/omniful/go_commons/constants"
	"github.com/omniful/go_commons/log"
	"net/http"
	"time"
)

type Option func(Config) Config
type Options []Option

func (opts Options) ToConfig() Config {
	cfg := defaultConfig()
	for _, opt := range opts {
		cfg = opt(cfg)
	}
	return cfg
}

func defaultConfig() Config {
	return Config{
		userAgent:               constants.Omniful,
		forceRequestIDInHeaders: true,
		panicHandler:            PanicLogger,
		maxRetries:              3,
	}
}

type LogConfig struct {
	Logger      *log.Logger
	LogLevel    string
	LogRequest  bool
	LogResponse bool
}

type Config struct {
	userAgent               string
	clientAuth              Auth
	requestAuthProvider     AuthProvider
	retryStrategies         []Retry
	rateLimiter             RateLimiter
	forceRequestIDInHeaders bool
	transport               http.RoundTripper
	panicHandler            PanicHandler
	maxRetries              int
	logConfig               *LogConfig
	contentType             string
	deadline                time.Duration
	timeout                 time.Duration

	// Callbacks
	beforeSendCallbacks    []BeforeSendCallback
	beforeAttemptCallbacks []BeforeAttemptCallback
	afterAttemptCallbacks  []AfterAttemptCallback
	afterSendCallbacks     []AfterSendCallback
	onError                OnErrorCallback
}

func WithUserAgent(ua string) Option {
	return func(c Config) Config {
		c.userAgent = ua
		return c
	}
}

func WithClientAuth(auth Auth) Option {
	return func(c Config) Config {
		c.clientAuth = auth
		return c
	}
}

func WithRequestAuthProvider(ap AuthProvider) Option {
	return func(c Config) Config {
		c.requestAuthProvider = ap
		return c
	}
}

func WithRetry(retry Retry) Option {
	return func(c Config) Config {
		c.retryStrategies = append(c.retryStrategies, retry)
		return c
	}
}

func WithRetryStrategies(rs []Retry) Option {
	return func(c Config) Config {
		c.retryStrategies = rs
		return c
	}
}

func WithRateLimiter(rl RateLimiter) Option {
	return func(c Config) Config {
		c.rateLimiter = rl
		return c
	}
}

func WithRequestIDInHeaders(v bool) Option {
	return func(c Config) Config {
		c.forceRequestIDInHeaders = v
		return c
	}
}

func WithTransport(tr http.RoundTripper) Option {
	return func(c Config) Config {
		c.transport = tr
		return c
	}
}

func WithPanicHandler(ph PanicHandler) Option {
	return func(c Config) Config {
		c.panicHandler = ph
		return c
	}
}

func WithLogConfig(lc LogConfig) Option {
	return func(c Config) Config {
		c.logConfig = &lc
		return c
	}
}

func WithContentType(ct string) Option {
	return func(c Config) Config {
		c.contentType = ct
		return c
	}
}

func WithBeforeSendCallback(cb BeforeSendCallback) Option {
	return func(c Config) Config {
		c.beforeSendCallbacks = append(c.beforeSendCallbacks, cb)
		return c
	}
}

func WithBeforeSendCallbacks(cbs []BeforeSendCallback) Option {
	return func(c Config) Config {
		c.beforeSendCallbacks = cbs
		return c
	}
}

func WithBeforeAttemptCallback(cb BeforeAttemptCallback) Option {
	return func(c Config) Config {
		c.beforeAttemptCallbacks = append(c.beforeAttemptCallbacks, cb)
		return c
	}
}

func WithBeforeAttemptCallbacks(cbs []BeforeAttemptCallback) Option {
	return func(c Config) Config {
		c.beforeAttemptCallbacks = cbs
		return c
	}
}

func WithAfterSendCallback(cb AfterSendCallback) Option {
	return func(c Config) Config {
		c.afterSendCallbacks = append(c.afterSendCallbacks, cb)
		return c
	}
}

func WithAfterSendCallbacks(cbs []AfterSendCallback) Option {
	return func(c Config) Config {
		c.afterSendCallbacks = cbs
		return c
	}
}

func WithAfterAttemptCallback(cb AfterAttemptCallback) Option {
	return func(c Config) Config {
		c.afterAttemptCallbacks = append(c.afterAttemptCallbacks, cb)
		return c
	}
}

func WithAfterAttemptCallbacks(cbs []AfterAttemptCallback) Option {
	return func(c Config) Config {
		c.afterAttemptCallbacks = cbs
		return c
	}
}

func WithOnErrorCallback(cb OnErrorCallback) Option {
	return func(c Config) Config {
		c.onError = cb
		return c
	}
}

func WithMaxRetries(i int) Option {
	return func(c Config) Config {
		c.maxRetries = i
		return c
	}
}

func WithDeadline(d time.Duration) Option {
	return func(c Config) Config {
		c.deadline = d
		return c
	}
}

func WithTimeout(d time.Duration) Option {
	return func(c Config) Config {
		c.timeout = d
		return c
	}
}
