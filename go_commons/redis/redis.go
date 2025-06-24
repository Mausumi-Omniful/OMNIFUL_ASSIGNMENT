package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/newrelic/go-agent/v3/integrations/nrredis-v8"
	"github.com/omniful/go_commons/log"
	"time"
)

type Cmdable interface {
	redis.Cmdable
	Close() error
}

// Client represents a Redis Client Connection
type Client struct {
	Nil error
	Cmdable
}

// Config defines the config structure for Redis
type Config struct {
	// Whether redis is running in cluster mode. If not, only a single host config
	// is expected in Hosts
	ClusterMode bool

	// Send read commands to slave nodes
	ServeReadsFromSlaves bool

	// Allows routing read-only commands randomly to master or slave node.
	// It automatically enables ReadOnly.
	ServeReadsFromMasterAndSlaves bool

	// Maximum number of socket connections
	// For Cluster Mode, the connections are per node and not the whole cluster
	PoolSize uint

	// Type of connection pool.
	// true for FIFO pool, false for LIFO pool.
	PoolFIFO bool

	// Minimum number of idle connections which is useful when establishing
	// new connection is slow. If 0, a default value is used.
	MinIdleConn uint

	// Database to be selected after connecting to the server
	// applicable in non cluster mode
	DB uint

	Hosts []string

	// Timeout for establishing the connection with the server
	DialTimeout time.Duration

	// Timeout for read requests
	ReadTimeout time.Duration

	// Timeout for write requests
	WriteTimeout time.Duration

	// By Default redis never closes the connection even if the client is idle
	// IdleTimeout can be used to configure the closing time for such connections
	IdleTimeout time.Duration
}

// NewClient returns a new Redis Client instance
func NewClient(cfg *Config) *Client {
	updateConfigWithDefaultValues(cfg)

	var redisdb Cmdable

	if cfg.ClusterMode {
		redisClusterClient := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:         cfg.Hosts,
			PoolSize:      int(cfg.PoolSize),
			PoolFIFO:      cfg.PoolFIFO,
			MinIdleConns:  int(cfg.MinIdleConn),
			ReadOnly:      cfg.ServeReadsFromSlaves,
			RouteRandomly: cfg.ServeReadsFromMasterAndSlaves,
			DialTimeout:   cfg.DialTimeout,
			ReadTimeout:   cfg.ReadTimeout,
			WriteTimeout:  cfg.WriteTimeout,
			IdleTimeout:   cfg.IdleTimeout,
		})

		// Adding new-relic tracking
		redisClusterClient.AddHook(nrredis.NewHook(nil))
		redisdb = redisClusterClient
	} else {
		if len(cfg.Hosts) > 1 || len(cfg.Hosts) <= 0 {
			log.Panicf("expecting a single host config for non cluster mode, found %d", len(cfg.Hosts))
		}

		opts := &redis.Options{
			Addr:         cfg.Hosts[0],
			PoolSize:     int(cfg.PoolSize),
			PoolFIFO:     cfg.PoolFIFO,
			DB:           int(cfg.DB),
			MinIdleConns: int(cfg.MinIdleConn),
		}

		// Adding new-relic tracking
		redisClient := redis.NewClient(opts)
		redisClient.AddHook(nrredis.NewHook(nil))
		redisdb = redisClient
	}

	client := &Client{
		Nil:     redis.Nil,
		Cmdable: redisdb,
	}

	return client
}

func updateConfigWithDefaultValues(cfg *Config) {
	if cfg.PoolSize == 0 {
		cfg.PoolSize = 50
	}

	if cfg.MinIdleConn == 0 {
		cfg.MinIdleConn = 6
	}

	if cfg.DialTimeout == 0 {
		cfg.DialTimeout = 500 * time.Millisecond
	}

	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = 2000 * time.Millisecond
	}

	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = 2000 * time.Millisecond
	}

	if cfg.IdleTimeout == 0 {
		cfg.IdleTimeout = 600 * time.Second
	}
}
