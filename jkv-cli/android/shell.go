package main

import (
	"context"
	"fmt"
	"image/color"
	"os"
	"os/exec"
	"runtime"
	"strings"

	_ "embed"

	"github.com/panduit-joeb/jkv"
	"github.com/panduit-joeb/jkv/pkg"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func genShell(rdb jkv.Client) fyne.CanvasObject {
	winWidth := float32(1024)
	objSize := fyne.NewSquareSize(winWidth / 4)
	bg := canvas.NewRectangle(color.White)
	bg.Resize(objSize)

	ctx := context.Background()
	initDb(ctx, rdb)
	msg := widget.NewLabel("")
	output := widget.NewMultiLineEntry()
	output.SetMinRowsVisible(20)
	input := widget.NewEntry()
	input.SetPlaceHolder("sh, pwd, ls, mount, cat, flushdb, set, get, del, keys, hset, hget, hdel, hkeys...")
	input.OnSubmitted = func(string) { runCmd(ctx, rdb, input, output, msg) }
	content := container.NewVBox(
		input,
		container.NewVBox(
			container.NewHBox(
				widget.NewButton("Go", func() { runCmd(ctx, rdb, input, output, msg) }),
				widget.NewButton("initdb", func() { initDb(ctx, rdb) }),
				widget.NewButton("flushdb", func() { doCmd(ctx, rdb, "flushdb", input, output, msg) }),
				widget.NewButton("set", func() { doCmd(ctx, rdb, "set", input, output, msg) }),
				widget.NewButton("get", func() { doCmd(ctx, rdb, "get", input, output, msg) }),
				widget.NewButton("del", func() { doCmd(ctx, rdb, "del", input, output, msg) }),
				widget.NewButton("keys", func() { doCmd(ctx, rdb, "keys", input, output, msg) }),
				widget.NewButton("hset", func() { doCmd(ctx, rdb, "hset", input, output, msg) }),
				widget.NewButton("hget", func() { doCmd(ctx, rdb, "hget", input, output, msg) }),
				widget.NewButton("hdel", func() { doCmd(ctx, rdb, "hdel", input, output, msg) }),
				widget.NewButton("hkeys", func() { doCmd(ctx, rdb, "hkeys", input, output, msg) }),
				widget.NewButton("on/off", func() { doCmd(ctx, rdb, "toggle", input, output, msg) }),
				widget.NewButton("getOffline", func() { doCmd(ctx, rdb, "getOffline", input, output, msg) }),
			),
			widget.NewButton("Bye", func() { os.Exit(0) }),
		),
		msg,
		output,
	)
	return content
}

func mkdir(dir string) {
	uri := storage.NewFileURI(dir)
	if err := storage.CreateListable(uri); err == nil {
		fmt.Println("mkdir", uri, "ok")
	} else {
		fmt.Println("mkdir", uri, "failed, err", err.Error())
	}
}

func doCmd(ctx context.Context, rdb jkv.Client, cmd string, input, output *widget.Entry, msg *widget.Label) {
	tokens := strings.Fields(input.Text)
	fmt.Printf("doCmd() called with cmd=%s, tokens=%s\n", cmd, strings.Join(tokens, " "))
	if cmd == "getOffline" {
		hashes := os.TempDir() + "/hashes/"
		mkdir(hashes)
		hash := hashes + "UserSelected"
		mkdir(hash)
		key := hash + "/Offline"
		uri := storage.NewFileURI(key)
		fmt.Println(uri)
		w, err := storage.Writer(uri)
		if err == nil {
			_, err = w.Write([]byte("1"))
			if err == nil {
				fmt.Println("Write good")
				w.Close()
				r, err := storage.Reader(uri)
				if err == nil {
					buf := make([]byte, 128)
					_, err := r.Read(buf)
					if err == nil {
						fmt.Printf("Got: %s\n", string(buf))
						r.Close()
					} else {
						fmt.Println("Read bad", err.Error())
					}
				} else {
					fmt.Println("Reader bad", err.Error())
				}
			} else {
				fmt.Println("Write bad", err.Error())
			}
		} else {
			fmt.Println("Writer bad", err.Error())
		}
	} else if len(tokens) == 0 {
		switch cmd {
		case "keys":
			reportCmd(rdb.Keys(ctx, "*"), output, msg)
		case "flushdb":
			reportCmd(rdb.FlushDB(ctx), output, msg)
		case "toggle":
			reportCmd(toggleOnline(ctx, rdb), output, msg)
		}
	} else if len(tokens) > 0 {
		switch cmd {
		case "set":
			reportCmd(rdb.Set(ctx, tokens[0], tokens[1], 0), output, msg)
		case "get":
			reportCmd(rdb.Get(ctx, tokens[0]), output, msg)
		case "del":
			reportCmd(rdb.Del(ctx, tokens[0]), output, msg)
		case "hset":
			reportCmd(rdb.HSet(ctx, tokens[0], tokens[1:]...), output, msg)
		case "hget":
			reportCmd(rdb.HGet(ctx, tokens[0], tokens[1]), output, msg)
		case "hdel":
			reportCmd(rdb.HDel(ctx, tokens[0], tokens[1:]...), output, msg)
		case "hkeys":
			reportCmd(rdb.HKeys(ctx, tokens[0]), output, msg)
		}
	}
}

func reportCmd(rc interface{}, output *widget.Entry, msg *widget.Label) {
	fmt.Printf("reportCmd() called with %#v\n", rc)
	switch t := rc.(type) {
	case *jkv.StringCmd:
		msg.Text = e(t.Err())
		output.Text = t.Val()
	case *jkv.IntCmd:
		msg.Text = e(t.Err())
		output.Text = fmt.Sprintf("%d", t.Val())
	case *jkv.StatusCmd:
		msg.Text = e(t.Err())
		output.Text = t.Val()
	case *jkv.BoolCmd:
		msg.Text = e(t.Err())
		output.Text = fmt.Sprintf("%t", t.Val())
	case *jkv.StringSliceCmd:
		msg.Text = e(t.Err())
		output.Text = strings.Join(t.Val(), "\n")
	default:
		msg.Text = fmt.Sprintf("%#v - not supported", rc)
	}
	msg.Refresh()
	output.Refresh()
}

func toggleOnline(ctx context.Context, rdb jkv.Client) *jkv.StatusCmd {
	rec := rdb.HGet(ctx, "UserSelected", "Offline")
	if rec.Err() != nil {
		return jkv.NewStatusCmd("(nil)", rec.Err())
	}
	if pkg.StringToBool(rec.Val()) {
		rc := rdb.HSet(ctx, "UserSelected", "Offline", pkg.BoolToString(false))
		if rc.Err() != nil {
			return jkv.NewStatusCmd("(nil)", rec.Err())
		}
		return jkv.NewStatusCmd("OK", nil)
	}
	rc := rdb.HSet(ctx, "UserSelected", "Offline", pkg.BoolToString(true))
	if rc.Err() != nil {
		return jkv.NewStatusCmd("(nil)", rec.Err())
	}
	return jkv.NewStatusCmd("OK", nil)
}

func e(err error) string {
	if err == nil {
		return "(nil)"
	}
	return err.Error()
}

func pathToCmd(cmd string) string {
	switch cmd {
	case "ls":
		fallthrough
	case "cat":
		fallthrough
	case "sh":
		if runtime.GOOS == "android" {
			return "/system/bin/" + cmd
		}
		return "/bin/" + cmd
	case "mount":
		if runtime.GOOS == "android" {
			return "/system/bin/" + cmd
		}
		return "/sbin/" + cmd
	}
	return cmd
}

func runCmd(ctx context.Context, rdb jkv.Client, input, output *widget.Entry, msg *widget.Label) {
	rdb.Open()
	msg.SetText(fmt.Sprintf("executing \"%s\"", input.Text))
	tokens := strings.Fields(input.Text)
	if len(tokens) > 0 {
		cmd := strings.ToLower(tokens[0])
		switch cmd {
		case "sh":
			cmd := exec.Command(pathToCmd(cmd), "-c", strings.Join(tokens[1:], " "))
			sh, err := cmd.CombinedOutput()
			if err == nil {
				output.SetText(string(sh))
			} else {
				output.SetText(fmt.Sprintf("%s\n%s", e(err), ""))
			}
		case "pwd":
			wd, err := os.Getwd()
			msg.Text = e(err)
			output.Text = wd
			msg.Refresh()
			output.Refresh()
		case "ls":
			if len(tokens) >= 2 {
				cmd := exec.Command(pathToCmd(cmd), tokens[1:]...)
				mount, err := cmd.CombinedOutput()
				if err == nil {
					output.SetText(string(mount))
				} else {
					output.SetText(fmt.Sprintf("%s\n%s", e(err), ""))
				}
			}
		case "mount":
			cmd := exec.Command(pathToCmd(cmd))
			mount, err := cmd.CombinedOutput()
			if err == nil {
				output.SetText(string(mount))
			} else {
				output.SetText(fmt.Sprintf("%s\n%s", e(err), ""))
			}
		case "cat":
			if len(tokens) >= 2 {
				cmd := exec.Command(pathToCmd(cmd), tokens[1])
				cat, err := cmd.CombinedOutput()
				if err == nil {
					output.SetText(string(cat))
				} else {
					output.SetText(fmt.Sprintf("%s\n%s", e(err), ""))
				}
			}
		default:
			doCmd(ctx, rdb, cmd, input, output, msg)
		}
	}
}

func logInt(msg string, rc *jkv.IntCmd) {
	if rc.Err() != nil {
		fmt.Println(msg, "failed, err:", rc.Err().Error())
	} else {
		fmt.Println(msg, "worked, val:", rc.Val())
	}
}

func TestStorage() {
	var (
		user_dir, db_dir string
	)
	user_dir = os.TempDir()
	for _, db_dir = range []string{user_dir + "/jkv_db/scalars", user_dir + "/jkv_db/hashes"} {
		if err := os.MkdirAll(db_dir, 0775); err == nil {
			fmt.Printf("MkdirAll(\"%s\") worked\n", db_dir)
		} else {
			fmt.Printf("MkdirAll(\"%s\") failed, err: %s\n", db_dir, err.Error())
		}
	}
	dir := storage.NewFileURI(db_dir)
	if err := storage.CreateListable(dir); err != nil {
		fmt.Println("creating directory", dir, "failed", err.Error())
	} else {
		fmt.Println("creating directory", dir, "worked")
	}
	file_name := db_dir + "/file"
	file := storage.NewFileURI(file_name)
	writer, err := storage.Writer(file)
	if err != nil {
		fmt.Println("creating writer for", file, "failed", err.Error())
	} else {
		fmt.Println("creating writer for", file, "worked")
	}
	var n int
	if n, err = writer.Write([]byte("hello world")); err != nil {
		fmt.Println("write failed", err.Error())
	} else {
		fmt.Println("wrote", n, "bytes to", file)
	}
	writer.Close()
}

func initDb(ctx context.Context, rdb jkv.Client) {
	// dup for testing
	logInt("UserSelected...", rdb.HSet(ctx, "UserSelected", "Offline", string(JKVDB_hashes_UserSelected_Offline)))

	// logInt("Networks...", rdb.HSet(ctx, "Networks", "default", string(JKVDB_hashes_Networks_default)))
	// logInt("Networks...", rdb.HSet(ctx, "Networks", "syscfg_ips", string(JKVDB_hashes_Networks_syscfg_ips)))
	// logInt("Networks...", rdb.HSet(ctx, "Networks", "static", string(JKVDB_hashes_Networks_static)))
	// logInt("Networks...", rdb.HSet(ctx, "Networks", "dhcp", string(JKVDB_hashes_Networks_dhcp)))
	// logInt("SuperScreens...", rdb.HSet(ctx, "SuperScreens", "test_mode", string(JKVDB_hashes_SuperScreens_test_mode)))
	// logInt("SuperScreens...", rdb.HSet(ctx, "SuperScreens", "qr", string(JKVDB_hashes_SuperScreens_qr)))
	// logInt("SuperScreens...", rdb.HSet(ctx, "SuperScreens", "system_config", string(JKVDB_hashes_SuperScreens_system_config)))
	// logInt("UserSelected...", rdb.HSet(ctx, "UserSelected", "Cloud", string(JKVDB_hashes_UserSelected_Cloud)))
	// logInt("UserSelected...", rdb.HSet(ctx, "UserSelected", "ScreenInverted", string(JKVDB_hashes_UserSelected_ScreenInverted)))
	// logInt("UserSelected...", rdb.HSet(ctx, "UserSelected", "ScreenKey", string(JKVDB_hashes_UserSelected_ScreenKey)))
	// logInt("UserSelected...", rdb.HSet(ctx, "UserSelected", "Offline", string(JKVDB_hashes_UserSelected_Offline)))
	// logInt("UserSelected...", rdb.HSet(ctx, "UserSelected", "ScreenCollection", string(JKVDB_hashes_UserSelected_ScreenCollection)))
	// logInt("UserSelected...", rdb.HSet(ctx, "UserSelected", "InternetEnabled", string(JKVDB_hashes_UserSelected_InternetEnabled)))
	// logInt("UserSelected...", rdb.HSet(ctx, "UserSelected", "Internet", string(JKVDB_hashes_UserSelected_Internet)))
	// logInt("UserSelected...", rdb.HSet(ctx, "UserSelected", "sleep", string(JKVDB_hashes_UserScreens_sleep)))
	// logInt("UserSelected...", rdb.HSet(ctx, "UserSelected", "passcode", string(JKVDB_hashes_UserScreens_passcode)))
	// logInt("SquareImages...", rdb.HSet(ctx, "SquareImages", "qr", string(JKVDB_hashes_SquareImages_qr)))
}
