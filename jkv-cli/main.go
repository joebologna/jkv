package main

import (
	"bufio"
	"flag"
	"fmt"
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

	var run_test, redis_cmd, fs_cmd, version, opt_x, prompt bool
	flag.BoolVar(&run_test, "t", false, "Run JKV tests")
	flag.BoolVar(&redis_cmd, "r", cmd == "redis-cli", "Run JKV tests using Redis")
	flag.BoolVar(&fs_cmd, "f", cmd == "jkv-cli", "Run JKV tests using FS")
	flag.BoolVar(&version, "v", false, "Print version")
	flag.BoolVar(&opt_x, "x", false, "Get value from stdin")
	flag.Parse()

	if version {
		fmt.Println(jkv.VERSION)
		os.Exit(0)
	}

	prompt = len(flag.Args()) == 0

	if run_test {
		c := fs.NewJKVClient()
		ok("open", c.Open)
		ok("set", func() error { return c.SET("this", "that") })
		c.FLUSHDB()

		r := redis.NewJKVClient()
		ok("open", r.Open)
		ok("set", func() error { return r.SET("this", "that") })
		r.FLUSHDB()
	} else if redis_cmd {
		r := redis.NewJKVClient()
		r.Open()

		if prompt {
			scanner := bufio.NewScanner(os.Stdin)

			fmt.Printf(r.DBDir + "> ")
			for scanner.Scan() {
				ProcessCmd(r, scanner.Text())
				fmt.Printf(r.DBDir + "> ")
			}

			if err := scanner.Err(); err != nil {
				fmt.Println("Error reading input:", err)
			}
		} else {
			ProcessCmd(r, strings.Join(flag.Args(), " "))
		}
	} else if fs_cmd {
		f := fs.NewJKVClient()
		f.Open()

		if prompt {
			scanner := bufio.NewScanner(os.Stdin)

			fmt.Printf(f.DBDir + "> ")
			for scanner.Scan() {
				ProcessCmd(f, scanner.Text())
				fmt.Printf(f.DBDir + "> ")
			}

			if err := scanner.Err(); err != nil {
				fmt.Println("Error reading input:", err)
			}
		} else {
			ProcessCmd(f, strings.Join(flag.Args(), " "))
		}
	}
}

func ok(label string, doit func() error) { fmt.Printf("%s returned %t\n", label, doit() == nil) }

func ProcessCmd(db interface{}, cmd string) {
	var (
		value  string
		values []string
		err    error
	)
	tokens := strings.Fields(cmd)
	if len(tokens) == 0 {
		return
	}
	switch strings.ToUpper(tokens[0]) {
	case "FLUSHDB":
		if len(tokens) == 1 {
			if r, ok := db.(*redis.JKV_DB); ok {
				r.FLUSHDB()
			} else {
				db.(*fs.JKV_DB).FLUSHDB()
			}
			fmt.Println("OK")
		} else {
			fmt.Println("(error) ERR syntax error")
		}
	case "HGET":
		if len(tokens) == 3 {
			if r, ok := db.(*redis.JKV_DB); ok {
				value, err = r.HGET(tokens[1], tokens[2])
			} else {
				value, err = db.(*fs.JKV_DB).HGET(tokens[1], tokens[2])
			}
			if err != nil {
				fmt.Println("(nil)")
			} else {
				fmt.Printf("\"%s\"\n", value)
			}
		} else {
			fmt.Println("(nil)")
		}
	case "HSET":
		if len(tokens) == 4 {
			var msg string
			if r, ok := db.(*redis.JKV_DB); ok {
				if r.HEXISTS(tokens[1], tokens[2]) {
					msg = "(integer) 0"
				} else {
					msg = "(integer) 1"
				}
				err = r.HSET(tokens[1], tokens[2], tokens[3])
			} else {
				if db.(*fs.JKV_DB).HEXISTS(tokens[1], tokens[2]) {
					msg = "(integer) 0"
				} else {
					msg = "(integer) 1"
				}
				err = db.(*fs.JKV_DB).HSET(tokens[1], tokens[2], tokens[3])
			}
			if err != nil {
				if strings.Contains(err.Error(), "exists as a scalar, cannot be a hash") || strings.Contains(err.Error(), "WRONGTYPE") {
					fmt.Println("(error) WRONGTYPE Operation against a key holding the wrong kind of value")
				} else {
					fmt.Println("(nil)")
				}
			} else {
				fmt.Println(msg)
			}
		} else {
			fmt.Println("(error) ERR wrong number of arguments for 'hset' command")
		}
	case "HDEL":
		if len(tokens) == 2 {
			if r, ok := db.(*redis.JKV_DB); ok {
				err = r.HDEL(tokens[1], tokens[2])
			} else {
				err = db.(*fs.JKV_DB).HDEL(tokens[1], tokens[2])
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
			if r, ok := db.(*redis.JKV_DB); ok {
				values, err = r.HKEYS(tokens[1])
			} else {
				values, err = db.(*fs.JKV_DB).HKEYS(tokens[1])
			}
			if err != nil {
				fmt.Println("(nil)")
			} else {
				for i, v := range values {
					fmt.Printf("%d) \"%s\"\n", i+1, v)
				}
			}
		} else {
			fmt.Println("(error) ERR wrong number of arguments for 'hkeys' command")
		}
	case "HEXISTS":
		if len(tokens) == 3 {
			var exists bool
			if r, ok := db.(*redis.JKV_DB); ok {
				exists = r.HEXISTS(tokens[1], tokens[2])
			} else {
				exists = db.(*fs.JKV_DB).HEXISTS(tokens[1], tokens[2])
			}
			if exists {
				fmt.Println("(integer) 1")
			} else {
				fmt.Println("(integer) 0")
			}
		} else {
			fmt.Println("(error) ERR wrong number of arguments for 'exists' command")
		}
	case "GET":
		if len(tokens) == 2 {
			if r, ok := db.(*redis.JKV_DB); ok {
				value, err = r.GET(tokens[1])
			} else {
				value, err = db.(*fs.JKV_DB).GET(tokens[1])
			}
			if err != nil {
				fmt.Println("(nil)")
			} else {
				fmt.Printf("\"%s\"\n", value)
			}
		} else {
			fmt.Println("(nil)")
		}
	case "SET":
		if len(tokens) == 3 {
			if r, ok := db.(*redis.JKV_DB); ok {
				err = r.SET(tokens[1], tokens[2])
			} else {
				err = db.(*fs.JKV_DB).SET(tokens[1], tokens[2])
			}
			if err != nil {
				fmt.Println("(nil)")
			} else {
				fmt.Println("OK")
			}
		} else {
			fmt.Println("(error) ERR wrong number of arguments for 'set' command")
		}
	case "DEL":
		if len(tokens) == 2 {
			if r, ok := db.(*redis.JKV_DB); ok {
				err = r.DEL(tokens[1])
			} else {
				err = db.(*fs.JKV_DB).DEL(tokens[1])
			}
			if err != nil {
				fmt.Println("(nil)")
			} else {
				fmt.Printf("\"%s\"\n", value)
			}
		} else {
			fmt.Println("(nil)")
		}
	case "KEYS":
		if len(tokens) == 2 {
			if r, ok := db.(*redis.JKV_DB); ok {
				values, err = r.KEYS(tokens[1])
			} else {
				values, err = db.(*fs.JKV_DB).KEYS(tokens[1])
			}
			if err != nil {
				fmt.Println("(nil)")
			} else {
				for i, v := range values {
					fmt.Printf("%d) \"%s\"\n", i+1, v)
				}
			}
		} else {
			fmt.Println("(error) ERR wrong number of arguments for 'keys' command")
		}
	case "EXISTS":
		if len(tokens) == 2 {
			var exists bool
			if r, ok := db.(*redis.JKV_DB); ok {
				exists = r.EXISTS(tokens[1])
			} else {
				exists = db.(*fs.JKV_DB).EXISTS(tokens[1])
			}
			if exists {
				fmt.Println("(integer) 1")
			} else {
				fmt.Println("(integer) 0")
			}
		} else {
			fmt.Println("(error) ERR wrong number of arguments for 'exists' command")
		}
	default:
		fmt.Printf("(error) ERR unknown command '%s', with args beginning with:\n", tokens[0])
	}
}
