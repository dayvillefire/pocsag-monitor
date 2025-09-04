module github.com/dayvillefire/pocsag-monitor/filter

go 1.23

replace (
	github.com/dayvillefire/pocsag-monitor => ../
	github.com/dayvillefire/pocsag-monitor/obj => ../obj
)

require (
	github.com/dayvillefire/pocsag-monitor/obj v0.0.0-20250828124955-f581123ea150
	github.com/jbuchbinder/shims v0.0.0-20250818154854-22c0ac83b788
)
