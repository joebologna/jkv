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
	os.Setenv("FYNE_THEME", "dark")
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
		c.SetContent(genShell(f))
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

func initdb(ctx context.Context, f jkv.Client) {
	f.HSet(ctx, "UserSelected", "Offline", string(JKVDB_hashes_UserSelected_Offline))
	// f.HSet(ctx, "Networks", "default", string(JKVDB_hashes_Networks_default))
	// f.HSet(ctx, "Networks", "syscfg_ips", string(JKVDB_hashes_Networks_syscfg_ips))
	// f.HSet(ctx, "Networks", "static", string(JKVDB_hashes_Networks_static))
	// f.HSet(ctx, "Networks", "dhcp", string(JKVDB_hashes_Networks_dhcp))
	// f.HSet(ctx, "SuperScreens", "test_mode", string(JKVDB_hashes_SuperScreens_test_mode))
	// f.HSet(ctx, "SuperScreens", "qr", string(JKVDB_hashes_SuperScreens_qr))
	// f.HSet(ctx, "SuperScreens", "system_config", string(JKVDB_hashes_SuperScreens_system_config))
	// f.HSet(ctx, "UserSelected", "Cloud", string(JKVDB_hashes_UserSelected_Cloud))
	// f.HSet(ctx, "UserSelected", "ScreenInverted", string(JKVDB_hashes_UserSelected_ScreenInverted))
	// f.HSet(ctx, "UserSelected", "ScreenKey", string(JKVDB_hashes_UserSelected_ScreenKey))
	// f.HSet(ctx, "UserSelected", "ScreenCollection", string(JKVDB_hashes_UserSelected_ScreenCollection))
	// f.HSet(ctx, "UserSelected", "InternetEnabled", string(JKVDB_hashes_UserSelected_InternetEnabled))
	// f.HSet(ctx, "UserSelected", "Internet", string(JKVDB_hashes_UserSelected_Internet))
	// f.HSet(ctx, "UserSelected", "sleep", string(JKVDB_hashes_UserScreens_sleep))
	// f.HSet(ctx, "UserSelected", "passcode", string(JKVDB_hashes_UserScreens_passcode))
	// f.HSet(ctx, "SquareImages", "qr", string(JKVDB_hashes_SquareImages_qr))
}
