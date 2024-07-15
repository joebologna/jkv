@test "HSET/KEYS -f" {
    [ "$(redis-cli flushdb)" = "OK" ]
    [ "$(redis-cli hset hash key1 one)" = "1" ]
    [ "$(redis-cli hset hash key1 one)" = "0" ]

    [ "$(jkv-cli flushdb)" = "OK" ]
    [ "$(jkv-cli hset hash key1 one)" = "1" ]
    [ "$(jkv-cli hset hash key1 one)" = "0" ]

    [ "$(redis-cli set scalar one)" = "OK" ]
    # for debugging
    # redis-cli keys '*' | sha1sum

    [ "$(jkv-cli set scalar one)" = "OK" ]
    [ "$(jkv-cli keys '*' | sha1sum)" = "$(redis-cli keys '*' | sha1sum)" ]
}

@test "HSET/KEYS -r" {
    [ "$(redis-cli flushdb)" = "OK" ]
    [ "$(redis-cli hset hash key1 one)" = "1" ]
    [ "$(redis-cli hset hash key1 one)" = "0" ]

    [ "$(jkv-cli -r flushdb)" = "OK" ]
    [ "$(jkv-cli -r hset hash key1 one)" = "1" ]
    [ "$(jkv-cli -r hset hash key1 one)" = "0" ]

    [ "$(redis-cli set scalar one)" = "OK" ]
    # for debugging
    # redis-cli keys '*' | sha1sum

    [ "$(jkv-cli -r set scalar one)" = "OK" ]
    [ "$(jkv-cli -r keys '*' | sha1sum)" = "$(redis-cli keys '*' | sha1sum)" ]
}
