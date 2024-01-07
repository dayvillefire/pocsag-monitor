all: binary

clean:
	go clean -v

binary: clean
	GOARM=5 GOARCH=arm go build -v

copy: binary
	rsync -rvutpP pocsag-monitor config.yaml -e 'ssh -p2222' jbuchbinder@localhost:
	rsync -rvutpP pocsag-monitor config.yaml pi@pi-pager:
