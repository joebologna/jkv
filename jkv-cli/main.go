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
	flag.BoolVar(&redis_cmd, "r", cmd == "redis-cli", "Run JKV tests using Redis")
	flag.BoolVar(&fs_cmd, "f", cmd == "jkv-cli", "Run JKV tests using FS")
	flag.BoolVar(&version, "v", false, "Print version")
	// flag.BoolVar(&opt_x, "x", false, "Get value from stdin")
	flag.Parse()

	if version {
		fmt.Println(jkv.VERSION)
		os.Exit(0)
	}

	if len(flag.Args()) > 1 && flag.Args()[1] == "-x" && ((len(flag.Args()) == 4 && strings.ToUpper(flag.Args()[0]) == "HSET") || (len(flag.Args()) == 3 && strings.ToUpper(flag.Args()[0]) == "SET")) {
		opt_x = true
	}

	prompt = len(flag.Args()) == 0

	if redis_cmd {
		r := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
		r.Open()

		if prompt {
			scanner := bufio.NewScanner(os.Stdin)

			fmt.Printf(r.DBDir + "> ")
			for scanner.Scan() {
				ProcessCmd(r, scanner.Text(), opt_x, isPipe())
				fmt.Printf(r.DBDir + "> ")
			}

			if err := scanner.Err(); err != nil {
				fmt.Println("Error reading input:", err)
			}
		} else {
			ProcessCmd(r, strings.Join(flag.Args(), " "), opt_x, isPipe())
		}
	} else if fs_cmd {
		f := fs.NewClient(&fs.Options{Addr: fs.DEFAULT_DB})
		f.Open()

		if prompt {
			scanner := bufio.NewScanner(os.Stdin)

			fmt.Printf(f.DBDir + "> ")
			for scanner.Scan() {
				ProcessCmd(f, scanner.Text(), opt_x, isPipe())
				fmt.Printf(f.DBDir + "> ")
			}

			if err := scanner.Err(); err != nil {
				fmt.Println("Error reading input:", err)
			}
		} else {
			ProcessCmd(f, strings.Join(flag.Args(), " "), opt_x, isPipe())
		}
	}
}

func ProcessCmd(db interface{}, cmd string, opt_x, is_pipe bool) {
	var (
		value string
		err   error
	)
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
			if r, ok := db.(*redis.Client); ok {
				r.FlushDB(ctx)
			} else {
				db.(*fs.Client).FlushDB(ctx)
			}
			fmt.Println("OK")
		} else {
			report("(error)", "ERR syntax error", is_pipe)
		}
	case "HGET":
		if len(tokens) == 3 {
			var rec *jkv.StringCmd
			if r, ok := db.(*redis.Client); ok {
				rec = r.HGet(ctx, tokens[1], tokens[2])
			} else {
				rec = db.(*fs.Client).HGet(ctx, tokens[1], tokens[2])
			}
			value = rec.Val()
			err = rec.Err()
			if err != nil {
				fmt.Println("(nil)")
			} else {
				fmt.Printf("\"%s\"\n", value)
			}
		} else {
			fmt.Println("(nil)")
		}
	case "HSET":
		ctx := context.Background()
		if opt_x {
			if len(tokens) == 4 {
				var buf = make([]byte, 1024*1024)
				var n = 0
				n, err = os.Stdin.Read(buf)
				if n == 0 {
					if err != io.EOF {
						panic(err.Error())
					}
					return
				}
				// buf = []byte("hello")

				var rec *jkv.IntCmd
				hash := tokens[2]
				key := tokens[3]

				if r, ok := db.(*redis.Client); ok {
					rec = r.HSet(ctx, hash, key, string(buf))
				} else {
					rec = db.(*fs.Client).HSet(ctx, hash, key, string(buf))
				}
				report("(integer)", fmt.Sprintf("%d", rec.Val()), is_pipe)
			} else {
				report("(error)", "ERR wrong number of arguments for 'hset' command", is_pipe)
			}
		} else {
			if len(tokens) > 2 {
				if r, ok := db.(*redis.Client); ok {
					rec := r.HSet(ctx, tokens[1], tokens[2:]...)
					if rec.Err() != nil {
						fmt.Println(rec.Err().Error())
						return
					}
					report("(integer)", fmt.Sprintf("%d", rec.Val()), is_pipe)
					return
				} else {
					if db.(*fs.Client).Exists(ctx, tokens[1]).Val() != 0 {
						report("(error)", "WRONGTYPE Operation against a key holding the wrong kind of value", is_pipe)
						return
					}
					if len(tokens) >= 4 && ((len(tokens)-2)%2 == 0) {
						var rec *jkv.IntCmd
						if r, ok := db.(*redis.Client); ok {
							rec = r.HSet(ctx, tokens[1], tokens[2:]...)
						} else {
							rec = db.(*fs.Client).HSet(ctx, tokens[1], tokens[2:]...)
						}
						if rec.Err() != nil {
							fmt.Println(rec.Err())
						} else {
							report("(integer)", fmt.Sprintf("%d", rec.Val()), is_pipe)
						}
					} else {
						report("(error)", "ERR wrong number of arguments for 'hset' command", is_pipe)
					}
				}
			} else {
				fmt.Println("(nil)")
			}
		}
	case "HDEL":
		if len(tokens) == 2 {
			ctx := context.Background()
			if r, ok := db.(*redis.Client); ok {
				err = r.HDel(ctx, tokens[1], tokens[2]).Err()
			} else {
				err = db.(*fs.Client).HDel(ctx, tokens[1], tokens[2]).Err()
			}
			if err != nil {
				fmt.Println("(nil)")
			} else {
				fmt.Printf("\"%s\"\n", value)
			}
		} else {
			fmt.Println("(nil)")
		}
	case "HKEYS":
		if len(tokens) == 2 {
			ctx := context.Background()
			var values []string
			var err error
			var rec *jkv.StringSliceCmd
			if r, ok := db.(*redis.Client); ok {
				rec = r.HKeys(ctx, tokens[1])
			} else {
				rec = db.(*fs.Client).HKeys(ctx, tokens[1])
			}
			values = rec.Val()
			err = rec.Err()
			if err != nil {
				fmt.Println("(nil)")
			} else {
				for i, v := range values {
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
			var rec *jkv.BoolCmd
			if r, ok := db.(*redis.Client); ok {
				rec = r.HExists(ctx, tokens[1], tokens[2])
			} else {
				rec = db.(*fs.Client).HExists(ctx, tokens[1], tokens[2])
			}
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
			var value string
			var err error
			var rec *jkv.StringCmd
			if r, ok := db.(*redis.Client); ok {
				rec = r.Get(ctx, tokens[1])
			} else {
				rec = db.(*fs.Client).Get(ctx, tokens[1])
			}
			value = rec.Val()
			err = rec.Err()
			if err != nil {
				fmt.Println("(nil)")
			} else {
				fmt.Printf("\"%s\"\n", value)
			}
		} else {
			fmt.Println("(nil)")
		}
	case "SET":
		if opt_x {
			if len(tokens) == 3 {
				ctx := context.Background()
				var buf = make([]byte, 1024*1024)
				var n = 0
				n, err = os.Stdin.Read(buf)
				if n == 0 {
					if err != io.EOF {
						panic(err.Error())
					}
					return
				}
				key := tokens[2]
				var rec *jkv.StatusCmd
				if r, ok := db.(*redis.Client); ok {
					rec = r.Set(ctx, key, string(buf[:n-1]))
				} else {
					rec = db.(*fs.Client).Set(ctx, key, string(buf[:n-1]))
				}
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
				var rec *jkv.StatusCmd
				if r, ok := db.(*redis.Client); ok {
					rec = r.Set(ctx, tokens[1], tokens[2])
				} else {
					rec = db.(*fs.Client).Set(ctx, tokens[1], tokens[2])
				}
				fmt.Println(rec.Val())
			} else {
				report("(error)", "ERR wrong number of arguments for 'set' command", is_pipe)
			}
		}
	case "DEL":
		if len(tokens) >= 2 {
			ctx := context.Background()
			var rec *jkv.IntCmd
			if r, ok := db.(*redis.Client); ok {
				rec = r.Del(ctx, tokens[1:]...)
			} else {
				rec = db.(*fs.Client).Del(ctx, tokens[1:]...)
			}
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
			var rec *jkv.StringSliceCmd
			if r, ok := db.(*redis.Client); ok {
				rec = r.Keys(ctx, tokens[1])
			} else {
				rec = db.(*fs.Client).Keys(ctx, tokens[1])
			}
			if rec.Err() != nil {
				fmt.Println("(nil)")
			} else {
				for i, v := range rec.Val() {
					if is_pipe {
						fmt.Println(v)
					} else {
						fmt.Printf("%d) \"%s\"\n", i+1, v)
					}
				}
			}
		} else {
			report("(error)", "ERR wrong number of arguments for 'keys' command", is_pipe)
		}
	case "EXISTS":
		if len(tokens) >= 2 {
			ctx := context.Background()
			var rec *jkv.IntCmd
			var n int64
			for _, token := range tokens[1:] {
				if r, ok := db.(*redis.Client); ok {
					rec = r.Exists(ctx, token)
				} else {
					rec = db.(*fs.Client).Exists(ctx, token)
				}
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
