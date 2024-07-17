@test "HSET/KEYS -f" {
    [ "$(redis-cli flushdb)" = "OK" ]
    [ "$(redis-cli hset hash key1 one)" = "1" ]
    [ "$(redis-cli hset hash key1 one)" = "0" ]

    [ "$(./jkv-cli flushdb)" = "OK" ]
    [ "$(./jkv-cli hset hash key1 one)" = "1" ]
    [ "$(./jkv-cli hset hash key1 one)" = "0" ]

    [ "$(redis-cli set scalar one)" = "OK" ]
    # for debugging
    # redis-cli keys '*' | sha1sum

    [ "$(./jkv-cli set scalar one)" = "OK" ]
    [ "$(./jkv-cli keys '*' | sha1sum)" = "$(redis-cli keys '*' | sha1sum)" ]
}

@test "HSET/KEYS -r" {
    [ "$(redis-cli flushdb)" = "OK" ]
    [ "$(redis-cli hset hash key1 one)" = "1" ]
    [ "$(redis-cli hset hash key1 one)" = "0" ]

    [ "$(./jkv-cli -r flushdb)" = "OK" ]
    [ "$(./jkv-cli -r hset hash key1 one)" = "1" ]
    [ "$(./jkv-cli -r hset hash key1 one)" = "0" ]

    [ "$(redis-cli set scalar one)" = "OK" ]
    # for debugging
    # redis-cli keys '*' | sha1sum

    [ "$(./jkv-cli -r set scalar one)" = "OK" ]
    [ "$(./jkv-cli -r keys '*' | sha1sum)" = "$(redis-cli keys '*' | sha1sum)" ]
}

@test "Test SET/GET -f vs. -r" {
    [ "$(./jkv-cli -f flushdb)" = "OK" ]
    [ "$(./jkv-cli -f set a b)" = "$(./jkv-cli -r set a b)" ]
    [ "$(./jkv-cli -f get a)" = "$(./jkv-cli -r get a)" ]
}

@test "Test SET/DEL -f vs. -r" {
    [ "$(./jkv-cli -f flushdb)" = "OK" ]
    [ "$(./jkv-cli -f set a b)" = "$(./jkv-cli -r set a b)" ]
    [ "$(./jkv-cli -f del a)" = "$(./jkv-cli -r del a)" ]
    [ "$(./jkv-cli -f get a)" = "$(./jkv-cli -r get a)" ]
    [ "$(./jkv-cli -f set key1 one)" = "$(./jkv-cli -r set key1 one)" ]
    [ "$(./jkv-cli -f set key2 two)" = "$(./jkv-cli -r set key2 two)" ]
    [ "$(./jkv-cli -f set key3 three)" = "$(./jkv-cli -r set key3 three)" ]
    [ "$(./jkv-cli -f del key1 key2 key3)" = "$(./jkv-cli -r del key1 key2 key3)" ]
}

@test "Test EXISTS -f vs. -r" {
    [ "$(./jkv-cli -f set a b)" = "$(./jkv-cli -r set a b)" ]
    [ "$(./jkv-cli -f set a b)" = "$(redis-cli set a b)" ]
    [ "$(./jkv-cli -f set b c)" = "$(./jkv-cli -r set b c)" ]
    [ "$(./jkv-cli -f set b c)" = "$(redis-cli set b c)" ]
    [ "$(./jkv-cli -f exists a)" = "$(./jkv-cli -r exists a)" ]
    [ "$(./jkv-cli -f exists a | cut -d' ' -f2)" = "$(redis-cli exists a)" ]
    [ "$(./jkv-cli -f exists b)" = "$(./jkv-cli -r exists b)" ]
    [ "$(./jkv-cli -f exists b | cut -d' ' -f2)" = "$(redis-cli exists b)" ]
    [ "$(./jkv-cli -f exists a b)" = "$(./jkv-cli -r exists a b)" ]
    [ "$(./jkv-cli -f exists a b | cut -d' ' -f2)" = "$(redis-cli exists a b)" ]
}

@test "Test SET Syntax Error" {
    [ "$(./jkv-cli -f set a b c)" = "$(./jkv-cli -r set a b c)" ]
}

@test "Test HEXISTS -f vs. -r" {
    [ "$(./jkv-cli -f flushdb)" = "OK" ]
    [ "$(./jkv-cli -r flushdb)" = "OK" ]
    [ "$(./jkv-cli -f hset hash key1 one key2 two key3 three)" = "$(./jkv-cli -r hset hash key1 one key2 two key3 three)" ]
    [ "$(./jkv-cli -f hexists hash key1)" = "$(./jkv-cli -r hexists hash key1)" ]
    [ "$(./jkv-cli -f hexists hash key2)" = "$(./jkv-cli -r hexists hash key2)" ]
    [ "$(./jkv-cli -f hexists hash key3)" = "$(./jkv-cli -r hexists hash key3)" ]
}

@test "Test HKEYS -f" {
    redis-cli flushdb
    redis-cli hset hash key1 one key2 two key3 three
    [ "$(./jkv-cli -f flushdb)" = "OK" ]
    [ "$(./jkv-cli -f hset hash key1 one key2 two key3 three)" = "3" ]
    [ "$(./jkv-cli -f hkeys hash | sha1sum)" = "$(redis-cli hkeys hash | sha1sum)" ]
}

@test "Test HKEYS -r" {
    redis-cli flushdb
    redis-cli hset hash key1 one key2 two key3 three
    [ "$(./jkv-cli -r flushdb)" = "OK" ]
    [ "$(./jkv-cli -r hset hash key1 one key2 two key3 three)" = "3" ]
    [ "$(./jkv-cli -r hkeys hash | sha1sum)" = "$(redis-cli hkeys hash | sha1sum)" ]
}

@test "Test -x here" {
    [ "$(./jkv-cli -f flushdb)" = "OK" ]
    [ "$(printf one | ./jkv-cli -f -x hset hash key1)" = "1" ]
    [ "$(printf two | ./jkv-cli -f -x hset hash key2)" = "1" ]
    [ "$(printf three | ./jkv-cli -f -x hset hash key3)" = "1" ]
    [ "$(printf three | ./jkv-cli -f -x hset hash key3)" = "0" ]

    [ "$(./jkv-cli -r flushdb)" = "OK" ]
    [ "$(printf one | ./jkv-cli -r -x hset hash key1)" = "1" ]
    [ "$(printf two | ./jkv-cli -r -x hset hash key2)" = "1" ]
    [ "$(printf three | ./jkv-cli -r -x hset hash key3)" = "1" ]
    [ "$(printf three | ./jkv-cli -r -x hset hash key3)" = "0" ]

    [ "$(./jkv-cli -f flushdb)" = "OK" ]
    [ "$(printf one | ./jkv-cli -f -x set key1)" = "OK" ]
    [ "$(printf two | ./jkv-cli -f -x set key2)" = "OK" ]
    [ "$(printf two | ./jkv-cli -f -x set key2)" = "OK" ]
    ./jkv-cli -f keys '*'

    [ "$(./jkv-cli -r flushdb)" = "OK" ]
    [ "$(printf one | ./jkv-cli -r -x set key1)" = "OK" ]
    [ "$(printf two | ./jkv-cli -r -x set key2)" = "OK" ]
    [ "$(printf two | ./jkv-cli -r -x set key2)" = "OK" ]
    ./jkv-cli -f keys '*'

    [ "$(redis-cli flushdb)" = "OK" ]
    [ "$(printf one | redis-cli -x set key1)" = "OK" ]
    [ "$(printf two | redis-cli -x set key2)" = "OK" ]
    [ "$(printf two | redis-cli -x set key2)" = "OK" ]
    redis-cli keys '*'
}

@test "PING" {
    [ "$(./jkv-cli ping)" = "PONG" ]
}

@test "DB Location" {
    [ "$(./jkv-cli flushdb)" = "OK" ]
    [ ! -d "${HOME}/jkv_db" ]
    [ "$(./jkv-cli set key1 one)" = "OK" ]
    [ -d "${HOME}/jkv_db" ]
    ./jkv-cli flushdb

    [ "$(./jkv-cli -d "${HOME}/db" flushdb)" = "OK" ]
    [ ! -d "${HOME}/db" ]

    [ "$(./jkv-cli -d "${HOME}/db" set key1 one)" = "OK" ]
    [ -d "${HOME}/db" ]

    ./jkv-cli -d "${HOME}/db" flushdb
}
