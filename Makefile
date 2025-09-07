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

deploy: copy
	@echo "Stopping pocsag-monitor ... "
	@ssh jbuchbinder@manage sudo systemctl stop pocsag-monitor
	@echo "Waiting for three seconds ... "
	@sleep 3
	@echo "Starting pocsag-monitor ... "
	@ssh jbuchbinder@manage sudo systemctl start pocsag-monitor
