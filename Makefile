.PHONY:

all: jkv-cli

jkv-cli: .PHONY
	cd jkv-cli; make clean all

set_version:
	printf "package jkv\nconst VERSION = \"%s\"\n" `git log --oneline --decorate|grep tag:|head -1|cut -d: -f2|cut -d, -f1` >version.go

clean:
	find . -type d -name jkv_db | xargs rm -fr
	cd jkv-cli; make clean
