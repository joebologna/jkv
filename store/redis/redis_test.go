package redis

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScalar(t *testing.T) {
	t.Run("Test Open() Default DB", func(t *testing.T) {
		var (
			c   *JKV_DB
			err error
		)

		a := assert.New(t)

		os.Remove(DEFAULT_DB)

		c = NewJKVClient()

		// this will create the directory
		err = c.Open()
		a.Nil(err)

		// this will use the directory and write/remove a file to verify it is writable
		err = c.Open()
		a.Nil(err)
	})

	t.Run("Test FLUSHDB()", func(t *testing.T) {
		var (
			c   *JKV_DB
			err error
		)

		a := assert.New(t)

		c = NewJKVClient()

		// this will remove the directory
		c.FLUSHDB()
		_, err = os.Stat(c.DBDir)
		a.True(os.IsNotExist(err))
	})

	t.Run("Set Key to Value", func(t *testing.T) {
		var (
			c = NewJKVClient()
		)

		os.Remove(DEFAULT_DB)

		a := assert.New(t)
		// this will create the directory
		err := c.Open()
		a.Nil(err)

		a.Nil(c.SET("this", "that"))
	})

	t.Run("Del Key", func(t *testing.T) {
		var (
			c = NewJKVClient()
		)

		os.Remove(DEFAULT_DB)

		a := assert.New(t)
		// this will create the directory
		err := c.Open()
		a.Nil(err)

		key := "this"
		value := "that"
		a.Nil(c.SET(key, value))
		a.Nil(c.DEL(key))
	})

	t.Run("Get Key", func(t *testing.T) {
		var (
			c = NewJKVClient()
		)

		os.Remove(DEFAULT_DB)

		a := assert.New(t)
		// this will create the directory
		err := c.Open()
		a.Nil(err)

		key := "this"
		value := "that"
		a.Nil(c.SET(key, value))
		data, err := c.GET(key)
		a.Nil(err)
		a.Equal(value, data)
	})

	t.Run("Key Exists?", func(t *testing.T) {
		var (
			c = NewJKVClient()
		)

		os.Remove(DEFAULT_DB)

		a := assert.New(t)
		// this will create the directory
		err := c.Open()
		a.Nil(err)

		key := "this"
		value := "that"
		a.Nil(c.SET(key, value))
		a.True(c.EXISTS(key))
	})

	t.Run("Keys", func(t *testing.T) {
		var (
			c = NewJKVClient()
		)

		os.Remove(DEFAULT_DB)

		a := assert.New(t)
		// this will create the directory
		err := c.Open()
		a.Nil(err)

		key := "this"
		value := "that"
		a.Nil(c.SET(key, value))
		a.True(c.EXISTS(key))
		keys, err := c.KEYS("*")
		a.Nil(err)
		a.Equal(1, len(keys))
		a.Equal(key, keys[0])
	})
}

func TestHash(t *testing.T) {
	t.Run("Set Hash", func(t *testing.T) {
		var (
			c = NewJKVClient()
		)

		os.Remove(DEFAULT_DB)

		a := assert.New(t)
		// this will create the directory
		err := c.Open()
		a.Nil(err)

		hash := "hashed"
		key := "this"
		value := "that"
		a.Nil(c.HSET(hash, key, value))
	})

	t.Run("Get Hash", func(t *testing.T) {
		var (
			c = NewJKVClient()
		)

		os.Remove(DEFAULT_DB)

		a := assert.New(t)
		// this will create the directory
		err := c.Open()
		a.Nil(err)

		hash := "hashed"
		key := "this"
		value := "that"
		a.Nil(c.HSET(hash, key, value))
		data, err := c.HGET(hash, key)
		a.Nil(err)
		a.Equal(value, data)
	})

	t.Run("Exists Hash", func(t *testing.T) {
		var (
			c = NewJKVClient()
		)

		os.Remove(DEFAULT_DB)

		a := assert.New(t)
		// this will create the directory
		err := c.Open()
		a.Nil(err)

		hash := "hashed"
		key := "this"
		value := "that"
		a.Nil(c.HSET(hash, key, value))
		a.True(c.HEXISTS(hash, key))
	})

	t.Run("Del Hash and it's dir", func(t *testing.T) {
		var (
			c = NewJKVClient()
		)

		os.Remove(DEFAULT_DB)

		a := assert.New(t)
		// this will create the directory
		err := c.Open()
		a.Nil(err)

		hash := "hashed"
		key := "this"
		value := "that"
		a.Nil(c.HSET(hash, key, value))
		a.Nil(c.HDEL(hash, key))
	})
}
