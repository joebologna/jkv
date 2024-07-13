package redis

import (
	"context"
	"errors"

	"github.com/panduit-joeb/jkv"

	real_redis "github.com/go-redis/redis/v8"
)

type Options struct {
	Addr, Password string
	DB             int
}

type Client struct {
	DBDir       string
	IsOpen      bool
	RedisClient *real_redis.Client
}

var _ jkv.Client = (*Client)(nil)

const DEFAULT_DB = "localhost:6379"

func notOpen() error { return errors.New("DB is not open") }

func NewClient(opts *Options) (db *Client) {
	return &Client{DBDir: opts.Addr, IsOpen: false, RedisClient: real_redis.NewClient(&real_redis.Options{Addr: opts.Addr, Password: opts.Password, DB: opts.DB})}
}

// Open a database by creating the directories required if they don't exist and mark the database open
func (c *Client) Open() error {
	c.IsOpen = true
	return nil
}

// Close a database, basically just mark it closed
func (c *Client) Close() { c.IsOpen = false; c.RedisClient.Close() }

// FLUSHDB a database by removing the j.dbDir and everything underneath, ignore errors for now
func (c *Client) FlushDB(ctx context.Context) *jkv.StatusCmd {
	rec := c.RedisClient.FlushDB(context.Background())
	return jkv.NewStatusCmd(rec.Val(), rec.Err())
}

// Return data in scalar key data, error is file is missing or inaccessible
func (c *Client) Get(ctx context.Context, key string) *jkv.StringCmd {
	if c.IsOpen {
		rec := c.RedisClient.Get(context.Background(), key)
		return jkv.NewStringCmd(rec.Val(), rec.Err())
	}
	return jkv.NewStringCmd("", notOpen())
}

// Set a scalar key to a value
func (c *Client) Set(ctx context.Context, key, value string) *jkv.StatusCmd {
	if c.IsOpen {
		rec := c.RedisClient.Set(ctx, key, value, 0)
		return jkv.NewStatusCmd(rec.Val(), rec.Err())
	}
	return jkv.NewStatusCmd("", notOpen())
}

// Delete a key by removing the scalar file
func (c *Client) Del(ctx context.Context, keys ...string) *jkv.IntCmd {
	if c.IsOpen {
		rec := c.RedisClient.Del(context.Background(), keys...)
		return jkv.NewIntCmd(rec.Val(), rec.Err())
	}
	return jkv.NewIntCmd(0, notOpen())
}

// KEYS return a list of keys
func (c *Client) Keys(ctx context.Context, pattern string) *jkv.StringSliceCmd {
	if c.IsOpen {
		rec := c.RedisClient.Keys(context.Background(), pattern)
		return jkv.NewStringSliceCmd(rec.Val(), rec.Err())
	}
	return jkv.NewStringSliceCmd([]string{}, notOpen())
}

// Return true if scalar key file exists, false otherwise
func (c *Client) Exists(ctx context.Context, keys ...string) *jkv.IntCmd {
	if c.IsOpen {
		rec := c.RedisClient.Exists(context.Background(), keys...)
		return jkv.NewIntCmd(rec.Val(), rec.Err())
	}
	return jkv.NewIntCmd(0, notOpen())
}

// Return data in hashed key data, error is file is missing or inaccessible
func (c *Client) HGet(ctx context.Context, hash, key string) *jkv.StringCmd {
	if c.IsOpen {
		rec := c.RedisClient.HGet(ctx, hash, key)
		return jkv.NewStringCmd(rec.Val(), rec.Err())
	}
	return jkv.NewStringCmd("", notOpen())
}

// Create a hash directory and store the data in a key file
func (c *Client) HSet(ctx context.Context, hash, key string, values ...string) *jkv.IntCmd {
	var rec *real_redis.IntCmd
	if c.IsOpen {
		rec = c.RedisClient.HSet(ctx, hash, key, values)
		return jkv.NewIntCmd(rec.Val(), rec.Err())
	}
	return jkv.NewIntCmd(0, notOpen())
}

// Delete a hashed key by removing the file, if no keys exist after the operation remove the hash directory
func (c *Client) HDel(ctx context.Context, hash, key string) *jkv.IntCmd {
	var rec *real_redis.IntCmd
	if c.IsOpen {
		rec = c.RedisClient.HDel(ctx, hash, key)
		return jkv.NewIntCmd(rec.Val(), rec.Err())
	}
	return jkv.NewIntCmd(0, notOpen())
}

// HKEYS return a list of keys for a hash
func (c *Client) HKeys(ctx context.Context, hash string) *jkv.StringSliceCmd {
	if c.IsOpen {
		rec := c.RedisClient.HKeys(ctx, hash)
		return jkv.NewStringSliceCmd(rec.Val(), rec.Err())
	}
	return jkv.NewStringSliceCmd([]string{}, notOpen())
}

// Return true if hashed key file exists, false otherwise
func (c *Client) HExists(ctx context.Context, hash, key string) *jkv.BoolCmd {
	if c.IsOpen {
		rec := c.RedisClient.HExists(context.Background(), hash, key)
		return jkv.NewBoolCmd(rec.Val(), rec.Err())
	}
	return jkv.NewBoolCmd(false, notOpen())
}

func (c *Client) Ping(ctx context.Context) *jkv.StatusCmd {
	rec := c.RedisClient.Ping(ctx)
	return jkv.NewStatusCmd(rec.Val(), rec.Err())
}
