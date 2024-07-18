#!/bin/bash
(./jkv-cli -f flushdb
./jkv-cli -f hset hash key1 one
./jkv-cli -f keys '*'
./jkv-cli -f hkeys hash
./jkv-cli -f hdel hash key1
./jkv-cli -f keys '*'
./jkv-cli -f hkeys hash) > hdel-f.log 2>&1

(redis-cli flushdb
redis-cli hset hash key1 one
redis-cli keys '*'
redis-cli hkeys hash
redis-cli hdel hash key1
redis-cli keys '*'
redis-cli hkeys hash) > hdel-r.log 2>&1

code hdel-?.log
