setup() {
    redis-cli flushdb
    jkv-cli flushdb
}

@test "Test HSET" {
    redis-cli flushdb
    jkv-cli flushdb

    [ "$(redis-cli set hash1 tmp)" = "OK" ]
    [ "$(redis-cli hset hash1 key1 one)" = "WRONGTYPE Operation against a key holding the wrong kind of value" ]
    [ "$(redis-cli del hash1)" = "1" ]

    [ "$(jkv-cli set hash1 tmp)" = "OK" ]
    [ "$(jkv-cli hset hash1 key1 one)" = "WRONGTYPE Operation against a key holding the wrong kind of value" ]
    [ "$(jkv-cli del hash1)" = "1" ]
}

@test "Test DEL" {
    redis-cli flushdb
    jkv-cli flushdb

    [ "$(redis-cli set key1 one)" = "OK" ]
    [ "$(redis-cli set key2 two)" = "OK" ]
    [ "$(redis-cli set key3 three)" = "OK" ]
    [ "$(redis-cli del key1 key2 key3)" = "3" ]

    [ "$(jkv-cli set key1 one)" = "OK" ]
    [ "$(jkv-cli set key2 two)" = "OK" ]
    [ "$(jkv-cli set key3 three)" = "OK" ]
    [ "$(jkv-cli del key1 key2 key3)" = "3" ]
}

@test "Test HSET -x" {
    jkv-cli -f hset hash1 key1 one key2 two
    jkv-cli -r hset hash1 key1 one key2 two
}

@test "Test SET -x" {
    jkv-cli -f hset hash1 key1 one key2 two
    jkv-cli -r hset hash1 key1 one key2 two
}

@test "Test HSET 2+ keys -f" {
    [ "$(jkv-cli -f hset hash1 key1 one key2 two)" = "2" ]
}

@test "Test HSET 2+ keys -r" {
    [ "$(jkv-cli -r hset hash1 key1 one key2 two)" = "2" ]
}
