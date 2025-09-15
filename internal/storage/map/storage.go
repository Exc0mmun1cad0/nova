// Package mapstorage contains implementation of key-value storage with ordinary map type.
package mapstorage

import (
	"context"
	"nova/internal/storage"
	ds "nova/pkg/datastructures"
	"strconv"
	"sync"
	"time"
)

// ValueType represents data type of item in storage
type ValueType int

const (
	ValueTypeString ValueType = iota
	ValueTypeInt
	ValueTypeList
)

// item represents a value in storage with expiration time.
type item struct {
	valueType ValueType
	value     any
	expiresAt time.Time
}

type Storage struct {
	mu sync.RWMutex

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

// Set adds key-value pair with specified time-to-live.
// If there is a record of different data type with this key, it would be deleted.
func (s *Storage) Set(key, value string, ttl time.Duration) {
	// if ttl is not specified (equal nil), expiresAt would be nil
	// which means that key-value record doesn't have expiration time
	expiresAt := time.Time{}
	if ttl != 0 {
		expiresAt = time.Now().Add(ttl)
	}

	valType := ValueTypeString
	var valToSet any = value
	if num, err := strconv.Atoi(value); err == nil {
		valType = ValueTypeInt
		valToSet = num
	}

	s.mu.Lock()

	// add new item
	s.data[key] = item{
		valueType: valType,
		value:     valToSet,
		expiresAt: expiresAt,
	}

	s.mu.Unlock()
}

// Get returns value via given key. If there is no such value, ErrKeyNotFound is returned.
func (s *Storage) Get(key string) (string, error) {
	s.mu.RLock()
	item, ok := s.data[key]
	s.mu.RUnlock()

	// cannot be executed with list type
	if ok && item.valueType == ValueTypeList {
		return "", storage.ErrWrongType
	}

	if !ok || isExpired(item) {
		return "", storage.ErrKeyNotFound
	}

	// cast to string in different ways according to value type
	var result string
	if item.valueType == ValueTypeString {
		result = item.value.(string)
	}
	if item.valueType == ValueTypeInt {
		result = strconv.Itoa(item.value.(int))
	}

	return result, nil
}

// DeleteMany deletes all records with specified keys. Returns count of deleted records
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

// RPush adds new elements to the end of the list available via given key.
// It returns length of list after addition. If there is no such list, it is created.
// If we push elements to non-list item, it becomes of list data type.
func (s *Storage) RPush(key string, values []string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if el, ok := s.data[key]; !ok {
		s.data[key] = item{
			valueType: ValueTypeList,
			value:     &ds.LinkedList{},
		}
	} else {
		list := ds.NewLinkedList()
		switch el.valueType {
		case ValueTypeInt:
			list.PushBack(strconv.Itoa(el.value.(int)))
		case ValueTypeString:
			list.PushBack(el.value.(string))
		}

		newItem := item{
			value:     list,
			valueType: ValueTypeList,
			expiresAt: el.expiresAt,
		}

		s.data[key] = newItem
	}

	var length int
	for _, value := range values {
		length = s.data[key].value.(*ds.LinkedList).PushBack(value)
	}

	return length, nil
}

// LPush adds new elements to the beginning of the list available via given key.
// It returns length of list after addition. If there is not such list, it is created.
func (s *Storage) LPush(key string, values []string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if el, ok := s.data[key]; !ok {
		s.data[key] = item{
			valueType: ValueTypeList,
			value:     &ds.LinkedList{},
		}
	} else {
		list := ds.NewLinkedList()
		switch el.valueType {
		case ValueTypeInt:
			list.PushBack(strconv.Itoa(el.value.(int)))
		case ValueTypeString:
			list.PushBack(el.value.(string))
		}

		newItem := item{
			value:     list,
			valueType: ValueTypeList,
			expiresAt: el.expiresAt,
		}

		s.data[key] = newItem
	}

	var length int
	for _, value := range values {
		length = s.data[key].value.(*ds.LinkedList).PushForward(value)
	}

	return length, nil
}

// LRange returns node values in range of indexes [start, stop].
func (s *Storage) LRange(key string, start, stop int) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RLock()

	el, ok := s.data[key]
	if !ok {
		return nil, storage.ErrKeyNotFound
	}
	if ok && !isExpired(el) && el.valueType != ValueTypeList {
		return nil, storage.ErrWrongType
	}

	list := s.data[key].value.(*ds.LinkedList)
	values := list.LRange(start, stop)

	return values, nil
}

// LPop pops first n elements from list via given key.
func (s *Storage) LPop(key string, n int) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	el, ok := s.data[key]
	if !ok {
		return []string{}, storage.ErrKeyNotFound
	}
	if ok && !isExpired(el) && el.valueType != ValueTypeList {
		return []string{}, storage.ErrWrongType
	}

	list := s.data[key].value.(*ds.LinkedList)
	values := list.PopForwardNTimes(n)

	return values, nil
}

func (s *Storage) ListLen(key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	el, ok := s.data[key]
	if !ok {
		return 0, storage.ErrKeyNotFound
	}
	if ok && !isExpired(el) && el.valueType != ValueTypeList {
		return 0, storage.ErrWrongType
	}

	list := s.data[key].value.(*ds.LinkedList)
	return list.Len(), nil
}
