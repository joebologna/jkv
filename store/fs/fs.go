package fs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/panduit-joeb/jkv"
)

type JKV_DB struct {
	DBDir  string
	IsOpen bool
}

var _ jkv.JKV_OP = (*JKV_DB)(nil)

const DEFAULT_DB = "jkv_db"

func (j *JKV_DB) ScalarDir() string { return j.DBDir + "/scalars/" }
func (j *JKV_DB) HashDir() string   { return j.DBDir + "/hashes/" }
func notOpen() error                { return errors.New("DB is not open") }

// default location of db is "./jkv_db"
func NewJKVClient(db_dir ...string) *JKV_DB {
	if len(db_dir) == 0 {
		return &JKV_DB{DBDir: DEFAULT_DB}
	}
	return &JKV_DB{DBDir: db_dir[0]}
}

// Open a database by creating the directories required if they don't exist and mark the database open
func (j *JKV_DB) Open() error {
	j.IsOpen = false
	for _, dir := range []string{j.ScalarDir(), j.HashDir()} {
		if err := os.MkdirAll(dir, 0775); err != nil {
			return err
		}
	}
	j.IsOpen = true
	return nil
}

// Close a database, basically just mark it closed
func (j *JKV_DB) Close() { j.IsOpen = false }

// FLUSHDB a database by removing the j.dbDir and everything underneath, ignore errors for now
func (j *JKV_DB) FLUSHDB() { os.RemoveAll(j.DBDir) }

// Return data in scalar key data, error is file is missing or inaccessible
func (j *JKV_DB) GET(key string) (value string, err error) {
	if j.IsOpen {
		var data []byte
		if data, err = os.ReadFile(j.ScalarDir() + key); err != nil {
			return "", err
		}
		return string(data), nil
	}
	return "", notOpen()
}

// Set a scalar key to a value
func (j *JKV_DB) SET(key, value string) (err error) {
	if j.IsOpen {
		return os.WriteFile(j.DBDir+"/scalars/"+key, []byte(value), 0660)
	}
	return notOpen()
}

// Delete a key by removing the scalar file
func (j *JKV_DB) DEL(key string) error {
	if j.IsOpen {
		return os.Remove(j.ScalarDir() + key)
	}
	return notOpen()
}

// KEYS returns the scalar and hash keys
func (j *JKV_DB) KEYS(pattern string) ([]string, error) {
	var files []string
	for _, dir := range []string{j.ScalarDir(), j.HashDir()} {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return []string{}, err
		}
		for _, file := range entries {
			files = append(files, file.Name())
		}
	}
	return files, nil
}

// Return true if scalar key file exists, false otherwise
func (j *JKV_DB) EXISTS(key string) bool {
	if j.IsOpen {
		_, err := os.Stat(j.ScalarDir() + key)
		return err == nil
	}
	return false
}

// Return data in hashed key data, error is file is missing or inaccessible
func (j *JKV_DB) HGET(hash, key string) (value string, err error) {
	if j.IsOpen {
		data, err := os.ReadFile(j.HashDir() + hash + "/" + key)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	return "", notOpen()
}

// Create a hash directory and store the data in a key file
// todo: reject a hash if a scalar key exists
func (j *JKV_DB) HSET(hash, key, value string) (err error) {
	if j.IsOpen {
		if !j.EXISTS(hash) {
			if err = os.MkdirAll(j.HashDir()+hash, 0775); err != nil {
				return err
			}
			return os.WriteFile(j.HashDir()+hash+"/"+key, []byte(value), 0664)
		}
		return fmt.Errorf("key \"%s\" exists as a scalar, cannot be a hash", hash)
	}
	return notOpen()
}

// Delete a hashed key by removing the file, if no keys exist after the operation remove the hash directory
func (j *JKV_DB) HDEL(hash, key string) (err error) {
	if j.IsOpen {
		if err = os.Remove(j.HashDir() + hash + "/" + key); err != nil {
			return err
		}
		var entries []fs.DirEntry
		if entries, err = os.ReadDir(j.HashDir() + hash); err != nil {
			return err
		}
		if len(entries) == 0 {
			return os.RemoveAll(j.HashDir() + hash)
		}
		return nil
	}
	return notOpen()
}

// HKEYS returns the hash keys
func (j *JKV_DB) HKEYS(hash string) (_ []string, err error) {
	if _, err = os.Stat(j.HashDir() + hash); err == nil {
		entries, err := os.ReadDir(j.HashDir() + hash)
		if err != nil {
			return []string{}, err
		}
		var files []string
		for _, file := range entries {
			files = append(files, file.Name())
		}
		return files, nil
	}
	return []string{}, err
}

// Return true if hashed key file exists, false otherwise
func (j *JKV_DB) HEXISTS(hash, key string) bool {
	if j.IsOpen {
		_, err := os.Stat(j.HashDir() + hash + "/" + key)
		return err == nil
	}
	return false
}
