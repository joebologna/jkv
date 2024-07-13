setup() {
    redis-cli flushdb
    jkv-cli flushdb
}

@test "Test HSET" {
    [ "$(redis-cli set hash1 tmp)" = "OK" ]
    [ "$(redis-cli hset hash1 key1 one)" = "WRONGTYPE Operation against a key holding the wrong kind of value" ]
    [ "$(redis-cli del hash1)" = "1" ]

    [ "$(jkv-cli set hash1 tmp)" = "OK" ]
    [ "$(jkv-cli hset hash1 key1 one)" = "WRONGTYPE Operation against a key holding the wrong kind of value" ]
    [ "$(jkv-cli del hash1)" = "1" ]
}
