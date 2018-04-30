package cache

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/coreos/bbolt"
)

// Cache represents the cache.
type Cache struct {
	bucket []byte
	db     *bolt.DB
}

// New creates a new cache.
func New(dir string) (*Cache, error) {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	db, err := bolt.Open(filepath.Join(dir, "cache.db"), 0600, &bolt.Options{Timeout: 1 * time.Second})

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
	return c.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(c.bucket))
	})
}

// Get gets a value from the cache and returns it or a error.
func (c *Cache) Get(key string) ([]byte, error) {
	var v []byte

	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(c.bucket)

		if b == nil {
			return bolt.ErrBucketNotFound
		}

		v = b.Get([]byte(key))
		if len(v) == 0 {
			return errors.New("empty cache value")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return v, nil
}

// Set sets a value and return a error if any.
func (c *Cache) Set(key string, value []byte) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(c.bucket)
		if err != nil {
			return err
		}

		return b.Put([]byte(key), value)
	})
}
