package redis

import (
	"context"
	"fmt"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/gomodule/redigo/redis"
	redisClient "github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redisClient.Client
	searchClient *redisearch.Client
}

// NewRedisClient creates a new RedisClient
func NewRedisClient(addr string, password string, db int) (*RedisClient, error) {
	client := redisClient.NewClient(&redisClient.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Test connection
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	pool := &redis.Pool{Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", addr, redis.DialPassword(password))
	}}

	// shop index
	searchClient := redisearch.NewClientFromPool(pool, "shopIdx")
	
	// todo: search results index?

	return &RedisClient{client: client, searchClient: searchClient}, nil
}