package main

import (
	"fmt"
	"os"
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func fstest() fyne.CanvasObject {
	cmd := binding.NewString()
	hashDir := os.TempDir() + "/jkv_db/hashes/UserSelected"
	cmd.Set("cat " + hashDir + "/Offline")
	cmdEntry := widget.NewEntryWithData(cmd)
	cmdEntry.OnSubmitted = func(s string) {
		j := exec.Command("/system/bin/sh", "-c", s)
		if out, err := j.CombinedOutput(); err == nil {
			fmt.Println(s, "returns:", "\n"+string(out))
		} else {
			fmt.Println(s, "failed:", err.Error())
		}
	}
	testBtn := widget.NewButton("Test", func() {
		// tmpDir := strings.Replace(os.TempDir(), "/cache", "", -1)
		// os.MkdirAll(tmpDir+"/jkv_db/hashes/UserSelected", 0775)
		// os.MkdirAll(tmpDir+"/jkv_db/scalars", 0775)
		// buf := make([]byte, 128)
		// os.WriteFile(tmpDir+"/jkv_db/hashes/UserSelected/Offline", buf, 0664)
		// buf, _ = os.ReadFile(tmpDir + "/jkv_db/hashes/UserSelected/Offline")
		hashDir := os.TempDir() + "/jkv_db/hashes/UserSelected"
		os.MkdirAll(hashDir, 0777)
		buf := []byte("hello")
		os.WriteFile(hashDir+"/Offline", buf, 0666)
		buf, err := os.ReadFile(hashDir + "/Offline")
		fmt.Printf("buf = %s, err = %#v\n", string(buf), err)
	})
	return container.NewVBox(cmdEntry, testBtn)
}

/*
/mnt/installer
/mnt/androidwritable
/metadata
/vendor
/product
/system_ext
/config
/data
/data/user/0
/storage
/storage/emulated
/mnt/installer/0/emulated
/mnt/androidwritable/0/emulated
/mnt/user/0/emulated
/mnt/pass_through/0/emulated
/mnt/androidwritable/0/emulated/0/Android/data
/mnt/installer/0/emulated/0/Android/data
/storage/emulated/0/Android/data
/mnt/user/0/emulated/0/Android/data
/mnt/androidwritable/0/emulated/0/Android/obb
/mnt/installer/0/emulated/0/Android/obb
/storage/emulated/0/Android/obb
/mnt/user/0/emulated/0/Android/obb
/mnt/media_rw/07FC-2819
/storage/07FC-2819
/mnt/installer/0/07FC-2819
/mnt/androidwritable/0/07FC-2819
/mnt/user/0/07FC-2819
/mnt/pass_through/0/07FC-2819
*/
