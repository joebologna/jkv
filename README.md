# Joe's Key/Value DB

I need an embedded key/value database that has the same behavior as Redis. Redis is implemented as a service, which runs in a separate process. This database runs in the same process and it therefore very portable. The database is basically a veneer over the file system where a key is a file and the value is the content of the file. Scalars are just a key with string data stored as the value. String is the only data type supported.

Hashes are supported as well. The full list of operations supported are listed in the [jkv.go](./jkv.go) interface.

JKV_OPs are simple versions of Redis operations exposed in the Redis Go API. Redis overloads responses with error and values with different data types. This simple approach uses the traditional (value, err) return from API calls instead.

# Data Stores

## jkv/store/fs

The jkv/store/fs package implements storage using files and directories. The implementation does not protect against go routines causing data corruption. This method is inherently persisent vs. the memcache approach taken by Redis.

## jkv/store/redis

The jkv/store/redis package implements storage using Redis. The implementation should be suitable for use with go routines because Redis is inherently designed to prevent data corruption during concurrent use of the database. It is not inherently persistent.

# Disclaimer

This project is a work in progress, it is far from complete. Using the jkv/store/fs implementation assumes a single process and a single thread of execution is used to prevent data corruption. Enhancements will likely be made when using this package in real life situtations. Or it may be abandoned as an experiment.
