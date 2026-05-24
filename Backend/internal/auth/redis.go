package auth

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

const permissionVersionPrefix = "auth:user:version:"
const revokedTokensPrefix = "refresh:revoked:"

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(host, port string) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
	})

	return &RedisClient{client: client}
}

func (r *RedisClient) CheckRefreshTokenRevokation(refreshToken string) bool {
	key := revokedTokensPrefix + refreshToken
	val := r.client.Get(context.Background(), key)

	if val != nil {
		return true
	}
	return false
}

func (r *RedisClient) GetPermissionVersion(userID string) (int, error) {
	key := permissionVersionPrefix + userID
	val, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get permission version: %w", err)
	}

	version, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("invalid permission version value: %w", err)
	}

	return version, nil
}

func (r *RedisClient) Ping() error {
	return r.client.Ping(context.Background()).Err()
}
