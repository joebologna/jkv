setup() {
    redis-cli flushdb
}

@test "Test HSET" {
    [ "$(redis-cli set hash1 tmp)" = "OK" ]
    # warning, interactively redis-cli outputs: (error) then WRONGTYPE...
    [ "$(redis-cli hset hash1 key1 one)" = "WRONGTYPE Operation against a key holding the wrong kind of value" ]
}
