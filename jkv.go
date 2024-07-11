package jkv

type JKV_OP interface {
	Open() error
	Close()
	FLUSHDB()
	GET(key string) (string, error)
	SET(key, value string) error
	DEL(key string) error
	KEYS(pattern string) ([]string, error)
	EXISTS(key string) bool
	HGET(hash, key string) (string, error)
	HSET(hash, key, value string) error
	HDEL(hash, key string) error
	HKEYS(hash string) ([]string, error)
	HEXISTS(hash, key string) bool
}
