module github.com/dayvillefire/pocsag-monitor/filter

go 1.23

replace (
	github.com/dayvillefire/pocsag-monitor => ../
	github.com/dayvillefire/pocsag-monitor/obj => ../obj
)

require github.com/dayvillefire/pocsag-monitor/obj v0.0.0-20241231155749-74823eee0d9a
