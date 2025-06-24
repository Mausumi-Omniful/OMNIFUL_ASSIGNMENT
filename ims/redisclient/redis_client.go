package redisclient

import (
	"context"
	"fmt"
	"time"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/redis"
)

var Client *redis.Client

func InitRedis(ctx context.Context) error {
	host := config.GetString(ctx, "REDIS_HOST")
	port := config.GetString(ctx, "REDIS_PORT")
	addr := host
	if port != "" {
		addr = fmt.Sprintf("%s:%s", host, port)
	}
	if addr == ":" || addr == "" {
		addr = "localhost:6379"
	}

	config := &redis.Config{
		Hosts:        []string{addr},
		PoolSize:     50,
		MinIdleConn:  10,
		DialTimeout:  500 * time.Millisecond,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		IdleTimeout:  600 * time.Second,
	}

	Client = redis.NewClient(config)
	ctxBg := context.Background()

	success, err := Client.Set(ctxBg, "test_connection", "ping", 10*time.Second)
	if err != nil {
		return fmt.Errorf("Redis connection test failed: %v", err)
	}

	if success {
		Client.Del(ctxBg, "test_connection")
		fmt.Println("Redis connection test successful")
		return nil
	} else {
		return fmt.Errorf("Redis connection test failed")
	}
}

func Close() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}