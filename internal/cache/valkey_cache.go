package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis"
)

// ValkeyMarketCache implements Redis/Valkey-based caching
type ValkeyMarketCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewValkeyMarketCache(addr string) (*ValkeyMarketCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr, // e.g., "localhost:6379"
		Password: "",   // no password set
		DB:       0,    // use default DB
	})

	ctx := context.Background()

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Valkey: %w", err)
	}

	return &ValkeyMarketCache{
		client: client,
		ctx:    ctx,
	}, nil
}

func (v *ValkeyMarketCache) SetMarkets(exchange string, markets interface{}, ttl time.Duration) error {
	data, err := json.Marshal(markets)
	if err != nil {
		return fmt.Errorf("failed to marshal markets: %w", err)
	}

	key := fmt.Sprintf("markets:%s", exchange)
	if err := v.client.Set(v.ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	log.Printf("[Cache] Saved %s markets to Valkey (TTL: %v)", exchange, ttl)
	return nil
}

func (v *ValkeyMarketCache) GetMarkets(exchange string, markets interface{}) error {
	key := fmt.Sprintf("markets:%s", exchange)

	data, err := v.client.Get(v.ctx, key).Bytes()
	if err == redis.Nil {
		return fmt.Errorf("cache miss")
	} else if err != nil {
		return fmt.Errorf("failed to get cache: %w", err)
	}

	if err := json.Unmarshal(data, markets); err != nil {
		return fmt.Errorf("failed to unmarshal cache: %w", err)
	}

	log.Printf("[Cache] Loaded %s markets from Valkey cache", exchange)
	return nil
}

func (v *ValkeyMarketCache) Exists(exchange string) (bool, error) {
	key := fmt.Sprintf("markets:%s", exchange)
	exists, err := v.client.Exists(v.ctx, key).Result()
	return exists > 0, err
}

func (v *ValkeyMarketCache) Close() error {
	return v.client.Close()
}
