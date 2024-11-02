module github.com/dayvillefire/pocsag-monitor/filter

go 1.23

replace (
	github.com/dayvillefire/pocsag-monitor => ../
	github.com/dayvillefire/pocsag-monitor/obj => ../obj
)

require github.com/dayvillefire/pocsag-monitor/obj v0.0.0-20241011175516-7f67fc8e798e
