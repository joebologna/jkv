setup() {
    echo Need a flag to verify output removes parens when interactive
    [ "$(jkv-cli -f flushdb)" = "OK" ]
    [ "$(jkv-cli -r flushdb)" = "OK" ]
    [ "$(redis-cli flushdb)" = "OK" ]
}

@test "KEYS/HEXISTS/HKEYS" {
    [ "$(redis-cli hset hash key1 one key2 two key3 three)" = "3" ]
    [ "$(redis-cli hexists hash key1)" = "1" ]
    [ "$(redis-cli hexists hash key2)" = "1" ]
    [ "$(redis-cli hexists hash key3)" = "1" ]
    [ "$(redis-cli set scalar one)" = "OK" ]
    SHA1="$(redis-cli keys '*' | sha1sum | cut -d' ' -f1)"
    SHA2="$(redis-cli hkeys hash | sha1sum | cut -d' ' -f1)"

    [ "$(jkv-cli -f hset hash key1 one key2 two key3 three)" = "3" ]
    [ "$(jkv-cli -f hexists hash key1)" = "1" ]
    [ "$(jkv-cli -f hexists hash key2)" = "1" ]
    [ "$(jkv-cli -f hexists hash key3)" = "1" ]
    [ "$(jkv-cli -f set scalar one)" = "OK" ]
    jkv-cli -f keys '*'
    echo under a pipe this command should emit a naked list
    [ "$(jkv-cli -f keys '*' | sha1sum | cut -d' ' -f1)" = "${SHA1}" ]
    jvk-cli -f hkeys hash
    echo under a pipe this command should emit a naked list
    [ "$(jvk-cli -f hkeys hash | sha1sum | cut -d' ' -f1)" = "${SHA2}" ]

    [ "$(jkv-cli -r hset hash key1 one key2 two key3 three)" = "3" ]
    [ "$(jkv-cli -r hexists hash key1)" = "1" ]
    [ "$(jkv-cli -r hexists hash key2)" = "1" ]
    [ "$(jkv-cli -r hexists hash key3)" = "1" ]
    [ "$(jkv-cli -r set scalar one)" = "OK" ]
    jkv-cli -r keys '*'
    [ "$(jkv-cli -r keys '*' | sha1sum | cut -d' ' -f1)" = "${SHA1}" ]
    jvk-cli -r hkeys hash
    [ "$(jvk-cli -r hkeys hash | sha1sum | cut -d' ' -f1)" = "${SHA2}" ]
}

@test "HGET/HSET" {
    [ "$(jkv-cli -f hset hash key1 one key2 two key3 three)" = "3" ]
    [ "$(jkv-cli -f hget hash key1)" = "\"one\"" ]
    false need to match list

    [ "$(jkv-cli -r flushdb)" = "OK" ]
    [ "$(jkv-cli -r hset hash key1 one key2 two key3 three)" = "3" ]
    [ "$(jkv-cli -r hget hash key1)" = "\"one\"" ]
    false need to match list
}

@test "HDEL" {
    [ "$(jkv-cli -f hset hash key1 one key2 two key3 three)" = "3" ]
    [ "$(jkv-cli -f hdel hash key1 key2 key3)" = "3" ]

    [ "$(jkv-cli -r flushdb)" = "OK" ]
    [ "$(jkv-cli -r hset hash key1 one key2 two key3 three)" = "3" ]
    [ "$(jkv-cli -r hdel hash key1 key2 key3)" = "3" ]
}

@test "TO DO: HKEYS, HEXISTS, KEYS and -x" {
    false
}
