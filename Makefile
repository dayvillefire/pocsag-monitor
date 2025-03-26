VERSION := $(shell date +%Y%m%d.%H%M)

all: binary

clean:
	go clean -v

binary: clean
	GOARM=5 GOARCH=arm go build -v -ldflags "-X main.Version=${VERSION}"

test-config:
	go test -run Test_LoadConfig -v

copy: binary test-config
	rsync -rvutpP pocsag-monitor config.yaml dynamic.yaml jbuchbinder@manage:
