// Package mapstorage contains implementation of key-value storage with ordinary map type.
package mapstorage

import (
	"context"
	"sync"
	"time"
)

var (
	defaultCleanupInterval = 1 * time.Minute
)

// item represents a value in storage with expiration time.
type item struct {
	value     string
	expiresAt time.Time
}

type Storage struct {
	mu              sync.RWMutex
	data            map[string]item
	cleanupInterval time.Duration
}

func New(ctx context.Context, opts ...Option) *Storage {
	storage := &Storage{
		data:            map[string]item{},
		cleanupInterval: defaultCleanupInterval,
	}

	for _, opt := range opts {
		opt(storage)
	}

	go storage.cleanup(ctx)

	return storage
}

func (s *Storage) Set(key, value string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	expiresAt := time.Time{}
	if ttl != 0 {
		expiresAt = time.Now().Add(ttl)
	}
	s.data[key] = item{
		value:     value,
		expiresAt: expiresAt,
	}
}

func (s *Storage) Get(key string) (string, bool) {
	s.mu.RLock()
	item, ok := s.data[key]
	s.mu.RUnlock()
	if !ok || isExpired(item) {
		return "", false
	}

	return item.value, true
}

func (s *Storage) DeleteMany(keys []string) int {
	count := 0

	s.mu.Lock()
	for _, key := range keys {
		if el, ok := s.data[key]; ok && !isExpired(el) {
			delete(s.data, key)
			count++
		}
	}
	s.mu.Unlock()

	return count
}

// cleanup is a background worker which deletes expired values every (cleanupInterval) seconds.
func (s *Storage) cleanup(ctx context.Context) {
	ticker := time.NewTicker(s.cleanupInterval)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			expiredKeys := []string{}

			s.mu.RLock()
			for key, item := range s.data {
				if isExpired(item) {
					expiredKeys = append(expiredKeys, key)
				}
			}
			s.mu.RUnlock()

			s.mu.Lock()
			for _, key := range expiredKeys {
				delete(s.data, key)
			}
			s.mu.Unlock()
		}
	}
}

func isExpired(el item) bool {
	expiresAt := el.expiresAt
	return time.Now().After(expiresAt) && !expiresAt.IsZero()
}
