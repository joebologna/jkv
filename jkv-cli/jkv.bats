# setup() {
	# jkv-cli -f flushdb
	# jkv-cli -r flushdb
# }

@test "Test SET/GET" {
	[ "$(jkv-cli -f set a b)" = "$(jkv-cli -r set a b)" ]
	[ "$(jkv-cli -f get a)" = "$(jkv-cli -r get a)" ]
}

@test "Test SET/DEL" {
	[ "$(jkv-cli -f set a b)" = "$(jkv-cli -r set a b)" ]
	[ "$(jkv-cli -f del a)" = "$(jkv-cli -r del a)" ]
	[ "$(jkv-cli -f get a)" = "$(jkv-cli -r get a)" ]
}

@test "Test EXISTS" {
	[ "$(jkv-cli -f set a b)" = "$(jkv-cli -r set a b)" ]
	[ "$(jkv-cli -f set a b)" = "$(/usr/local/bin/redis-cli set a b)" ]
	[ "$(jkv-cli -f set b c)" = "$(jkv-cli -r set b c)" ]
	[ "$(jkv-cli -f set b c)" = "$(/usr/local/bin/redis-cli set b c)" ]
	[ "$(jkv-cli -f exists a)" = "$(jkv-cli -r exists a)" ]
	[ "$(jkv-cli -f exists a | cut -d' ' -f2)" = "$(/usr/local/bin/redis-cli exists a)" ]
	[ "$(jkv-cli -f exists b)" = "$(jkv-cli -r exists b)" ]
	[ "$(jkv-cli -f exists b | cut -d' ' -f2)" = "$(/usr/local/bin/redis-cli exists b)" ]
	[ "$(jkv-cli -f exists a b)" = "$(jkv-cli -r exists a b)" ]
	skip "TODO: support multiple keys"
	[ "$(jkv-cli -f exists a b | cut -d' ' -f2)" = "$(/usr/local/bin/redis-cli exists a b)" ]
}

@test "Test SET Syntax Error" {
	[ "$(jkv-cli -f set a b c)" = "$(jkv-cli -r set a b c)" ]
}
