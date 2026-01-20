package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

func InitRedis(url string) (*redis.Client, error) {
	log.Println("Initializing Redis connection", "url", maskRedisURL(url))

	opt, err := redis.ParseURL(url)
	if err != nil {
		log.Println("Failed to parse Redis URL", "error", err, "url", maskRedisURL(url))
		return nil, fmt.Errorf("Failed to parse Redis URL: %w", err)
	}

	// Only enable TLS if the URL starts with 'rediss://' (secure redis)
	// This allows local 'redis://' to work without certificates
	if strings.HasPrefix(url, "rediss://") {
		log.Println("Enabling TLS for Redis connection")
		opt.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	// Set connection pool settings
	opt.PoolSize = 100
	opt.MinIdleConns = 10
	opt.MaxRetries = 3
	opt.DialTimeout = 5 * time.Second
	opt.ReadTimeout = 3 * time.Second
	opt.WriteTimeout = 3 * time.Second
	opt.PoolTimeout = 4 * time.Second

	client := redis.NewClient(opt)

	log.Println("Redis client configured",
		"pool_size", opt.PoolSize,
		"min_idle_conns", opt.MinIdleConns,
		"max_retries", opt.MaxRetries,
		"dial_timeout", opt.DialTimeout,
		"read_timeout", opt.ReadTimeout,
		"write_timeout", opt.WriteTimeout,
		"pool_timeout", opt.PoolTimeout,
		"use_tls", strings.HasPrefix(url, "rediss://"))

	// Test connection
	ctx := context.Background()
	start := time.Now()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Println("Redis connection failed",
			"error", err,
			"connection_time", time.Since(start),
			"url", maskRedisURL(url))
		return nil, fmt.Errorf("Redis connection failed: %w", err)
	}

	log.Println("Redis connected successfully",
		"connection_time", time.Since(start),
		"address", opt.Addr,
		"database", opt.DB)

	return client, nil
}
