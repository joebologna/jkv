package dbutil

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	redis "github.com/panduit-joeb/jkv/store/redis"
)

var LOCAL_REDIS = &redis.Options{Addr: "127.0.0.1:6379", Password: "", DB: 0}

func NewRDBClient(opts *redis.Options) *redis.Client {
	return redis.NewClient(opts)
}

// Utility function for accessing data in an associative array
func GetDataUsingField(RDB *redis.Client, collection string, key string) (string, error) {
	rec := RDB.HGet(context.Background(), collection, key)
	if rec.Err() != nil {
		return "", rec.Err()
	}
	return rec.Val(), nil
}

// Utility function for storing data in an associative array
func UpsertItem(RDB *redis.Client, collection string, key string, value string) error {
	rec := RDB.HSet(context.Background(), collection, key, value)
	if rec.Err() != nil {
		return rec.Err()
	}
	return nil
}

// WaitForRedis Utility function to spin waiting for Redis to (re)start
func WaitForRedis(RDB *redis.Client) {
	deadman := 60
	freq := 1
	for {
		rc := RDB.Ping(context.Background())
		err := rc.Err()
		if err == nil {
			// fmt.Println("Redis is up.")
			break
		}
		fmt.Printf("Redis is down, err = %s. Waiting (%d).\n", err, deadman*freq)
		time.Sleep(time.Second * time.Duration(freq))
		deadman--
		if deadman < 0 {
			fmt.Printf("Exited after waiting %ds\n", 20*freq)
			os.Exit(1)
		}
	}
}

func SplitPath(name string) (string, string, string) {
	path := strings.Split(name, "/")
	if len(path) != 3 {
		log.Fatal("bad zip")
	}
	return path[0], path[1], path[2]
}

func bye(err error) error {
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

/*
These utility functions are used to simplify and standardize using Redis
to store collections as associative arrays and providing the current or default
value of the collection.

Collections are stored as Hashes, where the value of each field contains
the data for the collection - basically an associative array.

The current item in a collection has the field ":current", therefore GetCurrentData will return Collection[":current"]

It's possible a current item is unset, if the flag defaultOK == true, the Collection[":default"] will be returned.
*/

func HasCurrentItem(RDB *redis.Client, collection string) bool {
	return FieldExists(RDB, collection, ":current")
}

func HasDefaultItem(RDB *redis.Client, collection string) bool {
	return FieldExists(RDB, collection, ":default")
}

// return the data for collection[:current] or :default item
func GetCurrentItem(RDB *redis.Client, collection string, defaultOK bool) (string, error) {
	if HasCurrentItem(RDB, collection) {
		return GetCurrentField(RDB, collection)
	}
	if defaultOK && HasDefaultItem(RDB, collection) {
		item, err := GetDataUsingField(RDB, collection, ":default")
		if err != nil {
			return "", fmt.Errorf(":current and :default are missing from %s", collection)
		}
		return item, nil
	}
	return "", fmt.Errorf(":current is missing from %s", collection)
}

// return the data for collection[:default] item
func GetDefaultItem(RDB *redis.Client, collection string) (string, error) {
	if HasDefaultItem(RDB, collection) {
		item, err := GetDataUsingField(RDB, collection, ":default")
		if err != nil {
			return "", fmt.Errorf(":default is missing from %s", collection)
		}
		return item, nil
	}
	return "", fmt.Errorf("%s has no :default", collection)
}

// return the data at collection[collection[":current"]]
func GetCurrentField(RDB *redis.Client, collection string) (string, error) {
	data, err := GetDataUsingField(RDB, collection, ":current")
	if err != nil {
		return "", err
	}
	data, err = GetDataUsingField(RDB, collection, data)
	if err != nil {
		return "", err
	}
	return data, nil
}

func FieldExists(RDB *redis.Client, collection string, field string) bool {
	rec := RDB.HExists(context.Background(), collection, field)
	if rec.Err() != nil {
		return false
	}
	return rec.Val()
}

// set the data at collection[field], if field is :current, moves the pointer
func SetCurrentItem(RDB *redis.Client, collection string, field string) error {
	if !FieldExists(RDB, collection, field) {
		return fmt.Errorf("field %s does not exist in %s cannot set it current", field, collection)
	}
	rec := RDB.HSet(context.Background(), collection, ":current", field)
	if rec.Err() != nil {
		return rec.Err()
	}
	return nil
}

func SetDefaultItem(RDB *redis.Client, collection string, field string) error {
	rec := RDB.HSet(context.Background(), collection, ":default", field)
	if rec.Err() != nil {
		return rec.Err()
	}
	return nil
}

// remove the pointer
func ClearCurrentItem(RDB *redis.Client, collection string) error {
	rec := RDB.HDel(context.Background(), collection, ":current")
	if rec.Err() != nil {
		return rec.Err()
	}
	return nil
}

// delete an item at collection[field]
func DeleteItem(RDB *redis.Client, collection string, field string) error {
	rec := RDB.HDel(context.Background(), collection, field)
	if rec.Err() != nil {
		return rec.Err()
	}
	val, err := GetDataUsingField(RDB, collection, ":current")
	if err == nil && val == field {
		rec := RDB.HDel(context.Background(), collection, ":current")
		if rec.Err() != nil {
			return rec.Err()
		}
	}
	return nil
}
