package redis

import (
	"context"
	"os"
	"testing"

	"github.com/panduit-joeb/jkv"
	"github.com/stretchr/testify/assert"

	real_redis "github.com/go-redis/redis/v8"
)

func TestAgainstRealRedis(t *testing.T) {
	ctx := context.Background()

	t.Run("HSET Twice (Real)", func(t *testing.T) {
		a := assert.New(t)
		var (
			rc     *real_redis.Client
			int_rc *real_redis.IntCmd
		)
		rc = real_redis.NewClient(&real_redis.Options{Addr: "localhost:6379"})

		a.Nil(rc.FlushDB(ctx).Err())
		int_rc = rc.HSet(ctx, "hash1", "key1", "one", "key2", "two")
		a.Nil(int_rc.Err())
		a.Equal(int64(2), int_rc.Val())
		int_rc = rc.HSet(ctx, "hash1", "key1", "one", "key2", "two")
		a.Nil(int_rc.Err())
		a.Equal(int64(0), int_rc.Val())

		t.Cleanup(func() { rc.Close() })
	})

	t.Run("HSET Twice (JKV -r)", func(t *testing.T) {
		a := assert.New(t)
		var (
			rc     *Client
			int_rc *jkv.IntCmd
		)
		rc = NewClient(&Options{Addr: "localhost:6379"})
		rc.Open()

		a.Nil(rc.FlushDB(ctx).Err())
		int_rc = rc.HSet(ctx, "hash1", "key1", "one", "key2", "two")
		a.Nil(int_rc.Err())
		a.Equal(int64(2), int_rc.Val())
		int_rc = rc.HSet(ctx, "hash1", "key1", "one", "key2", "two")
		a.Nil(int_rc.Err())
		a.Equal(int64(0), int_rc.Val())

		t.Cleanup(func() { rc.Close() })
	})
}

func TestScalar(t *testing.T) {
	t.Run("Test Open()", func(t *testing.T) {
		var c = NewClient(&Options{Addr: "localhost:6379", Password: "", DB: 0})
		var err error

		defer c.Close()

		a := assert.New(t)

		err = c.Open()
		a.Nil(err)
	})

	t.Run("Test FlushDB()", func(t *testing.T) {
		var c = NewClient(&Options{Addr: "localhost:6379", Password: "", DB: 0})

		defer c.Close()

		ctx := context.Background()

		a := assert.New(t)

		c.Open()
		c.FlushDB(ctx)
		rec := c.Keys(context.Background(), "*")
		a.Nil(rec.Err())
		a.Equal(0, len(rec.Val()))
	})

	t.Run("Set Key to Value", func(t *testing.T) {
		var c = NewClient(&Options{Addr: "localhost:6379", Password: "", DB: 0})
		defer c.Close()

		ctx := context.Background()

		os.Remove(DEFAULT_DB)

		a := assert.New(t)
		// this will create the directory
		err := c.Open()
		a.Nil(err)

		rec := c.Set(ctx, "this", "that", 0)
		a.Nil(rec.Err())
		a.Equal("OK", rec.Val())
	})

	t.Run("Del Key", func(t *testing.T) {
		var c = NewClient(&Options{Addr: "localhost:6379", Password: "", DB: 0})
		defer c.Close()

		ctx := context.Background()

		c.FlushDB(ctx)

		os.Remove(DEFAULT_DB)

		a := assert.New(t)
		// this will create the directory
		err := c.Open()
		a.Nil(err)

		key := "this"
		value := "that"
		rec := c.Set(ctx, key, value, 0)
		a.Nil(rec.Err())
		a.Equal("OK", rec.Val())

		rec2 := c.Del(ctx, key)
		a.Nil(rec2.Err())
		a.Equal(int64(1), rec2.Val())
	})

	t.Run("Get Key", func(t *testing.T) {
		var c = NewClient(&Options{Addr: "localhost:6379", Password: "", DB: 0})
		defer c.Close()

		ctx := context.Background()

		os.Remove(DEFAULT_DB)

		a := assert.New(t)
		// this will create the directory
		err := c.Open()
		a.Nil(err)

		key := "this"
		value := "that"
		rec := c.Set(ctx, key, value, 0)
		a.Nil(rec.Err())
		a.Equal("OK", rec.Val())

		rec2 := c.Get(ctx, key)
		a.Nil(rec2.Err())
		a.Equal(value, rec2.Val())
	})

	t.Run("Key Exists?", func(t *testing.T) {
		var c = NewClient(&Options{Addr: "localhost:6379", Password: "", DB: 0})
		defer c.Close()

		ctx := context.Background()

		os.Remove(DEFAULT_DB)

		a := assert.New(t)
		// this will create the directory
		err := c.Open()
		a.Nil(err)

		key := "this"
		value := "that"
		rec := c.Set(ctx, key, value, 0)
		a.Nil(rec.Err())
		a.Equal("OK", rec.Val())

		rec2 := c.Exists(ctx, key)
		a.Nil(rec2.Err())
		a.Equal(int64(1), rec2.Val())
	})

	t.Run("Keys", func(t *testing.T) {
		var c = NewClient(&Options{Addr: "localhost:6379", Password: "", DB: 0})
		defer c.Close()

		ctx := context.Background()

		os.Remove(DEFAULT_DB)

		a := assert.New(t)
		// this will create the directory
		err := c.Open()
		a.Nil(err)

		key := "this"
		value := "that"
		rec := c.Set(ctx, key, value, 0)
		a.Nil(rec.Err())
		a.Equal("OK", rec.Val())

		rec2 := c.Keys(ctx, "*")
		a.Nil(rec2.Err())
		a.Equal(1, len(rec2.Val()))
	})
}

func TestHash(t *testing.T) {
	t.Run("Set Hash", func(t *testing.T) {
		var c = NewClient(&Options{Addr: "localhost:6379", Password: "", DB: 0})
		defer c.Close()

		os.Remove(DEFAULT_DB)

		a := assert.New(t)
		// this will create the directory
		err := c.Open()
		a.Nil(err)

		hash := "hashed"
		key := "this"
		value := "that"
		a.Nil(c.HSet(context.Background(), hash, key, value).Err())
	})

	t.Run("Get Hash", func(t *testing.T) {
		var c = NewClient(&Options{Addr: "localhost:6379", Password: "", DB: 0})
		defer c.Close()

		ctx := context.Background()
		os.Remove(DEFAULT_DB)

		a := assert.New(t)
		// this will create the directory
		err := c.Open()
		a.Nil(err)

		hash := "hashed"
		key := "this"
		value := "that"
		a.Nil(c.HSet(ctx, hash, key, value).Err())
		rec := c.HGet(ctx, hash, key)
		a.Nil(rec.Err())
		a.Equal(value, rec.Val())
	})

	t.Run("Exists Hash", func(t *testing.T) {
		var c = NewClient(&Options{Addr: "localhost:6379", Password: "", DB: 0})
		defer c.Close()

		ctx := context.Background()

		os.Remove(DEFAULT_DB)

		a := assert.New(t)
		// this will create the directory
		err := c.Open()
		a.Nil(err)

		hash := "hashed"
		key := "this"
		value := "that"
		a.Nil(c.HSet(ctx, hash, key, value).Err())
		rec := c.HExists(ctx, hash, key)
		a.Nil(rec.Err())
		a.True(rec.Val())
	})

	t.Run("Del Hash and it's dir", func(t *testing.T) {
		var c = NewClient(&Options{Addr: "localhost:6379", Password: "", DB: 0})
		defer c.Close()
		ctx := context.Background()

		os.Remove(DEFAULT_DB)

		a := assert.New(t)
		// this will create the directory
		err := c.Open()
		a.Nil(err)

		hash := "hashed"
		key := "this"
		value := "that"
		a.Nil(c.HSet(ctx, hash, key, value).Err())
		a.Nil(c.HDel(ctx, hash, key).Err())
	})
}
