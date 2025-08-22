// Package mapstorage contains implementation of key-value storage with ordinary map type.
package mapstorage

import (
	"context"
	"nova/internal/storage"
	ds "nova/pkg/datastructures"
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
	mu sync.RWMutex

	data            map[string]item
	cleanupInterval time.Duration

	lists map[string]*ds.LinkedList
}

func New(ctx context.Context, opts ...Option) *Storage {
	storage := &Storage{
		data:            map[string]item{},
		lists:           map[string]*ds.LinkedList{},
		cleanupInterval: defaultCleanupInterval,
	}

	for _, opt := range opts {
		opt(storage)
	}

	go storage.cleanup(ctx)

	return storage
}

// Set adds key-value pair with specified time-to-live.
// If there is a record of different data type with this key, it would be deleted.
func (s *Storage) Set(key, value string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// if ttl is not specified (equal nil), expiresAt would be nil
	// which means that key-value record doesn't have expiration time
	expiresAt := time.Time{}
	if ttl != 0 {
		expiresAt = time.Now().Add(ttl)
	}

	// delete list with similar key. There can be only one
	delete(s.lists, key)

	// add new item
	s.data[key] = item{
		value:     value,
		expiresAt: expiresAt,
	}
}

// Get returns value via given key. If there is no such value, ErrKeyNotFound is returned.
func (s *Storage) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, ok := s.lists[key]; ok {
		return "", storage.ErrWrongType
	}

	item, ok := s.data[key]
	if !ok || isExpired(item) {
		return "", storage.ErrKeyNotFound
	}

	return item.value, nil
}

// DeleteMany deletes all records with specified keys. Returns count of deleted records
func (s *Storage) DeleteMany(keys []string) int {
	count := 0

	s.mu.Lock()

	// only one of if-blocks would be executed
	for _, key := range keys {
		if _, ok := s.lists[key]; ok {
			delete(s.lists, key)
			count++
		}

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

// RPush adds new elements to the end of the list available via given key.
// It returns length of list after addition. If there is not such list, it is created.
func (s *Storage) RPush(key string, values []string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if el, ok := s.data[key]; ok && !isExpired(el) {
		return 0, storage.ErrWrongType
	}

	if _, ok := s.lists[key]; !ok {
		s.lists[key] = ds.NewLinkedList()
	}

	var length int
	for _, value := range values {
		length = s.lists[key].PushBack(value)
	}
	return length, nil
}

// RPush adds new elements to the beginning of the list available via given key.
// It returns length of list after addition. If there is not such list, it is created.
func (s *Storage) LPush(key string, values []string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if el, ok := s.data[key]; ok && !isExpired(el) {
		return 0, storage.ErrWrongType
	}

	if _, ok := s.lists[key]; !ok {
		s.lists[key] = ds.NewLinkedList()
	}

	var length int
	for _, value := range values {
		length = s.lists[key].PushForward(value)
	}
	return length, nil
}

// LRange returns node values in range of indexes [start, stop].
func (s *Storage) LRange(key string, start, stop int) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RLock()

	if el, ok := s.data[key]; ok && !isExpired(el) {
		return nil, storage.ErrWrongType
	}

	if _, ok := s.lists[key]; !ok {
		return nil, storage.ErrKeyNotFound
	}

	values := s.lists[key].LRange(start, stop)
	return values, nil
}
