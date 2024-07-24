Thu 11 Jul 2024 01:39:25 PM CDT

# initdb

This script requires:

- [x] redis-cli FLUSHDB
- upload-icons
  - [x] redis-cli -x HSET hkey key <file
- upload-networks
  - [x] go run tmp2json.go | redis-cli -x HSET dhcp 80-wired-dhcp.tmpl
  - [x] go run tmp2json.go | redis-cli -x HSET dhcp 80-wired-static.tmpl
  - [x] redis-cli hset Networks :default dhcp
- screens/main.go
  - [x] dbutil.UpsertItem(key, screen_data.String()) -> convert to HSet
  - [x] rdb.Del()
- upload-backgrounds
  - redis-cli -x HSET hkey key <file
- [x] redis-cli HSET RunModes Production '{"Debug":true}' Dev '{"Debug":true}' :current Dev
  - returns `(integer) 3`
- backupdb
  - rdbrip -s rdb.sys.zip
  - rdbrip -c rdb-config.zip
