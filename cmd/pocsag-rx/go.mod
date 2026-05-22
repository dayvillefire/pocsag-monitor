module github.com/dayvillefire/pocsag-monitor/cmd/pocsag-rx

go 1.25.3

require (
	github.com/dayvillefire/pocsag-monitor/pocsag v0.0.0-00010101000000-000000000000
	github.com/dayvillefire/pocsag-monitor/sdr v0.0.0-00010101000000-000000000000
)

replace (
	github.com/dayvillefire/pocsag-monitor/pocsag => ../../pocsag
	github.com/dayvillefire/pocsag-monitor/sdr => ../../sdr
)
