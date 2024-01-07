module github.com/dayvillefire/pocsag-monitor

go 1.21

replace (
	github.com/dayvillefire/pocsag-monitor/obj => ./obj
	github.com/dayvillefire/pocsag-monitor/output => ./output
)

require (
	github.com/codegangsta/cli v1.20.0
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf
	github.com/creasty/defaults v1.7.0
	github.com/dayvillefire/pocsag-monitor/obj v0.0.0-20230902204846-5d7fdba5601a
	github.com/dayvillefire/pocsag-monitor/output v0.0.0-20230902204846-5d7fdba5601a
	github.com/dhogborg/go-pocsag v0.0.0-20151112204230-b07839f7ef4b
	github.com/fatih/color v1.15.0
	github.com/genjidb/genji v0.15.1
	github.com/jpoirier/gortlsdr v2.10.0+incompatible
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/DataDog/zstd v1.5.5 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/bwmarrin/discordgo v0.27.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cockroachdb/errors v1.9.1 // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/pebble v0.0.0-20230609131623-fe8461837052 // indirect
	github.com/cockroachdb/redact v1.1.4 // indirect
	github.com/getsentry/sentry-go v0.21.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/klauspost/compress v1.16.5 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.15.1 // indirect
	github.com/prometheus/client_model v0.4.0 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.10.1 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/exp v0.0.0-20230522175609-2e198f4a06a1 // indirect
	golang.org/x/net v0.19.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)
