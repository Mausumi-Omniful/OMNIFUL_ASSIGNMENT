package interservice_client

import (
	"github.com/omniful/go_commons/http"
	libhttp "net/http"
	"time"
)

type Config struct {
	ServiceName string
	BaseURL     string
	Timeout     time.Duration
	Transport   *libhttp.Transport
}

func HTTPTransport() *libhttp.Transport {
	return &libhttp.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     60 * time.Second,
		DisableCompression:  true,
	}
}

func NewClientWithConfig(config Config) (*Client, error) {
	transport := HTTPTransport()
	if config.Transport != nil {
		transport = config.Transport
	}

	c, err := http.NewHTTPClient(
		config.ServiceName,
		config.BaseURL,
		transport,
		http.WithTimeout(config.Timeout*time.Second),
	)
	if err != nil {
		return nil, err
	}

	return NewClient(c), nil
}
