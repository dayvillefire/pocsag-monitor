module github.com/dayvillefire/pocsag-monitor/filter

go 1.23

replace (
	github.com/dayvillefire/pocsag-monitor => ../
	github.com/dayvillefire/pocsag-monitor/obj => ../obj
)

require (
	github.com/dayvillefire/pocsag-monitor/obj v0.0.0-20250411125416-eecc4a8a5f6c
	github.com/jbuchbinder/shims v0.0.0-20250315180801-ea13cafaf717
)
