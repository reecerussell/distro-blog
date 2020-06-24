package caching

import (
	"fmt"

	"github.com/rainycape/memcache"

	"github.com/reecerussell/distro-blog/libraries/logging"
)

// Cache is a high-level interface which acts as a wrapper
// around a *memcache.Client.
type Cache interface {
	Set(key string, data []byte) error
	Get(key string) ([]byte, bool)
}

type client struct {
	mc *memcache.Client
}

// New returns a new instance of Cache.
func New(host string) (Cache, error) {
	mc, err := memcache.New(host)
	if err != nil {
		return nil, err
	}

	return &client{
		mc: mc,
	}, nil
}

// Set adds a given item to the backing cache with the key.
func (c *client) Set(key string, data []byte) error {
	if c == nil || c.mc == nil {
		return fmt.Errorf("cache not instantiated")
	}

	i := &memcache.Item{
		Key: key,
		Value: data,
	}

	return c.mc.Set(i)
}

// Get attempts to retrieve an item from the cache.
func (c *client) Get(key string) ([]byte, bool) {
	if c == nil || c.mc == nil {
		return nil, false
	}

	i, err := c.mc.Get(key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			logging.Debugf("item '%s' does not exist in cache\n", key)
			return nil, false
		}

		logging.Errorf("failed to fetch item '%s' from cache: %v\n", key, err)
		return nil, false
	}

	return i.Value, true
}