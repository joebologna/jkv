package main

import (
	"context"
	"fmt"

	"github.com/panduit-joeb/jkv/store/apk"
)

func runTests(ctx context.Context, db *apk.Client, opts apk.Options) []string {
	data := []string{}
	data = append(data, runTestA(ctx, db, opts)...)
	data = append(data, runTestB(ctx, db, opts)...)
	data = append(data, runTestC(ctx, db, opts)...)
	data = append(data, runTestD(ctx, db, opts)...)
	data = append(data, runTestE(ctx, db, opts)...)
	data = append(data, runTestG(ctx, db, opts)...)
	data = append(data, runTestH(ctx, db, opts)...)
	data = append(data, runTestK(ctx, db, opts)...)
	return data
}

func runTestA(ctx context.Context, db *apk.Client, opts apk.Options) []string {
	results := []string{"A: HSET/KEYS/HDEL"}
	results = append(results, func() []string {
		results := []string{
			fmt.Sprintf("  db.Open returns %#v, IsOpen: %t", db.Open(), db.IsOpen),
			fmt.Sprintf("  flushdb OK? %t", db.FlushDB(ctx).Val() == "OK"),
			fmt.Sprintf("  (re)db.Open returns %#v, IsOpen: %t", db.Open(), db.IsOpen),
			fmt.Sprintf("  hset hash key1 one == 1? %t", db.HSet(ctx, "hash", "key1", "one").Val() == 1),
			fmt.Sprintf("  hset hash key1 one == 0? %t", db.HSet(ctx, "hash", "key1", "one").Val() == 0),
			fmt.Sprintf("  set scalar one == OK? %t", db.Set(ctx, "scalar", "one", 0).Val() == "OK"),
			fmt.Sprintf("  get scalar == one? %t", db.Get(ctx, "scalar").Val() == "one"),
			fmt.Sprintf("  len(keys *) == 2? %t", len(db.Keys(ctx, "*").Val()) == 2),
		}
		return results
	}()...)
	return results
}

func runTestB(ctx context.Context, db *apk.Client, opts apk.Options) []string {
	results := []string{"B: HSET/KEYS"}
	results = append(results, func() []string {
		results := []string{
			fmt.Sprintf("  db.Open returns %#v, IsOpen: %t", db.Open(), db.IsOpen),
			fmt.Sprintf("  flushdb OK? %t", db.FlushDB(ctx).Val() == "OK"),
			fmt.Sprintf("  (re)db.Open returns %#v, IsOpen: %t", db.Open(), db.IsOpen),
			fmt.Sprintf("  hset hash key1 one == 1? %t", db.HSet(ctx, "hash", "key1", "one").Val() == 1),
			fmt.Sprintf("  hset hash key1 one == 0? %t", db.HSet(ctx, "hash", "key1", "one").Val() == 0),
			fmt.Sprintf("  set scalar one == OK? %t", db.Set(ctx, "scalar", "one", 0).Val() == "OK"),
			fmt.Sprintf("  get scalar == one? %t", db.Get(ctx, "scalar").Val() == "one"),
			fmt.Sprintf("  len(keys *) == 2? %t", len(db.Keys(ctx, "*").Val()) == 2),
			fmt.Sprintf("  hdel hash key1 == 1? %t", db.HDel(ctx, "hash", "key1").Val() == int64(1)),
			fmt.Sprintf("  len(keys *) == 1? %t", len(db.Keys(ctx, "*").Val()) == 1),
		}
		return results
	}()...)
	return results
}

func runTestC(ctx context.Context, db *apk.Client, opts apk.Options) []string {
	results := []string{"C: SET/GET"}
	results = append(results, func() []string {
		results := []string{
			fmt.Sprintf("  db.Open returns %#v, IsOpen: %t", db.Open(), db.IsOpen),
			fmt.Sprintf("  flushdb OK? %t", db.FlushDB(ctx).Val() == "OK"),
			fmt.Sprintf("  (re)db.Open returns %#v, IsOpen: %t", db.Open(), db.IsOpen),
			fmt.Sprintf("  set a b == 1? %t", db.Set(ctx, "a", "b", 0).Val() == "OK"),
			fmt.Sprintf("  get a == b? %t", db.Get(ctx, "a").Val() == "b"),
		}
		return results
	}()...)
	return results
}

func runTestD(ctx context.Context, db *apk.Client, opts apk.Options) []string {
	results := []string{"D: SET/DEL"}
	results = append(results, func() []string {
		results := []string{
			fmt.Sprintf("  db.Open returns %#v, IsOpen: %t", db.Open(), db.IsOpen),
			fmt.Sprintf("  flushdb OK? %t", db.FlushDB(ctx).Val() == "OK"),
			fmt.Sprintf("  (re)db.Open returns %#v, IsOpen: %t", db.Open(), db.IsOpen),
			fmt.Sprintf("  set a b == 1? %t", db.Set(ctx, "a", "b", 0).Val() == "OK"),
			fmt.Sprintf("  del a == 1? %t", db.Del(ctx, "a").Val() == 1),
			fmt.Sprintf("  get a == ''? %t", db.Get(ctx, "a").Val() == ""),
			fmt.Sprintf("  set key1 one == 1? %t", db.Set(ctx, "key1", "one", 0).Val() == "OK"),
			fmt.Sprintf("  set key2 two == 1? %t", db.Set(ctx, "key2", "two", 0).Val() == "OK"),
			fmt.Sprintf("  set key3 three == 1? %t", db.Set(ctx, "key3", "three", 0).Val() == "OK"),
			fmt.Sprintf("  del key1 key2 key3 == 3? %t", db.Del(ctx, "key1", "key2", "key3").Val() == 3),
		}
		return results
	}()...)
	return results
}

func runTestE(ctx context.Context, db *apk.Client, opts apk.Options) []string {
	results := []string{"E: EXISTS"}
	results = append(results, func() []string {
		results := []string{
			fmt.Sprintf("  db.Open returns %#v, IsOpen: %t", db.Open(), db.IsOpen),
			fmt.Sprintf("  flushdb OK? %t", db.FlushDB(ctx).Val() == "OK"),
			fmt.Sprintf("  (re)db.Open returns %#v, IsOpen: %t", db.Open(), db.IsOpen),
			fmt.Sprintf("  set a b == 1? %t", db.Set(ctx, "a", "b", 0).Val() == "OK"),
			fmt.Sprintf("  set b c == 1? %t", db.Set(ctx, "b", "c", 0).Val() == "OK"),
			fmt.Sprintf("  exists a == 1? %t", db.Exists(ctx, "a").Val() == 1),
			fmt.Sprintf("  exists b == 1? %t", db.Exists(ctx, "b").Val() == 1),
			fmt.Sprintf("  exists a b == 2? %t", db.Exists(ctx, "a", "b").Val() == 2),
		}
		return results
	}()...)
	return results
}

func runTestG(ctx context.Context, db *apk.Client, opts apk.Options) []string {
	results := []string{"G: HEXISTS"}
	results = append(results, func() []string {
		results := []string{
			fmt.Sprintf("  db.Open returns %#v, IsOpen: %t", db.Open(), db.IsOpen),
			fmt.Sprintf("  flushdb OK? %t", db.FlushDB(ctx).Val() == "OK"),
			fmt.Sprintf("  (re)db.Open returns %#v, IsOpen: %t", db.Open(), db.IsOpen),
			fmt.Sprintf("  hset hash key1 one key2 two key3 three == 3? %t", db.HSet(ctx, "hash", "key1", "one", "key2", "two", "key3", "three").Val() == 3),
			fmt.Sprintf("  hexists hash key1 == true? %t", db.HExists(ctx, "hash", "key1").Val()),
			fmt.Sprintf("  hexists hash key2 == true? %t", db.HExists(ctx, "hash", "key2").Val()),
			fmt.Sprintf("  hexists hash key3 == true? %t", db.HExists(ctx, "hash", "key3").Val()),
		}
		return results
	}()...)
	return results
}

func runTestH(ctx context.Context, db *apk.Client, opts apk.Options) []string {
	results := []string{"H: HKEYS"}
	results = append(results, func() []string {
		results := []string{
			fmt.Sprintf("  db.Open returns %#v, IsOpen: %t", db.Open(), db.IsOpen),
			fmt.Sprintf("  flushdb OK? %t", db.FlushDB(ctx).Val() == "OK"),
			fmt.Sprintf("  (re)db.Open returns %#v, IsOpen: %t", db.Open(), db.IsOpen),
			fmt.Sprintf("  hset hash key1 one key2 two key3 three == 3? %t", db.HSet(ctx, "hash", "key1", "one", "key2", "two", "key3", "three").Val() == 3),
			fmt.Sprintf("  len(hkeys hash) == 3? %t", len(db.HKeys(ctx, "hash").Val()) == 3),
		}
		return results
	}()...)
	return results
}

func runTestK(ctx context.Context, db *apk.Client, opts apk.Options) []string {
	results := []string{"K: PING"}
	results = append(results, func() []string {
		results := []string{
			fmt.Sprintf("  db.Ping == \"PONG\"? %t", db.Ping(ctx).Val() == "PONG"),
		}
		return results
	}()...)
	return results
}
