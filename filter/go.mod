module github.com/dayvillefire/pocsag-monitor/filter

go 1.25

replace (
	github.com/dayvillefire/pocsag-monitor => ../
	github.com/dayvillefire/pocsag-monitor/obj => ../obj
)

require (
	github.com/dayvillefire/pocsag-monitor/obj v0.0.0-20251208221133-403e471ee4eb
	github.com/jbuchbinder/shims v0.0.0-20251029164657-6c80f5d6bc01
)
