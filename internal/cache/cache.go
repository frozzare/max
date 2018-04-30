package cache

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/coreos/bbolt"
)

const (
	filePerm os.FileMode = 0644
)

// Cache represents the cache.
type Cache struct {
	bucket []byte
	db     *bolt.DB

	sync.Mutex
}

// New creates a new cache.
func New(dir string) (*Cache, error) {
	if err := os.MkdirAll(dir, 0750); err != nil {
		return nil, err
	}

	db, err := bolt.Open(filepath.Join(dir, "cache.db"), filePerm, &bolt.Options{Timeout: 1 * time.Second})

	if err != nil {
		return nil, err
	}

	return &Cache{
		bucket: []byte("cache"),
		db:     db,
	}, nil
}

// Close the cache.
func (c *Cache) Close() error {
	return c.db.Close()
}

// Delete delets a cache value.
func (c *Cache) Delete(key string) error {
	c.Lock()
	defer c.Unlock()

	return c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(c.bucket)

		if b == nil {
			return bolt.ErrBucketNotFound
		}

		return b.Delete([]byte(key))
	})
}

// Flush flushes the cache.
func (c *Cache) Flush() error {
	c.Lock()
	defer c.Unlock()

	return c.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(c.bucket))
	})
}

// Get gets a value from the cache and returns it or a error.
func (c *Cache) Get(key string) ([]byte, error) {
	var val []byte

	c.Lock()
	defer c.Unlock()

	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(c.bucket)

		if b == nil {
			return bolt.ErrBucketNotFound
		}

		v := b.Get([]byte(key))
		if len(v) == 0 {
			return errors.New("empty cache value")
		}

		val = make([]byte, len(v))
		copy(val, v)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return val, nil
}

// Set sets a value and return a error if any.
func (c *Cache) Set(key string, value []byte) error {
	c.Lock()
	defer c.Unlock()

	return c.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(c.bucket)
		if err != nil {
			return err
		}

		return b.Put([]byte(key), value)
	})
}
