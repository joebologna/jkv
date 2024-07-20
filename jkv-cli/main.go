package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/panduit-joeb/jkv"
	"github.com/panduit-joeb/jkv/store/fs"
	"github.com/panduit-joeb/jkv/store/redis"
)

func main() {
	cmd := os.Args[0]
	if strings.Contains(os.Args[0], "/") {
		a := strings.Split(os.Args[0], "/")
		cmd = a[len(a)-1]
	}
	// fmt.Println("cmd is", cmd)

	var redis_cmd, fs_cmd, version, opt_x, prompt bool
	var redis_host, db_dir string
	flag.BoolVar(&redis_cmd, "r", cmd == "redis-cli", "Run JKV tests using Redis")
	flag.BoolVar(&fs_cmd, "f", cmd == "jkv-cli", "Run JKV tests using FS")
	flag.BoolVar(&version, "v", false, "Print version")
	flag.BoolVar(&opt_x, "x", false, "Get value from stdin")
	flag.StringVar(&redis_host, "h", redis.DEFAULT_DB, "Redis server host and port")
	flag.StringVar(&db_dir, "d", fs.DEFAULT_DB, "Location of FS DB")
	flag.Parse()

	if version {
		fmt.Println(jkv.VERSION)
		os.Exit(0)
	}

	prompt = len(flag.Args()) == 0

	var db jkv.Client
	var db_loc string

	if redis_cmd {
		db_loc = redis_host
		db = redis.NewClient(&redis.Options{Addr: db_loc, Password: "", DB: 0})
	} else if fs_cmd {
		db_loc = db_dir
		db = fs.NewClient(&fs.Options{Addr: db_loc})
	}
	db.Open()

	if prompt {
		scanner := bufio.NewScanner(os.Stdin)

		fmt.Printf(db_loc + "> ")
		for scanner.Scan() {
			ProcessCmd(db, scanner.Text(), opt_x, isPipe())
			fmt.Printf(db_loc + "> ")
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading input:", err)
		}
	} else {
		ProcessCmd(db, strings.Join(flag.Args(), " "), opt_x, isPipe())
	}
}

func ProcessCmd(db jkv.Client, cmd string, opt_x, is_pipe bool) {
	tokens := strings.Fields(cmd)
	if len(tokens) == 0 {
		return
	}
	ctx := context.Background()
	switch strings.ToUpper(tokens[0]) {
	case "PING":
		fmt.Println("PONG")
	case "FLUSHDB":
		if len(tokens) == 1 {
			db.FlushDB(ctx)
			fmt.Println("OK")
		} else {
			report("(error)", "ERR syntax error", is_pipe)
		}
	case "HGET":
		if len(tokens) == 3 {
			rec := db.HGet(ctx, tokens[1], tokens[2])
			if rec.Err() != nil {
				fmt.Println("(nil)")
			} else {
				fmt.Printf("\"%s\"\n", rec.Val())
			}
		} else {
			fmt.Println("(nil)")
		}
	case "HSET":
		ctx := context.Background()
		if opt_x {
			if len(tokens) == 3 {
				var buf = make([]byte, 1024*1024)
				var n = 0
				n, err := os.Stdin.Read(buf)
				if n == 0 {
					if err != io.EOF {
						panic(err.Error())
					}
					return
				}
				// buf = []byte("hello")

				hash := tokens[1]
				key := tokens[2]
				rec := db.HSet(ctx, hash, key, string(buf))
				report("(integer)", fmt.Sprintf("%d", rec.Val()), is_pipe)
			} else {
				report("(error)", "ERR wrong number of arguments for 'hset' command", is_pipe)
			}
		} else {
			if len(tokens) > 2 {
				rec := db.HSet(ctx, tokens[1], tokens[2:]...)
				if rec.Err() != nil {
					fmt.Println(rec.Err().Error())
					return
				}
				report("(integer)", fmt.Sprintf("%d", rec.Val()), is_pipe)
				return
			} else {
				if db.Exists(ctx, tokens[1]).Val() != 0 {
					report("(error)", "WRONGTYPE Operation against a key holding the wrong kind of value", is_pipe)
					return
				}
				if len(tokens) >= 4 && ((len(tokens)-2)%2 == 0) {
					rec := db.HSet(ctx, tokens[1], tokens[2:]...)
					if rec.Err() != nil {
						fmt.Println(rec.Err())
					} else {
						report("(integer)", fmt.Sprintf("%d", rec.Val()), is_pipe)
					}
				} else {
					report("(error)", "ERR wrong number of arguments for 'hset' command", is_pipe)
				}
			}
		}
	case "HDEL":
		if len(tokens) == 3 {
			ctx := context.Background()
			rec := db.HDel(ctx, tokens[1], tokens[2])
			if rec.Err() != nil {
				fmt.Println("(nil)")
			} else {
				report("(integer)", fmt.Sprintf("%d", rec.Val()), is_pipe)
			}
		} else {
			fmt.Println("(nil)")
		}
	case "HKEYS":
		if len(tokens) == 2 {
			ctx := context.Background()
			rec := db.HKeys(ctx, tokens[1])
			if rec.Err() != nil {
				if os.IsNotExist(rec.Err()) {
					report("(empty array)", "", is_pipe)
				} else {
					fmt.Println("(nil)")
				}
			} else {
				if len(rec.Val()) == 0 {
					report("(empty array)", "", is_pipe)
				}
				for i, v := range rec.Val() {
					if is_pipe {
						fmt.Println(v)
					} else {
						fmt.Printf("%d) \"%s\"\n", i+1, v)
					}
				}
			}
		} else {
			report("(error)", "ERR wrong number of arguments for 'hkeys' command", is_pipe)
		}
	case "HEXISTS":
		if len(tokens) == 3 {
			ctx := context.Background()
			rec := db.HExists(ctx, tokens[1], tokens[2])
			if rec.Val() {
				report("(integer)", "1", is_pipe)
			} else {
				report("(integer)", "0", is_pipe)
			}
		} else {
			report("(error)", "ERR wrong number of arguments for 'exists' command", is_pipe)
		}
	case "GET":
		if len(tokens) == 2 {
			ctx := context.Background()
			rec := db.Get(ctx, tokens[1])
			if rec.Err() != nil {
				fmt.Println("(nil)")
			} else {
				fmt.Printf("\"%s\"\n", rec.Val())
			}
		} else {
			fmt.Println("(nil)")
		}
	case "SET":
		if opt_x {
			if len(tokens) == 2 {
				ctx := context.Background()
				var buf = make([]byte, 1024*1024)
				var n = 0
				n, err := os.Stdin.Read(buf)
				if n == 0 {
					if err != io.EOF {
						panic(err.Error())
					}
					return
				}
				key := tokens[1]
				rec := db.Set(ctx, key, string(buf[:n-1]), 0)
				if rec.Err() != nil {
					fmt.Println("(nil)")
				} else {
					fmt.Println("OK")
				}
			} else {
				fmt.Println("(error) ERR wrong number of arguments for 'set' command")
			}
		} else {
			if len(tokens) == 3 {
				ctx := context.Background()
				rec := db.Set(ctx, tokens[1], tokens[2], 0)
				if rec.Err() != nil {
					fmt.Println("(nil)")
				} else {
					fmt.Println(rec.Val())
				}
			} else {
				report("(error)", "ERR wrong number of arguments for 'set' command", is_pipe)
			}
		}
	case "DEL":
		if len(tokens) >= 2 {
			ctx := context.Background()
			rec := db.Del(ctx, tokens[1:]...)
			if rec.Err() != nil {
				fmt.Println("(nil)", rec.Err().Error())
			} else {
				report("(integer)", fmt.Sprintf("%d", rec.Val()), is_pipe)
			}
		} else {
			fmt.Println("(nil)")
		}
	case "KEYS":
		if len(tokens) == 2 {
			ctx := context.Background()
			rec := db.Keys(ctx, tokens[1])
			if rec.Err() != nil {
				fmt.Println("(nil)")
			} else {
				if len(rec.Val()) == 0 {
					report("(empty array)", "", is_pipe)
				} else {
					for i, v := range rec.Val() {
						if is_pipe {
							fmt.Println(v)
						} else {
							fmt.Printf("%d) \"%s\"\n", i+1, v)
						}
					}
				}
			}
		} else {
			report("(error)", "ERR wrong number of arguments for 'keys' command", is_pipe)
		}
	case "EXISTS":
		if len(tokens) >= 2 {
			ctx := context.Background()
			var n int64
			for _, token := range tokens[1:] {
				rec := db.Exists(ctx, token)
				if rec.Err() != nil {
					// fmt.Println("(nil)")
					break
				}
				n = n + rec.Val()
			}
			report("(integer)", fmt.Sprintf("%d", n), is_pipe)
		} else {
			report("(error)", "ERR wrong number of arguments for 'exists' command", is_pipe)
		}
	default:
		report("(error)", fmt.Sprintf("ERR unknown command '%s', with args beginning with:\n", tokens[0]), is_pipe)
	}
}

func isPipe() bool {
	fi, _ := os.Stdout.Stat()
	return (fi.Mode() & os.ModeCharDevice) == 0
}

func report(prefix, msg string, is_pipe bool) {
	if !is_pipe {
		msg = prefix + " " + msg
	}
	fmt.Println(msg)
}
