@test "HSET/KEYS" {
    [ "$(redis-cli flushdb)" = "OK" ]
    [ "$(redis-cli hset hash key1 one)" = "1" ]
    [ "$(redis-cli hset hash key1 one)" = "0" ]
    [ "$(jkv-cli flushdb)" = "OK" ]
    [ "$(jkv-cli hset hash key1 one)" = "1" ]
    [ "$(jkv-cli hset hash key1 one)" = "0" ]

    [ "$(redis-cli set scalar one)" = "OK" ]
    # for debugging
    # redis-cli keys '*' | sha1sum

    [ "$(jkv-cli -f set scalar one)" = "OK" ]
    [ "$(jkv-cli -f keys '*' | sha1sum)" = "$(redis-cli keys '*' | sha1sum)" ]
}
