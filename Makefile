VERSION := $(shell date +%Y%m%d.%H%M)

all: binary

clean:
	go clean -v
	rm -f pocsag-monitor.bin

binary-arm: clean
	GOARM=5 GOARCH=arm go build -v -ldflags "-X main.Version=${VERSION}"

binary: clean
	go build -v -ldflags "-X main.Version=${VERSION}" -o pocsag-monitor.bin

test-config:
	go test -run Test_LoadConfig -v

copy: binary-arm test-config
	rsync -rvutpP pocsag-monitor config.yaml dynamic.yaml jbuchbinder@manage:
