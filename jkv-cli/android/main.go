package main

import (
	"context"
	"fmt"
	"image/color"
	"os"

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
		f.Close()
	}()

	w.ShowAndRun()
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
