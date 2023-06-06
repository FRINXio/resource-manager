package cache

import (
	"encoding/json"
	"time"

	"github.com/allegro/bigcache/v3"
)

type BigCacheService[T any] interface {
	Set(key string, value *T) error
	Get(key string) (*T, error)
	Delete(key string) error
	Update(key string, value *T) error
}

type BigCacheServiceImpl[T any] struct {
	cache *bigcache.BigCache
}

func NewBigCacheService[T any](invalidateAfter time.Duration) *BigCacheServiceImpl[T] {
	cache, _ := bigcache.NewBigCache(bigcache.DefaultConfig(invalidateAfter))
	return &BigCacheServiceImpl[T]{cache: cache}
}

func (b *BigCacheServiceImpl[T]) Set(key string, value T) error {
	valueBuffer, err := json.Marshal(value)

	if err != nil {
		return err
	}

	return b.cache.Set(key, valueBuffer)
}

func (b *BigCacheServiceImpl[T]) Get(key string) (*T, error) {
	valueBuffer, err := b.cache.Get(key)

	if err != nil {
		return nil, err
	}

	var value T
	err = json.Unmarshal(valueBuffer, &value)

	if err != nil {
		return nil, err
	}

	return &value, nil
}

func (b *BigCacheServiceImpl[T]) Delete(key string) error {
	return b.cache.Delete(key)
}

func (b *BigCacheServiceImpl[T]) Update(key string, value T) error {
	return b.Set(key, value)
}
