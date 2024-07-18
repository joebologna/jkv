BINARY = jkv

.PHONY:

all: clean $(BINARY)-cli
	@echo all done

$(BINARY)-cli: main.go
	go build .

# git clone fyne; git checkout 6145be22d
apk:
	cd android; fyne package -os android -appID com.atlona.touchos.$(BINARY) -icon Icon.png

install-apk:
	adb devices
	# only needed for the emulator
	# adb forward tcp:9009 tcp:9009
	# remove with adb forward --remove tcp:9009
	-adb uninstall com.atlona.touchos.$(BINARY)
	cd android; adb install $(BINARY).apk
	adb shell monkey -p com.atlona.touchos.$(BINARY) -c android.intent.category.LAUNCHER 1

$(BINARY)-gui: getdb
	cd android; go build -o $@ .

getdb:
	cd android; unzip $(HOME)/TouchOS/tpgui/asset-mgmt/jkv_db.zip -d android

logcat: .PHONY
	adb logcat Fyne '*:S'

clean:
	rm -f android/$(BINARY).apk android/$(BINARY)-gui $(BINARY)-cli
