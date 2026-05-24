package auth

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

const jwtVersionPrefix = "auth:user:jwt:version:"

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(host, port string) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
	})

	return &RedisClient{client: client}
}

func (r *RedisClient) GetJwtVersion(userID string) (int, error) {
	key := jwtVersionPrefix + userID
	val, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get jwt version: %w", err)
	}

	version, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("invalid jwt version value: %w", err)
	}

	return version, nil
}

func (r *RedisClient) Ping() error {
	return r.client.Ping(context.Background()).Err()
}
