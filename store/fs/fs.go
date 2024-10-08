package fs

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/panduit-joeb/jkv"
)

type Options struct {
	Addr, Password string
	DB             int
}

type Client struct {
	DBDir  string
	IsOpen bool
}

var _ jkv.Client = (*Client)(nil)

var DEFAULT_DB = GetDBDir()

func GetDBDir() (dir string) {
	return os.TempDir() + "/jkv_db"
}

func (c *Client) ScalarDir() string { return c.DBDir + "/scalars/" }
func (c *Client) HashDir() string   { return c.DBDir + "/hashes/" }
func notOpen() error                { return errors.New("DB is not open") }

func (c *Client) GetDBDir() string {
	return c.DBDir
}

func NewClient(opts *Options) (db *Client) {
	return &Client{DBDir: opts.Addr, IsOpen: false}
}

// Open a database by creating the directories required if they don't exist and mark the database open
func (c *Client) Open() error {
	c.IsOpen = false
	for _, dir := range []string{c.ScalarDir(), c.HashDir()} {
		if err := os.MkdirAll(dir, 0775); err != nil {
			return err
		}
	}
	c.IsOpen = true
	return nil
}

// Close a database, basically just mark it closed
func (c *Client) Close() { c.IsOpen = false }

// FLUSHDB a database by removing the j.dbDir and everything underneath, ignore errors for now
func (j *Client) FlushDB(ctx context.Context) *jkv.StatusCmd {
	os.RemoveAll(j.DBDir)
	return jkv.NewStatusCmd("OK", nil)
}

// Return data in scalar key data, error is file is missing or inaccessible
func (c *Client) Get(ctx context.Context, key string) *jkv.StringCmd {
	if c.IsOpen {
		data, err := os.ReadFile(c.ScalarDir() + key)
		return jkv.NewStringCmd(string(data), err)
	}
	return jkv.NewStringCmd("", notOpen())
}

// Set a scalar key to a value
func (c *Client) Set(ctx context.Context, key, value string, expiration time.Duration) *jkv.StatusCmd {
	if c.IsOpen {
		return jkv.NewStatusCmd("OK", os.WriteFile(c.ScalarDir()+key, []byte(value), 0660))
	}
	return jkv.NewStatusCmd("(nil)", notOpen())
}

// Delete a key by removing the scalar file
func (c *Client) Del(ctx context.Context, keys ...string) *jkv.IntCmd {
	if c.IsOpen {
		n := 0
		for _, key := range keys {
			if os.Remove(c.ScalarDir()+key) == nil {
				n++
			}
		}
		return jkv.NewIntCmd(int64(n), nil)
	}
	return jkv.NewIntCmd(0, notOpen())
}

// KEYS returns the scalar and hash keys
func (c *Client) Keys(ctx context.Context, pattern string) *jkv.StringSliceCmd {
	var files []string
	for _, dir := range []string{c.HashDir(), c.ScalarDir()} {
		entries, err := os.ReadDir(dir)
		if err != nil {
			if os.IsNotExist(err) {
				return jkv.NewStringSliceCmd([]string{}, nil)
			}
			return jkv.NewStringSliceCmd([]string{}, err)
		}
		for _, file := range entries {
			files = append(files, file.Name())
		}
	}
	return jkv.NewStringSliceCmd(files, nil)
}

// Return true if scalar key file exists, false otherwise
func (c *Client) Exists(ctx context.Context, keys ...string) *jkv.IntCmd {
	if c.IsOpen {
		n := int64(0)
		for _, key := range keys {
			if _, err := os.Stat(c.ScalarDir() + key); err == nil {
				n++
			}
		}
		return jkv.NewIntCmd(n, nil)
	}
	return jkv.NewIntCmd(0, notOpen())
}

// Return data in hashed key data, error is file is missing or inaccessible
func (c *Client) HGet(ctx context.Context, hash, key string) *jkv.StringCmd {
	if c.IsOpen {
		data, err := os.ReadFile(c.HashDir() + hash + "/" + key)
		if err != nil {
			return jkv.NewStringCmd("", err)
		}
		return jkv.NewStringCmd(string(data), nil)
	}
	return jkv.NewStringCmd("", notOpen())
}

// Create a hash directory and store the data in a key file
// todo: reject a hash if a scalar key exists
func (c *Client) HSet(ctx context.Context, hash string, values ...string) *jkv.IntCmd {
	if c.IsOpen {
		rec := c.Exists(ctx, hash)
		if rec.Err() != nil {
			return jkv.NewIntCmd(0, rec.Err())
		}

		if rec.Val() > 0 {
			return jkv.NewIntCmd(0, fmt.Errorf("key \"%s\" exists as a scalar, cannot be a hash", hash))
		}

		if err := os.MkdirAll(c.HashDir()+hash, 0775); err != nil {
			return jkv.NewIntCmd(0, rec.Err())
		}

		n := 0
		for i := 0; i < len(values); i++ {
			key := values[i]
			f := c.HashDir() + hash + "/" + key
			info, err := os.Stat(f)
			if info == nil && os.IsNotExist(err) {
				n++
			}
			if err := os.WriteFile(f, []byte(values[i+1]), 0664); err != nil {
				fmt.Println("write file failed")
				return jkv.NewIntCmd(0, rec.Err())
			}
			i++
		}
		return jkv.NewIntCmd(int64(n), nil)
	}
	return jkv.NewIntCmd(0, notOpen())
}

// Delete a hashed key by removing the file, if no keys exist after the operation remove the hash directory
func (c *Client) HDel(ctx context.Context, hash string, keys ...string) *jkv.IntCmd {
	if c.IsOpen {
		rec := c.Exists(ctx, hash)
		if rec.Err() != nil {
			return jkv.NewIntCmd(0, rec.Err())
		}

		if rec.Val() > 0 {
			return jkv.NewIntCmd(0, fmt.Errorf("key \"%s\" exists as a scalar, cannot be a hash", hash))
		}

		n := int64(0)
		for _, key := range keys {
			f := c.HashDir() + hash + "/" + key
			info, err := os.Stat(f)
			if info == nil && os.IsNotExist(err) {
				continue
			}
			if err := os.Remove(f); err == nil {
				n++
			} else {
				return jkv.NewIntCmd(0, err)
			}
		}
		// remove the hash if no more keys exist
		if files, err := os.ReadDir(c.HashDir() + hash); err == nil {
			if len(files) == 0 {
				if err = os.Remove(c.HashDir() + hash); err != nil {
					fmt.Println("removing", c.HashDir()+hash, "failed, err", err.Error())
				}
			}
		}
		return jkv.NewIntCmd(n, nil)
	}
	return jkv.NewIntCmd(0, notOpen())
}

// HKEYS returns the hash keys
func (c *Client) HKeys(ctx context.Context, hash string) *jkv.StringSliceCmd {
	var err error
	if c.IsOpen {
		if _, err = os.Stat(c.HashDir() + hash); err == nil {
			entries, err := os.ReadDir(c.HashDir() + hash)
			if err != nil {
				return jkv.NewStringSliceCmd([]string{}, err)
			}
			var files []string
			for _, file := range entries {
				files = append(files, file.Name())
			}
			return jkv.NewStringSliceCmd(files, nil)
		}
		return jkv.NewStringSliceCmd([]string{}, err)
	}
	return jkv.NewStringSliceCmd([]string{}, notOpen())
}

// Return true if hashed key file exists, false otherwise
func (c *Client) HExists(ctx context.Context, hash, key string) *jkv.BoolCmd {
	if c.IsOpen {
		var err error
		if _, err = os.Stat(c.HashDir() + hash + "/" + key); err != nil {
			return jkv.NewBoolCmd(false, err)
		}
		return jkv.NewBoolCmd(true, nil)
	}
	return jkv.NewBoolCmd(false, notOpen())
}

func (c *Client) Ping(ctx context.Context) *jkv.StatusCmd {
	if c.IsOpen {
		return jkv.NewStatusCmd("PONG", nil)
	}
	return jkv.NewStatusCmd("", notOpen())
}
