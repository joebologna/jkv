package main

import (
	"context"
	"fmt"
	"image/color"
	"io/fs"
	"os"
	"os/exec"
	"strings"

	_ "embed"

	"github.com/panduit-joeb/jkv"
	"github.com/panduit-joeb/jkv/store/apk"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

//go:embed jkv_db/hashes/Networks/default
var JKVDB_hashes_Networks_default []byte

//go:embed jkv_db/hashes/Networks/syscfg_ips
var JKVDB_hashes_Networks_syscfg_ips []byte

//go:embed jkv_db/hashes/Networks/static
var JKVDB_hashes_Networks_static []byte

//go:embed jkv_db/hashes/Networks/dhcp
var JKVDB_hashes_Networks_dhcp []byte

//go:embed jkv_db/hashes/SuperScreens/test_mode
var JKVDB_hashes_SuperScreens_test_mode []byte

//go:embed jkv_db/hashes/SuperScreens/qr
var JKVDB_hashes_SuperScreens_qr []byte

//go:embed jkv_db/hashes/SuperScreens/system_config
var JKVDB_hashes_SuperScreens_system_config []byte

//go:embed jkv_db/hashes/UserSelected/Cloud
var JKVDB_hashes_UserSelected_Cloud []byte

//go:embed jkv_db/hashes/UserSelected/ScreenInverted
var JKVDB_hashes_UserSelected_ScreenInverted []byte

//go:embed jkv_db/hashes/UserSelected/ScreenKey
var JKVDB_hashes_UserSelected_ScreenKey []byte

//go:embed jkv_db/hashes/UserSelected/Offline
var JKVDB_hashes_UserSelected_Offline []byte

//go:embed jkv_db/hashes/UserSelected/ScreenCollection
var JKVDB_hashes_UserSelected_ScreenCollection []byte

//go:embed jkv_db/hashes/UserSelected/InternetEnabled
var JKVDB_hashes_UserSelected_InternetEnabled []byte

//go:embed jkv_db/hashes/UserSelected/Internet
var JKVDB_hashes_UserSelected_Internet []byte

//go:embed jkv_db/hashes/UserScreens/sleep
var JKVDB_hashes_UserScreens_sleep []byte

//go:embed jkv_db/hashes/UserScreens/passcode
var JKVDB_hashes_UserScreens_passcode []byte

//go:embed jkv_db/hashes/SquareImages/qr
var JKVDB_hashes_SquareImages_qr []byte

func main() {
	a := app.NewWithID("com.atlona.touchos.preferences")
	w := a.NewWindow("JKV-CLI")
	winWidth := float32(1024)
	winSize := fyne.NewSquareSize(winWidth)
	objSize := fyne.NewSquareSize(winWidth / 4)
	bg := canvas.NewRectangle(color.White)
	bg.Resize(objSize)
	c := w.Canvas()

	var objs = []fyne.CanvasObject{container.NewStack(bg, widget.NewLabel("Booting..."))}
	c.SetContent(container.NewWithoutLayout(objs...))

	w.Resize(winSize)

	go func() {
		objs := c.Content().(*fyne.Container).Objects
		label := objs[0].(*fyne.Container).Objects[1].(*widget.Label)
		label.SetText("Booted.")
		label.Refresh()
		c.Refresh(c.Content())
		f := apk.NewClient(&apk.Options{Addr: apk.GetDBDir()})
		f.Open()
		ctx := context.Background()
		logStatus("flushdb", f.FlushDB(ctx))
		logStatus("set key1 one", f.Set(ctx, "key1", "one", 0))
		logString("get key1", f.Get(ctx, "key1"))
		logInt("del key1", f.Del(ctx, "key1"))
		initdb(ctx, f)
		f.Close()
		msg := widget.NewLabel("")
		output := widget.NewMultiLineEntry()
		output.SetMinRowsVisible(20)
		input := widget.NewEntry()
		input.SetPlaceHolder("flushdb")
		content := container.NewVBox(input, widget.NewButton("Go", func() {
			f.Open()
			msg.SetText(fmt.Sprintf("executing \"%s\"", input.Text))
			tokens := strings.Fields(input.Text)
			if len(tokens) > 0 {
				switch strings.ToLower(tokens[0]) {
				case "ls":
					lsfmt := func(f fs.DirEntry) string {
						d := "---"
						if f.IsDir() {
							d = "dr-x"
						}
						return fmt.Sprintf("%s---- %s", d, f.Name())
					}
					var a []string
					fmt.Printf("cwd: %s\ntrying to ReadDir(\"%s\")\n", func() string {
						s, err := os.Getwd()
						if err == nil {
							return s
						}
						return err.Error()
					}(), tokens[1])
					files, err := os.ReadDir(tokens[1])
					if err == nil {
						for _, f := range files {
							a = append(a, lsfmt(f))
						}
						output.SetText(fmt.Sprintf("%s\n%s", e(err), strings.Join(a, "\n")))
					} else {
						output.SetText(fmt.Sprintf("%s\n%s", e(err), ""))
					}
				case "mount":
					cmd := exec.Command("/system/bin/mount")
					mount, err := cmd.CombinedOutput()
					if err == nil {
						output.SetText(string(mount))
					} else {
						output.SetText(fmt.Sprintf("%s\n%s", e(err), ""))
					}
				case "cat":
					cmd := exec.Command("/system/bin/cat", tokens[1])
					cat, err := cmd.CombinedOutput()
					if err == nil {
						output.SetText(string(cat))
					} else {
						output.SetText(fmt.Sprintf("%s\n%s", e(err), ""))
					}
				case "flushdb":
					rc := f.FlushDB(ctx)
					output.SetText(fmt.Sprintf("%s\n%s", e(rc.Err()), rc.Val()))
				case "set":
					if len(tokens) == 3 {
						rc := f.Set(ctx, tokens[1], tokens[2], 0)
						output.SetText(fmt.Sprintf("%s\n%s", e(rc.Err()), rc.Val()))
					}
				case "get":
					if len(tokens) == 2 {
						rc := f.Get(ctx, tokens[1])
						output.SetText(fmt.Sprintf("%s\n%s", e(rc.Err()), rc.Val()))
					}
				case "del":
					if len(tokens) == 2 {
						rc := f.Del(ctx, tokens[1])
						output.SetText(fmt.Sprintf("%s\n%d", e(rc.Err()), rc.Val()))
					}
				case "keys":
					if len(tokens) == 2 {
						rc := f.Keys(ctx, tokens[1])
						output.SetText(fmt.Sprintf("%s\n%s", e(rc.Err()), strings.Join(rc.Val(), "\n")))
					}
				case "hset":
					if len(tokens) == 4 {
						rc := f.HSet(ctx, tokens[1], tokens[2], tokens[3])
						output.SetText(fmt.Sprintf("%s\n%d", e(rc.Err()), rc.Val()))
					}
				case "hget":
					if len(tokens) == 3 {
						rc := f.HGet(ctx, tokens[1], tokens[2])
						output.SetText(fmt.Sprintf("%s\n%s", e(rc.Err()), rc.Val()))
					}
				case "hdel":
					if len(tokens) == 3 {
						rc := f.HDel(ctx, tokens[1], tokens[2])
						output.SetText(fmt.Sprintf("%s\n%d", e(rc.Err()), rc.Val()))
					}
				case "hkeys":
					if len(tokens) == 2 {
						rc := f.HKeys(ctx, tokens[1])
						output.SetText(fmt.Sprintf("%s\n%s", e(rc.Err()), strings.Join(rc.Val(), "\n")))
					}
				default:
					output.SetText("Invalid Command")
				}
			}
		}), msg, output)
		c.SetContent(content)
	}()

	w.ShowAndRun()
}

func e(err error) string {
	if err == nil {
		return "(nil)"
	}
	return err.Error()
}

func logStatus(msg string, rc *jkv.StatusCmd) {
	if rc.Err() != nil {
		fmt.Println(msg, "failed, err:", rc.Err().Error())
	} else {
		fmt.Println(msg, "worked")
	}
}

func logString(msg string, rc *jkv.StringCmd) {
	if rc.Err() != nil {
		fmt.Println(msg, "failed, err:", rc.Err().Error())
	} else {
		fmt.Println(msg, "worked, val:", rc.Val())
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
		if err := apk.MkdirAll(db_dir, 0775); err == nil {
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

func initdb(ctx context.Context, f *apk.Client) {
	f.HSet(ctx, "Networks", "default", string(JKVDB_hashes_Networks_default))
	f.HSet(ctx, "Networks", "syscfg_ips", string(JKVDB_hashes_Networks_syscfg_ips))
	f.HSet(ctx, "Networks", "static", string(JKVDB_hashes_Networks_static))
	f.HSet(ctx, "Networks", "dhcp", string(JKVDB_hashes_Networks_dhcp))
	f.HSet(ctx, "SuperScreens", "test_mode", string(JKVDB_hashes_SuperScreens_test_mode))
	f.HSet(ctx, "SuperScreens", "qr", string(JKVDB_hashes_SuperScreens_qr))
	f.HSet(ctx, "SuperScreens", "system_config", string(JKVDB_hashes_SuperScreens_system_config))
	f.HSet(ctx, "UserSelected", "Cloud", string(JKVDB_hashes_UserSelected_Cloud))
	f.HSet(ctx, "UserSelected", "ScreenInverted", string(JKVDB_hashes_UserSelected_ScreenInverted))
	f.HSet(ctx, "UserSelected", "ScreenKey", string(JKVDB_hashes_UserSelected_ScreenKey))
	f.HSet(ctx, "UserSelected", "Offline", string(JKVDB_hashes_UserSelected_Offline))
	f.HSet(ctx, "UserSelected", "ScreenCollection", string(JKVDB_hashes_UserSelected_ScreenCollection))
	f.HSet(ctx, "UserSelected", "InternetEnabled", string(JKVDB_hashes_UserSelected_InternetEnabled))
	f.HSet(ctx, "UserSelected", "Internet", string(JKVDB_hashes_UserSelected_Internet))
	f.HSet(ctx, "UserSelected", "sleep", string(JKVDB_hashes_UserScreens_sleep))
	f.HSet(ctx, "UserSelected", "passcode", string(JKVDB_hashes_UserScreens_passcode))
	f.HSet(ctx, "SquareImages", "qr", string(JKVDB_hashes_SquareImages_qr))
}
