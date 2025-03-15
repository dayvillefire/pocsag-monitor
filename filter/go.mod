module github.com/dayvillefire/pocsag-monitor/filter

go 1.23

replace (
	github.com/dayvillefire/pocsag-monitor => ../
	github.com/dayvillefire/pocsag-monitor/obj => ../obj
)

require github.com/dayvillefire/pocsag-monitor/obj v0.0.0-20241231155749-74823eee0d9a

require (
	github.com/dlclark/regexp2 v1.11.5 // indirect
	github.com/jessevdk/go-flags v1.4.0 // indirect
	github.com/mattn/go-shellwords v1.0.11 // indirect
	github.com/zyedidia/gpeg v0.0.0-20210126031808-c49c92955001 // indirect
	github.com/zyedidia/sregx v0.3.0 // indirect
)
