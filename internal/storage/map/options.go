package mapstorage

import "time"

type Option func(*Storage)

func WithCleanupInterval(interval time.Duration) Option {
	return func(s *Storage) {
		s.cleanupInterval = interval
	}
}
