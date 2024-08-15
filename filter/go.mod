module github.com/dayvillefire/pocsag-monitor/filter

go 1.22.0

replace (
	github.com/dayvillefire/pocsag-monitor => ../
	github.com/dayvillefire/pocsag-monitor/obj => ../obj
)

require github.com/dayvillefire/pocsag-monitor/obj v0.0.0-20240718135929-e5e4f1babde7
