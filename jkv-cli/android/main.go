package main

import (
	"fmt"
	"image/color"
	"os"

	_ "embed"

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
		var (
			user_dir, db_dir, file_name string
			dir, file                   fyne.URI
			writer                      fyne.URIWriteCloser
			err                         error
		)
		user_dir = os.TempDir()
		db_dir = user_dir + "/jkv_db"
		dir = storage.NewFileURI(db_dir)
		if err = storage.CreateListable(dir); err != nil {
			fmt.Println("creating directory", dir, "failed", err.Error())
		} else {
			fmt.Println("creating directory", dir, "worked")
		}
		file_name = db_dir + "/file"
		file = storage.NewFileURI(file_name)
		writer, err = storage.Writer(file)
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
	}()

	w.ShowAndRun()
}
