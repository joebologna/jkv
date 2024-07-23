package apk

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"fyne.io/fyne/v2/storage"
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

func (j *Client) ScalarDir() string { return j.DBDir + "/scalars" }
func (j *Client) HashDir() string   { return j.DBDir + "/hashes" }
func notOpen() error                { return errors.New("DB is not open") }

func NewClient(opts *Options) (db *Client) {
	return &Client{DBDir: opts.Addr, IsOpen: false}
}

func (j *Client) Open() (err error) {
	var uri = storage.NewFileURI(j.DBDir)
	storage.CreateListable(uri)
	uri = storage.NewFileURI(j.ScalarDir())
	storage.CreateListable(uri)
	uri = storage.NewFileURI(j.HashDir())
	storage.CreateListable(uri)
	j.IsOpen = true
	return nil
}

// Open a database by creating the directories required if they don't exist and mark the database open
// func (j *Client) Open() (err error) {
// 	uri := storage.NewFileURI(j.DBDir)
// 	ok, err := storage.CanList(uri)
// 	fmt.Printf("CanList(%s): ok = %t, err = %#v\n", uri, ok, err)
// 	ok, err = storage.Exists(uri)
// 	fmt.Printf("Exists(%s): ok = %t, err = %#v\n", uri, ok, err)
// 	err = storage.CreateListable(uri)
// 	fmt.Printf("CreateListable(%s): err = %#v\n", uri, err)
// 	ok, err = storage.CanList(uri)
// 	fmt.Printf("CanList(%s): ok = %t, err = %#v\n", uri, ok, err)
// 	if !j.IsOpen {
// 		fmt.Println("db is not open yet.")
// 		fmt.Println("making sub dirs")
// 		if err = j.makeSubDirs(); err != nil {
// 			return err
// 		}
// 		fmt.Println("making sub dirs worked, db is open")
// 		j.IsOpen = true
// 	} else {
// 		fmt.Println("db is already opened")
// 	}
// 	return nil
// }

// func (j *Client) makeSubDirs() (err error) {
// 	ok := true
// 	for _, dir := range []string{j.HashDir(), j.ScalarDir()} {
// 		uri := storage.NewFileURI(dir)
// 		ok, err = storage.Exists(uri)
// 		fmt.Printf("Exists(%s): ok = %t, err = %#v\n", uri, ok, err)
// 		err = storage.CreateListable(uri)
// 		fmt.Printf("CreateListable(%s): err = %#v\n", uri, err)
// 		ok, err = storage.CanList(uri)
// 		fmt.Printf("CanList(%s): ok = %t, err = %#v\n", uri, ok, err)
// 	}
// 	return err
// }

// func (j *Client) makeSubDirs() (err error) {
// 	ok := true
// 	for _, dir := range []string{j.HashDir(), j.ScalarDir()} {
// 		_, err = Stat(dir)
// 		if err != nil {
// 			if os.IsNotExist(err) {
// 				fmt.Println(dir, "does not exist, try to make it")
// 				if err = Mkdir(j.DBDir, 0775); err == nil {
// 					fmt.Println(j.DBDir, "created")
// 					ok = ok && true
// 				} else {
// 					fmt.Println("mkdir", j.DBDir, "failed", err.Error())
// 					ok = ok && false
// 				}
// 			} else if err == nil || os.IsExist(err) {
// 				fmt.Println(dir, "exists, skip making it")
// 				ok = ok && true
// 			} else {
// 				fmt.Println("stat", j.DBDir, "failed with unknown err:", err.Error())
// 				ok = ok && false
// 			}
// 		} else {
// 			fmt.Println("stat", j.DBDir, "returned nil, assuming this is good")
// 			ok = ok && true
// 		}
// 	}
// 	if ok {
// 		return nil
// 	}
// 	return errors.New("something bad happened, err: " + err.Error())
// }

// Close a database, basically just mark it closed
func (j *Client) Close() { j.IsOpen = false }

// FLUSHDB a database by removing the j.dbDir and everything underneath, ignore errors for now
func (j *Client) FlushDB(ctx context.Context) *jkv.StatusCmd {
	RemoveAll(j.DBDir)
	// need to recreate the directory structure
	j.Open()
	return jkv.NewStatusCmd("OK", nil)
}

// Return data in scalar key data, error is file is missing or inaccessible
func (c *Client) Get(ctx context.Context, key string) *jkv.StringCmd {
	if c.IsOpen {
		data, err := ReadFile(c.ScalarDir() + "/" + key)
		return jkv.NewStringCmd(string(data), err)
	}
	return jkv.NewStringCmd("", notOpen())
}

// Set a scalar key to a value
func (c *Client) Set(ctx context.Context, key, value string, expiration time.Duration) *jkv.StatusCmd {
	if c.IsOpen {
		return jkv.NewStatusCmd("OK", WriteFile(c.DBDir+"/scalars/"+key, []byte(value), 0660))
	}
	return jkv.NewStatusCmd("(nil)", notOpen())
}

// Delete a key by removing the scalar file
func (c *Client) Del(ctx context.Context, keys ...string) *jkv.IntCmd {
	if c.IsOpen {
		n := 0
		for _, key := range keys {
			if Remove(c.ScalarDir()+"/"+key) == nil {
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
		entries, err := ReadDir(dir)
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
		// todo: add a loop here
		if _, err := Stat(c.ScalarDir() + "/" + keys[0]); err != nil {
			if os.IsNotExist(err) {
				return jkv.NewIntCmd(int64(0), nil)
			}
			return jkv.NewIntCmd(0, err)
		}
		// return jkv.NewIntCmd(int64(len(keys)), nil)
		return jkv.NewIntCmd(1, nil)
	}
	return jkv.NewIntCmd(0, nil)
}

// Return data in hashed key data, error is file is missing or inaccessible
func (c *Client) HGet(ctx context.Context, hash, key string) *jkv.StringCmd {
	if c.IsOpen {
		f := c.HashDir() + "/" + hash + "/" + key
		data, err := ReadFile(f)
		if err != nil {
			fmt.Println("reading", f, "failed, err:", err.Error())
			return jkv.NewStringCmd("", err)
		}
		fmt.Println("reading", f, "succeeded, data:", string(data))
		return jkv.NewStringCmd(string(data), nil)
	}
	return jkv.NewStringCmd("", notOpen())
}

// Create a hash directory and store the data in a key file
// todo: reject a hash if a scalar key exists
func (c *Client) HSet(ctx context.Context, hash string, values ...string) *jkv.IntCmd {
	if c.IsOpen {
		n := int64(0)
		rec := c.Exists(ctx, hash)
		if rec.Err() == nil && rec.Val() == 1 {
			return jkv.NewIntCmd(0, errors.New("scalar exists, this is bad"))
		}
		// todo: add loop here
		f := c.HashDir() + "/" + hash + "/" + values[0]
		err := WriteFile(f, []byte(values[1]), 0644)
		if err == nil {
			n++
			// fmt.Printf("WriteFile(%s, %s), succeeded.\n", f, values[1])
		} else {
			fmt.Printf("WriteFile(%s, %s), failed. err: %#v\n", f, values[1], err)
		}
		return jkv.NewIntCmd(n, nil)
	}
	return jkv.NewIntCmd(0, errors.New("db is not open"))
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
			f := c.HashDir() + "/" + hash + "/" + key
			info, err := Stat(f)
			if info == nil && os.IsNotExist(err) {
				continue
			}
			if err := Remove(f); err == nil {
				n++
			} else {
				return jkv.NewIntCmd(0, err)
			}
		}
		// remove the hash if no more keys exist
		if files, err := ReadDir(c.HashDir() + "/" + hash); err == nil {
			if len(files) == 0 {
				if err = Remove(c.HashDir() + "/" + hash); err != nil {
					fmt.Println("removing", c.HashDir()+"/"+hash, "failed, err", err.Error())
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
		if _, err = Stat(c.HashDir() + "/" + hash); err == nil {
			entries, err := ReadDir(c.HashDir() + "/" + hash)
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
		if _, err = Stat(c.HashDir() + "/" + hash + "/" + key); err != nil {
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
