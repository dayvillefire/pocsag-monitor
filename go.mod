module github.com/dayvillefire/pocsag-monitor

go 1.23.0

toolchain go1.24.3

replace (
	github.com/dayvillefire/groupme => ../groupme
	github.com/dayvillefire/pocsag-monitor/config => ./config
	github.com/dayvillefire/pocsag-monitor/controllers => ./controllers
	github.com/dayvillefire/pocsag-monitor/obj => ./obj
	github.com/dayvillefire/pocsag-monitor/output => ./output
)

require (
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf
	github.com/dayvillefire/pocsag-monitor/config v0.0.0-20250609174951-ddbc9bda5e05
	github.com/dayvillefire/pocsag-monitor/obj v0.0.0-20250609174951-ddbc9bda5e05
	github.com/dayvillefire/pocsag-monitor/output v0.0.0-20250609174951-ddbc9bda5e05
	github.com/genjidb/genji v0.15.3
	github.com/gin-contrib/gzip v1.2.3
	github.com/gin-gonic/gin v1.10.1
	github.com/jbuchbinder/shims v0.0.0-20250315180801-ea13cafaf717
)

require (
	github.com/DataDog/zstd v1.5.5 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/bwmarrin/discordgo v0.29.0 // indirect
	github.com/bytedance/sonic v1.14.0 // indirect
	github.com/bytedance/sonic/loader v0.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.5 // indirect
	github.com/cockroachdb/errors v1.11.3 // indirect
	github.com/cockroachdb/fifo v0.0.0-20240616162244-4768e80dfb9a // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/pebble v1.1.1 // indirect
	github.com/cockroachdb/redact v1.1.5 // indirect
	github.com/cockroachdb/tokenbucket v0.0.0-20230807174530-cc333fc44b06 // indirect
	github.com/creasty/defaults v1.8.0 // indirect
	github.com/dayvillefire/groupme v0.5.1 // indirect
	github.com/eclipse/paho.mqtt.golang v1.5.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.9 // indirect
	github.com/getsentry/sentry-go v0.28.1 // indirect
	github.com/gin-contrib/sse v1.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.27.0 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.19.1 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.3.0 // indirect
	golang.org/x/arch v0.19.0 // indirect
	golang.org/x/crypto v0.40.0 // indirect
	golang.org/x/exp v0.0.0-20240613232115-7f521ea00fb8 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
