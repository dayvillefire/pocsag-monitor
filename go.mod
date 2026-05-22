module github.com/dayvillefire/pocsag-monitor

go 1.25.3

replace (
	github.com/dayvillefire/groupme => ../groupme
	github.com/dayvillefire/pocsag-monitor/config => ./config
	github.com/dayvillefire/pocsag-monitor/controllers => ./controllers
	github.com/dayvillefire/pocsag-monitor/obj => ./obj
	github.com/dayvillefire/pocsag-monitor/output => ./output
	github.com/dayvillefire/pocsag-monitor/pocsag => ./pocsag
	github.com/dayvillefire/pocsag-monitor/sdr => ./sdr
	github.com/dayvillefire/pocsag-router => ../pocsag-router
	github.com/dayvillefire/pocsag-router/client => ../pocsag-router/client
	github.com/dayvillefire/pocsag-router/obj => ../pocsag-router/obj
)

require github.com/dayvillefire/pocsag-monitor/config v0.0.0-20251213184454-5b837bf694aa

require (
	github.com/creasty/defaults v1.8.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
