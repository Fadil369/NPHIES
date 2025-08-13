package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

// Manager handles caching operations
type Manager struct {
	client *redis.Client
	ttl    time.Duration
}

// NewManager creates a new cache manager
func NewManager(client *redis.Client, ttl time.Duration) *Manager {
	return &Manager{
		client: client,
		ttl:    ttl,
	}
}

// Get retrieves a value from cache
func (m *Manager) Get(ctx context.Context, key string) (string, error) {
	return m.client.Get(ctx, key).Result()
}

// Set stores a value in cache with TTL
func (m *Manager) Set(ctx context.Context, key, value string) error {
	return m.client.Set(ctx, key, value, m.ttl).Err()
}

// SetWithTTL stores a value in cache with custom TTL
func (m *Manager) SetWithTTL(ctx context.Context, key, value string, ttl time.Duration) error {
	return m.client.Set(ctx, key, value, ttl).Err()
}

// Delete removes a value from cache
func (m *Manager) Delete(ctx context.Context, key string) error {
	return m.client.Del(ctx, key).Err()
}

// GetJSON retrieves and unmarshals a JSON value from cache
func (m *Manager) GetJSON(ctx context.Context, key string, dest interface{}) error {
	val, err := m.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// SetJSON marshals and stores a JSON value in cache
func (m *Manager) SetJSON(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return m.client.Set(ctx, key, data, m.ttl).Err()
}

// Exists checks if a key exists in cache
func (m *Manager) Exists(ctx context.Context, key string) (bool, error) {
	count, err := m.client.Exists(ctx, key).Result()
	return count > 0, err
}

// TTL returns the TTL of a key
func (m *Manager) TTL(ctx context.Context, key string) (time.Duration, error) {
	return m.client.TTL(ctx, key).Result()
}

// Expire sets a new TTL for a key
func (m *Manager) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return m.client.Expire(ctx, key, ttl).Err()
}

// FlushAll clears all cache entries
func (m *Manager) FlushAll(ctx context.Context) error {
	return m.client.FlushAll(ctx).Err()
}

// Keys returns all keys matching a pattern
func (m *Manager) Keys(ctx context.Context, pattern string) ([]string, error) {
	return m.client.Keys(ctx, pattern).Result()
}

// DeletePattern deletes all keys matching a pattern
func (m *Manager) DeletePattern(ctx context.Context, pattern string) error {
	keys, err := m.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	
	if len(keys) == 0 {
		return nil
	}
	
	return m.client.Del(ctx, keys...).Err()
}

// IncrementCounter increments a counter key
func (m *Manager) IncrementCounter(ctx context.Context, key string) (int64, error) {
	return m.client.Incr(ctx, key).Result()
}

// GetCounter gets the value of a counter key
func (m *Manager) GetCounter(ctx context.Context, key string) (int64, error) {
	result := m.client.Get(ctx, key)
	if result.Err() == redis.Nil {
		return 0, nil
	}
	if result.Err() != nil {
		return 0, result.Err()
	}
	
	return result.Int64()
}

// SetCounter sets a counter value
func (m *Manager) SetCounter(ctx context.Context, key string, value int64) error {
	return m.client.Set(ctx, key, value, m.ttl).Err()
}

// GetInfo returns cache information
func (m *Manager) GetInfo(ctx context.Context) (*redis.StringCmd, error) {
	return m.client.Info(ctx), nil
}