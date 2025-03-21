all: binary

clean:
	go clean -v

binary: clean
	GOARM=5 GOARCH=arm go build -v

test-config:
	go test -run Test_Config -v

copy: binary test-config
	rsync -rvutpP pocsag-monitor config.yaml dynamic.yaml jbuchbinder@manage:
