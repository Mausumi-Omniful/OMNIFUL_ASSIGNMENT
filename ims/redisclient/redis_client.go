package redisclient

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/omniful/go_commons/redis"
)

var Client *redis.Client

func InitRedis() error {
	host := os.Getenv("REDIS_ADDR")
	if host == "" {
		host = "localhost:6379" 
	}

	config := &redis.Config{
		Hosts:        []string{host},
		PoolSize:     50,
		MinIdleConn:  10,
		DialTimeout:  500 * time.Millisecond,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		IdleTimeout:  600 * time.Second,
	}

	Client = redis.NewClient(config)
	ctx := context.Background()

	success, err := Client.Set(ctx, "test_connection", "ping", 10*time.Second)
	if err != nil {
		return fmt.Errorf("Redis connection test failed: %v", err)
	}

	if success {
		Client.Del(ctx, "test_connection")
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
