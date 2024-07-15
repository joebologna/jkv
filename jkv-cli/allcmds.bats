@test "FLUSHDB" {
    [ "$(jkv-cli -f flushdb)" = "OK" ]
    [ "$(jkv-cli -r flushdb)" = "OK" ]
}

@test "EXISTS" {
    [ "$(jkv-cli -f hset hash key1 one key2 two)" = "2" ]
    [ "$(jkv-cli -f keys \*)" = "1) \"hash\"" ]
}

@test "HGET/HSET" {
    [ "$(jkv-cli -f flushdb)" = "OK" ]
    [ "$(jkv-cli -f hset hash key1 one key2 two)" = "2" ]
    [ "$(jkv-cli -f hget hash key1)" = "\"one\"" ]

    [ "$(jkv-cli -r flushdb)" = "OK" ]
    [ "$(jkv-cli -r hset hash key1 one key2 two)" = "2" ]
    [ "$(jkv-cli -r hget hash key1)" = "\"one\"" ]
}

@test "HDEL" {
    [ "$(jkv-cli -f flushdb)" = "OK" ]
    [ "$(jkv-cli -f hset hash key1 one key2 two)" = "2" ]
    [ "$(jkv-cli -f hdel hash key1 key2)" = "2" ]

    [ "$(jkv-cli -r flushdb)" = "OK" ]
    [ "$(jkv-cli -r hset hash key1 one key2 two)" = "2" ]
    [ "$(jkv-cli -r hdel hash key1 key2)" = "2" ]
}

@test "TO DO: HKEYS, HEXISTS, KEYS and -x" {
    false
}
