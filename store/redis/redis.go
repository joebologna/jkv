package redis

import (
	"context"
	"errors"

	"github.com/panduit-joeb/jkv"

	"github.com/go-redis/redis/v8"
)

type JKV_DB struct {
	DBDir  string
	IsOpen bool
}

var _ jkv.JKV_OP = (*JKV_DB)(nil)

const DEFAULT_DB = "localhost:6379"

var client *redis.Client

func notOpen() error { return errors.New("DB is not open") }

// default location of db is "./jkv_db"
func NewJKVClient(db_dir ...string) (db *JKV_DB) {
	if len(db_dir) == 0 {
		db = &JKV_DB{DBDir: DEFAULT_DB}
	} else {
		db = &JKV_DB{DBDir: db_dir[0]}
	}
	client = redis.NewClient(&redis.Options{Addr: db.DBDir})
	return db
}

// Open a database by creating the directories required if they don't exist and mark the database open
func (j *JKV_DB) Open() error {
	j.IsOpen = true
	return nil
}

// Close a database, basically just mark it closed
func (j *JKV_DB) Close() { j.IsOpen = false; client.Close() }

// FLUSHDB a database by removing the j.dbDir and everything underneath, ignore errors for now
func (j *JKV_DB) FLUSHDB() { client.FlushDB(context.Background()) }

// Return data in scalar key data, error is file is missing or inaccessible
func (j *JKV_DB) GET(key string) (value string, err error) {
	if j.IsOpen {
		rec := client.Get(context.Background(), key)
		if rec.Err() != nil {
			return "", rec.Err()
		}
		return rec.Val(), nil
	}
	return "", notOpen()
}

// Set a scalar key to a value
func (j *JKV_DB) SET(key, value string) (err error) {
	if j.IsOpen {
		return client.Set(context.Background(), key, value, 0).Err()
	}
	return notOpen()
}

// Delete a key by removing the scalar file
func (j *JKV_DB) DEL(key string) error {
	if j.IsOpen {
		return client.Del(context.Background(), key).Err()
	}
	return notOpen()
}

// KEYS return a list of keys
func (j *JKV_DB) KEYS(pattern string) ([]string, error) {
	if j.IsOpen {
		rec := client.Keys(context.Background(), pattern)
		if rec.Err() != nil {
			return []string{}, rec.Err()
		}
		return rec.Val(), nil
	}
	return []string{}, notOpen()
}

// Return true if scalar key file exists, false otherwise
func (j *JKV_DB) EXISTS(key string) bool {
	if j.IsOpen {
		return client.Exists(context.Background(), key).Val() == 1
	}
	return false
}

// Return data in hashed key data, error is file is missing or inaccessible
func (j *JKV_DB) HGET(hash, key string) (value string, err error) {
	if j.IsOpen {
		rec := client.HGet(context.Background(), hash, key)
		if rec.Err() != nil {
			return "", rec.Err()
		}
		return rec.Val(), nil
	}
	return "", notOpen()
}

// Create a hash directory and store the data in a key file
func (j *JKV_DB) HSET(hash, key, value string) (err error) {
	if j.IsOpen {
		return client.HSet(context.Background(), hash, key, value).Err()
	}
	return notOpen()
}

// Delete a hashed key by removing the file, if no keys exist after the operation remove the hash directory
func (j *JKV_DB) HDEL(hash, key string) (err error) {
	if j.IsOpen {
		return client.HDel(context.Background(), hash, key).Err()
	}
	return notOpen()
}

// HKEYS return a list of keys for a hash
func (j *JKV_DB) HKEYS(hash, pattern string) ([]string, error) {
	if j.IsOpen {
		rec := client.HKeys(context.Background(), hash)
		if rec.Err() != nil {
			return []string{}, rec.Err()
		}
		return rec.Val(), nil
	}
	return []string{}, notOpen()
}

// Return true if hashed key file exists, false otherwise
func (j *JKV_DB) HEXISTS(hash, key string) bool {
	if j.IsOpen {
		return client.HExists(context.Background(), hash, key).Val()
	}
	return false
}
