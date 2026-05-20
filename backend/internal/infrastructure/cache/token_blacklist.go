package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// TokenBlacklist 令牌黑名单服务
type TokenBlacklist struct {
	client *redis.Client
	prefix string
}

// NewTokenBlacklist 创建令牌黑名单服务
func NewTokenBlacklist(client *redis.Client) *TokenBlacklist {
	return &TokenBlacklist{
		client: client,
		prefix: "blacklist:token:",
	}
}

// Add 将令牌加入黑名单
func (t *TokenBlacklist) Add(ctx context.Context, token string, expiresAt time.Time) error {
	key := t.prefix + token
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return nil
	}
	return t.client.Set(ctx, key, "1", ttl).Err()
}

// IsBlacklisted 检查令牌是否在黑名单中
func (t *TokenBlacklist) IsBlacklisted(ctx context.Context, token string) bool {
	key := t.prefix + token
	result, err := t.client.Get(ctx, key).Result()
	if err != nil {
		return false
	}
	return result == "1"
}

// Remove 从黑名单移除令牌（用于刷新令牌场景，旧令牌失效后换新令牌）
func (t *TokenBlacklist) Remove(ctx context.Context, token string) error {
	key := t.prefix + token
	return t.client.Del(ctx, key).Err()
}
