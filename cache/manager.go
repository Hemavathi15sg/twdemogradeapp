package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheStats tracks cache performance metrics
type CacheStats struct {
	Hits    int64
	Misses  int64
	Updates int64
	Deletes int64
	mu      sync.RWMutex
}

// Manager handles Redis caching with statistics
type Manager struct {
	client *redis.Client
	ttl    time.Duration
	stats  *CacheStats
}

// NewManager creates a new cache manager
func NewManager(redisClient *redis.Client) *Manager {
	return &Manager{
		client: redisClient,
		ttl:    5 * time.Minute,
		stats: &CacheStats{
			Hits:    0,
			Misses:  0,
			Updates: 0,
			Deletes: 0,
		},
	}
}

// Get retrieves value from cache
func (m *Manager) Get(ctx context.Context, key string, target interface{}) (bool, error) {
	val, err := m.client.Get(ctx, key).Result()
	if err == redis.Nil {
		m.recordMiss()
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if err := json.Unmarshal([]byte(val), target); err != nil {
		m.recordMiss()
		return false, nil
	}

	m.recordHit()
	return true, nil
}

// Set stores value in cache
func (m *Manager) Set(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	m.recordUpdate()
	return m.client.Set(ctx, key, data, m.ttl).Err()
}

// Del removes keys from cache
func (m *Manager) Del(ctx context.Context, keys ...string) error {
	if len(keys) > 0 {
		m.recordDelete()
		return m.client.Del(ctx, keys...).Err()
	}
	return nil
}

// recordHit increments hit counter
func (m *Manager) recordHit() {
	m.stats.mu.Lock()
	defer m.stats.mu.Unlock()
	m.stats.Hits++
}

// recordMiss increments miss counter
func (m *Manager) recordMiss() {
	m.stats.mu.Lock()
	defer m.stats.mu.Unlock()
	m.stats.Misses++
}

// recordUpdate increments update counter
func (m *Manager) recordUpdate() {
	m.stats.mu.Lock()
	defer m.stats.mu.Unlock()
	m.stats.Updates++
}

// recordDelete increments delete counter
func (m *Manager) recordDelete() {
	m.stats.mu.Lock()
	defer m.stats.mu.Unlock()
	m.stats.Deletes++
}

// GetStats returns current cache statistics
func (m *Manager) GetStats() map[string]int64 {
	m.stats.mu.RLock()
	defer m.stats.mu.RUnlock()

	total := m.stats.Hits + m.stats.Misses
	hitRate := int64(0)
	if total > 0 {
		hitRate = (m.stats.Hits * 100) / total
	}

	return map[string]int64{
		"hits":     m.stats.Hits,
		"misses":   m.stats.Misses,
		"updates":  m.stats.Updates,
		"deletes":  m.stats.Deletes,
		"total":    total,
		"hit_rate": hitRate,
	}
}

// LogStats logs current cache statistics
func (m *Manager) LogStats() string {
	stats := m.GetStats()
	return fmt.Sprintf(
		"Cache Stats - Hits: %d, Misses: %d, Hit Rate: %d%%, Updates: %d, Deletes: %d",
		stats["hits"],
		stats["misses"],
		stats["hit_rate"],
		stats["updates"],
		stats["deletes"],
	)
}
