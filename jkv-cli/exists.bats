setup() {
    jkv-cli flushdb
    jkv-cli set key1 one
    jkv-cli set key2 two
    redis-cli flushdb
    redis-cli set key1 one
    redis-cli set key2 two
}

@test "Test jkv-cli -f EXISTS" {
    [ "$(jkv-cli exists)" = "ERR wrong number of arguments for 'exists' command" ]
    [ "$(jkv-cli exists key1)" = "1" ]
    [ "$(jkv-cli exists key2)" = "1" ]
    [ "$(jkv-cli exists key1 key2)" = "2" ]
    [ "$(jkv-cli exists key3)" = "0" ]
    [ "$(jkv-cli exists key1 key2 key3)" = "2" ]
}

@test "Test redis-cli EXISTS" {
    [ "$(redis-cli exists)" = "ERR wrong number of arguments for 'exists' command" ]
    [ "$(redis-cli exists key1)" = "1" ]
    [ "$(redis-cli exists key2)" = "1" ]
    [ "$(redis-cli exists key1 key2)" = "2" ]
    [ "$(redis-cli exists key3)" = "0" ]
    [ "$(redis-cli exists key1 key2 key3)" = "2" ]
}
