package mapstorage

import (
	"context"
	"time"
)

var (
	defaultCleanupInterval = 1 * time.Minute
)

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
